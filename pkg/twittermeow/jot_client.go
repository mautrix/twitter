package twittermeow

import (
	"context"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

type JotClient struct {
	client *Client
}

func (c *Client) newJotClient() *JotClient {
	return &JotClient{
		client: c,
	}
}

// sendClientLoggingEvent is disabled - JOT logging not currently used.
func (jc *JotClient) sendClientLoggingEvent(ctx context.Context, category payload.JotLoggingCategory, debug bool, body []interface{}) error {
	return nil
}
