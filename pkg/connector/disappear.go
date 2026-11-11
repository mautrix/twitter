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

// applyImplicitXChatDisappearTimer updates the portal timer from a live message, since every X Chat
// message carries the conversation's timer and the timer change itself may have been missed.
func applyImplicitXChatDisappearTimer(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *types.MessageData) {
	if !markXChatDisappearTimerApplied(portal, msg.ID) {
		return
	}
	setting := xchatDisappearSetting(msg.TTL)
	if portal.Disappear.Timer == setting.Timer {
		return
	}
	portal.UpdateDisappearingSetting(ctx, setting, bridgev2.UpdateDisappearingSettingOpts{
		Sender:     intent,
		Timestamp:  methods.ParseMsecTimestamp(msg.Time),
		Implicit:   true,
		Save:       true,
		SendNotice: true,
	})
}
