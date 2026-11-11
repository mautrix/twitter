package connector

import (
	"fmt"
	"time"

	"go.mau.fi/util/exfmt"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func convertXChatNotice(portal *bridgev2.Portal, notice *types.MessageNotice) *bridgev2.ConvertedMessage {
	content := &event.MessageEventContent{
		MsgType: event.MsgNotice,
		Body:    xchatNoticeText(portal, notice),
	}
	switch notice.Type {
	case types.MessageNoticeCallStarted, types.MessageNoticeCallEnded, types.MessageNoticeCallMissed:
		callType := event.BeeperActionMessageCallTypeVideo
		if notice.AudioOnly {
			callType = event.BeeperActionMessageCallTypeVoice
		}
		content.BeeperActionMessage = &event.BeeperActionMessage{
			Type:     event.BeeperActionMessageCall,
			CallType: callType,
		}
	}
	return &bridgev2.ConvertedMessage{
		Parts: []*bridgev2.ConvertedMessagePart{{
			Type:       event.EventMessage,
			Content:    content,
			DBMetadata: &MessageMetadata{},
		}},
	}
}

func xchatNoticeText(portal *bridgev2.Portal, notice *types.MessageNotice) string {
	callKind := "Video call"
	if portal != nil && portal.RoomType != database.RoomTypeDM {
		callKind = "Group call"
	} else if notice.AudioOnly {
		callKind = "Audio call"
	}
	switch notice.Type {
	case types.MessageNoticeCallStarted:
		return callKind + " started"
	case types.MessageNoticeCallEnded:
		if notice.DurationSeconds > 0 {
			return fmt.Sprintf("%s ended (%s)", callKind, exfmt.Duration(time.Duration(notice.DurationSeconds)*time.Second))
		}
		return callKind + " ended"
	case types.MessageNoticeCallMissed:
		return callKind + " missed"
	case types.MessageNoticeScreenshot:
		return "Took a screenshot"
	case types.MessageNoticeScreenRecording:
		return "Recorded the screen"
	case types.MessageNoticeUndecryptable:
		return "Unable to decrypt messages in this chat that were encrypted for another X Chat key."
	default:
		return "Unsupported X Chat event"
	}
}
