package twittermeow

import (
	"errors"
	"net/url"
	"strings"
)

var ErrJetfuelCastleTokenRequired = errors.New("x login needs a Castle token from the client webview")

type JetfuelCastleTokenInfo struct {
	ScriptURL string
	PublicKey string
}

func (info JetfuelCastleTokenInfo) IsValid() bool {
	return strings.TrimSpace(info.ScriptURL) != "" && strings.TrimSpace(info.PublicKey) != ""
}

func (c *Client) JetfuelCastleTokenInfo() JetfuelCastleTokenInfo {
	return c.jetfuelCastleInfo
}

// SetNextJetfuelCastleTokens provides browser-generated Castle tokens for the
// next Jetfuel form submissions. Tokens are consumed in order and only once.
func (c *Client) SetNextJetfuelCastleTokens(tokens []string) {
	c.jetfuelCastleTokenMu.Lock()
	defer c.jetfuelCastleTokenMu.Unlock()
	c.jetfuelCastleTokens = c.jetfuelCastleTokens[:0]
	for _, token := range tokens {
		if token = strings.TrimSpace(token); token != "" {
			c.jetfuelCastleTokens = append(c.jetfuelCastleTokens, token)
		}
	}
}

func (c *Client) takeNextJetfuelCastleToken() string {
	c.jetfuelCastleTokenMu.Lock()
	defer c.jetfuelCastleTokenMu.Unlock()
	if len(c.jetfuelCastleTokens) == 0 {
		return ""
	}
	token := c.jetfuelCastleTokens[0]
	copy(c.jetfuelCastleTokens, c.jetfuelCastleTokens[1:])
	c.jetfuelCastleTokens = c.jetfuelCastleTokens[:len(c.jetfuelCastleTokens)-1]
	return token
}

func (c *Client) HasNextJetfuelCastleToken() bool {
	c.jetfuelCastleTokenMu.Lock()
	defer c.jetfuelCastleTokenMu.Unlock()
	return len(c.jetfuelCastleTokens) > 0
}

func (c *Client) addJetfuelCastleTokenToForm(form url.Values) error {
	if form.Get("$castle_token") != "" {
		return nil
	}
	token := c.takeNextJetfuelCastleToken()
	if token == "" {
		return ErrJetfuelCastleTokenRequired
	}
	form.Set("$castle_token", token)
	return nil
}
