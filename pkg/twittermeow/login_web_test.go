package twittermeow

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

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

func TestJetfuelBrowserParityHeaderDecoding(t *testing.T) {
	if got := decodeTFEGuestCookie("v1%3A12345"); got != "12345" {
		t.Fatalf("decodeTFEGuestCookie() = %q, want 12345", got)
	}
	path := "/svc/route"
	if got := decodeDtabLocal(base64.RawURLEncoding.EncodeToString([]byte(path))); got != path {
		t.Fatalf("decodeDtabLocal() = %q, want %s", got, path)
	}
	if got := decodeDtabLocal("%2Fsvc%2Froute"); got != path {
		t.Fatalf("decodeDtabLocal(encoded path) = %q, want %s", got, path)
	}
}

func TestJetfuelTimezoneEnvOverride(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_TIMEZONE", "Europe/Paris")
	if got := jetfuelTimezone(); got != "Europe/Paris" {
		t.Fatalf("jetfuelTimezone() = %q, want Europe/Paris", got)
	}
}
