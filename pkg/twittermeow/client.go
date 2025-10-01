package twittermeow

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"golang.org/x/net/proxy"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

type EventHandler func(evt types.TwitterEvent, inbox *response.TwitterInboxData) bool
type StreamEventHandler func(evt response.StreamEvent)

type Client struct {
	Logger     zerolog.Logger
	cookies    *cookies.Cookies
	session    *CachedSession
	HTTP       *http.Client
	httpProxy  func(*http.Request) (*url.URL, error)
	socksProxy proxy.Dialer

	eventHandler       EventHandler
	streamEventHandler StreamEventHandler
	onCursorChanged    func(ctx context.Context)

	jot     *JotClient
	polling *PollingClient
	stream  *StreamClient
}

func NewClient(cookies *cookies.Cookies, logger zerolog.Logger) *Client {
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

	cli.polling = cli.newPollingClient()
	cli.stream = cli.newStreamClient()
	cli.jot = cli.newJotClient()
	cli.cookies = cookies
	cli.session = &CachedSession{
		ClientUUID: uuid.NewString(),
	}

	return &cli
}

func (c *Client) GetCookieString() string {
	return c.cookies.String()
}

func (c *Client) SetSession(sess *CachedSession) {
	c.session = sess
}

func (c *Client) GetSession() *CachedSession {
	return c.session
}

func (c *Client) Connect(ctx context.Context, cached bool) {
	if c.eventHandler == nil {
		panic(ErrConnectSetEventHandler)
	}

	c.polling.startPolling(c.Logger.WithContext(ctx))
	if cached {
		c.polling.doShortCircuit()
	}
}

func (c *Client) Disconnect() {
	c.polling.stopPolling()
	c.stream.stopStream()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.polling.pollingStopped.Wait(timeoutCtx)
	if err != nil {
		c.Logger.Warn().Msg("Timed out waiting for polling to stop")
		return
	}
	err = c.stream.streamStopped.Wait(timeoutCtx)
	if err != nil {
		c.Logger.Warn().Msg("Timed out waiting for stream to stop")
		return
	}
	c.Logger.Debug().Msg("Polling and stream clients stopped")
}

func (c *Client) Logout(ctx context.Context) error {
	defer c.cookies.Set(cookies.XAuthToken, "")
	return c.loadPage(ctx, endpoints.BASE_LOGOUT_URL)
}

func (c *Client) LoadMessagesPage(ctx context.Context) (*response.InboxInitialStateResponse, *response.AccountSettingsResponse, error) {
	err := c.loadPage(ctx, endpoints.BASE_MESSAGES_URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load messages page: %w", err)
	}

	data, err := c.GetAccountSettings(ctx, payload.AccountSettingsQuery{
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
		if IsAuthError(err) {
			return nil, nil, err
		}
		c.Logger.Warn().Err(err).Msg("Failed to get account settings")
		data = &response.AccountSettingsResponse{}
	}

	initialInboxState, err := c.GetInitialInboxState(ctx, ptr.Ptr(payload.DMRequestQuery{}.Default()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get initial inbox state: %w", err)
	}

	c.session.PollingCursor = initialInboxState.InboxInitialState.Cursor
	c.session.InitializedAt = time.Now()

	c.Logger.Info().
		Str("screen_name", data.ScreenName).
		Str("initial_inbox_cursor", initialInboxState.InboxInitialState.Cursor).
		Msg("Successfully loaded and authenticated as user")

	return initialInboxState, data, nil
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

func (c *Client) SetEventHandler(handler EventHandler, streamHandler StreamEventHandler, onCursorChanged func(context.Context)) {
	c.eventHandler = handler
	c.streamEventHandler = streamHandler
	c.onCursorChanged = onCursorChanged
}

func (c *Client) fetchScript(ctx context.Context, url string) ([]byte, error) {
	extraHeaders := map[string]string{
		"accept":         "*/*",
		"sec-fetch-site": "cross-site",
		"sec-fetch-mode": "cors",
		"sec-fetch-dest": "script",
		"origin":         endpoints.BASE_URL,
	}
	_, scriptRespBody, err := c.MakeRequest(ctx, url, http.MethodGet, c.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: endpoints.BASE_URL + "/"}), nil, types.ContentTypeNone)
	return scriptRespBody, err
}

func (c *Client) fetchAndParseMainScript(ctx context.Context, scriptURL string) {
	scriptRespBody, err := c.fetchScript(ctx, scriptURL)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Failed to fetch main script")
		return
	}
	authToken := methods.ParseBearerToken(scriptRespBody)
	if len(authToken) == 0 {
		zerolog.Ctx(ctx).Warn().Msg("No bearer token found in main script, using default")
		return
	}
	c.session.bearerToken = string(authToken[1])
	if c.session.bearerToken != BearerToken {
		zerolog.Ctx(ctx).Warn().Str("bearer_token", c.session.bearerToken).
			Msg("Hardcoded token doesn't match fetched one")
	}
}

func (c *Client) fetchAndParseSScript(ctx context.Context, scriptURL string) (*[4]int, error) {
	scriptRespBody, err := c.fetchScript(ctx, scriptURL)
	if err != nil {
		return nil, err
	}

	byteIndexes := methods.ParseVariableIndexes(scriptRespBody)
	if len(byteIndexes) < 5 {
		return nil, fmt.Errorf("failed to find variable indexes")
	}
	var indexes [4]int
	for i := 0; i < 4; i++ {
		index, err := strconv.Atoi(string(byteIndexes[i+1]))
		if err != nil {
			return nil, fmt.Errorf("failed to parse variable index %d (%s): %w", i, byteIndexes[i+1], err)
		}
		indexes[i] = index
	}
	return &indexes, nil
}

var nonNumbers = regexp.MustCompile(`\D+`)

func (c *Client) parseMainPageHTML(ctx context.Context, mainPageResp *http.Response, mainPageHTML string) error {
	country := methods.ParseCountry(mainPageHTML)
	if country == "" {
		return fmt.Errorf("country code not found (HTTP %d)", mainPageResp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(mainPageHTML))
	if err != nil {
		return fmt.Errorf("failed to parse main page html: %w", err)
	}
	verificationToken, ok := doc.Find("meta[name^=tw]").Attr("content")
	if !ok {
		return fmt.Errorf("failed to find meta tags in main page html")
	}
	var loadingAnims = new([4][16][11]int)
	idx := 0
	doc.Find("svg[id^=loading-x-anim]").Each(func(i int, svg *goquery.Selection) {
		if idx >= 4 {
			idx++
			return
		}
		pathVal, ok := svg.Find("path").Eq(1).Attr("d")
		if !ok {
			return
		}
		sets := strings.Split(pathVal[9:], "C")
		if len(sets) != 16 {
			return
		}
		var numSets [16][11]int
		for i, set := range sets {
			numbers := strings.Split(strings.TrimSpace(nonNumbers.ReplaceAllString(set, " ")), " ")
			if len(numbers) != 11 {
				return
			}
			for j, num := range numbers {
				parsed, err := strconv.Atoi(num)
				if err != nil {
					return
				}
				numSets[i][j] = parsed
			}
		}
		loadingAnims[idx] = numSets
		idx++
	})
	if idx != 4 {
		c.Logger.Warn().Int("found_count", idx).Msg("Didn't find 4 loading animations in main page HTML")
		loadingAnims = nil
	} else {
		c.Logger.Trace().
			Str("verification_token", verificationToken).
			Any("loading_animations", loadingAnims[:]).
			Msg("Found loading animations and verification token")
	}

	c.session.Country = country
	c.session.VerificationToken = verificationToken
	c.session.loadingAnims = loadingAnims

	guestToken := methods.ParseGuestToken(mainPageHTML)
	if guestToken == "" {
		if c.cookies.IsCookieEmpty(cookies.XGuestToken) && !c.IsLoggedIn() {
			// most likely means your cookies are invalid / expired
			return fmt.Errorf("guest token not found (HTTP %d)", mainPageResp.StatusCode)
		}
	} else {
		c.cookies.Set(cookies.XGuestToken, guestToken)
	}

	mainScriptURL := methods.ParseMainScriptURL(mainPageHTML)
	if mainScriptURL == "" {
		zerolog.Ctx(ctx).Warn().Int("status_code", mainPageResp.StatusCode).Msg("Main script URL not found in main page HTML")
	} else {
		c.fetchAndParseMainScript(ctx, mainScriptURL)
	}

	ondemandS := methods.ParseOndemandS(mainPageHTML)
	if ondemandS == "" {
		c.Logger.Warn().Msg("ondemand.s not found in main page HTML")
	} else if indexes, err := c.fetchAndParseSScript(ctx, fmt.Sprintf("https://abs.twimg.com/responsive-web/client-web/ondemand.s.%sa.js", ondemandS)); err != nil {
		c.Logger.Warn().Err(err).Msg("Failed to fetch and parse s script")
	} else {
		c.session.variableIndexes = indexes
		c.Logger.Trace().
			Any("variable_indexes", indexes[:]).
			Msg("Found variable indexes")
	}

	c.calculateAnimationToken()

	return nil
}

func (c *Client) performJotClientEvent(ctx context.Context, category payload.JotLoggingCategory, debug bool, body []interface{}) error {
	if c.jot == nil {
		return nil
	}
	return c.jot.sendClientLoggingEvent(ctx, category, debug, body)
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
	Headers        map[string]string
}

func (c *Client) makeAPIRequest(ctx context.Context, apiRequestOpts apiRequestOpts) (*http.Response, []byte, error) {
	clientTransactionID, err := crypto.SignTransaction(c.session.AnimationToken, c.session.VerificationToken, apiRequestOpts.URL, apiRequestOpts.Method)
	if err != nil {
		c.Logger.Trace().Err(err).Msg("Failed to create client transaction ID")
		clientTransactionID = base64.RawStdEncoding.EncodeToString([]byte("e:"))
	}

	extraHeaders := map[string]string{
		"x-client-transaction-id": clientTransactionID,
		"accept":                  "*/*",
		"sec-fetch-dest":          "empty",
		"sec-fetch-mode":          "cors",
		"sec-fetch-site":          "same-origin",
	}
	for k, v := range apiRequestOpts.Headers {
		extraHeaders[k] = v
	}

	headerOpts := HeaderOpts{
		WithAuthBearer:      true,
		WithCookies:         true,
		WithXTwitterHeaders: true,
		WithXCsrfToken:      true,
		Referer:             apiRequestOpts.Referer,
		Origin:              apiRequestOpts.Origin,
		Extra:               extraHeaders,
		WithXClientUUID:     apiRequestOpts.WithClientUUID,
	}
	headers := c.buildHeaders(headerOpts)

	return c.MakeRequest(ctx, apiRequestOpts.URL, apiRequestOpts.Method, headers, apiRequestOpts.Body, apiRequestOpts.ContentType)
}

func (c *Client) SetActiveConversation(conversationID string) {
	c.polling.SetActiveConversation(conversationID)
	c.stream.startOrUpdateEventStream(conversationID)
}

func (c *Client) PollConversation(conversationID string) {
	c.polling.pollConversation(conversationID)
}
