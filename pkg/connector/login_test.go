package connector

import (
	"context"
	"errors"
	"strings"
	"testing"

	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

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
