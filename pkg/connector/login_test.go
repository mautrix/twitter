package connector

import (
	"bytes"
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
	"maunium.net/go/mautrix/bridgev2/database"

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

func TestSubmitPINRejectsEmptyUserIDBeforeAPIRequest(t *testing.T) {
	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	requestCount := 0
	transport := connectorRoundTripFunc(func(*http.Request) (*http.Response, error) {
		requestCount++
		return nil, errors.New("unexpected request")
	})
	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		client:             client,
		loginHTTPTransport: transport,
	}

	step, err := login.SubmitUserInput(context.Background(), map[string]string{"pin": "1234"})
	if step != nil {
		t.Fatalf("SubmitUserInput() step = %#v, want nil", step)
	}
	if !errors.Is(err, ErrMissingUserID) {
		t.Fatalf("SubmitUserInput() error = %v, want ErrMissingUserID", err)
	}
	if requestCount != 0 {
		t.Fatalf("proxied API request count = %d, want 0", requestCount)
	}
}

func TestMigrationMissingUserIDFallsBackToCredentials(t *testing.T) {
	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	login := &TwitterLogin{
		User:   &bridgev2.User{Log: zerolog.Nop()},
		client: client,
	}
	override := &bridgev2.UserLogin{UserLogin: &database.UserLogin{
		Metadata: &UserLoginMetadata{Cookies: "auth_token=fake"},
	}}

	step, err := login.startWithOverride(context.Background(), override)
	if err != nil {
		t.Fatalf("startWithOverride() error = %v", err)
	}
	if step != nil && step.Type == bridgev2.LoginStepTypeCookies && step.CookiesParams.Hidden {
		step, err = login.SubmitCookies(context.Background(), map[string]string{loginFieldBrowserUserAgent: "test-native-user-agent"})
		if err != nil {
			t.Fatalf("SubmitCookies() error = %v", err)
		}
	}
	if step == nil || step.StepID != LoginStepIDCredentials {
		t.Fatalf("startWithOverride() step = %#v, want credentials step", step)
	}
	if login.client != nil {
		t.Fatal("startWithOverride() retained invalid PIN client")
	}
	if login.isMigration {
		t.Fatal("startWithOverride() marked missing-ID login as migration")
	}
}

func TestSubmitPINRoutesXChatRequestsThroughLoginHTTPAndRestoresDefault(t *testing.T) {
	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetCurrentUserID("123456789")
	defaultTransport := client.HTTP.Transport
	publicKeysBody, err := json.Marshal(makePublicKeysResponse())
	if err != nil {
		t.Fatalf("marshal public keys response: %v", err)
	}

	seenGetPublicKeys := false
	seenAddPublicKey := false
	transport := connectorRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch {
		case strings.Contains(req.URL.Path, "GetPublicKeys"):
			seenGetPublicKeys = true
			var variables struct {
				IDs []string `json:"ids"`
			}
			if err := json.Unmarshal([]byte(req.URL.Query().Get("variables")), &variables); err != nil {
				t.Errorf("decode GetPublicKeys variables: %v", err)
			} else if len(variables.IDs) != 1 || variables.IDs[0] != "123456789" {
				t.Errorf("GetPublicKeys IDs = %#v, want nonempty current user ID", variables.IDs)
			}
			return connectorTestHTTPResponse(string(publicKeysBody)), nil
		case strings.Contains(req.URL.Path, "AddXChatPublicKey"):
			seenAddPublicKey = true
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"blocked"}]}`)),
			}, nil
		default:
			t.Errorf("unexpected proxied request: %s %s", req.Method, req.URL.String())
			return connectorTestHTTPResponse("{}"), nil
		}
	})
	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		client:             client,
		loginHTTPTransport: transport,
	}

	step, err := login.SubmitUserInput(context.Background(), map[string]string{"pin": "1234"})
	if step != nil {
		t.Fatalf("SubmitUserInput() step = %#v, want nil on X API rejection", step)
	}
	if err == nil || !strings.Contains(err.Error(), "failed to register xchat public key") {
		t.Fatalf("SubmitUserInput() error = %v, want AddXChatPublicKey failure", err)
	}
	if !seenGetPublicKeys || !seenAddPublicKey {
		t.Fatalf("proxied calls: GetPublicKeys=%t AddXChatPublicKey=%t, want both", seenGetPublicKeys, seenAddPublicKey)
	}
	if client.HTTP.Transport != defaultTransport {
		t.Fatal("PIN failure left the login HTTP transport attached to the X client")
	}
	if login.loginHTTPTransport == nil {
		t.Fatal("PIN failure discarded the login HTTP transport needed for retry")
	}
}

func TestSubmitPINClientHTTPFailureIsRetryable(t *testing.T) {
	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetCurrentUserID("123456789")
	defaultTransport := client.HTTP.Transport
	transport := connectorRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("error from client: browser request failed")
	})
	login := &TwitterLogin{
		client:             client,
		loginHTTPTransport: transport,
		needsPINSetup:      true,
	}

	step, err := login.SubmitUserInput(context.Background(), map[string]string{"pin": "1234"})
	if err != nil {
		t.Fatalf("SubmitUserInput() error = %v, want retry step", err)
	}
	if step == nil || step.StepID != LoginStepJuiceboxPIN || !strings.Contains(step.Instructions, clientHTTPFailureInstructions) {
		t.Fatalf("SubmitUserInput() step = %#v, want retryable PIN step", step)
	}
	if login.client != client || login.loginHTTPTransport == nil {
		t.Fatal("retryable client HTTP failure discarded the authenticated login session")
	}
	if client.HTTP.Transport != defaultTransport {
		t.Fatal("retryable client HTTP failure left login transport attached between attempts")
	}
}

func TestFinishLoginHTTPTransportDetachesPersistentClient(t *testing.T) {
	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	defaultTransport := client.HTTP.Transport
	transport := connectorRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return connectorTestHTTPResponse("{}"), nil
	})
	login := &TwitterLogin{client: client, loginHTTPTransport: transport}
	login.attachLoginHTTPTransportForPIN()

	login.finishLoginHTTPTransport()

	if client.HTTP.Transport != defaultTransport {
		t.Fatal("finishLoginHTTPTransport() did not restore the persistent client transport")
	}
	if login.loginHTTPTransport != nil {
		t.Fatal("finishLoginHTTPTransport() retained the login-scoped transport")
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

func TestHandleWebLoginCredentialsErrorRetriesRecoverableErrors(t *testing.T) {
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
	if err != nil {
		t.Fatalf("handleWebLoginCredentialsError(temporary limit) error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDCredentials || !strings.Contains(step.Instructions, "Wait a bit") {
		t.Fatalf("handleWebLoginCredentialsError(temporary limit) step = %#v, want retryable credentials step", step)
	}
}

func TestHandleWebLoginVerificationErrorRetriesOnlyCodeErrors(t *testing.T) {
	challenge := &twittermeow.WebLoginChallenge{Description: "Enter the verification code from X."}
	step, err := handleWebLoginVerificationError(challenge, &twittermeow.WebLoginError{
		Code:    32,
		Message: "The verification code is incorrect",
	})
	if err != nil {
		t.Fatalf("handleWebLoginVerificationError(wrong code) error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDVerification {
		t.Fatalf("handleWebLoginVerificationError(wrong code) step = %#v, want verification step", step)
	}

	step, err = handleWebLoginVerificationError(challenge, &twittermeow.WebLoginError{
		Code:    399,
		Message: "We've temporarily limited your login. Please try again later.",
	})
	if step != nil {
		t.Fatalf("handleWebLoginVerificationError(temporary limit) step = %#v, want nil", step)
	}
	var respErr bridgev2.RespError
	if !errors.As(err, &respErr) || respErr.ErrCode != ErrWebLoginFailed.ErrCode {
		t.Fatalf("handleWebLoginVerificationError(temporary limit) error = %#v, want ErrWebLoginFailed", err)
	}
}

func TestHandleWebLoginAuthMethodErrorRetriesOnlyUnsupportedSelection(t *testing.T) {
	methods := []twittermeow.WebLoginAuthMethod{{ID: "Totp", Name: "Authenticator App", Supported: true}}
	step, err := handleWebLoginAuthMethodError(methods, fmt.Errorf("%w: Security Key", twittermeow.ErrWebLoginUnsupportedAuthMethod))
	if err != nil {
		t.Fatalf("handleWebLoginAuthMethodError(unsupported method) error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDAuthMethod {
		t.Fatalf("handleWebLoginAuthMethodError(unsupported method) step = %#v, want auth method step", step)
	}

	step, err = handleWebLoginAuthMethodError(methods, &twittermeow.WebLoginError{
		Code:    399,
		Message: "We've temporarily limited your login. Please try again later.",
	})
	if step != nil {
		t.Fatalf("handleWebLoginAuthMethodError(temporary limit) step = %#v, want nil", step)
	}
	var respErr bridgev2.RespError
	if !errors.As(err, &respErr) || respErr.ErrCode != ErrWebLoginFailed.ErrCode {
		t.Fatalf("handleWebLoginAuthMethodError(temporary limit) error = %#v, want ErrWebLoginFailed", err)
	}

	step, err = handleWebLoginAuthMethodError(methods, twittermeow.ErrWebLoginMissingAuthMethodState)
	if err != nil {
		t.Fatalf("handleWebLoginAuthMethodError(missing state) error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDCredentials {
		t.Fatalf("handleWebLoginAuthMethodError(missing state) step = %#v, want credentials step", step)
	}
}

func TestHandleWebCastleStageErrorMakesPreludeFailureTerminal(t *testing.T) {
	step, err := handleWebCastleStageError(
		webLoginCastleStageBeginTwoFactor,
		nil,
		nil,
		fmt.Errorf("%w: two-factor prelude response", twittermeow.ErrWebLoginUnexpectedSubtask),
	)
	if step != nil {
		t.Fatalf("handleWebCastleStageError() step = %#v, want nil", step)
	}
	var respErr bridgev2.RespError
	if !errors.As(err, &respErr) || respErr.ErrCode != ErrWebLoginFailed.ErrCode {
		t.Fatalf("handleWebCastleStageError() error = %#v, want ErrWebLoginFailed", err)
	}
}

func TestCreateLoginAcceptsSupportedFlows(t *testing.T) {
	tests := []struct {
		name        string
		flowID      string
		wantType    bridgev2.LoginStepType
		wantStepID  string
		wantInvalid bool
	}{
		{
			name:       "empty flow defaults to cookies",
			wantType:   bridgev2.LoginStepTypeCookies,
			wantStepID: LoginStepIDCookies,
		},
		{
			name:       "cookies",
			flowID:     LoginFlowIDCookies,
			wantType:   bridgev2.LoginStepTypeCookies,
			wantStepID: LoginStepIDCookies,
		},
		{
			name:       "native",
			flowID:     LoginFlowIDPassword,
			wantType:   bridgev2.LoginStepTypeCookies,
			wantStepID: "fi.mau.twitter.login.browser_identity",
		},
		{
			name:        "unknown",
			flowID:      "unknown",
			wantInvalid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			process, err := (&TwitterConnector{}).CreateLogin(context.Background(), nil, test.flowID)
			if test.wantInvalid {
				if !errors.Is(err, bridgev2.ErrInvalidLoginFlowID) {
					t.Fatalf("CreateLogin() error = %v, want ErrInvalidLoginFlowID", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("CreateLogin() error = %v", err)
			}

			step, err := process.Start(context.Background())
			if err != nil {
				t.Fatalf("Start() error = %v", err)
			}
			if step.Type != test.wantType || step.StepID != test.wantStepID {
				t.Fatalf("Start() step = (%s, %s), want (%s, %s)", step.Type, step.StepID, test.wantType, test.wantStepID)
			}
			if test.wantType == bridgev2.LoginStepTypeCookies {
				if step.CookiesParams == nil || step.UserInputParams != nil {
					t.Fatalf("Start() params = cookies %#v, user input %#v", step.CookiesParams, step.UserInputParams)
				}
			} else if step.UserInputParams == nil || step.CookiesParams != nil {
				t.Fatalf("Start() params = user input %#v, cookies %#v", step.UserInputParams, step.CookiesParams)
			}
			if test.flowID == LoginFlowIDPassword {
				step, err = process.(bridgev2.LoginProcessCookies).SubmitCookies(context.Background(), map[string]string{loginFieldBrowserUserAgent: "test-native-user-agent"})
				if err != nil || step == nil || step.StepID != LoginStepIDCredentials {
					t.Fatalf("browser identity result = %#v, %v, want credentials", step, err)
				}
				if got := process.(*TwitterLogin).newLoginClient().GetBrowserHeaders().UserAgent; got != "test-native-user-agent" {
					t.Fatalf("bootstrap user agent = %q, want captured browser identity", got)
				}
			}
		})
	}
}

func TestWebCastleLoginUsesJetfuelDocumentMetadataAndReturnsCode399AsRetry(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_VIEWER_CONTEXT", "0")
	const mainPageHTML = `<html><head><meta name="twitter-site-verification" content="verification-token"></head><body><script>
{"country": "US"}
gt=123456789
</script></body></html>`
	const jetfuelActionResponse = endpoints.JETFUEL_BEGIN_LOGIN_PATH + "\x00username_or_email"
	const jetfuelDocument = `{"responsive_web_castle_public_key":{"value":"test-public-key"}}123:"ondemand.castle",{123:"abcdef"}`

	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	identifierRequestCount := 0
	documentRequestCount := 0
	client.HTTP = &http.Client{Transport: connectorRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodGet && req.URL.String() == endpoints.JETFUEL_LOGIN_DOCUMENT_URL:
			if req.Header.Get("sec-fetch-dest") != "document" {
				t.Fatal("Castle bootstrap request must use document headers")
			}
			documentRequestCount++
			if documentRequestCount == 1 {
				resp := connectorTestHTTPResponse(mainPageHTML)
				resp.Header.Add("Set-Cookie", "guest_id=v1%3A123456789; Path=/; Secure")
				return resp, nil
			}
			resp := connectorTestHTTPResponse(jetfuelDocument)
			resp.Header.Add("Set-Cookie", "ct0=test-csrf; Path=/; Secure")
			return resp, nil
		case req.Method == http.MethodGet && req.URL.Path == endpoints.JETFUEL_LANDING_PATH:
			return connectorTestHTTPResponse("landing"), nil
		case req.Method == http.MethodGet && req.URL.Path == "/onboarding/web" && req.URL.Query().Get("mode") == "login":
			if req.URL.Host != "jf.x.com" {
				t.Fatal("login action graph did not use native host")
			}
			return connectorTestHTTPResponse(jetfuelActionResponse), nil
		case req.Method == http.MethodPost && req.URL.Path == endpoints.JETFUEL_BEGIN_LOGIN_PATH:
			identifierRequestCount++
			if req.URL.Host != "jf.x.com" || req.Header.Get("Referer") != "https://x.com/" || req.Header.Get("sec-fetch-site") != "same-site" || req.Header.Get("Origin") != "https://x.com" {
				t.Fatal("identifier request did not use native host and browser headers")
			}
			if req.Header.Get("x-csrf-token") != "test-csrf" || !strings.Contains(req.Header.Get("Cookie"), "ct0=test-csrf") {
				t.Fatalf("identifier request did not include CSRF cookie and header from Jetfuel document")
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll(identifier request) error = %v", err)
			}
			for _, value := range []string{"username_or_email=test-user", "%24castle_token=identifier-token"} {
				if !strings.Contains(string(body), value) {
					t.Fatalf("identifier request body missing %q", value)
				}
			}
			if strings.Contains(string(body), "password=") {
				t.Fatal("identifier request included password")
			}
			return connectorTestHTTPResponse("We've temporarily limited your login. Please try again later."), nil
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
	if info := client.JetfuelCastleTokenInfo(); info.PublicKey != "test-public-key" ||
		info.ScriptURL != "https://abs.twimg.com/responsive-web/client-web/ondemand.castle.abcdefa.js" {
		t.Fatalf("Castle metadata from Jetfuel document = %#v", info)
	}
	if documentRequestCount != 2 {
		t.Fatalf("document request count = %d, want bootstrap and metadata fallback", documentRequestCount)
	}

	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		webLogin:           session,
		webLoginIdentifier: "test-user",
		webLoginPassword:   "test-password",
	}
	step, err := login.continueStartedCredentialsLogin(context.Background(), result)
	if err != nil {
		t.Fatalf("continueStartedCredentialsLogin() error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDCastleToken {
		t.Fatalf("continueStartedCredentialsLogin() step = %#v, want Castle token step", step)
	}
	if login.webLoginCastleStage != webLoginCastleStageIdentifier {
		t.Fatalf("initial Castle stage = %q, want identifier", login.webLoginCastleStage)
	}

	client.SetNextJetfuelCastleTokens([]string{"identifier-token", "unused-token"})
	var logs bytes.Buffer
	logger := zerolog.New(&logs).With().Str("login_id", "test-login").Logger()
	ctx := logger.WithContext(context.Background())
	step, err = login.continueWebCastleLogin(ctx)
	if err != nil {
		t.Fatalf("continueWebCastleLogin() error = %v", err)
	}
	if step == nil || step.StepID != LoginStepIDCredentials || !strings.Contains(step.Instructions, "Wait a bit") {
		t.Fatalf("continueWebCastleLogin() step = %#v, want retryable credentials step", step)
	}
	if login.webLoginCastleStage != "" {
		t.Fatalf("Castle stage = %q, want cleared", login.webLoginCastleStage)
	}
	if identifierRequestCount != 1 {
		t.Fatalf("identifier request count = %d, want 1", identifierRequestCount)
	}
	if !client.HasNextJetfuelCastleToken() {
		t.Fatal("unused Castle token was consumed after code 399")
	}
	logged := logs.String()
	for _, field := range []string{"\"login_id\":\"test-login\"", "\"stage\":\"identifier\"", "\"error_kind\":\"x_response\"", "\"error_code\":399"} {
		if !strings.Contains(logged, field) {
			t.Fatalf("safe Castle failure log missing %q: %s", field, logged)
		}
	}
	if strings.Contains(logged, "temporarily limited") || strings.Contains(logged, "test-password") || strings.Contains(logged, "identifier-token") {
		t.Fatalf("Castle failure log leaked response or submitted values: %s", logged)
	}
}

func TestBrowserIdentityRejectsInvalidHeadersAndClearsOnCancel(t *testing.T) {
	for _, ua := range []string{"", "bad\r\nheader", strings.Repeat("a", 1025)} {
		login := &TwitterLogin{}
		requests := 0
		transport := connectorRoundTripFunc(func(*http.Request) (*http.Response, error) { requests++; return nil, errors.New("unexpected request") })
		step, err := login.StartWithParams(context.Background(), bridgev2.LoginStartParams{HTTP: transport})
		if err != nil || step.Type != bridgev2.LoginStepTypeCookies || !step.CookiesParams.Hidden {
			t.Fatalf("identity start: %v", err)
		}
		step, err = login.SubmitCookies(context.Background(), map[string]string{loginFieldBrowserUserAgent: ua})
		if err == nil || step != nil || requests != 0 || login.browserHeaders.UserAgent != "" {
			t.Fatal("invalid identity accepted or issued HTTP")
		}
	}
	login := &TwitterLogin{}
	_, _ = login.Start(context.Background())
	login.Cancel()
	if login.waitingForBrowserIdentity || login.webLogin != nil || login.webLoginPassword != "" {
		t.Fatal("cancel retained pending identity or password")
	}
}

func TestCastleBrowserIdentityMismatchStopsBeforePOST(t *testing.T) {
	login := &TwitterLogin{}
	_, err := login.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = login.SubmitCookies(context.Background(), map[string]string{loginFieldBrowserUserAgent: "Mozilla/5.0 Chrome/159.0.0.0"})
	if err != nil {
		t.Fatal(err)
	}
	client := login.newLoginClient()
	requests := 0
	client.HTTP = &http.Client{Transport: connectorRoundTripFunc(func(*http.Request) (*http.Response, error) { requests++; return nil, errors.New("unexpected request") })}
	login.webLogin = twittermeow.NewWebLoginSession(client)
	login.webLoginIdentifier = "synthetic-user"
	login.webLoginPassword = "synthetic-password"
	login.webLoginCastleStage = webLoginCastleStageIdentifier
	step, err := login.SubmitCookies(context.Background(), map[string]string{loginFieldBrowserUserAgent: "Mozilla/5.0 Chrome/140.0.0.0", loginFieldCastleToken: strings.Repeat("A", 128)})
	if err == nil || step != nil || requests != 0 || login.webLoginPassword != "" {
		t.Fatal("mismatched identity accepted or retained credentials")
	}
}

func TestCredentialsBootstrapScriptFailureIsRetryable(t *testing.T) {
	const mainURL = "https://abs.twimg.com/responsive-web/client-web/main.abcdef.js"
	const html = `<meta name="twitter-site-verification" content="verification-token"><script>{"country":"US","responsive_web_castle_public_key":{"value":"test-key"}};gt=123456789;{123:"ondemand.castle"};{123:"abcdef"}</script><script src="` + mainURL + `"></script>`
	failScript := true
	scriptRequests := 0
	transport := connectorRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodGet && req.URL.String() == endpoints.JETFUEL_LOGIN_DOCUMENT_URL:
			resp := connectorTestHTTPResponse(html)
			resp.Header.Add("Set-Cookie", "guest_id=test-guest; Path=/; Secure")
			return resp, nil
		case req.Method == http.MethodGet && req.URL.String() == mainURL:
			scriptRequests++
			if failScript {
				return nil, errors.New("error from client: browser request failed")
			}
			return connectorTestHTTPResponse("script"), nil
		case req.Method == http.MethodGet && req.URL.Host == "jf.x.com":
			return connectorTestHTTPResponse(endpoints.JETFUEL_BEGIN_LOGIN_PATH + "\x00username_or_email"), nil
		default:
			t.Fatal("unexpected request before Castle step")
			return nil, nil
		}
	})
	login := &TwitterLogin{loginHTTPTransport: transport}
	input := map[string]string{loginFieldIdentifier: "test-user", loginFieldPassword: "test-password"}
	step, err := login.SubmitUserInput(context.Background(), input)
	if err != nil || step == nil || step.StepID != LoginStepIDCredentials || !strings.Contains(step.Instructions, clientHTTPFailureInstructions) || scriptRequests != 1 {
		t.Fatalf("script failure did not return credentials retry: step=%v err=%v requests=%d", step, err, scriptRequests)
	}
	if login.webLogin != nil || login.webLoginPassword != "" || login.loginHTTPTransport == nil {
		t.Fatal("retry did not clear failed bootstrap while preserving client transport")
	}
	failScript = false
	step, err = login.SubmitUserInput(context.Background(), input)
	if err != nil || step == nil || step.StepID != LoginStepIDCastleToken || scriptRequests != 2 {
		t.Fatalf("credentials retry did not reach Castle: step=%v err=%v requests=%d", step, err, scriptRequests)
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
		case req.Method == http.MethodGet && req.URL.String() == endpoints.JETFUEL_LOGIN_DOCUMENT_URL:
			resp := connectorTestHTTPResponse(mainPageHTML)
			resp.Header.Add("Set-Cookie", "guest_id=v1%3A123456789; Path=/; Secure")
			return resp, nil
		case req.Method == http.MethodGet && req.URL.Path == endpoints.JETFUEL_LANDING_PATH:
			return connectorTestHTTPResponse("landing"), nil
		case req.Method == http.MethodGet && req.URL.Path == "/onboarding/web" && req.URL.Query().Get("mode") == "login":
			return connectorTestHTTPResponse(endpoints.JETFUEL_BEGIN_LOGIN_PATH + "\x00username_or_email"), nil
		case req.Method == http.MethodPost && req.URL.Path == endpoints.JETFUEL_BEGIN_LOGIN_PATH:
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("ReadAll(identifier request) error = %v", err)
			}
			if !strings.Contains(string(body), "%24castle_token=identifier-token") {
				t.Fatalf("identifier request body = %q", body)
			}
			if !strings.Contains(string(body), "username_or_email=test-user") || strings.Contains(string(body), "password=") {
				t.Fatal("identifier request must include username without password")
			}
			return connectorTestHTTPResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
		case req.Method == http.MethodPost && req.URL.Path == endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH:
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

	login := &TwitterLogin{
		User:               &bridgev2.User{Log: zerolog.Nop()},
		webLogin:           session,
		webLoginIdentifier: "test-user",
		webLoginPassword:   "test-password",
	}
	step, err := login.continueStartedCredentialsLogin(context.Background(), result)
	if err != nil || step == nil || step.StepID != LoginStepIDCastleToken || login.webLoginCastleStage != webLoginCastleStageIdentifier {
		t.Fatalf("initial step = %#v, error = %v, want identifier Castle step", step, err)
	}
	client.SetNextJetfuelCastleTokens([]string{"identifier-token"})
	step, err = login.continueWebCastleLogin(context.Background())
	if err != nil || step == nil || step.StepID != LoginStepIDCastleToken || login.webLoginCastleStage != webLoginCastleStagePassword {
		t.Fatalf("identifier result = %#v, error = %v, want password Castle step", step, err)
	}
	client.SetNextJetfuelCastleTokens([]string{"first-password-token"})
	step, err = login.continueWebCastleLogin(context.Background())
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
