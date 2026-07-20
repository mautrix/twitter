package twittermeow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (rtf roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rtf(req)
}

func readJetfuelTestForm(t *testing.T, req *http.Request) url.Values {
	t.Helper()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("ReadAll(request body) error = %v", err)
	}
	form, err := url.ParseQuery(string(body))
	if err != nil {
		t.Fatalf("ParseQuery(request body) error = %v", err)
	}
	return form
}

func jetfuelTestResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestSettingsListIdentifierPayloadShape(t *testing.T) {
	payload := onboardingTaskRequest{
		FlowToken: "flow-token",
		SubtaskInputs: []onboardingSubtaskInput{{
			SubtaskID: webLoginSubtaskIdentifier,
			SettingsList: &settingsListInput{
				SettingResponses: []settingResponseInput{{
					Key: "user_identifier",
					ResponseData: map[string]resultInput{
						"text_data": {Result: "example"},
					},
				}},
				Link: webLoginLinkNext,
			},
		}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var got map[string]any
	if err = json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	inputs := got["subtask_inputs"].([]any)
	settingsList := inputs[0].(map[string]any)["settings_list"].(map[string]any)
	responses := settingsList["setting_responses"].([]any)
	firstResponse := responses[0].(map[string]any)
	responseData := firstResponse["response_data"].(map[string]any)
	textData := responseData["text_data"].(map[string]any)

	if firstResponse["key"] != "user_identifier" {
		t.Fatalf("setting response key = %v, want user_identifier", firstResponse["key"])
	}
	if textData["result"] != "example" {
		t.Fatalf("text_data.result = %v, want example", textData["result"])
	}
	if settingsList["link"] != webLoginLinkNext {
		t.Fatalf("settings_list.link = %v, want %s", settingsList["link"], webLoginLinkNext)
	}
	if _, ok := settingsList["castle_token"]; ok {
		t.Fatalf("castle_token should be omitted when empty")
	}
}

func TestWebLoginStartPayloadShape(t *testing.T) {
	payload := newWebLoginStartPayload("US")
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var got map[string]any
	if err = json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	inputFlowData := got["input_flow_data"].(map[string]any)
	flowContext := inputFlowData["flow_context"].(map[string]any)
	startLocation := flowContext["start_location"].(map[string]any)
	subtaskVersions := got["subtask_versions"].(map[string]any)

	if inputFlowData["country_code"] != "US" {
		t.Fatalf("country_code = %v, want US", inputFlowData["country_code"])
	}
	if startLocation["location"] != "manual_link" {
		t.Fatalf("start_location.location = %v, want manual_link", startLocation["location"])
	}
	if subtaskVersions["settings_list"] != float64(7) {
		t.Fatalf("settings_list version = %v, want 7", subtaskVersions["settings_list"])
	}
}

func TestWebLoginResultClassifiesTwoFactor(t *testing.T) {
	session := &WebLoginSession{
		subtasks: []onboardingSubtask{{
			SubtaskID: webLoginSubtaskTwoFactor,
			EnterText: &enterTextSubtask{
				HintText: "Enter your code",
				NextLink: &navigationLink{LinkID: webLoginLinkNext},
			},
		}},
	}
	result := session.result()
	if result.Status != WebLoginStatusNeedsText {
		t.Fatalf("Status = %s, want %s", result.Status, WebLoginStatusNeedsText)
	}
	if result.Challenge == nil || !result.Challenge.IsTwoFactor {
		t.Fatalf("Challenge = %#v, want two-factor challenge", result.Challenge)
	}
	if result.Challenge.Hint != "Enter your code" {
		t.Fatalf("Challenge.Hint = %q", result.Challenge.Hint)
	}
}

func TestParseJetfuelLoginResponseFindsActionsAndFields(t *testing.T) {
	body := []byte{
		0x03, 0x00, 0xff,
	}
	body = append(body, []byte("username_or_email\x00/onboarding/web/actions/begin_login\x00password\x00/onboarding/web/actions/login_enter_password\x00session_token\x0012345678-1234-1234-1234-123456789abc")...)

	parsed := parseJetfuelLoginResponse(body)
	if !parsed.hasField("username_or_email") {
		t.Fatalf("username_or_email field not found in %#v", parsed.fields)
	}
	if !parsed.hasPath("/onboarding/web/actions/begin_login") {
		t.Fatalf("begin_login path not found in %#v", parsed.paths)
	}
	if action := parsed.passwordAction(); action != endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
		t.Fatalf("passwordAction() = %q", action)
	}
	if token := parsed.uuidValue("session_token"); token != "12345678-1234-1234-1234-123456789abc" {
		t.Fatalf("uuidValue(session_token) = %q", token)
	}
}

func TestJetfuelLoginResponseSeparatesTwoFactorPreludeFromCodeAction(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("prelude_dispatch_id\x00aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee\x00/onboarding/web/actions/begin_two_factor_auth\x00/onboarding/web/actions/two_factor_code"))

	if action := parsed.beginTwoFactorAction(); action != endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("beginTwoFactorAction() = %q", action)
	}
	if action := parsed.verificationAction(); action != "/onboarding/web/actions/two_factor_code" {
		t.Fatalf("verificationAction() = %q", action)
	}
	if id := parsed.uuidValue("prelude_dispatch_id"); id != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Fatalf("uuidValue(prelude_dispatch_id) = %q", id)
	}
}

func TestJetfuelLoginResponseExpandsBareTwoFactorActions(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("prelude_dispatch_id\x00aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee\x00begin_two_factor_auth\x00finish_two_factor_auth\x00challenge_response"))

	if action := parsed.beginTwoFactorAction(); action != endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("beginTwoFactorAction() = %q", action)
	}
	if action := parsed.verificationAction(); action != "/onboarding/web/actions/finish_two_factor_auth" {
		t.Fatalf("verificationAction() = %q", action)
	}
}

func TestJetfuelLoginResponseFindsSMSPhoneAction(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte(
		"Confirm the phone number associated with your X account.\x00phone_number\x00/onboarding/web/actions/send_sms_code\x00session_token",
	))

	method := WebLoginAuthMethod{Kind: WebLoginAuthMethodKindSMS}
	if action := parsed.verificationActionForMethod(method); action != "/onboarding/web/actions/send_sms_code" {
		t.Fatalf("verificationActionForMethod(SMS) = %q, want send_sms_code action", action)
	}
}

func TestJetfuelLoginResponseFindsSMSCodeAction(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte(
		"We sent a text message with a verification code to your phone.\x00challenge_response\x00/onboarding/web/actions/enter_sms_pin\x00session_token",
	))

	method := WebLoginAuthMethod{Kind: WebLoginAuthMethodKindSMS}
	if action := parsed.verificationActionForMethod(method); action != "/onboarding/web/actions/enter_sms_pin" {
		t.Fatalf("verificationActionForMethod(SMS) = %q, want enter_sms_pin action", action)
	}
}

func TestJetfuelLoginResponseFindsTwoFactorCodeFields(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("authentication code\x00session_token\x00challenge_response\x00verification_code\x00two_factor_code\x00prelude_dispatch_id"))
	fields := parsed.verificationCodeFields()
	want := []string{"challenge_response", "verification_code", "two_factor_code"}
	if strings.Join(fields, ",") != strings.Join(want, ",") {
		t.Fatalf("verificationCodeFields() = %#v, want %#v", fields, want)
	}
}

func TestJetfuelLoginResponseFindsBackupCodeFields(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Enter a backup code from X.\x00backup_code\x00challenge_response\x00session_token\x00castle_token"))
	fields := parsed.verificationCodeFields()
	want := []string{"backup_code", "challenge_response"}
	if strings.Join(fields, ",") != strings.Join(want, ",") {
		t.Fatalf("verificationCodeFields() = %#v, want %#v", fields, want)
	}
}

func TestJetfuelLoginResponseFindsPhoneNumberFields(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Enter the phone number associated with your X account.\x00phone_number\x00challenge_response\x00session_token"))
	fields := parsed.verificationCodeFields()
	want := []string{"phone_number"}
	if strings.Join(fields, ",") != strings.Join(want, ",") {
		t.Fatalf("verificationCodeFields() = %#v, want %#v", fields, want)
	}
}

func TestJetfuelLoginResponseBuildsPhoneNumberChallenge(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Confirm the phone number associated with your X account.\x00phone_number\x00session_token"))
	challenge := parsed.verificationChallenge()
	if challenge.Description != "Enter the phone number associated with your X account." {
		t.Fatalf("verificationChallenge().Description = %q", challenge.Description)
	}
	if challenge.Hint != "Phone number" {
		t.Fatalf("verificationChallenge().Hint = %q", challenge.Hint)
	}
	if challenge.InputKind != WebLoginChallengeInputKindPhoneNumber {
		t.Fatalf("verificationChallenge().InputKind = %q, want phone_number", challenge.InputKind)
	}
	if challenge.IsTwoFactor {
		t.Fatal("verificationChallenge().IsTwoFactor = true, want false for phone number input")
	}
}

func TestJetfuelLoginResponseKeepsSMSCodeAsTwoFactorCode(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("We sent a text message with a verification code to your phone.\x00challenge_response\x00session_token"))
	challenge := parsed.verificationChallenge()
	if challenge.InputKind != WebLoginChallengeInputKindCode {
		t.Fatalf("verificationChallenge().InputKind = %q, want code", challenge.InputKind)
	}
	if challenge.Description != "Enter the code sent to your phone number." {
		t.Fatalf("verificationChallenge().Description = %q", challenge.Description)
	}
	if !challenge.IsTwoFactor {
		t.Fatal("verificationChallenge().IsTwoFactor = false, want true")
	}
}

func TestJetfuelLoginResponseUsesAuthenticatorAppCopyForBackupCodePrompt(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Enter a backup code from X.\x00challenge_response\x00Totp\x00BackupCode\x00session_token"))
	challenge := parsed.verificationChallenge()
	if challenge.Description != "Enter the code from your authentication app." {
		t.Fatalf("verificationChallenge().Description = %q", challenge.Description)
	}
	if !challenge.IsTwoFactor {
		t.Fatal("verificationChallenge().IsTwoFactor = false, want true")
	}
}

func TestJetfuelLoginResponseDoesNotTreatChallengeModesAsFields(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Enter a backup code from X.\x00challenge_response\x00Totp\x00BackupCode\x00session_token"))
	fields := parsed.verificationCodeFields()
	want := []string{"challenge_response"}
	if strings.Join(fields, ",") != strings.Join(want, ",") {
		t.Fatalf("verificationCodeFields() = %#v, want %#v", fields, want)
	}
}

func TestJetfuelLoginResponseFindsAuthMethodChoice(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte(
		"Select a method to authenticate\x00Choose the method you prefer to use for 2-step verification.\x00two_factor_method\x00" +
			"Totp\x00Sms\x00BackupCode\x00U2fSecurityKey\x00user_id\x001127993589949243392\x00session_token\x0012345678-1234-1234-1234-123456789abc\x00begin_two_factor_auth",
	))
	methods := parsed.authMethods()
	if len(methods) != 4 {
		t.Fatalf("authMethods() length = %d, want 4: %#v", len(methods), methods)
	}
	wantIDs := []string{"Totp", "Sms", "BackupCode", "U2fSecurityKey"}
	for i, want := range wantIDs {
		if methods[i].ID != want || methods[i].Index != i {
			t.Fatalf("method[%d] = %#v, want ID %s index %d", i, methods[i], want, i)
		}
	}
	if !methods[0].Supported || methods[1].Supported || !methods[2].Supported {
		t.Fatalf("only authenticator app and backup code should be supported: %#v", methods)
	}
	if methods[1].Kind != WebLoginAuthMethodKindSMS {
		t.Fatalf("sms method = %#v, want SMS method", methods[1])
	}
	if methods[3].Supported || methods[3].Kind != WebLoginAuthMethodKindUnknown {
		t.Fatalf("security key method = %#v, want known unsupported method", methods[3])
	}
	supported := supportedWebLoginAuthMethods(methods)
	if len(supported) != 2 {
		t.Fatalf("supportedWebLoginAuthMethods() length = %d, want 2: %#v", len(supported), supported)
	}
	if got := parsed.numericValue("user_id"); got != "1127993589949243392" {
		t.Fatalf("numericValue(user_id) = %q", got)
	}
	if action := parsed.beginTwoFactorAction(); action != endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("beginTwoFactorAction() = %q", action)
	}
}

func TestJetfuelLoginResponseFindsBareAuthMethodChoice(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte(
		"Totp\x00Text\x00BackupCode\x00U2fSecurityKey\x00user_id\x001127993589949243392\x00begin_two_factor_auth",
	))
	session := &WebLoginSession{jetfuel: &jetfuelLoginState{}}
	session.updateJetfuelState(parsed)

	result := session.jetfuelAuthMethodChoiceResult(parsed)
	if result == nil {
		t.Fatal("jetfuelAuthMethodChoiceResult() returned nil")
	}
	if result.Status != WebLoginStatusNeedsAuthMethod {
		t.Fatalf("Status = %s, want %s", result.Status, WebLoginStatusNeedsAuthMethod)
	}
	if len(result.AuthMethods) != 2 {
		t.Fatalf("AuthMethods length = %d, want 2 supported methods: %#v", len(result.AuthMethods), result.AuthMethods)
	}
	wantIDs := []string{"Totp", "BackupCode"}
	for i, want := range wantIDs {
		if result.AuthMethods[i].ID != want {
			t.Fatalf("AuthMethods[%d].ID = %q, want %q", i, result.AuthMethods[i].ID, want)
		}
	}
	if session.jetfuel.twoFactorAction != endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("twoFactorAction = %q, want begin_two_factor_auth", session.jetfuel.twoFactorAction)
	}
}

func TestJetfuelLoginResponsePrefersSMSSubmitToken(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte(
		"Select a method to authenticate\x00two_factor_method\x00Text Message\x00Totp\x00Text\x00BackupCode\x00begin_two_factor_auth",
	))
	methods := parsed.authMethods()
	if len(methods) != 3 {
		t.Fatalf("authMethods() length = %d, want 3: %#v", len(methods), methods)
	}
	if methods[0].ID != "Sms" {
		t.Fatalf("authMethods()[0].ID = %q, want Sms", methods[0].ID)
	}
	if methods[0].SubmitID != "Text" {
		t.Fatalf("authMethods()[0].SubmitID = %q, want Text", methods[0].SubmitID)
	}
	if methods[0].Index != 0 {
		t.Fatalf("authMethods()[0].Index = %d, want 0", methods[0].Index)
	}
}

func TestJetfuelLoginResponseDoesNotTreatCodePromptAsAuthMethodChoice(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Enter your two factor code\x00Use your authenticator app to generate the code.\x00challenge_response\x00Totp\x00BackupCode\x00session_token"))
	if methods := parsed.authMethods(); len(methods) != 0 {
		t.Fatalf("authMethods() = %#v, want none for code prompt", methods)
	}
}

func TestJetfuelLoginResponseClassifiesPhoneAuthMethodAliases(t *testing.T) {
	tests := []struct {
		raw      string
		submitID string
	}{
		{raw: "Sms"},
		{raw: "SMS"},
		{raw: "Text", submitID: "Text"},
		{raw: "Text message"},
		{raw: "Phone number"},
		{raw: "TextMessage", submitID: "TextMessage"},
		{raw: "PhoneNumber", submitID: "PhoneNumber"},
	}
	for _, test := range tests {
		method, ok := classifyJetfuelAuthMethod(test.raw)
		if !ok {
			t.Fatalf("classifyJetfuelAuthMethod(%q) returned false", test.raw)
		}
		if method.ID != "Sms" || method.Kind != WebLoginAuthMethodKindSMS || method.Supported {
			t.Fatalf("classifyJetfuelAuthMethod(%q) = %#v, want coming-soon SMS method", test.raw, method)
		}
		if method.SubmitID != test.submitID {
			t.Fatalf("classifyJetfuelAuthMethod(%q).SubmitID = %q, want %q", test.raw, method.SubmitID, test.submitID)
		}
	}
}

func TestJetfuelAuthMethodFormShape(t *testing.T) {
	state := &jetfuelLoginState{
		sessionToken:      "session-token",
		preludeDispatchID: "prelude-id",
		userID:            "1127993589949243392",
	}
	form := state.authMethodForm(WebLoginAuthMethod{ID: "Totp", Index: 0})

	if got := form.Get("two_factor_auth_method_type"); got != "Totp" {
		t.Fatalf("two_factor_auth_method_type = %q", got)
	}
	if got := form.Get("_selected_method_idx"); got != "0" {
		t.Fatalf("_selected_method_idx = %q", got)
	}
	if got := form.Get("user_id"); got != "1127993589949243392" {
		t.Fatalf("user_id = %q", got)
	}
	if got := form.Get("session_token"); got != "session-token" {
		t.Fatalf("session_token = %q", got)
	}
	if got := form.Get("prelude_dispatch_id"); got != "" {
		t.Fatalf("prelude_dispatch_id = %q, want omitted", got)
	}
}

func TestJetfuelAuthMethodFormUsesSubmitID(t *testing.T) {
	state := &jetfuelLoginState{}
	form := state.authMethodForm(WebLoginAuthMethod{ID: "Sms", SubmitID: "TextMessage", Index: 1})

	if got := form.Get("two_factor_auth_method_type"); got != "TextMessage" {
		t.Fatalf("two_factor_auth_method_type = %q, want TextMessage", got)
	}
	if got := form.Get("_selected_method_idx"); got != "1" {
		t.Fatalf("_selected_method_idx = %q, want 1", got)
	}
}

func TestJetfuelCastleTokenUsesClientProvidedOneShotToken(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	form := make(url.Values)

	if err := client.addJetfuelCastleTokenToForm(form); !errors.Is(err, ErrJetfuelCastleTokenRequired) {
		t.Fatalf("addJetfuelCastleTokenToForm() error = %v, want ErrJetfuelCastleTokenRequired", err)
	}

	client.SetNextJetfuelCastleTokens([]string{" castle-from-webview "})
	if err := client.addJetfuelCastleTokenToForm(form); err != nil {
		t.Fatalf("addJetfuelCastleTokenToForm() error = %v", err)
	}
	if got := form.Get("$castle_token"); got != "castle-from-webview" {
		t.Fatalf("$castle_token = %q, want webview token", got)
	}
	if err := client.addJetfuelCastleTokenToForm(make(url.Values)); !errors.Is(err, ErrJetfuelCastleTokenRequired) {
		t.Fatalf("second addJetfuelCastleTokenToForm() error = %v, want one-shot token to be consumed", err)
	}
}

func TestJetfuelCastleTokenConsumesQueuedTokensInOrder(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{" first-token ", "", "second-token"})

	firstForm := make(url.Values)
	if err := client.addJetfuelCastleTokenToForm(firstForm); err != nil {
		t.Fatalf("first addJetfuelCastleTokenToForm() error = %v", err)
	}
	if got := firstForm.Get("$castle_token"); got != "first-token" {
		t.Fatalf("first $castle_token = %q, want first-token", got)
	}
	if !client.HasNextJetfuelCastleToken() {
		t.Fatal("HasNextJetfuelCastleToken() = false, want queued second token")
	}

	secondForm := make(url.Values)
	if err := client.addJetfuelCastleTokenToForm(secondForm); err != nil {
		t.Fatalf("second addJetfuelCastleTokenToForm() error = %v", err)
	}
	if got := secondForm.Get("$castle_token"); got != "second-token" {
		t.Fatalf("second $castle_token = %q, want second-token", got)
	}
	if client.HasNextJetfuelCastleToken() {
		t.Fatal("HasNextJetfuelCastleToken() = true, want queue exhausted")
	}
}

func TestSetCookiesPreservesExistingLoginCookies(t *testing.T) {
	client := NewClient(cookies.NewCookies(map[string]string{
		"att":      "native-login-cookie",
		"guest_id": "old-guest-id",
	}), nil, zerolog.Nop())

	client.SetCookies(map[string]string{
		"guest_id": "browser-guest-id",
		"__cf_bm":  "browser-cf-cookie",
	})

	if got := client.cookies.Get(cookies.XAtt); got != "native-login-cookie" {
		t.Fatalf("att cookie = %q, want native login cookie to be preserved", got)
	}
	if got := client.cookies.Get(cookies.XGuestID); got != "browser-guest-id" {
		t.Fatalf("guest_id cookie = %q, want webview cookie to update existing value", got)
	}
	if got := client.cookies.Get(cookies.XCookieName("__cf_bm")); got != "browser-cf-cookie" {
		t.Fatalf("__cf_bm cookie = %q, want webview cookie", got)
	}
}

func TestJetfuelIdentifierNoSupportedActionClassification(t *testing.T) {
	if !errors.Is(ErrJetfuelIdentifierNoSupportedAction, ErrWebLoginUnexpectedSubtask) {
		t.Fatalf("ErrJetfuelIdentifierNoSupportedAction must wrap ErrWebLoginUnexpectedSubtask")
	}
	if !isJetfuelPrePasswordParityError(ErrJetfuelIdentifierNoSupportedAction) {
		t.Fatal("identifier no-action error must enable the combined-credentials fallback")
	}
}

func TestSubmitJetfuelCredentialsFallsBackToCombinedAfterIdentifierNoAction(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_VIEWER_CONTEXT", "0")
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"identifier-token", "combined-token"})

	requestCount := 0
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_BEGIN_LOGIN_PATH {
			t.Fatalf("request path = %s, want begin_login", req.URL.Path)
		}
		requestCount++
		form := readJetfuelTestForm(t, req)
		switch requestCount {
		case 1:
			if got := form.Get("username_or_email"); got != "test-user" {
				t.Fatalf("identifier username_or_email = %q", got)
			}
			if _, ok := form["password"]; ok {
				t.Fatalf("identifier form unexpectedly contains password: %#v", form)
			}
			if got := form.Get("$castle_token"); got != "identifier-token" {
				t.Fatalf("identifier Castle token = %q", got)
			}
			return jetfuelTestResponse("identifier accepted"), nil
		case 2:
			if got := form.Get("username_or_email"); got != "test-user" {
				t.Fatalf("combined username_or_email = %q", got)
			}
			if got := form.Get("password"); got != "test-password" {
				t.Fatalf("combined password = %q", got)
			}
			if got := form.Get("$castle_token"); got != "combined-token" {
				t.Fatalf("combined Castle token = %q", got)
			}
			return jetfuelTestResponse("/home"), nil
		default:
			t.Fatalf("unexpected request %d", requestCount)
			return nil, nil
		}
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{}

	result, err := session.SubmitCredentials(context.Background(), "test-user", "test-password")
	if err != nil {
		t.Fatalf("SubmitCredentials() error = %v", err)
	}
	if result == nil || result.Status != WebLoginStatusComplete {
		t.Fatalf("SubmitCredentials() result = %#v, want complete", result)
	}
	if requestCount != 2 {
		t.Fatalf("request count = %d, want 2", requestCount)
	}
}

func TestUnsupportedJetfuelResponseLoggingOmitsResponseValues(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_VIEWER_CONTEXT", "0")
	var logs bytes.Buffer
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.New(&logs).Level(zerolog.DebugLevel))
	client.SetNextJetfuelCastleTokens([]string{"castle-secret-marker"})
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return jetfuelTestResponse("response-secret-marker\x00opaque_field"), nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{}

	result, err := session.SubmitIdentifier(context.Background(), "identifier-secret-marker")
	if result != nil {
		t.Fatalf("SubmitIdentifier() result = %#v, want nil", result)
	}
	if !errors.Is(err, ErrJetfuelIdentifierNoSupportedAction) {
		t.Fatalf("SubmitIdentifier() error = %v, want ErrJetfuelIdentifierNoSupportedAction", err)
	}
	logged := logs.String()
	for _, secret := range []string{"response-secret-marker", "opaque_field", "identifier-secret-marker", "castle-secret-marker"} {
		if strings.Contains(logged, secret) {
			t.Fatalf("sanitized diagnostics contain %q: %s", secret, logged)
		}
	}
	for _, field := range []string{"\"stage\":\"identifier\"", "\"response_bytes\":", "\"string_count\":", "\"path_count\":", "\"field_count\":"} {
		if !strings.Contains(logged, field) {
			t.Fatalf("sanitized diagnostics missing %q: %s", field, logged)
		}
	}
}

func TestSubmitJetfuelCombinedCredentialsReturnsPasswordAction(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"combined-token"})
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_BEGIN_LOGIN_PATH {
			t.Fatalf("request path = %s, want begin_login", req.URL.Path)
		}
		form := readJetfuelTestForm(t, req)
		if form.Get("username_or_email") != "test-user" || form.Get("password") != "test-password" {
			t.Fatalf("combined form = %#v", form)
		}
		return jetfuelTestResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{}

	result, err := session.SubmitCombinedCredentials(context.Background(), "test-user", "test-password")
	if err != nil {
		t.Fatalf("SubmitCombinedCredentials() error = %v", err)
	}
	if result == nil || result.Status != WebLoginStatusNeedsPassword {
		t.Fatalf("SubmitCombinedCredentials() result = %#v, want needs password", result)
	}
	if got := session.jetfuel.passwordAction; got != endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
		t.Fatalf("password action = %q", got)
	}
}

func TestSubmitJetfuelPasswordAllowsOneBoundedReplay(t *testing.T) {
	structuredActionlessBody := "/onboarding/web/actions/persist_login_state\x00opaque_field"
	tests := []struct {
		name         string
		replayBody   string
		terminalBody string
		wantStatus   WebLoginStatus
	}{
		{
			name:         "explicit password action then complete",
			replayBody:   endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password",
			terminalBody: "/home",
			wantStatus:   WebLoginStatusComplete,
		},
		{
			name:         "explicit password action then verification challenge",
			replayBody:   endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password",
			terminalBody: endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH + "\x00challenge_response\x00Enter your verification code",
			wantStatus:   WebLoginStatusNeedsText,
		},
		{
			name:         "structured actionless response then complete",
			replayBody:   structuredActionlessBody,
			terminalBody: "/home",
			wantStatus:   WebLoginStatusComplete,
		},
		{
			name:         "structured actionless response then verification challenge",
			replayBody:   structuredActionlessBody,
			terminalBody: endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH + "\x00challenge_response\x00Enter your verification code",
			wantStatus:   WebLoginStatusNeedsText,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
			client.SetNextJetfuelCastleTokens([]string{"combined-token", "first-password-token", "replay-password-token"})

			passwordRequestCount := 0
			client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				form := readJetfuelTestForm(t, req)
				if req.URL.Path == "/i/jfapi"+endpoints.JETFUEL_BEGIN_LOGIN_PATH {
					if form.Get("username_or_email") != "test-user" || form.Get("password") != "test-password" {
						t.Fatalf("combined form = %#v", form)
					}
					if got := form.Get("$castle_token"); got != "combined-token" {
						t.Fatalf("combined Castle token = %q", got)
					}
					return jetfuelTestResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
				}
				if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
					t.Fatalf("request path = %s, want login_enter_password", req.URL.Path)
				}
				passwordRequestCount++
				if got := form.Get("password"); got != "test-password" {
					t.Fatalf("password = %q", got)
				}
				wantToken := "first-password-token"
				if passwordRequestCount == 2 {
					wantToken = "replay-password-token"
				}
				if got := form.Get("$castle_token"); got != wantToken {
					t.Fatalf("password request %d Castle token = %q, want %q", passwordRequestCount, got, wantToken)
				}
				if passwordRequestCount == 1 {
					return jetfuelTestResponse(tc.replayBody), nil
				}
				return jetfuelTestResponse(tc.terminalBody), nil
			})}
			session := NewWebLoginSession(client)
			session.backend = webLoginBackendJetfuel
			session.jetfuel = &jetfuelLoginState{}

			result, err := session.SubmitCombinedCredentials(context.Background(), "test-user", "test-password")
			if err != nil {
				t.Fatalf("SubmitCombinedCredentials() error = %v", err)
			}
			if result == nil || result.Status != WebLoginStatusNeedsPassword {
				t.Fatalf("SubmitCombinedCredentials() result = %#v, want needs password", result)
			}
			if session.jetfuel.passwordReplayUsed {
				t.Fatal("combined password action incorrectly consumed the response-driven replay")
			}

			result, err = session.SubmitPassword(context.Background(), "test-password")
			if err != nil {
				t.Fatalf("first SubmitPassword() error = %v", err)
			}
			if result == nil || result.Status != WebLoginStatusNeedsPassword {
				t.Fatalf("first SubmitPassword() result = %#v, want needs password", result)
			}

			result, err = session.SubmitPassword(context.Background(), "test-password")
			if err != nil {
				t.Fatalf("replay SubmitPassword() error = %v", err)
			}
			if result == nil || result.Status != tc.wantStatus {
				t.Fatalf("replay SubmitPassword() result = %#v, want status %s", result, tc.wantStatus)
			}
			if passwordRequestCount != 2 {
				t.Fatalf("password request count = %d, want 2", passwordRequestCount)
			}
		})
	}
}

func TestSubmitJetfuelPasswordRejectsSecondResponseDrivenReplay(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"combined-token", "first-password-token", "replay-password-token", "unused-token"})

	passwordRequestCount := 0
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/i/jfapi"+endpoints.JETFUEL_BEGIN_LOGIN_PATH {
			return jetfuelTestResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
		}
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
			t.Fatalf("request path = %s, want login_enter_password", req.URL.Path)
		}
		passwordRequestCount++
		return jetfuelTestResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{}

	result, err := session.SubmitCombinedCredentials(context.Background(), "test-user", "test-password")
	if err != nil || result == nil || result.Status != WebLoginStatusNeedsPassword {
		t.Fatalf("SubmitCombinedCredentials() result = %#v, error = %v", result, err)
	}

	result, err = session.SubmitPassword(context.Background(), "test-password")
	if err != nil || result == nil || result.Status != WebLoginStatusNeedsPassword {
		t.Fatalf("first SubmitPassword() result = %#v, error = %v", result, err)
	}
	result, err = session.SubmitPassword(context.Background(), "test-password")
	if result != nil {
		t.Fatalf("replay SubmitPassword() result = %#v, want nil", result)
	}
	if !errors.Is(err, ErrWebLoginUnexpectedSubtask) {
		t.Fatalf("replay SubmitPassword() error = %v, want ErrWebLoginUnexpectedSubtask", err)
	}
	if passwordRequestCount != 2 {
		t.Fatalf("password request count = %d, want 2", passwordRequestCount)
	}
	if !client.HasNextJetfuelCastleToken() {
		t.Fatal("third Castle token was consumed; password replay was not bounded")
	}
}

func TestSubmitJetfuelPasswordRejectsSecondStructuredActionlessResponse(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"combined-token", "first-password-token", "replay-password-token", "unused-token"})

	passwordRequestCount := 0
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/i/jfapi"+endpoints.JETFUEL_BEGIN_LOGIN_PATH {
			return jetfuelTestResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00password"), nil
		}
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
			t.Fatalf("request path = %s, want login_enter_password", req.URL.Path)
		}
		passwordRequestCount++
		return jetfuelTestResponse("/onboarding/web/actions/persist_login_state\x00opaque_field"), nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{}

	result, err := session.SubmitCombinedCredentials(context.Background(), "test-user", "test-password")
	if err != nil || result == nil || result.Status != WebLoginStatusNeedsPassword {
		t.Fatalf("SubmitCombinedCredentials() result = %#v, error = %v", result, err)
	}

	result, err = session.SubmitPassword(context.Background(), "test-password")
	if err != nil || result == nil || result.Status != WebLoginStatusNeedsPassword {
		t.Fatalf("first SubmitPassword() result = %#v, error = %v", result, err)
	}
	result, err = session.SubmitPassword(context.Background(), "test-password")
	if result != nil {
		t.Fatalf("replay SubmitPassword() result = %#v, want nil", result)
	}
	if !errors.Is(err, ErrWebLoginUnexpectedSubtask) {
		t.Fatalf("replay SubmitPassword() error = %v, want ErrWebLoginUnexpectedSubtask", err)
	}
	if passwordRequestCount != 2 {
		t.Fatalf("password request count = %d, want 2", passwordRequestCount)
	}
	if !client.HasNextJetfuelCastleToken() {
		t.Fatal("third Castle token was consumed; actionless password replay was not bounded")
	}
}

func TestSubmitJetfuelPasswordDoesNotReplayUnstructuredResponse(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"password-token", "unused-token"})

	passwordRequestCount := 0
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
			t.Fatalf("request path = %s, want login_enter_password", req.URL.Path)
		}
		passwordRequestCount++
		return jetfuelTestResponse("opaque response with no action path"), nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{
		identifier:     "test-user",
		passwordAction: endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH,
	}

	result, err := session.SubmitPassword(context.Background(), "test-password")
	if result != nil {
		t.Fatalf("SubmitPassword() result = %#v, want nil", result)
	}
	if !errors.Is(err, ErrWebLoginUnexpectedSubtask) {
		t.Fatalf("SubmitPassword() error = %v, want ErrWebLoginUnexpectedSubtask", err)
	}
	if session.jetfuel.passwordReplayUsed {
		t.Fatal("unstructured response incorrectly consumed the password replay")
	}
	if passwordRequestCount != 1 {
		t.Fatalf("password request count = %d, want 1", passwordRequestCount)
	}
	if !client.HasNextJetfuelCastleToken() {
		t.Fatal("unused Castle token was consumed after an unstructured response")
	}
}

func TestSubmitJetfuelPasswordDefersTwoFactorPreludeUntilNextCastleToken(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"password-castle-token"})

	var paths []string
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(request body) error = %v", err)
		}
		paths = append(paths, req.URL.Path)
		form := string(body)
		switch req.URL.Path {
		case "/i/jfapi" + endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH:
			if !strings.Contains(form, "%24castle_token=password-castle-token") {
				t.Fatalf("password request body = %q, want password Castle token", form)
			}
			responseBody := "begin_two_factor_auth\x00session_token\x0012345678-1234-1234-1234-123456789abc\x00prelude_dispatch_id\x00abcdefab-1234-1234-1234-abcdefabcdef"
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		case "/i/jfapi" + endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH:
			if !strings.Contains(form, "%24castle_token=twofactor-castle-token") {
				t.Fatalf("two-factor request body = %q, want second Castle token", form)
			}
			if !strings.Contains(form, "session_token=12345678-1234-1234-1234-123456789abc") {
				t.Fatalf("two-factor request body = %q, want session token", form)
			}
			responseBody := "Select a method to authenticate\x00two_factor_method\x00Totp\x00BackupCode"
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		default:
			t.Fatalf("unexpected request path: %s", req.URL.Path)
			return nil, nil
		}
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{
		identifier:     "test-user",
		passwordAction: endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH,
	}

	result, err := session.SubmitPassword(context.Background(), "password")
	if !errors.Is(err, ErrJetfuelCastleTokenRequired) {
		t.Fatalf("SubmitPassword() error = %v, want ErrJetfuelCastleTokenRequired", err)
	}
	if result != nil {
		t.Fatalf("SubmitPassword() result = %#v, want nil while waiting for next Castle token", result)
	}
	if got := session.jetfuel.twoFactorAction; got != endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("twoFactorAction = %q, want begin two-factor action", got)
	}
	if len(paths) != 1 || paths[0] != "/i/jfapi"+endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH {
		t.Fatalf("paths after password = %#v, want only password request", paths)
	}

	client.SetNextJetfuelCastleTokens([]string{"twofactor-castle-token"})
	result, err = session.SubmitPendingTwoFactor(context.Background())
	if err != nil {
		t.Fatalf("SubmitPendingTwoFactor() error = %v", err)
	}
	if result == nil || result.Status != WebLoginStatusNeedsAuthMethod {
		t.Fatalf("SubmitPendingTwoFactor() result = %#v, want auth method chooser", result)
	}
	if len(paths) != 2 || paths[1] != "/i/jfapi"+endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("paths after pending two-factor = %#v", paths)
	}
}

func TestSubmitJetfuelAuthMethodPrefersVerificationChallenge(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"castle-from-webview"})
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
			t.Fatalf("request path = %s", req.URL.Path)
		}
		if got := req.Header.Get("x-jf-client-theme"); got != jetfuelHeaderTheme {
			t.Fatalf("x-jf-client-theme = %q, want %q", got, jetfuelHeaderTheme)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(request body) error = %v", err)
		}
		form := string(body)
		if !strings.Contains(form, "two_factor_auth_method_type=Totp") {
			t.Fatalf("request body = %q, want Totp method", form)
		}
		if !strings.Contains(form, "%24castle_token=castle-from-webview") {
			t.Fatalf("request body = %q, want webview Castle token", form)
		}
		responseBody := "Select a method to authenticate\x00two_factor_method\x00Totp\x00BackupCode\x00U2fSecurityKey\x00" +
			"Enter the code from your authentication app.\x00challenge_response\x00finish_two_factor_auth\x00session_token"
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(responseBody)),
		}, nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{
		sessionToken:    "session-token",
		twoFactorAction: endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH,
		twoFactorMethods: []WebLoginAuthMethod{
			{ID: "Totp", Name: "Authenticator App", Kind: WebLoginAuthMethodKindCode, Supported: true},
			{ID: "BackupCode", Name: "Backup Code", Kind: WebLoginAuthMethodKindBackupCode, Supported: true, Index: 1},
			{ID: "U2fSecurityKey", Name: "Security Key PC", Kind: WebLoginAuthMethodKindUnknown, Supported: false, Index: 2},
		},
	}

	result, err := session.SubmitAuthMethod(context.Background(), "Authenticator App")
	if err != nil {
		t.Fatalf("SubmitAuthMethod() error = %v", err)
	}
	if result.Status != WebLoginStatusNeedsText {
		t.Fatalf("SubmitAuthMethod() status = %s, want %s", result.Status, WebLoginStatusNeedsText)
	}
	if result.Challenge == nil || result.Challenge.Description != "Enter the code from your authentication app." {
		t.Fatalf("Challenge = %#v, want authenticator app code prompt", result.Challenge)
	}
	if session.jetfuel.verificationAction != endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("verificationAction = %q", session.jetfuel.verificationAction)
	}
}

func TestSubmitJetfuelSMSAuthMethodReturnsPhoneChallenge(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"sms-method-castle-token"})
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
			t.Fatalf("request path = %s", req.URL.Path)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(request body) error = %v", err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("ParseQuery(request body) error = %v", err)
		}
		if got := form.Get("two_factor_auth_method_type"); got != "Sms" {
			t.Fatalf("two_factor_auth_method_type = %q, want Sms", got)
		}
		if got := form.Get("_selected_method_idx"); got != "1" {
			t.Fatalf("_selected_method_idx = %q, want 1", got)
		}
		responseBody := "We sent a text message with a verification code to your phone.\x00challenge_response\x00finish_two_factor_auth\x00session_token"
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(responseBody)),
		}, nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{
		sessionToken:    "session-token",
		twoFactorAction: endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH,
		twoFactorMethods: []WebLoginAuthMethod{
			{ID: "Totp", Name: "Authenticator App", Kind: WebLoginAuthMethodKindCode, Supported: true},
			{ID: "Sms", Name: "Text Message", Kind: WebLoginAuthMethodKindSMS, Supported: true, Index: 1},
			{ID: "BackupCode", Name: "Backup Code", Kind: WebLoginAuthMethodKindBackupCode, Supported: true, Index: 2},
		},
	}

	result, err := session.SubmitAuthMethod(context.Background(), "Text Message")
	if err != nil {
		t.Fatalf("SubmitAuthMethod() error = %v", err)
	}
	if result.Status != WebLoginStatusNeedsText {
		t.Fatalf("SubmitAuthMethod() status = %s, want %s", result.Status, WebLoginStatusNeedsText)
	}
	if result.Challenge == nil || result.Challenge.Description != "Enter the code sent to your phone number." {
		t.Fatalf("Challenge = %#v, want phone code prompt", result.Challenge)
	}
	if session.jetfuel.verificationAction != endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("verificationAction = %q", session.jetfuel.verificationAction)
	}
}

func TestSubmitJetfuelSMSAuthMethodDefaultsActionForPhoneChallenge(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"sms-phone-challenge-castle-token"})
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH {
			t.Fatalf("request path = %s", req.URL.Path)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(request body) error = %v", err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("ParseQuery(request body) error = %v", err)
		}
		if got := form.Get("two_factor_auth_method_type"); got != "Text" {
			t.Fatalf("two_factor_auth_method_type = %q, want Text", got)
		}
		responseBody := "Confirm the phone number associated with your X account.\x00phone_number\x00session_token"
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(responseBody)),
		}, nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{
		sessionToken:    "session-token",
		twoFactorAction: endpoints.JETFUEL_BEGIN_TWO_FACTOR_AUTH_PATH,
		twoFactorMethods: []WebLoginAuthMethod{
			{ID: "Totp", Name: "Authenticator App", Kind: WebLoginAuthMethodKindCode, Supported: true},
			{ID: "Sms", SubmitID: "Text", Name: "Text Message", Kind: WebLoginAuthMethodKindSMS, Supported: true, Index: 1},
			{ID: "BackupCode", Name: "Backup Code", Kind: WebLoginAuthMethodKindBackupCode, Supported: true, Index: 2},
		},
	}

	result, err := session.SubmitAuthMethod(context.Background(), "Text Message")
	if err != nil {
		t.Fatalf("SubmitAuthMethod() error = %v", err)
	}
	if result.Status != WebLoginStatusNeedsText {
		t.Fatalf("SubmitAuthMethod() status = %s, want %s", result.Status, WebLoginStatusNeedsText)
	}
	if result.Challenge == nil || result.Challenge.InputKind != WebLoginChallengeInputKindPhoneNumber {
		t.Fatalf("Challenge = %#v, want phone-number prompt", result.Challenge)
	}
	if session.jetfuel.verificationAction != endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH {
		t.Fatalf("verificationAction = %q", session.jetfuel.verificationAction)
	}
}

func TestSubmitJetfuelPhoneNumberVerificationPostsPhoneField(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.SetNextJetfuelCastleTokens([]string{"phone-number-castle-token"})
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/i/jfapi"+endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH {
			t.Fatalf("request path = %s", req.URL.Path)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(request body) error = %v", err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("ParseQuery(request body) error = %v", err)
		}
		if got := form.Get("phone_number"); got != "+15551234567" {
			t.Fatalf("phone_number = %q, want test phone number", got)
		}
		if got := form.Get("challenge_response"); got != "" {
			t.Fatalf("challenge_response = %q, want omitted when phone_number is known", got)
		}
		if got := form.Get("verification_code"); got != "" {
			t.Fatalf("verification_code = %q, want omitted when phone_number is known", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("/home")),
		}, nil
	})}
	session := NewWebLoginSession(client)
	session.backend = webLoginBackendJetfuel
	session.jetfuel = &jetfuelLoginState{
		verificationAction: endpoints.JETFUEL_FINISH_TWO_FACTOR_AUTH_PATH,
		verificationFields: []string{"phone_number"},
	}

	result, err := session.SubmitText(context.Background(), "+15551234567")
	if err != nil {
		t.Fatalf("SubmitText() error = %v", err)
	}
	if result.Status != WebLoginStatusComplete {
		t.Fatalf("SubmitText() status = %s, want %s", result.Status, WebLoginStatusComplete)
	}
}

func TestJetfuelAuthMethodChoiceWithOnlyUnsupportedMethod(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Select a method to authenticate\x00two_factor_method\x00U2fSecurityKey\x00begin_two_factor_auth"))
	session := &WebLoginSession{jetfuel: &jetfuelLoginState{}}

	result := session.jetfuelAuthMethodChoiceResult(parsed)
	if result == nil {
		t.Fatal("jetfuelAuthMethodChoiceResult() returned nil")
	}
	if result.Status != WebLoginStatusUnsupported {
		t.Fatalf("Status = %s, want %s", result.Status, WebLoginStatusUnsupported)
	}
	if len(session.jetfuel.twoFactorMethods) != 1 || session.jetfuel.twoFactorMethods[0].Supported {
		t.Fatalf("twoFactorMethods = %#v, want one unsupported method", session.jetfuel.twoFactorMethods)
	}
}

func TestJetfuelAuthMethodChoiceWithOnlySMSIsComingSoon(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Select a method to authenticate\x00two_factor_method\x00Text\x00begin_two_factor_auth"))
	session := &WebLoginSession{jetfuel: &jetfuelLoginState{}}

	result := session.jetfuelAuthMethodChoiceResult(parsed)
	if result == nil {
		t.Fatal("jetfuelAuthMethodChoiceResult() returned nil")
	}
	if result.Status != WebLoginStatusUnsupported {
		t.Fatalf("Status = %s, want %s", result.Status, WebLoginStatusUnsupported)
	}
	if len(session.jetfuel.twoFactorMethods) != 1 {
		t.Fatalf("twoFactorMethods length = %d, want 1: %#v", len(session.jetfuel.twoFactorMethods), session.jetfuel.twoFactorMethods)
	}
	method := session.jetfuel.twoFactorMethods[0]
	if method.ID != "Sms" || method.Supported || method.Description != "Text message verification is coming soon." {
		t.Fatalf("SMS method = %#v, want coming-soon unsupported method", method)
	}
	if result.Challenge == nil || result.Challenge.Description != "Text message verification is coming soon." {
		t.Fatalf("Challenge = %#v, want coming-soon description", result.Challenge)
	}
}

func TestSubmitJetfuelAuthMethodRejectsUnsupportedMethod(t *testing.T) {
	session := &WebLoginSession{
		backend: webLoginBackendJetfuel,
		jetfuel: &jetfuelLoginState{
			twoFactorMethods: []WebLoginAuthMethod{{
				ID:        "U2fSecurityKey",
				Name:      "Security Key PC",
				Supported: false,
			}},
		},
	}

	_, err := session.SubmitAuthMethod(context.Background(), "U2fSecurityKey")
	if !errors.Is(err, ErrWebLoginUnsupportedAuthMethod) {
		t.Fatalf("SubmitAuthMethod() error = %v, want ErrWebLoginUnsupportedAuthMethod", err)
	}
}

func TestJetfuelLoginResponseClassifiesTemporaryLimit(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("We've temporarily limited your login. Please try again later."))
	err := parsed.loginError()
	if err == nil {
		t.Fatal("loginError() returned nil")
	}
	webErr, ok := err.(*WebLoginError)
	if !ok {
		t.Fatalf("loginError() = %T, want *WebLoginError", err)
	}
	if webErr.Code != 399 {
		t.Fatalf("WebLoginError.Code = %d, want 399", webErr.Code)
	}
}

func TestJetfuelLoginResponseClassifiesTooManyAttempts(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("errors.Too many attempts. Try again in a few minutes.\x00message.Too many attempts. Try again in a few minutes."))
	err := parsed.loginError()
	if err == nil {
		t.Fatal("loginError() returned nil")
	}
	webErr, ok := err.(*WebLoginError)
	if !ok {
		t.Fatalf("loginError() = %T, want *WebLoginError", err)
	}
	if webErr.Code != 399 {
		t.Fatalf("WebLoginError.Code = %d, want 399", webErr.Code)
	}
}

func TestJetfuelLoginResponseClassifiesOfficialClientError(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("Please use X.com or official X apps to proceed with log in/sign up."))
	err := parsed.loginError()
	if err == nil {
		t.Fatal("loginError() returned nil")
	}
	webErr, ok := err.(*WebLoginError)
	if !ok {
		t.Fatalf("loginError() = %T, want *WebLoginError", err)
	}
	if webErr.Code != 399 {
		t.Fatalf("WebLoginError.Code = %d, want 399", webErr.Code)
	}
	if !strings.Contains(webErr.Message, "official X apps") {
		t.Fatalf("WebLoginError.Message = %q", webErr.Message)
	}
}

func TestJetfuelLoginResponseClassifiesMissingAccount(t *testing.T) {
	parsed := parseJetfuelLoginResponse([]byte("missing_account This email or username is not registered yet"))
	err := parsed.loginError()
	if err == nil {
		t.Fatal("loginError() returned nil")
	}
	webErr, ok := err.(*WebLoginError)
	if !ok {
		t.Fatalf("loginError() = %T, want *WebLoginError", err)
	}
	if webErr.Code != 32 {
		t.Fatalf("WebLoginError.Code = %d, want 32", webErr.Code)
	}
}

func TestJetfuelLoginResponseClassifiesBadCredentials(t *testing.T) {
	tests := []string{
		"Wrong password",
		"The password you entered is incorrect.",
		"The username and password you entered did not match our records.",
		"Invalid username or password",
		"Invalid credentials",
	}
	for _, body := range tests {
		parsed := parseJetfuelLoginResponse([]byte(body))
		err := parsed.loginError()
		if err == nil {
			t.Fatalf("loginError(%q) returned nil", body)
		}
		webErr, ok := err.(*WebLoginError)
		if !ok {
			t.Fatalf("loginError(%q) = %T, want *WebLoginError", body, err)
		}
		if webErr.Code != 32 {
			t.Fatalf("WebLoginError.Code for %q = %d, want 32", body, webErr.Code)
		}
	}
}
func TestJetfuelTimezoneEnvOverride(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_TIMEZONE", "Europe/Paris")
	if got := jetfuelTimezone(); got != "Europe/Paris" {
		t.Fatalf("jetfuelTimezone() = %q, want Europe/Paris", got)
	}
}
