package twittermeow

import (
	"testing"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
)

func TestClientUsesCapturedBrowserHeaders(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	captured := BrowserHeaders{
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/150.0.0.0 Safari/537.36",
		SecCHUserAgent: `"Chromium";v="150", "Google Chrome";v="150", "Not_A Brand";v="99"`,
		SecCHPlatform:  `"Windows"`,
		SecCHMobile:    "?0",
	}
	if !client.SetBrowserHeaders(captured) {
		t.Fatal("SetBrowserHeaders() rejected valid captured headers")
	}

	headers := client.GetBaseHeaders()
	if got := headers.Get("User-Agent"); got != captured.UserAgent {
		t.Fatalf("User-Agent = %q, want %q", got, captured.UserAgent)
	}
	if got := headers.Get("Sec-Ch-Ua"); got != captured.SecCHUserAgent {
		t.Fatalf("Sec-Ch-Ua = %q, want %q", got, captured.SecCHUserAgent)
	}
	if got := headers.Get("Sec-Ch-Ua-Platform"); got != captured.SecCHPlatform {
		t.Fatalf("Sec-Ch-Ua-Platform = %q, want %q", got, captured.SecCHPlatform)
	}
	if got := headers.Get("Sec-Ch-Ua-Mobile"); got != captured.SecCHMobile {
		t.Fatalf("Sec-Ch-Ua-Mobile = %q, want %q", got, captured.SecCHMobile)
	}
	if got := BaseHeaders.Get("User-Agent"); got != UserAgent {
		t.Fatalf("captured fingerprint mutated fallback headers: %q", got)
	}
}

func TestClientOmitsUnavailableClientHints(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	userAgent := "Mozilla/5.0 (Android 15; Mobile; rv:141.0) Gecko/141.0 Firefox/141.0"
	if !client.SetBrowserHeaders(BrowserHeaders{UserAgent: userAgent}) {
		t.Fatal("SetBrowserHeaders() rejected Firefox user agent")
	}

	headers := client.GetBaseHeaders()
	if got := headers.Get("User-Agent"); got != userAgent {
		t.Fatalf("User-Agent = %q, want %q", got, userAgent)
	}
	for _, name := range []string{"Sec-Ch-Ua", "Sec-Ch-Ua-Platform", "Sec-Ch-Ua-Mobile"} {
		if got := headers.Get(name); got != "" {
			t.Fatalf("%s = %q, want omitted", name, got)
		}
	}
}

func TestClientRejectsUnsafeCapturedBrowserHeaders(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	if client.SetBrowserHeaders(BrowserHeaders{UserAgent: "browser\r\nX-Injected: true"}) {
		t.Fatal("SetBrowserHeaders() accepted a header injection")
	}
	if got := client.GetUserAgent(); got != UserAgent {
		t.Fatalf("GetUserAgent() = %q, want fallback %q", got, UserAgent)
	}
}
