package connector

import (
	"context"
	"fmt"
	"slices"
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
			Type:    event.EventMessage,
			Content: content,
			DBMetadata: &MessageMetadata{Notice: &XChatNoticeMetadata{
				Type:     notice.Type,
				SenderID: notice.SenderID,
				CallCard: notice.CallCard,
			}},
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

func claimUndecryptableNotice(ctx context.Context, portal *bridgev2.Portal, keyVersion string) error {
	meta, ok := portal.Metadata.(*PortalMetadata)
	if !ok || keyVersion == "" {
		return nil
	} else if slices.Contains(meta.UndecryptableKeyVersions, keyVersion) {
		return fmt.Errorf("%w: undecryptable notice already sent for key version %s", bridgev2.ErrIgnoringRemoteEvent, keyVersion)
	}
	meta.UndecryptableKeyVersions = append(meta.UndecryptableKeyVersions, keyVersion)
	return portal.Save(ctx)
}

func isUndecryptableNotice(existing []*database.Message) bool {
	if len(existing) == 0 {
		return false
	}
	meta, ok := existing[0].Metadata.(*MessageMetadata)
	return ok && meta != nil && meta.Notice != nil && meta.Notice.Type == types.MessageNoticeUndecryptable
}

func (tc *TwitterClient) convertNowDecryptableMessage(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, existing []*database.Message, data *types.MessageData) (*bridgev2.ConvertedEdit, error) {
	converted := tc.convertToMatrix(ctx, portal, intent, data)
	if converted == nil || len(converted.Parts) == 0 {
		return nil, fmt.Errorf("%w: decrypted message has no parts", bridgev2.ErrIgnoringRemoteEvent)
	}
	edit := &bridgev2.ConvertedEdit{}
	for i, part := range converted.Parts {
		if i >= len(existing) {
			edit.AddedParts = &bridgev2.ConvertedMessage{Parts: converted.Parts[i:]}
			break
		}
		editPart := part.ToEditPart(existing[i])
		if meta, ok := editPart.Part.Metadata.(*MessageMetadata); ok && meta != nil {
			meta.Notice = nil
		}
		edit.ModifiedParts = append(edit.ModifiedParts, editPart)
	}
	if len(existing) > len(converted.Parts) {
		edit.DeletedParts = existing[len(converted.Parts):]
	}
	return edit, nil
}
