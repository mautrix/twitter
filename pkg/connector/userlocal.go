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

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) portalHasRoom(ctx context.Context, conversationID string) bool {
	portal, err := tc.connector.br.GetExistingPortalByKey(ctx, tc.MakePortalKeyFromID(conversationID))
	return err == nil && portal != nil && portal.MXID != ""
}

// queueUserLocalPortalInfo applies per-user chat state from X to portals that already have rooms.
// bridgev2 only applies it to existing rooms when mute_only_on_create/tag_only_on_create are off.
func (tc *TwitterClient) queueUserLocalPortalInfo(ctx context.Context, conversationIDs []string, ts time.Time, info *bridgev2.UserLocalPortalInfo) bool {
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
				ChatInfo: &bridgev2.ChatInfo{UserLocal: info},
			},
		}).Success && success
	}
	return success
}

func (tc *TwitterClient) isConversationPinned(conversationID string) bool {
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)
	return slices.Contains(meta.PinnedConversationIDs, NormalizeConversationID(conversationID))
}

// handleXChatPinnedConversations diffs a full pinned list from the user metadata against the last applied one.
func (tc *TwitterClient) handleXChatPinnedConversations(ctx context.Context, evt *types.XChatPinnedConversations) bool {
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)
	if evt.SequenceID == "" || compareIntStrings(evt.SequenceID, meta.PinnedConversationsSequenceID) <= 0 {
		return true
	}
	pinned := make([]string, 0, len(evt.ConversationIDs))
	for _, conversationID := range evt.ConversationIDs {
		pinned = append(pinned, NormalizeConversationID(conversationID))
	}
	previous := meta.PinnedConversationIDs
	meta.PinnedConversationIDs = pinned
	meta.PinnedConversationsSequenceID = evt.SequenceID
	if err := tc.userLogin.Save(ctx); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save pinned conversations")
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
	addOK := tc.queueUserLocalPortalInfo(ctx, added, evt.Timestamp, &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTagFavourite)})
	removeOK := tc.queueUserLocalPortalInfo(ctx, removed, evt.Timestamp, &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTag(""))})
	return addOK && removeOK
}
