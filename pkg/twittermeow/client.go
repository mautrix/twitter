package twittermeow

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"

	"github.com/rs/zerolog"
	"golang.org/x/net/proxy"
)

type ClientOpts struct {
	PollingInterval *time.Duration
	Cookies         *cookies.Cookies
	Session         *SessionLoader
	EventHandler    EventHandler
	WithJOTClient   bool
}
type EventHandler func(evt any)
type Client struct {
	Logger       zerolog.Logger
	cookies      *cookies.Cookies
	session      *SessionLoader
	HTTP         *http.Client
	httpProxy    func(*http.Request) (*url.URL, error)
	socksProxy   proxy.Dialer
	eventHandler EventHandler

	jot     *JotClient
	polling *PollingClient
}

func NewClient(opts *ClientOpts, logger zerolog.Logger) *Client {
	cli := Client{
		HTTP: &http.Client{
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 40 * time.Second,
				ForceAttemptHTTP2:     true,
			},
			Timeout: 60 * time.Second,
		},
		Logger: logger,
	}

	cli.polling = cli.newPollingClient(opts.PollingInterval)

	if opts.WithJOTClient {
		cli.jot = cli.newJotClient()
	}

	if opts.EventHandler != nil {
		cli.SetEventHandler(opts.EventHandler)
	}

	if opts.Cookies != nil {
		cli.cookies = opts.Cookies
	} else {
		cli.cookies = cookies.NewCookies(nil)
	}

	if opts.Session != nil {
		cli.session = opts.Session
	} else {
		cli.session = cli.newSessionLoader()
	}

	return &cli
}

func (c *Client) GetCookieString() string {
	return c.cookies.String()
}

func (c *Client) Connect() error {
	if c.eventHandler == nil {
		return ErrConnectPleaseSetEventHandler
	}

	if !c.isAuthenticated() {
		return ErrNotAuthenticatedYet
	}

	return c.polling.startPolling()
}

func (c *Client) Disconnect() error {
	return c.polling.stopPolling()
}

func (c *Client) Logout() (bool, error) {
	if !c.isAuthenticated() {
		return false, ErrNotAuthenticatedYet
	}
	err := c.session.LoadPage(endpoints.BASE_LOGOUT_URL)
	if err != nil {
		return false, err
	}
	c.cookies.Set(cookies.XAuthToken, "")
	return true, nil
}

func (c *Client) LoadMessagesPage() (*response.XInboxData, *response.AccountSettingsResponse, error) {
	err := c.session.LoadPage(endpoints.BASE_MESSAGES_URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load messages page: %w", err)
	}

	data, err := c.GetAccountSettings(payload.AccountSettingsQuery{
		IncludeExtSharingAudiospacesListeningDataWithFollowers: true,
		IncludeMentionFilter:        true,
		IncludeNSFWUserFlag:         true,
		IncludeNSFWAdminFlag:        true,
		IncludeRankedTimeline:       true,
		IncludeAltTextCompose:       true,
		Ext:                         "ssoConnections",
		IncludeCountryCode:          true,
		IncludeExtDMNSFWMediaFilter: true,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account settings: %w", err)
	}

	initialInboxState, err := c.GetInitialInboxState(ptr.Ptr(payload.DMRequestQuery{}.Default()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get initial inbox state: %w", err)
	}

	c.session.SetCurrentUser(data)
	c.polling.SetCurrentCursor(initialInboxState.InboxInitialState.Cursor)

	c.Logger.Info().
		Str("screen_name", data.ScreenName).
		Str("initial_inbox_cursor", initialInboxState.InboxInitialState.Cursor).
		Msg("Successfully loaded and authenticated as user")

	return &initialInboxState.InboxInitialState, data, nil
}

func (c *Client) GetCurrentUser() *response.AccountSettingsResponse {
	return c.session.GetCurrentUser()
}

func (c *Client) GetCurrentUserID() string {
	twid := c.cookies.Get(cookies.XTwid)
	return strings.Replace(strings.Replace(twid, "u%3D", "", -1), "u=", "", -1)
}

func (c *Client) SetProxy(proxyAddr string) error {
	proxyParsed, err := url.Parse(proxyAddr)
	if err != nil {
		return err
	}

	if proxyParsed.Scheme == "http" || proxyParsed.Scheme == "https" {
		c.httpProxy = http.ProxyURL(proxyParsed)
		c.HTTP.Transport.(*http.Transport).Proxy = c.httpProxy
	} else if proxyParsed.Scheme == "socks5" {
		c.socksProxy, err = proxy.FromURL(proxyParsed, &net.Dialer{Timeout: 20 * time.Second})
		if err != nil {
			return err
		}
		c.HTTP.Transport.(*http.Transport).DialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return c.socksProxy.Dial(network, addr)
		}
		contextDialer, ok := c.socksProxy.(proxy.ContextDialer)
		if ok {
			c.HTTP.Transport.(*http.Transport).DialContext = contextDialer.DialContext
		}
	}

	c.Logger.Debug().
		Str("scheme", proxyParsed.Scheme).
		Str("host", proxyParsed.Host).
		Msg("Using proxy")
	return nil
}

func (c *Client) IsLoggedIn() bool {
	return !c.cookies.IsCookieEmpty(cookies.XAuthToken)
}

func (c *Client) isAuthenticated() bool {
	return c.session.isAuthenticated()
}

func (c *Client) SetEventHandler(handler EventHandler) {
	c.eventHandler = handler
}

func (c *Client) fetchAndParseMainScript(scriptURL string) error {
	extraHeaders := map[string]string{
		"accept":         "*/*",
		"sec-fetch-site": "cross-site",
		"sec-fetch-mode": "cors",
		"sec-fetch-dest": "script",
		"origin":         endpoints.BASE_URL,
	}
	_, scriptRespBody, err := c.MakeRequest(scriptURL, http.MethodGet, c.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: endpoints.BASE_URL + "/"}), nil, types.NONE)
	if err != nil {
		return err
	}

	scriptText := string(scriptRespBody)

	authTokens := methods.ParseBearerTokens(scriptText)
	if len(authTokens) < 2 {
		return fmt.Errorf("failed to find auth tokens in main script response body")
	}

	authenticatedToken, notAuthenticatedToken := authTokens[0], authTokens[1]
	c.session.SetAuthTokens(authenticatedToken, notAuthenticatedToken)

	return nil
}

func (c *Client) parseMainPageHTML(mainPageResp *http.Response, mainPageHTML string) error {
	country := methods.ParseCountry(mainPageHTML)
	if country == "" {
		return fmt.Errorf("failed to find session country by regex in redirected html response body (response_body=%s, status_code=%d)", mainPageHTML, mainPageResp.StatusCode)
	}

	verificationToken := methods.ParseVerificationToken(mainPageHTML)
	if verificationToken == "" {
		return fmt.Errorf("failed to find twitter verification token by regex in redirected html response body (response_body=%s, status_code=%d)", mainPageHTML, mainPageResp.StatusCode)
	}

	c.session.SetCountry(country)
	c.session.SetVerificationToken(verificationToken)

	guestToken := methods.ParseGuestToken(mainPageHTML)
	if guestToken == "" {
		if c.cookies.IsCookieEmpty(cookies.XGuestToken) && !c.IsLoggedIn() {
			// most likely means your cookies are invalid / expired
			return fmt.Errorf("failed to find guest token by regex in redirected html response body (response_body=%s, status_code=%d)", mainPageHTML, mainPageResp.StatusCode)
		}
	} else {
		c.cookies.Set(cookies.XGuestToken, guestToken)
	}

	mainScriptURL := methods.ParseMainScriptURL(mainPageHTML)
	if mainScriptURL == "" {
		return fmt.Errorf("failed to find main script url by regex in redirected html response body (response_body=%s, status_code=%d)", mainPageHTML, mainPageResp.StatusCode)
	}

	err := c.fetchAndParseMainScript(mainScriptURL)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) performJotClientEvent(category payload.JotLoggingCategory, debug bool, body []interface{}) error {
	if c.jot == nil {
		return nil
	}
	return c.jot.sendClientLoggingEvent(category, debug, body)
}

func (c *Client) enableRedirects() {
	c.HTTP.CheckRedirect = nil
}

func (c *Client) disableRedirects() {
	c.HTTP.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return ErrRedirectAttempted
	}
}

type apiRequestOpts struct {
	URL            string
	Referer        string
	Origin         string
	Method         string
	Body           []byte
	ContentType    types.ContentType
	WithClientUUID bool
}

func (c *Client) makeAPIRequest(apiRequestOpts apiRequestOpts) (*http.Response, []byte, error) {
	clientTransactionID, err := crypto.SignTransaction(c.session.verificationToken, apiRequestOpts.URL, apiRequestOpts.Method)
	if err != nil {
		return nil, nil, err
	}

	headerOpts := HeaderOpts{
		WithAuthBearer:      true,
		WithCookies:         true,
		WithXTwitterHeaders: true,
		WithXCsrfToken:      true,
		Referer:             apiRequestOpts.Referer,
		Origin:              apiRequestOpts.Origin,
		Extra: map[string]string{
			"x-client-transaction-id": clientTransactionID,
			"accept":                  "*/*",
			"sec-fetch-dest":          "empty",
			"sec-fetch-mode":          "cors",
			"sec-fetch-site":          "same-site",
		},
		WithXClientUUID: apiRequestOpts.WithClientUUID,
	}
	headers := c.buildHeaders(headerOpts)

	return c.MakeRequest(apiRequestOpts.URL, apiRequestOpts.Method, headers, apiRequestOpts.Body, apiRequestOpts.ContentType)
}
