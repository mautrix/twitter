package connector

import (
	"strings"
	"testing"

	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

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
