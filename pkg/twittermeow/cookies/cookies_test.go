package cookies

import (
	"net/http"
	"testing"
)

func TestNewCookiesFromStringParsesCookieHeader(t *testing.T) {
	store := NewCookiesFromString("auth_token=auth-value; ct0=csrf-value")

	if got := store.Get(XAuthToken); got != "auth-value" {
		t.Fatalf("auth_token = %q, want auth-value", got)
	}
	if got := store.Get(XCt0); got != "csrf-value" {
		t.Fatalf("ct0 = %q, want csrf-value", got)
	}
}

func TestNewCookiesFromStringParsesSingleSetCookieHeader(t *testing.T) {
	store := NewCookiesFromString("ct0=csrf-value; Path=/; Secure; HttpOnly")

	if got := store.Get(XCt0); got != "csrf-value" {
		t.Fatalf("ct0 = %q, want csrf-value", got)
	}
	if got := store.Get(XCookieName("Path")); got != "" {
		t.Fatalf("Path pseudo-cookie = %q, want omitted", got)
	}
}

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
