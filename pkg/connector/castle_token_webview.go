package connector

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

const castleTokenJSConfigPlaceholder = "__MAUTRIX_TWITTER_CASTLE_CONFIG__"

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
