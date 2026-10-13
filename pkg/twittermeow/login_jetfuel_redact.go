package twittermeow

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"go.mau.fi/util/redact"
)

// Jetfuel strings carry no keys and its prompts can echo the user's identifier or name,
// so only the protocol structure is kept and all other text is hashed.
var jetfuelRedactPolicy = func() redact.Policy {
	p := redact.MakeDefaultPolicy()
	p.KeepPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^.?/onboarding/web(?:/[a-z0-9_]+)*(?:\?mode=[a-z]+)?(?:#/s/[a-z0-9_]+(?:/r-[a-z0-9]+)?)?$`),
		regexp.MustCompile(`^[^A-Za-z0-9]?[a-z]+(?:_[a-z]+)+$`),
		regexp.MustCompile(`^(?:password|phone|email|code|token|response|text|LoginAcid)$`),
	}
	return p
}()

func redactJetfuelStrings(body []byte, identifier string) []string {
	policy := jetfuelRedactPolicy
	if identifier = strings.ToLower(strings.TrimSpace(identifier)); len(identifier) >= 3 {
		policy.ClassifyString = func(s string) redact.Decision {
			if strings.Contains(strings.ToLower(s), identifier) {
				return redact.Hash
			}
			return redact.Auto
		}
	}
	strs := extractJetfuelStrings(body)
	for i, str := range strs {
		strs[i] = policy.String(str)
	}
	return strs
}

func (c *Client) logRedactedJetfuelResponse(ctx context.Context, method, path string, reqBody []byte, resp *http.Response, body []byte) {
	logger := c.Logger
	if contextLogger := zerolog.Ctx(ctx); contextLogger != nil && contextLogger.GetLevel() != zerolog.Disabled {
		logger = *contextLogger
	}
	evt := logger.Debug()
	if !evt.Enabled() {
		return
	}
	if form, err := url.ParseQuery(string(reqBody)); err == nil {
		for _, key := range []string{"username_or_email", "username"} {
			if value := strings.TrimSpace(form.Get(key)); value != "" {
				c.jetfuelLoginIdentifier = value
			}
		}
	}
	evt = evt.
		Str("method", method).
		Str("path", path).
		Int("response_bytes", len(body)).
		Strs("response_strings", redactJetfuelStrings(body, c.jetfuelLoginIdentifier))
	if resp != nil {
		evt = evt.Int("status_code", resp.StatusCode)
	}
	evt.Msg("Redacted Jetfuel login response")
}
