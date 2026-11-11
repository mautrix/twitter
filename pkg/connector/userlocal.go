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
			userLocal.Tag = ptr.Ptr(event.RoomTagFavourite)
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
	archived := slices.Clone(meta.ArchivedConversationIDs)
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
	addOK := tc.queueUserLocalPortalInfo(ctx, added, evt.Timestamp, staticUserLocalInfo(&bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTagFavourite)}))
	removeOK := tc.queueUserLocalPortalInfo(ctx, removed, evt.Timestamp, func(conversationID string) *bridgev2.UserLocalPortalInfo {
		if slices.Contains(archived, conversationID) {
			return &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTagLowPriority)}
		}
		return &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTag(""))}
	})
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
	meta := tc.loginMetadata()
	if meta == nil {
		return true
	}
	tc.userLocalLock.Lock()
	pinned := slices.Clone(meta.PinnedConversationIDs)
	if meta.ArchiveSequenceIDs == nil {
		meta.ArchiveSequenceIDs = make(map[string]string)
	}
	var changed []string
	for _, conversationID := range evt.ConversationIDs {
		conversationID = NormalizeConversationID(conversationID)
		if evt.SequenceID == "" || compareIntStrings(evt.SequenceID, meta.ArchiveSequenceIDs[conversationID]) <= 0 {
			continue
		}
		meta.ArchiveSequenceIDs[conversationID] = evt.SequenceID
		changed = append(changed, conversationID)
		meta.ArchivedConversationIDs = slices.DeleteFunc(meta.ArchivedConversationIDs, func(id string) bool { return id == conversationID })
		if evt.Archived {
			meta.ArchivedConversationIDs = append(meta.ArchivedConversationIDs, conversationID)
		}
	}
	err := tc.userLogin.Save(ctx)
	tc.userLocalLock.Unlock()
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save archived conversations")
	}
	return tc.queueUserLocalPortalInfo(ctx, changed, evt.Timestamp, func(conversationID string) *bridgev2.UserLocalPortalInfo {
		switch {
		case evt.Archived:
			return &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTagLowPriority)}
		case slices.Contains(pinned, NormalizeConversationID(conversationID)):
			return &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTagFavourite)}
		default:
			return &bridgev2.UserLocalPortalInfo{Tag: ptr.Ptr(event.RoomTag(""))}
		}
	})
}

// handleXChatMarkedUnread bridges the user's own mark-unread from another device once, since inbox
// snapshots replay the latest events of each conversation on every reconnect.
func (tc *TwitterClient) handleXChatMarkedUnread(ctx context.Context, evt *types.XChatConversationMarkedUnread) bool {
	meta := tc.loginMetadata()
	if meta == nil {
		return true
	}
	key := NormalizeConversationID(evt.ConversationID)
	tc.userLocalLock.Lock()
	lastSequenceID, seen := meta.MarkedUnreadSequenceIDs[key]
	tc.userLocalLock.Unlock()
	if evt.SequenceID == "" || compareIntStrings(evt.SequenceID, lastSequenceID) <= 0 {
		return true
	}
	// The first mark-unread seen from a snapshot may be long stale, and rooms that don't exist yet only need
	// the sequence recorded so later replays are ignored.
	if (evt.FromSnapshot && !seen) || !tc.portalHasRoom(ctx, evt.ConversationID) {
		tc.recordMarkedUnreadSequence(ctx, meta, key, evt.SequenceID)
		return true
	}
	if !tc.userLogin.QueueRemoteEvent(&simplevent.MarkUnread{
		EventMeta: simplevent.EventMeta{
			Type:      bridgev2.RemoteEventMarkUnread,
			PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
			Sender:    tc.MakeEventSender(tc.client.GetCurrentUserID()),
			Timestamp: evt.Timestamp,
		},
		Unread: true,
	}).Success {
		return false
	}
	tc.recordMarkedUnreadSequence(ctx, meta, key, evt.SequenceID)
	return true
}

func (tc *TwitterClient) recordMarkedUnreadSequence(ctx context.Context, meta *UserLoginMetadata, key, sequenceID string) {
	tc.userLocalLock.Lock()
	defer tc.userLocalLock.Unlock()
	if compareIntStrings(sequenceID, meta.MarkedUnreadSequenceIDs[key]) <= 0 {
		return
	}
	if meta.MarkedUnreadSequenceIDs == nil {
		meta.MarkedUnreadSequenceIDs = make(map[string]string)
	}
	meta.MarkedUnreadSequenceIDs[key] = sequenceID
	if err := tc.userLogin.Save(ctx); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save marked unread sequence")
	}
}

// setConversationPinned updates the stored list after our own pin or unpin, since X sends no snapshot back for it.
func (tc *TwitterClient) setConversationPinned(ctx context.Context, conversationID string, pinned bool) {
	meta := tc.loginMetadata()
	if meta == nil {
		return
	}
	conversationID = NormalizeConversationID(conversationID)
	tc.userLocalLock.Lock()
	defer tc.userLocalLock.Unlock()
	meta.PinnedConversationIDs = slices.DeleteFunc(meta.PinnedConversationIDs, func(id string) bool { return id == conversationID })
	if pinned {
		meta.PinnedConversationIDs = append(meta.PinnedConversationIDs, conversationID)
	}
	if err := tc.userLogin.Save(ctx); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to save pinned conversations")
	}
}
