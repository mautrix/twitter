package twittermeow

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var defaultPollingInterval = 10 * time.Second

type PollingClient struct {
	client        *Client
	currentCursor string
	stop          atomic.Pointer[context.CancelFunc]
}

// interval is the delay inbetween checking for new updates
// default interval will be 10s
func (c *Client) newPollingClient() *PollingClient {
	return &PollingClient{
		client: c,
	}
}

func (pc *PollingClient) startPolling(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	if oldCancel := pc.stop.Swap(&cancel); oldCancel != nil {
		(*oldCancel)()
	}
	go pc.doPoll(ctx)
}

func (pc *PollingClient) doPoll(ctx context.Context) {
	userUpdatesQuery := (&payload.DMRequestQuery{}).Default()
	tick := time.NewTicker(defaultPollingInterval)
	defer tick.Stop()
	log := zerolog.Ctx(ctx)
	for {
		select {
		case <-tick.C:
			if pc.currentCursor != "" {
				userUpdatesQuery.Cursor = pc.currentCursor
			}

			userUpdatesResponse, err := pc.client.GetDMUserUpdates(ctx, &userUpdatesQuery)
			if err != nil {
				log.Err(err).Msg("Failed to get user updates")
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
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("Polling context canceled")
			return
		}
	}
}

func (pc *PollingClient) SetCurrentCursor(cursor string) {
	pc.currentCursor = cursor
}

func (pc *PollingClient) stopPolling() {
	if cancel := pc.stop.Swap(nil); cancel != nil {
		(*cancel)()
	}
}
