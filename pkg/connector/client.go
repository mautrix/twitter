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

package connector

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	thrifter "github.com/thrift-iterator/go"
	thriftergeneral "github.com/thrift-iterator/go/general"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/status"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type TwitterClient struct {
	connector *TwitterConnector
	client    *twittermeow.Client

	userLogin *bridgev2.UserLogin

	userCache     map[string]*types.User
	userCacheLock sync.RWMutex

	participantCache map[string][]types.Participant

	matrixParser *format.HTMLParser

	reconnectAttempted atomic.Bool

	ensurePortalLocks sync.Map
}

var _ bridgev2.NetworkAPI = (*TwitterClient)(nil)

func NewTwitterClient(login *bridgev2.UserLogin, connector *TwitterConnector, client *twittermeow.Client) *TwitterClient {
	tc := &TwitterClient{
		connector:        connector,
		client:           client,
		userLogin:        login,
		userCache:        make(map[string]*types.User),
		participantCache: make(map[string][]types.Participant),
	}
	client.SetXChatEventHandler(tc.HandleXChatEvent)
	// Ensure current user ID is available even if cookies omit twid
	meta := ensureUserLoginMetadata(login)
	if meta.UserID != "" {
		client.SetCurrentUserID(meta.UserID)
	} else {
		client.SetCurrentUserID(string(login.ID))
	}
	tc.matrixParser = &format.HTMLParser{
		TabsToSpaces:   4,
		Newline:        "\n",
		HorizontalLine: "\n---\n",
		PillConverter: func(displayname, mxid, eventID string, ctx format.Context) string {
			userID, ok := tc.connector.br.Matrix.ParseGhostMXID(id.UserID(mxid))
			if !ok {
				return displayname
			}
			ghost, err := tc.connector.br.GetGhostByID(context.TODO(), userID)
			if err != nil || len(ghost.Identifiers) < 1 {
				return displayname
			}
			id := ghost.Identifiers[0]
			return "@" + strings.TrimPrefix(id, "twitter:")
		},
	}
	return tc
}

func (tc *TwitterConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	meta := ensureUserLoginMetadata(login)
	c := cookies.NewCookiesFromString(meta.Cookies)
	log := login.Log.With().Str("component", "twitter_client").Logger()
	client := twittermeow.NewClient(c, log)
	if meta.UserID != "" {
		client.SetCurrentUserID(meta.UserID)
	} else {
		client.SetCurrentUserID(string(login.ID))
	}
	client.SetKeyStore(newUserLoginKeyStore(login))
	login.Client = NewTwitterClient(login, tc, client)
	return nil
}

const (
	sessionMaxTimeSinceSave = 24 * time.Hour
	sessionMaxTimeSinceInit = 48 * time.Hour
)

func (tc *TwitterClient) Connect(ctx context.Context) {
	log := zerolog.Ctx(ctx)

	if tc.client == nil {
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateBadCredentials,
			Error:      "twitter-not-logged-in",
		})
		return
	}

	tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnecting})
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)

	// Check for cached session
	useCachedSession := tc.connector.Config.CacheSession &&
		meta.Session != nil &&
		meta.Session.LastSaved.Add(sessionMaxTimeSinceSave).After(time.Now()) &&
		meta.Session.InitializedAt.Add(48*time.Hour).After(time.Now()) &&
		meta.Session.CacheVersion == twittermeow.CurrentCacheVersion

	if useCachedSession {
		log.Debug().
			Time("session_ts", meta.Session.LastSaved).
			Time("session_init_ts", meta.Session.InitializedAt).
			Msg("Connecting with cached session")
		tc.client.SetSession(meta.Session)
	} else {
		// Load messages page to initialize session (populates cookies, tokens, etc.)
		_, err := tc.client.LoadMessagesPage(ctx)
		if err != nil {
			log.Err(err).Msg("Failed to load messages page")
			if twittermeow.IsAuthError(err) {
				tc.userLogin.BridgeState.Send(status.BridgeState{
					StateEvent: status.StateBadCredentials,
					Error:      "twitter-invalid-credentials",
					Message:    err.Error(),
				})
			} else {
				tc.userLogin.BridgeState.Send(status.BridgeState{
					StateEvent: status.StateUnknownError,
					Error:      "twitter-load-error",
				})
			}
			return
		}
	}

	// Set up XChat processor and sequence ID tracking
	processor := tc.client.GetXChatProcessor()
	var maxSeqID string
	var maxSeqIDLock sync.Mutex

	processor.SetSequenceIDCallback(func(seqID string) {
		maxSeqIDLock.Lock()
		defer maxSeqIDLock.Unlock()
		if parseSequenceID(seqID) > parseSequenceID(maxSeqID) {
			maxSeqID = seqID
		}
	})

	// Fetch XChat inbox pages
	fetchLog := log.With().Str("component", "xchat_fetch").Logger()
	seqID := meta.MaxUserSequenceID
	if seqID == "" {
		seqID = "null"
	}
	msgPullVersion := meta.MessagePullVersion

	users := make(map[string]*types.User)
	totalItems := 0
	var pages []response.XChatInboxPage

	// Initial page fetch
	vars := payload.NewInitialXChatPageQueryVariables(seqID)
	if msgPullVersion != nil {
		vars.MessagePullVersion = msgPullVersion
	}

	fetchLog.Info().
		Str("request_cursor_id", "").
		Str("request_graph_snapshot_id", "").
		Msg("Fetching initial XChat inbox page")

	initialResp, err := tc.client.GetInitialXChatPage(ctx, vars)
	if err != nil {
		fetchLog.Err(err).
			Msg("Failed to fetch initial XChat inbox page")
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateUnknownError,
			Error:      "twitter-xchat-fetch-error",
			Message:    err.Error(),
		})
		return
	}

	page := initialResp.Data.GetInboxPage
	fetchLog.Info().
		Str("response_cursor_id", page.InboxCursor.CursorID).
		Str("response_graph_snapshot_id", page.InboxCursor.GraphSnapshotID).
		Bool("pull_finished", page.InboxCursor.PullFinished).
		Msg("Received XChat inbox page")

	if page.MaxUserSequenceID != nil && parseSequenceID(*page.MaxUserSequenceID) > parseSequenceID(maxSeqID) {
		maxSeqID = *page.MaxUserSequenceID
	}
	if page.MessagePullVersion != nil {
		msgPullVersion = page.MessagePullVersion
	}

	pages = append(pages, page)

	var cursor *payload.XChatCursor
	if !page.InboxCursor.PullFinished && page.InboxCursor.CursorID != "" && page.InboxCursor.GraphSnapshotID != "" {
		cursor = &payload.XChatCursor{
			CursorId:        page.InboxCursor.CursorID,
			GraphSnapshotId: page.InboxCursor.GraphSnapshotID,
		}
	}

	// Subsequent pages via GetInboxPageRequest
	for cursor != nil {
		fetchLog.Info().
			Str("cursor_id", cursor.CursorId).
			Str("graph_snapshot_id", cursor.GraphSnapshotId).
			Msg("Fetching XChat inbox page")

		inboxVars := payload.NewInboxPageRequestQueryVariables(cursor)
		resp, err := tc.client.GetInboxPageRequest(ctx, inboxVars)
		if err != nil {
			fetchLog.Err(err).
				Msg("Failed to fetch XChat inbox page")
			tc.userLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateUnknownError,
				Error:      "twitter-xchat-fetch-error",
				Message:    err.Error(),
			})
			return
		}

		page := resp.Data.GetInboxPage
		fetchLog.Info().
			Str("response_cursor_id", page.InboxCursor.CursorID).
			Str("response_graph_snapshot_id", page.InboxCursor.GraphSnapshotID).
			Bool("pull_finished", page.InboxCursor.PullFinished).
			Msg("Received XChat inbox page")

		if page.MaxUserSequenceID != nil && parseSequenceID(*page.MaxUserSequenceID) > parseSequenceID(maxSeqID) {
			maxSeqID = *page.MaxUserSequenceID
		}
		if page.MessagePullVersion != nil {
			msgPullVersion = page.MessagePullVersion
		}

		pages = append(pages, page)

		if page.InboxCursor.PullFinished || page.InboxCursor.CursorID == "" || page.InboxCursor.GraphSnapshotID == "" {
			cursor = nil
			break
		}
		if page.InboxCursor.CursorID == cursor.CursorId && page.InboxCursor.GraphSnapshotID == cursor.GraphSnapshotId {
			fetchLog.Debug().
				Str("cursor_id", page.InboxCursor.CursorID).
				Msg("Cursor did not advance, stopping inbox pagination")
			cursor = nil
			break
		}

		cursor = &payload.XChatCursor{
			CursorId:        page.InboxCursor.CursorID,
			GraphSnapshotId: page.InboxCursor.GraphSnapshotID,
		}
	}

	// Collect users from all pages
	var missingUserIDs []string
	for _, pg := range pages {
		for _, item := range pg.Items {
			for _, p := range item.ConversationDetail.ParticipantsResults {
				if p.Result != nil {
					users[p.RestID] = twittermeow.ConvertXChatUserToUser(p.Result)
				} else if p.RestID != "" {
					missingUserIDs = append(missingUserIDs, p.RestID)
				}
			}
			for _, p := range item.ConversationDetail.GroupMembersResults {
				if p.Result != nil {
					users[p.RestID] = twittermeow.ConvertXChatUserToUser(p.Result)
				} else if p.RestID != "" {
					missingUserIDs = append(missingUserIDs, p.RestID)
				}
			}
			for _, p := range item.ConversationDetail.GroupAdminsResults {
				if p.Result != nil {
					users[p.RestID] = twittermeow.ConvertXChatUserToUser(p.Result)
				} else if p.RestID != "" {
					missingUserIDs = append(missingUserIDs, p.RestID)
				}
			}
		}
	}

	if err := tc.ensureUsersInCacheByID(ctx, missingUserIDs); err != nil {
		log.Err(err).Msg("Failed to fetch missing XChat users")
	} else {
		tc.userCacheLock.RLock()
		for _, id := range missingUserIDs {
			if users[id] != nil {
				continue
			}
			if u := tc.userCache[id]; u != nil {
				users[id] = u
			}
		}
		tc.userCacheLock.RUnlock()
	}

	tc.userCacheLock.Lock()
	for userID, user := range users {
		tc.userCache[userID] = user
	}
	tc.userCacheLock.Unlock()

	// First, process all conversation key events (needed for decryption)
	for pIdx := range pages {
		for i := range pages[pIdx].Items {
			item := &pages[pIdx].Items[i]
			conversationID := item.ConversationDetail.ConversationID

			if err := processor.ProcessKeyChangeEvents(ctx, item); err != nil {
				log.Warn().
					Err(err).
					Str("conversation_id", conversationID).
					Msg("Failed to process key change events")
			}
		}
	}

	// Start XChat websocket for real-time events after keys are in place
	if err := tc.client.StartXChatWebsocket(ctx); err != nil {
		log.Err(err).Msg("Failed to start XChat websocket")
	}

	// Then process message/read events and sync channels
	for pIdx := range pages {
		for i := range pages[pIdx].Items {
			item := &pages[pIdx].Items[i]

			tc.syncXChatChannel(ctx, item, users)

			if err := processor.ProcessMessageAndReadEvents(ctx, item); err != nil {
				log.Warn().
					Err(err).
					Str("conversation_id", item.ConversationDetail.ConversationID).
					Msg("Failed to process message/read events")
			}
		}

		totalItems += len(pages[pIdx].Items)
	}

	log.Info().
		Int("conversations", totalItems).
		Int("users", len(users)).
		Msg("Finished fetching XChat inbox")

	tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})

	// Update remote profile from cached user data
	currentUserID := tc.client.GetCurrentUserID()
	if MakeUserLoginID(currentUserID) != tc.userLogin.ID {
		log.Warn().
			Str("user_login_id", string(tc.userLogin.ID)).
			Str("current_user_id", currentUserID).
			Msg("User login ID mismatch")
	}

	if selfUser, ok := users[currentUserID]; ok && selfUser != nil {
		remoteProfile := tc.makeXChatRemoteProfile(selfUser)
		if tc.userLogin.RemoteName != remoteProfile.Username ||
			tc.userLogin.RemoteProfile != *remoteProfile {
			tc.userLogin.RemoteName = remoteProfile.Username
			tc.userLogin.RemoteProfile = *remoteProfile
			if err := tc.userLogin.Save(ctx); err != nil {
				log.Err(err).Msg("Failed to save user login after updating remote profile")
			}
		}
	}

	// Save max sequence ID if updated
	if maxSeqID != "" && parseSequenceID(maxSeqID) > parseSequenceID(meta.MaxUserSequenceID) {
		log.Debug().
			Str("old_max_seq", meta.MaxUserSequenceID).
			Str("new_max_seq", maxSeqID).
			Msg("Updating max sequence ID")
		meta.MaxUserSequenceID = maxSeqID
		meta.MaxSequenceID = ""
	}

	// Persist message pull version if received
	if msgPullVersion != nil {
		meta.MessagePullVersion = msgPullVersion
	}

	// Save session state
	tc.HandleCursorChange(ctx)
}

// makeXChatRemoteProfile creates a RemoteProfile from XChat user data.
func (tc *TwitterClient) makeXChatRemoteProfile(user *types.User) *status.RemoteProfile {
	var avatarMXC id.ContentURIString
	ownGhost, err := tc.connector.br.GetGhostByID(context.Background(), MakeUserID(user.IDStr))
	if err == nil && ownGhost != nil {
		avatarMXC = ownGhost.AvatarMXC
	}
	return &status.RemoteProfile{
		Username: user.ScreenName,
		Name:     user.Name,
		Avatar:   avatarMXC,
	}
}

// parseSequenceID parses a sequence ID string to int64.
func parseSequenceID(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func (tc *TwitterClient) DoConnect(ctx context.Context) {
	tc.Connect(ctx)
}

func (tc *TwitterClient) HandleCursorChange(ctx context.Context) {
	if !tc.connector.Config.CacheSession {
		return
	}
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)
	meta.Session = tc.client.GetSession()
	meta.Session.LastSaved = time.Now()
	err := tc.userLogin.Save(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save user login after cursor change")
	}
}

func (tc *TwitterClient) startPolling(ctx context.Context, cached bool) {
	if ctx.Err() != nil {
		return
	}
	zerolog.Ctx(ctx).Info().Msg("Starting polling")
	tc.client.Connect(tc.connector.br.BackgroundCtx, cached)
}

func (tc *TwitterClient) Disconnect() {
	tc.client.Disconnect()
}

func (tc *TwitterClient) IsLoggedIn() bool {
	return tc.client.IsLoggedIn()
}

func (tc *TwitterClient) LogoutRemote(ctx context.Context) {
	log := zerolog.Ctx(ctx)
	err := tc.client.Logout(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to log out")
	}
}

func (tc *TwitterClient) IsThisUser(_ context.Context, userID networkid.UserID) bool {
	return UserLoginIDToUserID(tc.userLogin.ID) == userID
}

func (tc *TwitterClient) FullReconnect() {
	tc.Disconnect()
	tc.userLogin.Metadata.(*UserLoginMetadata).Session = nil
	tc.Connect(tc.userLogin.Log.WithContext(tc.connector.br.BackgroundCtx))
}

// TEMP TEST FUNCTION - DELETE AFTER TESTING
func (tc *TwitterClient) tempTestSendEncrypted(ctx context.Context) {
	// Wait 5 seconds before sending
	time.Sleep(5 * time.Second)

	log := zerolog.Ctx(ctx)

	// First, test signature with exact Python test data
	testSignature(log)

	// Test data with correct timestamp-based key versions
	conversationID := "1374864718591098886:1374864718591098886"
	convKeyB64 := "ST4hdGP1V3A5az3dM4GgjvoPRiaBl9TZUwGkfPy9/0E="
	privKeyB64 := "rdGJfa2OxTXLoUKp5D2O3fyjjiEFHrpCtLHO2TP4dbE="
	keyVersion := "1763143771470"    // Conversation key timestamp
	sigKeyVersion := "1749110687780" // Signing key timestamp

	// Decode conversation key
	convKey, err := base64.StdEncoding.DecodeString(convKeyB64)
	if err != nil {
		log.Error().Err(err).Msg("TEMP TEST: Failed to decode conversation key")
		return
	}

	// Parse signing key
	signingKey, err := crypto.ParsePrivateKeyScalar(privKeyB64)
	if err != nil {
		log.Error().Err(err).Msg("TEMP TEST: Failed to parse signing key")
		return
	}

	// Generate message ID
	messageID := uuid.NewString()

	// Build message using MessageBuilder with the provided key and signing key
	ownUserID := string(tc.userLogin.ID)
	builder := crypto.NewMessageBuilder(nil, ownUserID).
		SetMessageID(messageID).
		SetConversationID(conversationID).
		SetText("Hello from forged encrypted message!").
		WithConversationKey(convKey, keyVersion).
		WithSigningKey(signingKey, sigKeyVersion)

	// Build the message event
	event, err := builder.Build(ctx)
	if err != nil {
		log.Error().Err(err).Msg("TEMP TEST: Failed to build message")
		return
	}

	// Encode MessageCreateEvent to base64 Thrift
	encodedMCE, err := crypto.EncodeMessageCreateEvent(event.Detail.MessageCreateEvent)
	if err != nil {
		log.Error().Err(err).Msg("TEMP TEST: Failed to encode MCE")
		return
	}

	// Debug: decode and print raw bytes
	rawMCE, _ := base64.StdEncoding.DecodeString(encodedMCE)
	log.Info().
		Str("hex", fmt.Sprintf("%x", rawMCE)).
		Int("len", len(rawMCE)).
		Msg("TEMP TEST: Encoded MCE bytes")

	// Encode signature
	encodedSig := ""
	if event.MessageEventSignature != nil {
		encodedSig, err = crypto.EncodeMessageEventSignature(event.MessageEventSignature)
		if err != nil {
			log.Error().Err(err).Msg("TEMP TEST: Failed to encode signature")
			return
		}
	}

	// Generate conversation token (HS256 JWT)
	currentUserID := string(tc.userLogin.ID)
	recipient := crypto.GetRecipientFromConversationID(conversationID, currentUserID)
	validSinceMSec := int64(1745884800000)
	token := crypto.GenerateConversationToken(currentUserID, recipient, validSinceMSec, convKey)

	log.Info().
		Str("message_id", messageID).
		Str("conversation_id", conversationID).
		Str("encoded_mce", encodedMCE).
		Str("encoded_sig", encodedSig).
		Str("token", token).
		Msg("TEMP TEST: Sending forged encrypted message")

	// Send the message
	resp, err := tc.client.SendRawEncryptedMessage(ctx, conversationID, messageID, token, encodedMCE, encodedSig)
	if err != nil {
		log.Error().Err(err).Msg("TEMP TEST: Failed to send encrypted message")
		return
	}

	log.Info().Any("response", resp).Msg("TEMP TEST: Successfully sent encrypted message")
}

// testSignature tests our signature implementation against the Python test data
func testSignature(log *zerolog.Logger) {
	log.Info().Msg("=== SIGNATURE TEST ===")

	// New valid payload data
	messageID := "cb88a1ac-606d-4c58-ad66-8b865c2650b4"
	senderID := "1374864718591098886"
	conversationID := "1374864718591098886:1374864718591098886"
	convKeyVersion := "1763143771470"
	contentsHex := "3a84581ad14fd439151e9b7bd011a8f55260e4b4c26135dc1fd90acebc2ca8aec2e1fcb3d063422cdb924f1427c6b64e09b2a01e82d5b4627ebaddde64d482f4aa5b2b8645a10120244998fc5c471877ded191"
	expectedSig := "U7ZhqITL83z1/cME18jSJLEnrM1s3+n8Cd6FzKUpQx2oAI2NVVIO9bjid4DNc/tnbtq67pWNNOePTX/IL9O7oQ"
	pubKeyB64 := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAER7OEz60En69/Dr1z3vodgWk2UHn3Kxxs/qoSyaii/i5njCDyOMW6G95O7+Hcbv4iFGz33IWIqmXwBBnz3OBwcw=="
	privKeyB64 := "rdGJfa2OxTXLoUKp5D2O3fyjjiEFHrpCtLHO2TP4dbE="
	convKeyB64 := "ST4hdGP1V3A5az3dM4GgjvoPRiaBl9TZUwGkfPy9/0E="

	// Decode contents from hex
	contentsBytes, err := hex.DecodeString(contentsHex)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode contents hex")
		return
	}

	// Decrypt the test ciphertext to see what valid MessageContents looks like
	convKey, _ := base64.StdEncoding.DecodeString(convKeyB64)
	decrypted, err := crypto.SecretboxDecrypt(contentsBytes, convKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt test ciphertext")
	} else {
		log.Info().
			Str("decrypted_hex", fmt.Sprintf("%x", decrypted)).
			Int("len", len(decrypted)).
			Msg("Decrypted valid MessageContents (raw Thrift binary)")

		// Decode Thrift bytes generically to see the structure
		decoder := thrifter.NewDecoder(bytes.NewReader(decrypted), nil)
		var decoded thriftergeneral.Struct
		if err := decoder.Decode(&decoded); err != nil {
			log.Error().Err(err).Msg("Failed to decode MessageContents thrift")
		} else {
			log.Info().Interface("decoded", decoded).Msg("Decoded MessageContents structure")
		}
	}

	// Use "test" as that's what the valid payload contains
	text := "test"
	log.Info().Str("encoding_text", text).Msg("Encoding MessageContents with same text")

	ourContents, err := crypto.EncodeMessageContents(&payload.MessageContents{
		MessageText: &text,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to encode our MessageContents")
	} else {
		log.Info().
			Str("our_hex", fmt.Sprintf("%x", ourContents)).
			Int("len", len(ourContents)).
			Msg("Our encoded MessageContents")

		// Compare with decrypted valid payload
		if decrypted != nil {
			log.Info().
				Str("valid_hex", fmt.Sprintf("%x", decrypted)).
				Str("our_hex", fmt.Sprintf("%x", ourContents)).
				Bool("match", fmt.Sprintf("%x", decrypted) == fmt.Sprintf("%x", ourContents)).
				Msg("Comparison: valid vs our encoding")
		}
	}

	// Build preimage (same as Python)
	contentsB64 := base64.RawStdEncoding.EncodeToString(contentsBytes)
	preimage := fmt.Sprintf("MessageCreateEvent,%s,%s,%s,%s,%s",
		messageID, senderID, conversationID, convKeyVersion, contentsB64)

	log.Info().
		Str("preimage", preimage).
		Int("preimage_len", len(preimage)).
		Msg("Built preimage")

	// Parse keys
	privKey, err := crypto.ParsePrivateKeyScalar(privKeyB64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse private key")
		return
	}

	pubKey, err := crypto.ParsePublicKeySPKI(pubKeyB64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse public key")
		return
	}

	// Sign
	ourSig, err := crypto.SignMessage(privKey, messageID, senderID, conversationID, convKeyVersion, contentsBytes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign")
		return
	}

	log.Info().
		Str("our_sig", ourSig).
		Str("expected_sig", expectedSig).
		Msg("Generated signature")

	// Verify our signature
	err = crypto.VerifyMessage(pubKey, messageID, senderID, conversationID, convKeyVersion, contentsBytes, ourSig)
	if err != nil {
		log.Error().Err(err).Msg("Our signature INVALID")
	} else {
		log.Info().Msg("Our signature VALID")
	}

	// Verify expected signature
	err = crypto.VerifyMessage(pubKey, messageID, senderID, conversationID, convKeyVersion, contentsBytes, expectedSig)
	if err != nil {
		log.Error().Err(err).Msg("Expected signature INVALID")
	} else {
		log.Info().Msg("Expected signature VALID")
	}

	log.Info().Msg("=== END SIGNATURE TEST ===")
}
