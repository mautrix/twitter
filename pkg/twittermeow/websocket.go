package twittermeow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

// Reconnect configuration
const (
	initialReconnectDelay      = 1 * time.Second
	maxReconnectDelay          = 5 * time.Minute
	reconnectBackoffMultiplier = 2.0
	stableConnectionDuration   = 30 * time.Second
	xchatPingInterval          = 30 * time.Second
	xchatPingTimeout           = 10 * time.Second
)

var ErrXChatWebsocketNotConnected = errors.New("xchat websocket not connected")

// decodeXChatPayload decodes binary thrift data.
func decodeXChatPayload(data []byte) (out *payload.Message, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("thrift decode panic: %v", r)
		}
	}()

	var decoded payload.Message
	if err = payload.Decode(data, &decoded); err != nil {
		return nil, fmt.Errorf("thrift binary decode failed (no envelope expected): %w", err)
	}

	return &decoded, nil
}

func encodeXChatPayload(msg *payload.Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("xchat payload is nil")
	}
	data, err := payload.Encode(msg)
	if err != nil {
		return nil, fmt.Errorf("thrift encode: %w", err)
	}
	return data, nil
}

type xchatWebsocketClient struct {
	client *Client

	shouldStop atomic.Pointer[context.CancelFunc]
	conn       atomic.Pointer[websocket.Conn]
	writeMu    sync.Mutex

	tokenProvider     func(context.Context, bool) (string, error)
	connectionRunner  func(context.Context, string, zerolog.Logger, func(context.Context, XChatLiveDrain) error) (bool, error)
	initialRetryDelay time.Duration
	maximumRetryDelay time.Duration
}

func (c *Client) newXChatWebsocketClient() *xchatWebsocketClient {
	xc := &xchatWebsocketClient{
		client:            c,
		initialRetryDelay: initialReconnectDelay,
		maximumRetryDelay: maxReconnectDelay,
	}
	xc.tokenProvider = xc.getToken
	xc.connectionRunner = xc.runConnection
	return xc
}

func (c *Client) StartXChatWebsocket(ctx context.Context) error {
	if c.xchat == nil {
		return errors.New("xchat websocket not initialized")
	}
	return c.xchat.start(ctx)
}

func (c *Client) SendXChatPayload(ctx context.Context, msg *payload.Message) error {
	if c.xchat == nil {
		return errors.New("xchat websocket not initialized")
	}
	return c.xchat.send(ctx, msg)
}

// stopXChatWebsocket stops any active XChat websocket connection.
func (c *Client) stopXChatWebsocket() {
	if c.xchat != nil {
		c.xchat.stop()
	}
}

func (xc *xchatWebsocketClient) start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	if old := xc.shouldStop.Swap(&cancel); old != nil {
		(*old)()
	}

	log := xc.client.Logger.With().Str("component", "xchat_websocket").Logger()
	ready := make(chan struct{})

	go func() {
		defer func() {
			xc.shouldStop.CompareAndSwap(&cancel, nil)
		}()

		var readyOnce sync.Once
		var token string
		initialRetryDelay := xc.initialRetryDelay
		if initialRetryDelay <= 0 {
			initialRetryDelay = initialReconnectDelay
		}
		reconnectDelay := initialRetryDelay
		maximumRetryDelay := xc.maximumRetryDelay
		if maximumRetryDelay <= 0 {
			maximumRetryDelay = maxReconnectDelay
		}
		forceRefresh := false
		for {
			// Check if we should stop before attempting connection
			if ctx.Err() != nil {
				log.Debug().Msg("XChat websocket stopping (context cancelled)")
				return
			}

			if token == "" || forceRefresh {
				newToken, err := xc.tokenProvider(ctx, forceRefresh)
				if err == nil && newToken == "" {
					forceRefresh = true
					err = errors.New("xchat token is empty")
				}
				if err != nil {
					log.Err(err).
						Bool("force_refresh", forceRefresh).
						Dur("retry_in", reconnectDelay).
						Msg("Failed to get XChat token, will retry")
					if !waitForXChatReconnect(ctx, reconnectDelay) {
						log.Debug().Msg("XChat websocket stopping during reconnect wait")
						return
					}
					reconnectDelay = nextXChatReconnectDelay(reconnectDelay, maximumRetryDelay)
					continue
				}
				token = newToken
				forceRefresh = false
			}

			// Run connection (blocks until disconnect)
			readyThisAttempt := false
			var connectedAt time.Time
			refreshToken, err := xc.runConnectionAttempt(ctx, token, log, func(connectionCtx context.Context, drain XChatLiveDrain) error {
				if xc.client.xchatProcessor != nil {
					xc.client.xchatProcessor.MarkReconnected()
				}
				if connectHandler := xc.client.getXChatConnectHandler(); connectHandler != nil {
					log.Info().Msg("Running XChat socket handoff catch-up")
					if err := connectHandler(connectionCtx, drain); err != nil {
						return fmt.Errorf("socket handoff catch-up failed: %w", err)
					}
					log.Info().Msg("Finished XChat socket handoff catch-up")
				}
				readyThisAttempt = true
				connectedAt = time.Now()
				readyOnce.Do(func() {
					close(ready)
				})
				return nil
			})

			// Check if intentionally stopped
			if ctx.Err() != nil {
				log.Debug().Msg("XChat websocket stopped (context cancelled)")
				return
			}

			// A successful handshake followed by an immediate close is not a
			// healthy connection. Only reset backoff after a stable live period.
			waitDelay, nextDelay := xchatReconnectBackoff(
				reconnectDelay,
				initialRetryDelay,
				maximumRetryDelay,
				readyThisAttempt,
				connectedAt,
				time.Now(),
			)
			log.Warn().Err(err).Dur("retry_in", waitDelay).Msg("XChat websocket disconnected, will reconnect")

			if !waitForXChatReconnect(ctx, waitDelay) {
				log.Debug().Msg("XChat websocket stopping during reconnect wait")
				return
			}
			reconnectDelay = nextDelay

			if refreshToken {
				forceRefresh = true
				token = ""
				continue
			}

			token = ""
		}
	}()

	select {
	case <-ready:
		log.Info().Msg("Initialized XChat Websocket Connection")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (xc *xchatWebsocketClient) runConnectionAttempt(
	ctx context.Context,
	token string,
	log zerolog.Logger,
	onConnected func(context.Context, XChatLiveDrain) error,
) (refreshToken bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Error().
				Bytes("stack", debug.Stack()).
				Str("panic_type", fmt.Sprintf("%T", recovered)).
				Msg("XChat connection attempt panicked")
			refreshToken = false
			err = fmt.Errorf("xchat connection attempt panicked (%T)", recovered)
		}
	}()
	return xc.connectionRunner(ctx, token, log, onConnected)
}

func (xc *xchatWebsocketClient) getToken(ctx context.Context, forceRefresh bool) (string, error) {
	if forceRefresh {
		return xc.client.refreshXChatToken(ctx)
	}
	return xc.client.GetXChatToken(ctx)
}

func waitForXChatReconnect(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func nextXChatReconnectDelay(current, maximum time.Duration) time.Duration {
	return min(time.Duration(float64(current)*reconnectBackoffMultiplier), maximum)
}

func xchatReconnectBackoff(
	current, initial, maximum time.Duration,
	connected bool,
	connectedAt, disconnectedAt time.Time,
) (wait, next time.Duration) {
	if connected && !connectedAt.IsZero() && disconnectedAt.Sub(connectedAt) >= stableConnectionDuration {
		return initial, initial
	}
	return current, nextXChatReconnectDelay(current, maximum)
}

func (xc *xchatWebsocketClient) send(ctx context.Context, msg *payload.Message) error {
	conn := xc.conn.Load()
	if conn == nil {
		return ErrXChatWebsocketNotConnected
	}
	bytes, err := encodeXChatPayload(msg)
	if err != nil {
		return err
	}
	xc.writeMu.Lock()
	err = conn.Write(ctx, websocket.MessageBinary, bytes)
	xc.writeMu.Unlock()
	if err != nil {
		return fmt.Errorf("xchat websocket write failed: %w", err)
	}
	return nil
}

// runConnection handles a single WebSocket connection lifecycle.
// It returns when the connection is lost and whether the next attempt should
// force-refresh the short-lived XChat token.
func (xc *xchatWebsocketClient) runConnection(
	ctx context.Context,
	token string,
	log zerolog.Logger,
	onConnected func(context.Context, XChatLiveDrain) error,
) (bool, error) {
	wsURL, err := url.Parse(endpoints.XCHAT_WEBSOCKET_URL)
	if err != nil {
		return false, fmt.Errorf("failed to parse websocket URL: %w", err)
	}
	q := wsURL.Query()
	q.Set("token", token)
	wsURL.RawQuery = q.Encode()

	headers := xc.client.buildHeaders(HeaderOpts{
		WithCookies:         true,
		WithAuthBearer:      true,
		WithXCsrfToken:      true,
		WithXGuestToken:     true,
		WithXTwitterHeaders: true,
		Origin:              endpoints.BASE_URL,
		Referer:             endpoints.BASE_MESSAGES_URL,
	})

	conn, resp, err := websocket.Dial(ctx, wsURL.String(), &websocket.DialOptions{
		HTTPHeader: headers,
	})
	if err != nil {
		if resp != nil {
			log.Err(err).Int("status", resp.StatusCode).Msg("Failed to dial XChat websocket")
			return resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
				fmt.Errorf("dial failed with HTTP %d: %w", resp.StatusCode, err)
		} else {
			log.Err(err).Msg("Failed to dial XChat websocket")
		}
		return false, fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "shutdown")

	log.Info().Msg("Connected to XChat websocket")
	xc.conn.Store(conn)
	defer func() {
		xc.conn.CompareAndSwap(conn, nil)
	}()

	// Start keepalives before catch-up: a large inbox repair must not make the
	// newly established socket look idle and get disconnected before reads begin.
	pingCtx, pingCancel := context.WithCancel(ctx)
	defer pingCancel()

	pingTicker := time.NewTicker(xchatPingInterval)
	defer pingTicker.Stop()

	go func() {
		for {
			select {
			case <-pingCtx.Done():
				return
			case <-pingTicker.C:
				instruction := payload.Message{
					MessageInstruction: &payload.MessageInstruction{
						KeepAliveInstruction: &payload.KeepAliveInstruction{},
					},
				}
				data, err := payload.Encode(&instruction)
				if err != nil {
					log.Err(err).Msg("Failed to encode XChat Ping Instruction")
					continue
				}
				writeCtx, writeCancel := context.WithTimeout(pingCtx, xchatPingTimeout)
				xc.writeMu.Lock()
				err = conn.Write(writeCtx, websocket.MessageBinary, data)
				xc.writeMu.Unlock()
				writeCancel()
				if err != nil {
					log.Warn().Err(err).Msg("Failed to send XChat ping frame")
					pingCancel()
					return
				}
				log.Debug().Int("bytes", len(data)).Msg("Sent XChat ping frame")
			}
		}
	}()

	conn.SetReadLimit(32768)
	readCtx, cancelRead := context.WithCancelCause(pingCtx)
	frames := make(chan []byte, 32)
	done := make(chan struct{})
	var readErr error
	var refreshToken bool
	go func() {
		defer close(done)
		defer close(frames)
		for {
			typ, data, err := conn.Read(readCtx)
			if err != nil {
				readErr = fmt.Errorf("read XChat websocket: %w", err)
				status := websocket.CloseStatus(err)
				refreshToken = status != -1
				if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
					log.Debug().Err(err).Uint32("status", uint32(status)).Msg("XChat websocket closed by server")
				} else if readCtx.Err() != nil {
					log.Debug().Err(err).Msg("XChat websocket read stopped by context")
				} else {
					log.Error().Err(err).Uint32("status", uint32(status)).Msg("XChat websocket read failed")
				}
				cancelRead(readErr)
				return
			}
			log.Debug().Str("type", typ.String()).Int("bytes", len(data)).Msg("Received XChat websocket message")
			if typ != websocket.MessageBinary || len(data) == 0 {
				continue
			}
			select {
			case frames <- data:
			case <-readCtx.Done():
				return
			}
		}
	}()
	defer func() { cancelRead(context.Canceled); <-done }()
	result := func(err error) (bool, error) {
		cancelRead(err)
		<-done
		if readErr != nil && (refreshToken || !errors.Is(readErr, context.Canceled)) {
			return refreshToken, readErr
		}
		return false, err
	}
	process := func(data []byte, before func(*payload.Message) error) error {
		if err := readCtx.Err(); err != nil {
			return err
		}
		decoded, err := decodeXChatPayload(data)
		if err != nil {
			log.Warn().Err(err).Int("bytes", len(data)).Msg("Failed to decode XChat websocket payload")
			return err
		}
		batchedEventCount := 0
		if decoded.BatchedMessageEvents != nil {
			batchedEventCount = len(decoded.BatchedMessageEvents.MessageEvents)
		}
		log.Trace().
			Bool("has_message_event", decoded.MessageEvent != nil).
			Bool("has_instruction", decoded.MessageInstruction != nil).
			Int("batched_event_count", batchedEventCount).
			Msg("Decoded XChat websocket payload")
		if before != nil {
			if err = before(decoded); err != nil {
				return err
			}
		}
		err = xc.client.xchatProcessor.ProcessMessage(readCtx, decoded)
		if err != nil {
			log.Err(err).Msg("Failed to process XChat message")
		}
		return err
	}
	receive := func(before func(*payload.Message) error) error {
		select {
		case data, ok := <-frames:
			if !ok {
				return context.Cause(readCtx)
			}
			return process(data, before)
		case <-readCtx.Done():
			return context.Cause(readCtx)
		}
	}
	drain := func(before func(*payload.Message) error) error {
		for range len(frames) {
			if err := receive(before); err != nil {
				return err
			}
		}
		return readCtx.Err()
	}
	if onConnected != nil {
		if err := onConnected(readCtx, drain); err != nil {
			return result(err)
		}
	}
	for {
		if err := receive(nil); err != nil {
			return result(err)
		}
	}
}

func (xc *xchatWebsocketClient) stop() {
	if cancel := xc.shouldStop.Load(); cancel != nil {
		(*cancel)()
	}
}

// XChatLiveDrain calls before, then dispatches, for a finite FIFO prefix.
type XChatLiveDrain func(before func(*payload.Message) error) error

// XChatMessageConversations returns nil when a payload needs a global barrier.
func XChatMessageConversations(message *payload.Message) []string {
	return xchatMessageConversations(message, "")
}

func xchatMessageConversations(message *payload.Message, snapshotConversation string) []string {
	if message == nil {
		return nil
	}
	ids := []string{}
	add := func(id *string) bool {
		if id == nil || *id == "" || (snapshotConversation != "" && *id != snapshotConversation) {
			return false
		}
		ids = append(ids, *id)
		return true
	}
	events := []*payload.MessageEvent{}
	if message.MessageEvent != nil {
		events = append(events, message.MessageEvent)
	}
	if message.BatchedMessageEvents != nil {
		events = append(events, message.BatchedMessageEvents.MessageEvents...)
	}
	if len(events) == 0 && message.MessageInstruction == nil {
		return nil
	}
	for _, event := range events {
		if event == nil || !add(event.ConversationId) {
			return nil
		}
		if event.Detail == nil {
			continue
		}
		if deletion := event.Detail.ConversationDeleteEvent; deletion != nil && (snapshotConversation == "" || !add(deletion.ConversationId)) {
			return nil
		}
		if typing := event.Detail.MessageTypingEvent; typing != nil && typing.ConversationId != nil && !add(typing.ConversationId) {
			return nil
		}
	}
	return ids
}

func XChatInboxItemIsConversationScoped(item *response.XChatInboxItem) bool {
	id := item.ConversationDetail.ConversationID
	if id == "" {
		return false
	}
	encoded := append([]string{}, item.LatestMessageEvents...)
	encoded = append(encoded, item.EncodedMessageEvents...)
	encoded = append(encoded, item.LatestConversationKeyChangeEvents...)
	encoded = append(encoded, item.LatestNotifiableMessageCreateEvent)
	for _, read := range item.LatestReadEventsPerParticipant {
		encoded = append(encoded, read.LatestMarkConversationReadEvent)
	}
	for _, raw := range encoded {
		if raw == "" {
			continue
		}
		event, err := DecodeMessageEvent(raw)
		if err != nil {
			return false
		}
		if event.ConversationId == nil || *event.ConversationId == "" {
			event.ConversationId = &id
		}
		if xchatMessageConversations(&payload.Message{MessageEvent: event}, id) == nil {
			return false
		}
	}
	return true
}
