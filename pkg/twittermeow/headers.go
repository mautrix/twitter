package twittermeow

import (
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

const BrowserName = "Chrome"
const ChromeVersion = "141"
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + ChromeVersion + ".0.0.0 Safari/537.36"
const SecCHUserAgent = `"Chromium";v="` + ChromeVersion + `", "Google Chrome";v="` + ChromeVersion + `", "Not-A.Brand";v="99"`
const OSName = "Linux"
const SecCHPlatform = `"` + OSName + `"`
const SecCHMobile = "?0"

const UDID = OSName + "/" + BrowserName

var BaseHeaders = http.Header{
	"Accept":             []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	"Accept-Language":    []string{"en-US,en;q=0.9"},
	"User-Agent":         []string{UserAgent},
	"Sec-Ch-Ua":          []string{SecCHUserAgent},
	"Sec-Ch-Ua-Platform": []string{SecCHPlatform},
	"Sec-Ch-Ua-Mobile":   []string{SecCHMobile},
	"Referer":            []string{endpoints.BASE_URL + "/"},
	"Origin":             []string{endpoints.BASE_URL},
}

type HeaderOpts struct {
	WithAuthBearer      bool
	WithNonAuthBearer   bool
	WithCookies         bool
	WithXGuestToken     bool
	WithXTwitterHeaders bool
	WithXCsrfToken      bool
	WithXClientUUID     bool
	Referer             string
	Origin              string
	Extra               map[string]string
}

const BearerToken = "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

func (c *Client) buildHeaders(opts HeaderOpts) http.Header {
	if opts.Extra == nil {
		opts.Extra = make(map[string]string)
	}

	headers := BaseHeaders.Clone()
	if opts.WithCookies {
		opts.Extra["cookie"] = c.cookies.String()
	}

	if opts.WithAuthBearer || opts.WithNonAuthBearer {
		if c.session.bearerToken != "" {
			opts.Extra["authorization"] = c.session.bearerToken
		} else {
			opts.Extra["authorization"] = BearerToken
		}
	}

	if opts.WithXGuestToken {
		opts.Extra["x-guest-token"] = c.cookies.Get(cookies.XGuestToken)
	}

	if opts.WithXClientUUID {
		opts.Extra["x-client-uuid"] = c.session.ClientUUID
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
