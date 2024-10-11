package connector

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	bridgeEvt "maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (tc *TwitterClient) MessagesToBackfillMessages(ctx context.Context, messages []types.Message, conv types.Conversation) ([]*bridgev2.BackfillMessage, error) {
	backfilledMessages := make([]*bridgev2.BackfillMessage, 0)
	selfUserId := tc.client.GetCurrentUserID()
	for _, msg := range messages {
		backfilledMessage, err := tc.MessageToBackfillMessage(ctx, msg, conv, selfUserId)
		if err != nil {
			return nil, err
		}
		backfilledMessages = append(backfilledMessages, backfilledMessage)
	}

	return backfilledMessages, nil
}

func (tc *TwitterClient) MessageToBackfillMessage(ctx context.Context, message types.Message, conv types.Conversation, selfUserId string) (*bridgev2.BackfillMessage, error) {
	messageReactions, err := tc.MessageReactionsToBackfillReactions(message.MessageReactions, selfUserId)
	if err != nil {
		return nil, err
	}

	sentAt, err := methods.UnixStringMilliToTime(message.MessageData.Time)
	if err != nil {
		return nil, err
	}

	partId := networkid.PartID("")
	parts := make([]*bridgev2.ConvertedMessagePart, 0)

	textPart := &bridgev2.ConvertedMessagePart{
		ID:   partId,
		Type: bridgeEvt.EventMessage,
		Content: &bridgeEvt.MessageEventContent{
			MsgType: bridgeEvt.MsgText,
			Body:    message.MessageData.Text,
		},
	}

	replyData := message.MessageData.ReplyData
	var replyTo *networkid.MessageOptionalPartID
	if replyData.ID != "" {
		replyTo = &networkid.MessageOptionalPartID{
			MessageID: networkid.MessageID(replyData.ID),
			PartID:    &partId,
		}
	}

	if message.MessageData.Attachment != nil {
		portal, err := tc.connector.br.GetPortalByKey(ctx, tc.MakePortalKey(conv))
		if err != nil {
			return nil, err
		}

		convertedAttachmentPart, indices, err := tc.TwitterAttachmentToMatrix(ctx, portal, tc.userLogin.Bridge.Matrix.BotIntent(), message.MessageData.Attachment)
		if err != nil {
			return nil, err
		}

		RemoveEntityLinkFromText(textPart, indices)
		parts = append(parts, convertedAttachmentPart)
	}

	if len(textPart.Content.Body) > 0 {
		parts = append(parts, textPart)
	}

	return &bridgev2.BackfillMessage{
		ConvertedMessage: &bridgev2.ConvertedMessage{
			ReplyTo: replyTo,
			Parts:   parts,
		},
		Sender: bridgev2.EventSender{
			IsFromMe: message.MessageData.SenderID == selfUserId,
			Sender:   networkid.UserID(message.MessageData.SenderID),
		},
		ID:        networkid.MessageID(message.MessageData.ID),
		Timestamp: sentAt,
		Reactions: messageReactions,
	}, nil
}

func RemoveEntityLinkFromText(msgPart *bridgev2.ConvertedMessagePart, indices []int) {
	start, end := indices[0], indices[1]
	msgPart.Content.Body = msgPart.Content.Body[:start-1] + msgPart.Content.Body[end:]
}

func (tc *TwitterClient) MessageReactionsToBackfillReactions(reactions []types.MessageReaction, selfUserId string) ([]*bridgev2.BackfillReaction, error) {
	backfillReactions := make([]*bridgev2.BackfillReaction, 0)
	for _, reaction := range reactions {
		reactionTime, err := methods.UnixStringMilliToTime(reaction.Time)
		if err != nil {
			return nil, err
		}

		backfillReaction := &bridgev2.BackfillReaction{
			Timestamp: reactionTime,
			Sender: bridgev2.EventSender{
				IsFromMe: reaction.SenderID == selfUserId,
				Sender:   networkid.UserID(reaction.SenderID),
			},
			EmojiID: "",
			Emoji:   reaction.EmojiReaction,
		}
		backfillReactions = append(backfillReactions, backfillReaction)
	}
	return backfillReactions, nil
}

func (tc *TwitterClient) ConversationToChatInfo(conv *types.Conversation) *bridgev2.ChatInfo {
	memberList := tc.ParticipantsToMemberList(conv.Participants)
	return &bridgev2.ChatInfo{
		Name:        &conv.Name,
		Avatar:      MakeAvatar(conv.AvatarImageHttps),
		Members:     memberList,
		Type:        tc.ConversationTypeToRoomType(conv.Type),
		CanBackfill: true,
	}
}

func (tc *TwitterClient) ConversationTypeToRoomType(convType types.ConversationType) *database.RoomType {
	var roomType database.RoomType
	switch convType {
	case types.ONE_TO_ONE:
		roomType = database.RoomTypeDM
	case types.GROUP_DM:
		roomType = database.RoomTypeGroupDM
	}

	return &roomType
}

func (tc *TwitterClient) UsersToMemberList(users []types.User) *bridgev2.ChatMemberList {
	selfUserId := tc.client.GetCurrentUserID()
	chatMembers := make([]bridgev2.ChatMember, len(users)-1)
	for _, user := range users {
		chatMembers = append(chatMembers, tc.UserToChatMember(user, user.IDStr == selfUserId))
	}

	return &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(users),
		Members:          chatMembers, // TODO use membermap instead
	}
}

func (tc *TwitterClient) ParticipantsToMemberList(participants []types.Participant) *bridgev2.ChatMemberList {
	selfUserId := tc.client.GetCurrentUserID()
	chatMembers := make([]bridgev2.ChatMember, len(participants)-1)
	for _, participant := range participants {
		chatMembers = append(chatMembers, tc.ParticipantToChatMember(participant, participant.UserID == selfUserId))
	}

	return &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(participants),
		Members:          chatMembers, // TODO use membermap instead
	}
}

func (tc *TwitterClient) UserToChatMember(user types.User, isFromMe bool) bridgev2.ChatMember {
	return bridgev2.ChatMember{
		EventSender: bridgev2.EventSender{
			IsFromMe: isFromMe,
			Sender:   networkid.UserID(user.IDStr),
		},
		UserInfo: &bridgev2.UserInfo{
			Name:   &user.Name,
			Avatar: MakeAvatar(user.ProfileImageURL),
		},
	}
}

func (tc *TwitterClient) ParticipantToChatMember(participant types.Participant, isFromMe bool) bridgev2.ChatMember {
	return bridgev2.ChatMember{
		EventSender: bridgev2.EventSender{
			IsFromMe: isFromMe,
			Sender:   networkid.UserID(participant.UserID),
			//SenderLogin: networkid.UserLoginID(participant.UserID),
		},
		UserInfo: tc.GetUserInfoBridge(participant.UserID),
	}
}

func (tc *TwitterClient) GetUserInfoBridge(userId string) *bridgev2.UserInfo {
	var userinfo *bridgev2.UserInfo
	if userCacheEntry, ok := tc.userCache[userId]; ok {
		userinfo = &bridgev2.UserInfo{
			Name:   &userCacheEntry.Name,
			Avatar: MakeAvatar(userCacheEntry.ProfileImageURL),
			Identifiers: []string{fmt.Sprintf("twitter:%s", userCacheEntry.ScreenName)},
		}
	}
	return userinfo
}

func (tc *TwitterClient) TwitterAttachmentToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, attachment *types.Attachment) (*bridgev2.ConvertedMessagePart, []int, error) {
	var attachmentInfo *types.AttachmentInfo
	var attachmentURL string
	var mimeType string
	var indices []int
	var msgType bridgeEvt.MessageType
	if attachment.Photo != nil {
		// image attachment
		attachmentInfo = attachment.Photo
		mimeType = "image/jpeg" // attachment doesn't include this specifically
		msgType = bridgeEvt.MsgImage
		attachmentURL = attachmentInfo.MediaURLHTTPS
		indices = attachmentInfo.Indices
	} else if attachment.Video != nil {
		// video attachment
		attachmentInfo = attachment.Video
		mimeType = "video/mp4"
		msgType = bridgeEvt.MsgVideo

		highestBitRateVariant, err := attachmentInfo.VideoInfo.GetHighestBitrateVariant()
		if err != nil {
			return nil, nil, err
		}
		attachmentURL = highestBitRateVariant.URL
		indices = attachmentInfo.Indices
	} else if attachment.AnimatedGif != nil {
		// gif attachment
		attachmentInfo = attachment.AnimatedGif
		mimeType = "image/gif"
		msgType = bridgeEvt.MsgVideo

		highestBitRateVariant, err := attachmentInfo.VideoInfo.GetHighestBitrateVariant()
		if err != nil {
			return nil, nil, err
		}
		attachmentURL = highestBitRateVariant.URL
		indices = attachmentInfo.Indices
	}

	attachmentBytes, err := DownloadPlainFile(ctx, tc.client.GetCookieString(), attachmentURL, "twitter attachment")
	if err != nil {
		return nil, nil, err
	}

	content := convertTwitterAttachmentMetadata(attachmentInfo, mimeType, msgType, attachmentBytes)
	err = uploadMedia(ctx, portal, intent, attachmentBytes, &content)
	if err != nil {
		return nil, nil, err
	}

	return &bridgev2.ConvertedMessagePart{
		ID:      networkid.PartID(fmt.Sprintf("attachment-%s", attachmentInfo.IDStr)),
		Type:    bridgeEvt.EventMessage,
		Content: &content,
	}, indices, nil
}

func uploadMedia(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data []byte, content *bridgeEvt.MessageEventContent) error {
	mxc, file, err := intent.UploadMedia(ctx, portal.MXID, data, "", content.Info.MimeType)
	if err != nil {
		return err
	}
	if file != nil {
		content.File = file
	} else {
		content.URL = mxc
	}
	return nil
}

func convertTwitterAttachmentMetadata(attachmentInfo *types.AttachmentInfo, mimeType string, msgType bridgeEvt.MessageType, attachmentBytes []byte) bridgeEvt.MessageEventContent {
	content := bridgeEvt.MessageEventContent{
		Info: &bridgeEvt.FileInfo{
			MimeType: mimeType,
			Size:     len(attachmentBytes),
		},
		MsgType: msgType,
		Body:    attachmentInfo.IDStr,
	}

	originalInfo := attachmentInfo.OriginalInfo
	if originalInfo.Width != 0 {
		content.Info.Width = originalInfo.Width
	}
	if originalInfo.Height != 0 {
		content.Info.Height = originalInfo.Height
	}

	if attachmentInfo.VideoInfo.DurationMillis != 0 {
		content.Info.Duration = attachmentInfo.VideoInfo.DurationMillis
	}

	return content
}

func MakeAvatar(avatarURL string) *bridgev2.Avatar {
	// idk if this check is needed, maybe the Remove field is enough. Change later if it isn't needed
	if avatarURL == "" {
		return nil
	}
	return &bridgev2.Avatar{
		ID: networkid.AvatarID(avatarURL),
		Get: func(ctx context.Context) ([]byte, error) {
			return DownloadPlainFile(ctx, "", avatarURL, "avatar")
		},
		Remove: avatarURL == "",
	}
}

func DownloadPlainFile(ctx context.Context, cookies, url, thing string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	if cookies != "" {
		req.Header.Add("cookie", cookies)
	}

	getResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", thing, err)
	}

	data, err := io.ReadAll(getResp.Body)
	_ = getResp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read %s data: %w", thing, err)
	}
	return data, nil
}
