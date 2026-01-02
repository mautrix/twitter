// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2025 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package realm implements HTTP communication with Juicebox realms.
package realm

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/noise"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/requests"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/types"
)

const (
	juiceboxVersionHeader = "X-Juicebox-Version"
	juiceboxVersion       = "0.3.4"
	userAgent             = "JuiceboxSdk-Go/0.3.4"
)

// AuthTokenProvider retrieves auth tokens for realms.
type AuthTokenProvider func(realmID types.RealmID) (types.AuthToken, bool)

// Client handles communication with a single realm.
type Client struct {
	realm             types.Realm
	httpClient        *http.Client
	authTokenProvider AuthTokenProvider

	// Session state for hardware realms
	sessionMu sync.Mutex
	session   *Session
}

// Session holds an active Noise session with a realm.
type Session struct {
	ID        types.SessionID
	Transport *noise.Transport
	Lifetime  time.Duration
	LastUsed  time.Time
}

// NewClient creates a new realm client.
func NewClient(realm types.Realm, httpClient *http.Client, authTokenProvider AuthTokenProvider) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		realm:             realm,
		httpClient:        httpClient,
		authTokenProvider: authTokenProvider,
	}
}

// MakeRequest sends a SecretsRequest to the realm.
func (c *Client) MakeRequest(ctx context.Context, req *requests.SecretsRequest) (*requests.SecretsResponse, error) {
	if c.realm.PublicKey != nil {
		return c.makeHardwareRealmRequest(ctx, req)
	}
	return c.makeSoftwareRealmRequest(ctx, req)
}

// makeSoftwareRealmRequest sends a direct request to a software realm.
func (c *Client) makeSoftwareRealmRequest(ctx context.Context, req *requests.SecretsRequest) (*requests.SecretsResponse, error) {
	authToken, ok := c.authTokenProvider(c.realm.ID)
	if !ok {
		return nil, types.ErrInvalidAuth
	}

	body, err := requests.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.realm.Address, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/cbor")
	httpReq.Header.Set("Authorization", "Bearer "+string(authToken))
	httpReq.Header.Set(juiceboxVersionHeader, juiceboxVersion)
	httpReq.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, types.ErrTransient
	}
	defer resp.Body.Close()

	if err := checkStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.ErrTransient
	}

	var secretsResp requests.SecretsResponse
	if err := requests.Unmarshal(respBody, &secretsResp); err != nil {
		return nil, types.ErrAssertion
	}

	return &secretsResp, nil
}

// makeHardwareRealmRequest sends a Noise-encrypted request to a hardware realm.
func (c *Client) makeHardwareRealmRequest(ctx context.Context, req *requests.SecretsRequest) (*requests.SecretsResponse, error) {
	needsForwardSecrecy := req.Recover2 != nil || req.Recover3 != nil
	reqBytes, err := requests.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()

	// Try up to 5 times (with potential session re-establishment)
	for attempt := 1; attempt <= 5; attempt++ {
		var respBytes []byte
		var err error

		// Check if we have a valid session
		session := c.session
		if session != nil && time.Since(session.LastUsed) >= session.Lifetime {
			session = nil
			c.session = nil
		}

		if session == nil && needsForwardSecrecy {
			// Need forward secrecy but no session - establish session first
			session, _, err = c.makeHandshakeRequest(ctx, nil)
			if err != nil {
				if err == types.ErrTransient && attempt < 5 {
					time.Sleep(time.Duration(attempt*5) * time.Millisecond)
					continue
				}
				return nil, err
			}
			c.session = session
		}

		if session == nil {
			// No session and don't need forward secrecy - piggyback on handshake
			var newSession *Session
			newSession, respBytes, err = c.makeHandshakeRequest(ctx, reqBytes)
			if err != nil {
				if err == types.ErrTransient && attempt < 5 {
					time.Sleep(time.Duration(attempt*5) * time.Millisecond)
					continue
				}
				return nil, err
			}
			c.session = newSession
		} else {
			// Have session - use transport
			respBytes, err = c.makeTransportRequest(ctx, session, reqBytes)
			if err == errMissingSession {
				c.session = nil
				continue
			}
			if err != nil {
				if err == types.ErrTransient && attempt < 5 {
					time.Sleep(time.Duration(attempt*5) * time.Millisecond)
					continue
				}
				return nil, err
			}
			session.LastUsed = time.Now()
		}

		// Parse padded response
		var paddedResp requests.PaddedSecretsResponse
		if err := requests.Unmarshal(respBytes, &paddedResp); err != nil {
			return nil, types.ErrAssertion
		}

		secretsResp, err := requests.UnmarshalSecretsResponse(&paddedResp)
		if err != nil {
			return nil, types.ErrAssertion
		}

		return secretsResp, nil
	}

	return nil, types.ErrTransient
}

var errMissingSession = fmt.Errorf("missing session")

func (c *Client) makeHandshakeRequest(ctx context.Context, payload []byte) (*Session, []byte, error) {
	authToken, ok := c.authTokenProvider(c.realm.ID)
	if !ok {
		return nil, nil, types.ErrInvalidAuth
	}

	handshake, handshakeReq, err := noise.Start(c.realm.PublicKey, payload)
	if err != nil {
		return nil, nil, types.ErrAssertion
	}

	sessionID := randomSessionID()
	kind := requests.ClientRequestKindHandshakeOnly
	if len(payload) > 0 {
		kind = requests.ClientRequestKindSecretsRequest
	}

	clientReq := &requests.ClientRequest{
		Realm:     c.realm.ID,
		AuthToken: authToken,
		SessionID: sessionID,
		Kind:      kind,
		Encrypted: requests.NoiseRequest{
			Handshake: handshakeReq,
		},
	}

	clientResp, err := c.sendClientRequest(ctx, clientReq)
	if err != nil {
		return nil, nil, err
	}

	if clientResp.Ok == nil || clientResp.Ok.Handshake == nil {
		return nil, nil, types.ErrAssertion
	}

	transport, respPayload, err := handshake.Finish(&clientResp.Ok.Handshake.Handshake)
	if err != nil {
		return nil, nil, types.ErrAssertion
	}

	session := &Session{
		ID:        sessionID,
		Transport: transport,
		Lifetime:  clientResp.Ok.Handshake.SessionLifetime,
		LastUsed:  time.Now(),
	}

	return session, respPayload, nil
}

func (c *Client) makeTransportRequest(ctx context.Context, session *Session, payload []byte) ([]byte, error) {
	authToken, ok := c.authTokenProvider(c.realm.ID)
	if !ok {
		return nil, types.ErrInvalidAuth
	}

	ciphertext, err := session.Transport.Encrypt(payload)
	if err != nil {
		return nil, types.ErrAssertion
	}

	clientReq := &requests.ClientRequest{
		Realm:     c.realm.ID,
		AuthToken: authToken,
		SessionID: session.ID,
		Kind:      requests.ClientRequestKindSecretsRequest,
		Encrypted: requests.NoiseRequest{
			Transport: &requests.NoiseTransportRequest{
				Ciphertext: ciphertext,
			},
		},
	}

	clientResp, err := c.sendClientRequest(ctx, clientReq)
	if err != nil {
		return nil, err
	}

	if clientResp.MissingSession {
		return nil, errMissingSession
	}
	if clientResp.Ok == nil || clientResp.Ok.Transport == nil {
		return nil, types.ErrAssertion
	}

	return session.Transport.Decrypt(clientResp.Ok.Transport.Ciphertext)
}

func (c *Client) sendClientRequest(ctx context.Context, req *requests.ClientRequest) (*requests.ClientResponse, error) {
	body, err := requests.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.realm.Address, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/cbor")
	httpReq.Header.Set(juiceboxVersionHeader, juiceboxVersion)
	httpReq.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, types.ErrTransient
	}
	defer resp.Body.Close()

	if err := checkStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.ErrTransient
	}

	var clientResp requests.ClientResponse
	if err := requests.Unmarshal(respBody, &clientResp); err != nil {
		return nil, types.ErrAssertion
	}

	if clientResp.InvalidAuth {
		return nil, types.ErrInvalidAuth
	}
	if clientResp.RateLimitExceeded {
		return nil, types.ErrRateLimitExceeded
	}
	if clientResp.Unavailable {
		return nil, types.ErrTransient
	}

	return &clientResp, nil
}

func checkStatusCode(code int) error {
	switch code {
	case 200:
		return nil
	case 401:
		return types.ErrInvalidAuth
	case 426:
		return types.ErrUpgradeRequired
	case 429:
		return types.ErrRateLimitExceeded
	default:
		return types.ErrTransient
	}
}

func randomSessionID() types.SessionID {
	var buf [4]byte
	rand.Read(buf[:])
	return types.SessionID(binary.LittleEndian.Uint32(buf[:]))
}
