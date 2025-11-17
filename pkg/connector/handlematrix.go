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

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mau.fi/util/variationselector"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
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
	sendDMPayload := &payload.SendDirectMessagePayload{
		ConversationID:    conversationID,
		IncludeCards:      1,
		IncludeQuoteCount: true,
		RecipientIDs:      false,
		DMUsers:           false,
		CardsPlatform:     "Web-12",
		RequestID:         string(msg.InputTransactionID),
	}
	if sendDMPayload.RequestID == "" {
		sendDMPayload.RequestID = uuid.NewString()
	}

	if msg.ReplyTo != nil {
		sendDMPayload.ReplyToDMID = string(msg.ReplyTo.ID)
	}

	content := msg.Content
	if content.Format == event.FormatHTML {
		sendDMPayload.Text = tc.matrixParser.Parse(content.FormattedBody, format.NewContext(ctx))
	} else {
		sendDMPayload.Text = content.Body
	}

	switch content.MsgType {
	case event.MsgText:
		break
	case event.MsgVideo, event.MsgImage, event.MsgAudio:
		if content.FileName == "" || content.Body == content.FileName {
			sendDMPayload.Text = ""
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
			sendDMPayload.AudioOnlyMediaAttachment = true
			uploadMediaParams.MediaCategory = "dm_video"
			if content.Info.MimeType != "video/mp4" {
				converted, err := tc.client.ConvertAudioPayload(ctx, data, content.Info.MimeType)
				if err != nil {
					return nil, err
				} else {
					data = converted
				}
			}
		}
		uploadedMediaResponse, err := tc.client.UploadMedia(ctx, uploadMediaParams, data)
		if err != nil {
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Any("media_info", uploadedMediaResponse).Msg("Successfully uploaded media to twitter's servers")
		sendDMPayload.MediaID = uploadedMediaResponse.MediaIDString
	default:
		return nil, fmt.Errorf("%w %s", bridgev2.ErrUnsupportedMessageType, content.MsgType)
	}

	txnID := networkid.TransactionID(sendDMPayload.RequestID)
	msg.AddPendingToIgnore(txnID)
	resp, err := tc.client.SendDirectMessage(ctx, sendDMPayload)
	if err != nil {
		return nil, err
	} else if len(resp.Entries) == 0 {
		return nil, fmt.Errorf("no entries in send response")
	} else if len(resp.Entries) > 1 {
		zerolog.Ctx(ctx).Warn().
			Int("entry_count", len(resp.Entries)).
			Msg("Unexpected number of entries in send response")
	}
	entry, ok := resp.Entries[0].ParseWithErrorLog(zerolog.Ctx(ctx)).(*types.Message)
	if !ok {
		return nil, fmt.Errorf("unexpected response data: not a message")
	}
	return &bridgev2.MatrixMessageResponse{
		DB: &database.Message{
			ID:        networkid.MessageID(entry.MessageData.ID),
			MXID:      msg.Event.ID,
			Room:      msg.Portal.PortalKey,
			SenderID:  UserLoginIDToUserID(tc.userLogin.ID),
			Timestamp: methods.ParseSnowflake(entry.MessageData.ID),
			Metadata:  &MessageMetadata{},
		},
		StreamOrder:   methods.ParseSnowflakeInt(entry.MessageData.ID),
		RemovePending: txnID,
	}, nil
}

func (tc *TwitterClient) HandleMatrixReactionRemove(ctx context.Context, msg *bridgev2.MatrixReactionRemove) error {
	return tc.doHandleMatrixReaction(ctx, true, string(msg.Portal.ID), string(msg.TargetReaction.MessageID), msg.TargetReaction.Emoji)
}

func (tc *TwitterClient) PreHandleMatrixReaction(_ context.Context, msg *bridgev2.MatrixReaction) (bridgev2.MatrixReactionPreResponse, error) {
	return bridgev2.MatrixReactionPreResponse{
		SenderID:     UserLoginIDToUserID(tc.userLogin.ID),
		Emoji:        variationselector.FullyQualify(msg.Content.RelatesTo.Key),
		MaxReactions: 1,
	}, nil
}

func (tc *TwitterClient) HandleMatrixReaction(ctx context.Context, msg *bridgev2.MatrixReaction) (reaction *database.Reaction, err error) {
	return nil, tc.doHandleMatrixReaction(ctx, false, string(msg.Portal.ID), string(msg.TargetMessage.ID), msg.PreHandleResp.Emoji)
}

func (tc *TwitterClient) doHandleMatrixReaction(ctx context.Context, remove bool, conversationID, messageID, emoji string) error {
	reactionPayload := &payload.ReactionActionPayload{
		ConversationID: conversationID,
		MessageID:      messageID,
		ReactionTypes:  []string{"Emoji"},
		EmojiReactions: []string{emoji},
	}
	reactionResponse, err := tc.client.React(ctx, reactionPayload, remove)
	if err != nil {
		return err
	}
	tc.client.Logger.Debug().Any("reactionResponse", reactionResponse).Any("payload", reactionPayload).Msg("Reaction response")
	if reactionResponse.Data.CreateDmReaction.Typename == "CreateDMReactionFailure" {
		return fmt.Errorf("server rejected reaction")
	}
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
