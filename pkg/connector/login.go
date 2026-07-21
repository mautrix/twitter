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
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
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
	useCookieLogin    bool

	client              *twittermeow.Client
	webLogin            *twittermeow.WebLoginSession
	webLoginIdentifier  string
	webLoginPassword    string
	webLoginCastleStage string
	webLoginAuthMethod  string
	webLoginText        string
	webLoginChallenge   *twittermeow.WebLoginChallenge
	webLoginMethods     []twittermeow.WebLoginAuthMethod
	browserHeaders      twittermeow.BrowserHeaders
	profile             twittermeow.CurrentUserProfile
}

var (
	LoginFlowIDPassword        = "password"
	LoginFlowIDCookies         = "cookies"
	LoginStepIDCredentials     = "fi.mau.twitter.login.enter_credentials"
	LoginStepIDCastleToken     = "fi.mau.twitter.login.castle_token"
	LoginStepIDVerification    = "fi.mau.twitter.login.enter_verification"
	LoginStepIDAuthMethod      = "fi.mau.twitter.login.select_auth_method"
	LoginStepIDCookies         = "fi.mau.twitter.login.enter_cookies"
	LoginStepJuiceboxPIN       = "fi.mau.twitter.login.juicebox_pin"
	LoginStepIDComplete        = "fi.mau.twitter.login.complete"
	loginFieldIdentifier       = "identifier"
	loginFieldPassword         = "password"
	loginFieldCastleToken      = "castle_token"
	loginFieldBrowserUserAgent = "browser_user_agent"
	loginFieldBrowserSecCHUA   = "browser_sec_ch_ua"
	loginFieldBrowserPlatform  = "browser_sec_ch_ua_platform"
	loginFieldBrowserMobile    = "browser_sec_ch_ua_mobile"
	loginFieldVerificationCode = "verification_code"
	loginFieldAuthMethod       = "auth_method"
)

const (
	webLoginCastleStageIdentifier     = "identifier"
	webLoginCastleStagePassword       = "password"
	webLoginCastleStageCombined       = "combined"
	webLoginCastleStageBeginTwoFactor = "begin_two_factor"
	webLoginCastleStageAuthMethod     = "auth_method"
	webLoginCastleStageText           = "text"
	castleTokenWebviewURL             = "https://x.com/robots.txt"
	castleTokenContextURL             = "https://x.com/i/jf/onboarding/web?mode=login"
	castleTokenHeaderPrefix           = "text/plain, application/x-mautrix-twitter-castle;v="
	castleTokenBatchSize              = 8
)

var castleTokenCookieNames = []string{
	"__cf_bm",
	"__cuid",
	"gt",
	"guest_id",
	"guest_id_ads",
	"guest_id_marketing",
	"personalization_id",
}

type browserHeaderField struct {
	ID         string
	HeaderName string
	Required   bool
	Pattern    string
}

var browserHeaderFields = []browserHeaderField{
	{
		ID:         loginFieldBrowserUserAgent,
		HeaderName: "user-agent",
		Required:   true,
		Pattern:    `^[^\r\n]{1,1024}$`,
	},
	{
		ID:         loginFieldBrowserSecCHUA,
		HeaderName: "sec-ch-ua",
		Pattern:    `^[^\r\n]{1,1024}$`,
	},
	{
		ID:         loginFieldBrowserPlatform,
		HeaderName: "sec-ch-ua-platform",
		Pattern:    `^[^\r\n]{1,1024}$`,
	},
	{
		ID:         loginFieldBrowserMobile,
		HeaderName: "sec-ch-ua-mobile",
		Pattern:    `^\?[01]$`,
	},
}

func decodeCastleTokenInput(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	if idx := strings.Index(value, castleTokenHeaderPrefix); idx >= 0 {
		encoded := removeCastleTokenWhitespace(value[idx+len(castleTokenHeaderPrefix):])
		decoded, err := base64.RawURLEncoding.DecodeString(encoded)
		if err != nil {
			return "", fmt.Errorf("decode Castle header token: %w", err)
		}
		value = string(decoded)
	}
	return removeCastleTokenWhitespace(value), nil
}

func removeCastleTokenWhitespace(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, value)
}

func castleTokenFieldID(index int) string {
	if index <= 1 {
		return loginFieldCastleToken
	}
	return fmt.Sprintf("%s_%d", loginFieldCastleToken, index)
}

func castleTokenCookieFields() []bridgev2.LoginCookieField {
	fields := []bridgev2.LoginCookieField{
		{
			ID:       loginFieldCastleToken,
			Required: true,
			Pattern:  `^[\s\S]{128,}$`,
			Sources: []bridgev2.LoginCookieFieldSource{
				{
					Type: bridgev2.LoginCookieTypeSpecial,
					Name: "fi.mau.twitter.castle_token",
				},
				{
					Type: bridgev2.LoginCookieTypeLocalStorage,
					Name: "fi.mau.twitter.castle_token",
				},
			},
		},
	}
	for index := 2; index <= castleTokenBatchSize; index++ {
		fieldID := castleTokenFieldID(index)
		storageKey := fmt.Sprintf("fi.mau.twitter.castle_token_%d", index)
		fields = append(fields, bridgev2.LoginCookieField{
			ID:       fieldID,
			Required: false,
			Pattern:  `^[\s\S]{128,}$`,
			Sources: []bridgev2.LoginCookieFieldSource{
				{
					Type: bridgev2.LoginCookieTypeSpecial,
					Name: storageKey,
				},
				{
					Type: bridgev2.LoginCookieTypeLocalStorage,
					Name: storageKey,
				},
			},
		})
	}
	for _, field := range browserHeaderFields {
		fields = append(fields, bridgev2.LoginCookieField{
			ID:       field.ID,
			Required: field.Required,
			Pattern:  field.Pattern,
			Sources: []bridgev2.LoginCookieFieldSource{
				{
					Type:            bridgev2.LoginCookieTypeRequestHeader,
					Name:            field.HeaderName,
					RequestURLRegex: `^https://x\.com/`,
				},
				{
					Type: bridgev2.LoginCookieTypeSpecial,
					Name: field.ID,
				},
			},
		})
	}
	for _, name := range castleTokenCookieNames {
		fields = append(fields, bridgev2.LoginCookieField{
			ID:       name,
			Required: false,
			Sources: []bridgev2.LoginCookieFieldSource{
				{
					Type:         bridgev2.LoginCookieTypeCookie,
					Name:         name,
					CookieDomain: "x.com",
				},
				{
					Type:         bridgev2.LoginCookieTypeCookie,
					Name:         name,
					CookieDomain: ".x.com",
				},
				{
					Type: bridgev2.LoginCookieTypeSpecial,
					Name: "fi.mau.twitter.cookie." + name,
				},
				{
					Type: bridgev2.LoginCookieTypeLocalStorage,
					Name: "fi.mau.twitter.cookie." + name,
				},
			},
		})
	}
	return fields
}

var _ bridgev2.LoginProcessCookies = (*TwitterLogin)(nil)
var _ bridgev2.LoginProcessUserInput = (*TwitterLogin)(nil)
var _ bridgev2.LoginProcessWithOverride = (*TwitterLogin)(nil)

const (
	pinRegex            = "^[0-9]{4}$"
	passcodeBodyRecover = "To retrieve your encrypted messages, please enter your passcode below. For more information see: https://help.x.com/en/using-x/about-chat"
	passcodeBodySetup   = "No PIN code is registered yet. Register by creating your PIN code below or using the X app."
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
	ErrMissingLoginInput = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.MISSING_LOGIN_INPUT",
		Err:        "Missing required login input.",
		StatusCode: http.StatusBadRequest,
	}
	ErrWebLoginFailed = bridgev2.RespError{
		ErrCode:    "FI.MAU.TWITTER.LOGIN_FAILED",
		Err:        "X login failed.",
		StatusCode: http.StatusBadGateway,
	}
)

func (tc *TwitterConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{
		{
			Name:        "Username/password",
			Description: "Log in with your X username, email, or phone number and password",
			ID:          LoginFlowIDPassword,
		},
	}
}

func (tc *TwitterConnector) CreateLogin(_ context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID == "" {
		flowID = LoginFlowIDPassword
	}
	if flowID != LoginFlowIDPassword && flowID != LoginFlowIDCookies {
		return nil, bridgev2.ErrInvalidLoginFlowID
	}
	return &TwitterLogin{User: user, tc: tc, useCookieLogin: flowID == LoginFlowIDCookies}, nil
}

func (t *TwitterLogin) Start(_ context.Context) (*bridgev2.LoginStep, error) {
	if !t.useCookieLogin {
		return makeCredentialsStep(""), nil
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       LoginStepIDCookies,
		Instructions: "Open the Login URL in an Incognito/Private browsing mode. Then, extract the cookies as a JSON object/cURL command copied from the Network tab of your browser's DevTools. After that, close the browser **before** pasting the cookies.\n\nFor example: `{\"ct0\":\"123466-...\",\"auth_token\":\"abcde-...\"}`",
		CookiesParams: &bridgev2.LoginCookiesParams{
			URL: "https://x.com/i/flow/login",
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

func makeCredentialsStep(errorLine string) *bridgev2.LoginStep {
	instructions := "Enter your X username, email, or phone number and password."
	if errorLine != "" {
		instructions = fmt.Sprintf("%s\n\n%s", errorLine, instructions)
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       LoginStepIDCredentials,
		Instructions: instructions,
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type:        bridgev2.LoginInputFieldTypeUsername,
					ID:          loginFieldIdentifier,
					Name:        "Username, email, or phone",
					Description: "The identifier you use to sign in to X.",
				},
				{
					Type: bridgev2.LoginInputFieldTypePassword,
					ID:   loginFieldPassword,
					Name: "Password",
				},
			},
		},
	}
}

func makeCastleTokenStep(info twittermeow.JetfuelCastleTokenInfo, identifier, errorLine string) *bridgev2.LoginStep {
	instructions := "Generating an X browser token for this login."
	if errorLine != "" {
		instructions = fmt.Sprintf("%s\n\n%s", errorLine, instructions)
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       LoginStepIDCastleToken,
		Instructions: instructions,
		CookiesParams: &bridgev2.LoginCookiesParams{
			URL:               castleTokenWebviewURL,
			ExtractJS:         castleTokenExtractJS(info, identifier),
			WaitForURLPattern: `^https://x\.com/robots\.txt$`,
			Fields:            castleTokenCookieFields(),
			Hidden:            true,
		},
	}
}

func (t *TwitterLogin) makeWebLoginCastleTokenStep(errorLine string) *bridgev2.LoginStep {
	if t.webLogin == nil || t.webLogin.Client() == nil {
		t.webLoginCastleStage = ""
		return makeCredentialsStep("The X login session expired. Enter your X login details again.")
	}
	info := t.webLogin.Client().JetfuelCastleTokenInfo()
	if !info.IsValid() {
		t.webLoginCastleStage = ""
		t.User.Log.Warn().Msg("X Castle web metadata missing from login page bootstrap")
		return makeCredentialsStep("X did not provide the browser-token metadata needed for native login. Try again.")
	}
	return makeCastleTokenStep(info, t.webLoginIdentifier, errorLine)
}

func makeVerificationStep(challenge *twittermeow.WebLoginChallenge, errorLine string) *bridgev2.LoginStep {
	instructions := "X needs additional verification for this login."
	fieldName := "Verification"
	fieldType := bridgev2.LoginInputFieldTypeToken
	if challenge != nil {
		if challenge.Description != "" {
			instructions = challenge.Description
		} else if challenge.Hint != "" {
			instructions = challenge.Hint
		}
		switch challenge.InputKind {
		case twittermeow.WebLoginChallengeInputKindPhoneNumber:
			fieldName = "Phone number"
			fieldType = bridgev2.LoginInputFieldTypePhoneNumber
			if instructions == "" {
				instructions = "Enter the phone number associated with your X account."
			}
		case twittermeow.WebLoginChallengeInputKindCode:
			fieldName = "Verification code"
			fieldType = bridgev2.LoginInputFieldType2FACode
			if instructions == "" {
				instructions = "Enter the verification code from X."
			}
		default:
			if challenge.IsTwoFactor {
				fieldName = "Verification code"
				fieldType = bridgev2.LoginInputFieldType2FACode
				if instructions == "" {
					instructions = "Enter the verification code from X."
				}
			}
		}
	}
	if errorLine != "" {
		instructions = fmt.Sprintf("%s\n\n%s", errorLine, instructions)
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       LoginStepIDVerification,
		Instructions: instructions,
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type: fieldType,
					ID:   loginFieldVerificationCode,
					Name: fieldName,
				},
			},
		},
	}
}

func makeAuthMethodStep(methods []twittermeow.WebLoginAuthMethod, errorLine string) *bridgev2.LoginStep {
	instructions := "Choose how to verify this X login."
	if errorLine != "" {
		instructions = fmt.Sprintf("%s\n\n%s", errorLine, instructions)
	}
	options := make([]string, 0, len(methods))
	for _, method := range methods {
		if !method.Supported {
			continue
		}
		if method.Name == "" {
			options = append(options, method.ID)
		} else {
			options = append(options, method.Name)
		}
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       LoginStepIDAuthMethod,
		Instructions: instructions,
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type:        bridgev2.LoginInputFieldTypeSelect,
					ID:          loginFieldAuthMethod,
					Name:        "Verification method",
					Description: "Choose the method X should use for this login.",
					Options:     options,
				},
			},
		},
	}
}

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
	if meta.BrowserHeaders != nil {
		t.browserHeaders = *meta.BrowserHeaders
	}
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
	if t.isWaitingForWebLoginCastleToken() {
		return t.submitWebCastleTokenInput(ctx, cookies)
	}
	cookieStruct := twitCookies.NewCookies(cookies)
	t.Cookies = cookieStruct.String()

	client := twittermeow.NewClient(cookieStruct, nil, t.User.Log.With().Str("component", "login_twitter_client").Logger())

	profile, err := client.LoadMessagesPage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load messages page after submitting cookies: %w", err)
	}
	t.client = client
	t.profile = profile
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
	if t.browserHeaders.UserAgent != "" {
		t.client.SetBrowserHeaders(t.browserHeaders)
	}
	profile, err := t.client.LoadMessagesPage(ctx)
	if err != nil {
		return fmt.Errorf("failed to load messages page: %w", err)
	}
	t.profile = profile
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
	_, hasJuiceboxTokens := latestPublicKeyWithJuiceboxTokens(publicKeysResp)
	needsSetup := !hasJuiceboxTokens
	t.User.Log.Debug().
		Bool("has_juicebox_tokens", hasJuiceboxTokens).
		Bool("needs_pin_setup", needsSetup).
		Msg("Determined PIN setup state from GetPublicKeys")
	return needsSetup, nil
}

func latestPublicKeyWithJuiceboxTokens(data *response.GetPublicKeysResponse) (response.PublicKeyWithTokenMap, bool) {
	if data == nil || len(data.Data.UserResultsByRestIDs) == 0 {
		return response.PublicKeyWithTokenMap{}, false
	}
	withTokens := data.Data.UserResultsByRestIDs[0].Result.GetPublicKeys.PublicKeysWithTokenMap
	var latest response.PublicKeyWithTokenMap
	found := false
	for _, keyData := range withTokens {
		if len(keyData.TokenMap.TokenMap) == 0 {
			continue
		}
		if !found || methods.CompareSnowflake(keyData.PublicKeyWithMetadata.Version, latest.PublicKeyWithMetadata.Version) > 0 {
			latest = keyData
			found = true
		}
	}
	return latest, found
}

func resolveJuiceboxConfigAndTokens(tokenMap response.KeyStoreTokenMap) (string, map[string]string) {
	juiceboxConfigJSON := strings.TrimSpace(tokenMap.KeyStoreTokenMapJSON)
	authTokens := make(map[string]string, len(tokenMap.TokenMap))
	for _, entry := range tokenMap.TokenMap {
		realmID := strings.ToLower(strings.TrimSpace(entry.Key))
		token := strings.TrimSpace(entry.Value.Token)
		if realmID != "" && token != "" {
			authTokens[realmID] = token
		}
	}
	return juiceboxConfigJSON, authTokens
}

type juiceboxRegistrationChecker func(context.Context, string, map[string]string) error

func selectRegisteredJuiceboxKey(
	ctx context.Context,
	data *response.GetPublicKeysResponse,
	checkRegistration juiceboxRegistrationChecker,
) (response.PublicKeyWithTokenMap, bool, error) {
	if data == nil || len(data.Data.UserResultsByRestIDs) == 0 {
		return response.PublicKeyWithTokenMap{}, false, nil
	}

	candidates := slices.Clone(data.Data.UserResultsByRestIDs[0].Result.GetPublicKeys.PublicKeysWithTokenMap)
	slices.SortStableFunc(candidates, func(a, b response.PublicKeyWithTokenMap) int {
		return methods.CompareSnowflake(b.PublicKeyWithMetadata.Version, a.PublicKeyWithMetadata.Version)
	})

	hasCandidates := false
	for _, keyData := range candidates {
		configJSON, authTokens := resolveJuiceboxConfigAndTokens(keyData.TokenMap)
		if configJSON == "" || len(authTokens) == 0 {
			continue
		}
		hasCandidates = true
		err := checkRegistration(ctx, configJSON, authTokens)
		switch {
		case err == nil:
			return keyData, true, nil
		case errors.Is(err, juiceboxgo.ErrNotRegistered):
			continue
		default:
			return response.PublicKeyWithTokenMap{}, true, err
		}
	}
	if hasCandidates {
		return response.PublicKeyWithTokenMap{}, true, juiceboxgo.ErrNotRegistered
	}
	return response.PublicKeyWithTokenMap{}, false, nil
}

func (t *TwitterLogin) selectRegisteredJuiceboxKey(ctx context.Context, data *response.GetPublicKeysResponse) (response.PublicKeyWithTokenMap, bool, error) {
	logger := t.User.Log.With().Str("component", "juicebox").Logger()
	return selectRegisteredJuiceboxKey(ctx, data, func(ctx context.Context, configJSON string, authTokens map[string]string) error {
		return CheckJuiceboxRegistration(ctx, configJSON, authTokens, logger)
	})
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
		Int("juicebox_config_len", len(juiceboxConfigJSON)).
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
		Int("juicebox_config_len", len(juiceboxConfigJSON)).
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
	if _, ok := input["pin"]; ok {
		return t.submitPINInput(ctx, input)
	}
	if _, ok := input[loginFieldAuthMethod]; ok {
		return t.submitWebAuthMethodInput(ctx, input)
	}
	if _, ok := input[loginFieldVerificationCode]; ok || t.webLoginChallenge != nil {
		return t.submitWebVerificationInput(ctx, input)
	}
	if _, ok := input[loginFieldIdentifier]; ok || input[loginFieldPassword] != "" {
		return t.submitCredentialsInput(ctx, input)
	}
	return nil, ErrMissingLoginInput
}

func (t *TwitterLogin) submitCredentialsInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	identifier := strings.TrimSpace(input[loginFieldIdentifier])
	password := input[loginFieldPassword]
	if identifier == "" || password == "" {
		return nil, ErrMissingLoginInput
	}

	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, t.User.Log.With().Str("component", "login_twitter_client").Logger())
	t.webLogin = twittermeow.NewWebLoginSession(client)
	t.webLoginIdentifier = identifier
	t.webLoginPassword = password
	t.webLoginCastleStage = ""
	t.webLoginAuthMethod = ""
	t.webLoginText = ""
	t.webLoginChallenge = nil
	t.webLoginMethods = nil

	result, err := t.webLogin.Start(ctx)
	if err != nil {
		return nil, webLoginFailureError(err)
	}
	return t.continueStartedCredentialsLogin(ctx, result)
}

func (t *TwitterLogin) continueStartedCredentialsLogin(ctx context.Context, result *twittermeow.WebLoginResult) (*bridgev2.LoginStep, error) {
	if result.Status != twittermeow.WebLoginStatusNeedsIdentifier {
		return t.handleWebLoginResult(ctx, result)
	}
	if t.webLogin.UsesJetfuel() {
		return t.startWebCastleLogin(), nil
	}

	result, err := t.webLogin.SubmitCredentials(ctx, t.webLoginIdentifier, t.webLoginPassword)
	if err != nil {
		return handleWebLoginCredentialsError(err)
	}
	return t.handleWebLoginResult(ctx, result)
}

func (t *TwitterLogin) startWebCastleLogin() *bridgev2.LoginStep {
	// The browser sends username and password together in begin_login. Starting with the
	// combined form avoids advancing the same Jetfuel session with an identifier-only POST.
	t.webLoginCastleStage = webLoginCastleStageCombined
	return t.makeWebLoginCastleTokenStep("")
}

func (t *TwitterLogin) isWaitingForWebLoginCastleToken() bool {
	return t.webLogin != nil && t.webLoginCastleStage != "" && t.webLoginIdentifier != "" && t.webLoginPassword != ""
}

func (t *TwitterLogin) submitWebCastleTokenInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	if !t.isWaitingForWebLoginCastleToken() {
		t.webLoginCastleStage = ""
		return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
	}
	castleTokens, err := decodeCastleTokenBatchInput(input)
	if err != nil {
		return t.makeWebLoginCastleTokenStep("The X webview returned an invalid browser token."), nil
	}
	if len(castleTokens) == 0 {
		return t.makeWebLoginCastleTokenStep("The X webview did not return a browser token."), nil
	}
	client := t.webLogin.Client()
	if !client.SetBrowserHeaders(browserHeadersFromInput(input)) {
		return t.makeWebLoginCastleTokenStep("The X webview did not return a valid browser fingerprint."), nil
	}
	t.browserHeaders = client.GetBrowserHeaders()
	client.SetCookies(castleWebviewCookies(input))
	client.SetNextJetfuelCastleTokens(castleTokens)
	return t.continueWebCastleLogin(ctx)
}

func browserHeadersFromInput(input map[string]string) twittermeow.BrowserHeaders {
	return twittermeow.BrowserHeaders{
		UserAgent:      input[loginFieldBrowserUserAgent],
		SecCHUserAgent: input[loginFieldBrowserSecCHUA],
		SecCHPlatform:  input[loginFieldBrowserPlatform],
		SecCHMobile:    input[loginFieldBrowserMobile],
	}
}

func decodeCastleTokenBatchInput(input map[string]string) ([]string, error) {
	tokens := make([]string, 0, castleTokenBatchSize)
	seen := make(map[string]struct{}, castleTokenBatchSize)
	for index := 1; index <= castleTokenBatchSize; index++ {
		token, err := decodeCastleTokenInput(input[castleTokenFieldID(index)])
		if err != nil {
			return nil, err
		}
		if token == "" {
			continue
		}
		if len(token) < 128 || len(token) > 20000 {
			return nil, fmt.Errorf("invalid Castle token length")
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (t *TwitterLogin) continueWebCastleLogin(ctx context.Context) (*bridgev2.LoginStep, error) {
	for attempts := 0; attempts < castleTokenBatchSize; attempts++ {
		stage := t.webLoginCastleStage
		zerolog.Ctx(ctx).Debug().
			Str("stage", stage).
			Bool("castle_token_available", t.webLogin.Client().HasNextJetfuelCastleToken()).
			Msg("Processing X Jetfuel Castle stage")

		var (
			result *twittermeow.WebLoginResult
			err    error
		)
		switch stage {
		case webLoginCastleStageIdentifier:
			result, err = t.webLogin.SubmitIdentifier(ctx, t.webLoginIdentifier)
			if err != nil && twittermeow.IsWebLoginPrePasswordParityError(err) {
				t.webLoginCastleStage = webLoginCastleStageCombined
				continue
			}
		case webLoginCastleStagePassword:
			result, err = t.webLogin.SubmitPassword(ctx, t.webLoginPassword)
		case webLoginCastleStageCombined:
			result, err = t.webLogin.SubmitCombinedCredentials(ctx, t.webLoginIdentifier, t.webLoginPassword)
		case webLoginCastleStageBeginTwoFactor:
			result, err = t.webLogin.SubmitPendingTwoFactor(ctx)
		case webLoginCastleStageAuthMethod:
			methodID := t.webLoginAuthMethod
			if methodID == "" {
				t.webLoginCastleStage = ""
				return makeAuthMethodStep(t.webLoginMethods, "Choose a verification method."), nil
			}
			result, err = t.webLogin.SubmitAuthMethod(ctx, methodID)
		case webLoginCastleStageText:
			text := t.webLoginText
			if text == "" {
				t.webLoginCastleStage = ""
				return makeVerificationStep(t.webLoginChallenge, "Enter the X verification code."), nil
			}
			result, err = t.webLogin.SubmitText(ctx, text)
		default:
			t.webLoginCastleStage = ""
			return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
		}
		if err != nil {
			if errors.Is(err, twittermeow.ErrJetfuelCastleTokenRequired) {
				if stage == webLoginCastleStagePassword || stage == webLoginCastleStageText {
					t.webLoginCastleStage = webLoginCastleStageBeginTwoFactor
					continue
				}
				return t.makeWebLoginCastleTokenStep("X needs a fresh browser token to continue this login."), nil
			}
			logWebCastleFailure(ctx, stage, err)
			t.webLoginCastleStage = ""
			return handleWebCastleStageError(stage, t.webLoginChallenge, t.webLoginMethods, err)
		}
		if stage == webLoginCastleStageAuthMethod {
			t.webLoginAuthMethod = ""
		}
		if stage == webLoginCastleStageText {
			t.webLoginText = ""
		}
		if result != nil && result.Status == twittermeow.WebLoginStatusNeedsPassword && t.webLoginPassword != "" {
			t.webLoginCastleStage = webLoginCastleStagePassword
			if !t.webLogin.Client().HasNextJetfuelCastleToken() {
				return t.makeWebLoginCastleTokenStep("X needs a fresh browser token to continue this login."), nil
			}
			continue
		}
		t.webLoginCastleStage = ""
		return t.handleWebLoginResult(ctx, result)
	}
	return t.makeWebLoginCastleTokenStep("X needs a fresh browser token to continue this login."), nil
}

func logWebCastleFailure(ctx context.Context, stage string, err error) {
	event := zerolog.Ctx(ctx).Debug().Str("stage", stage)
	var webErr *twittermeow.WebLoginError
	switch {
	case errors.As(err, &webErr):
		event = event.Str("error_kind", "x_response").Int("error_code", webErr.Code)
	case errors.Is(err, twittermeow.ErrWebLoginUnexpectedSubtask):
		event = event.Str("error_kind", "unsupported_response")
	default:
		event = event.Str("error_kind", "request")
	}
	event.Msg("X Jetfuel Castle stage failed")
}

func handleWebCastleStageError(stage string, challenge *twittermeow.WebLoginChallenge, methods []twittermeow.WebLoginAuthMethod, err error) (*bridgev2.LoginStep, error) {
	switch stage {
	case webLoginCastleStageIdentifier, webLoginCastleStagePassword, webLoginCastleStageCombined:
		return handleWebLoginCredentialsError(err)
	case webLoginCastleStageText:
		return handleWebLoginVerificationError(challenge, err)
	case webLoginCastleStageAuthMethod:
		return handleWebLoginAuthMethodError(methods, err)
	default:
		return nil, webLoginFailureError(err)
	}
}

func castleWebviewCookies(input map[string]string) map[string]string {
	out := make(map[string]string)
	for _, name := range castleTokenCookieNames {
		if value := strings.TrimSpace(input[name]); value != "" {
			out[name] = value
		}
	}
	return out
}

func (t *TwitterLogin) submitWebAuthMethodInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	if t.webLogin == nil {
		t.webLoginMethods = nil
		return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
	}
	methodID := strings.TrimSpace(input[loginFieldAuthMethod])
	if methodID == "" {
		return nil, ErrMissingLoginInput
	}
	if method, ok := findWebLoginAuthMethod(t.webLoginMethods, methodID); ok {
		methodID = method.ID
		if methodID == "" {
			methodID = method.Name
		}
	}
	if t.webLogin.UsesJetfuel() {
		t.webLoginAuthMethod = methodID
		t.webLoginCastleStage = webLoginCastleStageAuthMethod
		if t.webLogin.Client().HasNextJetfuelCastleToken() {
			return t.continueWebCastleLogin(ctx)
		}
		return t.makeWebLoginCastleTokenStep(""), nil
	}
	result, err := t.webLogin.SubmitAuthMethod(ctx, methodID)
	if err != nil {
		return nil, webLoginFailureError(err)
	}
	return t.handleWebLoginResult(ctx, result)
}

func findWebLoginAuthMethod(methods []twittermeow.WebLoginAuthMethod, selected string) (twittermeow.WebLoginAuthMethod, bool) {
	selected = normalizeLoginChoice(selected)
	for _, method := range methods {
		if normalizeLoginChoice(method.ID) == selected || normalizeLoginChoice(method.Name) == selected {
			return method, true
		}
	}
	return twittermeow.WebLoginAuthMethod{}, false
}

func normalizeLoginChoice(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, ".", "")
	return value
}

func (t *TwitterLogin) submitWebVerificationInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	if t.webLogin == nil {
		t.webLoginChallenge = nil
		t.webLoginPassword = ""
		t.webLoginMethods = nil
		return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
	}
	text := strings.TrimSpace(input[loginFieldVerificationCode])
	if text == "" {
		return nil, ErrMissingLoginInput
	}
	if t.webLogin.UsesJetfuel() {
		t.webLoginText = text
		t.webLoginCastleStage = webLoginCastleStageText
		if t.webLogin.Client().HasNextJetfuelCastleToken() {
			return t.continueWebCastleLogin(ctx)
		}
		return t.makeWebLoginCastleTokenStep(""), nil
	}
	result, err := t.webLogin.SubmitText(ctx, text)
	if err != nil {
		return handleWebLoginVerificationError(t.webLoginChallenge, err)
	}
	if result.Status == twittermeow.WebLoginStatusNeedsPassword && t.webLoginPassword != "" {
		if t.webLogin.UsesJetfuel() {
			t.webLoginCastleStage = webLoginCastleStagePassword
			return t.makeWebLoginCastleTokenStep(""), nil
		}
		result, err = t.webLogin.SubmitPassword(ctx, t.webLoginPassword)
		if err != nil {
			return handleWebLoginCredentialsError(err)
		}
	}
	return t.handleWebLoginResult(ctx, result)
}

func (t *TwitterLogin) handleWebLoginResult(ctx context.Context, result *twittermeow.WebLoginResult) (*bridgev2.LoginStep, error) {
	if result == nil {
		return makeCredentialsStep("X did not return a login step. Try again."), nil
	}
	switch result.Status {
	case twittermeow.WebLoginStatusComplete:
		return t.completeWebLogin(ctx)
	case twittermeow.WebLoginStatusNeedsAuthMethod:
		if len(result.AuthMethods) == 0 {
			return makeCredentialsStep("X returned a verification method chooser without any methods. Try again."), nil
		}
		t.webLoginChallenge = nil
		t.webLoginMethods = result.AuthMethods
		return makeAuthMethodStep(result.AuthMethods, ""), nil
	case twittermeow.WebLoginStatusNeedsText:
		t.webLoginChallenge = result.Challenge
		t.webLoginMethods = nil
		return makeVerificationStep(result.Challenge, ""), nil
	case twittermeow.WebLoginStatusNeedsPassword:
		if t.webLogin != nil && t.webLoginPassword != "" {
			if t.webLogin.UsesJetfuel() {
				t.webLoginCastleStage = webLoginCastleStagePassword
				return t.makeWebLoginCastleTokenStep(""), nil
			}
			next, err := t.webLogin.SubmitPassword(ctx, t.webLoginPassword)
			if err != nil {
				return handleWebLoginCredentialsError(err)
			}
			if next != nil && next.Status != twittermeow.WebLoginStatusNeedsPassword {
				return t.handleWebLoginResult(ctx, next)
			}
		}
		return makeCredentialsStep("X still needs your password. Enter your login details again."), nil
	case twittermeow.WebLoginStatusNeedsIdentifier:
		return makeCredentialsStep("X still needs your username, email, or phone. Enter your login details again."), nil
	default:
		t.User.Log.Warn().
			Str("subtask_id", result.CurrentSubtaskID).
			Str("status", string(result.Status)).
			Msg("X returned unsupported login subtask")
		return makeCredentialsStep(webLoginUnsupportedInstructions(result)), nil
	}
}

func webLoginUnsupportedInstructions(result *twittermeow.WebLoginResult) string {
	if result != nil && result.Challenge != nil {
		description := strings.TrimSpace(result.Challenge.Description)
		if description != "" {
			return description
		}
	}
	return "X returned a login challenge this bridge does not support yet."
}

func handleWebLoginCredentialsError(err error) (*bridgev2.LoginStep, error) {
	if isWebLoginCredentialsInputError(err) {
		return makeCredentialsStep(webLoginErrorInstructions(err)), nil
	}
	return nil, webLoginFailureError(err)
}

func handleWebLoginVerificationError(challenge *twittermeow.WebLoginChallenge, err error) (*bridgev2.LoginStep, error) {
	if isWebLoginVerificationInputError(err) {
		return makeVerificationStep(challenge, webLoginErrorInstructions(err)), nil
	}
	return nil, webLoginFailureError(err)
}

func handleWebLoginAuthMethodError(methods []twittermeow.WebLoginAuthMethod, err error) (*bridgev2.LoginStep, error) {
	if errors.Is(err, twittermeow.ErrWebLoginUnsupportedAuthMethod) {
		return makeAuthMethodStep(methods, webLoginErrorInstructions(err)), nil
	}
	if errors.Is(err, twittermeow.ErrWebLoginMissingAuthMethodState) {
		return makeCredentialsStep(webLoginErrorInstructions(err)), nil
	}
	return nil, webLoginFailureError(err)
}

func isWebLoginCredentialsInputError(err error) bool {
	var webErr *twittermeow.WebLoginError
	if !errors.As(err, &webErr) {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(webErr.Message))
	if webErr.Code != 32 {
		return false
	}
	return strings.Contains(msg, "wrong password") ||
		strings.Contains(msg, "incorrect password") ||
		strings.Contains(msg, "invalid password") ||
		strings.Contains(msg, "password you entered") ||
		strings.Contains(msg, "password is incorrect") ||
		strings.Contains(msg, "username and password") && strings.Contains(msg, "did not match") ||
		strings.Contains(msg, "invalid username or password") ||
		strings.Contains(msg, "invalid credentials") ||
		strings.Contains(msg, "missing_account") ||
		strings.Contains(msg, "not registered")
}

func isWebLoginVerificationInputError(err error) bool {
	var webErr *twittermeow.WebLoginError
	if !errors.As(err, &webErr) {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(webErr.Message))
	return strings.Contains(msg, "wrong code") ||
		strings.Contains(msg, "incorrect code") ||
		strings.Contains(msg, "invalid code") ||
		strings.Contains(msg, "code is incorrect") ||
		strings.Contains(msg, "verification code") && strings.Contains(msg, "incorrect") ||
		strings.Contains(msg, "authentication code") && strings.Contains(msg, "incorrect")
}

func webLoginFailureError(err error) error {
	if err == nil {
		return ErrWebLoginFailed
	}
	return ErrWebLoginFailed.WithMessage(webLoginErrorInstructions(err))
}

func (t *TwitterLogin) completeWebLogin(ctx context.Context) (*bridgev2.LoginStep, error) {
	if t.webLogin == nil || t.webLogin.Client() == nil {
		return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
	}
	client := t.webLogin.Client()
	profile, err := client.LoadMessagesPage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load authenticated X messages page after login: %w", err)
	}
	t.client = client
	t.profile = profile
	t.persistClientCookiesAndUserID()
	t.refreshPINSetupState(ctx, "Failed to determine PIN setup state after native login, using recovery prompt")
	t.webLoginIdentifier = ""
	t.webLoginPassword = ""
	t.webLoginCastleStage = ""
	t.webLoginAuthMethod = ""
	t.webLoginText = ""
	t.webLoginChallenge = nil
	t.webLoginMethods = nil

	return makePINStep("", t.needsPINSetup), nil
}

func webLoginErrorInstructions(err error) string {
	if err == nil {
		return "X rejected this login. Please check the details and try again."
	}
	var webErr *twittermeow.WebLoginError
	if errors.As(err, &webErr) {
		return webErr.UserMessage()
	}
	if errors.Is(err, twittermeow.ErrWebLoginUnexpectedSubtask) {
		return "X returned a login challenge this bridge does not support yet."
	}
	if errors.Is(err, twittermeow.ErrWebLoginUnsupportedAuthMethod) {
		return "That X verification method is not available for this login."
	}
	if errors.Is(err, twittermeow.ErrWebLoginMissingAuthMethodState) {
		return "The X verification method selection expired. Enter your X login details again."
	}
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "X rejected this login. Please check the details and try again."
	}
	return fmt.Sprintf("X login failed: %s", msg)
}

func (t *TwitterLogin) submitPINInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
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

	keyData, hasJuiceboxCandidates, selectionErr := t.selectRegisteredJuiceboxKey(ctx, publicKeysResp)
	needsPINSetup := !hasJuiceboxCandidates
	if needsPINSetup {
		keys, signingKeyVersion, err = t.bootstrapJuiceboxPIN(ctx, pin)
		if err != nil {
			return nil, err
		}
	} else {
		if selectionErr != nil {
			if mappedErr := mapCommonJuiceboxError(selectionErr); mappedErr != nil {
				return nil, mappedErr
			}
			if errors.Is(selectionErr, juiceboxgo.ErrNotRegistered) {
				return nil, ErrJuiceboxNotRegistered
			}
			return nil, fmt.Errorf("failed to select registered juicebox key: %w", selectionErr)
		}
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
	if browserHeaders := t.client.GetBrowserHeaders(); browserHeaders.UserAgent != "" {
		meta.BrowserHeaders = &browserHeaders
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

	currentUserID := strings.TrimSpace(t.client.GetCurrentUserID())
	if currentUserID == "" {
		return nil, ErrMissingUserID
	}

	remoteProfile := &status.RemoteProfile{
		Username: strings.TrimSpace(t.profile.ScreenName),
		Name:     strings.TrimSpace(t.profile.Name),
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

	if profile := t.profile; profile.AvatarURL != "" {
		updatedProfile := ul.Client.(*TwitterClient).makeXChatRemoteProfile(ctx, &types.User{
			IDStr:                currentUserID,
			ScreenName:           profile.ScreenName,
			Name:                 profile.Name,
			ProfileImageURLHTTPS: profile.AvatarURL,
		})
		if ul.UserLogin.RemoteName != updatedProfile.Username || ul.UserLogin.RemoteProfile != *updatedProfile {
			ul.UserLogin.RemoteName = updatedProfile.Username
			ul.UserLogin.RemoteProfile = *updatedProfile
			if err := ul.Save(ctx); err != nil {
				t.User.Log.Warn().Err(err).Msg("Failed to save login profile after syncing avatar")
			}
		}
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
