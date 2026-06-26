package cookies

import (
	"net/http"
	"testing"
)

func TestUpdateFromResponseKeepsSessionCookies(t *testing.T) {
	store := NewCookies(map[string]string{"deleted": "old"})
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Add("Set-Cookie", "__cf_bm=session-value; Path=/; Secure; HttpOnly")
	resp.Header.Add("Set-Cookie", "deleted=gone; Max-Age=0; Path=/")

	store.UpdateFromResponse(resp)

	if got := store.Get(XCookieName("__cf_bm")); got != "session-value" {
		t.Fatalf("__cf_bm = %q, want session-value", got)
	}
	if got := store.Get(XCookieName("deleted")); got != "" {
		t.Fatalf("deleted cookie still present: %q", got)
	}
}
