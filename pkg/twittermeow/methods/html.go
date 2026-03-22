package methods

import (
	"regexp"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var (
	metaTagRegex           = regexp.MustCompile(`<meta\s+http-equiv=["']refresh["']\s+content=["'][^;]+;\s*url\s*=\s*([^"']+)["']\s*/?>`)
	migrateFormDataRegex   = regexp.MustCompile(`<form[^>]* action="([^"]+)"[^>]*>[\s\S]*?<input[^>]* name="tok" value="([^"]+)"[^>]*>[\s\S]*?<input[^>]* name="data" value="([^"]+)"[^>]*>`)
	mainScriptURLRegex     = regexp.MustCompile(`https:\/\/(?:[A-Za-z0-9.-]+)\/responsive-web\/client-web\/main\.[0-9A-Za-z]+\.js`)
	bearerTokenRegex       = regexp.MustCompile(`Bearer\s[A-Za-z0-9%]{16,}`)
	guestTokenRegex        = regexp.MustCompile(`gt=([0-9]+)`)
	verificationTokenRegex = regexp.MustCompile(`meta name="twitter-site-verification" content="([^"]+)"`)
	countryCodeRegex       = regexp.MustCompile(`"country":\s*"([A-Z]{2})"`)
	ondemandSChunkIDRegex  = regexp.MustCompile(`(\d+):"ondemand\.s"`)
	variableIndexesRegex   = regexp.MustCompile(`\[.+?\(\w{1,2}\[(\d{1,2})],16\).+?\(\w{1,2}\[(\d{1,2})],16\).+?\(\w{1,2}\[(\d{1,2})],16\).+?\(\w{1,2}\[(\d{1,2})],16\)`)
)

func ParseMigrateURL(html string) (string, bool) {
	match := metaTagRegex.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1], true
	}
	return "", false
}

func ParseMigrateRequestData(html string) (string, *payload.MigrationRequestPayload) {
	match := migrateFormDataRegex.FindStringSubmatch(html)
	if len(match) < 4 {
		return "", nil
	}

	return match[1], &payload.MigrationRequestPayload{Tok: match[2], Data: match[3]}
}

func ParseMainScriptURL(html string) string {
	match := mainScriptURLRegex.FindStringSubmatch(html)
	if len(match) < 1 {
		return ""
	}
	return match[0]
}

func ParseBearerToken(js []byte) [][]byte {
	return bearerTokenRegex.FindAll(js, -1)
}

func ParseVariableIndexes(js []byte) [][]byte {
	return variableIndexesRegex.FindSubmatch(js)
}

func ParseGuestToken(html string) string {
	match := guestTokenRegex.FindStringSubmatch(html)
	if len(match) < 1 {
		return ""
	}
	return match[1]
}

func ParseVerificationToken(html string) string {
	match := verificationTokenRegex.FindStringSubmatch(html)
	if len(match) < 1 {
		return ""
	}
	return match[1]
}

func ParseCountry(html string) string {
	match := countryCodeRegex.FindStringSubmatch(html)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func ParseOndemandSURLFromScript(js []byte) string {
	chunkIDMatch := ondemandSChunkIDRegex.FindSubmatchIndex(js)
	if len(chunkIDMatch) < 4 {
		return ""
	}

	chunkID := string(js[chunkIDMatch[2]:chunkIDMatch[3]])
	hashRegex := regexp.MustCompile(`(?:^|[,{])` + regexp.QuoteMeta(chunkID) + `:"([0-9a-f]+)"`)
	jsAfterNameMap := js[chunkIDMatch[1]:]
	hashMatch := hashRegex.FindSubmatchIndex(jsAfterNameMap)
	if len(hashMatch) < 4 {
		return ""
	}

	hash := string(jsAfterNameMap[hashMatch[2]:hashMatch[3]])
	return "https://abs.twimg.com/responsive-web/client-web/ondemand.s." + hash + "a.js"
}
