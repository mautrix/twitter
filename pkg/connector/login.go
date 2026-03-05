// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

type TwitterLogin struct {
	User              *bridgev2.User
	Cookies           string
	SecretKey         string
	SigningKey        string
	SigningKeyVersion string
	tc                *TwitterConnector
	isMigration       bool // True if upgrading from main branch (had cookies but no encryption keys)
	needsPINSetup     bool

	client   *twittermeow.Client
	settings *response.AccountSettingsResponse
}

var (
	LoginStepIDCookies   = "fi.mau.twitter.login.enter_cookies"
	LoginStepSecretKey   = "fi.mau.twitter.login.secret_key"
	LoginStepJuiceboxPIN = "fi.mau.twitter.login.juicebox_pin"
	LoginStepIDComplete  = "fi.mau.twitter.login.complete"
)

var _ bridgev2.LoginProcessCookies = (*TwitterLogin)(nil)
var _ bridgev2.LoginProcessUserInput = (*TwitterLogin)(nil)
var _ bridgev2.LoginProcessWithOverride = (*TwitterLogin)(nil)

const (
	pinRegex            = "^[0-9]{4}$"
	passcodeBodyRecover = "To retrieve your encrypted messages, please enter your passcode below. For more information see: https://help.x.com/en/using-x/about-chat."
	passcodeBodySetup   = "No PIN code is registered yet. Register by creating your PIN code below."
)

var (
	ErrJuiceboxLocked = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.JUICEBOX_LOCKED",
		Err:        "Too many incorrect passcode attempts. X Chat is locked.",
		StatusCode: http.StatusForbidden,
	}
	ErrJuiceboxRateLimited = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.JUICEBOX_RATE_LIMITED",
		Err:        "Too many attempts. Please try again later.",
		StatusCode: http.StatusTooManyRequests,
	}
	ErrJuiceboxInvalidAuth = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.JUICEBOX_INVALID_AUTH",
		Err:        "Couldn't verify your passcode. Please try again.",
		StatusCode: http.StatusBadRequest,
	}
	ErrJuiceboxNotRegistered = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.JUICEBOX_NOT_REGISTERED",
		Err:        "Passcode isn't set up for X Chat. Set it up on x.com.",
		StatusCode: http.StatusBadRequest,
	}
	ErrJuiceboxUpgradeRequired = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.JUICEBOX_UPGRADE_REQUIRED",
		Err:        "This bridge is out of date. Update and try again.",
		StatusCode: http.StatusUpgradeRequired,
	}
	ErrJuiceboxTransient = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.JUICEBOX_TRANSIENT",
		Err:        "Temporary error. Try again.",
		StatusCode: http.StatusServiceUnavailable,
	}
	ErrMissingUserID = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.MISSING_USER_ID",
		Err:        "Couldn't read your X account ID. Please try again.",
		StatusCode: http.StatusInternalServerError,
	}
)

func (tc *TwitterConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{
		{
			Name:        "Cookies",
			Description: "Log in with your X account using your cookies",
			ID:          "cookies",
		},
	}
}

func (tc *TwitterConnector) CreateLogin(_ context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != "cookies" {
		return nil, bridgev2.ErrInvalidLoginFlowID
	}
	return &TwitterLogin{User: user, tc: tc}, nil
}

func (t *TwitterLogin) Start(_ context.Context) (*bridgev2.LoginStep, error) {
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       LoginStepIDCookies,
		Instructions: "Open the Login URL in an Incognito/Private browsing mode. Then, extract the cookies as a JSON object/cURL command copied from the Network tab of your browser's DevTools. After that, close the browser **before** pasting the cookies.\n\nFor example: `{\"ct0\":\"123466-...\",\"auth_token\":\"abcde-...\"}`",
		CookiesParams: &bridgev2.LoginCookiesParams{
			URL:       "https://x.com/i/flow/login",
			UserAgent: "",
			Fields: []bridgev2.LoginCookieField{
				{
					ID:       "ct0",
					Required: true,
					Sources: []bridgev2.LoginCookieFieldSource{
						{Type: bridgev2.LoginCookieTypeCookie, Name: "ct0"},
					},
				},
				{
					ID:       "auth_token",
					Required: true,
					Sources: []bridgev2.LoginCookieFieldSource{
						{Type: bridgev2.LoginCookieTypeCookie, Name: "auth_token"},
					},
				},
			},
		},
	}, nil
}

func (t *TwitterLogin) Cancel() {}

func makePINStep(errorLine string, isSetup bool) *bridgev2.LoginStep {
	instructions := passcodeBodyRecover
	fieldName := "Passcode"
	if isSetup {
		instructions = passcodeBodySetup
		fieldName = "Create your PIN code"
	}
	if errorLine != "" {
		instructions = fmt.Sprintf("%s\n\n%s", errorLine, instructions)
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       LoginStepJuiceboxPIN,
		Instructions: instructions,
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type:    bridgev2.LoginInputFieldType2FACode,
					ID:      "pin",
					Name:    fieldName,
					Pattern: pinRegex,
				},
			},
		},
	}
}

// StartWithOverride is called when re-authenticating an existing login.
// For migration users (cookies but no encryption keys), this skips to passcode step.
func (t *TwitterLogin) StartWithOverride(ctx context.Context, override *bridgev2.UserLogin) (*bridgev2.LoginStep, error) {
	meta, ok := override.Metadata.(*UserLoginMetadata)
	if !ok || meta == nil || meta.Cookies == "" {
		return t.Start(ctx)
	}

	// Migration case: validate existing cookies and skip to passcode
	t.Cookies = meta.Cookies
	if err := t.ensureClientForPIN(ctx); err != nil {
		// Cookies expired, fall back to normal flow
		t.User.Log.Warn().Err(err).Msg("Migration: cookies invalid, falling back to full login")
		return t.Start(ctx)
	}
	t.persistClientCookiesAndUserID()
	t.isMigration = true

	t.refreshPINSetupState(ctx, "Migration: failed to determine PIN setup state, using recovery prompt")

	t.User.Log.Info().Msg("Migration: cookies validated, skipping to passcode step")
	return makePINStep("", t.needsPINSetup), nil
}

func (t *TwitterLogin) SubmitCookies(ctx context.Context, cookies map[string]string) (*bridgev2.LoginStep, error) {
	cookieStruct := twitCookies.NewCookies(cookies)
	t.Cookies = cookieStruct.String()

	client := twittermeow.NewClient(cookieStruct, nil, t.User.Log.With().Str("component", "login_twitter_client").Logger())

	settings, err := client.LoadMessagesPage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load messages page after submitting cookies: %w", err)
	}
	t.client = client
	t.settings = settings
	t.persistClientCookiesAndUserID()

	t.refreshPINSetupState(ctx, "Failed to determine PIN setup state, using recovery prompt")

	return makePINStep("", t.needsPINSetup), nil
}

func parsePINInput(input map[string]string) (string, error) {
	pin, ok := input["pin"]
	if !ok {
		return "", fmt.Errorf("passcode input is required")
	}
	pin = strings.TrimSpace(pin)
	if pin == "" {
		return "", fmt.Errorf("passcode cannot be empty")
	}
	return pin, nil
}

func (t *TwitterLogin) ensureClientForPIN(ctx context.Context) error {
	if t.client != nil {
		return nil
	}

	if t.Cookies == "" {
		return fmt.Errorf("cookies must be submitted before passcode")
	}
	cookieStruct := twitCookies.NewCookiesFromString(t.Cookies)
	t.client = twittermeow.NewClient(cookieStruct, nil, t.User.Log.With().Str("component", "login_twitter_client").Logger())
	settings, err := t.client.LoadMessagesPage(ctx)
	if err != nil {
		return fmt.Errorf("failed to load messages page: %w", err)
	}
	t.settings = settings
	return nil
}

func (t *TwitterLogin) persistClientCookiesAndUserID() {
	// Persist any cookies set by LoadMessagesPage so subsequent sessions include them.
	t.Cookies = t.client.GetCookieString()
	t.client.SetCurrentUserID(t.client.GetCurrentUserID())
}

func (t *TwitterLogin) refreshPINSetupState(ctx context.Context, warnMsg string) {
	needsPINSetup, err := t.detectPINSetupNeeded(ctx)
	if err != nil {
		t.User.Log.Warn().Err(err).Msg(warnMsg)
		return
	}
	t.needsPINSetup = needsPINSetup
}

func (t *TwitterLogin) detectPINSetupNeeded(ctx context.Context) (bool, error) {
	currentUserID := strings.TrimSpace(t.client.GetCurrentUserID())
	if currentUserID == "" {
		return false, ErrMissingUserID
	}
	publicKeysResp, err := t.client.GetPublicKeys(ctx, []string{currentUserID})
	if err != nil {
		return false, err
	}
	_, hasJuiceboxTokens := firstPublicKeyWithJuiceboxTokens(publicKeysResp)
	needsSetup := !hasJuiceboxTokens
	t.User.Log.Debug().
		Bool("has_juicebox_tokens", hasJuiceboxTokens).
		Bool("needs_pin_setup", needsSetup).
		Msg("Determined PIN setup state from GetPublicKeys")
	return needsSetup, nil
}

func firstPublicKeyWithJuiceboxTokens(data *response.GetPublicKeysResponse) (response.PublicKeyWithTokenMap, bool) {
	if len(data.Data.UserResultsByRestIDs) == 0 {
		return response.PublicKeyWithTokenMap{}, false
	}
	withTokens := data.Data.UserResultsByRestIDs[0].Result.GetPublicKeys.PublicKeysWithTokenMap
	for _, keyData := range withTokens {
		if len(keyData.TokenMap.TokenMap) == 0 {
			continue
		}
		return keyData, true
	}
	return response.PublicKeyWithTokenMap{}, false
}

func resolveJuiceboxConfigAndTokens(tokenMap response.KeyStoreTokenMap) (string, map[string]string) {
	juiceboxConfigJSON := tokenMap.KeyStoreTokenMapJSON
	authTokens := make(map[string]string, len(tokenMap.TokenMap))
	for _, entry := range tokenMap.TokenMap {
		authTokens[strings.ToLower(entry.Key)] = entry.Value.Token
	}
	return juiceboxConfigJSON, authTokens
}

func validateRecoveredKey(name, key string) error {
	if key == "" {
		return nil
	}
	if _, err := crypto.ParsePrivateKeyScalar(key); err != nil {
		return fmt.Errorf("recovered invalid %s: %w", name, err)
	}
	return nil
}

func handleRecoverPasscodeError(err error) (*bridgev2.LoginStep, error, bool) {
	var recoverErr *juiceboxgo.RecoverError
	if !errors.As(err, &recoverErr) {
		return nil, nil, false
	}

	if recoverErr.GuessesRemaining == nil {
		return makePINStep("Invalid passcode.", false), nil, true
	}

	guessesLeft := *recoverErr.GuessesRemaining
	if guessesLeft > 0 {
		guessWord := "guesses"
		if guessesLeft == 1 {
			guessWord = "guess"
		}
		return makePINStep(fmt.Sprintf("Invalid passcode. You have %d %s remaining.", guessesLeft, guessWord), false), nil, true
	}

	return nil, ErrJuiceboxLocked, true
}

func mapCommonJuiceboxError(err error) error {
	switch {
	case errors.Is(err, juiceboxgo.ErrRateLimitExceeded):
		return ErrJuiceboxRateLimited
	case errors.Is(err, juiceboxgo.ErrInvalidAuth):
		return ErrJuiceboxInvalidAuth
	case errors.Is(err, juiceboxgo.ErrUpgradeRequired):
		return ErrJuiceboxUpgradeRequired
	case errors.Is(err, juiceboxgo.ErrTransient):
		return ErrJuiceboxTransient
	default:
		return nil
	}
}

func (t *TwitterLogin) bootstrapJuiceboxPIN(ctx context.Context, pin string) (*KeyBackupData, string, error) {
	bootstrapData, err := GenerateFirstTimePINBootstrapData()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate first-time key material: %w", err)
	}

	requestVersion := strconv.FormatInt(time.Now().UnixMilli(), 10)
	addResp, err := t.client.AddXChatPublicKey(ctx, &payload.AddXChatPublicKeyMutationVariables{
		Version:         requestVersion,
		GenerateVersion: true,
		PublicKey: payload.AddXChatPublicKeyInput{
			PublicKey:                  bootstrapData.PublicKeySPKI,
			SigningPublicKey:           bootstrapData.SigningPublicKeySPKI,
			IdentityPublicKeySignature: bootstrapData.IdentityPublicKeySignature,
			RegistrationMethod:         payload.AddXChatPublicKeyRegistrationMethodCustomPin,
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to register xchat public key: %w", err)
	}
	if len(addResp.Errors) > 0 {
		return nil, "", fmt.Errorf("failed to register xchat public key: %s", addResp.Errors[0].Message)
	}

	juiceboxConfigJSON, authTokens := resolveJuiceboxConfigAndTokens(addResp.Data.UserAddPublicKey.TokenMap)
	if juiceboxConfigJSON == "" {
		return nil, "", fmt.Errorf("xchat public key registration returned empty juicebox config")
	}
	if len(authTokens) == 0 {
		return nil, "", fmt.Errorf("xchat public key registration returned no juicebox auth tokens")
	}

	signingKeyVersion := strings.TrimSpace(addResp.Data.UserAddPublicKey.Version)
	if signingKeyVersion == "" {
		return nil, "", fmt.Errorf("xchat public key registration returned empty key version")
	}

	juiceboxLogger := t.User.Log.With().Str("component", "juicebox").Logger()
	juiceboxLogger.Debug().
		Str("juicebox_config", juiceboxConfigJSON).
		Int("juicebox_config_len", len(juiceboxConfigJSON)).
		Any("auth_tokens", authTokens).
		Int("auth_tokens_count", len(authTokens)).
		Int("max_guess_count", addResp.Data.UserAddPublicKey.TokenMap.MaxGuessCount).
		Msg("Juicebox bootstrap parameters")

	err = RegisterSecretToJuicebox(
		ctx,
		juiceboxConfigJSON,
		authTokens,
		pin,
		"",
		bootstrapData.RawSecret,
		addResp.Data.UserAddPublicKey.TokenMap.MaxGuessCount,
		juiceboxLogger,
	)
	if err != nil {
		if mappedErr := mapCommonJuiceboxError(err); mappedErr != nil {
			return nil, "", mappedErr
		}
		return nil, "", fmt.Errorf("failed to bootstrap juicebox secret: %w", err)
	}

	return &KeyBackupData{
		SecretKey:  bootstrapData.SecretKey,
		SigningKey: bootstrapData.SigningKey,
	}, signingKeyVersion, nil
}

func (t *TwitterLogin) recoverJuiceboxPIN(
	ctx context.Context,
	pin string,
	keyData response.PublicKeyWithTokenMap,
) (*KeyBackupData, string, *bridgev2.LoginStep, error) {
	juiceboxConfigJSON, authTokens := resolveJuiceboxConfigAndTokens(keyData.TokenMap)
	if juiceboxConfigJSON == "" || len(authTokens) == 0 {
		return nil, "", nil, ErrJuiceboxNotRegistered
	}

	juiceboxLogger := t.User.Log.With().Str("component", "juicebox").Logger()
	juiceboxLogger.Debug().
		Str("juicebox_config", juiceboxConfigJSON).
		Int("juicebox_config_len", len(juiceboxConfigJSON)).
		Any("auth_tokens", authTokens).
		Int("auth_tokens_count", len(authTokens)).
		Msg("Juicebox recovery parameters")

	keys, err := RecoverKeysFromJuicebox(ctx, juiceboxConfigJSON, authTokens, pin, "", juiceboxLogger)
	if err != nil {
		if retryStep, handledErr, handled := handleRecoverPasscodeError(err); handled {
			return nil, "", retryStep, handledErr
		}
		if mappedErr := mapCommonJuiceboxError(err); mappedErr != nil {
			return nil, "", nil, mappedErr
		}
		if errors.Is(err, juiceboxgo.ErrNotRegistered) {
			return nil, "", nil, ErrJuiceboxNotRegistered
		}
		return nil, "", nil, fmt.Errorf("failed to recover keys: %w", err)
	}

	// SigningKeyVersion comes from the API response, not Juicebox (binary data doesn't include it).
	return keys, keyData.PublicKeyWithMetadata.Version, nil, nil
}

func (t *TwitterLogin) SubmitUserInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	pin, err := parsePINInput(input)
	if err != nil {
		return nil, err
	}

	if err = t.ensureClientForPIN(ctx); err != nil {
		return nil, err
	}
	t.persistClientCookiesAndUserID()

	// Get recovery config from X API
	publicKeysResp, err := t.client.GetPublicKeys(ctx, []string{t.client.GetCurrentUserID()})
	if err != nil {
		return nil, fmt.Errorf("failed to get public keys: %w", err)
	}

	var (
		keys              *KeyBackupData
		signingKeyVersion string
	)

	keyData, hasJuiceboxTokens := firstPublicKeyWithJuiceboxTokens(publicKeysResp)
	needsPINSetup := !hasJuiceboxTokens
	if needsPINSetup {
		keys, signingKeyVersion, err = t.bootstrapJuiceboxPIN(ctx, pin)
		if err != nil {
			return nil, err
		}
	} else {
		var retryStep *bridgev2.LoginStep
		keys, signingKeyVersion, retryStep, err = t.recoverJuiceboxPIN(ctx, pin, keyData)
		if err != nil {
			return nil, err
		}
		if retryStep != nil {
			return retryStep, nil
		}
	}

	t.SecretKey = keys.SecretKey
	t.SigningKey = keys.SigningKey
	t.SigningKeyVersion = signingKeyVersion

	if err := validateRecoveredKey("secret key", t.SecretKey); err != nil {
		return nil, err
	}
	if err := validateRecoveredKey("signing key", t.SigningKey); err != nil {
		return nil, err
	}

	meta := &UserLoginMetadata{
		Cookies:           t.Cookies,
		SecretKey:         t.SecretKey,
		SigningKey:        t.SigningKey,
		SigningKeyVersion: t.SigningKeyVersion,
	}

	// If this is a migration, mark it and flag for full encrypted room sync
	if t.isMigration {
		now := time.Now()
		meta.MigratedAt = &now
		meta.PendingEncryptedSync = true
		meta.Session = nil          // Clear cached session to force full resync
		meta.MaxUserSequenceID = "" // Reset sequence ID to fetch all messages
		t.User.Log.Info().Msg("Migration: flagged for full encrypted room backfill")
	}

	remoteProfile := &status.RemoteProfile{
		Username: t.settings.ScreenName,
	}
	currentUserID := strings.TrimSpace(t.client.GetCurrentUserID())
	if currentUserID == "" {
		return nil, ErrMissingUserID
	}
	id := MakeUserLoginID(currentUserID)
	ul, err := t.User.NewLogin(
		ctx,
		&database.UserLogin{
			ID:            id,
			Metadata:      meta,
			RemoteName:    remoteProfile.Username,
			RemoteProfile: *remoteProfile,
		},
		&bridgev2.NewLoginParams{
			DeleteOnConflict:  true,
			DontReuseExisting: false,
			LoadUserLogin: func(ctx context.Context, login *bridgev2.UserLogin) error {
				if t.client != nil {
					t.client.SetKeyStore(newUserLoginKeyStore(login, t.tc))
				}
				t.client.Logger = login.Log.With().Str("component", "twitter_client").Logger()
				login.Client = NewTwitterClient(login, t.tc, t.client)
				return nil
			},
		},
	)
	if err != nil {
		return nil, err
	}

	go func(ctx context.Context, client *TwitterClient) {
		client.DoConnect(ctx)
	}(context.WithoutCancel(ctx), ul.Client.(*TwitterClient))

	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeComplete,
		StepID:       LoginStepIDComplete,
		Instructions: fmt.Sprintf("Successfully logged into X as @%s", ul.UserLogin.RemoteName),
		CompleteParams: &bridgev2.LoginCompleteParams{
			UserLoginID: ul.ID,
			UserLogin:   ul,
		},
	}, nil
}
