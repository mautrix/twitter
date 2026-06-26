//go:build liveprobe

package connector

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
)

func TestLiveNativeLoginFlowProbe(t *testing.T) {
	identifier := strings.TrimSpace(os.Getenv("TWITTER_LIVE_IDENTIFIER"))
	password := os.Getenv("TWITTER_LIVE_PASSWORD")
	if identifier == "" || password == "" {
		t.Skip("TWITTER_LIVE_IDENTIFIER and TWITTER_LIVE_PASSWORD are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	tc := &TwitterConnector{}
	user := &bridgev2.User{
		User: &database.User{
			MXID: id.UserID("@highest:beeper.com"),
		},
		Log: zerolog.Nop(),
	}
	process, err := tc.CreateLogin(ctx, user, LoginFlowIDPassword)
	if err != nil {
		t.Fatalf("CreateLogin() failed: %v", err)
	}
	defer process.Cancel()

	first, err := process.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	if first == nil {
		t.Fatal("Start() returned nil step")
	}
	t.Logf("first step: type=%s id=%s", first.Type, first.StepID)
	if first.Type != bridgev2.LoginStepTypeUserInput || first.StepID != LoginStepIDCredentials {
		t.Fatalf("unexpected first step: type=%s id=%s", first.Type, first.StepID)
	}

	next, err := process.(bridgev2.LoginProcessUserInput).SubmitUserInput(ctx, map[string]string{
		loginFieldIdentifier: identifier,
		loginFieldPassword:   password,
	})
	if err != nil {
		t.Fatalf("SubmitUserInput(credentials) failed: %v", err)
	}
	if next == nil {
		t.Fatal("SubmitUserInput(credentials) returned nil step")
	}
	t.Logf("next step: type=%s id=%s instructions=%q", next.Type, next.StepID, next.Instructions)

	switch next.StepID {
	case LoginStepJuiceboxPIN, LoginStepIDVerification, LoginStepIDComplete:
		verificationCode := strings.TrimSpace(os.Getenv("TWITTER_LIVE_VERIFICATION_CODE"))
		if next.StepID == LoginStepIDVerification {
			if verificationCode == "" {
				return
			}
			next, err = process.(bridgev2.LoginProcessUserInput).SubmitUserInput(ctx, map[string]string{
				loginFieldVerificationCode: verificationCode,
			})
			if err != nil {
				t.Fatalf("SubmitUserInput(verification) failed: %v", err)
			}
			if next == nil {
				t.Fatal("SubmitUserInput(verification) returned nil step")
			}
			t.Logf("after verification step: type=%s id=%s instructions=%q", next.Type, next.StepID, next.Instructions)
		}
		if verificationCode != "" && next.StepID == LoginStepIDVerification {
			t.Fatalf("verification code was provided but flow stayed on verification: %s", next.Instructions)
		}
		if next.StepID != LoginStepJuiceboxPIN {
			return
		}
	case LoginStepIDCredentials:
		if strings.Contains(next.Instructions, "Wait a bit") ||
			strings.Contains(next.Instructions, "cannot log") ||
			strings.Contains(next.Instructions, "could not log") ||
			strings.Contains(next.Instructions, "returned a login challenge") {
			t.Fatalf("native credential submission returned retry/error step: %s", next.Instructions)
		}
		t.Fatalf("native credential submission did not advance past credentials: %s", next.Instructions)
	default:
		t.Fatalf("unexpected next step after credentials: type=%s id=%s instructions=%q", next.Type, next.StepID, next.Instructions)
	}

	pin := strings.TrimSpace(os.Getenv("TWITTER_LIVE_PIN"))
	if pin == "" {
		return
	}
	if next.StepID != LoginStepJuiceboxPIN {
		t.Fatalf("PIN was provided but flow did not reach PIN step: type=%s id=%s instructions=%q", next.Type, next.StepID, next.Instructions)
	}
	final, err := process.(bridgev2.LoginProcessUserInput).SubmitUserInput(ctx, map[string]string{
		"pin": pin,
	})
	if err != nil {
		t.Fatalf("SubmitUserInput(pin) failed: %v", err)
	}
	if final == nil {
		t.Fatal("SubmitUserInput(pin) returned nil step")
	}
	t.Logf("final step: type=%s id=%s", final.Type, final.StepID)
	if final.Type != bridgev2.LoginStepTypeComplete || final.StepID != LoginStepIDComplete {
		t.Fatalf("unexpected final step after pin: type=%s id=%s instructions=%q", final.Type, final.StepID, final.Instructions)
	}
}

func TestLiveNativeLoginStageProbe(t *testing.T) {
	identifier := strings.TrimSpace(os.Getenv("TWITTER_LIVE_IDENTIFIER"))
	password := os.Getenv("TWITTER_LIVE_PASSWORD")
	if identifier == "" || password == "" {
		t.Skip("TWITTER_LIVE_IDENTIFIER and TWITTER_LIVE_PASSWORD are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	webLogin := twittermeow.NewWebLoginSession(client)

	result, err := webLogin.Start(ctx)
	logStage(t, "start", result, err)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	result, err = webLogin.SubmitCredentials(ctx, identifier, password)
	logStage(t, "credentials", result, err)
	if err != nil {
		t.Fatalf("SubmitCredentials() failed: %v", err)
	}
	switch result.Status {
	case twittermeow.WebLoginStatusComplete, twittermeow.WebLoginStatusNeedsText, twittermeow.WebLoginStatusNeedsPassword:
	default:
		t.Fatalf("SubmitCredentials() returned status %s, want complete, password, or text challenge", result.Status)
	}
}

func TestLiveIdentifierOnlyProbe(t *testing.T) {
	identifier := strings.TrimSpace(os.Getenv("TWITTER_IDENTIFIER_PROBE"))
	if identifier == "" {
		t.Skip("TWITTER_IDENTIFIER_PROBE is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, zerolog.Nop())
	webLogin := twittermeow.NewWebLoginSession(client)

	result, err := webLogin.Start(ctx)
	logStage(t, "start", result, err)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	result, err = webLogin.SubmitIdentifier(ctx, identifier)
	logStage(t, "identifier", result, err)
	if err == nil && result.Status != twittermeow.WebLoginStatusNeedsPassword {
		t.Fatalf("SubmitIdentifier() status = %s, want password step or a login error", result.Status)
	}
}

func logStage(t *testing.T, stage string, result *twittermeow.WebLoginResult, err error) {
	t.Helper()
	if result != nil {
		t.Logf("%s result: status=%s subtask=%s", stage, result.Status, result.CurrentSubtaskID)
	}
	if err != nil {
		var webErr *twittermeow.WebLoginError
		if errors.As(err, &webErr) {
			t.Logf("%s error: code=%d message=%q", stage, webErr.Code, webErr.Message)
			return
		}
		t.Logf("%s error: %T", stage, err)
	}
}
