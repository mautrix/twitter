package connector

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

const clientHTTPTestMainPage = `<html><head><meta name="twitter-site-verification" content="verification-token"></head><body><script>
{"country": "US", "responsive_web_castle_public_key":{"value":"test-public-key"}}
gt=123456789
123:"ondemand.castle",{123:"abcdef"}
</script></body></html>`

func TestClientHTTPLoginUsesCastleWebviewThenCapturesLocalRequest(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_VIEWER_CONTEXT", "0")
	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		useClientHTTPLogin: true,
	}

	step, err := login.submitCredentialsInput(context.Background(), map[string]string{
		loginFieldIdentifier: "test-user",
		loginFieldPassword:   "test-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	if step.StepID != LoginStepIDClientHTTPRequest || !step.CookiesParams.Hidden {
		t.Fatalf("first step = %#v, want hidden client HTTP request", step)
	}

	for range 20 {
		if step.StepID == LoginStepIDCastleToken {
			break
		}
		pending := login.clientHTTPTransport.PendingRequest()
		if pending == nil {
			t.Fatal("client HTTP request is not pending")
		}
		var body string
		switch {
		case pending.URL == endpoints.JETFUEL_LOGIN_REFERER_URL:
			body = clientHTTPTestMainPage
		case strings.Contains(pending.URL, "/i/jfapi"+endpoints.JETFUEL_LANDING_PATH):
			body = "landing"
		case strings.Contains(pending.URL, "/i/jfapi/onboarding/web?mode=login"):
			body = endpoints.JETFUEL_BEGIN_LOGIN_PATH + "\x00username_or_email"
		default:
			body = "{}"
		}
		input := clientHTTPTestResponseInput(pending, body)
		input["guest_id"] = "v1%3A123456789"
		input["gt"] = "123456789"
		step, err = login.submitClientHTTPInput(context.Background(), input)
		if err != nil {
			t.Fatalf("submit response for %s: %v", pending.URL, err)
		}
		if step.StepID != LoginStepIDClientHTTPRequest && step.StepID != LoginStepIDCastleToken {
			t.Fatalf("step after %s = %#v", pending.URL, step)
		}
	}

	if step.StepID != LoginStepIDCastleToken || step.CookiesParams == nil || !step.CookiesParams.Hidden {
		t.Fatalf("step = %#v, want native hidden Castle token step", step)
	}
	step, err = login.submitWebCastleTokenInput(context.Background(), clientHTTPTestCastleInput())
	if err != nil {
		t.Fatal(err)
	}
	pending := login.clientHTTPTransport.PendingRequest()
	if step.StepID != LoginStepIDClientHTTPRequest || pending == nil {
		t.Fatalf("step = %#v pending = %#v, want direct client request", step, pending)
	}
	form, err := url.ParseQuery(string(pending.Body))
	if err != nil {
		t.Fatal(err)
	}
	if form.Get("username_or_email") != "test-user" || form.Get("password") != "test-password" {
		t.Fatalf("combined form did not preserve credentials")
	}
	if form.Get("$castle_token") != clientHTTPTestCastleToken(1) {
		t.Fatalf("Castle form value = %q", form.Get("$castle_token"))
	}
	if strings.Contains(step.CookiesParams.ExtractJS, "test-password") {
		t.Fatal("client HTTP request envelope contains the plaintext password")
	}
	if strings.Contains(step.CookiesParams.ExtractJS, "createRequestToken") ||
		strings.Contains(step.CookiesParams.ExtractJS, "fetch(") {
		t.Fatal("client HTTP request envelope contains WebView networking code")
	}

	input := clientHTTPTestResponseInput(pending, "login accepted")
	input["guest_id"] = "v1%3A123456789"
	input["gt"] = "123456789"
	input["auth_token"] = "authenticated-cookie"
	input["ct0"] = "csrf-cookie"
	input["twid"] = "u%3D123456789"
	step, err = login.submitClientHTTPInput(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	pending = login.clientHTTPTransport.PendingRequest()
	if step.StepID != LoginStepIDClientHTTPRequest || pending == nil ||
		pending.URL != endpoints.BASE_MESSAGES_URL {
		t.Fatalf("post-login step = %#v pending = %#v, want client-side messages request", step, pending)
	}

	for range 30 {
		if step.StepID == LoginStepIDCastleToken {
			step, err = login.submitWebCastleTokenInput(context.Background(), clientHTTPTestCastleInput())
			if err != nil {
				t.Fatal(err)
			}
		}
		if step.StepID != LoginStepIDClientHTTPRequest {
			break
		}
		pending = login.clientHTTPTransport.PendingRequest()
		if pending == nil {
			break
		}
		body := "{}"
		if pending.URL == endpoints.BASE_MESSAGES_URL {
			body = clientHTTPTestMainPage
		}
		input = clientHTTPTestResponseInput(pending, body)
		input["guest_id"] = "v1%3A123456789"
		input["gt"] = "123456789"
		input["auth_token"] = "authenticated-cookie"
		input["ct0"] = "csrf-cookie"
		input["twid"] = "u%3D123456789"
		step, err = login.submitClientHTTPInput(context.Background(), input)
		if err != nil {
			t.Fatalf("submit post-login response for %s: %v", pending.URL, err)
		}
		if step.StepID == LoginStepIDCastleToken {
			continue
		}
		if step.StepID != LoginStepIDClientHTTPRequest {
			break
		}
	}
	if step.StepID != LoginStepJuiceboxPIN {
		t.Fatalf("final step = %#v, want PIN step", step)
	}
	if login.clientHTTPTransport != nil || login.webLogin.Client().IsClientHTTPEnabled() {
		t.Fatal("client HTTP transport remained enabled after login completion")
	}
}

func TestClientHTTPLoginPreservesVerificationCodeDuringPostLoginRequests(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_VIEWER_CONTEXT", "0")
	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		useClientHTTPLogin: true,
	}
	ctx := context.Background()

	step, err := login.submitCredentialsInput(ctx, map[string]string{
		loginFieldIdentifier: "test-user",
		loginFieldPassword:   "test-password",
	})
	if err != nil {
		t.Fatal(err)
	}

	for range 20 {
		if step.StepID == LoginStepIDCastleToken {
			break
		}
		pending := login.clientHTTPTransport.PendingRequest()
		if pending == nil {
			t.Fatal("client HTTP request is not pending")
		}
		body := "{}"
		switch {
		case pending.URL == endpoints.JETFUEL_LOGIN_REFERER_URL:
			body = clientHTTPTestMainPage
		case strings.Contains(pending.URL, "/i/jfapi"+endpoints.JETFUEL_LANDING_PATH):
			body = "landing"
		case strings.Contains(pending.URL, "/i/jfapi/onboarding/web?mode=login"):
			body = endpoints.JETFUEL_BEGIN_LOGIN_PATH + "\x00username_or_email"
		}
		input := clientHTTPTestResponseInput(pending, body)
		input["guest_id"] = "v1%3A123456789"
		input["gt"] = "123456789"
		step, err = login.submitClientHTTPInput(ctx, input)
		if err != nil {
			t.Fatalf("submit bootstrap response for %s: %v", pending.URL, err)
		}
	}

	if step.StepID != LoginStepIDCastleToken {
		t.Fatalf("step = %#v, want native Castle token step", step)
	}
	step, err = login.submitWebCastleTokenInput(ctx, clientHTTPTestCastleInput())
	if err != nil {
		t.Fatal(err)
	}
	pending := login.clientHTTPTransport.PendingRequest()
	if step.StepID != LoginStepIDClientHTTPRequest || pending == nil {
		t.Fatalf("step = %#v pending = %#v, want credentials request", step, pending)
	}
	chooserBody := "Select a method to authenticate\x00Choose the method you prefer to use for 2-step verification.\x00" +
		"two_factor_method\x00Totp\x00BackupCode\x00U2fSecurityKey\x00" +
		"user_id\x001127993589949243392\x00" +
		"session_token\x0012345678-1234-1234-1234-123456789abc\x00" +
		"prelude_dispatch_id\x00abcdefab-1234-1234-1234-abcdefabcdef\x00" +
		"begin_two_factor_auth"
	input := clientHTTPTestResponseInput(pending, chooserBody)
	input["guest_id"] = "v1%3A123456789"
	input["gt"] = "123456789"
	step, err = login.submitClientHTTPInput(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if step.StepID != LoginStepIDAuthMethod {
		t.Fatalf("credentials response step = %#v, want auth method", step)
	}

	step, err = login.submitWebAuthMethodInput(ctx, map[string]string{
		loginFieldAuthMethod: "Totp",
	})
	if err != nil {
		t.Fatal(err)
	}
	pending = login.clientHTTPTransport.PendingRequest()
	if step.StepID != LoginStepIDClientHTTPRequest || pending == nil {
		t.Fatalf("auth method submit step = %#v pending = %#v", step, pending)
	}
	challengeBody := "Enter the code from your authentication app.\x00challenge_response\x00" +
		endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH + "\x00" +
		"session_token\x0012345678-1234-1234-1234-123456789abc"
	input = clientHTTPTestResponseInput(pending, challengeBody)
	input["guest_id"] = "v1%3A123456789"
	input["gt"] = "123456789"
	step, err = login.submitClientHTTPInput(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if step.StepID != LoginStepIDVerification {
		t.Fatalf("auth method response step = %#v, want verification", step)
	}

	const verificationCode = "123456"
	step, err = login.submitWebVerificationInput(ctx, map[string]string{
		loginFieldVerificationCode: verificationCode,
	})
	if err != nil {
		t.Fatal(err)
	}
	pending = login.clientHTTPTransport.PendingRequest()
	if step.StepID != LoginStepIDClientHTTPRequest || pending == nil {
		t.Fatalf("verification submit step = %#v pending = %#v", step, pending)
	}
	input = clientHTTPTestResponseInput(pending, "/home")
	input["guest_id"] = "v1%3A123456789"
	input["gt"] = "123456789"
	input["auth_token"] = "authenticated-cookie"
	input["ct0"] = "csrf-cookie"
	input["twid"] = "u%3D123456789"
	step, err = login.submitClientHTTPInput(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	pending = login.clientHTTPTransport.PendingRequest()
	if step.StepID != LoginStepIDClientHTTPRequest || pending == nil ||
		pending.URL != endpoints.BASE_MESSAGES_URL {
		t.Fatalf("post-verification step = %#v pending = %#v, want messages request", step, pending)
	}
	if login.webLoginText != verificationCode {
		t.Fatalf("verification code was cleared while its client HTTP operation is still pending")
	}

	for range 30 {
		if step.StepID == LoginStepIDCastleToken {
			step, err = login.submitWebCastleTokenInput(ctx, clientHTTPTestCastleInput())
			if err != nil {
				t.Fatal(err)
			}
		}
		if step.StepID != LoginStepIDClientHTTPRequest {
			break
		}
		pending = login.clientHTTPTransport.PendingRequest()
		if pending == nil {
			break
		}
		body := "{}"
		if pending.URL == endpoints.BASE_MESSAGES_URL {
			body = clientHTTPTestMainPage
		}
		input = clientHTTPTestResponseInput(pending, body)
		input["guest_id"] = "v1%3A123456789"
		input["gt"] = "123456789"
		input["auth_token"] = "authenticated-cookie"
		input["ct0"] = "csrf-cookie"
		input["twid"] = "u%3D123456789"
		step, err = login.submitClientHTTPInput(ctx, input)
		if err != nil {
			t.Fatalf("submit post-login response for %s: %v", pending.URL, err)
		}
		if step.StepID == LoginStepIDCastleToken {
			continue
		}
		if step.StepID != LoginStepIDClientHTTPRequest {
			break
		}
	}
	if step.StepID != LoginStepJuiceboxPIN {
		t.Fatalf("final step = %#v, want PIN step", step)
	}
	if login.webLoginText != "" {
		t.Fatal("verification code remained after its client HTTP operation completed")
	}
}

func TestClientHTTPRequestFailureStopsAutomaticRetry(t *testing.T) {
	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		useClientHTTPLogin: true,
	}
	step, err := login.submitCredentialsInput(context.Background(), map[string]string{
		loginFieldIdentifier: "test-user",
		loginFieldPassword:   "test-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	if step.StepID != LoginStepIDClientHTTPRequest {
		t.Fatalf("initial step = %#v", step)
	}

	step, err = login.submitClientHTTPInput(context.Background(), map[string]string{
		loginFieldClientHTTPError: "Failed to fetch",
	})
	if err != nil {
		t.Fatal(err)
	}
	if step.StepID != LoginStepIDCredentials {
		t.Fatalf("failure step = %#v, want credentials step", step)
	}
	if !strings.Contains(step.Instructions, "Please try again") {
		t.Fatalf("failure instructions = %q", step.Instructions)
	}
	if login.clientHTTPTransport != nil || login.webLogin != nil ||
		login.webLoginIdentifier != "" || login.webLoginPassword != "" {
		t.Fatal("client HTTP failure retained retryable login state")
	}
}

func TestClientHTTPRequestStepContract(t *testing.T) {
	request := &twittermeow.ClientHTTPRequest{
		ID:           "client-http-7",
		Method:       http.MethodPost,
		URL:          "https://x.com/i/jfapi/onboarding/web/actions/begin_login",
		Referrer:     endpoints.JETFUEL_LOGIN_REFERER_URL,
		Headers:      map[string]string{"authorization": "Bearer public"},
		CookieHeader: "guest_id=v1%3A123; ct0=csrf",
		Body:         []byte("password=secret"),
	}
	step, err := makeClientHTTPRequestStep(request)
	if err != nil {
		t.Fatal(err)
	}
	if step.Type != bridgev2.LoginStepTypeCookies || step.StepID != LoginStepIDClientHTTPRequest ||
		step.CookiesParams == nil || !step.CookiesParams.Hidden {
		t.Fatalf("step = %#v", step)
	}
	if step.CookiesParams.URL != clientHTTPDispatchURL ||
		!strings.Contains(step.CookiesParams.WaitForURLPattern, "beeper-client-http") {
		t.Fatalf("client HTTP step does not carry its client-dispatch marker: %#v", step.CookiesParams)
	}
	fields := make(map[string]bridgev2.LoginCookieField)
	for _, field := range step.CookiesParams.Fields {
		fields[field.ID] = field
	}
	for _, fieldID := range []string{
		loginFieldClientHTTPRequestID,
		loginFieldClientHTTPStatus,
		loginFieldClientHTTPResponseHeaders,
		loginFieldClientHTTPResponseBody,
		loginFieldBrowserUserAgent,
		"auth_token",
		"ct0",
	} {
		if _, ok := fields[fieldID]; !ok {
			t.Fatalf("client HTTP fields missing %q", fieldID)
		}
	}
	if !fields[loginFieldClientHTTPRequestID].Required || !fields[loginFieldBrowserUserAgent].Required {
		t.Fatal("request ID and browser user agent must be required")
	}
	if !strings.HasPrefix(step.CookiesParams.ExtractJS, clientHTTPConfigMarkerPrefix) {
		t.Fatal("client HTTP request envelope is missing its protocol marker")
	}
	if strings.Contains(step.CookiesParams.ExtractJS, "password=secret") {
		t.Fatal("client HTTP request envelope contains plaintext request body")
	}
	if strings.Contains(step.CookiesParams.ExtractJS, "fetch(") ||
		strings.Contains(step.CookiesParams.ExtractJS, "createRequestToken") {
		t.Fatal("client HTTP request envelope contains WebView networking code")
	}
	config := decodeClientHTTPTestConfig(t, step.CookiesParams.ExtractJS)
	if config.RequestID != request.ID || config.RequestURL != request.URL ||
		config.CookieHeader != request.CookieHeader ||
		config.Body != base64.RawURLEncoding.EncodeToString(request.Body) {
		t.Fatalf("decoded client HTTP config = %#v", config)
	}
}

func TestClientHTTPResponseDecodeRejectsMalformedData(t *testing.T) {
	valid := clientHTTPTestResponseInput(&twittermeow.ClientHTTPRequest{
		ID:  "client-http-1",
		URL: "https://x.com/login",
	}, "response")
	if _, err := decodeClientHTTPResponse(valid); err != nil {
		t.Fatalf("valid response rejected: %v", err)
	}
	for _, mutate := range []func(map[string]string){
		func(input map[string]string) { delete(input, loginFieldClientHTTPRequestID) },
		func(input map[string]string) { input[loginFieldClientHTTPStatus] = "not-a-status" },
		func(input map[string]string) { input[loginFieldClientHTTPResponseBody] = "bad" },
		func(input map[string]string) { input[loginFieldClientHTTPResponseHeaders] = "v1.invalid" },
		func(input map[string]string) {
			headers, _ := json.Marshal(map[string]string{"bad header": "value"})
			input[loginFieldClientHTTPResponseHeaders] = clientHTTPEncodedFieldPrefix +
				base64.RawURLEncoding.EncodeToString(headers)
		},
	} {
		input := make(map[string]string, len(valid))
		for key, value := range valid {
			input[key] = value
		}
		mutate(input)
		if _, err := decodeClientHTTPResponse(input); err == nil {
			t.Fatalf("malformed response accepted: %#v", input)
		}
	}
}

func clientHTTPTestResponseInput(request *twittermeow.ClientHTTPRequest, body string) map[string]string {
	headers, _ := json.Marshal(map[string]string{"content-type": "application/octet-stream"})
	return map[string]string{
		loginFieldClientHTTPRequestID:       request.ID,
		loginFieldClientHTTPStatus:          "200",
		loginFieldClientHTTPResponseURL:     request.URL,
		loginFieldClientHTTPResponseHeaders: clientHTTPEncodedFieldPrefix + base64.RawURLEncoding.EncodeToString(headers),
		loginFieldClientHTTPResponseBody:    clientHTTPEncodedFieldPrefix + base64.RawURLEncoding.EncodeToString([]byte(body)),
		loginFieldBrowserUserAgent:          "Mozilla/5.0 test client",
		loginFieldBrowserSecCHUA:            `"Chromium";v="150"`,
		loginFieldBrowserPlatform:           `"Windows"`,
		loginFieldBrowserMobile:             "?0",
	}
}

func clientHTTPTestCastleToken(index int) string {
	return strings.Repeat("client-http-castle-"+string(rune('0'+index))+"-", 12)
}

func clientHTTPTestCastleInput() map[string]string {
	input := map[string]string{
		loginFieldBrowserUserAgent: "Mozilla/5.0 test client",
	}
	for index := 1; index <= castleTokenBatchSize; index++ {
		input[castleTokenFieldID(index)] = clientHTTPTestCastleToken(index)
	}
	return input
}

func decodeClientHTTPTestConfig(t *testing.T, extractJS string) clientHTTPConfig {
	t.Helper()
	if !strings.HasPrefix(extractJS, clientHTTPConfigMarkerPrefix) {
		t.Fatal("client HTTP config marker prefix is missing")
	}
	end := strings.Index(extractJS, clientHTTPConfigMarkerSuffix)
	if end < 0 {
		t.Fatal("client HTTP config marker suffix is missing")
	}
	encoded := strings.TrimPrefix(extractJS[:end], clientHTTPConfigMarkerPrefix)
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode client HTTP config: %v", err)
	}
	var config clientHTTPConfig
	if err = json.Unmarshal(raw, &config); err != nil {
		t.Fatalf("unmarshal client HTTP config: %v", err)
	}
	return config
}
