package connector

import (
	"context"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/simplevent"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func xchatDisappearSetting(ttl time.Duration) database.DisappearingSetting {
	if ttl <= 0 {
		return database.DisappearingSetting{}
	}
	return database.DisappearingSetting{Type: event.DisappearingTypeAfterRead, Timer: ttl}
}

// markXChatDisappearTimerApplied records the sequence ID of a timer change or message whose timer was applied,
// so older changes replayed from snapshots, catch-up pages or reconnects are ignored.
func markXChatDisappearTimerApplied(portal *bridgev2.Portal, sequenceID string) bool {
	meta, ok := portal.Metadata.(*PortalMetadata)
	if !ok || meta == nil || sequenceID == "" || compareIntStrings(sequenceID, meta.DisappearTimerSequenceID) <= 0 {
		return false
	}
	meta.DisappearTimerSequenceID = sequenceID
	return true
}

func isNewerXChatDisappearTimer(portal *bridgev2.Portal, sequenceID string) bool {
	meta, ok := portal.Metadata.(*PortalMetadata)
	return ok && meta != nil && sequenceID != "" && compareIntStrings(sequenceID, meta.DisappearTimerSequenceID) > 0
}

type xchatDisappearTimerChange struct {
	simplevent.EventMeta
	SequenceID string
	TTL        time.Duration
}

var _ bridgev2.RemoteChatInfoChange = (*xchatDisappearTimerChange)(nil)

func (tc *TwitterClient) makeXChatDisappearTimerChange(evt *types.XChatConversationTTLChange) *xchatDisappearTimerChange {
	return &xchatDisappearTimerChange{
		EventMeta: simplevent.EventMeta{
			Type:        bridgev2.RemoteEventChatInfoChange,
			Sender:      tc.MakeEventSender(evt.SenderID),
			PortalKey:   tc.MakePortalKeyFromID(evt.ConversationID),
			StreamOrder: methods.ParseSnowflakeInt(evt.ID),
			Timestamp:   evt.Timestamp,
		},
		SequenceID: evt.ID,
		TTL:        evt.TTL,
	}
}

func (evt *xchatDisappearTimerChange) GetChatInfoChange(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfoChange, error) {
	if !isNewerXChatDisappearTimer(portal, evt.SequenceID) {
		return nil, nil
	}
	setting := xchatDisappearSetting(evt.TTL)
	return &bridgev2.ChatInfoChange{
		ChatInfo: &bridgev2.ChatInfo{
			Disappear: &setting,
			ExtraUpdates: func(ctx context.Context, portal *bridgev2.Portal) bool {
				return markXChatDisappearTimerApplied(portal, evt.SequenceID)
			},
		},
	}, nil
}

func (tc *TwitterClient) xchatPortalHasRoom(ctx context.Context, conversationID string) bool {
	portal, err := tc.connector.br.GetExistingPortalByKey(ctx, tc.MakePortalKeyFromID(conversationID))
	return err == nil && portal != nil && portal.MXID != ""
}

type pendingXChatDisappearTimer struct {
	SequenceID string
	TTL        time.Duration
}

// rememberPendingXChatDisappearTimer keeps the newest timer change for a conversation without a room,
// because bridgev2 drops chat info changes for portals that don't exist yet.
func (tc *TwitterClient) rememberPendingXChatDisappearTimer(evt *types.XChatConversationTTLChange) {
	key := NormalizeConversationID(evt.ConversationID)
	next := &pendingXChatDisappearTimer{SequenceID: evt.ID, TTL: evt.TTL}
	for {
		existing, loaded := tc.pendingDisappearTimers.LoadOrStore(key, next)
		if !loaded {
			return
		}
		prev := existing.(*pendingXChatDisappearTimer)
		if compareIntStrings(next.SequenceID, prev.SequenceID) <= 0 || tc.pendingDisappearTimers.CompareAndSwap(key, prev, next) {
			return
		}
	}
}

func (tc *TwitterClient) applyPendingXChatDisappearTimer(ctx context.Context, conversationID string, info *bridgev2.ChatInfo) {
	key := NormalizeConversationID(conversationID)
	value, ok := tc.pendingDisappearTimers.Load(key)
	if !ok || info.Disappear != nil {
		return
	} else if tc.xchatPortalHasRoom(ctx, conversationID) {
		tc.pendingDisappearTimers.CompareAndDelete(key, value)
		return
	}
	pending := value.(*pendingXChatDisappearTimer)
	setting := xchatDisappearSetting(pending.TTL)
	info.Disappear = &setting
	info.ExtraUpdates = bridgev2.MergeExtraUpdaters(info.ExtraUpdates, func(_ context.Context, portal *bridgev2.Portal) bool {
		tc.pendingDisappearTimers.CompareAndDelete(key, pending)
		return markXChatDisappearTimerApplied(portal, pending.SequenceID)
	})
}
