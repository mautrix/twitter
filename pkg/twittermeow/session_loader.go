package twittermeow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"

	"github.com/google/go-querystring/query"
)

var (
	errCookieGuestIDNotFound = errors.New("failed to retrieve and set 'guest_id' cookie from Set-Cookie response headers")
)

// retrieved from main page resp, its a 2 year old timestamp; looks constant
const fetchedTime = 1661971138705

type CachedSession struct {
	InitializedAt time.Time `json:"initialized_at"`
	LastSaved     time.Time `json:"last_saved"`
	PollingCursor string    `json:"polling_cursor"`

	VerificationToken string `json:"verification_token"`
	AnimationToken    string `json:"animation_token"`
	Country           string `json:"country"`
	ClientUUID        string `json:"client_uuid"`

	bearerToken string

	loadingAnims    *[4][16][11]int
	variableIndexes *[4]int
}

func (c *Client) loadPage(ctx context.Context, url string) error {
	mainPageURL, err := neturl.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse URL %q: %w", url, err)
	}
	extraHeaders := map[string]string{
		"upgrade-insecure-requests": "1",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"sec-fetch-dest":            "document",
	}
	mainPageResp, mainPageRespBody, err := c.MakeRequest(ctx, url, http.MethodGet, c.buildHeaders(HeaderOpts{Extra: extraHeaders, WithCookies: true}), nil, types.ContentTypeNone)
	if err != nil {
		return fmt.Errorf("failed to send main page request: %w", err)
	}

	c.cookies.UpdateFromResponse(mainPageResp)
	if c.cookies.IsCookieEmpty(cookies.XGuestID) {
		c.Logger.Err(errCookieGuestIDNotFound).Msg("No GuestID found in response headers")
		return errCookieGuestIDNotFound
	}

	mainPageHTML := string(mainPageRespBody)
	migrationURL, migrationRequired := methods.ParseMigrateURL(mainPageHTML)
	if migrationRequired {
		c.Logger.Debug().Msg("Migrating session from twitter.com")
		extraHeaders = map[string]string{
			"upgrade-insecure-requests": "1",
			"sec-fetch-site":            "cross-site",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-dest":            "document",
		}
		migrationPageResp, migrationPageRespBody, err := c.MakeRequest(ctx, migrationURL, http.MethodGet, c.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: fmt.Sprintf("https://%s/", mainPageURL.Host), WithCookies: true}), nil, types.ContentTypeNone)
		if err != nil {
			return fmt.Errorf("failed to send migration request: %w", err)
		}

		migrationPageHTML := string(migrationPageRespBody)
		migrationFormURL, migrationFormPayload := methods.ParseMigrateRequestData(migrationPageHTML)
		if migrationFormPayload != nil {
			migrationForm, err := query.Values(migrationFormPayload)
			if err != nil {
				return fmt.Errorf("failed to parse migration form data: %w", err)
			}
			migrationPayload := []byte(migrationForm.Encode())
			extraHeaders["origin"] = endpoints.TWITTER_BASE_URL

			c.disableRedirects()
			mainPageResp, _, err = c.MakeRequest(ctx, migrationFormURL, http.MethodPost, c.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: endpoints.TWITTER_BASE_URL + "/", WithCookies: true}), migrationPayload, types.ContentTypeForm)
			if err == nil || !errors.Is(err, ErrRedirectAttempted) {
				return fmt.Errorf("failed to make request to main page, server did not respond with a redirect response")
			}
			c.enableRedirects()
			c.cookies.UpdateFromResponse(mainPageResp) // update the cookies received from the redirected response headers

			migrationFormURL = endpoints.BASE_URL + mainPageResp.Header.Get("Location")
			mainPageResp, mainPageRespBody, err = c.MakeRequest(ctx, migrationFormURL, http.MethodGet, c.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: endpoints.TWITTER_BASE_URL + "/", WithCookies: true}), migrationPayload, types.ContentTypeForm)
			if err != nil {
				return fmt.Errorf("failed to send main page request after migration: %w", err)
			}

			mainPageHTML := string(mainPageRespBody)
			err = c.parseMainPageHTML(ctx, mainPageResp, mainPageHTML)
			if err != nil {
				return fmt.Errorf("failed to parse main page HTML after migration: %w", err)
			}

			err = c.doInitialClientLoggingEvents(ctx)
			if err != nil {
				return fmt.Errorf("failed to perform initial client logging events after migration: %w", err)
			}

		} else {
			return fmt.Errorf("failed to find form request data in migration response (HTTP %d)", migrationPageResp.StatusCode)
		}
	} else {
		// most likely means... already authenticated
		mainPageHTML := string(mainPageRespBody)
		err = c.parseMainPageHTML(ctx, mainPageResp, mainPageHTML)
		if err != nil {
			return fmt.Errorf("failed to parse main page HTML: %w", err)
		}

		err = c.doInitialClientLoggingEvents(ctx)
		if err != nil {
			return fmt.Errorf("failed to perform initial client logging events after migration: %w", err)
		}
	}
	return nil
}

func (c *Client) doCookiesMetaDataLoad(ctx context.Context) error {
	logData := []interface{}{
		&payload.JotLogPayload{Description: "rweb:cookiesMetadata:load", Product: "rweb", EventValue: time.Until(time.UnixMilli(fetchedTime)).Milliseconds()},
	}
	return c.performJotClientEvent(ctx, payload.JotLoggingCategoryPerftown, false, logData)
}

func (c *Client) doInitialClientLoggingEvents(ctx context.Context) error {
	err := c.doCookiesMetaDataLoad(ctx)
	if err != nil {
		return err
	}
	logData := []interface{}{
		&payload.JotLogPayload{
			Description: "rweb:init:storePrepare",
			Product:     "rweb",
			DurationMS:  9,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:perfSupported",
			Product:     "rweb",
			DurationMS:  1,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:perfSupported:" + c.session.Country,
			Product:     "rweb",
			DurationMS:  1,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:connect",
			Product:     "rweb",
			DurationMS:  165,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:connect:" + c.session.Country,
			Product:     "rweb",
			DurationMS:  165,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:process",
			Product:     "rweb",
			DurationMS:  177,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:process:" + c.session.Country,
			Product:     "rweb",
			DurationMS:  177,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:response",
			Product:     "rweb",
			DurationMS:  212,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:response:" + c.session.Country,
			Product:     "rweb",
			DurationMS:  212,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:interactivity",
			Product:     "rweb",
			DurationMS:  422,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:interactivity:" + c.session.Country,
			Product:     "rweb",
			DurationMS:  422,
		},
	}
	err = c.performJotClientEvent(ctx, payload.JotLoggingCategoryPerftown, false, logData)
	if err != nil {
		return err
	}

	triggeredTimestamp := time.Now().UnixMilli()
	logData = []interface{}{
		&payload.JotDebugLogPayload{
			Category:      payload.JotDebugLoggingCategoryClientEvent,
			TriggeredOn:   triggeredTimestamp,
			FormatVersion: 2,
			Items:         []interface{}{},
			EventNamespace: payload.EventNamespace{
				Page:   "cookie_compliance_banner",
				Action: "impression",
				Client: "m5",
			},
			ClientEventSequenceStartTimestamp: triggeredTimestamp,
			ClientEventSequenceNumber:         0,
			ClientAppID:                       "3033300",
		},
	}

	err = c.performJotClientEvent(ctx, "", true, logData)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) calculateAnimationToken() {
	if c.session.variableIndexes != nil && c.session.loadingAnims != nil && c.session.VerificationToken != "" {
		c.session.AnimationToken = crypto.GenerateAnimationState(c.session.variableIndexes, c.session.loadingAnims, c.session.VerificationToken)
	}
}
