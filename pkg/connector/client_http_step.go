package connector

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

const (
	clientHTTPConfigMarkerPrefix        = "/*__BEEPER_CLIENT_HTTP_V1__:"
	clientHTTPConfigMarkerSuffix        = "__*/"
	clientHTTPDispatchURL               = "https://x.com/robots.txt#beeper-client-http"
	loginFieldClientHTTPRequestID       = "client_http_request_id"
	loginFieldClientHTTPStatus          = "client_http_status"
	loginFieldClientHTTPResponseURL     = "client_http_response_url"
	loginFieldClientHTTPResponseHeaders = "client_http_response_headers"
	loginFieldClientHTTPResponseBody    = "client_http_response_body"
	loginFieldClientHTTPError           = "client_http_error"
	clientHTTPEncodedFieldPrefix        = "v1."
)

var clientHTTPCookieNames = []string{
	"__cf_bm",
	"__cuid",
	"auth_token",
	"att",
	"ct0",
	"dtab_local",
	"gt",
	"guest_id",
	"guest_id_ads",
	"guest_id_marketing",
	"kdt",
	"lang",
	"night_mode",
	"personalization_id",
	"twid",
}

type clientHTTPConfig struct {
	RequestID    string            `json:"requestID"`
	Method       string            `json:"method"`
	RequestURL   string            `json:"requestURL"`
	Referrer     string            `json:"referrer,omitempty"`
	Headers      map[string]string `json:"headers"`
	CookieHeader string            `json:"cookieHeader,omitempty"`
	Body         string            `json:"body"`
}

func makeClientHTTPRequestStep(request *twittermeow.ClientHTTPRequest) (*bridgev2.LoginStep, error) {
	if request == nil {
		return nil, fmt.Errorf("client HTTP request is missing")
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       LoginStepIDClientHTTPRequest,
		Instructions: "Sending an X sign-in request from this device. Client HTTP is in beta.",
		CookiesParams: &bridgev2.LoginCookiesParams{
			URL:               clientHTTPDispatchURL,
			ExtractJS:         clientHTTPRequestEnvelope(request),
			WaitForURLPattern: `^https://x\.com/robots\.txt#beeper-client-http$`,
			Fields:            clientHTTPCookieFields(),
			Hidden:            true,
		},
	}, nil
}

func clientHTTPRequestEnvelope(request *twittermeow.ClientHTTPRequest) string {
	config, err := json.Marshal(clientHTTPConfig{
		RequestID:    request.ID,
		Method:       request.Method,
		RequestURL:   request.URL,
		Referrer:     request.Referrer,
		Headers:      request.Headers,
		CookieHeader: request.CookieHeader,
		Body:         base64.RawURLEncoding.EncodeToString(request.Body),
	})
	if err != nil {
		panic(fmt.Errorf("marshal client HTTP request config: %w", err))
	}

	encodedConfig := base64.RawURLEncoding.EncodeToString(config)
	requestID, _ := json.Marshal(request.ID)
	return clientHTTPConfigMarkerPrefix + encodedConfig + clientHTTPConfigMarkerSuffix + "\n" +
		"(async () => ({client_http_request_id:" + string(requestID) +
		",client_http_error:\"Client HTTP requires a supported Beeper client\"}))()"
}

func clientHTTPCookieFields() []bridgev2.LoginCookieField {
	fields := []bridgev2.LoginCookieField{
		clientHTTPSpecialField(loginFieldClientHTTPRequestID, true, `^client-http-[0-9]+$`),
		clientHTTPSpecialField(loginFieldClientHTTPStatus, false, `^[1-5][0-9]{2}$`),
		clientHTTPSpecialField(loginFieldClientHTTPResponseURL, false, `^https://[^\r\n]{1,4096}$`),
		clientHTTPSpecialField(loginFieldClientHTTPResponseHeaders, false, `^v1\.[A-Za-z0-9_-]*$`),
		clientHTTPSpecialField(loginFieldClientHTTPResponseBody, false, `^v1\.[A-Za-z0-9_-]*$`),
		clientHTTPSpecialField(loginFieldClientHTTPError, false, `^[^\r\n]{1,300}$`),
	}
	for _, field := range browserHeaderFields {
		fields = append(fields, bridgev2.LoginCookieField{
			ID:       field.ID,
			Required: field.Required,
			Pattern:  field.Pattern,
			Sources: []bridgev2.LoginCookieFieldSource{{
				Type: bridgev2.LoginCookieTypeSpecial,
				Name: field.ID,
			}},
		})
	}
	for _, name := range clientHTTPCookieNames {
		fields = append(fields, bridgev2.LoginCookieField{
			ID:       name,
			Required: false,
			Sources: []bridgev2.LoginCookieFieldSource{
				{Type: bridgev2.LoginCookieTypeCookie, Name: name, CookieDomain: "x.com"},
				{Type: bridgev2.LoginCookieTypeCookie, Name: name, CookieDomain: ".x.com"},
			},
		})
	}
	return fields
}

func clientHTTPSpecialField(id string, required bool, pattern string) bridgev2.LoginCookieField {
	return bridgev2.LoginCookieField{
		ID:       id,
		Required: required,
		Pattern:  pattern,
		Sources: []bridgev2.LoginCookieFieldSource{{
			Type: bridgev2.LoginCookieTypeSpecial,
			Name: id,
		}},
	}
}

func decodeClientHTTPResponse(input map[string]string) (twittermeow.ClientHTTPResponse, error) {
	requestID := strings.TrimSpace(input[loginFieldClientHTTPRequestID])
	if requestID == "" {
		return twittermeow.ClientHTTPResponse{}, fmt.Errorf("client HTTP response request ID is missing")
	}
	status, err := strconv.Atoi(strings.TrimSpace(input[loginFieldClientHTTPStatus]))
	if err != nil {
		return twittermeow.ClientHTTPResponse{}, fmt.Errorf("client HTTP response status is invalid")
	}
	body, err := decodeClientHTTPEncodedField(
		input[loginFieldClientHTTPResponseBody],
		twittermeow.ClientHTTPMaxResponseBodySize,
	)
	if err != nil {
		return twittermeow.ClientHTTPResponse{}, fmt.Errorf("client HTTP response body is invalid: %w", err)
	}
	headersJSON, err := decodeClientHTTPEncodedField(
		input[loginFieldClientHTTPResponseHeaders],
		twittermeow.ClientHTTPMaxResponseHeadersSize,
	)
	if err != nil {
		return twittermeow.ClientHTTPResponse{}, fmt.Errorf("client HTTP response headers are invalid: %w", err)
	}
	var headerValues map[string]string
	if err = json.Unmarshal(headersJSON, &headerValues); err != nil {
		return twittermeow.ClientHTTPResponse{}, fmt.Errorf("client HTTP response headers are invalid: %w", err)
	}
	headers := make(http.Header, len(headerValues))
	for name, value := range headerValues {
		if !httpguts.ValidHeaderFieldName(name) || !httpguts.ValidHeaderFieldValue(value) {
			return twittermeow.ClientHTTPResponse{}, fmt.Errorf("client HTTP response header is invalid")
		}
		headers.Set(name, value)
	}
	return twittermeow.ClientHTTPResponse{
		RequestID: requestID,
		Status:    status,
		Headers:   headers,
		Body:      body,
		FinalURL:  strings.TrimSpace(input[loginFieldClientHTTPResponseURL]),
	}, nil
}

func decodeClientHTTPEncodedField(value string, maxDecodedSize int) ([]byte, error) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, clientHTTPEncodedFieldPrefix) {
		return nil, fmt.Errorf("unsupported encoding")
	}
	encoded := strings.TrimPrefix(value, clientHTTPEncodedFieldPrefix)
	if base64.RawURLEncoding.DecodedLen(len(encoded)) > maxDecodedSize {
		return nil, fmt.Errorf("encoded value is too large")
	}
	return base64.RawURLEncoding.DecodeString(encoded)
}

func clientHTTPCookies(input map[string]string) map[string]string {
	out := make(map[string]string)
	for _, name := range clientHTTPCookieNames {
		if value := strings.TrimSpace(input[name]); value != "" {
			out[name] = value
		}
	}
	return out
}
