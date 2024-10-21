package methods

import (
	"regexp"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var (
	metaTagRegex           = regexp.MustCompile(`<meta\s+http-equiv=["']refresh["']\s+content=["'][^;]+;\s*url\s*=\s*([^"']+)["']\s*/?>`)
	migrateFormDataRegex   = regexp.MustCompile(`<form[^>]* action="([^"]+)"[^>]*>[\s\S]*?<input[^>]* name="tok" value="([^"]+)"[^>]*>[\s\S]*?<input[^>]* name="data" value="([^"]+)"[^>]*>`)
	mainScriptUrlRegex     = regexp.MustCompile(`https:\/\/(?:[A-Za-z0-9.-]+)\/responsive-web\/client-web\/main\.[0-9A-Za-z]+\.js`)
	bearerTokenRegex       = regexp.MustCompile(`(Bearer\s[A-Za-z0-9%]+)`)
	guestTokenRegex        = regexp.MustCompile(`gt=([0-9]+)`)
	verificationTokenRegex = regexp.MustCompile(`meta name="twitter-site-verification" content="([^"]+)"`)
	countryCodeRegex       = regexp.MustCompile(`"country":\s*"([A-Z]{2})"`)
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
	match := mainScriptUrlRegex.FindStringSubmatch(html)
	if len(match) < 1 {
		return ""
	}
	return match[0]
}

func ParseBearerTokens(html string) []string {
	return bearerTokenRegex.FindAllString(html, -1)
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
