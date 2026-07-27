package twittermeow

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/http/httpguts"
)

const (
	ClientHTTPMaxResponseBodySize    = 8 << 20
	ClientHTTPMaxResponseHeadersSize = 64 << 10
	clientHTTPMaxRequestBodySize     = 4 << 20
	clientHTTPMaxExchanges           = 64
)

var (
	ErrClientHTTPRequestPending   = errors.New("client HTTP request is pending")
	errClientHTTPOperationChanged = errors.New("client HTTP operation changed before completion")
	errClientHTTPReplayMismatch   = errors.New("client HTTP request did not match the recorded operation")
)

type ClientHTTPRequest struct {
	ID      string
	Method  string
	URL     string
	Headers http.Header
	Body    []byte
}

func (req *ClientHTTPRequest) clone() *ClientHTTPRequest {
	if req == nil {
		return nil
	}
	cloned := *req
	cloned.Headers = req.Headers.Clone()
	cloned.Body = bytes.Clone(req.Body)
	return &cloned
}

type ClientHTTPResponse struct {
	RequestID string
	Status    int
	Headers   http.Header
	Body      []byte
	FinalURL  string
}

type clientHTTPExchange struct {
	method   string
	url      string
	response ClientHTTPResponse
}

// ClientHTTPTransport turns each outgoing HTTP request into a client login step.
//
// A connector starts an operation, runs the normal twittermeow method, submits
// the client response, and reruns the same operation. Previously submitted
// responses are replayed in order until the next real request is reached. This
// keeps request construction and response parsing in Go while the network I/O
// happens in the Beeper client's network stack.
type ClientHTTPTransport struct {
	mu sync.Mutex

	operation         string
	exchanges         []clientHTTPExchange
	cursor            int
	pending           *ClientHTTPRequest
	nextID            uint64
	baseTransport     http.RoundTripper
	baseCheckRedirect func(*http.Request, []*http.Request) error
}

func newClientHTTPTransport() *ClientHTTPTransport {
	return &ClientHTTPTransport{}
}

func (t *ClientHTTPTransport) BeginOperation(operation string) error {
	operation = strings.TrimSpace(operation)
	if operation == "" {
		return fmt.Errorf("client HTTP operation name is required")
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.operation == "" {
		t.operation = operation
	} else if t.operation != operation {
		return fmt.Errorf("%w: have %q, got %q", errClientHTTPOperationChanged, t.operation, operation)
	}
	t.cursor = 0
	return nil
}

func (t *ClientHTTPTransport) EndOperation() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.pending != nil {
		return ErrClientHTTPRequestPending
	}
	if t.cursor != len(t.exchanges) {
		return fmt.Errorf("%w: consumed %d of %d responses", errClientHTTPReplayMismatch, t.cursor, len(t.exchanges))
	}
	t.resetOperationLocked()
	return nil
}

func (t *ClientHTTPTransport) ResetOperation() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.resetOperationLocked()
}

func (t *ClientHTTPTransport) resetOperationLocked() {
	t.operation = ""
	t.exchanges = nil
	t.cursor = 0
	t.pending = nil
}

func (t *ClientHTTPTransport) PendingRequest() *ClientHTTPRequest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.pending.clone()
}

func (t *ClientHTTPTransport) SubmitResponse(response ClientHTTPResponse) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.pending == nil {
		return fmt.Errorf("no client HTTP request is pending")
	}
	if response.RequestID != t.pending.ID {
		return fmt.Errorf("client HTTP response request ID does not match")
	}
	if response.Status < 100 || response.Status > 599 {
		return fmt.Errorf("invalid client HTTP response status %d", response.Status)
	}
	if len(response.Body) > ClientHTTPMaxResponseBodySize {
		return fmt.Errorf("client HTTP response body is too large")
	}
	if len(t.exchanges) >= clientHTTPMaxExchanges {
		return fmt.Errorf("client HTTP operation has too many requests")
	}
	headerSize := 0
	for name, values := range response.Headers {
		if !httpguts.ValidHeaderFieldName(name) {
			return fmt.Errorf("client HTTP response header name is invalid")
		}
		headerSize += len(name)
		for _, value := range values {
			if !httpguts.ValidHeaderFieldValue(value) {
				return fmt.Errorf("client HTTP response header value is invalid")
			}
			headerSize += len(value)
		}
	}
	if headerSize > ClientHTTPMaxResponseHeadersSize {
		return fmt.Errorf("client HTTP response headers are too large")
	}
	if response.FinalURL != "" {
		finalURL, err := url.Parse(response.FinalURL)
		if err != nil || !isAllowedClientHTTPURL(finalURL) {
			return fmt.Errorf("client HTTP response URL is not allowed")
		}
	}
	response.Headers = response.Headers.Clone()
	response.Body = bytes.Clone(response.Body)
	t.exchanges = append(t.exchanges, clientHTTPExchange{
		method:   t.pending.Method,
		url:      t.pending.URL,
		response: response,
	})
	t.pending = nil
	return nil
}

func (t *ClientHTTPTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	body, err := readClientHTTPRequestBody(request)
	if err != nil {
		return nil, err
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.operation == "" {
		return nil, fmt.Errorf("client HTTP request made outside an operation")
	}
	if t.cursor < len(t.exchanges) {
		exchange := t.exchanges[t.cursor]
		if request.Method != exchange.method || request.URL.String() != exchange.url {
			return nil, fmt.Errorf(
				"%w: response %d expected %s %s, got %s %s",
				errClientHTTPReplayMismatch,
				t.cursor,
				exchange.method,
				exchange.url,
				request.Method,
				request.URL.String(),
			)
		}
		t.cursor++
		return replayClientHTTPResponse(request, exchange.response), nil
	}
	if t.pending != nil {
		if request.Method != t.pending.Method || request.URL.String() != t.pending.URL {
			return nil, fmt.Errorf("%w: a different request arrived while waiting for the client", errClientHTTPReplayMismatch)
		}
		return nil, ErrClientHTTPRequestPending
	}
	if err = validateClientHTTPRequest(request, body); err != nil {
		return nil, err
	}

	t.nextID++
	t.pending = &ClientHTTPRequest{
		ID:      "client-http-" + strconv.FormatUint(t.nextID, 10),
		Method:  request.Method,
		URL:     request.URL.String(),
		Headers: clientHTTPFetchHeaders(request.Header),
		Body:    bytes.Clone(body),
	}
	return nil, ErrClientHTTPRequestPending
}

func readClientHTTPRequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(io.LimitReader(request.Body, clientHTTPMaxRequestBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("read client HTTP request body: %w", err)
	}
	request.Body = io.NopCloser(bytes.NewReader(body))
	if len(body) > clientHTTPMaxRequestBodySize {
		return nil, fmt.Errorf("client HTTP request body is too large")
	}
	return body, nil
}

func validateClientHTTPRequest(request *http.Request, body []byte) error {
	if request.Method != http.MethodGet && request.Method != http.MethodPost {
		return fmt.Errorf("client HTTP method %s is not allowed", request.Method)
	}
	if !isAllowedClientHTTPURL(request.URL) {
		return fmt.Errorf("client HTTP URL is not allowed")
	}
	if len(body) > clientHTTPMaxRequestBodySize {
		return fmt.Errorf("client HTTP request body is too large")
	}
	return nil
}

func isAllowedClientHTTPURL(parsed *url.URL) bool {
	if parsed == nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Port() != "" {
		return false
	}
	switch strings.ToLower(parsed.Hostname()) {
	case "x.com", "api.x.com", "twitter.com", "abs.twimg.com":
		return true
	default:
		return false
	}
}

func clientHTTPFetchHeaders(headers http.Header) http.Header {
	out := make(http.Header)
	for key, values := range headers {
		key = strings.TrimSpace(key)
		if key == "" || isClientManagedHeader(strings.ToLower(key)) {
			continue
		}
		for _, value := range values {
			if value == "" || strings.ContainsAny(value, "\r\n") {
				continue
			}
			out.Add(key, value)
		}
	}
	return out
}

func isClientManagedHeader(name string) bool {
	switch name {
	case "accept-encoding", "connection", "content-length", "expect",
		"host", "keep-alive", "te", "trailer", "transfer-encoding", "upgrade", "via":
		return true
	}
	return strings.HasPrefix(name, "proxy-")
}

func replayClientHTTPResponse(request *http.Request, response ClientHTTPResponse) *http.Response {
	statusText := http.StatusText(response.Status)
	status := strconv.Itoa(response.Status)
	if statusText != "" {
		status += " " + statusText
	}
	finalRequest := request
	if response.FinalURL != "" && response.FinalURL != request.URL.String() {
		if finalURL, err := url.Parse(response.FinalURL); err == nil {
			cloned := request.Clone(request.Context())
			cloned.URL = finalURL
			finalRequest = cloned
		}
	}
	return &http.Response{
		Status:        status,
		StatusCode:    response.Status,
		Header:        response.Headers.Clone(),
		Body:          io.NopCloser(bytes.NewReader(response.Body)),
		ContentLength: int64(len(response.Body)),
		Request:       finalRequest,
	}
}

func (c *Client) EnableClientHTTP() *ClientHTTPTransport {
	if c.clientHTTPTransport != nil {
		return c.clientHTTPTransport
	}
	transport := newClientHTTPTransport()
	transport.baseTransport = c.HTTP.Transport
	transport.baseCheckRedirect = c.HTTP.CheckRedirect
	c.HTTP.Transport = transport
	c.HTTP.CheckRedirect = nil
	c.clientHTTPTransport = transport
	return transport
}

func (c *Client) DisableClientHTTP() {
	if c.clientHTTPTransport == nil {
		return
	}
	transport := c.clientHTTPTransport
	transport.ResetOperation()
	c.HTTP.Transport = transport.baseTransport
	c.HTTP.CheckRedirect = transport.baseCheckRedirect
	c.clientHTTPTransport = nil
}

func (c *Client) IsClientHTTPEnabled() bool {
	return c.clientHTTPTransport != nil
}
