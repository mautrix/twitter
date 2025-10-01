package methods

import (
	"regexp"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var (
	metaTagRegex           = regexp.MustCompile(`<meta\s+http-equiv=["']refresh["']\s+content=["'][^;]+;\s*url\s*=\s*([^"']+)["']\s*/?>`)
	migrateFormDataRegex   = regexp.MustCompile(`<form[^>]* action="([^"]+)"[^>]*>[\s\S]*?<input[^>]* name="tok" value="([^"]+)"[^>]*>[\s\S]*?<input[^>]* name="data" value="([^"]+)"[^>]*>`)
	mainScriptURLRegex     = regexp.MustCompile(`https:\/\/(?:[A-Za-z0-9.-]+)\/responsive-web\/client-web\/main\.[0-9A-Za-z]+\.js`)
	bearerTokenRegex       = regexp.MustCompile(`(Bearer\s[A-Za-z0-9%]{16,})`)
	guestTokenRegex        = regexp.MustCompile(`gt=([0-9]+)`)
	verificationTokenRegex = regexp.MustCompile(`meta name="twitter-site-verification" content="([^"]+)"`)
	countryCodeRegex       = regexp.MustCompile(`"country":\s*"([A-Z]{2})"`)
	ondemandSRegex         = regexp.MustCompile(`"ondemand.s":"([a-f0-9]+)"`)
	variableIndexesRegex   = regexp.MustCompile(`const\[\w{1,2},\w{1,2}]=\[.+?\(\w{1,2}\[(\d{1,2})],16\).+?\(\w{1,2}\[(\d{1,2})],16\).+?\(\w{1,2}\[(\d{1,2})],16\).+?\(\w{1,2}\[(\d{1,2})],16\)`)
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
	return bearerTokenRegex.FindSubmatch(js)
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

func ParseOndemandS(html string) string {
	match := ondemandSRegex.FindStringSubmatch(html)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}
