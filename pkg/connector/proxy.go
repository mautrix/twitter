package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"maunium.net/go/mautrix"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

var getProxyClient = &http.Client{Timeout: 30 * time.Second}

func (tc *TwitterConnector) getProxy(ctx context.Context, reason string) (string, error) {
	if tc.Config.GetProxyURL == "" {
		return tc.Config.Proxy, nil
	}
	parsed, err := url.Parse(tc.Config.GetProxyURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse get_proxy_url: %w", err)
	}
	q := parsed.Query()
	q.Set("reason", reason)
	parsed.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header.Set("User-Agent", mautrix.DefaultUserAgent)
	resp, err := getProxyClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var respData struct {
		ProxyURL string `json:"proxy_url"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return respData.ProxyURL, nil
}

// applyProxy sets the configured proxy on a client. reason is passed to get_proxy_url.
func (tc *TwitterConnector) applyProxy(ctx context.Context, client *twittermeow.Client, reason string) error {
	if tc == nil || (tc.Config.Proxy == "" && tc.Config.GetProxyURL == "") {
		return nil
	}
	proxy, err := tc.getProxy(ctx, reason)
	if err != nil {
		return fmt.Errorf("failed to get proxy address: %w", err)
	}
	return client.SetProxy(proxy)
}
