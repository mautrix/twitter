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
	"strings"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
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
	pinRegex      = "^[0-9]{4}$"
	passcodeTitle = "Enter Passcode"
	passcodeBody  = "To retrieve your encrypted messages, please enter your passcode below."
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

func makePINStep(errorLine string) *bridgev2.LoginStep {
	instructions := fmt.Sprintf("%s\n\n%s", passcodeTitle, passcodeBody)
	if errorLine != "" {
		instructions = fmt.Sprintf("**%s**\n\n%s", errorLine, instructions)
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
					Name:    "Passcode",
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
	cookieStruct := twitCookies.NewCookiesFromString(meta.Cookies)
	t.client = twittermeow.NewClient(cookieStruct, nil, t.User.Log.With().Str("component", "login_twitter_client").Logger())

	settings, err := t.client.LoadMessagesPage(ctx)
	if err != nil {
		// Cookies expired, fall back to normal flow
		t.User.Log.Warn().Err(err).Msg("Migration: cookies invalid, falling back to full login")
		return t.Start(ctx)
	}

	t.settings = settings
	t.Cookies = t.client.GetCookieString()
	t.client.SetCurrentUserID(t.client.GetCurrentUserID())
	t.isMigration = true

	t.User.Log.Info().Msg("Migration: cookies validated, skipping to passcode step")
	return makePINStep(""), nil
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
	// Refresh cookies with any values set during LoadMessagesPage (e.g., twid)
	t.Cookies = t.client.GetCookieString()
	t.client.SetCurrentUserID(t.client.GetCurrentUserID())

	return makePINStep(""), nil
}

func (t *TwitterLogin) SubmitUserInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	pin, ok := input["pin"]
	if !ok {
		return nil, fmt.Errorf("passcode input is required")
	}
	pin = strings.TrimSpace(pin)
	if pin == "" {
		return nil, fmt.Errorf("passcode cannot be empty")
	}

	if t.client == nil {
		if t.Cookies == "" {
			return nil, fmt.Errorf("cookies must be submitted before passcode")
		}
		cookieStruct := twitCookies.NewCookiesFromString(t.Cookies)
		t.client = twittermeow.NewClient(cookieStruct, nil, t.User.Log.With().Str("component", "login_twitter_client").Logger())
		settings, err := t.client.LoadMessagesPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load messages page: %w", err)
		}
		t.settings = settings
	}

	// Persist any cookies set by LoadMessagesPage so subsequent sessions include them.
	t.Cookies = t.client.GetCookieString()
	t.client.SetCurrentUserID(t.client.GetCurrentUserID())

	// Get recovery config from X API
	publicKeysResp, err := t.client.GetPublicKeys(ctx, []string{t.client.GetCurrentUserID()})
	if err != nil {
		return nil, fmt.Errorf("failed to get public keys: %w", err)
	}

	if len(publicKeysResp.Data.UserResultsByRestIDs) == 0 {
		return nil, ErrJuiceboxNotRegistered
	}

	userResult := publicKeysResp.Data.UserResultsByRestIDs[0]
	if len(userResult.Result.GetPublicKeys.PublicKeysWithTokenMap) == 0 {
		return nil, ErrJuiceboxNotRegistered
	}

	keyData := userResult.Result.GetPublicKeys.PublicKeysWithTokenMap[0]
	juiceboxConfigJSON := keyData.TokenMap.KeyStoreTokenMapJSON

	// Validate config JSON is not empty
	if juiceboxConfigJSON == "" {
		return nil, ErrJuiceboxNotRegistered
	}

	// Build auth tokens map from token_map entries
	// Maps realm ID (hex string, lowercase) to pre-fetched JWT auth token
	authTokens := make(map[string]string)
	for _, entry := range keyData.TokenMap.TokenMap {
		authTokens[strings.ToLower(entry.Key)] = entry.Value.Token
	}

	if len(authTokens) == 0 {
		return nil, ErrJuiceboxNotRegistered
	}

	juiceboxLogger := t.User.Log.With().Str("component", "juicebox").Logger()
	juiceboxLogger.Debug().
		Str("juicebox_config", juiceboxConfigJSON).
		Int("juicebox_config_len", len(juiceboxConfigJSON)).
		Any("auth_tokens", authTokens).
		Int("auth_tokens_count", len(authTokens)).
		Msg("Juicebox recovery parameters")

	// Recover keys from Juicebox (user info must be empty)
	keys, err := RecoverKeysFromJuicebox(ctx, juiceboxConfigJSON, authTokens, pin, "", juiceboxLogger)
	if err != nil {
		// Check if this is an invalid passcode error that allows retry
		var recoverErr *juiceboxgo.RecoverError
		if errors.As(err, &recoverErr) && recoverErr.GuessesRemaining != nil {
			guessesLeft := *recoverErr.GuessesRemaining
			if guessesLeft > 0 {
				guessWord := "guesses"
				if guessesLeft == 1 {
					guessWord = "guess"
				}
				// Return the same step with error message to allow retry
				return makePINStep(
					fmt.Sprintf("Invalid Passcode. You have %d %s remaining.", guessesLeft, guessWord),
				), nil
			}
			// No guesses remaining - user is locked out
			return nil, ErrJuiceboxLocked
		}
		if errors.As(err, &recoverErr) && recoverErr.GuessesRemaining == nil {
			return makePINStep("Invalid Passcode."), nil
		}
		if errors.Is(err, juiceboxgo.ErrRateLimitExceeded) {
			return nil, ErrJuiceboxRateLimited
		} else if errors.Is(err, juiceboxgo.ErrInvalidAuth) {
			return nil, ErrJuiceboxInvalidAuth
		} else if errors.Is(err, juiceboxgo.ErrNotRegistered) {
			return nil, ErrJuiceboxNotRegistered
		} else if errors.Is(err, juiceboxgo.ErrUpgradeRequired) {
			return nil, ErrJuiceboxUpgradeRequired
		} else if errors.Is(err, juiceboxgo.ErrTransient) {
			return nil, ErrJuiceboxTransient
		}
		// Other errors (network, auth, etc.) - return as-is
		return nil, fmt.Errorf("failed to recover keys: %w", err)
	}

	t.SecretKey = keys.SecretKey
	t.SigningKey = keys.SigningKey
	// SigningKeyVersion comes from the API response, not Juicebox (binary data doesn't include it)
	t.SigningKeyVersion = keyData.PublicKeyWithMetadata.Version

	// Validate recovered keys
	if t.SecretKey != "" {
		if _, err := crypto.ParsePrivateKeyScalar(t.SecretKey); err != nil {
			return nil, fmt.Errorf("recovered invalid secret key: %w", err)
		}
	}
	if t.SigningKey != "" {
		if _, err := crypto.ParsePrivateKeyScalar(t.SigningKey); err != nil {
			return nil, fmt.Errorf("recovered invalid signing key: %w", err)
		}
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
