package twittermeow

import (
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

const BrowserName = "Chrome"
const ChromeVersion = "131"
const ChromeVersionFull = ChromeVersion + ".0.6778.85"
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + ChromeVersion + ".0.0.0 Safari/537.36"
const SecCHUserAgent = `"Chromium";v="` + ChromeVersion + `", "Google Chrome";v="` + ChromeVersion + `", "Not-A.Brand";v="99"`
const SecCHFullVersionList = `"Chromium";v="` + ChromeVersionFull + `", "Google Chrome";v="` + ChromeVersionFull + `", "Not-A.Brand";v="99.0.0.0"`
const OSName = "Linux"
const OSVersion = "6.8.0"
const SecCHPlatform = `"` + OSName + `"`
const SecCHPlatformVersion = `"` + OSVersion + `"`
const SecCHMobile = "?0"
const SecCHModel = ""
const SecCHPrefersColorScheme = "light"

const UDID = OSName + "/" + BrowserName

var BaseHeaders = http.Header{
	"accept":             []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	"accept-language":    []string{"en-US,en;q=0.9"},
	"user-agent":         []string{UserAgent},
	"sec-ch-ua":          []string{SecCHUserAgent},
	"sec-ch-ua-platform": []string{SecCHPlatform},
	"sec-ch-ua-mobile":   []string{SecCHMobile},
	"referer":            []string{endpoints.BASE_URL + "/"},
	"origin":             []string{endpoints.BASE_URL},
	//"sec-ch-prefers-color-scheme": []string{SecCHPrefersColorScheme},
	//"sec-ch-ua-full-version-list": []string{SecCHFullVersionList},
	//"sec-ch-ua-model":             []string{SecCHModel},
	//"sec-ch-ua-platform-version":  []string{SecCHPlatformVersion},
}

type HeaderOpts struct {
	WithAuthBearer      bool
	WithCookies         bool
	WithXGuestToken     bool
	WithXTwitterHeaders bool
	WithXCsrfToken      bool
	WithXClientUUID     bool
	Referer             string
	Origin              string
	Extra               map[string]string
}

func (c *Client) buildHeaders(opts HeaderOpts) http.Header {
	if opts.Extra == nil {
		opts.Extra = make(map[string]string)
	}

	headers := BaseHeaders.Clone()
	if opts.WithCookies {
		opts.Extra["cookie"] = c.cookies.String()
	}

	if opts.WithAuthBearer {
		authTokens := c.session.GetAuthTokens()
		// check if client is authenticated here
		var bearerToken string
		if c.isAuthenticated() {
			bearerToken = authTokens.authenticated
		} else {
			bearerToken = authTokens.notAuthenticated
		}
		opts.Extra["authorization"] = bearerToken
	}

	if opts.WithXGuestToken {
		opts.Extra["x-guest-token"] = c.cookies.Get(cookies.XGuestToken)
	}

	if opts.WithXClientUUID {
		opts.Extra["x-client-uuid"] = c.session.clientUUID
	}

	if opts.WithXTwitterHeaders {
		opts.Extra["x-twitter-active-user"] = "yes"
		opts.Extra["x-twitter-client-language"] = "en"
	}

	if opts.WithXCsrfToken {
		opts.Extra["x-csrf-token"] = c.cookies.Get(cookies.XCt0)
		opts.Extra["x-twitter-auth-type"] = "OAuth2Session"
	}

	if opts.Origin != "" {
		opts.Extra["origin"] = opts.Origin
	}

	if opts.Referer != "" {
		opts.Extra["referer"] = opts.Referer
	}

	for k, v := range opts.Extra {
		headers.Set(k, v)
	}

	return headers
}
