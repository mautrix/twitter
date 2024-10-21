package twittermeow

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

var (
	errCookieGuestIDNotFound = errors.New("failed to retrieve and set 'guest_id' cookie from Set-Cookie response headers")
)

// retrieved from main page resp, its a 2 year old timestamp; looks constant
const fetchedTime = 1661971138705

type SessionAuthTokens struct {
	authenticated    string
	notAuthenticated string
}

type SessionLoader struct {
	client            *Client
	currentUser       *response.AccountSettingsResponse
	verificationToken string
	country           string
	clientUUID        string
	authTokens        *SessionAuthTokens
}

func (c *Client) newSessionLoader() *SessionLoader {
	return &SessionLoader{
		client:     c,
		clientUUID: uuid.NewString(),
		authTokens: &SessionAuthTokens{},
	}
}

func (s *SessionLoader) SetCurrentUser(data *response.AccountSettingsResponse) {
	s.currentUser = data
}

func (s *SessionLoader) GetCurrentUser() *response.AccountSettingsResponse {
	return s.currentUser
}

func (s *SessionLoader) isAuthenticated() bool {
	return s.currentUser != nil && s.currentUser.ScreenName != ""
}

func (s *SessionLoader) LoadPage(url string) error {
	mainPageUrl, err := neturl.Parse(url)
	if err != nil {
		return err
	}
	extraHeaders := map[string]string{
		"upgrade-insecure-requests": "1",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"sec-fetch-dest":            "document",
	}
	mainPageResp, mainPageRespBody, err := s.client.MakeRequest(url, http.MethodGet, s.client.buildHeaders(HeaderOpts{Extra: extraHeaders, WithCookies: true}), nil, types.NONE)
	if err != nil {
		return err
	}

	s.client.cookies.UpdateFromResponse(mainPageResp)
	if s.client.cookies.IsCookieEmpty(cookies.XGuestID) {
		s.client.Logger.Err(errCookieGuestIDNotFound).Msg("No GuestID found in response headers")
		return errCookieGuestIDNotFound
	}

	mainPageHTML := string(mainPageRespBody)
	migrationUrl, migrationRequired := methods.ParseMigrateURL(mainPageHTML)
	if migrationRequired {
		log.Println("migration is required...")
		extraHeaders = map[string]string{
			"upgrade-insecure-requests": "1",
			"sec-fetch-site":            "cross-site",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-dest":            "document",
		}
		migrationPageResp, migrationPageRespBody, err := s.client.MakeRequest(migrationUrl, http.MethodGet, s.client.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: fmt.Sprintf("https://%s/", mainPageUrl.Host), WithCookies: true}), nil, types.NONE)
		if err != nil {
			return err
		}

		migrationPageHTML := string(migrationPageRespBody)
		migrationFormUrl, migrationFormPayload := methods.ParseMigrateRequestData(migrationPageHTML)
		if migrationFormPayload != nil {
			migrationForm, err := query.Values(migrationFormPayload)
			if err != nil {
				return err
			}
			migrationPayload := []byte(migrationForm.Encode())
			extraHeaders["origin"] = endpoints.TWITTER_BASE_URL

			s.client.disableRedirects()
			mainPageResp, _, err = s.client.MakeRequest(migrationFormUrl, http.MethodPost, s.client.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: endpoints.TWITTER_BASE_URL + "/", WithCookies: true}), migrationPayload, types.FORM)
			if err == nil && !errors.Is(err, ErrRedirectAttempted) {
				return fmt.Errorf("failed to make request to main page, server did not respond with a redirect response")
			}
			s.client.enableRedirects()
			s.client.cookies.UpdateFromResponse(mainPageResp) // update the cookies received from the redirected response headers

			migrationFormUrl = endpoints.BASE_URL + mainPageResp.Header.Get("Location")
			mainPageResp, mainPageRespBody, err = s.client.MakeRequest(migrationFormUrl, http.MethodGet, s.client.buildHeaders(HeaderOpts{Extra: extraHeaders, Referer: endpoints.TWITTER_BASE_URL + "/", WithCookies: true}), migrationPayload, types.FORM)
			if err != nil {
				return err
			}

			mainPageHTML := string(mainPageRespBody)
			err = s.client.parseMainPageHTML(mainPageResp, mainPageHTML)
			if err != nil {
				return err
			}

			err = s.doInitialClientLoggingEvents()
			if err != nil {
				return err
			}

		} else {
			return fmt.Errorf("failed to find form request data in migration response: (response_body=%s, status_code=%d)", migrationPageHTML, migrationPageResp.StatusCode)
		}
	} else {
		// most likely means... already authenticated
		mainPageHTML := string(mainPageRespBody)
		err = s.client.parseMainPageHTML(mainPageResp, mainPageHTML)
		if err != nil {
			return err
		}

		err = s.doInitialClientLoggingEvents()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionLoader) doCookiesMetaDataLoad() error {
	logData := []interface{}{
		&payload.JotLogPayload{Description: "rweb:cookiesMetadata:load", Product: "rweb", EventValue: time.Until(time.UnixMilli(fetchedTime)).Milliseconds()},
	}
	return s.client.performJotClientEvent(payload.JotLoggingCategoryPerftown, false, logData)
}

func (s *SessionLoader) doInitialClientLoggingEvents() error {
	err := s.doCookiesMetaDataLoad()
	if err != nil {
		return err
	}
	country := s.GetCountry()
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
			Description: "rweb:ttft:perfSupported:" + country,
			Product:     "rweb",
			DurationMS:  1,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:connect",
			Product:     "rweb",
			DurationMS:  165,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:connect:" + country,
			Product:     "rweb",
			DurationMS:  165,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:process",
			Product:     "rweb",
			DurationMS:  177,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:process:" + country,
			Product:     "rweb",
			DurationMS:  177,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:response",
			Product:     "rweb",
			DurationMS:  212,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:response:" + country,
			Product:     "rweb",
			DurationMS:  212,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:interactivity",
			Product:     "rweb",
			DurationMS:  422,
		},
		&payload.JotLogPayload{
			Description: "rweb:ttft:interactivity:" + country,
			Product:     "rweb",
			DurationMS:  422,
		},
	}
	err = s.client.performJotClientEvent(payload.JotLoggingCategoryPerftown, false, logData)
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

	err = s.client.performJotClientEvent("", true, logData)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionLoader) GetCountry() string {
	return s.country
}

func (s *SessionLoader) SetCountry(country string) {
	s.country = country
}

func (s *SessionLoader) SetVerificationToken(verificationToken string) {
	s.verificationToken = verificationToken
}

func (s *SessionLoader) GetVerificationToken() string {
	return s.verificationToken
}

func (s *SessionLoader) SetAuthTokens(authenticatedToken, notAuthenticatedToken string) {
	s.authTokens.authenticated = authenticatedToken
	s.authTokens.notAuthenticated = notAuthenticatedToken
}

func (s *SessionLoader) GetAuthTokens() *SessionAuthTokens {
	return s.authTokens
}
