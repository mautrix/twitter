package methods

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestParseDocumentCookieAssignments(t *testing.T) {
	html := `<script>
document.cookie = "ct0=csrf-value; Path=/; Secure";
document.cookie = "guest_id=guest-value; Domain=.x.com; SameSite=None";
</script>`

	cookies := ParseDocumentCookieAssignments(html)
	if got := cookies["ct0"]; got != "csrf-value" {
		t.Fatalf("ct0 = %q, want csrf-value", got)
	}
	if got := cookies["guest_id"]; got != "guest-value" {
		t.Fatalf("guest_id = %q, want guest-value", got)
	}
	if got := cookies["Path"]; got != "" {
		t.Fatalf("Path pseudo-cookie = %q, want omitted", got)
	}
}

func TestParseCastleBootstrapInfo(t *testing.T) {
	html := `{"responsive_web_castle_public_key":{"value":"castle-public-key"},"lang":"en"}`
	if got := ParseResponsiveWebCastlePublicKey(html); got != "castle-public-key" {
		t.Fatalf("ParseResponsiveWebCastlePublicKey() = %q", got)
	}

	js := []byte(`{100:"bundle.Home",15793:"ondemand.castle",16000:"other"};{100:"abc1234",15793:"1ff15ff",16000:"def5678"}`)
	gotURL := ParseOndemandCastleURLFromScript(js)
	wantURL := "https://abs.twimg.com/responsive-web/client-web/ondemand.castle.1ff15ffa.js"
	if gotURL != wantURL {
		t.Fatalf("ParseOndemandCastleURLFromScript() = %q, want %q", gotURL, wantURL)
	}
	if strings.Contains(gotURL, "ondemand.s") {
		t.Fatalf("ParseOndemandCastleURLFromScript() = %q, want Castle chunk URL", gotURL)
	}
}

func TestParseOndemandSURLFromScript(t *testing.T) {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://x.com/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("failed to fetch x.com: %v", err)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read x.com response: %v", err)
	}

	ondemandURL := ParseOndemandSURLFromScript(html)
	if ondemandURL == "" {
		mainScriptURL := ParseMainScriptURL(string(html))
		if mainScriptURL == "" {
			t.Fatalf("failed to locate main script URL from x.com response")
		}

		req, err = http.NewRequest(http.MethodGet, mainScriptURL, nil)
		if err != nil {
			t.Fatalf("failed to create main script request: %v", err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("failed to fetch main script: %v", err)
		}
		defer resp.Body.Close()

		script, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read main script response: %v", err)
		}
		ondemandURL = ParseOndemandSURLFromScript(script)
	}

	if ondemandURL == "" {
		t.Fatalf("failed to resolve ondemand.s URL from live x.com bootstrap")
	}
}
