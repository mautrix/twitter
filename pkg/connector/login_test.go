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
		{ID: "BackupCode", Name: "Backup Code", Supported: true},
		{ID: "U2fSecurityKey", Name: "Security Key PC", Supported: true},
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
	if strings.Join(field.Options, ",") != "Authenticator App,Backup Code,Security Key PC" {
		t.Fatalf("field.Options = %#v", field.Options)
	}
	if strings.Contains(step.Instructions, "not supported") {
		t.Fatalf("Instructions = %q, want no unsupported caveat", step.Instructions)
	}
}

func TestFindWebLoginAuthMethodMatchesNameOrID(t *testing.T) {
	methods := []twittermeow.WebLoginAuthMethod{
		{ID: "Totp", Name: "Authenticator App", Supported: true},
		{ID: "U2fSecurityKey", Name: "Security Key PC", Supported: true},
	}
	if method, ok := findWebLoginAuthMethod(methods, "Authenticator App"); !ok || method.ID != "Totp" {
		t.Fatalf("find by label = %#v %t, want Totp", method, ok)
	}
	if method, ok := findWebLoginAuthMethod(methods, "u2f_security_key"); !ok || method.ID != "U2fSecurityKey" {
		t.Fatalf("find by normalized ID = %#v %t, want U2fSecurityKey", method, ok)
	}
}
