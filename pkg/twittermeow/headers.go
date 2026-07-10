package twittermeow

import (
	"net/http"
	"strings"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
)

// These defaults bootstrap requests made before BrowserAuth captures the
// client webview's actual HTTP headers.
const BrowserName = "Chrome"
const ChromeVersion = "141"
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + ChromeVersion + ".0.0.0 Safari/537.36"
const SecCHUserAgent = `"Chromium";v="` + ChromeVersion + `", "Google Chrome";v="` + ChromeVersion + `", "Not-A.Brand";v="99"`
const OSName = "Linux"
const SecCHPlatform = `"` + OSName + `"`
const SecCHMobile = "?0"

const UDID = OSName + "/" + BrowserName

// BrowserHeaders contains the browser fingerprint captured from the client
// webview. Empty client-hint values are intentional for non-Chromium browsers.
type BrowserHeaders struct {
	UserAgent      string `json:"user_agent"`
	SecCHUserAgent string `json:"sec_ch_ua,omitempty"`
	SecCHPlatform  string `json:"sec_ch_ua_platform,omitempty"`
	SecCHMobile    string `json:"sec_ch_ua_mobile,omitempty"`
}

var BaseHeaders = http.Header{
	"Accept":             []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	"Accept-Language":    []string{"en"},
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

	headers := c.GetBaseHeaders()
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

func normalizeBrowserHeader(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 1024 {
		return ""
	}
	for _, char := range value {
		if char < 0x20 || char == 0x7f {
			return ""
		}
	}
	return value
}

func (headers BrowserHeaders) normalized() BrowserHeaders {
	headers.UserAgent = normalizeBrowserHeader(headers.UserAgent)
	headers.SecCHUserAgent = normalizeBrowserHeader(headers.SecCHUserAgent)
	headers.SecCHPlatform = normalizeBrowserHeader(headers.SecCHPlatform)
	headers.SecCHMobile = normalizeBrowserHeader(headers.SecCHMobile)
	if headers.SecCHMobile != "?0" && headers.SecCHMobile != "?1" {
		headers.SecCHMobile = ""
	}
	if headers.UserAgent == "" {
		return BrowserHeaders{}
	}
	return headers
}

// SetBrowserHeaders installs a webview-captured browser fingerprint. It returns
// false when the required User-Agent is missing or unsafe to use as a header.
func (c *Client) SetBrowserHeaders(headers BrowserHeaders) bool {
	headers = headers.normalized()
	if headers.UserAgent == "" {
		return false
	}
	c.browserHeaders = headers
	return true
}

func (c *Client) GetBrowserHeaders() BrowserHeaders {
	return c.browserHeaders
}

func (c *Client) GetBaseHeaders() http.Header {
	headers := BaseHeaders.Clone()
	browser := c.browserHeaders
	if browser.UserAgent == "" {
		return headers
	}
	headers.Set("User-Agent", browser.UserAgent)
	setOptionalHeader(headers, "Sec-Ch-Ua", browser.SecCHUserAgent)
	setOptionalHeader(headers, "Sec-Ch-Ua-Platform", browser.SecCHPlatform)
	setOptionalHeader(headers, "Sec-Ch-Ua-Mobile", browser.SecCHMobile)
	return headers
}

func setOptionalHeader(headers http.Header, name, value string) {
	if value == "" {
		headers.Del(name)
	} else {
		headers.Set(name, value)
	}
}

func (c *Client) GetUserAgent() string {
	if c.browserHeaders.UserAgent != "" {
		return c.browserHeaders.UserAgent
	}
	return UserAgent
}
