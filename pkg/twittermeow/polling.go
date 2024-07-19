package twittermeow

import (
	"log"
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
	userUpdatesQuery := (&payload.DmRequestQuery{}).Default()
	for range pc.ticker.C {
		if pc.currentCursor != "" {
			userUpdatesQuery.Cursor = pc.currentCursor
		}

		userUpdatesResponse, err := pc.client.GetDMUserUpdates(userUpdatesQuery)
		if err != nil {
			log.Fatal(err)
		}

		userEvents := userUpdatesResponse.UserEvents
		if len(userEvents.Entries) > 0 {
			pc.client.processEventEntries(userUpdatesResponse)
		}

		// pc.client.logger.Info().Any("user_events", userUpdatesResponse.UserEvents).Any("inbox_initial_state", userUpdatesResponse.InboxInitialState).Msg("Got polling update response")

		pc.SetCurrentCursor(userUpdatesResponse.UserEvents.Cursor)
	}
}

func (pc *PollingClient) SetCurrentCursor(cursor string) {
	pc.currentCursor = cursor
}

//lint:ignore U1000 TODO fix unused method
func (pc *PollingClient) stopPolling() error {
	if pc.ticker == nil {
		return ErrNotPollingUpdates
	}

	pc.ticker.Stop()
	pc.ticker = nil

	return nil
}
