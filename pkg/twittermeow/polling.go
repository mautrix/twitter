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
	client                *Client
	currentCursor         string
	stop                  atomic.Pointer[context.CancelFunc]
	activeConversationID  string
	includeConversationID bool
	shortCircuit          chan struct{}
}

// interval is the delay inbetween checking for new updates
// default interval will be 10s
func (c *Client) newPollingClient() *PollingClient {
	return &PollingClient{
		client:       c,
		shortCircuit: make(chan struct{}, 1),
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
	tick := time.NewTicker(defaultPollingInterval)
	defer tick.Stop()
	log := zerolog.Ctx(ctx)
	for {
		select {
		case <-tick.C:
			pc.poll(ctx)
		case <-pc.shortCircuit:
			tick.Reset(defaultPollingInterval)
			pc.poll(ctx)
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("Polling context canceled")
			return
		}
	}
}

func (pc *PollingClient) poll(ctx context.Context) {
	userUpdatesQuery := (&payload.DMRequestQuery{}).Default()
	log := zerolog.Ctx(ctx)
	if pc.includeConversationID {
		userUpdatesQuery.ActiveConversationID = pc.activeConversationID
	} else {
		userUpdatesQuery.ActiveConversationID = ""
	}
	pc.includeConversationID = !pc.includeConversationID
	if pc.currentCursor != "" {
		userUpdatesQuery.Cursor = pc.currentCursor
	}

	userUpdatesResponse, err := pc.client.GetDMUserUpdates(ctx, &userUpdatesQuery)
	if err != nil {
		log.Err(err).Msg("Failed to get user updates")
		time.Sleep(1 * time.Minute)
		return
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

func (pc *PollingClient) SetCurrentCursor(cursor string) {
	pc.currentCursor = cursor
}

func (pc *PollingClient) stopPolling() {
	if cancel := pc.stop.Swap(nil); cancel != nil {
		(*cancel)()
	}
}

func (pc *PollingClient) SetActiveConversation(conversationID string) {
	pc.activeConversationID = conversationID
	pc.pollConversation(conversationID)
}

func (pc *PollingClient) pollConversation(conversationID string) {
	if pc.activeConversationID == conversationID {
		pc.includeConversationID = true
		select {
		case <-pc.shortCircuit:
		default:
		}
		pc.shortCircuit <- struct{}{}
	}
}
