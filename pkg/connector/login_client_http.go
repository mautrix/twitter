package connector

import (
	"context"
	"errors"
	"strings"

	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

const (
	clientHTTPStageStart          = "start"
	clientHTTPMaxImmediateActions = 16
)

func (t *TwitterLogin) isWaitingForClientHTTPRequest() bool {
	return t.useClientHTTPLogin &&
		t.webLogin != nil &&
		t.webLogin.Client() != nil &&
		t.clientHTTPTransport != nil &&
		t.clientHTTPTransport.PendingRequest() != nil &&
		t.clientHTTPStage != ""
}

func (t *TwitterLogin) submitClientHTTPInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	if !t.isWaitingForClientHTTPRequest() {
		t.stopClientHTTPLogin()
		return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
	}
	if strings.TrimSpace(input[loginFieldClientHTTPError]) != "" {
		return t.failClientHTTPRequest("The request did not complete on this device. Please try again.")
	}
	response, err := decodeClientHTTPResponse(input)
	if err != nil {
		return t.failClientHTTPRequest("The client returned an invalid HTTP response. Please try again.")
	}
	if !t.webLogin.Client().SetBrowserHeaders(browserHeadersFromInput(input)) {
		return t.failClientHTTPRequest("The client did not return a valid browser fingerprint. Please try again.")
	}
	if err = t.clientHTTPTransport.SubmitResponse(response); err != nil {
		return t.failClientHTTPRequest("The client HTTP response did not match this login. Please try again.")
	}
	t.browserHeaders = t.webLogin.Client().GetBrowserHeaders()
	t.webLogin.Client().SetCookies(clientHTTPCookies(input))
	return t.continueClientHTTPLogin(ctx)
}

func (t *TwitterLogin) failClientHTTPRequest(message string) (*bridgev2.LoginStep, error) {
	t.stopClientHTTPLogin()
	t.webLogin = nil
	t.webLoginIdentifier = ""
	t.webLoginPassword = ""
	t.webLoginCastleStage = ""
	t.webLoginAuthMethod = ""
	t.webLoginText = ""
	t.webLoginChallenge = nil
	t.webLoginMethods = nil
	return makeCredentialsStep(message), nil
}

func (t *TwitterLogin) continueClientHTTPLogin(ctx context.Context) (*bridgev2.LoginStep, error) {
	if t.webLogin == nil || t.webLogin.Client() == nil || t.clientHTTPTransport == nil || t.clientHTTPStage == "" {
		t.stopClientHTTPLogin()
		return makeCredentialsStep("The X login session expired. Enter your X login details again."), nil
	}

	for range clientHTTPMaxImmediateActions {
		stage := t.clientHTTPStage
		if err := t.clientHTTPTransport.BeginOperation(stage); err != nil {
			t.stopClientHTTPLogin()
			return nil, webLoginFailureError(err)
		}

		result, err := t.runClientHTTPStage(ctx, stage)
		if errors.Is(err, twittermeow.ErrClientHTTPRequestPending) {
			return makeClientHTTPRequestStep(t.clientHTTPTransport.PendingRequest())
		}
		if errors.Is(err, twittermeow.ErrJetfuelCastleTokenRequired) {
			t.webLoginCastleStage = stage
			return t.makeWebLoginCastleTokenStep(""), nil
		}
		if err != nil {
			logWebCastleFailure(ctx, stage, err)
			step, handledErr := handleClientHTTPStageError(stage, t.webLoginChallenge, t.webLoginMethods, err)
			t.clientHTTPStage = ""
			if handledErr != nil {
				t.stopClientHTTPLogin()
			} else {
				t.clientHTTPTransport.ResetOperation()
			}
			return step, handledErr
		}
		switch {
		case stage == clientHTTPStageStart && result != nil && result.Status == twittermeow.WebLoginStatusNeedsIdentifier:
			if err = t.clientHTTPTransport.EndOperation(); err != nil {
				t.stopClientHTTPLogin()
				return nil, webLoginFailureError(err)
			}
			t.clientHTTPStage = webLoginCastleStageCombined
			continue
		case result != nil && result.Status == twittermeow.WebLoginStatusNeedsPassword && t.webLoginPassword != "":
			if err = t.clientHTTPTransport.EndOperation(); err != nil {
				t.stopClientHTTPLogin()
				return nil, webLoginFailureError(err)
			}
			t.clientHTTPStage = webLoginCastleStagePassword
			continue
		}

		step, err := t.handleWebLoginResult(ctx, result)
		if errors.Is(err, twittermeow.ErrClientHTTPRequestPending) {
			return makeClientHTTPRequestStep(t.clientHTTPTransport.PendingRequest())
		}
		if errors.Is(err, twittermeow.ErrJetfuelCastleTokenRequired) {
			t.webLoginCastleStage = stage
			return t.makeWebLoginCastleTokenStep(""), nil
		}
		if err != nil {
			t.stopClientHTTPLogin()
			return nil, err
		}
		if err = t.clientHTTPTransport.EndOperation(); err != nil {
			t.stopClientHTTPLogin()
			return nil, webLoginFailureError(err)
		}
		if stage == webLoginCastleStageAuthMethod {
			t.webLoginAuthMethod = ""
		}
		if stage == webLoginCastleStageText {
			t.webLoginText = ""
		}
		t.clientHTTPStage = ""
		if result != nil && result.Status == twittermeow.WebLoginStatusComplete {
			t.stopClientHTTPLogin()
		}
		return step, nil
	}

	t.stopClientHTTPLogin()
	return nil, webLoginFailureError(errors.New("client HTTP login exceeded the immediate action limit"))
}

func (t *TwitterLogin) runClientHTTPStage(ctx context.Context, stage string) (*twittermeow.WebLoginResult, error) {
	switch stage {
	case clientHTTPStageStart:
		return t.webLogin.Start(ctx)
	case webLoginCastleStagePassword:
		return t.webLogin.SubmitPassword(ctx, t.webLoginPassword)
	case webLoginCastleStageCombined:
		return t.webLogin.SubmitCombinedCredentials(ctx, t.webLoginIdentifier, t.webLoginPassword)
	case webLoginCastleStageAuthMethod:
		return t.webLogin.SubmitAuthMethod(ctx, t.webLoginAuthMethod)
	case webLoginCastleStageText:
		return t.webLogin.SubmitText(ctx, t.webLoginText)
	default:
		return nil, errors.New("unknown client HTTP login stage")
	}
}

func handleClientHTTPStageError(
	stage string,
	challenge *twittermeow.WebLoginChallenge,
	methods []twittermeow.WebLoginAuthMethod,
	err error,
) (*bridgev2.LoginStep, error) {
	if stage == clientHTTPStageStart {
		return nil, webLoginFailureError(err)
	}
	return handleWebCastleStageError(stage, challenge, methods, err)
}

func (t *TwitterLogin) stopClientHTTPLogin() {
	if t.webLogin != nil && t.webLogin.Client() != nil {
		t.webLogin.Client().DisableClientHTTP()
	} else if t.clientHTTPTransport != nil {
		t.clientHTTPTransport.ResetOperation()
	}
	t.clientHTTPTransport = nil
	t.clientHTTPStage = ""
}
