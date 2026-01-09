package connector

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type ensurePortalContextKey struct{}

// fetchConversationData retrieves conversation details via the conversation data endpoint
// and converts them into an inbox item plus a user cache map to feed into existing sync logic.
func (tc *TwitterClient) fetchConversationData(ctx context.Context, conversationID string) (*response.XChatInboxItem, map[string]*types.User, error) {
	vars := payload.NewInboxPageConversationDataQueryVariables(conversationID, true)
	resp, err := tc.client.GetConversationData(ctx, vars)
	if err != nil {
		return nil, nil, err
	}

	data := resp.Data.GetInboxPageConversationData.Data
	users := make(map[string]*types.User)
	missingIDs := make([]string, 0)

	collect := func(results []response.XChatUserResult) {
		for _, r := range results {
			if r.Result != nil {
				users[r.RestID] = twittermeow.ConvertXChatUserToUser(r.Result)
			} else if r.RestID != "" {
				missingIDs = append(missingIDs, r.RestID)
			}
		}
	}

	collect(data.ConversationDetail.ParticipantsResults)
	collect(data.ConversationDetail.GroupMembersResults)
	collect(data.ConversationDetail.GroupAdminsResults)

	if err := tc.ensureUsersInCacheByID(ctx, missingIDs); err != nil {
		return nil, nil, err
	}

	tc.userCacheLock.RLock()
	for _, id := range missingIDs {
		if users[id] != nil {
			continue
		}
		if u := tc.userCache[id]; u != nil {
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

	portalKey := tc.MakePortalKeyFromID(conversationID)
	log := zerolog.Ctx(ctx).With().
		Str("conversation_id", conversationID).
		Str("required_key_version", requiredKeyVersion).
		Logger()

	portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
	if err == nil && portal != nil && portal.MXID != "" {
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
	tc.syncXChatChannel(ctx, item, users)

	// Process messages/read events to backfill and register any keys embedded there
	bootstrapCtx := context.WithValue(ctx, ensurePortalContextKey{}, true)
	if err := processor.ProcessMessageAndReadEvents(bootstrapCtx, item); err != nil {
		log.Warn().Err(err).Msg("Failed to process message/read events for fetched conversation data")
	}

	portal, err = tc.connector.br.GetPortalByKey(ctx, portalKey)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to find portal after fetching conversation data")
		return nil, err
	}

	if portal == nil || portal.MXID == "" {
		return nil, fmt.Errorf("portal not found for conversation %s after sync", conversationID)
	}

	if requiredKeyVersion != "" && !tc.hasConversationKey(ctx, conversationID, requiredKeyVersion) {
		return portal, fmt.Errorf("required conversation key %s still missing after sync", requiredKeyVersion)
	}

	return portal, nil
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
