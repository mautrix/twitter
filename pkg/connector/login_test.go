package connector

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

type connectorRoundTripFunc func(*http.Request) (*http.Response, error)

func (rtf connectorRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rtf(req)
}

func connectorTestHTTPResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestSubmitUserInputRejectsMissingRequiredCredentialFields(t *testing.T) {
	login := &TwitterLogin{}
	tests := []map[string]string{
		{},
		{loginFieldIdentifier: "alice"},
		{loginFieldPassword: "secret"},
		{loginFieldIdentifier: "   ", loginFieldPassword: "secret"},
		{loginFieldIdentifier: "alice", loginFieldPassword: ""},
	}

	for _, input := range tests {
		step, err := login.SubmitUserInput(context.Background(), input)
		if step != nil {
			t.Fatalf("SubmitUserInput(%#v) step = %#v, want nil", input, step)
		}
		if !errors.Is(err, ErrMissingLoginInput) {
			t.Fatalf("SubmitUserInput(%#v) error = %v, want ErrMissingLoginInput", input, err)
		}
	}
}

func TestHandleWebLoginCredentialsErrorRetriesOnlyCredentialErrors(t *testing.T) {
	step, err := handleWebLoginCredentialsError(&twittermeow.WebLoginError{
		Code:    32,
		Message: "Wrong password",
	})
	if err != nil {
		t.Fatalf("handleWebLoginCredentialsError(wrong password) error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDCredentials {
		t.Fatalf("handleWebLoginCredentialsError(wrong password) step = %#v, want credentials step", step)
	}

	step, err = handleWebLoginCredentialsError(&twittermeow.WebLoginError{
		Code:    399,
		Message: "We've temporarily limited your login. Please try again later.",
	})
	if step != nil {
		t.Fatalf("handleWebLoginCredentialsError(temporary limit) step = %#v, want nil", step)
	}
	var respErr bridgev2.RespError
	if !errors.As(err, &respErr) || respErr.ErrCode != ErrWebLoginFailed.ErrCode {
		t.Fatalf("handleWebLoginCredentialsError(temporary limit) error = %#v, want ErrWebLoginFailed response", err)
	}
}

func TestGetLoginFlowsAdvertisesNativePasswordOnly(t *testing.T) {
	var tc TwitterConnector
	flows := tc.GetLoginFlows()

	if len(flows) != 1 {
		t.Fatalf("len(flows) = %d, want 1", len(flows))
	}
	if flows[0].ID != LoginFlowIDPassword {
		t.Fatalf("flow ID = %s, want %s", flows[0].ID, LoginFlowIDPassword)
	}
	if strings.Contains(strings.ToLower(flows[0].Description), "cookie") {
		t.Fatalf("flow description = %q, want native login flow", flows[0].Description)
	}
}

func TestContinueWebCastleLoginRequestsFreshTokenForActionlessPasswordReplay(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_VIEWER_CONTEXT", "0")
	const mainPageHTML = `<html><head><meta name="twitter-site-verification" content="verification-token"></head><body><script>
{"country": "US", "responsive_web_castle_public_key":{"value":"test-public-key"}}
gt=123456789
123:"ondemand.castle",{123:"abcdef"}
</script></body></html>`

	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	passwordRequestCount := 0
	client.HTTP = &http.Client{Transport: connectorRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/i/jf/onboarding/web":
			resp := connectorTestHTTPResponse(mainPageHTML)
			resp.Header.Add("Set-Cookie", "guest_id=v1%3A123456789; Path=/; Secure")
			return resp, nil
		case req.Method == http.MethodGet && req.URL.Path == "/i/jfapi"+endpoints.JETFUEL_LANDING_PATH:
			return connectorTestHTTPResponse("landing"), nil
		case req.Method == http.MethodGet && req.URL.Path == "/i/jfapi/onboarding/web" && req.URL.Query().Get("mode") == "login":
			return connectorTestHTTPResponse(endpoints.JETFUEL_BEGIN_LOGIN_PATH + "\x00username_or_email"), nil
		case req.Method == http.MethodPost && req.URL.Path == "/i/jfapi"+endpoints.JETFUEL_BEGIN_LOGIN_PATH:
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll(identifier request) error = %v", err)
			}
			if !strings.Contains(string(body), "%24castle_token=identifier-token") {
				t.Fatalf("identifier request body = %q", body)
			}
			return connectorTestHTTPResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
		case req.Method == http.MethodPost && req.URL.Path == "/i/jfapi"+endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH:
			passwordRequestCount++
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll(password request) error = %v", err)
			}
			wantToken := "first-password-token"
			if passwordRequestCount == 2 {
				wantToken = "replay-password-token"
			}
			if !strings.Contains(string(body), "%24castle_token="+wantToken) {
				t.Fatalf("password request %d body = %q, want token %q", passwordRequestCount, body, wantToken)
			}
			if passwordRequestCount == 1 {
				return connectorTestHTTPResponse("/onboarding/web/actions/persist_login_state\x00opaque_field"), nil
			}
			return connectorTestHTTPResponse(endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH + "\x00challenge_response\x00Enter your verification code"), nil
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
			return nil, nil
		}
	})}
	session := twittermeow.NewWebLoginSession(client)
	result, err := session.Start(context.Background())
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if !session.UsesJetfuel() || result == nil || result.Status != twittermeow.WebLoginStatusNeedsIdentifier {
		t.Fatalf("Start() result = %#v, UsesJetfuel = %t", result, session.UsesJetfuel())
	}

	client.SetNextJetfuelCastleTokens([]string{"identifier-token"})
	result, err = session.SubmitIdentifier(context.Background(), "test-user")
	if err != nil || result == nil || result.Status != twittermeow.WebLoginStatusNeedsPassword {
		t.Fatalf("SubmitIdentifier() result = %#v, error = %v", result, err)
	}

	client.SetNextJetfuelCastleTokens([]string{"first-password-token"})
	login := &TwitterLogin{
		User:                &bridgev2.User{Log: zerolog.Nop()},
		webLogin:            session,
		webLoginIdentifier:  "test-user",
		webLoginPassword:    "test-password",
		webLoginCastleStage: webLoginCastleStagePassword,
	}
	step, err := login.continueWebCastleLogin(context.Background())
	if err != nil {
		t.Fatalf("first continueWebCastleLogin() error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDCastleToken {
		t.Fatalf("first continueWebCastleLogin() step = %#v, want Castle token step", step)
	}
	if login.webLoginCastleStage != webLoginCastleStagePassword {
		t.Fatalf("Castle stage = %q, want password", login.webLoginCastleStage)
	}
	if passwordRequestCount != 1 {
		t.Fatalf("password request count = %d, want 1", passwordRequestCount)
	}

	client.SetNextJetfuelCastleTokens([]string{"replay-password-token"})
	step, err = login.continueWebCastleLogin(context.Background())
	if err != nil {
		t.Fatalf("second continueWebCastleLogin() error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDVerification {
		t.Fatalf("second continueWebCastleLogin() step = %#v, want verification step", step)
	}
	if passwordRequestCount != 2 {
		t.Fatalf("password request count = %d, want 2", passwordRequestCount)
	}
}

func TestCookieLoginRemainsVisible(t *testing.T) {
	login := &TwitterLogin{useCookieLogin: true}
	step, err := login.Start(context.Background())
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if step.CookiesParams == nil {
		t.Fatal("CookiesParams = nil")
	}
	if step.CookiesParams.Hidden {
		t.Fatal("CookiesParams.Hidden = true, want user-driven cookie login to remain visible")
	}
}

func TestMakeAuthMethodStepUsesNativeSelect(t *testing.T) {
	methods := []twittermeow.WebLoginAuthMethod{
		{ID: "Totp", Name: "Authenticator App", Supported: true},
		{ID: "Sms", Name: "Text Message", Supported: false},
		{ID: "BackupCode", Name: "Backup Code", Supported: true},
		{ID: "U2fSecurityKey", Name: "Security Key PC", Supported: false},
	}
	step := makeAuthMethodStep(methods, "")

	if step.Type != bridgev2.LoginStepTypeUserInput {
		t.Fatalf("Type = %s, want user input", step.Type)
	}
	if step.StepID != LoginStepIDAuthMethod {
		t.Fatalf("StepID = %s, want %s", step.StepID, LoginStepIDAuthMethod)
	}
	if step.UserInputParams == nil || len(step.UserInputParams.Fields) != 1 {
		t.Fatalf("UserInputParams = %#v, want one field", step.UserInputParams)
	}
	field := step.UserInputParams.Fields[0]
	if field.Type != bridgev2.LoginInputFieldTypeSelect {
		t.Fatalf("field.Type = %s, want select", field.Type)
	}
	if field.ID != loginFieldAuthMethod {
		t.Fatalf("field.ID = %s, want %s", field.ID, loginFieldAuthMethod)
	}
	if strings.Join(field.Options, ",") != "Authenticator App,Backup Code" {
		t.Fatalf("field.Options = %#v", field.Options)
	}
	if strings.Contains(step.Instructions, "not supported") {
		t.Fatalf("Instructions = %q, want no unsupported caveat", step.Instructions)
	}
}

func TestWebLoginUnsupportedInstructionsUsesChallengeDescription(t *testing.T) {
	result := &twittermeow.WebLoginResult{
		Status: twittermeow.WebLoginStatusUnsupported,
		Challenge: &twittermeow.WebLoginChallenge{
			Description: "Text message verification is coming soon.",
		},
	}

	if got := webLoginUnsupportedInstructions(result); got != "Text message verification is coming soon." {
		t.Fatalf("webLoginUnsupportedInstructions() = %q", got)
	}
}

func TestMakeCastleTokenStepUsesClientWebviewExtraction(t *testing.T) {
	if castleTokenBatchSize < 6 {
		t.Fatalf("castleTokenBatchSize = %d, want enough tokens for the six-request 2FA login path", castleTokenBatchSize)
	}
	info := twittermeow.JetfuelCastleTokenInfo{
		ScriptURL: "https://abs.twimg.com/responsive-web/client-web/ondemand.castle.1ff15ffa.js",
		PublicKey: "castle-public-key",
	}
	step := makeCastleTokenStep(info, "test-user", "")

	if step.Type != bridgev2.LoginStepTypeCookies {
		t.Fatalf("Type = %s, want cookies webview step", step.Type)
	}
	if step.StepID != LoginStepIDCastleToken {
		t.Fatalf("StepID = %s, want %s", step.StepID, LoginStepIDCastleToken)
	}
	if step.CookiesParams == nil {
		t.Fatal("CookiesParams = nil")
	}
	if !step.CookiesParams.Hidden {
		t.Fatal("CookiesParams.Hidden = false, want Castle token acquisition to run in a hidden webview")
	}
	if step.CookiesParams.UserAgent != "" {
		t.Fatalf("UserAgent = %q, want the webview's native user agent", step.CookiesParams.UserAgent)
	}
	if step.CookiesParams.URL != castleTokenWebviewURL {
		t.Fatalf("URL = %q", step.CookiesParams.URL)
	}
	if strings.Contains(step.CookiesParams.URL, "/i/flow/login") {
		t.Fatalf("URL = %q, want neutral webview page", step.CookiesParams.URL)
	}
	if !strings.Contains(step.CookiesParams.WaitForURLPattern, "robots") {
		t.Fatalf("WaitForURLPattern = %q, want neutral X URL", step.CookiesParams.WaitForURLPattern)
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, info.ScriptURL) ||
		!strings.Contains(step.CookiesParams.ExtractJS, "createRequestToken") {
		t.Fatalf("ExtractJS does not load X Castle token generator")
	}
	if strings.Contains(step.CookiesParams.ExtractJS, castleTokenJSConfigPlaceholder) {
		t.Fatal("ExtractJS still contains the embedded script config placeholder")
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, castleTokenContextURL) {
		t.Fatalf("ExtractJS does not include the X login context")
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, "showBrowserLoginStatus") ||
		!strings.Contains(step.CookiesParams.ExtractJS, "Signing in to X") ||
		!strings.Contains(step.CookiesParams.ExtractJS, "mautrix-twitter-login-status") ||
		!strings.Contains(step.CookiesParams.ExtractJS, "body.replaceChildren(container)") {
		t.Fatalf("ExtractJS does not replace robots.txt with the visible X login status")
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, "__BEEP_BEEP_AUTH_RESULTS__") {
		t.Fatalf("ExtractJS does not store the BrowserAuth result for Desktop polling")
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, "__MAUTRIX_TWITTER_CASTLE_IN_PROGRESS__") {
		t.Fatalf("ExtractJS does not guard against repeated BrowserAuth navigation runs")
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, "castleTokenBatchSize") {
		t.Fatalf("ExtractJS does not generate a Castle token batch")
	}
	fields := map[string]bridgev2.LoginCookieField{}
	for _, field := range step.CookiesParams.Fields {
		fields[field.ID] = field
	}
	if field, ok := fields[loginFieldCastleToken]; !ok || !field.Required {
		t.Fatalf("Fields = %#v, want required Castle token field", step.CookiesParams.Fields)
	} else if !hasCookieFieldSource(field, bridgev2.LoginCookieTypeLocalStorage, "fi.mau.twitter.castle_token") {
		t.Fatalf("Castle token field sources = %#v, want local_storage fallback", field.Sources)
	}
	for index := 2; index <= castleTokenBatchSize; index++ {
		fieldID := castleTokenFieldID(index)
		field, ok := fields[fieldID]
		if !ok {
			t.Fatalf("Fields missing optional Castle token batch field %q", fieldID)
		}
		if field.Required {
			t.Fatalf("Castle token batch field %q is required", fieldID)
		}
		if !hasCookieFieldSource(field, bridgev2.LoginCookieTypeLocalStorage, "fi.mau.twitter.castle_token_"+strconv.Itoa(index)) {
			t.Fatalf("Castle token batch field %q sources = %#v, want local_storage fallback", fieldID, field.Sources)
		}
	}
	for _, expected := range browserHeaderFields {
		field, ok := fields[expected.ID]
		if !ok {
			t.Fatalf("Fields missing browser header %q", expected.ID)
		}
		if field.Required != expected.Required {
			t.Fatalf("browser header %q required = %t, want %t", expected.ID, field.Required, expected.Required)
		}
		if !hasCookieFieldSource(field, bridgev2.LoginCookieTypeRequestHeader, expected.HeaderName) {
			t.Fatalf("browser header %q sources = %#v, want request header %q", expected.ID, field.Sources, expected.HeaderName)
		}
		if field.Sources[0].RequestURLRegex == "" || !hasCookieFieldSource(field, bridgev2.LoginCookieTypeSpecial, expected.ID) {
			t.Fatalf("browser header %q sources = %#v, want x.com request extraction and webview JS fallback", expected.ID, field.Sources)
		}
	}
	if !strings.Contains(step.CookiesParams.ExtractJS, "navigator.userAgent") ||
		!strings.Contains(step.CookiesParams.ExtractJS, "navigator.userAgentData") {
		t.Fatal("ExtractJS does not include the webview browser-header fallback")
	}
	for _, name := range castleTokenCookieNames {
		field, ok := fields[name]
		if !ok {
			t.Fatalf("Fields missing optional browser cookie %q", name)
		}
		if field.Required {
			t.Fatalf("browser cookie field %q is required", name)
		}
		if !strings.Contains(step.CookiesParams.ExtractJS, name) {
			t.Fatalf("ExtractJS does not include browser cookie %q", name)
		}
		if !hasCookieFieldSource(field, bridgev2.LoginCookieTypeLocalStorage, "fi.mau.twitter.cookie."+name) {
			t.Fatalf("browser cookie field %q sources = %#v, want local_storage fallback", name, field.Sources)
		}
	}
	for _, name := range []string{"auth_token", "ct0", "twid", "kdt"} {
		if _, ok := fields[name]; ok {
			t.Fatalf("Fields include auth cookie %q", name)
		}
	}
}

func TestBrowserHeadersFromInput(t *testing.T) {
	input := map[string]string{
		loginFieldBrowserUserAgent: "test user agent",
		loginFieldBrowserSecCHUA:   `"Chromium";v="150"`,
		loginFieldBrowserPlatform:  `"Android"`,
		loginFieldBrowserMobile:    "?1",
	}
	got := browserHeadersFromInput(input)
	if got.UserAgent != input[loginFieldBrowserUserAgent] ||
		got.SecCHUserAgent != input[loginFieldBrowserSecCHUA] ||
		got.SecCHPlatform != input[loginFieldBrowserPlatform] ||
		got.SecCHMobile != input[loginFieldBrowserMobile] {
		t.Fatalf("browserHeadersFromInput() = %#v", got)
	}
}

func TestUserLoginMetadataPersistsBrowserHeaders(t *testing.T) {
	meta := UserLoginMetadata{
		Cookies: "ct0=test",
		BrowserHeaders: &twittermeow.BrowserHeaders{
			UserAgent:      "test user agent",
			SecCHUserAgent: `"Chromium";v="150"`,
			SecCHPlatform:  `"Windows"`,
			SecCHMobile:    "?0",
		},
	}
	encoded, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	var decoded UserLoginMetadata
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if decoded.BrowserHeaders == nil || *decoded.BrowserHeaders != *meta.BrowserHeaders {
		t.Fatalf("decoded BrowserHeaders = %#v, want %#v", decoded.BrowserHeaders, meta.BrowserHeaders)
	}
}

func hasCookieFieldSource(field bridgev2.LoginCookieField, sourceType bridgev2.LoginCookieFieldSourceType, name string) bool {
	for _, source := range field.Sources {
		if source.Type == sourceType && source.Name == name {
			return true
		}
	}
	return false
}

func TestDecodeCastleTokenInputAcceptsHeaderCaptureValue(t *testing.T) {
	token := strings.Repeat("castle-token-", 64)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(token))
	got, err := decodeCastleTokenInput(castleTokenHeaderPrefix + encoded)
	if err != nil {
		t.Fatalf("decodeCastleTokenInput() failed: %v", err)
	}
	if got != token {
		t.Fatalf("decodeCastleTokenInput() = %q, want original token", got)
	}

	got, err = decodeCastleTokenInput(token)
	if err != nil {
		t.Fatalf("decodeCastleTokenInput(raw) failed: %v", err)
	}
	if got != token {
		t.Fatalf("decodeCastleTokenInput(raw) = %q, want original token", got)
	}
}

func TestDecodeCastleTokenInputStripsTransportWhitespace(t *testing.T) {
	token := strings.Repeat("castle-token-", 64)
	wrapped := token[:80] + "\r\n" + token[80:160] + "\n\t " + token[160:]
	got, err := decodeCastleTokenInput(wrapped)
	if err != nil {
		t.Fatalf("decodeCastleTokenInput(wrapped) failed: %v", err)
	}
	if got != token {
		t.Fatalf("decodeCastleTokenInput(wrapped) = %q, want original token", got)
	}

	encoded := base64.RawURLEncoding.EncodeToString([]byte(token))
	wrappedEncoded := encoded[:80] + "\n" + encoded[80:]
	got, err = decodeCastleTokenInput(castleTokenHeaderPrefix + wrappedEncoded)
	if err != nil {
		t.Fatalf("decodeCastleTokenInput(wrapped header) failed: %v", err)
	}
	if got != token {
		t.Fatalf("decodeCastleTokenInput(wrapped header) = %q, want original token", got)
	}
}

func TestDecodeCastleTokenBatchInput(t *testing.T) {
	tokensByIndex := make([]string, castleTokenBatchSize)
	input := make(map[string]string, castleTokenBatchSize)
	for index := 1; index <= castleTokenBatchSize; index++ {
		token := strings.Repeat(fmt.Sprintf("castle-%d-", index), 64)
		tokensByIndex[index-1] = token
		input[castleTokenFieldID(index)] = token
	}
	input[loginFieldCastleToken] = tokensByIndex[0][:80] + "\n" + tokensByIndex[0][80:]

	tokens, err := decodeCastleTokenBatchInput(input)
	if err != nil {
		t.Fatalf("decodeCastleTokenBatchInput() error = %v", err)
	}
	if len(tokens) != castleTokenBatchSize {
		t.Fatalf("len(tokens) = %d, want %d", len(tokens), castleTokenBatchSize)
	}
	for index := range tokens {
		if tokens[index] != tokensByIndex[index] {
			t.Fatalf("tokens[%d] = %q, want token %d", index, tokens[index], index+1)
		}
	}
}

func TestCastleTokenFieldPatternAcceptsWrappedTransportValue(t *testing.T) {
	fields := castleTokenCookieFields()
	if len(fields) == 0 || fields[0].ID != loginFieldCastleToken {
		t.Fatalf("first field = %#v, want Castle token field", fields)
	}
	token := strings.Repeat("castle-token-", 64)
	wrapped := token[:80] + "\n" + token[80:]
	if !regexp.MustCompile(fields[0].Pattern).MatchString(wrapped) {
		t.Fatalf("Castle token field pattern %q rejected wrapped token transport", fields[0].Pattern)
	}
	if strings.Contains(fields[0].Pattern, "(?s)") {
		t.Fatalf("Castle token field pattern %q must be JavaScript RegExp compatible", fields[0].Pattern)
	}
	if !strings.Contains(fields[0].Pattern, `[\s\S]`) {
		t.Fatalf("Castle token field pattern %q should match newlines without JS-only-invalid flags", fields[0].Pattern)
	}
}

func TestFindWebLoginAuthMethodMatchesNameOrID(t *testing.T) {
	methods := []twittermeow.WebLoginAuthMethod{
		{ID: "Totp", Name: "Authenticator App", Supported: true},
		{ID: "Sms", Name: "Text Message", Supported: true},
		{ID: "BackupCode", Name: "Backup Code", Supported: true},
	}
	if method, ok := findWebLoginAuthMethod(methods, "Authenticator App"); !ok || method.ID != "Totp" {
		t.Fatalf("find by label = %#v %t, want Totp", method, ok)
	}
	if method, ok := findWebLoginAuthMethod(methods, "backup_code"); !ok || method.ID != "BackupCode" {
		t.Fatalf("find by normalized ID = %#v %t, want BackupCode", method, ok)
	}
	if method, ok := findWebLoginAuthMethod(methods, "text_message"); !ok || method.ID != "Sms" {
		t.Fatalf("find by normalized ID = %#v %t, want Sms", method, ok)
	}
}

func TestMakeVerificationStepUsesPhoneNumberInput(t *testing.T) {
	step := makeVerificationStep(&twittermeow.WebLoginChallenge{
		Description: "Enter the phone number associated with your X account.",
		InputKind:   twittermeow.WebLoginChallengeInputKindPhoneNumber,
	}, "")

	if step.UserInputParams == nil || len(step.UserInputParams.Fields) != 1 {
		t.Fatalf("UserInputParams = %#v, want one field", step.UserInputParams)
	}
	field := step.UserInputParams.Fields[0]
	if field.Type != bridgev2.LoginInputFieldTypePhoneNumber {
		t.Fatalf("field.Type = %s, want phone_number", field.Type)
	}
	if field.Name != "Phone number" {
		t.Fatalf("field.Name = %q, want Phone number", field.Name)
	}
	if !strings.Contains(step.Instructions, "phone number") {
		t.Fatalf("Instructions = %q, want phone number prompt", step.Instructions)
	}
}
