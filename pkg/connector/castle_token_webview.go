package connector

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

const castleTokenJSConfigPlaceholder = "__MAUTRIX_TWITTER_CASTLE_CONFIG__"
const browserHeadersJSPlaceholder = "__MAUTRIX_TWITTER_BROWSER_HEADERS__"

const browserHeadersExtractJS = `function captureBrowserHeaders() {
  function quoteClientHint(value) {
    return '"' + String(value).replace(/\\/g, "\\\\").replace(/"/g, '\\"') + '"';
  }
  const headers = { browser_user_agent: String(navigator.userAgent || "") };
  const userAgentData = navigator.userAgentData;
  if (!userAgentData) {
    return headers;
  }
  const brands = Array.from(userAgentData.brands || []);
  if (brands.length > 0) {
    headers.browser_sec_ch_ua = brands.map(item =>
      quoteClientHint(item.brand) + ";v=" + quoteClientHint(item.version)
    ).join(", ");
  }
  if (userAgentData.platform) {
    headers.browser_sec_ch_ua_platform = quoteClientHint(userAgentData.platform);
  }
  headers.browser_sec_ch_ua_mobile = userAgentData.mobile ? "?1" : "?0";
  return headers;
}`

const browserIdentityExtractJS = `(() => {
` + browserHeadersExtractJS + `
  const result = captureBrowserHeaders();
  const chromium = Array.from(navigator.userAgentData?.brands || [])
    .find(item => item.brand === "Chromium");
  if (chromium && /^\d+$/.test(String(chromium.version))) {
    result.browser_user_agent = result.browser_user_agent.replace(/\bChrome\/\d+\.\d+\.\d+\.\d+\b/,
      "Chrome/" + chromium.version + ".0.0.0");
  }
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
	if strings.Count(script, castleTokenJSConfigPlaceholder) != 1 || strings.Count(script, browserHeadersJSPlaceholder) != 1 {
		panic("Castle token extraction script must contain exactly one config and browser header placeholder")
	}
	return strings.NewReplacer(castleTokenJSConfigPlaceholder, string(config), browserHeadersJSPlaceholder, browserHeadersExtractJS).Replace(script)
}
