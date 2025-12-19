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
	"context"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/status"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
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

	// Flatten all items for parallel processing
	var allItems []*response.XChatInboxItem
	for pIdx := range pages {
		for i := range pages[pIdx].Items {
			allItems = append(allItems, &pages[pIdx].Items[i])
		}
	}
	totalItems = len(allItems)

	// Process key change events in parallel (bounded concurrency)
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	for _, item := range allItems {
		g.Go(func() error {
			if err := processor.ProcessKeyChangeEvents(gCtx, item); err != nil {
				log.Warn().
					Err(err).
					Str("conversation_id", item.ConversationDetail.ConversationID).
					Msg("Failed to process key change events")
			}
			return nil
		})
	}
	_ = g.Wait()

	// Start XChat websocket for real-time events after keys are in place
	if err := tc.client.StartXChatWebsocket(ctx); err != nil {
		log.Err(err).Msg("Failed to start XChat websocket")
	}

	// Process message/read events and sync channels in parallel
	g2, gCtx2 := errgroup.WithContext(ctx)
	g2.SetLimit(10)

	for _, item := range allItems {
		g2.Go(func() error {
			tc.syncXChatChannel(gCtx2, item, users)

			if err := processor.ProcessMessageAndReadEvents(gCtx2, item); err != nil {
				log.Warn().
					Err(err).
					Str("conversation_id", item.ConversationDetail.ConversationID).
					Msg("Failed to process message/read events")
			}
			return nil
		})
	}
	_ = g2.Wait()

	log.Info().
		Int("conversations", totalItems).
		Int("users", len(users)).
		Msg("Finished fetching XChat inbox")

	tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})

	// Update remote profile from cached user data
	currentUserID := tc.client.GetCurrentUserID()
	if networkid.UserLoginID(currentUserID) != tc.userLogin.ID {
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
	ownGhost, err := tc.connector.br.GetGhostByID(context.Background(), networkid.UserID(user.IDStr))
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
	return networkid.UserID(tc.userLogin.ID) == userID
}

func (tc *TwitterClient) FullReconnect() {
	tc.Disconnect()
	tc.userLogin.Metadata.(*UserLoginMetadata).Session = nil
	tc.Connect(tc.userLogin.Log.WithContext(tc.connector.br.BackgroundCtx))
}
