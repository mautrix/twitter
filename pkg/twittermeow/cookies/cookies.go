package cookies

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type XCookieName string

const (
	XAuthToken         XCookieName = "auth_token"
	XGuestID           XCookieName = "guest_id"
	XNightMode         XCookieName = "night_mode"
	XGuestToken        XCookieName = "gt"
	XCt0               XCookieName = "ct0"
	XKdt               XCookieName = "kdt"
	XTwid              XCookieName = "twid"
	XLang              XCookieName = "lang"
	XAtt               XCookieName = "att"
	XPersonalizationID XCookieName = "personalization_id"
	XGuestIDMarketing  XCookieName = "guest_id_marketing"
)

type Cookies struct {
	store map[string]string
	lock  sync.RWMutex
}

func NewCookies(store map[string]string) *Cookies {
	if store == nil {
		store = make(map[string]string)
	}
	return &Cookies{
		store: store,
		lock:  sync.RWMutex{},
	}
}

func NewCookiesFromString(cookieStr string) *Cookies {
	c := NewCookies(nil)
	cookieStrings := strings.Split(cookieStr, ";")
	fakeHeader := http.Header{}
	for _, cookieStr := range cookieStrings {
		trimmedCookieStr := strings.TrimSpace(cookieStr)
		if trimmedCookieStr != "" {
			fakeHeader.Add("Set-Cookie", trimmedCookieStr)
		}
	}
	fakeResponse := &http.Response{Header: fakeHeader}

	for _, cookie := range fakeResponse.Cookies() {
		c.store[cookie.Name] = cookie.Value
	}

	return c
}

func (c *Cookies) String() string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var out []string
	for k, v := range c.store {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(out, "; ")
}

func (c *Cookies) IsCookieEmpty(key XCookieName) bool {
	return c.Get(key) == ""
}

func (c *Cookies) Get(key XCookieName) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[string(key)]
}

func (c *Cookies) Set(key XCookieName, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.store[string(key)] = value
}

func (c *Cookies) UpdateFromResponse(r *http.Response) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, cookie := range r.Cookies() {
		if cookie.MaxAge == 0 || cookie.Expires.Before(time.Now()) {
			delete(c.store, cookie.Name)
		} else {
			//log.Println(fmt.Sprintf("updated cookie %s to value %s", cookie.Name, cookie.Value))
			c.store[cookie.Name] = cookie.Value
		}
	}
}
