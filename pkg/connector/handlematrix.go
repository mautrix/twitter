package connector

import (
	"context"

	"github.com/google/uuid"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var (
	MSG_TYPE_TO_MEDIA_TYPE = map[event.MessageType]payload.MediaType{
		event.MsgVideo: payload.MEDIA_TYPE_VIDEO_MP4,
		event.MsgImage: payload.MEDIA_TYPE_IMAGE_JPEG,
	}
	MSG_TYPE_TO_MEDIA_CATEGORY = map[event.MessageType]payload.MediaCategory{
		event.MsgVideo: payload.MEDIA_CATEGORY_DM_VIDEO,
		event.MsgImage: payload.MEDIA_CATEGORY_DM_IMAGE,
	}
)

var (
	_ bridgev2.ReactionHandlingNetworkAPI    = (*TwitterClient)(nil)
	_ bridgev2.ReadReceiptHandlingNetworkAPI = (*TwitterClient)(nil)
	_ bridgev2.EditHandlingNetworkAPI        = (*TwitterClient)(nil)
)

func (tc *TwitterClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	conversationId := string(msg.Portal.ID)
	sendDMPayload := &payload.SendDirectMessagePayload{
		ConversationID:    conversationId,
		IncludeCards:      1,
		IncludeQuoteCount: true,
		RecipientIds:      false,
		DmUsers:           false,
	}

	if msg.ReplyTo != nil {
		sendDMPayload.ReplyToDmID = string(msg.ReplyTo.ID)
	}

	content := msg.Content
	if content.FileName != "" && content.Body != content.FileName {
		sendDMPayload.Text = content.Body
	}

	switch content.MsgType {
	case event.MsgText:
		sendDMPayload.Text = content.Body
	case event.MsgVideo, event.MsgImage:
		file := content.GetFile()
		data, err := tc.connector.br.Bot.DownloadMedia(ctx, file.URL, file)
		if err != nil {
			return nil, err
		}

		uploadMediaParams := &payload.UploadMediaQuery{
			MediaType:     MSG_TYPE_TO_MEDIA_TYPE[content.MsgType],
			MediaCategory: MSG_TYPE_TO_MEDIA_CATEGORY[content.MsgType],
		}
		if content.Info.MimeType == "image/gif" {
			uploadMediaParams.MediaType = "image/gif"
			uploadMediaParams.MediaCategory = "dm_gif"
		}
		uploadedMediaResponse, err := tc.client.UploadMedia(uploadMediaParams, data)
		if err != nil {
			return nil, err
		}

		tc.client.Logger.Debug().Any("media_info", uploadedMediaResponse).Msg("Successfully uploaded media to twitter's servers")
		sendDMPayload.MediaID = uploadedMediaResponse.MediaIDString
	default:
		tc.client.Logger.Warn().Any("msg_type", content.MsgType).Msg("Found unhandled MsgType in HandleMatrixMessage function")
	}

	resp, err := tc.client.SendDirectMessage(sendDMPayload)
	if err != nil {
		return nil, err
	}

	messageData, err := resp.PrettifyMessages(conversationId)
	if err != nil {
		return nil, err
	}

	respMessageData := messageData[0]
	return &bridgev2.MatrixMessageResponse{
		DB: &database.Message{
			ID:        networkid.MessageID(respMessageData.MessageID),
			MXID:      msg.Event.ID,
			Room:      msg.Portal.PortalKey,
			SenderID:  networkid.UserID(tc.client.GetCurrentUserID()),
			Timestamp: respMessageData.SentAt,
		},
	}, nil
}

func (tc *TwitterClient) HandleMatrixReactionRemove(_ context.Context, msg *bridgev2.MatrixReactionRemove) error {
	return tc.doHandleMatrixReaction(true, string(msg.Portal.ID), string(msg.TargetReaction.MessageID), msg.TargetReaction.Emoji)
}

func (tc *TwitterClient) PreHandleMatrixReaction(_ context.Context, msg *bridgev2.MatrixReaction) (bridgev2.MatrixReactionPreResponse, error) {
	return bridgev2.MatrixReactionPreResponse{
		SenderID:     networkid.UserID(tc.userLogin.ID),
		Emoji:        msg.Content.RelatesTo.Key,
		MaxReactions: 1,
	}, nil
}

func (tc *TwitterClient) HandleMatrixReaction(_ context.Context, msg *bridgev2.MatrixReaction) (reaction *database.Reaction, err error) {
	return nil, tc.doHandleMatrixReaction(false, string(msg.Portal.ID), string(msg.TargetMessage.ID), msg.PreHandleResp.Emoji)
}

func (tc *TwitterClient) doHandleMatrixReaction(remove bool, conversationId, messageId, emoji string) error {
	reactionPayload := &payload.ReactionActionPayload{
		ConversationID: conversationId,
		MessageID:      messageId,
		ReactionTypes:  []string{"Emoji"},
		EmojiReactions: []string{emoji},
	}
	reactionResponse, err := tc.client.React(reactionPayload, remove)
	if err != nil {
		return err
	}

	tc.client.Logger.Debug().Any("reactionResponse", reactionResponse).Any("payload", reactionPayload).Msg("Reaction response")
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

	return tc.client.MarkConversationRead(params)
}

func (tc *TwitterClient) HandleMatrixEdit(_ context.Context, edit *bridgev2.MatrixEdit) error {
	_, err := tc.client.EditDirectMessage(&payload.EditDirectMessagePayload{
		ConversationID: string(edit.Portal.ID),
		RequestID:      uuid.New().String(),
		DmID:           string(edit.EditTarget.ID),
		Text:           edit.Content.Body,
	})
	return err
}
