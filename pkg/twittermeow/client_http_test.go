package twittermeow

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
)

func TestClientHTTPTransportCapturesAndReplaysRequest(t *testing.T) {
	transport := newClientHTTPTransport()
	if err := transport.BeginOperation("credentials"); err != nil {
		t.Fatal(err)
	}

	body := "$castle_token=" + url.QueryEscape("castle-from-webview") + "&password=secret"
	request, err := http.NewRequest(http.MethodPost, "https://x.com/i/jfapi/onboarding/web/actions/begin_login", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("authorization", "Bearer public")
	request.Header.Set("cookie", "auth_token=secret")
	request.Header.Set("origin", "https://x.com")
	request.Header.Set("referer", "https://x.com/i/jf/onboarding/web?mode=login")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("user-agent", "bridge user agent")
	request.Header.Add("x-multi", "one")
	request.Header.Add("x-multi", "two")

	_, err = transport.RoundTrip(request)
	if !errors.Is(err, ErrClientHTTPRequestPending) {
		t.Fatalf("RoundTrip() error = %v, want pending", err)
	}
	pending := transport.PendingRequest()
	if pending == nil || pending.Method != http.MethodPost || pending.URL != request.URL.String() {
		t.Fatalf("PendingRequest() = %#v", pending)
	}
	if pending.Headers.Get("cookie") != "auth_token=secret" {
		t.Fatalf("Cookie header = %q", pending.Headers.Get("cookie"))
	}
	if pending.Headers.Get("authorization") != "Bearer public" {
		t.Fatalf("authorization header missing: %#v", pending.Headers)
	}
	for _, forwarded := range []string{"cookie", "origin", "referer", "sec-fetch-site", "user-agent"} {
		if pending.Headers.Get(forwarded) == "" {
			t.Fatalf("end-to-end header %q was not forwarded", forwarded)
		}
	}
	if got := pending.Headers.Values("x-multi"); len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("multi-value header was not preserved: %#v", pending.Headers)
	}

	err = transport.SubmitResponse(ClientHTTPResponse{
		RequestID: pending.ID,
		Status:    http.StatusOK,
		Headers: http.Header{
			"Content-Type": {"application/octet-stream"},
			"Set-Cookie":   {"first=1; Path=/", "second=2; Path=/"},
		},
		Body:     []byte("response body"),
		FinalURL: "https://x.com/home",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = transport.BeginOperation("credentials"); err != nil {
		t.Fatal(err)
	}
	replayRequest := request.Clone(context.Background())
	replayRequest.Body = io.NopCloser(strings.NewReader(body))
	response, err := transport.RoundTrip(replayRequest)
	if err != nil {
		t.Fatal(err)
	}
	replayedBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || string(replayedBody) != "response body" {
		t.Fatalf("replayed response = %d %q", response.StatusCode, replayedBody)
	}
	if response.Request.URL.String() != "https://x.com/home" {
		t.Fatalf("replayed final URL = %q", response.Request.URL)
	}
	if got := response.Header.Values("Set-Cookie"); len(got) != 2 {
		t.Fatalf("replayed Set-Cookie headers = %#v", got)
	}
	if err = transport.EndOperation(); err != nil {
		t.Fatal(err)
	}
}

func TestClientHTTPTransportRejectsUnsafeRequestsAndResponses(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
	}{
		{name: "plain HTTP", method: http.MethodGet, url: "http://x.com/login"},
		{name: "credentials in URL", method: http.MethodGet, url: "https://user:pass@x.com/login"},
		{name: "unknown host", method: http.MethodGet, url: "https://example.com/login"},
		{name: "explicit port", method: http.MethodGet, url: "https://x.com:443/login"},
		{name: "unsupported method", method: http.MethodDelete, url: "https://x.com/login"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := newClientHTTPTransport()
			if err := transport.BeginOperation("test"); err != nil {
				t.Fatal(err)
			}
			request, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = transport.RoundTrip(request); err == nil || errors.Is(err, ErrClientHTTPRequestPending) {
				t.Fatalf("RoundTrip() error = %v, want safety rejection", err)
			}
		})
	}

	transport := newClientHTTPTransport()
	if err := transport.BeginOperation("test"); err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://x.com/login", nil)
	_, _ = transport.RoundTrip(request)
	pending := transport.PendingRequest()
	for _, response := range []ClientHTTPResponse{
		{RequestID: "wrong", Status: http.StatusOK},
		{RequestID: pending.ID, Status: 99},
		{RequestID: pending.ID, Status: http.StatusOK, FinalURL: "https://example.com/"},
		{RequestID: pending.ID, Status: http.StatusOK, Body: make([]byte, ClientHTTPMaxResponseBodySize+1)},
		{RequestID: pending.ID, Status: http.StatusOK, Headers: http.Header{"Bad Header": {"value"}}},
		{RequestID: pending.ID, Status: http.StatusOK, Headers: http.Header{"X-Test": {"bad\rvalue"}}},
		{RequestID: pending.ID, Status: http.StatusOK, Headers: http.Header{
			"X-Test": {strings.Repeat("a", ClientHTTPMaxResponseHeadersSize+1)},
		}},
	} {
		if err := transport.SubmitResponse(response); err == nil {
			t.Fatalf("SubmitResponse(%#v) succeeded, want rejection", response)
		}
	}

	largeBodyTransport := newClientHTTPTransport()
	if err := largeBodyTransport.BeginOperation("test"); err != nil {
		t.Fatal(err)
	}
	largeRequest, _ := http.NewRequest(
		http.MethodPost,
		"https://x.com/login",
		strings.NewReader(strings.Repeat("a", clientHTTPMaxRequestBodySize+1)),
	)
	if _, err := largeBodyTransport.RoundTrip(largeRequest); err == nil {
		t.Fatal("oversized client HTTP request body was accepted")
	}
}

func TestClientHTTPTransportOperationAndReplayGuards(t *testing.T) {
	transport := newClientHTTPTransport()
	if err := transport.BeginOperation("first"); err != nil {
		t.Fatal(err)
	}
	if err := transport.BeginOperation("second"); !errors.Is(err, errClientHTTPOperationChanged) {
		t.Fatalf("BeginOperation() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://x.com/login", nil)
	_, _ = transport.RoundTrip(request)
	if err := transport.EndOperation(); !errors.Is(err, ErrClientHTTPRequestPending) {
		t.Fatalf("EndOperation() error = %v", err)
	}
	pending := transport.PendingRequest()
	if err := transport.SubmitResponse(ClientHTTPResponse{RequestID: pending.ID, Status: http.StatusOK}); err != nil {
		t.Fatal(err)
	}
	if err := transport.SubmitResponse(ClientHTTPResponse{RequestID: pending.ID, Status: http.StatusOK}); err == nil {
		t.Fatal("duplicate client HTTP response was accepted")
	}
	if err := transport.BeginOperation("first"); err != nil {
		t.Fatal(err)
	}
	different, _ := http.NewRequest(http.MethodGet, "https://x.com/messages", nil)
	if _, err := transport.RoundTrip(different); !errors.Is(err, errClientHTTPReplayMismatch) {
		t.Fatalf("RoundTrip() error = %v", err)
	}
}

func TestClientHTTPModeUsesNativeCastleTokensAndRestoresTransport(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	originalTransport := client.HTTP.Transport
	transport := client.EnableClientHTTP()
	if transport == nil || !client.IsClientHTTPEnabled() || client.HTTP.Transport != transport {
		t.Fatal("client HTTP mode was not enabled")
	}
	form := url.Values{}
	if err := client.addJetfuelCastleTokenToForm(form); !errors.Is(err, ErrJetfuelCastleTokenRequired) {
		t.Fatalf("addJetfuelCastleTokenToForm() error = %v, want Castle token request", err)
	}
	client.SetNextJetfuelCastleTokens([]string{"castle-from-webview"})
	if err := client.addJetfuelCastleTokenToForm(form); err != nil {
		t.Fatal(err)
	}
	if got := form.Get("$castle_token"); got != "castle-from-webview" {
		t.Fatalf("$castle_token = %q", got)
	}
	client.DisableClientHTTP()
	if client.IsClientHTTPEnabled() || client.HTTP.Transport != originalTransport {
		t.Fatal("client HTTP mode did not restore the original transport")
	}
}

func TestMakeRequestReturnsClientHTTPPendingWithoutRetry(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	transport := client.EnableClientHTTP()
	if err := transport.BeginOperation("test"); err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	_, _, err := client.MakeRequest(
		context.Background(),
		"https://x.com/login",
		http.MethodGet,
		http.Header{},
		nil,
		"",
	)
	if !errors.Is(err, ErrClientHTTPRequestPending) {
		t.Fatalf("MakeRequest() error = %v", err)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("MakeRequest() took %s, likely retried pending request", elapsed)
	}
}
