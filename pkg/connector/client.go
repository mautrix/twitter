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

	ensurePortalLocks sync.Map

	pollingChatResyncLast sync.Map
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
	client.SetEventHandler(tc.HandlePollingEvent, tc.HandleStreamEvent, tc.HandleCursorChange)
	client.SetConversationDataCallback(tc.HandleConversationDataRefresh)
	// Ensure current user ID is available even if cookies omit twid
	client.SetCurrentUserID(ParseUserLoginID(login.ID))
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
	meta := login.Metadata.(*UserLoginMetadata)
	c := cookies.NewCookiesFromString(meta.Cookies)
	log := login.Log.With().Str("component", "twitter_client").Logger()
	client := twittermeow.NewClient(c, newUserLoginKeyStore(login, tc), log)
	client.SetCurrentUserID(ParseUserLoginID(login.ID))
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

	// Migration detection: user has valid cookies but is missing encryption keys.
	// This happens when upgrading from the main branch (non-encrypted) to xchat/juicebox branch.
	if meta.Cookies != "" && meta.SecretKey == "" && meta.SigningKey == "" {
		log.Info().
			Str("user_id", ParseUserLoginID(tc.userLogin.ID)).
			Msg("Migration detected: user has cookies but missing encryption keys, triggering passcode-only reauth")
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateBadCredentials,
			Error:      "twitter-migration-reauth",
			Message:    "Please re-authenticate to enable X Chat. You only need to enter your passcode.",
		})
		return
	}

	// If pending encrypted sync after migration, force full resync
	if meta.PendingEncryptedSync {
		log.Info().Msg("Post-migration: forcing full resync for encrypted rooms")
		meta.Session = nil          // Clear cached session
		meta.MaxUserSequenceID = "" // Reset sequence to fetch all messages
	}

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
					Info: map[string]any{
						"go_error": err.Error(),
					},
				})
			}
			return
		}
	}

	// Full XChat inbox sync (migration / fresh login) can take a while. Mark as connected once we
	// know credentials are valid, and let the initial sync finish in the background.
	if meta.PendingEncryptedSync || meta.MaxUserSequenceID == "" {
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateConnected,
			Reason:     "sync_in_progress",
			Info: map[string]any{
				"sync": "in_progress",
			},
		})
	}

	// Start REST API sync for untrusted conversations in parallel
	go func() {
		tc.syncUntrustedChannels(ctx)
		// Start REST API polling for untrusted conversation updates
		tc.client.StartPolling(ctx)
	}()

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
	msgPullVersion := meta.MessagePullVersion

	// Errgroup for processing pages in parallel as they're fetched
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	var totalItems atomic.Int32
	var missingUserIDs []string
	var missingUserIDsMu sync.Mutex

	// processPage spawns goroutines to process each item in a page immediately
	processPage := func(page response.XChatInboxPage) {
		var pageMissing []string

		for i := range page.Items {
			item := &page.Items[i]
			totalItems.Add(1)
			missing := tc.cacheUsersFromItem(item)
			if len(missing) > 0 {
				pageMissing = append(pageMissing, missing...)
			}
		}

		if len(pageMissing) > 0 {
			if err := tc.ensureUsersInCacheByID(gCtx, pageMissing); err != nil {
				log.Warn().
					Err(err).
					Int("missing_users", len(pageMissing)).
					Msg("Failed to prefetch missing users for inbox page")
			}
			missingUserIDsMu.Lock()
			missingUserIDs = append(missingUserIDs, pageMissing...)
			missingUserIDsMu.Unlock()
		}

		for i := range page.Items {
			item := &page.Items[i]

			g.Go(func() error {
				// Process key changes first (needed for decryption)
				if err := processor.ProcessKeyChangeEvents(gCtx, item); err != nil {
					log.Warn().
						Err(err).
						Str("conversation_id", item.ConversationDetail.ConversationID).
						Msg("Failed to process key change events")
				}

				// Sync channel and process messages
				tc.syncXChatChannel(gCtx, item, nil)

				if err := processor.ProcessMessageAndReadEvents(gCtx, item); err != nil {
					log.Warn().
						Err(err).
						Str("conversation_id", item.ConversationDetail.ConversationID).
						Msg("Failed to process message/read events")
				}
				return nil
			})
		}
	}

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
			Info: map[string]any{
				"go_error": err.Error(),
			},
		})
		return
	}

	page := initialResp.Data.GetInboxPage
	fetchLog.Info().
		Str("response_cursor_id", page.InboxCursor.CursorID).
		Str("response_graph_snapshot_id", page.InboxCursor.GraphSnapshotID).
		Bool("pull_finished", page.InboxCursor.PullFinished).
		Int("items", len(page.Items)).
		Msg("Received XChat inbox page")

	if page.MaxUserSequenceID != nil && parseSequenceID(*page.MaxUserSequenceID) > parseSequenceID(maxSeqID) {
		maxSeqID = *page.MaxUserSequenceID
	}
	if page.MessagePullVersion != nil {
		msgPullVersion = page.MessagePullVersion
	}

	// Process initial page immediately
	processPage(page)

	var cursor *payload.XChatCursor
	if !page.InboxCursor.PullFinished && page.InboxCursor.CursorID != "" && page.InboxCursor.GraphSnapshotID != "" {
		cursor = &payload.XChatCursor{
			CursorId:        page.InboxCursor.CursorID,
			GraphSnapshotId: page.InboxCursor.GraphSnapshotID,
		}
	}

	// Subsequent pages via GetInboxPageRequest - process each page as it's fetched
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
				Info: map[string]any{
					"go_error": err.Error(),
				},
			})
			break
		}

		page := resp.Data.GetInboxPage
		fetchLog.Info().
			Str("response_cursor_id", page.InboxCursor.CursorID).
			Str("response_graph_snapshot_id", page.InboxCursor.GraphSnapshotID).
			Bool("pull_finished", page.InboxCursor.PullFinished).
			Int("items", len(page.Items)).
			Msg("Received XChat inbox page")

		if page.MaxUserSequenceID != nil && parseSequenceID(*page.MaxUserSequenceID) > parseSequenceID(maxSeqID) {
			maxSeqID = *page.MaxUserSequenceID
		}
		if page.MessagePullVersion != nil {
			msgPullVersion = page.MessagePullVersion
		}

		// Process this page immediately while fetching continues
		processPage(page)

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

	// Wait for all page processing to complete
	_ = g.Wait()

	// Batch fetch any users that only had RestID without inline data
	if len(missingUserIDs) > 0 {
		log.Info().
			Int("count", len(missingUserIDs)).
			Msg("Fetching missing user info")
		if err := tc.ensureUsersInCacheByID(ctx, missingUserIDs); err != nil {
			log.Warn().Err(err).Msg("Failed to fetch some missing users")
		}
	}

	// Start XChat websocket for real-time events after initial sync
	if err := tc.client.StartXChatWebsocket(ctx); err != nil {
		log.Err(err).Msg("Failed to start XChat websocket")
	}

	log.Info().
		Int("conversations", int(totalItems.Load())).
		Msg("Finished fetching XChat inbox")

	tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})

	// Update remote profile from cached user data
	currentUserID := tc.currentUserID()
	if MakeUserLoginID(currentUserID) != tc.userLogin.ID {
		log.Warn().
			Str("user_login_id", ParseUserLoginID(tc.userLogin.ID)).
			Str("current_user_id", currentUserID).
			Msg("User login ID mismatch")
	}
	if err := tc.forceRefreshUserInCacheByID(ctx, currentUserID); err != nil {
		log.Warn().
			Err(err).
			Str("current_user_id", currentUserID).
			Msg("Failed to refresh current user profile")
	}

	tc.userCacheLock.RLock()
	selfUser := tc.userCache[currentUserID]
	tc.userCacheLock.RUnlock()
	if selfUser != nil {
		remoteProfile := tc.makeXChatRemoteProfile(ctx, selfUser)
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
	}

	// Persist message pull version if received
	if msgPullVersion != nil {
		meta.MessagePullVersion = msgPullVersion
	}

	// Clear pending encrypted sync flag after successful sync
	if meta.PendingEncryptedSync {
		meta.PendingEncryptedSync = false
		log.Info().Msg("Post-migration: encrypted room sync completed")
	}

	// Save session state
	tc.HandleCursorChange(ctx)
}

// makeXChatRemoteProfile creates a RemoteProfile from XChat user data.
func (tc *TwitterClient) makeXChatRemoteProfile(ctx context.Context, user *types.User) *status.RemoteProfile {
	avatarMXC := tc.syncOwnAvatarFromUser(ctx, user)
	return &status.RemoteProfile{
		Username: user.ScreenName,
		Name:     user.Name,
		Avatar:   avatarMXC,
	}
}

func (tc *TwitterClient) syncOwnAvatarFromUser(ctx context.Context, user *types.User) id.ContentURIString {
	if user == nil || user.IDStr == "" {
		return ""
	}
	if user.IDStr != tc.client.GetCurrentUserID() {
		return ""
	}
	ownGhost, err := tc.connector.br.GetGhostByID(ctx, MakeUserID(user.IDStr))
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get own ghost by ID for avatar sync")
		return ""
	}
	ownGhost.UpdateInfo(ctx, tc.connector.wrapUserInfo(tc.client, user))
	return ownGhost.AvatarMXC
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
	return MakeUserID(ParseUserLoginID(tc.userLogin.ID)) == userID
}

func (tc *TwitterClient) FullReconnect() {
	tc.Disconnect()
	tc.userLogin.Metadata.(*UserLoginMetadata).Session = nil
	tc.Connect(tc.userLogin.Log.WithContext(tc.connector.br.BackgroundCtx))
}

// collectAndCacheUserResults processes a slice of XChatUserResult, caching users with Result data
// and returning IDs of users that only have RestID without inline Result data.
// Must be called with userCacheLock held.
func (tc *TwitterClient) collectAndCacheUserResults(results []response.XChatUserResult) []string {
	var missingIDs []string
	for _, p := range results {
		if p.Result != nil {
			tc.userCache[p.RestID] = twittermeow.ConvertXChatUserToUser(p.Result)
		} else if p.RestID != "" {
			if _, ok := tc.userCache[p.RestID]; !ok {
				missingIDs = append(missingIDs, p.RestID)
			}
		}
	}
	return missingIDs
}

// cacheUsersFromItem extracts user info from an XChatInboxItem and caches them.
// Returns a list of user IDs that only have RestID without inline Result data.
func (tc *TwitterClient) cacheUsersFromItem(item *response.XChatInboxItem) []string {
	tc.userCacheLock.Lock()
	defer tc.userCacheLock.Unlock()

	missingIDs := tc.collectAndCacheUserResults(item.ConversationDetail.ParticipantsResults)
	missingIDs = append(missingIDs, tc.collectAndCacheUserResults(item.ConversationDetail.GroupMembersResults)...)
	missingIDs = append(missingIDs, tc.collectAndCacheUserResults(item.ConversationDetail.GroupAdminsResults)...)
	return missingIDs
}

// HandleConversationDataRefresh is called when conversation data is fetched on-demand.
// It syncs the room data (members, name, avatar, etc.) from the fetched conversation data.
func (tc *TwitterClient) HandleConversationDataRefresh(ctx context.Context, conversationID string, item *response.XChatInboxItem) {
	if item == nil {
		return
	}

	log := zerolog.Ctx(ctx).With().
		Str("conversation_id", conversationID).
		Logger()

	// Build users map from item and collect missing IDs
	users := make(map[string]*types.User)
	var missingIDs []string

	collect := func(results []response.XChatUserResult) {
		for _, r := range results {
			if r.Result != nil {
				users[r.RestID] = twittermeow.ConvertXChatUserToUser(r.Result)
			} else if r.RestID != "" {
				missingIDs = append(missingIDs, r.RestID)
			}
		}
	}
	collect(item.ConversationDetail.ParticipantsResults)
	collect(item.ConversationDetail.GroupMembersResults)
	collect(item.ConversationDetail.GroupAdminsResults)

	// Fallback for 1:1 DMs: if no participants in response, parse from conversation ID
	if len(item.ConversationDetail.ParticipantsResults) == 0 && !strings.HasPrefix(conversationID, "g") {
		parts := strings.Split(conversationID, ":")
		if len(parts) == 2 {
			for _, userID := range parts {
				if userID != "" && users[userID] == nil {
					missingIDs = append(missingIDs, userID)
				}
			}
			// Populate ParticipantsResults so syncXChatChannel can build member list
			for _, userID := range parts {
				if userID != "" {
					item.ConversationDetail.ParticipantsResults = append(
						item.ConversationDetail.ParticipantsResults,
						response.XChatUserResult{RestID: userID},
					)
				}
			}
			log.Debug().
				Strs("parsed_user_ids", parts).
				Msg("Parsed user IDs from conversation ID (no participants in response)")
		}
	}

	// Fetch missing users via API
	if err := tc.ensureUsersInCacheByID(ctx, missingIDs); err != nil {
		log.Warn().Err(err).Msg("Failed to fetch missing users for conversation data refresh")
	}

	// Pull missing users from cache
	tc.userCacheLock.RLock()
	for _, id := range missingIDs {
		if users[id] == nil {
			if u := tc.userCache[id]; u != nil {
				users[id] = u
			}
		}
	}
	tc.userCacheLock.RUnlock()

	log.Debug().
		Int("users_count", len(users)).
		Int("missing_fetched", len(missingIDs)).
		Msg("Syncing conversation data from refresh callback")

	tc.syncXChatChannel(ctx, item, users)
}
