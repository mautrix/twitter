package twittermeow

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type StreamClient struct {
	client *Client

	stop              atomic.Pointer[context.CancelFunc]
	oldConversationID string
	conversationID    string
	sessionID         string
	heartbeatInterval time.Duration
	shortCircuit      chan struct{}
}

func (c *Client) newStreamClient() *StreamClient {
	return &StreamClient{
		client:            c,
		heartbeatInterval: 25 * time.Second,
		shortCircuit:      make(chan struct{}, 1),
	}
}

func (sc *StreamClient) startOrUpdateEventStream(conversationID string) {
	ctx := sc.client.Logger.With().Str("action", "event stream").Logger().WithContext(context.Background())
	if sc.conversationID == "" {
		sc.conversationID = conversationID
		go sc.start(ctx)
	} else {
		sc.oldConversationID = sc.conversationID
		sc.conversationID = conversationID
		select {
		case <-sc.shortCircuit:
		default:
		}
		sc.shortCircuit <- struct{}{}
	}
}

func (sc *StreamClient) start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	if oldCancel := sc.stop.Swap(&cancel); oldCancel != nil {
		(*oldCancel)()
	}
	eventsURL, err := url.Parse(endpoints.PIPELINE_EVENTS_URL)
	if err != nil {
		zerolog.Ctx(ctx).Err(err)
		return
	}

	q := url.Values{
		"topic": []string{getSubscriptionTopic(sc.conversationID)},
	}
	eventsURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, eventsURL.String(), nil)
	if err != nil {
		zerolog.Ctx(ctx).Err(err)
		return
	}

	extraHeaders := map[string]string{
		"accept":        "text/event-stream",
		"cache-control": "no-cache",
	}
	headerOpts := HeaderOpts{
		WithCookies: true,
		Extra:       extraHeaders,
	}
	req.Header = sc.client.buildHeaders(headerOpts)

	resp, err := sc.client.HTTP.Do(req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to connect to event stream")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		zerolog.Ctx(ctx).Debug().Int("code", resp.StatusCode).Str("status", resp.Status).Msg("failed to connect to event stream")
		return
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line == ":" {
			continue
		}
		index := strings.Index(line, ":")
		field := line[:index]
		value := line[index+1:]
		if field != "data" {
			zerolog.Ctx(ctx).Warn().Str("field", field).Str("value", value).Msg("unhandled stream event")
			continue
		}
		var evt response.StreamEvent
		err := json.Unmarshal([]byte(value), &evt)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Str("value", value).Msg("error decoding stream event")
			continue
		}
		zerolog.Ctx(ctx).Trace().Any("evt", evt).Msg("stream event")

		config := evt.Payload.Config
		if config != nil {
			if config.HeartbeatMillis > 0 {
				sc.heartbeatInterval = time.Duration(config.HeartbeatMillis) * time.Millisecond
			}
			if config.SessionID != "" {
				noHeartbeat := sc.sessionID == ""
				sc.sessionID = config.SessionID
				if noHeartbeat {
					go sc.startHeartbeat(ctx)
				}
			}
		} else {
			sc.client.streamEventHandler(evt)
		}
	}
}

func getSubscriptionTopic(conversationID string) string {
	return "/dm_update/" + conversationID + ",/dm_typing/" + conversationID
}

func (sc *StreamClient) startHeartbeat(ctx context.Context) {
	tick := time.NewTicker(sc.heartbeatInterval)
	defer tick.Stop()

	sc.heartbeat(ctx)
	for {
		select {
		case <-tick.C:
			sc.heartbeat(ctx)
		case <-sc.shortCircuit:
			tick.Reset(sc.heartbeatInterval)
			sc.heartbeat(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (sc *StreamClient) heartbeat(ctx context.Context) {
	payload := &payload.UpdateSubscriptionsPayload{
		SubTopics:   getSubscriptionTopic(sc.conversationID),
		UnsubTopics: getSubscriptionTopic(sc.oldConversationID),
	}
	if sc.oldConversationID != "" {
		sc.oldConversationID = ""
	}
	encodedPayload, err := payload.Encode()
	if err != nil {
		zerolog.Ctx(ctx).Err(err)
		return
	}

	_, _, err = sc.client.makeAPIRequest(ctx, apiRequestOpts{
		URL:         endpoints.PIPELINE_UPDATE_URL,
		Method:      http.MethodPost,
		ContentType: types.ContentTypeForm,
		Body:        encodedPayload,
		Headers: map[string]string{
			"livepipeline-session": sc.sessionID,
		},
	})
	if err != nil {
		zerolog.Ctx(ctx).Err(err)
	}
}

func (sc *StreamClient) stopStream() {
	if cancel := sc.stop.Swap(nil); cancel != nil {
		(*cancel)()
	}
	sc.conversationID = ""
	sc.sessionID = ""
}
