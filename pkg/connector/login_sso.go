// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2026 Radon Rosborough
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"encoding/json"
	"net/url"
	"strings"

	"go.mau.fi/util/l10n"
	"go.mau.fi/util/random"
	"maunium.net/go/mautrix/bridgev2"
)

const (
	LoginFlowIDGoogle = "google"
	LoginFlowIDApple  = "apple"

	LoginStepIDSSO = "fi.mau.twitter.login.sso"

	loginFieldSSOResponse = "sso_response"

	// The client IDs x.com's own sign-in buttons use.
	xGoogleClientID = "49625052041-kgt0hghf445lmcmhijv46b715m2mpbct.apps.googleusercontent.com"
	xAppleClientID  = "com.twitter.twitter.siwa"
)

// x.com's Google and Apple buttons only work as popups that post the result back to the opener,
// which login webviews don't support. Instead, the sign-in step opens the provider's page directly
// and reads the result there, and the bridge submits it like x.com's button would.

// With redirect_uri=postmessage and no opener, Google's last page stays on "Processing your
// request..." with the ID token embedded in the page.
const googleSSOExtractJS = `new Promise(resolve => {
  const find = () => {
    if (location.hostname !== "accounts.google.com") return;
    const match = document.documentElement.outerHTML.match(/id_token(?:\\u003d|=)([A-Za-z0-9_\-]+\.[A-Za-z0-9_\-]+\.[A-Za-z0-9_\-]+)/);
    if (!match) return;
    clearInterval(interval);
    resolve({` + loginFieldSSOResponse + `: JSON.stringify({
      clientId: "` + xGoogleClientID + `",
      client_id: "` + xGoogleClientID + `",
      credential: match[1],
      select_by: "btn",
    })});
  };
  const interval = setInterval(find, 500);
  find();
})`

// With response_mode=form_post, Apple's last step submits a form to x.com. This takes the result
// from that form and stops the submission.
const appleSSOExtractJS = `new Promise(resolve => {
  if (location.hostname !== "appleid.apple.com") return;
  const capture = form => {
    if (!(form instanceof HTMLFormElement) || !form.elements.namedItem("id_token")) return false;
    const field = name => form.elements.namedItem(name)?.value || undefined;
    const data = {
      authorization: { code: field("code"), id_token: field("id_token"), state: field("state") },
    };
    const user = field("user");
    if (user) {
      try {
        data.user = JSON.parse(user);
      } catch (_) {}
    }
    resolve({` + loginFieldSSOResponse + `: JSON.stringify(data)});
    return true;
  };
  for (const name of ["submit", "requestSubmit"]) {
    const original = HTMLFormElement.prototype[name];
    HTMLFormElement.prototype[name] = function(...args) {
      if (!capture(this)) return original.apply(this, args);
    };
  }
  window.addEventListener("submit", evt => {
    if (capture(evt.target)) evt.preventDefault();
  }, true);
})`

func ssoProviderForFlow(flowID string) string {
	switch flowID {
	case LoginFlowIDGoogle:
		return "google"
	case LoginFlowIDApple:
		return "apple"
	default:
		return ""
	}
}

func (t *TwitterLogin) makeSSOStep(errorLine l10n.Message) *bridgev2.LoginStep {
	t.waitingForSSOResponse = true
	var loginURL, extractJS string
	instructions, errorInstructions := msgStepSSOGoogleInstructions, msgStepSSOGoogleErrorInstructions
	switch t.ssoProvider {
	case "google":
		loginURL = "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
			"client_id":     {xGoogleClientID},
			"redirect_uri":  {"postmessage"},
			"response_type": {"id_token"},
			"scope":         {"openid email profile"},
			"prompt":        {"select_account"},
			"nonce":         {random.String(16)},
			"origin":        {"https://x.com"},
		}.Encode()
		extractJS = googleSSOExtractJS
	case "apple":
		instructions, errorInstructions = msgStepSSOAppleInstructions, msgStepSSOAppleErrorInstructions
		loginURL = "https://appleid.apple.com/auth/authorize?" + url.Values{
			"client_id":     {xAppleClientID},
			"redirect_uri":  {"https://x.com"},
			"response_type": {"code id_token"},
			"scope":         {"name email"},
			"response_mode": {"form_post"},
		}.Encode()
		extractJS = appleSSOExtractJS
	}
	stepInstructions := instructions.Message()
	if !errorLine.IsZero() {
		stepInstructions = errorInstructions.Render(l10n.Params{"error": errorLine})
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       LoginStepIDSSO,
		Instructions: stepInstructions,
		CookiesParams: &bridgev2.LoginCookiesParams{
			URL: loginURL,
			Fields: []bridgev2.LoginCookieField{{
				ID:       loginFieldSSOResponse,
				Required: true,
				Sources: []bridgev2.LoginCookieFieldSource{{
					Type: bridgev2.LoginCookieTypeSpecial,
					Name: "fi.mau.twitter." + loginFieldSSOResponse,
				}},
			}},
			ExtractJS:  extractJS,
			JSRequired: true,
		},
	}
}

func (t *TwitterLogin) submitSSOResponse(input map[string]string) *bridgev2.LoginStep {
	response := strings.TrimSpace(input[loginFieldSSOResponse])
	if response == "" || !json.Valid([]byte(response)) {
		return t.makeSSOStep(msgErrSSOInvalidResult.Message())
	}
	t.waitingForSSOResponse = false
	t.resetWebLoginState()
	t.webLoginSSOResponse = response
	// The Castle step returns the browser headers, so unlike the password flow this doesn't need
	// the browser identity step first.
	return makePrepareLoginStep()
}

// restartStep is where a failed login goes back to: the sign-in window for SSO flows, the
// credentials form otherwise.
func (t *TwitterLogin) restartStep(errorLine l10n.Message) *bridgev2.LoginStep {
	if t.ssoProvider != "" {
		t.resetWebLoginState()
		return t.makeSSOStep(errorLine)
	}
	return makeCredentialsStep(errorLine)
}
