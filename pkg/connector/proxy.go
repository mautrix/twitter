package connector

import (
	"context"
	"fmt"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

func (tc *TwitterConnector) SupportsProxy() {}

func (tc *TwitterConnector) initProxy(client *twittermeow.Client) error {
	if tc == nil || !tc.br.Config.Proxy.IsConfigured() {
		return nil
	}
	client.RequireProxy()
	if tc.br.Config.Proxy.LookupURL == "" {
		return client.SetProxy(tc.br.Config.Proxy.URL)
	}
	return nil
}

func (tc *TwitterConnector) applyProxy(ctx context.Context, client *twittermeow.Client, reason string) error {
	if tc == nil || !tc.br.Config.Proxy.IsConfigured() {
		return nil
	}
	proxy, err := tc.br.GetProxy(ctx, reason)
	if err != nil {
		return fmt.Errorf("failed to get proxy address: %w", err)
	}
	return client.SetProxy(proxy)
}
