package twittermeow

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"
	thrifter "github.com/thrift-iterator/go"
	thriftergeneral "github.com/thrift-iterator/go/general"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

// Reconnect configuration
const (
	initialReconnectDelay      = 1 * time.Second
	maxReconnectDelay          = 5 * time.Minute
	reconnectBackoffMultiplier = 2.0
)

// decodeXChatPayload tries binary first, then compact thrift decoding.
func decodeXChatPayload(data []byte) (out *payload.Message, err error) {
	decoder := thrifter.NewDecoder(bytes.NewReader(data), nil)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("thrift decode panic: %v", r)
		}
	}()

	var decoded payload.Message
	if err = decoder.Decode(&decoded); err != nil {
		return nil, fmt.Errorf("thrift binary decode failed (no envelope expected): %w", err)
	}

	return &decoded, nil
}

// decodeXChatPayloadGeneric decodes without a schema for debugging.
func decodeXChatPayloadGeneric(data []byte) (out thriftergeneral.Struct, err error) {
	decoder := thrifter.NewDecoder(bytes.NewReader(data), nil)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("generic thrift decode panic: %v", r)
		}
	}()
	err = decoder.Decode(&out)
	return out, err
}

func encodeXChatPayload(msg *payload.Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("xchat payload is nil")
	}
	var buffer bytes.Buffer
	encoder := thrifter.NewEncoder(&buffer)
	if err := encoder.Encode(msg); err != nil {
		return nil, fmt.Errorf("thrift encode: %w", err)
	}
	return buffer.Bytes(), nil
}

type xchatWebsocketClient struct {
	client *Client

	shouldStop atomic.Pointer[context.CancelFunc]
	conn       atomic.Pointer[websocket.Conn]
	writeMu    sync.Mutex
}

func (c *Client) newXChatWebsocketClient() *xchatWebsocketClient {
	return &xchatWebsocketClient{client: c}
}

func (c *Client) StartXChatWebsocket(ctx context.Context) error {
	token, err := c.GenerateXChatToken(ctx)

	if err != nil {
		return err
	}

	return c.xchat.start(ctx, token.Data.UserGetXChatAuthToken.Token)
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

func (xc *xchatWebsocketClient) start(ctx context.Context, initialToken string) error {
	if initialToken == "" {
		return fmt.Errorf("xchat token must not be empty")
	}

	ctx, cancel := context.WithCancel(ctx)
	if old := xc.shouldStop.Swap(&cancel); old != nil {
		(*old)()
	}

	log := xc.client.Logger.With().Str("component", "xchat_websocket").Logger()

	go func() {
		defer func() {
			if xc.shouldStop.Load() == &cancel {
				xc.shouldStop.Swap(nil)
			}
		}()

		token := initialToken
		reconnectDelay := initialReconnectDelay
		forceRefresh := false

		for {
			// Check if we should stop before attempting connection
			if ctx.Err() != nil {
				log.Debug().Msg("XChat websocket stopping (context cancelled)")
				return
			}

			if forceRefresh {
				newToken, err := xc.client.refreshXChatToken(ctx)
				if err != nil {
					log.Err(err).Msg("Failed to generate new XChat token after close frame, will retry")
					reconnectDelay = min(time.Duration(float64(reconnectDelay)*reconnectBackoffMultiplier), maxReconnectDelay)
					select {
					case <-ctx.Done():
						log.Debug().Msg("XChat websocket stopping during reconnect wait")
						return
					case <-time.After(reconnectDelay):
					}
					continue
				}
				token = newToken
				reconnectDelay = initialReconnectDelay
				forceRefresh = false
			}

			// Run connection (blocks until disconnect)
			closeReceived, err := xc.runConnection(ctx, token, log)

			// Check if intentionally stopped
			if ctx.Err() != nil {
				log.Debug().Msg("XChat websocket stopped (context cancelled)")
				return
			}

			// Log disconnect, wait with backoff before reconnecting
			log.Warn().Err(err).Dur("retry_in", reconnectDelay).Msg("XChat websocket disconnected, will reconnect")

			select {
			case <-ctx.Done():
				log.Debug().Msg("XChat websocket stopping during reconnect wait")
				return
			case <-time.After(reconnectDelay):
			}

			if closeReceived {
				forceRefresh = true
				continue
			}

			// Get fresh token before reconnecting (tokens expire)
			newToken, err := xc.client.GetXChatToken(ctx)
			if err != nil {
				log.Err(err).Msg("Failed to get new XChat token for reconnect, using previous token")
				// Increase delay since we couldn't get a new token
				reconnectDelay = min(time.Duration(float64(reconnectDelay)*reconnectBackoffMultiplier), maxReconnectDelay)
			} else {
				token = newToken
				// Reset delay on successful token fetch
				reconnectDelay = initialReconnectDelay
			}
		}
	}()

	log.Info().Msg("Initialized XChat Websocket Connection")

	return nil
}

func (xc *xchatWebsocketClient) send(ctx context.Context, msg *payload.Message) error {
	conn := xc.conn.Load()
	if conn == nil {
		return errors.New("xchat websocket not connected")
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
// It returns when the connection is lost for any reason and whether a close frame was received.
func (xc *xchatWebsocketClient) runConnection(ctx context.Context, token string, log zerolog.Logger) (bool, error) {
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
		} else {
			log.Err(err).Msg("Failed to dial XChat websocket")
		}
		return false, fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "shutdown")

	log.Info().Str("url", wsURL.String()).Msg("Connected to XChat websocket")
	xc.conn.Store(conn)
	defer func() {
		if xc.conn.Load() == conn {
			xc.conn.Store(nil)
		}
	}()

	// Create a context for the ping goroutine that we can cancel on read error
	pingCtx, pingCancel := context.WithCancel(ctx)
	defer pingCancel()

	pingTicker := time.NewTicker(30 * time.Second)
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
				var buffer bytes.Buffer
				encoder := thrifter.NewEncoder(&buffer)
				if err := encoder.Encode(instruction); err != nil {
					log.Err(err).Msg("Failed to encode XChat Ping Instruction")
					continue
				}
				bytes := buffer.Bytes()
				xc.writeMu.Lock()
				err := conn.Write(pingCtx, websocket.MessageBinary, bytes)
				xc.writeMu.Unlock()
				if err != nil {
					log.Warn().Err(err).Msg("Failed to send XChat ping frame")
					return
				}
				log.Debug().Int("bytes", len(bytes)).Msg("Sent XChat ping frame")
			}
		}
	}()

	for {
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			status := websocket.CloseStatus(err)
			closeReceived := status != -1
			if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
				log.Debug().Err(err).Uint32("status", uint32(status)).Msg("XChat websocket closed by server")
			} else if ctx.Err() != nil {
				log.Debug().Err(err).Msg("XChat websocket read stopped by context")
			} else {
				log.Error().Err(err).Uint32("status", uint32(status)).Msg("XChat websocket read failed")
			}
			return closeReceived, fmt.Errorf("read failed: %w", err)
		}

		// Handle message. Currently we just trace-log; hook processing here as formats become known.
		log.Debug().
			Str("type", msgType.String()).
			Int("bytes", len(data)).
			Msg("Received XChat websocket message")

		if msgType != websocket.MessageBinary {
			log.Debug().Str("text", string(data)).Msg("Skipping non-binary XChat websocket frame")
			continue
		}
		if len(data) == 0 {
			log.Debug().Msg("Skipping empty XChat websocket frame")
			continue
		}

		decoded, err := decodeXChatPayload(data)
		if err != nil {
			prefixLen := min(32, len(data))
			log.Warn().
				Err(err).
				Int("bytes", len(data)).
				Str("hex_prefix", hex.EncodeToString(data[:prefixLen])).
				Msg("Failed to decode XChat websocket payload")

			if gen, gerr := decodeXChatPayloadGeneric(data); gerr == nil {
				log.Debug().
					Interface("generic", gen).
					Msg("XChat websocket payload (generic decode)")
			} else {
				log.Trace().Err(gerr).Msg("Generic thrift decode also failed")
			}
			continue
		}

		// Log the decoded message for debugging
		if log.Debug().Enabled() {
			if pretty, err := json.MarshalIndent(decoded, "", "  "); err != nil {
				log.Debug().
					Err(err).
					Interface("event", decoded).
					Msg("Decoded XChat websocket payload (failed to format JSON)")
			} else {
				log.Debug().Msgf("Decoded XChat websocket payload:\n%s", pretty)
			}
		}

		// Process the message through the XChat processor
		if err := xc.client.xchatProcessor.ProcessMessage(ctx, decoded); err != nil {
			log.Err(err).Msg("Failed to process XChat message")
		}
	}
}

func (xc *xchatWebsocketClient) stop() {
	if cancel := xc.shouldStop.Load(); cancel != nil {
		(*cancel)()
	}
}
