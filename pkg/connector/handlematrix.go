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
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var mediaCategoryMap = map[event.MessageType]payload.MediaCategory{
	event.MsgVideo: payload.MEDIA_CATEGORY_DM_VIDEO,
	event.MsgImage: payload.MEDIA_CATEGORY_DM_IMAGE,
}

var (
	_ bridgev2.ReactionHandlingNetworkAPI    = (*TwitterClient)(nil)
	_ bridgev2.ReadReceiptHandlingNetworkAPI = (*TwitterClient)(nil)
	_ bridgev2.EditHandlingNetworkAPI        = (*TwitterClient)(nil)
	_ bridgev2.TypingHandlingNetworkAPI      = (*TwitterClient)(nil)
	_ bridgev2.ChatViewingNetworkAPI         = (*TwitterClient)(nil)
	_ bridgev2.DeleteChatHandlingNetworkAPI  = (*TwitterClient)(nil)
	_ bridgev2.MembershipHandlingNetworkAPI  = (*TwitterClient)(nil)
	_ bridgev2.RoomAvatarHandlingNetworkAPI  = (*TwitterClient)(nil)
	_ bridgev2.RoomNameHandlingNetworkAPI    = (*TwitterClient)(nil)
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

	// Replies: best-effort include target ID.
	if msg.ReplyTo != nil {
		replyID := string(msg.ReplyTo.ID)
		opts.ReplyTo = &payload.ReplyingToPreview{
			ReplyingToMessageId: &replyID,
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

		uploadMediaParams := &payload.UploadMediaQuery{
			MediaCategory: mediaCategoryMap[content.MsgType],
			MediaType:     content.Info.MimeType,
		}
		if content.Info.MimeType == "image/gif" || content.Info.MauGIF {
			uploadMediaParams.MediaCategory = "dm_gif"
		}
		if content.MsgType == event.MsgAudio {
			uploadMediaParams.MediaCategory = "dm_video"
			if content.Info.MimeType != "video/mp4" {
				converted, err := tc.client.ConvertAudioPayload(ctx, data, content.Info.MimeType)
				if err != nil {
					return nil, err
				}
				data = converted
			}
		}

		uploadedMediaResponse, err := tc.client.UploadMedia(ctx, uploadMediaParams, data)
		if err != nil {
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Any("media_info", uploadedMediaResponse).Msg("Successfully uploaded media to twitter's servers")

		attType := payload.MediaTypeImage
		switch content.MsgType {
		case event.MsgVideo:
			attType = payload.MediaTypeVideo
		case event.MsgAudio:
			attType = payload.MediaTypeAudio
		}
		width := int64(content.Info.Width)
		height := int64(content.Info.Height)
		size := int64(uploadedMediaResponse.Size)
		opts.Attachments = append(opts.Attachments, &payload.MessageAttachment{
			Media: &payload.MediaAttachment{
				MediaHashKey: &uploadedMediaResponse.MediaKey,
				Type:         ptr.Ptr(int32(attType)),
				Dimensions: &payload.MediaDimensions{
					Width:  &width,
					Height: &height,
				},
				FilesizeBytes: &size,
				Filename:      ptr.Ptr(content.FileName),
				AttachmentId:  ptr.Ptr(uploadedMediaResponse.MediaIDString),
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
		SenderID:  UserLoginIDToUserID(tc.userLogin.ID),
		Timestamp: time.Now(),
		Metadata:  &MessageMetadata{XChatClientMsgID: messageID},
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

func (tc *TwitterClient) HandleMatrixReactionRemove(ctx context.Context, msg *bridgev2.MatrixReactionRemove) error {
	var senderID string
	if msg.TargetReaction != nil {
		senderID = string(msg.TargetReaction.SenderID)
	}
	zerolog.Ctx(ctx).Info().
		Str("conversation_id", string(msg.Portal.ID)).
		Str("target_message_id", string(msg.TargetReaction.MessageID)).
		Str("emoji", msg.TargetReaction.Emoji).
		Str("sender_id", senderID).
		Str("sender_mxid", msg.Event.Sender.String()).
		Msg("Handling Matrix reaction removal")
	return tc.doHandleMatrixReaction(ctx, true, string(msg.Portal.ID), string(msg.TargetReaction.MessageID), msg.TargetReaction.Emoji)
}

func (tc *TwitterClient) PreHandleMatrixReaction(_ context.Context, msg *bridgev2.MatrixReaction) (bridgev2.MatrixReactionPreResponse, error) {
	emoji := variationselector.FullyQualify(msg.Content.RelatesTo.Key)
	return bridgev2.MatrixReactionPreResponse{
		SenderID:     UserLoginIDToUserID(tc.userLogin.ID),
		EmojiID:      networkid.EmojiID(emoji),
		Emoji:        emoji,
		MaxReactions: 1,
	}, nil
}

func (tc *TwitterClient) HandleMatrixReaction(ctx context.Context, msg *bridgev2.MatrixReaction) (reaction *database.Reaction, err error) {
	zerolog.Ctx(ctx).Info().
		Str("conversation_id", string(msg.Portal.ID)).
		Str("target_message_id", string(msg.TargetMessage.ID)).
		Str("emoji", msg.PreHandleResp.Emoji).
		Str("sender_id", string(msg.PreHandleResp.SenderID)).
		Str("sender_mxid", msg.Event.Sender.String()).
		Msg("Handling Matrix reaction")
	if err := tc.doHandleMatrixReaction(ctx, false, string(msg.Portal.ID), string(msg.TargetMessage.ID), msg.PreHandleResp.Emoji); err != nil {
		return nil, err
	}

	return &database.Reaction{
		Room:          msg.Portal.PortalKey,
		MessageID:     msg.TargetMessage.ID,
		MessagePartID: msg.TargetMessage.PartID,
		SenderID:      msg.PreHandleResp.SenderID,
		SenderMXID:    msg.Event.Sender,
		EmojiID:       msg.PreHandleResp.EmojiID,
		MXID:          msg.Event.ID,
		Timestamp:     time.Now(),
		Emoji:         msg.PreHandleResp.Emoji,
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
	params := &payload.MarkConversationReadQuery{
		ConversationID: string(msg.Portal.ID),
	}

	if msg.ExactMessage != nil {
		params.LastReadEventID = string(msg.ExactMessage.ID)
	} else {
		lastMessage, err := tc.userLogin.Bridge.DB.Message.GetLastPartAtOrBeforeTime(ctx, msg.Portal.PortalKey, msg.ReadUpTo)
		if err != nil {
			return err
		}
		params.LastReadEventID = string(lastMessage.ID)
	}

	return tc.client.MarkConversationRead(ctx, params)
}

func (tc *TwitterClient) HandleMatrixEdit(ctx context.Context, edit *bridgev2.MatrixEdit) error {
	req := &payload.EditDirectMessagePayload{
		ConversationID: string(edit.Portal.ID),
		RequestID:      string(edit.InputTransactionID),
		DMID:           string(edit.EditTarget.ID),
		Text:           edit.Content.Body,
	}
	if req.RequestID == "" {
		req.RequestID = uuid.NewString()
	}
	resp, err := tc.client.EditDirectMessage(ctx, req)
	if err != nil {
		return err
	}
	edit.EditTarget.Metadata.(*MessageMetadata).EditCount = resp.MessageData.EditCount
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
