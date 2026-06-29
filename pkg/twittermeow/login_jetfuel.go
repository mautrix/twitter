package twittermeow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

const (
	webLoginBackendOCF     webLoginBackend = "ocf"
	webLoginBackendJetfuel webLoginBackend = "jetfuel"

	jetfuelHeaderVersion = "JP-5"
)

var (
	jetfuelActionPathRegex = regexp.MustCompile(`/onboarding/web/actions/[A-Za-z0-9_./-]+`)
	jetfuelFieldRegex      = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$-]{1,80}$`)
	jetfuelNumericIDRegex  = regexp.MustCompile(`\b[0-9]{5,30}\b`)
	jetfuelUUIDRegex       = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	jetfuelActionAliases   = map[string]string{
		"begin_login":            endpoints.JETFUEL_BEGIN_LOGIN_PATH,
		"login_enter_password":   endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH,
		"begin_two_factor_auth":  endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH,
		"finish_two_factor_auth": endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH,
		"two_factor_code":        "/onboarding/web/actions/two_factor_code",
	}
)

type webLoginBackend string

type jetfuelLoginState struct {
	identifier         string
	passwordAction     string
	verificationAction string
	verificationFields []string
	twoFactorAction    string
	twoFactorMethods   []WebLoginAuthMethod
	sessionToken       string
	preludeDispatchID  string
	userID             string
}

type jetfuelLoginResponse struct {
	strings []string
	paths   []string
	fields  []string
	raw     []byte
}

func (wls *WebLoginSession) startJetfuel(ctx context.Context) (*WebLoginResult, error) {
	if err := wls.client.loadPage(ctx, endpoints.BASE_FLOW_LOGIN_URL); err != nil {
		return nil, fmt.Errorf("failed to load X login page: %w", err)
	}
	if _, err := wls.client.jetfuelGet(ctx, endpoints.JETFUEL_LANDING_PATH); err != nil {
		wls.client.Logger.Debug().Err(err).Msg("Jetfuel login landing preflight failed")
	}
	body, err := wls.client.jetfuelGet(ctx, endpoints.JETFUEL_LOGIN_PATH)
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if !parsed.hasPath(endpoints.JETFUEL_BEGIN_LOGIN_PATH) && !parsed.hasField("username_or_email") {
		return nil, fmt.Errorf("%w: jetfuel login page did not expose a username action", ErrWebLoginUnexpectedSubtask)
	}
	wls.backend = webLoginBackendJetfuel
	wls.jetfuel = &jetfuelLoginState{}
	return &WebLoginResult{
		Status:           WebLoginStatusNeedsIdentifier,
		CurrentSubtaskID: "JetfuelBeginLogin",
		Challenge: &WebLoginChallenge{
			SubtaskID: "JetfuelBeginLogin",
			Hint:      "Phone, email, or username",
		},
	}, nil
}

func (wls *WebLoginSession) submitJetfuelIdentifier(ctx context.Context, identifier string) (*WebLoginResult, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, fmt.Errorf("x username, email, or phone is required")
	}
	if err := wls.client.sendJetfuelViewerContextEvent(ctx); err != nil {
		wls.client.Logger.Debug().Err(err).Msg("Jetfuel viewer-context preflight failed")
	}
	body, err := wls.client.jetfuelPostForm(ctx, endpoints.JETFUEL_BEGIN_LOGIN_PATH, url.Values{
		"username_or_email": {identifier},
	})
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if err := parsed.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(parsed)
	wls.jetfuel.identifier = identifier
	if wls.client.IsLoggedIn() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if action := parsed.passwordAction(); action != "" {
		wls.jetfuel.passwordAction = action
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsPassword,
			CurrentSubtaskID: "JetfuelPassword",
			Challenge: &WebLoginChallenge{
				SubtaskID: "JetfuelPassword",
				Hint:      "Password",
			},
		}, nil
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        parsed.verificationChallenge(),
		}, nil
	}
	return nil, fmt.Errorf("%w: jetfuel identifier response did not expose a supported next action", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) submitJetfuelCredentials(ctx context.Context, identifier, password string) (*WebLoginResult, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, fmt.Errorf("x username, email, or phone is required")
	}
	if password == "" {
		return nil, fmt.Errorf("x password is required")
	}
	result, err := wls.submitJetfuelIdentifier(ctx, identifier)
	if err != nil {
		if !isJetfuelPrePasswordParityError(err) {
			return nil, err
		}
		wls.client.Logger.Debug().Err(err).Msg("Jetfuel sequential identifier submit failed, trying combined credentials submit")
		return wls.submitJetfuelCombinedCredentials(ctx, identifier, password)
	}
	if result.Status != WebLoginStatusNeedsPassword {
		return result, nil
	}
	return wls.submitJetfuelPassword(ctx, password)
}

func isJetfuelPrePasswordParityError(err error) bool {
	var webErr *WebLoginError
	if !errors.As(err, &webErr) {
		return false
	}
	text := strings.ToLower(webErr.Message)
	return webErr.Code == 399 && (strings.Contains(text, "temporarily limited") ||
		strings.Contains(text, "official x apps") || strings.Contains(text, "use x.com"))
}

func (wls *WebLoginSession) submitJetfuelCombinedCredentials(ctx context.Context, identifier, password string) (*WebLoginResult, error) {
	body, err := wls.client.jetfuelPostForm(ctx, endpoints.JETFUEL_BEGIN_LOGIN_PATH, url.Values{
		"username_or_email": {identifier},
		"password":          {password},
	})
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if err := parsed.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(parsed)
	if wls.jetfuel != nil {
		wls.jetfuel.identifier = identifier
	}
	if wls.client.IsLoggedIn() || parsed.isComplete() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if result := wls.jetfuelAuthMethodChoiceResult(parsed); result != nil {
		return result, nil
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        parsed.verificationChallenge(),
		}, nil
	}
	if action := parsed.passwordAction(); action != "" {
		wls.jetfuel.passwordAction = action
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsPassword,
			CurrentSubtaskID: "JetfuelPassword",
			Challenge: &WebLoginChallenge{
				SubtaskID: "JetfuelPassword",
				Hint:      "Password",
			},
		}, nil
	}
	return nil, fmt.Errorf("%w: jetfuel credentials response did not complete or expose a supported challenge", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) submitJetfuelPassword(ctx context.Context, password string) (*WebLoginResult, error) {
	if password == "" {
		return nil, fmt.Errorf("x password is required")
	}
	if wls.jetfuel == nil || wls.jetfuel.passwordAction == "" {
		return nil, fmt.Errorf("%w: jetfuel password action is missing", ErrWebLoginUnexpectedSubtask)
	}
	form := url.Values{
		"password": {password},
	}
	if wls.jetfuel.identifier != "" {
		form.Set("username", wls.jetfuel.identifier)
	}
	if wls.jetfuel.sessionToken != "" {
		form.Set("session_token", wls.jetfuel.sessionToken)
	}
	body, err := wls.client.jetfuelPostForm(ctx, wls.jetfuel.passwordAction, form)
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if err := parsed.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(parsed)
	if wls.client.IsLoggedIn() || parsed.isComplete() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if result := wls.jetfuelAuthMethodChoiceResult(parsed); result != nil {
		return result, nil
	}
	if action := parsed.beginTwoFactorAction(); action != "" {
		return wls.submitJetfuelBeginTwoFactor(ctx, action)
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        parsed.verificationChallenge(),
		}, nil
	}
	return nil, fmt.Errorf("%w: jetfuel password response did not complete or expose a supported challenge", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) submitJetfuelBeginTwoFactor(ctx context.Context, action string) (*WebLoginResult, error) {
	if wls.jetfuel == nil {
		return nil, fmt.Errorf("%w: jetfuel session state is missing", ErrWebLoginUnexpectedSubtask)
	}
	form := url.Values{}
	if wls.jetfuel.preludeDispatchID != "" {
		form.Set("prelude_dispatch_id", wls.jetfuel.preludeDispatchID)
	}
	if wls.jetfuel.sessionToken != "" {
		form.Set("session_token", wls.jetfuel.sessionToken)
	}
	body, err := wls.client.jetfuelPostForm(ctx, action, form)
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if err := parsed.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(parsed)
	if wls.client.IsLoggedIn() || parsed.isComplete() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if result := wls.jetfuelAuthMethodChoiceResult(parsed); result != nil {
		return result, nil
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        parsed.verificationChallenge(),
		}, nil
	}
	return nil, fmt.Errorf("%w: jetfuel two-factor prelude did not expose a verification challenge", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) submitJetfuelText(ctx context.Context, text string) (*WebLoginResult, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("x verification code is required")
	}
	if wls.jetfuel == nil || wls.jetfuel.verificationAction == "" {
		return nil, fmt.Errorf("%w: jetfuel verification action is missing", ErrWebLoginUnexpectedSubtask)
	}
	form := url.Values{}
	for _, field := range wls.jetfuelVerificationFields() {
		form.Set(field, text)
	}
	if wls.jetfuel.sessionToken != "" {
		form.Set("session_token", wls.jetfuel.sessionToken)
	}
	if wls.jetfuel.preludeDispatchID != "" {
		form.Set("prelude_dispatch_id", wls.jetfuel.preludeDispatchID)
	}
	body, err := wls.client.jetfuelPostForm(ctx, wls.jetfuel.verificationAction, form)
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if err := parsed.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(parsed)
	if wls.client.IsLoggedIn() || parsed.isComplete() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if result := wls.jetfuelAuthMethodChoiceResult(parsed); result != nil {
		return result, nil
	}
	if action := parsed.passwordAction(); action != "" {
		wls.jetfuel.passwordAction = action
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsPassword,
			CurrentSubtaskID: "JetfuelPassword",
			Challenge: &WebLoginChallenge{
				SubtaskID: "JetfuelPassword",
				Hint:      "Password",
			},
		}, nil
	}
	if action := parsed.beginTwoFactorAction(); action != "" {
		return wls.submitJetfuelBeginTwoFactor(ctx, action)
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        parsed.verificationChallenge(),
		}, nil
	}
	return nil, fmt.Errorf("%w: jetfuel verification response did not complete login", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) updateJetfuelState(parsed jetfuelLoginResponse) {
	if wls.jetfuel == nil {
		wls.jetfuel = &jetfuelLoginState{}
	}
	if token := parsed.uuidValue("session_token"); token != "" {
		wls.jetfuel.sessionToken = token
	}
	if id := parsed.uuidValue("prelude_dispatch_id"); id != "" {
		wls.jetfuel.preludeDispatchID = id
	}
	if id := parsed.numericValue("user_id"); id != "" {
		wls.jetfuel.userID = id
	}
	if action := parsed.passwordAction(); action != "" {
		wls.jetfuel.passwordAction = action
	}
	if action := parsed.beginTwoFactorAction(); action != "" {
		wls.jetfuel.twoFactorAction = action
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
	}
	if methods := parsed.authMethods(); len(methods) > 0 {
		wls.jetfuel.twoFactorMethods = methods
	}
}

func (wls *WebLoginSession) jetfuelVerificationFields() []string {
	if wls.jetfuel != nil && len(wls.jetfuel.verificationFields) > 0 {
		return wls.jetfuel.verificationFields
	}
	return defaultJetfuelVerificationFields()
}

func (wls *WebLoginSession) jetfuelAuthMethodChoiceResult(parsed jetfuelLoginResponse) *WebLoginResult {
	methods := parsed.authMethods()
	if len(methods) == 0 && wls.jetfuel != nil {
		methods = wls.jetfuel.twoFactorMethods
	}
	if len(methods) == 0 || !parsed.isAuthMethodChoice() {
		return nil
	}
	if wls.jetfuel != nil {
		wls.jetfuel.twoFactorMethods = methods
	}
	return &WebLoginResult{
		Status:           WebLoginStatusNeedsAuthMethod,
		CurrentSubtaskID: "JetfuelTwoFactorMethod",
		Challenge: &WebLoginChallenge{
			SubtaskID:   "JetfuelTwoFactorMethod",
			Hint:        "Verification method",
			Description: "Choose how to verify this X login.",
			IsTwoFactor: true,
		},
		AuthMethods: methods,
	}
}

func (wls *WebLoginSession) submitJetfuelAuthMethod(ctx context.Context, methodID string) (*WebLoginResult, error) {
	if wls.jetfuel == nil || len(wls.jetfuel.twoFactorMethods) == 0 {
		return nil, ErrWebLoginMissingAuthMethodState
	}
	method, ok := wls.jetfuel.findAuthMethod(methodID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrWebLoginUnsupportedAuthMethod, strings.TrimSpace(methodID))
	}
	action := wls.jetfuel.twoFactorAction
	if action == "" {
		action = endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH
	}
	body, err := wls.client.jetfuelPostForm(ctx, action, wls.jetfuel.authMethodForm(method))
	if err != nil {
		return nil, err
	}
	parsed := parseJetfuelLoginResponse(body)
	if err := parsed.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(parsed)
	if wls.client.IsLoggedIn() || parsed.isComplete() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if method.Kind == WebLoginAuthMethodKindSecurityKey {
		return wls.submitJetfuelSecurityKeyChallenge(ctx, method, parsed)
	}
	if action := parsed.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = parsed.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        parsed.verificationChallengeForMethod(method),
		}, nil
	}
	if result := wls.jetfuelAuthMethodChoiceResult(parsed); result != nil {
		return result, nil
	}
	return nil, fmt.Errorf("%w: jetfuel auth method response did not expose a verification challenge", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) submitJetfuelSecurityKeyChallenge(ctx context.Context, method WebLoginAuthMethod, parsed jetfuelLoginResponse) (*WebLoginResult, error) {
	challenge, ok := parsed.webAuthnChallenge()
	if !ok {
		return nil, missingWebAuthnChallengeError(method.Name)
	}
	challengeResponse, err := createWebAuthnChallengeResponse(ctx, challenge)
	if err != nil {
		return nil, err
	}
	action := parsed.verificationAction()
	if action == "" {
		action = endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH
	}
	form := webAuthnChallengeSubmitForm(challengeResponse, parsed.verificationCodeFields(), wls.jetfuel.sessionToken, wls.jetfuel.preludeDispatchID)
	body, err := wls.client.jetfuelPostForm(ctx, action, form)
	if err != nil {
		return nil, err
	}
	next := parseJetfuelLoginResponse(body)
	if err := next.loginError(); err != nil {
		return nil, err
	}
	wls.updateJetfuelState(next)
	if wls.client.IsLoggedIn() || next.isComplete() {
		return &WebLoginResult{Status: WebLoginStatusComplete}, nil
	}
	if result := wls.jetfuelAuthMethodChoiceResult(next); result != nil {
		return result, nil
	}
	if action := next.verificationAction(); action != "" {
		wls.jetfuel.verificationAction = action
		wls.jetfuel.verificationFields = next.verificationCodeFields()
		return &WebLoginResult{
			Status:           WebLoginStatusNeedsText,
			CurrentSubtaskID: "JetfuelVerification",
			Challenge:        next.verificationChallengeForMethod(method),
		}, nil
	}
	return nil, fmt.Errorf("%w: jetfuel security-key response did not complete login", ErrWebLoginUnexpectedSubtask)
}

func (jls *jetfuelLoginState) findAuthMethod(methodID string) (WebLoginAuthMethod, bool) {
	methodID = normalizeJetfuelMethodID(methodID)
	for _, method := range jls.twoFactorMethods {
		if normalizeJetfuelMethodID(method.ID) == methodID || normalizeJetfuelMethodID(method.Name) == methodID {
			return method, true
		}
	}
	return WebLoginAuthMethod{}, false
}

func (jls *jetfuelLoginState) authMethodForm(method WebLoginAuthMethod) url.Values {
	form := url.Values{
		"two_factor_auth_method_type": {method.ID},
		"_selected_method_idx":        {fmt.Sprintf("%d", method.Index)},
	}
	if jls.userID != "" {
		form.Set("user_id", jls.userID)
	}
	if jls.sessionToken != "" {
		form.Set("session_token", jls.sessionToken)
	}
	return form
}

func (c *Client) jetfuelGet(ctx context.Context, path string) ([]byte, error) {
	return c.jetfuelRequest(ctx, path, http.MethodGet, nil)
}

func (c *Client) jetfuelPostForm(ctx context.Context, path string, form url.Values) ([]byte, error) {
	if err := c.addJetfuelCastleTokenToForm(form); err != nil {
		c.Logger.Trace().Err(err).Msg("Failed to create Castle request token")
	}
	return c.jetfuelRequest(ctx, path, http.MethodPost, []byte(form.Encode()))
}

func (c *Client) addJetfuelCastleTokenToForm(form url.Values) error {
	if form.Get("$castle_token") != "" {
		return nil
	}
	token, err := createCurrentCastleRequestToken(c.session.ClientUUID)
	if err == nil {
		form.Set("$castle_token", token)
		return nil
	}
	if fallbackErr := addCastleTokenToForm(form); fallbackErr != nil {
		return fmt.Errorf("current Castle: %w; fallback Castle: %w", err, fallbackErr)
	}
	c.Logger.Trace().Err(err).Msg("Fell back to legacy Castle request token")
	return nil
}

func (c *Client) jetfuelRequest(ctx context.Context, path, method string, body []byte) ([]byte, error) {
	fullURL := endpoints.JETFUEL_BASE_URL + ensureLeadingSlash(path)
	txID, err := crypto.SignTransaction(c.session.AnimationToken, c.session.VerificationToken, fullURL, method)
	if err != nil {
		c.Logger.Trace().Err(err).Msg("Failed to create X Jetfuel client transaction ID")
		txID = "e:"
	}
	extra := map[string]string{
		"accept":                  "*/*",
		"origin":                  endpoints.BASE_URL,
		"priority":                "u=1, i",
		"sec-fetch-dest":          "empty",
		"sec-fetch-mode":          "cors",
		"sec-fetch-site":          "same-origin",
		"timezone":                jetfuelTimezone(),
		"x-client-transaction-id": txID,
		"x-jf-client-theme":       "light",
		"x-jf-v":                  jetfuelHeaderVersion,
		"x-twitter-active-user":   "yes",
	}
	c.addJetfuelBrowserParityHeaders(extra)
	if csrfToken := c.cookies.Get(cookies.XCt0); csrfToken != "" {
		extra["x-csrf-token"] = csrfToken
	}
	headers := c.buildHeaders(HeaderOpts{
		WithNonAuthBearer: true,
		WithCookies:       true,
		WithXGuestToken:   true,
		Referer:           endpoints.BASE_URL + "/",
		Extra:             extra,
	})
	headers.Del("x-twitter-client-language")
	contentType := types.ContentTypeNone
	if method == http.MethodPost {
		contentType = types.ContentTypeForm
	}
	resp, respBody, err := c.makeRequestDirect(ctx, fullURL, method, headers, body, contentType)
	if resp != nil {
		c.cookies.UpdateFromResponse(resp)
	}
	if err != nil {
		return respBody, err
	}
	return respBody, nil
}

type jetfuelViewerContextEvent struct {
	Category                          string                      `json:"_category_"`
	FormatVersion                     int                         `json:"format_version"`
	TriggeredOn                       int64                       `json:"triggered_on"`
	Items                             []any                       `json:"items"`
	EventNamespace                    jetfuelViewerEventNamespace `json:"event_namespace"`
	ClientEventSequenceStartTimestamp int64                       `json:"client_event_sequence_start_timestamp"`
	ClientEventSequenceNumber         int                         `json:"client_event_sequence_number"`
	ClientAppID                       string                      `json:"client_app_id"`
}

type jetfuelViewerEventNamespace struct {
	Page    string `json:"page"`
	Action  string `json:"action"`
	Element string `json:"element"`
	Client  string `json:"client"`
}

func (c *Client) sendJetfuelViewerContextEvent(ctx context.Context) error {
	now := time.Now().UnixMilli()
	event := jetfuelViewerContextEvent{
		Category:                          "client_event",
		FormatVersion:                     2,
		TriggeredOn:                       now,
		Items:                             []any{},
		EventNamespace:                    jetfuelViewerEventNamespace{Page: "front", Action: "click", Element: "continue", Client: "m5"},
		ClientEventSequenceStartTimestamp: now,
		ClientEventSequenceNumber:         1,
		ClientAppID:                       "3033300",
	}
	logPayload, err := json.Marshal([]jetfuelViewerContextEvent{event})
	if err != nil {
		return err
	}
	form := url.Values{
		"debug": {"true"},
		"log":   {string(logPayload)},
	}
	txID, err := crypto.SignTransaction(c.session.AnimationToken, c.session.VerificationToken, endpoints.VIEWER_CONTEXT_URL, http.MethodPost)
	if err != nil {
		c.Logger.Trace().Err(err).Msg("Failed to create X viewer-context client transaction ID")
		txID = "e:"
	}
	headers := c.buildHeaders(HeaderOpts{
		WithNonAuthBearer:   true,
		WithCookies:         true,
		WithXGuestToken:     true,
		WithXTwitterHeaders: true,
		Origin:              endpoints.BASE_URL,
		Referer:             endpoints.BASE_URL + "/",
		Extra: map[string]string{
			"accept":                  "*/*",
			"priority":                "u=1, i",
			"sec-fetch-dest":          "empty",
			"sec-fetch-mode":          "cors",
			"sec-fetch-site":          "same-site",
			"x-client-transaction-id": txID,
		},
	})
	resp, _, err := c.makeRequestDirect(ctx, endpoints.VIEWER_CONTEXT_URL, http.MethodPost, headers, []byte(form.Encode()), types.ContentTypeForm)
	if resp != nil {
		c.cookies.UpdateFromResponse(resp)
	}
	return err
}

func jetfuelTimezone() string {
	if timezone := strings.TrimSpace(os.Getenv("TWITTER_JETFUEL_TIMEZONE")); timezone != "" {
		return timezone
	}
	if local := time.Local.String(); strings.Contains(local, "/") {
		return local
	}
	_, offset := time.Now().Zone()
	switch offset {
	case -10 * 60 * 60:
		return "Pacific/Honolulu"
	case -9 * 60 * 60:
		return "America/Anchorage"
	case -8 * 60 * 60:
		return "America/Los_Angeles"
	case -7 * 60 * 60:
		return "America/Denver"
	case -6 * 60 * 60, -5 * 60 * 60:
		return "America/Chicago"
	case -4 * 60 * 60:
		return "America/New_York"
	default:
		return "UTC"
	}
}

func (c *Client) addJetfuelBrowserParityHeaders(extra map[string]string) {
	addTFEGuestCookieHeader(extra, "x-tfe-guest-cookie-id", c.cookies.Get(cookies.XGuestID))
	addTFEGuestCookieHeader(extra, "x-tfe-guest-cookie-id-ads", c.cookies.Get(cookies.XGuestIDAds))
	addTFEGuestCookieHeader(extra, "x-tfe-guest-cookie-id-marketing", c.cookies.Get(cookies.XGuestIDMarketing))
	if dtab := decodeDtabLocal(c.cookies.Get(cookies.XDtabLocal)); dtab != "" {
		extra["dtab-local"] = dtab
	}
}

func addTFEGuestCookieHeader(extra map[string]string, headerName, cookieValue string) {
	if value := decodeTFEGuestCookie(cookieValue); value != "" {
		extra[headerName] = value
	}
}

func decodeTFEGuestCookie(cookieValue string) string {
	value, err := url.QueryUnescape(cookieValue)
	if err != nil {
		value = cookieValue
	}
	return strings.TrimPrefix(value, "v1:")
}

func decodeDtabLocal(cookieValue string) string {
	if cookieValue == "" {
		return ""
	}
	value, err := url.QueryUnescape(cookieValue)
	if err != nil {
		value = cookieValue
	}
	if strings.HasPrefix(value, "/") && isPrintableASCII(value) {
		return value
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(value)
	}
	if err == nil && strings.HasPrefix(string(decoded), "/") && isPrintableASCII(string(decoded)) {
		return string(decoded)
	}
	if isPrintableASCII(value) {
		return value
	}
	return ""
}

func isPrintableASCII(value string) bool {
	for _, r := range value {
		if r < 0x20 || r > 0x7e {
			return false
		}
	}
	return value != ""
}

func ensureLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func parseJetfuelLoginResponse(body []byte) jetfuelLoginResponse {
	strs := extractJetfuelStrings(body)
	paths := make([]string, 0)
	fields := make([]string, 0)
	for _, str := range strs {
		for _, path := range jetfuelActionPathRegex.FindAllString(str, -1) {
			paths = appendJetfuelPath(paths, path)
		}
		if path := canonicalJetfuelActionPath(str); path != "" {
			paths = appendJetfuelPath(paths, path)
		}
		if jetfuelFieldRegex.MatchString(str) && !slices.Contains(fields, str) {
			fields = append(fields, str)
		}
	}
	return jetfuelLoginResponse{strings: strs, paths: paths, fields: fields, raw: body}
}

func canonicalJetfuelActionPath(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "/onboarding/web/actions/") {
		return value
	}
	return jetfuelActionAliases[value]
}

func appendJetfuelPath(paths []string, path string) []string {
	if path == "" || slices.Contains(paths, path) {
		return paths
	}
	return append(paths, path)
}

func extractJetfuelStrings(body []byte) []string {
	var out []string
	seen := make(map[string]struct{})
	start := -1
	for i := 0; i < len(body); {
		r, size := utf8.DecodeRune(body[i:])
		if r == utf8.RuneError && size == 1 {
			if start >= 0 {
				addJetfuelString(body[start:i], seen, &out)
				start = -1
			}
			i++
			continue
		}
		if isJetfuelStringRune(r) {
			if start < 0 {
				start = i
			}
		} else if start >= 0 {
			addJetfuelString(body[start:i], seen, &out)
			start = -1
		}
		i += size
	}
	if start >= 0 {
		addJetfuelString(body[start:], seen, &out)
	}
	return out
}

func isJetfuelStringRune(r rune) bool {
	return r == '\n' || r == '\r' || r == '\t' || r >= 0x20 && r != utf8.RuneError
}

func addJetfuelString(raw []byte, seen map[string]struct{}, out *[]string) {
	str := strings.TrimSpace(string(bytes.Trim(raw, "\x00")))
	if len(str) < 3 {
		return
	}
	if _, ok := seen[str]; ok {
		return
	}
	seen[str] = struct{}{}
	*out = append(*out, str)
}

func (jfr jetfuelLoginResponse) text() string {
	return strings.ToLower(strings.Join(jfr.strings, "\n"))
}

func (jfr jetfuelLoginResponse) hasPath(path string) bool {
	return slices.Contains(jfr.paths, path)
}

func (jfr jetfuelLoginResponse) hasField(field string) bool {
	return slices.Contains(jfr.fields, field)
}

func (jfr jetfuelLoginResponse) passwordAction() string {
	for _, path := range jfr.paths {
		lower := strings.ToLower(path)
		if strings.Contains(lower, "password") {
			return path
		}
	}
	if jfr.hasField("password") || strings.Contains(jfr.text(), "password") {
		return endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH
	}
	return ""
}

func (jfr jetfuelLoginResponse) beginTwoFactorAction() string {
	for _, path := range jfr.paths {
		if strings.Contains(strings.ToLower(path), "begin_two_factor_auth") {
			return path
		}
	}
	return ""
}

func (jfr jetfuelLoginResponse) verificationAction() string {
	for _, path := range jfr.paths {
		lower := strings.ToLower(path)
		if strings.Contains(lower, "begin_two_factor_auth") {
			continue
		}
		if strings.Contains(lower, "two_factor") || strings.Contains(lower, "2fa") ||
			strings.Contains(lower, "challenge") || strings.Contains(lower, "verification") {
			return path
		}
	}
	return ""
}

func (jfr jetfuelLoginResponse) verificationCodeFields() []string {
	fields := make([]string, 0, 8)
	text := jfr.text()
	preferred := []string{
		"challenge_response",
		"verification_code",
		"two_factor_code",
		"code",
		"token",
	}
	if strings.Contains(text, "backup code") {
		preferred = append([]string{"backup_code"}, preferred...)
	}
	for _, field := range preferred {
		if jfr.hasField(field) {
			fields = appendJetfuelVerificationField(fields, field)
		}
	}
	for _, field := range jfr.fields {
		lower := strings.ToLower(field)
		if field != lower {
			continue
		}
		if lower == "session_token" || lower == "prelude_dispatch_id" ||
			strings.Contains(lower, "csrf") || strings.Contains(lower, "oauth") ||
			strings.Contains(lower, "castle") {
			continue
		}
		if strings.Contains(lower, "code") || strings.Contains(lower, "otp") ||
			strings.Contains(lower, "challenge") || strings.Contains(lower, "response") ||
			strings.Contains(lower, "token") {
			fields = appendJetfuelVerificationField(fields, field)
		}
	}
	if len(fields) == 0 {
		return defaultJetfuelVerificationFields()
	}
	return fields
}

func appendJetfuelVerificationField(fields []string, field string) []string {
	field = strings.TrimSpace(field)
	if field == "" || slices.Contains(fields, field) {
		return fields
	}
	return append(fields, field)
}

func defaultJetfuelVerificationFields() []string {
	return []string{"challenge_response", "verification_code", "two_factor_code", "backup_code", "code"}
}

func (jfr jetfuelLoginResponse) uuidValue(field string) string {
	field = strings.ToLower(field)
	for i, str := range jfr.strings {
		if !strings.Contains(strings.ToLower(str), field) {
			continue
		}
		if uuid := firstJetfuelUUID(str); uuid != "" {
			return uuid
		}
		for next := i + 1; next < len(jfr.strings) && next <= i+6; next++ {
			if uuid := firstJetfuelUUID(jfr.strings[next]); uuid != "" {
				return uuid
			}
		}
	}
	return ""
}

func (jfr jetfuelLoginResponse) numericValue(field string) string {
	field = strings.ToLower(field)
	for i, str := range jfr.strings {
		if !strings.Contains(strings.ToLower(str), field) {
			continue
		}
		if id := firstJetfuelNumericID(str); id != "" {
			return id
		}
		for next := i + 1; next < len(jfr.strings) && next <= i+6; next++ {
			if id := firstJetfuelNumericID(jfr.strings[next]); id != "" {
				return id
			}
		}
	}
	return ""
}

func firstJetfuelUUID(value string) string {
	return jetfuelUUIDRegex.FindString(value)
}

func firstJetfuelNumericID(value string) string {
	return jetfuelNumericIDRegex.FindString(value)
}

func (jfr jetfuelLoginResponse) isAuthMethodChoice() bool {
	text := jfr.text()
	return strings.Contains(text, "select a method") ||
		strings.Contains(text, "choose the method") ||
		strings.Contains(text, "two_factor_method") ||
		strings.Contains(text, "two factor method")
}

func (jfr jetfuelLoginResponse) authMethods() []WebLoginAuthMethod {
	if !jfr.isAuthMethodChoice() {
		return nil
	}
	methods := make([]WebLoginAuthMethod, 0, 3)
	for _, str := range jfr.strings {
		method, ok := classifyJetfuelAuthMethod(str)
		if !ok || containsJetfuelAuthMethod(methods, method.ID) {
			continue
		}
		method.Index = len(methods)
		methods = append(methods, method)
	}
	return methods
}

func classifyJetfuelAuthMethod(value string) (WebLoginAuthMethod, bool) {
	normalized := normalizeJetfuelMethodID(value)
	switch normalized {
	case "totp", "authenticatorapp", "authenticationapp":
		return WebLoginAuthMethod{
			ID:          "Totp",
			Name:        "Authenticator App",
			Description: "Use the code from your authentication app.",
			Kind:        WebLoginAuthMethodKindCode,
			Supported:   true,
		}, true
	case "backupcode":
		return WebLoginAuthMethod{
			ID:          "BackupCode",
			Name:        "Backup Code",
			Description: "Use a backup code from your X account settings.",
			Kind:        WebLoginAuthMethodKindBackupCode,
			Supported:   true,
		}, true
	case "u2fsecuritykey", "securitykey", "securitykeypc", "passkey":
		return WebLoginAuthMethod{
			ID:          "U2fSecurityKey",
			Name:        "Security Key PC",
			Description: "Use a security key or passkey.",
			Kind:        WebLoginAuthMethodKindSecurityKey,
			Supported:   true,
		}, true
	default:
		return WebLoginAuthMethod{}, false
	}
}

func containsJetfuelAuthMethod(methods []WebLoginAuthMethod, id string) bool {
	normalized := normalizeJetfuelMethodID(id)
	return slices.ContainsFunc(methods, func(method WebLoginAuthMethod) bool {
		return normalizeJetfuelMethodID(method.ID) == normalized
	})
}

func normalizeJetfuelMethodID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, ".", "")
	return value
}

func (jfr jetfuelLoginResponse) verificationChallenge() *WebLoginChallenge {
	text := jfr.text()
	return &WebLoginChallenge{
		SubtaskID:   "JetfuelVerification",
		Hint:        "Verification code",
		Description: jetfuelChallengeDescription(text),
		IsTwoFactor: strings.Contains(text, "two-factor") || strings.Contains(text, "two factor") ||
			strings.Contains(text, "authentication code") || strings.Contains(text, "verification code") ||
			strings.Contains(text, "backup code") || strings.Contains(text, "totp"),
	}
}

func (jfr jetfuelLoginResponse) verificationChallengeForMethod(method WebLoginAuthMethod) *WebLoginChallenge {
	challenge := jfr.verificationChallenge()
	switch method.Kind {
	case WebLoginAuthMethodKindCode:
		challenge.Description = "Enter the code from your authentication app."
	case WebLoginAuthMethodKindBackupCode:
		challenge.Hint = "Backup code"
		challenge.Description = "Enter a backup code from X."
	}
	challenge.IsTwoFactor = true
	return challenge
}

func jetfuelChallengeDescription(text string) string {
	switch {
	case strings.Contains(text, "backup code"):
		return "Enter the code from your authentication app."
	case strings.Contains(text, "authentication code"):
		return "Enter the authentication code from X."
	case strings.Contains(text, "verification code"):
		return "Enter the verification code from X."
	default:
		return "X needs additional verification for this login."
	}
}

func (jfr jetfuelLoginResponse) isComplete() bool {
	text := jfr.text()
	return strings.Contains(text, "/home") || strings.Contains(text, "open_account")
}

func (jfr jetfuelLoginResponse) loginError() error {
	text := jfr.text()
	switch {
	case strings.Contains(text, "official x apps") || strings.Contains(text, "use x.com"):
		return &WebLoginError{Code: 399, Message: "Please use X.com or official X apps to proceed with log in/sign up."}
	case strings.Contains(text, "temporarily limited") || strings.Contains(text, "try again later"):
		return &WebLoginError{Code: 399, Message: "We've temporarily limited your login. Please try again later."}
	case strings.Contains(text, "too many attempts") || strings.Contains(text, "try again in a few minutes"):
		return &WebLoginError{Code: 399, Message: "Too many attempts. Try again in a few minutes."}
	case strings.Contains(text, "missing_account") || strings.Contains(text, "not registered"):
		return &WebLoginError{Code: 32, Message: "This email or username is not registered yet."}
	case strings.Contains(text, "wrong password") || strings.Contains(text, "incorrect password"):
		return &WebLoginError{Code: 32, Message: "Wrong password"}
	case strings.Contains(text, "could not log you in") || strings.Contains(text, "couldn't log you in"):
		return &WebLoginError{Code: 399, Message: "Could not log you in now. Please try again later."}
	default:
		return nil
	}
}
