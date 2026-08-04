package connector

import (
	"context"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
)

// clientHTTPFailureInstructions is shown when the Beeper client couldn't execute a login
// request that was proxied to it (network error, blocked request, closed webview, ...).
const clientHTTPFailureInstructions = "The request did not complete on this device. Please try again."

func (t *TwitterLogin) StartWithParams(ctx context.Context, params bridgev2.LoginStartParams) (*bridgev2.LoginStep, error) {
	if params.HTTP != nil {
		t.loginHTTPTransport = params.HTTP
	}
	if params.Override != nil {
		return t.startWithOverride(ctx, params.Override)
	}
	return t.start(ctx)
}

func (t *TwitterLogin) Start(ctx context.Context) (*bridgev2.LoginStep, error) {
	return t.StartWithParams(ctx, bridgev2.LoginStartParams{})
}

func (t *TwitterLogin) StartWithOverride(ctx context.Context, override *bridgev2.UserLogin) (*bridgev2.LoginStep, error) {
	return t.StartWithParams(ctx, bridgev2.LoginStartParams{Override: override})
}

func (t *TwitterLogin) log() zerolog.Logger {
	if t.User == nil {
		return zerolog.Nop()
	}
	return t.User.Log
}

func (t *TwitterLogin) newLoginClient() *twittermeow.Client {
	log := t.log().With().Str("component", "login_twitter_client").Logger()
	client := twittermeow.NewClient(twitCookies.NewCookies(nil), nil, log)
	client.SetLoginHTTPTransport(t.loginHTTPTransport) // nil resets if not provided
	return client
}

func (t *TwitterLogin) mapClientHTTPFailure(step *bridgev2.LoginStep, err error) (*bridgev2.LoginStep, error) {
	if err == nil || !twittermeow.IsClientHTTPError(err) {
		return step, err
	}
	log := t.log()
	log.Warn().Err(err).Msg("Login client failed to execute a proxied X request")
	t.resetWebLoginState()
	return makeCredentialsStep(clientHTTPFailureInstructions), nil
}

func (t *TwitterLogin) clearWebLoginInputs() {
	t.webLoginIdentifier = ""
	t.webLoginPassword = ""
	t.webLoginCastleStage = ""
	t.webLoginAuthMethod = ""
	t.webLoginText = ""
	t.webLoginChallenge = nil
	t.webLoginMethods = nil
}

func (t *TwitterLogin) resetWebLoginState() {
	t.clearWebLoginInputs()
	t.webLogin = nil
}
