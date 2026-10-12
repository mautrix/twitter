package twittermeow

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

func TestRedactJetfuelStringsKeepsOnlyProtocolStructure(t *testing.T) {
	body := strings.Join([]string{
		"/onboarding/web/actions/begin_login_acid",
		"#/onboarding/web/actions/login_enter_password",
		"session_token",
		"password",
		"Enter your password",
		"Welcome back, Nick",
		"nick mills",
		"nickmb",
		"john_doe",
		"Is this you? Confirm the code we sent to someone@example.com",
		"https://x.com/nickmb",
		"7f3c9b2e-4d1a-4c6b-9e2f-1a2b3c4d5e6f",
	}, "\x00")

	got := redactJetfuelStrings([]byte(body), "John_Doe")
	for _, kept := range []string{
		"/onboarding/web/actions/begin_login_acid",
		"#/onboarding/web/actions/login_enter_password",
		"session_token",
		"password",
	} {
		if !slices.Contains(got, kept) {
			t.Errorf("redactJetfuelStrings() = %q, want %q kept", got, kept)
		}
	}
	joined := strings.Join(got, "\n")
	for _, secret := range []string{
		"Welcome back", "Nick", "nick mills", "nickmb", "john_doe", "someone@example.com",
		"7f3c9b2e-4d1a-4c6b-9e2f-1a2b3c4d5e6f",
	} {
		if strings.Contains(joined, secret) {
			t.Errorf("redactJetfuelStrings() = %q, leaked %q", got, secret)
		}
	}
}

func TestJetfuelResponseLoggingRequiresFlag(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		var logs bytes.Buffer
		client := NewClient(cookies.NewCookies(nil), nil, zerolog.New(&logs).Level(zerolog.DebugLevel))
		client.LogRedactedLoginResponses = enabled
		client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return jetfuelTestResponse(endpoints.JETFUEL_LOGIN_ENTER_PASSWORD_PATH + "\x00Welcome back, someone_else"), nil
		})}
		client.SetNextJetfuelCastleTokens([]string{"castle-token"})

		form := url.Values{"username_or_email": {"someone_else"}}
		if _, err := client.jetfuelPostForm(context.Background(), endpoints.JETFUEL_BEGIN_LOGIN_PATH, form); err != nil {
			t.Fatalf("jetfuelPostForm() error = %v", err)
		}
		logged := strings.Contains(logs.String(), "Redacted Jetfuel login response")
		if logged != enabled {
			t.Fatalf("flag %t: logged = %t, logs = %s", enabled, logged, logs.String())
		}
		if strings.Contains(logs.String(), "someone_else") {
			t.Fatalf("flag %t: identifier leaked into logs: %s", enabled, logs.String())
		}
	}
}
