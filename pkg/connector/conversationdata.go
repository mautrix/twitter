package connector

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type ensurePortalContextKey struct{}

// NormalizeConversationID converts a conversation ID from dash format (REST API / portal ID)
// to colon format (XChat API). Group chat IDs (prefixed with "g") are returned unchanged.
func NormalizeConversationID(id string) string {
	if strings.HasPrefix(id, "g") {
		return id
	}
	return strings.ReplaceAll(id, "-", ":")
}

func conversationDataResultID(requestedID, returnedID string) string {
	if returnedID == "" ||
		(strings.HasPrefix(requestedID, "g") && strings.TrimPrefix(requestedID, "g") == returnedID) {
		return requestedID
	}
	return returnedID
}

// fetchConversationData retrieves conversation details via the conversation data endpoint
// and converts them into an inbox item plus a user cache map to feed into existing sync logic.
func (tc *TwitterClient) fetchConversationData(ctx context.Context, conversationID string) (*response.XChatInboxItem, map[string]*types.User, error) {
	xchatConvID := NormalizeConversationID(conversationID)
	vars := payload.NewInboxPageConversationDataQueryVariables(xchatConvID, true)
	resp, err := tc.client.GetConversationData(ctx, vars)
	if err != nil {
		return nil, nil, err
	}

	data := resp.Data.GetInboxPageConversationData.Data
	data.ConversationDetail.ConversationID = conversationDataResultID(
		xchatConvID,
		data.ConversationDetail.ConversationID,
	)
	users := make(map[string]*types.User)
	missingIDs := make([]string, 0)

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

	collect(data.ConversationDetail.ParticipantsResults)
	collect(data.ConversationDetail.GroupMembersResults)
	collect(data.ConversationDetail.GroupAdminsResults)

	if err := tc.ensureUsersInCacheByID(ctx, missingIDs); err != nil {
		// Profile metadata is useful for room naming, but it must not prevent a
		// newly received message from creating/syncing its portal. The member-list
		// builder retries the lookup and can update metadata on a later resync.
		zerolog.Ctx(ctx).Warn().
			Err(err).
			Int("missing_users", len(missingIDs)).
			Msg("Failed to prefetch users for XChat conversation data")
	}

	tc.userCacheLock.RLock()
	for _, id := range missingIDs {
		if u := tc.userCache[id]; hasXChatUserDisplayInfo(u) {
			users[id] = u
		}
	}
	tc.userCacheLock.RUnlock()

	return &data, users, nil
}

// ensurePortalForConversation makes sure a portal exists for the given conversation and that the
// required key version (if provided) is available. If either is missing, it will fetch conversation
// data and sync the channel to create the portal and store keys.
func (tc *TwitterClient) ensurePortalForConversation(ctx context.Context, conversationID, requiredKeyVersion string) (*bridgev2.Portal, error) {
	lock := tc.getEnsurePortalLock(conversationID)
	lock.Lock()
	defer lock.Unlock()
	return tc.ensurePortalForConversationLocked(ctx, conversationID, requiredKeyVersion)
}

func (tc *TwitterClient) ensurePortalForConversationLocked(ctx context.Context, conversationID, requiredKeyVersion string) (*bridgev2.Portal, error) {
	portalKey := tc.MakePortalKeyFromID(conversationID)
	log := zerolog.Ctx(ctx).With().
		Str("conversation_id", conversationID).
		Str("required_key_version", requiredKeyVersion).
		Logger()

	portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
	if err != nil {
		return nil, err
	} else if portal.MXID != "" {
		if requiredKeyVersion == "" || tc.hasConversationKey(ctx, conversationID, requiredKeyVersion) {
			log.Debug().Msg("Portal already exists and required key is present")
			return portal, nil
		}
		log.Info().Msg("Portal exists but required key missing; fetching conversation data")
	} else {
		log.Info().Msg("Portal missing; fetching conversation data")
	}

	item, users, err := tc.fetchConversationData(ctx, conversationID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch conversation data")
		return nil, err
	}

	processor := tc.client.GetXChatProcessor()

	// Process key change events first (needed for decryption)
	if err := processor.ProcessKeyChangeEvents(ctx, item); err != nil {
		log.Warn().Err(err).Msg("Failed to process key change events for fetched conversation data")
	}

	// Sync channel (creates portal if needed)
	if err := tc.syncXChatChannel(ctx, item, users); err != nil {
		log.Warn().Err(err).Msg("Failed to sync channel for fetched conversation data")
		return portal, err
	}

	// Process messages/read events to backfill and register any keys embedded there
	bootstrapCtx := context.WithValue(ctx, ensurePortalContextKey{}, true)
	if err := processor.ProcessMessageAndReadEvents(bootstrapCtx, item); err != nil {
		log.Warn().Err(err).Msg("Failed to process message/read events for fetched conversation data")
	}

	if requiredKeyVersion != "" && !tc.hasConversationKey(ctx, conversationID, requiredKeyVersion) {
		return portal, fmt.Errorf("required conversation key %s still missing after sync", requiredKeyVersion)
	}

	return portal, nil
}

func isStalePortalMembershipError(err error) bool {
	if err == nil || (!errors.Is(err, mautrix.MNotFound) && !errors.Is(err, mautrix.MForbidden)) {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "user is not in the room")
}

func (tc *TwitterClient) recreateStaleXChatPortal(ctx context.Context, conversationID, failedRoomID string) error {
	lock := tc.getEnsurePortalLock(conversationID)
	lock.Lock()
	defer lock.Unlock()

	portal, err := tc.connector.br.GetPortalByKey(ctx, tc.MakePortalKeyFromID(conversationID))
	if err != nil {
		return err
	}
	if portal.MXID != "" {
		if string(portal.MXID) != failedRoomID {
			return nil
		}
		if err = portal.Delete(ctx); err != nil {
			return fmt.Errorf("delete stale XChat portal: %w", err)
		}
	}
	_, err = tc.ensurePortalForConversationLocked(ctx, conversationID, "")
	return err
}

func (tc *TwitterClient) queueXChatRemoteEventWithPortalRepair(
	ctx context.Context,
	conversationID string,
	evt bridgev2.RemoteEvent,
) bool {
	failedRoomID := ""
	if ctx != nil {
		portal, portalErr := tc.connector.br.GetExistingPortalByKey(ctx, tc.MakePortalKeyFromID(conversationID))
		if portalErr == nil && portal != nil {
			failedRoomID = string(portal.MXID)
		}
	}
	result := tc.userLogin.QueueRemoteEvent(evt)
	if xchatRemoteEventHandled(result) {
		return true
	}
	if ctx == nil || failedRoomID == "" || ctx.Value(ensurePortalContextKey{}) != nil || !isStalePortalMembershipError(result.Error) {
		return false
	}

	log := zerolog.Ctx(ctx).With().Str("conversation_id", conversationID).Logger()
	log.Warn().Msg("Recreating stale XChat portal room after Matrix membership failure")
	if err := tc.recreateStaleXChatPortal(ctx, conversationID, failedRoomID); err != nil {
		log.Warn().Err(err).Msg("Failed to recreate stale XChat portal room")
		return false
	}

	result = tc.userLogin.QueueRemoteEvent(evt)
	if !xchatRemoteEventHandled(result) {
		log.Warn().Err(result.Error).Msg("XChat event failed after recreating stale portal room")
		return false
	}
	return true
}

func (tc *TwitterClient) getEnsurePortalLock(conversationID string) *sync.Mutex {
	lock, _ := tc.ensurePortalLocks.LoadOrStore(conversationID, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func (tc *TwitterClient) hasConversationKey(ctx context.Context, conversationID, keyVersion string) bool {
	if keyVersion == "" {
		return true
	}
	keyManager := tc.client.GetKeyManager()
	if keyManager == nil {
		return false
	}
	_, err := keyManager.GetConversationKey(ctx, conversationID, keyVersion)
	return err == nil
}
