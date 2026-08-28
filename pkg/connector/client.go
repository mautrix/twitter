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
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
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

	connectLock    sync.Mutex
	connectCancel  context.CancelFunc
	connectRunLock sync.Mutex

	xchatInboxSyncLock    sync.Mutex
	xchatGapCatchupStates sync.Map
	xchatRequestsRepaired bool
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
	client.SetXChatGapHandler(tc.catchupXChatConversationGap)
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
	if meta.BrowserHeaders != nil {
		client.SetBrowserHeaders(*meta.BrowserHeaders)
	}
	client.SetCurrentUserID(ParseUserLoginID(login.ID))
	login.Client = NewTwitterClient(login, tc, client)
	return nil
}

const (
	sessionMaxTimeSinceSave = 24 * time.Hour
	sessionMaxTimeSinceInit = 48 * time.Hour
)

func (tc *TwitterClient) Connect(ctx context.Context) {
	// Bridge startup waits for every NetworkAPI.Connect call, while inbox import may take minutes.
	tc.startConnect(ctx, tc.connect)
}

func (tc *TwitterClient) startConnect(ctx context.Context, connect func(context.Context)) {
	connectCtx, cancel := context.WithCancel(ctx)
	tc.connectLock.Lock()
	previousCancel := tc.connectCancel
	tc.connectCancel = cancel
	tc.connectLock.Unlock()
	if previousCancel != nil {
		previousCancel()
	}
	go func() {
		tc.connectRunLock.Lock()
		defer tc.connectRunLock.Unlock()
		if connectCtx.Err() == nil {
			connect(connectCtx)
		}
	}()
}

func (tc *TwitterClient) cancelConnect() {
	tc.connectLock.Lock()
	cancel := tc.connectCancel
	tc.connectCancel = nil
	tc.connectLock.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (tc *TwitterClient) connect(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}
	// A replacement Connect call cancels the previous socket context. Wait for
	// any in-flight socket handoff to release shared inbox checkpoint state.
	tc.xchatInboxSyncLock.Lock()
	if ctx.Err() != nil {
		tc.xchatInboxSyncLock.Unlock()
		return
	}
	tc.xchatInboxSyncLock.Unlock()
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

	// If pending encrypted sync after migration, force full resync. Once a
	// resumable inbox cursor exists, keep it so reconnects can continue where
	// the previous import stopped.
	if meta.PendingEncryptedSync && meta.XChatInboxCursor == nil {
		log.Info().Msg("Post-migration: forcing full resync for encrypted rooms")
		meta.Session = nil          // Clear cached session
		meta.MaxUserSequenceID = "" // Reset sequence to fetch all messages
	}
	fullInboxSyncInProgress := meta.PendingEncryptedSync || meta.MaxUserSequenceID == "" || meta.XChatInboxCursor != nil

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
			if ctx.Err() != nil {
				return
			}
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
	if fullInboxSyncInProgress {
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateConnected,
			Reason:     "sync_in_progress",
			Info: map[string]any{
				"sync": "in_progress",
			},
		})
	}

	if err := tc.reconcileExistingGroupPortalAliases(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to reconcile existing REST/XChat group portal aliases")
	}
	if err := tc.repairExistingXChatMessageRequests(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to repair existing XChat message-request state")
	}

	// Set up XChat processor and sequence ID tracking
	processor := tc.client.GetXChatProcessor()
	var maxSeqID string
	var maxSeqIDLock sync.Mutex

	setMaxSeqID := func(seqID string) {
		if seqID == "" {
			return
		}
		maxSeqIDLock.Lock()
		defer maxSeqIDLock.Unlock()
		if compareIntStrings(seqID, maxSeqID) > 0 {
			maxSeqID = seqID
		}
	}

	getMaxSeqID := func() string {
		maxSeqIDLock.Lock()
		defer maxSeqIDLock.Unlock()
		return maxSeqID
	}

	processor.ResetSequenceState()
	// Inbox pages contain multiple conversations and aren't globally ordered by
	// sequence ID. Keep live callback publication disabled until the page-level
	// checkpoint logic confirms that no conversation on the import is blocked.
	processor.SetSequenceIDCallback(nil)

	// Fetch and checkpoint the XChat inbox one page at a time. Transient
	// failures retry from the last persisted cursor instead of starting the
	// websocket with a known hole in the initial import.
	fetchLog := log.With().Str("component", "xchat_fetch").Logger()
	var totalItems atomic.Int32
	missingUserIDs := make(map[string]struct{})
	// A non-empty checkpoint means this is an incremental import. In that mode,
	// an inbox item with has_more must be repaired before its latest events can
	// safely advance the global sequence. Fresh full imports use backfill tasks
	// for older history instead of eagerly walking every conversation here.
	repairTruncatedItems := meta.MaxUserSequenceID != ""
	var catchupResult xchatInboxCatchupResult
	retryDelay := time.Second
	for {
		var err error
		catchupResult, err = runXChatInboxCatchup(ctx, xchatInboxCatchupState{
			MaxSequenceID:      meta.MaxUserSequenceID,
			MessagePullVersion: meta.MessagePullVersion,
			Cursor:             xchatInboxCursorFromMetadata(meta.XChatInboxCursor),
		}, xchatInboxCatchupOps{
			FetchInitial: func(ctx context.Context, variables *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
				fetchLog.Info().
					Str("max_sequence_id", variables.MaxLocalSequenceId).
					Msg("Fetching initial XChat inbox page")
				resp, err := tc.client.GetInitialXChatPage(ctx, variables)
				if err != nil {
					return response.XChatInboxPage{}, err
				}
				return resp.Data.GetInboxPage, nil
			},
			FetchNext: func(ctx context.Context, variables *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
				fetchLog.Info().
					Str("cursor_id", variables.ContinueCursor.CursorId).
					Str("graph_snapshot_id", variables.ContinueCursor.GraphSnapshotId).
					Msg("Fetching XChat inbox page")
				resp, err := tc.client.GetInboxPageRequest(ctx, variables)
				if err != nil {
					return response.XChatInboxPage{}, err
				}
				return resp.Data.GetInboxPage, nil
			},
			ProcessPage: func(ctx context.Context, page response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
				pageMissing, err := tc.processXChatInboxPage(ctx, page, &totalItems, repairTruncatedItems)
				for _, userID := range pageMissing {
					missingUserIDs[userID] = struct{}{}
				}
				return xchatInboxPageProcessResult{
					MaxSequenceID:     processor.MaxHandledSequenceID(),
					CheckpointBlocked: processor.SequenceCheckpointBlocked(),
				}, err
			},
			Checkpoint: tc.saveXChatInboxCheckpoint,
		})
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			return
		}
		fetchLog.Err(err).
			Dur("retry_in", retryDelay).
			Msg("XChat inbox sync failed, retrying from saved checkpoint")
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateUnknownError,
			Error:      "twitter-xchat-sync-error",
			Info: map[string]any{
				"go_error": err.Error(),
			},
		})
		if !waitForXChatInboxRetry(ctx, retryDelay) {
			return
		}
		retryDelay = min(retryDelay*2, time.Minute)
	}

	setMaxSeqID(catchupResult.MaxSequenceID)
	processor.SetSequenceIDCallback(setMaxSeqID)
	msgPullVersion := cloneXChatInt(catchupResult.MessagePullVersion)
	if catchupResult.CheckpointBlocked {
		fetchLog.Warn().Msg("XChat inbox import finished with an unresolved conversation gap; checkpoint remains pending")
	}

	// Batch fetch any users that only had RestID without inline data
	if len(missingUserIDs) > 0 {
		missingUserIDList := make([]string, 0, len(missingUserIDs))
		for userID := range missingUserIDs {
			missingUserIDList = append(missingUserIDList, userID)
		}
		log.Info().
			Int("count", len(missingUserIDs)).
			Msg("Fetching missing user info")
		if err := tc.ensureUsersInCacheByID(ctx, missingUserIDList); err != nil {
			log.Warn().Err(err).Msg("Failed to fetch some missing users")
		}
	}
	if ctx.Err() != nil {
		return
	}

	getMessagePullVersion := func() *int {
		if msgPullVersion == nil {
			return nil
		}
		value := *msgPullVersion
		return &value
	}
	setMessagePullVersion := func(value *int) {
		if value == nil {
			return
		}
		copy := *value
		msgPullVersion = &copy
	}
	tc.client.SetXChatConnectHandler(func(connectCtx context.Context) error {
		return tc.syncXChatInboxAfterConnect(
			connectCtx,
			getMaxSeqID,
			setMaxSeqID,
			getMessagePullVersion,
			setMessagePullVersion,
		)
	})

	// Start XChat websocket for real-time events after initial sync
	if err := tc.client.StartXChatWebsocket(ctx); err != nil {
		if ctx.Err() != nil {
			return
		}
		log.Err(err).Msg("Failed to start XChat websocket")
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateUnknownError,
			Error:      "twitter-xchat-connect-error",
			Info: map[string]any{
				"go_error": err.Error(),
			},
		})
		return
	}
	log.Info().
		Int("conversations", int(totalItems.Load())).
		Msg("Finished fetching XChat inbox")

	go func() {
		tc.syncUntrustedChannels(ctx)
		tc.client.StartPolling(ctx)
	}()

	if ctx.Err() != nil {
		return
	}
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
	var remoteProfile *status.RemoteProfile
	if selfUser != nil {
		remoteProfile = tc.makeXChatRemoteProfile(ctx, selfUser)
	}

	// Serialize final user-login changes with reconnect catch-up. Page
	// checkpoints have already persisted max sequence, cursor and pull version;
	// rewriting them here could clobber a newer reconnect checkpoint.
	tc.xchatInboxSyncLock.Lock()
	defer tc.xchatInboxSyncLock.Unlock()
	if remoteProfile != nil {
		if tc.userLogin.RemoteName != remoteProfile.Username ||
			tc.userLogin.RemoteProfile != *remoteProfile {
			tc.userLogin.RemoteName = remoteProfile.Username
			tc.userLogin.RemoteProfile = *remoteProfile
		}
	}
	inboxSyncComplete := !processor.SequenceCheckpointBlocked()
	if meta.PendingEncryptedSync && inboxSyncComplete {
		meta.PendingEncryptedSync = false
		log.Info().Msg("Post-migration: encrypted room sync completed")
	} else if meta.PendingEncryptedSync {
		log.Warn().Msg("Post-migration encrypted room sync paused; will resume from saved inbox cursor")
	}

	// Save session state
	tc.HandleCursorChange(ctx)
}

func waitForXChatInboxRetry(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
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

func xchatInboxCursorFromMetadata(cursor *XChatInboxCursorData) *payload.XChatCursor {
	if cursor == nil || cursor.CursorID == "" || cursor.GraphSnapshotID == "" {
		return nil
	}
	return &payload.XChatCursor{
		CursorId:        cursor.CursorID,
		GraphSnapshotId: cursor.GraphSnapshotID,
	}
}

func nextXChatInboxCursor(page response.XChatInboxPage) *payload.XChatCursor {
	if page.InboxCursor.PullFinished || page.InboxCursor.CursorID == "" || page.InboxCursor.GraphSnapshotID == "" {
		return nil
	}
	return &payload.XChatCursor{
		CursorId:        page.InboxCursor.CursorID,
		GraphSnapshotId: page.InboxCursor.GraphSnapshotID,
	}
}

func validatedNextXChatInboxCursor(page response.XChatInboxPage) (*payload.XChatCursor, error) {
	// Terminal pages may omit the continuation cursor without setting pull_finished,
	// while still echoing the graph snapshot ID.
	if page.InboxCursor.PullFinished || page.InboxCursor.CursorID == "" {
		return nil, nil
	}
	cursor := nextXChatInboxCursor(page)
	if cursor == nil {
		return nil, errors.New("XChat inbox page is unfinished but has an incomplete cursor")
	}
	return cursor, nil
}

func (tc *TwitterClient) saveXChatInboxCheckpoint(ctx context.Context, cursor *payload.XChatCursor, maxSeqID string, msgPullVersion *int) error {
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)
	previousCursor := meta.XChatInboxCursor
	previousMaxSequenceID := meta.MaxUserSequenceID
	previousMessagePullVersion := meta.MessagePullVersion
	if cursor == nil {
		meta.XChatInboxCursor = nil
	} else {
		meta.XChatInboxCursor = &XChatInboxCursorData{
			CursorID:        cursor.CursorId,
			GraphSnapshotID: cursor.GraphSnapshotId,
		}
	}
	if maxSeqID != "" && compareIntStrings(maxSeqID, meta.MaxUserSequenceID) > 0 {
		meta.MaxUserSequenceID = maxSeqID
	}
	if msgPullVersion != nil {
		meta.MessagePullVersion = cloneXChatInt(msgPullVersion)
	}
	if err := tc.saveUserLoginState(ctx); err != nil {
		meta.XChatInboxCursor = previousCursor
		meta.MaxUserSequenceID = previousMaxSequenceID
		meta.MessagePullVersion = previousMessagePullVersion
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save XChat inbox import checkpoint")
		return err
	}

	evt := zerolog.Ctx(ctx).Info().
		Bool("complete", cursor == nil).
		Str("max_user_sequence_id", meta.MaxUserSequenceID)
	if msgPullVersion != nil {
		evt.Int("message_pull_version", *msgPullVersion)
	}
	if cursor == nil {
		evt.Msg("Cleared XChat inbox import checkpoint")
	} else {
		evt.
			Str("cursor_id", cursor.CursorId).
			Str("graph_snapshot_id", cursor.GraphSnapshotId).
			Msg("Saved XChat inbox import checkpoint")
	}
	return nil
}

func (tc *TwitterClient) DoConnect(ctx context.Context) {
	tc.Connect(ctx)
}

func (tc *TwitterClient) HandleCursorChange(ctx context.Context) {
	err := tc.saveUserLoginState(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save user login after cursor change")
	}
}

func (tc *TwitterClient) saveUserLoginState(ctx context.Context) error {
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)
	if tc.connector.Config.CacheSession {
		meta.Session = tc.client.GetSession()
		if meta.Session != nil {
			meta.Session.LastSaved = time.Now()
		}
	}
	return tc.userLogin.Save(ctx)
}

func (tc *TwitterClient) Disconnect() {
	tc.cancelConnect()
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
	tc.xchatInboxSyncLock.Lock()
	tc.userLogin.Metadata.(*UserLoginMetadata).Session = nil
	tc.xchatInboxSyncLock.Unlock()
	tc.Connect(tc.userLogin.Log.WithContext(tc.connector.br.BackgroundCtx))
}

// Must be called with userCacheLock held.
func (tc *TwitterClient) collectAndCacheUserResults(results []response.XChatUserResult) []string {
	var missingIDs []string
	for _, p := range results {
		userID, user := xchatUserFromResult(p)
		if userID == "" {
			continue
		}
		cached := mergeXChatUser(tc.userCache, userID, user)
		if !hasCompleteXChatUserProfile(cached) {
			missingIDs = append(missingIDs, userID)
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
	if tc.connector.br.BackgroundCtx != nil {
		ctx = tc.userLogin.Log.WithContext(tc.connector.br.BackgroundCtx)
	}
	if ctx.Err() != nil {
		return
	}

	log := zerolog.Ctx(ctx).With().
		Str("conversation_id", conversationID).
		Logger()
	returnedConversationID := item.ConversationDetail.ConversationID
	item.ConversationDetail.ConversationID = conversationDataResultID(
		conversationID,
		returnedConversationID,
	)
	if item.ConversationDetail.ConversationID != returnedConversationID {
		log.Debug().
			Str("returned_conversation_id", returnedConversationID).
			Str("resolved_conversation_id", item.ConversationDetail.ConversationID).
			Msg("Normalized conversation ID from conversation data refresh")
	}

	// Build users map from item and collect missing IDs
	users := make(map[string]*types.User)
	var missingIDs []string

	collect := func(results []response.XChatUserResult) {
		for _, r := range results {
			userID, user := xchatUserFromResult(r)
			if userID == "" {
				continue
			}
			mergeXChatUser(users, userID, user)
			if !hasCompleteXChatUserProfile(user) {
				missingIDs = append(missingIDs, userID)
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
		if u := tc.userCache[id]; hasXChatUserDisplayInfo(u) {
			users[id] = u
		}
	}
	tc.userCacheLock.RUnlock()

	log.Debug().
		Int("users_count", len(users)).
		Int("missing_fetched", len(missingIDs)).
		Msg("Syncing conversation data from refresh callback")

	if err := tc.syncXChatChannel(ctx, item, users); err != nil {
		log.Warn().Err(err).Msg("Failed to sync conversation data from refresh callback")
	}
}
