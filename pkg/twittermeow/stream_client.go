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

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type StreamClient struct {
	client *Client

	stopStream        atomic.Pointer[context.CancelFunc]
	stopHeartbeat     atomic.Pointer[context.CancelFunc]
	heartbeatInterval time.Duration
}

func (c *Client) newStreamClient() *StreamClient {
	return &StreamClient{
		client:            c,
		heartbeatInterval: 25 * time.Second,
	}
}

func (sc *StreamClient) startEventStream(conversationID string) {
	ctx, cancel := context.WithCancel(context.Background())
	if oldCancel := sc.stopStream.Swap(&cancel); oldCancel != nil {
		(*oldCancel)()
	}
	go sc.subscribe(ctx, conversationID)
}

func (sc *StreamClient) subscribe(ctx context.Context, conversationID string) {
	eventsURL, err := url.Parse(endpoints.PIPELINE_EVENTS_URL)
	if err != nil {
		sc.client.Logger.Err(err)
		return
	}

	q := url.Values{
		"topic": []string{getSubscriptionTopic(conversationID)},
	}
	eventsURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, eventsURL.String(), nil)
	if err != nil {
		sc.client.Logger.Err(err)
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
		sc.client.Logger.Err(err).Msg("failed to connect to event stream")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sc.client.Logger.Debug().Int("code", resp.StatusCode).Str("status", resp.Status).Msg("failed to connect to event stream")
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
			sc.client.Logger.Warn().Str("field", field).Str("value", value).Msg("unhandled stream event")
			continue
		}
		var evt response.StreamEvent
		err := json.Unmarshal([]byte(value), &evt)
		if err != nil {
			sc.client.Logger.Err(err).Str("value", value).Msg("error decoding stream event")
			continue
		}
		sc.client.Logger.Trace().Any("evt", evt).Msg("stream event")

		config := evt.Payload.Config
		if config != nil {
			if config.HeartbeatMillis > 0 {
				sc.heartbeatInterval = time.Duration(config.HeartbeatMillis) * time.Millisecond
			}
			if config.SessionID != "" {
				go sc.startHeartbeat(config.SessionID, conversationID)
			}
		} else {
			sc.client.streamEventHandler(evt)
		}
	}
}

func getSubscriptionTopic(conversationID string) string {
	return "/dm_update/" + conversationID + ",/dm_typing/" + conversationID
}

func (sc *StreamClient) startHeartbeat(sessionID, conversationID string) {
	ctx, cancel := context.WithCancel(context.Background())
	if oldCancel := sc.stopHeartbeat.Swap(&cancel); oldCancel != nil {
		(*oldCancel)()
	}

	tick := time.NewTicker(defaultPollingInterval)
	defer tick.Stop()

	sc.heartbeat(ctx, sessionID, conversationID)
	for {
		select {
		case <-tick.C:
			sc.heartbeat(ctx, sessionID, conversationID)
		case <-ctx.Done():
			return
		}
	}
}

func (sc *StreamClient) heartbeat(ctx context.Context, sessionID, conversationID string) {
	payload := &payload.UpdateSubscriptionsPayload{
		SubTopics: getSubscriptionTopic(conversationID),
	}
	encodedPayload, err := payload.Encode()
	if err != nil {
		sc.client.Logger.Err(err)
		return
	}

	_, _, err = sc.client.makeAPIRequest(ctx, apiRequestOpts{
		URL:         endpoints.PIPELINE_UPDATE_URL,
		Method:      http.MethodPost,
		ContentType: types.ContentTypeForm,
		Body:        encodedPayload,
		Headers: map[string]string{
			"livepipeline-session": sessionID,
		},
	})
	if err != nil {
		sc.client.Logger.Err(err)
	}
}
