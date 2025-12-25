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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"go.mau.fi/util/variationselector"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var (
	_ bridgev2.ReactionHandlingNetworkAPI    = (*TwitterClient)(nil)
	_ bridgev2.ReadReceiptHandlingNetworkAPI = (*TwitterClient)(nil)
	_ bridgev2.EditHandlingNetworkAPI        = (*TwitterClient)(nil)
	_ bridgev2.TypingHandlingNetworkAPI      = (*TwitterClient)(nil)
	_ bridgev2.ChatViewingNetworkAPI         = (*TwitterClient)(nil)
	_ bridgev2.DeleteChatHandlingNetworkAPI  = (*TwitterClient)(nil)
	_ bridgev2.RedactionHandlingNetworkAPI   = (*TwitterClient)(nil)
	_ bridgev2.MembershipHandlingNetworkAPI  = (*TwitterClient)(nil)
	_ bridgev2.RoomAvatarHandlingNetworkAPI  = (*TwitterClient)(nil)
	_ bridgev2.RoomNameHandlingNetworkAPI    = (*TwitterClient)(nil)
	_ bridgev2.TagHandlingNetworkAPI         = (*TwitterClient)(nil)
	_ bridgev2.MuteHandlingNetworkAPI        = (*TwitterClient)(nil)
)

var _ bridgev2.TransactionIDGeneratingNetwork = (*TwitterConnector)(nil)

func (tc *TwitterClient) HandleMatrixTyping(ctx context.Context, msg *bridgev2.MatrixTyping) error {
	if msg.IsTyping && msg.Type == bridgev2.TypingTypeText {
		return tc.client.SendTypingNotification(ctx, string(msg.Portal.ID))
	}
	return nil
}

func (tc *TwitterConnector) GenerateTransactionID(userID id.UserID, roomID id.RoomID, eventType event.Type) networkid.RawTransactionID {
	return networkid.RawTransactionID(uuid.NewString())
}

func (tc *TwitterClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	conversationID := string(msg.Portal.ID)
	content := msg.Content

	text := content.Body
	if content.Format == event.FormatHTML {
		text = tc.matrixParser.Parse(content.FormattedBody, format.NewContext(ctx))
	}

	messageID := string(msg.InputTransactionID)
	if messageID == "" {
		messageID = uuid.NewString()
	}

	opts := twittermeow.SendEncryptedMessageOpts{
		ConversationID: conversationID,
		MessageID:      messageID,
		Text:           text,
	}

	if msg.ReplyTo != nil {
		replySeqID := string(msg.ReplyTo.ID)
		replyMsgID := replySeqID
		var replyText string
		var replyDisplayName string
		var senderIDStr string

		var metaCopy MessageMetadata
		if meta, ok := msg.ReplyTo.Metadata.(*MessageMetadata); ok && meta != nil {
			metaCopy = *meta
		}
		if extra := tc.lookupReplyMetadata(ctx, msg.Portal.PortalKey, msg.ReplyTo.ID); extra != nil {
			metaCopy.CopyFrom(extra)
		}

		if metaCopy.XChatSequenceID != "" {
			replySeqID = metaCopy.XChatSequenceID
		}
		if metaCopy.XChatClientMsgID != "" {
			replyMsgID = metaCopy.XChatClientMsgID
		}
		replyText = metaCopy.MessageText
		replyDisplayName = metaCopy.SenderDisplayName
		senderIDStr = metaCopy.SenderID

		if senderIDStr == "" && msg.ReplyTo.SenderID != "" {
			senderIDStr = string(msg.ReplyTo.SenderID)
		}
		if replyDisplayName == "" && senderIDStr != "" {
			replyDisplayName = tc.getDisplayNameForUser(ctx, senderIDStr)
			if replyDisplayName == "" {
				replyDisplayName = senderIDStr
			}
		}
		var senderIDPtr *int64
		if senderIDStr != "" {
			if parsed, err := strconv.ParseInt(senderIDStr, 10, 64); err == nil {
				senderIDPtr = &parsed
			} else {
				zerolog.Ctx(ctx).Debug().
					Str("raw_sender_id", senderIDStr).
					Err(err).
					Msg("Failed to parse sender_id for reply preview")
			}
		}
		if replyMsgID == "" {
			replyMsgID = replySeqID
		}
		opts.ReplyTo = &payload.ReplyingToPreview{
			ReplyingToMessageId:         &replyMsgID,
			ReplyingToMessageSequenceId: &replySeqID,
			MessageText:                 &replyText,
			SenderDisplayName:           ptr.Ptr(replyDisplayName),
			SenderId:                    senderIDPtr,
		}
	}

	switch content.MsgType {
	case event.MsgText:
		// nothing extra
	case event.MsgVideo, event.MsgImage, event.MsgAudio:
		if content.FileName != "" && (content.Body == content.FileName || content.Body == "") {
			opts.Text = ""
		}

		data, err := tc.connector.br.Bot.DownloadMedia(ctx, content.URL, content.File)
		if err != nil {
			return nil, err
		}

		// Convert audio to mp4 if needed
		if content.MsgType == event.MsgAudio && content.Info.MimeType != "video/mp4" {
			converted, err := tc.client.ConvertAudioPayload(ctx, data, content.Info.MimeType)
			if err != nil {
				return nil, err
			}
			data = converted
		}

		// Upload media using encrypted XChat flow
		uploadResult, err := tc.client.UploadXChatMedia(ctx, conversationID, messageID, data)
		if err != nil {
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().
			Str("media_hash_key", uploadResult.MediaHashKey).
			Msg("Successfully uploaded encrypted media to XChat")

		attType := payload.MediaTypeImage
		switch content.MsgType {
		case event.MsgVideo:
			attType = payload.MediaTypeVideo
		case event.MsgAudio:
			attType = payload.MediaTypeAudio
		}
		width := int64(content.Info.Width)
		height := int64(content.Info.Height)
		size := int64(len(data))
		opts.Attachments = append(opts.Attachments, &payload.MessageAttachment{
			Media: &payload.MediaAttachment{
				MediaHashKey: &uploadResult.MediaHashKey,
				Type:         ptr.Ptr(int32(attType)),
				Dimensions: &payload.MediaDimensions{
					Width:  &width,
					Height: &height,
				},
				FilesizeBytes: &size,
				Filename:      ptr.Ptr(content.FileName),
			},
		})
	default:
		return nil, fmt.Errorf("%w %s", bridgev2.ErrUnsupportedMessageType, content.MsgType)
	}

	txnID := networkid.TransactionID(messageID)
	dbMsg := &database.Message{
		ID:        networkid.MessageID(messageID),
		MXID:      msg.Event.ID,
		Room:      msg.Portal.PortalKey,
		SenderID:  networkid.UserID(tc.userLogin.ID),
		Timestamp: time.Now(),
		Metadata: &MessageMetadata{
			XChatClientMsgID:  messageID,
			MessageText:       text,
			SenderID:          string(tc.userLogin.ID),
			SenderDisplayName: tc.getDisplayNameForUser(ctx, string(tc.userLogin.ID)),
		},
	}
	msg.AddPendingToSave(dbMsg, txnID, func(remote bridgev2.RemoteMessage, db *database.Message) (bool, error) {
		// Store the real (numeric) XChat message ID when the remote echo arrives.
		if remote != nil {
			db.ID = remote.GetID()
			if meta, ok := db.Metadata.(*MessageMetadata); ok && meta != nil {
				meta.XChatSequenceID = string(remote.GetID())
			}
		}
		return true, nil
	})

	if _, err := tc.client.SendEncryptedMessage(ctx, opts); err != nil {
		return nil, err
	}

	return &bridgev2.MatrixMessageResponse{
		DB:      dbMsg,
		Pending: true,
	}, nil
}

// lookupReplyMetadata fetches message metadata for a given message ID across all parts and merges them.
func (tc *TwitterClient) lookupReplyMetadata(ctx context.Context, portalKey networkid.PortalKey, msgID networkid.MessageID) *MessageMetadata {
	msgs, err := tc.connector.br.DB.Message.GetAllPartsByID(ctx, portalKey.Receiver, msgID)
	if err != nil {
		zerolog.Ctx(ctx).Debug().
			Err(err).
			Str("conversation_id", string(portalKey.ID)).
			Str("reply_to_id", string(msgID)).
			Msg("Failed to load reply target metadata from DB")
		return nil
	}
	var merged MessageMetadata
	for _, m := range msgs {
		if meta, ok := m.Metadata.(*MessageMetadata); ok && meta != nil {
			merged.CopyFrom(meta)
		}
	}
	return &merged
}

// lookupMessageSequenceID returns the XChat sequence ID for a message, looking up from metadata.
func (tc *TwitterClient) lookupMessageSequenceID(ctx context.Context, portalKey networkid.PortalKey, msgID networkid.MessageID) string {
	msgs, err := tc.connector.br.DB.Message.GetAllPartsByID(ctx, portalKey.Receiver, msgID)
	if err != nil {
		return ""
	}
	for _, m := range msgs {
		if meta, ok := m.Metadata.(*MessageMetadata); ok && meta != nil && meta.XChatSequenceID != "" {
			return meta.XChatSequenceID
		}
	}
	return ""
}

func (tc *TwitterClient) HandleMatrixReactionRemove(ctx context.Context, msg *bridgev2.MatrixReactionRemove) error {
	var senderID string
	if msg.TargetReaction != nil {
		senderID = string(msg.TargetReaction.SenderID)
	}
	conversationID := string(msg.Portal.ID)
	targetMessageID := string(msg.TargetReaction.MessageID)

	// Look up the XChat sequence ID for deduplication and sending
	if seqID := tc.lookupMessageSequenceID(ctx, msg.Portal.PortalKey, msg.TargetReaction.MessageID); seqID != "" {
		targetMessageID = seqID
	}

	emoji := variationselector.FullyQualify(msg.TargetReaction.Emoji)
	zerolog.Ctx(ctx).Info().
		Str("conversation_id", conversationID).
		Str("target_message_id", targetMessageID).
		Str("emoji", emoji).
		Str("sender_id", senderID).
		Str("sender_mxid", msg.Event.Sender.String()).
		Msg("Handling Matrix reaction removal")
	return tc.doHandleMatrixReaction(ctx, true, conversationID, targetMessageID, emoji)
}

func (tc *TwitterClient) PreHandleMatrixReaction(_ context.Context, msg *bridgev2.MatrixReaction) (bridgev2.MatrixReactionPreResponse, error) {
	emoji := variationselector.FullyQualify(msg.Content.RelatesTo.Key)
	return bridgev2.MatrixReactionPreResponse{
		SenderID:     networkid.UserID(tc.client.GetCurrentUserID()),
		EmojiID:      networkid.EmojiID(emoji),
		Emoji:        emoji,
		MaxReactions: 1,
	}, nil
}

func (tc *TwitterClient) HandleMatrixReaction(ctx context.Context, msg *bridgev2.MatrixReaction) (reaction *database.Reaction, err error) {
	conversationID := string(msg.Portal.ID)
	targetMessageID := string(msg.TargetMessage.ID)

	// Look up the XChat sequence ID for deduplication and sending
	if msg.TargetMessage.Metadata != nil {
		if meta, ok := msg.TargetMessage.Metadata.(*MessageMetadata); ok && meta.XChatSequenceID != "" {
			targetMessageID = meta.XChatSequenceID
		}
	}
	if seqID := tc.lookupMessageSequenceID(ctx, msg.Portal.PortalKey, msg.TargetMessage.ID); seqID != "" {
		targetMessageID = seqID
	}

	emoji := msg.PreHandleResp.Emoji
	senderID := msg.PreHandleResp.SenderID
	if senderID == "" {
		senderID = networkid.UserID(tc.client.GetCurrentUserID())
	}
	zerolog.Ctx(ctx).Info().
		Str("conversation_id", conversationID).
		Str("target_message_id", targetMessageID).
		Str("emoji", emoji).
		Str("sender_id", string(senderID)).
		Str("sender_mxid", msg.Event.Sender.String()).
		Msg("Handling Matrix reaction")
	if err := tc.doHandleMatrixReaction(ctx, false, conversationID, targetMessageID, emoji); err != nil {
		return nil, err
	}

	return &database.Reaction{
		Room:          msg.Portal.PortalKey,
		MessageID:     msg.TargetMessage.ID,
		MessagePartID: msg.TargetMessage.PartID,
		SenderID:      senderID,
		SenderMXID:    msg.Event.Sender,
		EmojiID:       msg.PreHandleResp.EmojiID,
		MXID:          msg.Event.ID,
		Timestamp:     time.Now(),
		Emoji:         emoji,
	}, nil
}

func (tc *TwitterClient) doHandleMatrixReaction(ctx context.Context, remove bool, conversationID, messageID, emoji string) error {
	// XChat reactions are sent as encrypted MessageCreateEvents (reaction_add/reaction_remove).
	resp, err := tc.client.SendEncryptedReaction(ctx, conversationID, messageID, emoji, remove)
	if err != nil {
		return err
	}
	tc.client.Logger.Debug().Any("reactionResponse", resp).Msg("Reaction response")
	return nil
}

func (tc *TwitterClient) HandleMatrixReadReceipt(ctx context.Context, msg *bridgev2.MatrixReadReceipt) error {
	conversationID := string(msg.Portal.ID)
	lastReadEventID := ""

	if msg.ExactMessage != nil {
		lastReadEventID = string(msg.ExactMessage.ID)
	} else {
		lastMessage, err := tc.userLogin.Bridge.DB.Message.GetLastPartAtOrBeforeTime(ctx, msg.Portal.PortalKey, msg.ReadUpTo)
		if err != nil {
			return err
		}
		lastReadEventID = string(lastMessage.ID)
	}

	readAt := msg.ReadUpTo
	if readAt.IsZero() {
		readAt = time.Now()
	}

	if err := tc.client.SendXChatReadReceipt(ctx, conversationID, lastReadEventID, readAt); err != nil {
		if errors.Is(err, crypto.ErrKeyNotFound) {
			params := &payload.MarkConversationReadQuery{
				ConversationID:  conversationID,
				LastReadEventID: lastReadEventID,
			}
			return tc.client.MarkConversationRead(ctx, params)
		}
		return err
	}

	return nil
}

func (tc *TwitterClient) HandleMatrixEdit(ctx context.Context, edit *bridgev2.MatrixEdit) error {
	targetMessageID := string(edit.EditTarget.ID)
	var meta *MessageMetadata
	if edit.EditTarget != nil {
		if typedMeta, ok := edit.EditTarget.Metadata.(*MessageMetadata); ok {
			meta = typedMeta
		}
	}
	if meta != nil && meta.XChatSequenceID != "" {
		targetMessageID = meta.XChatSequenceID
	}

	messageID := string(edit.InputTransactionID)
	if messageID == "" {
		messageID = uuid.NewString()
	}

	resp, err := tc.client.SendEncryptedEdit(ctx, twittermeow.SendEncryptedEditOpts{
		ConversationID:          string(edit.Portal.ID),
		MessageID:               messageID,
		TargetMessageSequenceID: targetMessageID,
		UpdatedText:             edit.Content.Body,
	})
	if err != nil {
		return err
	}
	tc.client.Logger.Debug().Any("editResponse", resp).Msg("Edit response")
	if meta != nil {
		meta.EditCount++
		meta.MessageText = edit.Content.Body
	}
	return nil
}

func (tc *TwitterClient) HandleMatrixViewingChat(ctx context.Context, chat *bridgev2.MatrixViewingChat) error {
	conversationID := ""
	if chat.Portal != nil {
		conversationID = string(chat.Portal.ID)
	}
	tc.client.SetActiveConversation(conversationID)
	return nil
}

func (tc *TwitterClient) HandleMatrixDeleteChat(ctx context.Context, chat *bridgev2.MatrixDeleteChat) error {
	if chat.Content.DeleteForEveryone {
		return errors.New("delete for everyone is not supported")
	}
	conversationID := string(chat.Portal.ID)
	reqQuery := payload.DMRequestQuery{}.Default()
	return tc.client.DeleteConversation(ctx, conversationID, &reqQuery)
}

func (tc *TwitterClient) HandleMatrixMessageRemove(ctx context.Context, msg *bridgev2.MatrixMessageRemove) error {
	conversationID := string(msg.Portal.ID)
	if msg.TargetMessage == nil {
		return errors.New("target message not found")
	}

	sequenceID := string(msg.TargetMessage.ID)
	if sequenceID == "" {
		return errors.New("message sequence ID not found")
	}

	return tc.client.DeleteXChatMessage(ctx, twittermeow.DeleteXChatMessageOpts{
		ConversationID: conversationID,
		SequenceIDs:    []string{sequenceID},
		DeleteForAll:   false, // Delete for self only
	})
}

func (tc *TwitterClient) HandleMatrixRoomAvatar(ctx context.Context, msg *bridgev2.MatrixRoomAvatar) (bool, error) {
	if msg.Portal.RoomType == database.RoomTypeDM {
		return false, errors.New("cannot set room avatar for DM")
	}

	if msg.Content.URL != "" || msg.Content.MSC3414File != nil {
		data, err := msg.Portal.Bridge.Bot.DownloadMedia(ctx, msg.Content.URL, msg.Content.MSC3414File)
		if err != nil {
			return false, fmt.Errorf("failed to download avatar: %w", err)
		}

		var mediaType string
		if msg.Content.Info != nil {
			mediaType = msg.Content.Info.MimeType
		} else {
			mediaType = http.DetectContentType(data)
		}

		uploadMediaParams := &payload.UploadMediaQuery{
			MediaType: mediaType,
		}
		uploadedMediaResponse, err := tc.client.UploadMedia(ctx, uploadMediaParams, data)
		if err != nil {
			return false, err
		}

		updateAvatarParams := &payload.DMRequestQuery{
			AvatarID: uploadedMediaResponse.MediaIDString,
		}
		err = tc.client.UpdateConversationAvatar(ctx, string(msg.Portal.ID), updateAvatarParams)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("avatar not found")
}

func (tc *TwitterClient) HandleMatrixRoomName(ctx context.Context, msg *bridgev2.MatrixRoomName) (bool, error) {
	if msg.Portal.RoomType == database.RoomTypeDM {
		return false, errors.New("cannot set room name for DM")
	}

	updateNameParams := &payload.DMRequestQuery{
		Name: msg.Content.Name,
	}
	err := tc.client.UpdateConversationName(ctx, string(msg.Portal.ID), updateNameParams)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (tc *TwitterClient) HandleMatrixMembership(ctx context.Context, msg *bridgev2.MatrixMembershipChange) (bool, error) {
	if msg.Type != bridgev2.Invite {
		return false, errors.New("unsupported membership change type")
	}
	if msg.Portal.RoomType == database.RoomTypeDM {
		return false, errors.New("cannot change members for DM")
	}

	var participantID string
	switch target := msg.Target.(type) {
	case *bridgev2.Ghost:
		participantID = string(target.ID)
	case *bridgev2.UserLogin:
		participantID = string(target.ID)
	}
	_, err := tc.client.AddParticipants(ctx, &payload.AddParticipantsPayload{
		ConversationID:    string(msg.Portal.ID),
		AddedParticipants: []string{participantID},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (tc *TwitterClient) HandleRoomTag(ctx context.Context, msg *bridgev2.MatrixRoomTag) error {
	conversationID := string(msg.Portal.ID)
	_, isFavourite := msg.Content.Tags[event.RoomTagFavourite]

	if isFavourite {
		return tc.client.SendXChatPinConversation(ctx, conversationID)
	}
	return tc.client.SendXChatUnpinConversation(ctx, conversationID)
}

func (tc *TwitterClient) HandleMute(ctx context.Context, msg *bridgev2.MatrixMute) error {
	conversationID := string(msg.Portal.ID)
	if msg.Content.IsMuted() {
		return tc.client.MuteConversation(ctx, conversationID)
	}
	return tc.client.UnmuteConversation(ctx, conversationID)
}
