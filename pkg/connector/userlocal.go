package connector

import (
	"context"
	"slices"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/simplevent"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) portalHasRoom(ctx context.Context, conversationID string) bool {
	portal, err := tc.connector.br.GetExistingPortalByKey(ctx, tc.MakePortalKeyFromID(conversationID))
	return err == nil && portal != nil && portal.MXID != ""
}

func (tc *TwitterClient) loginMetadata() *UserLoginMetadata {
	if tc.userLogin == nil {
		return nil
	}
	meta, _ := tc.userLogin.Metadata.(*UserLoginMetadata)
	return meta
}

// queueUserLocalPortalInfo applies per-user chat state from X to portals that already have rooms.
// bridgev2 only applies it to existing rooms when mute_only_on_create/tag_only_on_create are off.
func (tc *TwitterClient) queueUserLocalPortalInfo(ctx context.Context, conversationIDs []string, ts time.Time, info func(conversationID string) *bridgev2.UserLocalPortalInfo) bool {
	success := true
	for _, conversationID := range conversationIDs {
		if !tc.portalHasRoom(ctx, conversationID) {
			continue
		}
		success = tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventChatInfoChange,
				PortalKey: tc.MakePortalKeyFromID(conversationID),
				Timestamp: ts,
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{UserLocal: info(conversationID)},
			},
		}).Success && success
	}
	return success
}

func staticUserLocalInfo(info *bridgev2.UserLocalPortalInfo) func(string) *bridgev2.UserLocalPortalInfo {
	return func(string) *bridgev2.UserLocalPortalInfo { return info }
}

func (tc *TwitterClient) isConversationPinned(conversationID string) bool {
	meta := tc.loginMetadata()
	if meta == nil {
		return false
	}
	tc.userLocalLock.Lock()
	defer tc.userLocalLock.Unlock()
	return slices.Contains(meta.PinnedConversationIDs, NormalizeConversationID(conversationID))
}

// applyPinAtCreation tags a new room from the pinned list. ExtraUpdates runs in the portal event loop right
// before bridgev2 applies UserLocal to a new room, so the list read here includes snapshots handled while
// the room was being created.
func (tc *TwitterClient) applyPinAtCreation(ctx context.Context, conversationID string, info *bridgev2.ChatInfo) {
	if tc.loginMetadata() == nil || tc.portalHasRoom(ctx, conversationID) {
		return
	}
	if info.UserLocal == nil {
		info.UserLocal = &bridgev2.UserLocalPortalInfo{}
	}
	userLocal := info.UserLocal
	info.ExtraUpdates = bridgev2.MergeExtraUpdaters(info.ExtraUpdates, func(_ context.Context, portal *bridgev2.Portal) bool {
		if portal.MXID == "" && tc.isConversationPinned(conversationID) {
			userLocal.Tag = ptr.Ptr(tc.connector.Config.PinnedTag)
		}
		return false
	})
}

// handleXChatPinnedConversations diffs a full pinned list from the user metadata against the last applied one.
func (tc *TwitterClient) handleXChatPinnedConversations(ctx context.Context, evt *types.XChatPinnedConversations) bool {
	meta := tc.loginMetadata()
	if meta == nil {
		return true
	}
	tc.userLocalLock.Lock()
	if evt.SequenceID == "" || compareIntStrings(evt.SequenceID, meta.PinnedConversationsSequenceID) <= 0 {
		tc.userLocalLock.Unlock()
		return true
	}
	previous := slices.Clone(meta.PinnedConversationIDs)
	tc.userLocalLock.Unlock()
	pinned := make([]string, 0, len(evt.ConversationIDs))
	for _, conversationID := range evt.ConversationIDs {
		pinned = append(pinned, NormalizeConversationID(conversationID))
	}
	var added, removed []string
	for _, conversationID := range pinned {
		if !slices.Contains(previous, conversationID) {
			added = append(added, conversationID)
		}
	}
	for _, conversationID := range previous {
		if !slices.Contains(pinned, conversationID) {
			removed = append(removed, conversationID)
		}
	}
	addOK := tc.queueUserLocalPortalInfo(ctx, added, evt.Timestamp, staticUserLocalInfo(&bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(tc.connector.Config.PinnedTag)}))
	removeOK := tc.queueUserLocalPortalInfo(ctx, removed, evt.Timestamp, staticUserLocalInfo(&bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTag(""))}))
	if !addOK || !removeOK {
		return false
	}
	tc.userLocalLock.Lock()
	defer tc.userLocalLock.Unlock()
	if compareIntStrings(evt.SequenceID, meta.PinnedConversationsSequenceID) > 0 {
		meta.PinnedConversationIDs = pinned
		meta.PinnedConversationsSequenceID = evt.SequenceID
		if err := tc.userLogin.Save(ctx); err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to save pinned conversations")
		}
	}
	return true
}

func (tc *TwitterClient) handleXChatArchiveChange(ctx context.Context, evt *types.XChatConversationArchiveChange) bool {
	if tc.connector.Config.ArchiveTag == "" {
		return true
	}
	var tag event.RoomTag
	if evt.Archived {
		tag = tc.connector.Config.ArchiveTag
	}
	if evt.ConversationID == "" || tc.isConversationPinned(evt.ConversationID) {
		return true
	}
	conversationIDs := []string{NormalizeConversationID(evt.ConversationID)}
	return tc.queueUserLocalPortalInfo(ctx, conversationIDs, evt.Timestamp, staticUserLocalInfo(&bridgev2.UserLocalPortalInfo{Tag: &tag}))
}

func (tc *TwitterClient) xchatItemArchived(item *response.XChatInboxItem) bool {
	if item.ConversationDetail.IsArchived != nil {
		return *item.ConversationDetail.IsArchived
	}
	var archived bool
	var latestSequenceID string
	for _, raw := range slices.Concat(item.LatestMessageEvents, item.EncodedMessageEvents) {
		evt, err := twittermeow.DecodeMessageEvent(raw)
		if err != nil || evt.Detail == nil || evt.Detail.ConversationMetadataChangeEvent == nil {
			continue
		}
		change := evt.Detail.ConversationMetadataChangeEvent.ConversationMetadataChange
		if change == nil || (change.ArchiveConversation == nil && change.UnarchiveConversation == nil) ||
			ptr.Val(evt.SenderId) != tc.client.GetCurrentUserID() || compareIntStrings(ptr.Val(evt.SequenceId), latestSequenceID) <= 0 {
			continue
		}
		latestSequenceID = ptr.Val(evt.SequenceId)
		archived = change.ArchiveConversation != nil
	}
	return archived
}
