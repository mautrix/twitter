package methods

import (
	"io"
	"net/http"
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
