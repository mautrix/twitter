// Deprecated: This file contains legacy Twitter DM polling.
// Use XChat WebSocket (websocket.go) for real-time messaging.
package twittermeow

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/exsync"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

var defaultPollingInterval = 10 * time.Second

type PollingClient struct {
	client                *Client
	stop                  atomic.Pointer[context.CancelFunc]
	activeConversationID  string
	includeConversationID bool
	shortCircuit          chan struct{}
	pollingStopped        *exsync.Event
}

// interval is the delay inbetween checking for new updates
// default interval will be 10s
func (c *Client) newPollingClient() *PollingClient {
	pc := &PollingClient{
		client:         c,
		shortCircuit:   make(chan struct{}, 1),
		pollingStopped: exsync.NewEvent(),
	}
	pc.pollingStopped.Set()
	return pc
}

func (pc *PollingClient) startPolling(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	if oldCancel := pc.stop.Swap(&cancel); oldCancel != nil {
		(*oldCancel)()
	}
	pc.pollingStopped.Clear()
	go func() {
		defer func() {
			if pc.stop.Load() == &cancel {
				pc.pollingStopped.Set()
			}
		}()
		pc.doPoll(ctx)
	}()
}

func (pc *PollingClient) doPoll(ctx context.Context) {
	tick := time.NewTicker(defaultPollingInterval)
	defer tick.Stop()
	log := zerolog.Ctx(ctx)
	var failing bool
	backoffInterval := defaultPollingInterval / 2
	for {
		select {
		case <-tick.C:
		case <-pc.shortCircuit:
			tick.Reset(defaultPollingInterval)
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("Polling context canceled")
			return
		}
		err := pc.poll(ctx)
		if ctx.Err() != nil {
			log.Debug().Err(err).Msg("Polling context canceled during event handling")
			return
		} else if err != nil {
			log.Err(err).Msg("Failed to poll for updates")
			authError := IsAuthError(err)
			pc.client.eventHandler(&types.PollingError{Error: err, IsAuth: authError}, nil)
			if authError {
				return
			}
			failing = true
			backoffInterval = min(backoffInterval*2, 180*time.Second)
			if errors.Is(err, ErrRatelimitExceeded) {
				backoffInterval *= 2
			}
			tick.Reset(backoffInterval)
		} else if failing {
			failing = false
			pc.client.eventHandler(&types.PollingError{}, nil)
			backoffInterval = defaultPollingInterval / 2
			tick.Reset(defaultPollingInterval)
		}
	}
}

var errEventHandlerFailed = errors.New("event handler failed")

func (pc *PollingClient) poll(ctx context.Context) error {
	userUpdatesQuery := (&payload.DMRequestQuery{}).Default()
	if pc.includeConversationID {
		userUpdatesQuery.ActiveConversationID = pc.activeConversationID
	} else {
		userUpdatesQuery.ActiveConversationID = ""
	}
	pc.includeConversationID = !pc.includeConversationID
	userUpdatesQuery.Cursor = pc.client.session.PollingCursor

	userUpdatesResponse, err := pc.client.GetDMUserUpdates(ctx, &userUpdatesQuery)
	if err != nil {
		return err
	}

	if !pc.client.eventHandler(nil, userUpdatesResponse.UserEvents) {
		return errEventHandlerFailed
	}
	for _, entry := range userUpdatesResponse.UserEvents.Entries {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		parsed := entry.ParseWithErrorLog(&pc.client.Logger)
		if parsed != nil {
			if !pc.client.eventHandler(parsed, userUpdatesResponse.UserEvents) {
				return errEventHandlerFailed
			}
		}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if userUpdatesResponse.UserEvents.Cursor != pc.client.session.PollingCursor {
		pc.client.session.PollingCursor = userUpdatesResponse.UserEvents.Cursor
		pc.client.onCursorChanged(ctx)
	}
	return nil
}

func (pc *PollingClient) stopPolling() {
	if cancel := pc.stop.Load(); cancel != nil {
		(*cancel)()
	}
	pc.activeConversationID = ""
}

func (pc *PollingClient) SetActiveConversation(conversationID string) {
	pc.activeConversationID = conversationID
	pc.pollConversation(conversationID)
}

func (pc *PollingClient) doShortCircuit() {
	select {
	case pc.shortCircuit <- struct{}{}:
	default:
	}
}

func (pc *PollingClient) pollConversation(conversationID string) {
	if pc.activeConversationID == conversationID {
		pc.includeConversationID = true
		pc.doShortCircuit()
	}
}
