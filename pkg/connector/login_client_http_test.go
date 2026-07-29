package connector

import (
	"bytes"
	"context"
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
	if step.Type != bridgev2.LoginStepTypeClientHTTP || step.StepID != LoginStepIDClientHTTPRequest ||
		step.ClientHTTPParams == nil {
		t.Fatalf("first step = %#v, want typed client HTTP request", step)
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
		step, err = login.SubmitClientHTTPResponse(context.Background(), input)
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
	if step.ClientHTTPParams == nil || !bytes.Equal(step.ClientHTTPParams.Body, pending.Body) {
		t.Fatal("client HTTP step did not carry the exact request body")
	}
	if step.CookiesParams != nil {
		t.Fatal("client HTTP request unexpectedly used a cookie/WebView step")
	}

	input := clientHTTPTestResponseInput(pending, "login accepted")
	clientHTTPTestAuthenticate(input)
	step, err = login.SubmitClientHTTPResponse(context.Background(), input)
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
		clientHTTPTestAuthenticate(input)
		step, err = login.SubmitClientHTTPResponse(context.Background(), input)
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
		step, err = login.SubmitClientHTTPResponse(ctx, input)
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
	step, err = login.SubmitClientHTTPResponse(ctx, input)
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
	step, err = login.SubmitClientHTTPResponse(ctx, input)
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
	clientHTTPTestAuthenticate(input)
	step, err = login.SubmitClientHTTPResponse(ctx, input)
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
		clientHTTPTestAuthenticate(input)
		step, err = login.SubmitClientHTTPResponse(ctx, input)
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

	step, err = login.SubmitClientHTTPResponse(context.Background(), &bridgev2.LoginClientHTTPResponse{
		RequestID: step.ClientHTTPParams.RequestID,
		Error:     "Failed to fetch",
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
		ID:     "client-http-7",
		Method: http.MethodPost,
		URL:    "https://x.com/i/jfapi/onboarding/web/actions/begin_login",
		Headers: http.Header{
			"Authorization": {"Bearer public"},
			"Cookie":        {"guest_id=v1%3A123; ct0=csrf"},
			"Referer":       {endpoints.JETFUEL_LOGIN_REFERER_URL},
			"X-Multi":       {"one", "two"},
		},
		Body: []byte("password=secret"),
	}
	step, err := makeClientHTTPRequestStep(request)
	if err != nil {
		t.Fatal(err)
	}
	if step.Type != bridgev2.LoginStepTypeClientHTTP || step.StepID != LoginStepIDClientHTTPRequest ||
		step.ClientHTTPParams == nil || step.CookiesParams != nil {
		t.Fatalf("step = %#v", step)
	}
	params := step.ClientHTTPParams
	if params.RequestID != request.ID || params.Method != request.Method || params.URL != request.URL {
		t.Fatalf("client HTTP params = %#v", params)
	}
	if !bytes.Equal(params.Body, request.Body) {
		t.Fatalf("client HTTP body = %q", params.Body)
	}
	if got := params.Headers.Values("X-Multi"); len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("multi-value headers = %#v", params.Headers)
	}
	if params.Headers.Get("Cookie") != request.Headers.Get("Cookie") ||
		params.Headers.Get("Referer") != request.Headers.Get("Referer") {
		t.Fatalf("cookies or referrer were not preserved: %#v", params.Headers)
	}
	request.Headers.Set("Cookie", "changed")
	request.Body[0] = 'X'
	if params.Headers.Get("Cookie") == "changed" || params.Body[0] == 'X' {
		t.Fatal("client HTTP step did not clone mutable request data")
	}
}

func clientHTTPTestResponseInput(
	request *twittermeow.ClientHTTPRequest,
	body string,
) *bridgev2.LoginClientHTTPResponse {
	return &bridgev2.LoginClientHTTPResponse{
		RequestID:  request.ID,
		StatusCode: http.StatusOK,
		FinalURL:   request.URL,
		Headers: http.Header{
			"Content-Type": {"application/octet-stream"},
			"Set-Cookie": {
				"guest_id=v1%3A123456789; Domain=.x.com; Path=/",
				"gt=123456789; Domain=.x.com; Path=/",
			},
		},
		Body: []byte(body),
	}
}

func clientHTTPTestAuthenticate(response *bridgev2.LoginClientHTTPResponse) {
	response.Headers.Add("Set-Cookie", "auth_token=authenticated-cookie; Domain=.x.com; Path=/")
	response.Headers.Add("Set-Cookie", "ct0=csrf-cookie; Domain=.x.com; Path=/")
	response.Headers.Add("Set-Cookie", "twid=u%3D123456789; Domain=.x.com; Path=/")
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
