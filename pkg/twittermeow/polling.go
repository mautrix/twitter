package twittermeow

import (
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var defaultPollingInterval = 10 * time.Second

type PollingClient struct {
	client        *Client
	interval      *time.Duration
	ticker        *time.Ticker
	currentCursor string
}

// interval is the delay inbetween checking for new updates
// default interval will be 10s
func (c *Client) newPollingClient(interval *time.Duration) *PollingClient {
	if interval == nil {
		interval = &defaultPollingInterval
	}
	return &PollingClient{
		client:   c,
		interval: interval,
	}
}

func (pc *PollingClient) startPolling() error {
	if pc.ticker != nil {
		return ErrAlreadyPollingUpdates
	}

	pc.ticker = time.NewTicker(*pc.interval)
	go pc.startListening()

	return nil
}

func (pc *PollingClient) startListening() {
	userUpdatesQuery := (&payload.DMRequestQuery{}).Default()
	for range pc.ticker.C {
		if pc.currentCursor != "" {
			userUpdatesQuery.Cursor = pc.currentCursor
		}

		userUpdatesResponse, err := pc.client.GetDMUserUpdates(&userUpdatesQuery)
		if err != nil {
			pc.client.Logger.Err(err).Msg("Failed to get user updates")
			time.Sleep(1 * time.Minute)
			continue
		}

		pc.client.eventHandler(nil, userUpdatesResponse.UserEvents)
		for _, entry := range userUpdatesResponse.UserEvents.Entries {
			parsed := entry.ParseWithErrorLog(&pc.client.Logger)
			if parsed != nil {
				pc.client.eventHandler(parsed, userUpdatesResponse.UserEvents)
			}
		}

		pc.SetCurrentCursor(userUpdatesResponse.UserEvents.Cursor)
	}
}

func (pc *PollingClient) SetCurrentCursor(cursor string) {
	pc.currentCursor = cursor
}

func (pc *PollingClient) stopPolling() error {
	if pc.ticker == nil {
		return ErrNotPollingUpdates
	}

	pc.ticker.Stop()
	pc.ticker = nil

	return nil
}
