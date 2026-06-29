package twittermeow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const webAuthnOrigin = "https://x.com"

type webAuthnChallenge struct {
	RequestOptions webAuthnRequestOptions
}

type webAuthnRequestOptions struct {
	RPID             string                   `json:"rpId"`
	Challenge        string                   `json:"challenge"`
	AllowCredentials []webAuthnCredentialHint `json:"allowCredentials"`
	Extensions       webAuthnExtensions       `json:"extensions"`
}

type webAuthnCredentialHint struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Transports []string `json:"transports"`
}

type webAuthnExtensions struct {
	AppID string `json:"appid"`
}

type webAuthnAssertion struct {
	CredentialID      []byte
	ClientDataJSON    []byte
	AuthenticatorData []byte
	Signature         []byte
}

type webAuthnClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin"`
}

func (challenge webAuthnChallenge) rpID() string {
	rpID := strings.TrimSpace(challenge.RequestOptions.RPID)
	if rpID == "" {
		return "x.com"
	}
	return rpID
}

func (challenge webAuthnChallenge) clientDataJSON() ([]byte, error) {
	challengeValue := strings.TrimSpace(challenge.RequestOptions.Challenge)
	if challengeValue == "" {
		return nil, errors.New("x security-key challenge is missing")
	}
	return json.Marshal(webAuthnClientData{
		Type:        "webauthn.get",
		Challenge:   challengeValue,
		Origin:      webAuthnOrigin,
		CrossOrigin: false,
	})
}

func createWebAuthnChallengeResponse(ctx context.Context, challenge webAuthnChallenge) (string, error) {
	clientDataJSON, err := challenge.clientDataJSON()
	if err != nil {
		return "", err
	}
	assertion, err := platformWebAuthnGetAssertion(ctx, challenge, clientDataJSON)
	if err != nil {
		return "", err
	}
	return marshalWebAuthnChallengeResponse(assertion)
}

func marshalWebAuthnChallengeResponse(assertion *webAuthnAssertion) (string, error) {
	if assertion == nil {
		return "", errors.New("x security-key assertion was empty")
	}
	payload := struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Response struct {
			ClientDataJSON    string `json:"clientDataJSON"`
			AuthenticatorData string `json:"authenticatorData"`
			Signature         string `json:"signature"`
			UserHandle        string `json:"userHandle"`
		} `json:"response"`
		ClientExtensionResults struct {
			AppID bool `json:"appid"`
		} `json:"clientExtensionResults"`
	}{
		ID:   webAuthnBase64(assertion.CredentialID),
		Type: "public-key",
	}
	payload.Response.ClientDataJSON = webAuthnBase64(assertion.ClientDataJSON)
	payload.Response.AuthenticatorData = webAuthnBase64(assertion.AuthenticatorData)
	payload.Response.Signature = webAuthnBase64(assertion.Signature)
	payload.Response.UserHandle = ""
	payload.ClientExtensionResults.AppID = true
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func webAuthnBase64(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func webAuthnDecodeBase64(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	if data, err := base64.RawURLEncoding.DecodeString(value); err == nil {
		return data, nil
	}
	return base64.URLEncoding.DecodeString(value)
}

func (jfr jetfuelLoginResponse) webAuthnChallenge() (webAuthnChallenge, bool) {
	for _, str := range jfr.strings {
		if challenge, ok := findWebAuthnChallengeInText(str); ok {
			return challenge, true
		}
	}
	if len(jfr.raw) > 0 {
		return findWebAuthnChallengeInText(string(jfr.raw))
	}
	return webAuthnChallenge{}, false
}

func findWebAuthnChallengeInText(text string) (webAuthnChallenge, bool) {
	if !strings.Contains(text, "publicKeyCredentialRequestOptions") && !strings.Contains(text, "rpId") {
		return webAuthnChallenge{}, false
	}
	for _, candidate := range extractJSONObjectsContaining(text, "publicKeyCredentialRequestOptions") {
		if challenge, ok := parseWebAuthnChallengeCandidate(candidate); ok {
			return challenge, true
		}
	}
	for _, candidate := range extractJSONObjectsContaining(text, "rpId") {
		if challenge, ok := parseWebAuthnChallengeCandidate(candidate); ok {
			return challenge, true
		}
	}
	return webAuthnChallenge{}, false
}

func parseWebAuthnChallengeCandidate(candidate string) (webAuthnChallenge, bool) {
	for _, variant := range webAuthnJSONVariants(candidate) {
		if challenge, ok := parseWebAuthnChallengeJSON([]byte(variant)); ok {
			return challenge, true
		}
	}
	return webAuthnChallenge{}, false
}

func webAuthnJSONVariants(candidate string) []string {
	candidate = strings.TrimSpace(candidate)
	var out []string
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		for _, existing := range out {
			if existing == value {
				return
			}
		}
		out = append(out, value)
	}
	add(candidate)
	if unescaped, err := url.QueryUnescape(candidate); err == nil {
		add(unescaped)
	}
	if strings.Contains(candidate, `\"`) {
		add(strings.ReplaceAll(candidate, `\"`, `"`))
	}
	var quoted string
	if err := json.Unmarshal([]byte(candidate), &quoted); err == nil {
		add(quoted)
	}
	return out
}

func parseWebAuthnChallengeJSON(data []byte) (webAuthnChallenge, bool) {
	var envelope struct {
		Challenge                         json.RawMessage         `json:"challenge"`
		PublicKeyCredentialRequestOptions *webAuthnRequestOptions `json:"publicKeyCredentialRequestOptions"`
		RPID                              string                  `json:"rpId"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return webAuthnChallenge{}, false
	}
	if envelope.PublicKeyCredentialRequestOptions != nil && envelope.PublicKeyCredentialRequestOptions.Challenge != "" {
		return webAuthnChallenge{RequestOptions: *envelope.PublicKeyCredentialRequestOptions}, true
	}
	if len(envelope.Challenge) > 0 {
		challengeData := []byte(envelope.Challenge)
		var challengeString string
		if err := json.Unmarshal(envelope.Challenge, &challengeString); err == nil {
			challengeData = []byte(challengeString)
		}
		if challenge, ok := parseWebAuthnChallengeCandidate(string(challengeData)); ok {
			return challenge, true
		}
	}
	var direct webAuthnRequestOptions
	if err := json.Unmarshal(data, &direct); err == nil && direct.Challenge != "" {
		return webAuthnChallenge{RequestOptions: direct}, true
	}
	return webAuthnChallenge{}, false
}

func extractJSONObjectsContaining(text, needle string) []string {
	var out []string
	for searchStart := 0; ; {
		idx := strings.Index(text[searchStart:], needle)
		if idx < 0 {
			return out
		}
		idx += searchStart
		for start := idx; start >= 0; start-- {
			if text[start] != '{' {
				continue
			}
			end, ok := balancedJSONEnd(text, start)
			if !ok || end <= idx {
				continue
			}
			out = append(out, text[start:end])
			break
		}
		searchStart = idx + len(needle)
	}
}

func balancedJSONEnd(text string, start int) (int, bool) {
	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(text); i++ {
		ch := text[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			switch ch {
			case '\\':
				escaped = true
			case '"':
				inString = false
			}
			continue
		}
		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1, true
			}
			if depth < 0 {
				return 0, false
			}
		}
	}
	return 0, false
}

func webAuthnChallengeSubmitForm(challengeResponse string, fields []string, sessionToken, preludeDispatchID string) url.Values {
	form := url.Values{}
	for _, field := range fields {
		lower := strings.ToLower(field)
		if strings.Contains(lower, "challenge") || strings.Contains(lower, "response") || strings.Contains(lower, "webauthn") {
			form.Set(field, challengeResponse)
		}
	}
	if form.Get("challenge_response") == "" {
		form.Set("challenge_response", challengeResponse)
	}
	if sessionToken != "" {
		form.Set("session_token", sessionToken)
	}
	if preludeDispatchID != "" {
		form.Set("prelude_dispatch_id", preludeDispatchID)
	}
	return form
}

func missingWebAuthnChallengeError(methodName string) error {
	return fmt.Errorf("%w: %s did not expose a security-key challenge", ErrWebLoginUnexpectedSubtask, methodName)
}
