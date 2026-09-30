package connector

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

const castleTokenJSConfigPlaceholder = "__MAUTRIX_TWITTER_CASTLE_CONFIG__"

const browserIdentityExtractJS = `(() => {
  let userAgent = String(navigator.userAgent || "");
  const chromium = Array.from(navigator.userAgentData?.brands || [])
    .find(item => item.brand === "Chromium");
  if (chromium && /^\d+$/.test(String(chromium.version))) {
    userAgent = userAgent.replace(/\bChrome\/\d+\.\d+\.\d+\.\d+\b/,
      "Chrome/" + chromium.version + ".0.0.0");
  }
  const result = { browser_user_agent: userAgent };
  globalThis.__BEEP_BEEP_AUTH_RESULTS__ = result;
  return result;
})()`

//go:embed castle_token.js
var castleTokenExtractJSSource string

type castleTokenJSConfig struct {
	ScriptURL   string   `json:"scriptURL"`
	PublicKey   string   `json:"publicKey"`
	CookieNames []string `json:"cookieNames"`
	ContextURL  string   `json:"contextURL"`
	Identifier  string   `json:"identifier"`
	BatchSize   int      `json:"castleTokenBatchSize"`
}

func castleTokenExtractJS(info twittermeow.JetfuelCastleTokenInfo, identifier string) string {
	config, err := json.Marshal(castleTokenJSConfig{
		ScriptURL:   info.ScriptURL,
		PublicKey:   info.PublicKey,
		CookieNames: castleTokenCookieNames,
		ContextURL:  castleTokenContextURL,
		Identifier:  identifier,
		BatchSize:   castleTokenBatchSize,
	})
	if err != nil {
		panic(fmt.Errorf("marshal Castle token extraction config: %w", err))
	}

	script := strings.TrimRight(castleTokenExtractJSSource, "\r\n")
	if strings.Count(script, castleTokenJSConfigPlaceholder) != 1 {
		panic("Castle token extraction script must contain exactly one config placeholder")
	}
	return strings.Replace(script, castleTokenJSConfigPlaceholder, string(config), 1)
}
