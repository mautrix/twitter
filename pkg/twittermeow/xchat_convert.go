package twittermeow

import (
	"strconv"

	"go.mau.fi/util/ptr"
	"go.mau.fi/util/variationselector"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// convertXChatMessageToTwitterMessage converts an XChat MessageEvent with decrypted contents to a types.Message.
// keyVersion is the conversation key version attached to the message, if any.
func convertXChatMessageToTwitterMessage(evt *payload.MessageEvent, contents *payload.MessageContents, keyVersion string) *types.Message {
	seqID := ptr.Val(evt.SequenceId)
	msgData := types.MessageData{
		ID:                     seqID,
		Time:                   ptr.Val(evt.CreatedAtMsec),
		SenderID:               ptr.Val(evt.SenderId),
		Text:                   ptr.Val(contents.MessageText),
		ConversationKeyVersion: keyVersion,
	}

	// Convert entities
	if len(contents.Entities) > 0 {
		msgData.Entities = convertXChatEntities(contents.Entities)
	}

	// Convert attachments
	if len(contents.Attachments) > 0 {
		msgData.Attachment = convertXChatAttachments(contents.Attachments)
		msgData.OriginalAttachments = contents.Attachments
	}

	// Convert reply data
	if contents.ReplyingToPreview != nil {
		msgData.ReplyData = convertXChatReplyPreview(contents.ReplyingToPreview)
	}

	return &types.Message{
		ID:                     seqID,
		Time:                   ptr.Val(evt.CreatedAtMsec),
		SequenceID:             seqID,
		RequestID:              ptr.Val(evt.MessageId),
		ConversationID:         ptr.Val(evt.ConversationId),
		ConversationKeyVersion: keyVersion,
		MessageData:            msgData,
	}
}

// ConvertXChatMessageContentsToMessage converts an XChat MessageEvent with decrypted contents to a types.Message.
// keyVersion is the conversation key version attached to the message, if any.
func ConvertXChatMessageContentsToMessage(evt *payload.MessageEvent, contents *payload.MessageContents, keyVersion string) *types.Message {
	return convertXChatMessageToTwitterMessage(evt, contents, keyVersion)
}

// convertXChatMessageEdit converts an XChat MessageEdit to types.MessageEdit.
func convertXChatMessageEdit(evt *payload.MessageEvent, edit *payload.MessageEdit, keyVersion string) *types.MessageEdit {
	targetMsgID := ptr.Val(edit.MessageSequenceId)
	if targetMsgID == "" {
		targetMsgID = ptr.Val(evt.MessageId)
	}
	if targetMsgID == "" {
		targetMsgID = ptr.Val(evt.SequenceId)
	}

	msgData := types.MessageData{
		ID:                     targetMsgID,
		Time:                   ptr.Val(evt.CreatedAtMsec),
		SenderID:               ptr.Val(evt.SenderId),
		Text:                   ptr.Val(edit.UpdatedText),
		ConversationKeyVersion: keyVersion,
	}
	if len(edit.Entities) > 0 {
		msgData.Entities = convertXChatEntities(edit.Entities)
	}

	return (*types.MessageEdit)(&types.Message{
		ID:                     ptr.Val(evt.SequenceId),
		Time:                   ptr.Val(evt.CreatedAtMsec),
		SequenceID:             ptr.Val(evt.SequenceId),
		RequestID:              ptr.Val(evt.MessageId),
		ConversationID:         ptr.Val(evt.ConversationId),
		ConversationKeyVersion: keyVersion,
		MessageData:            msgData,
	})
}

// convertXChatEntities converts XChat RichTextEntity to types.Entities
func convertXChatEntities(entities []*payload.RichTextEntity) *types.Entities {
	result := &types.Entities{
		Hashtags:     []any{},
		Symbols:      []any{},
		UserMentions: []types.UserMention{},
		URLs:         []types.URLs{},
	}

	for _, entity := range entities {
		if entity == nil || entity.Content == nil {
			continue
		}

		startIdx := int(ptr.Val(entity.StartIndex))
		endIdx := int(ptr.Val(entity.EndIndex))
		indices := []int{startIdx, endIdx}

		content := entity.Content

		// Handle different entity types
		if content.Url != nil {
			result.URLs = append(result.URLs, types.URLs{
				Indices: indices,
			})
		}
		if content.Mention != nil {
			result.UserMentions = append(result.UserMentions, types.UserMention{
				Indices: indices,
			})
		}
		if content.Hashtag != nil {
			result.Hashtags = append(result.Hashtags, map[string]any{
				"indices": indices,
			})
		}
		if content.Cashtag != nil {
			result.Symbols = append(result.Symbols, map[string]any{
				"indices": indices,
			})
		}
	}

	return result
}

// convertXChatAttachments converts XChat MessageAttachment to types.Attachment
func convertXChatAttachments(attachments []*payload.MessageAttachment) *types.Attachment {
	if len(attachments) == 0 {
		return nil
	}

	result := &types.Attachment{}

	for _, att := range attachments {
		if att == nil {
			continue
		}

		if att.Media != nil {
			media := att.Media
			info := types.AttachmentInfo{
				IDStr:         ptr.Val(media.AttachmentId),
				MediaURLHTTPS: ptr.Val(media.LegacyMediaUrlHttps),
				MediaHashKey:  ptr.Val(media.MediaHashKey),
				Filename:      ptr.Val(media.Filename),
				FilesizeBytes: ptr.Val(media.FilesizeBytes),
			}

			if media.Dimensions != nil {
				info.OriginalInfo = types.OriginalInfo{
					Width:  int(ptr.Val(media.Dimensions.Width)),
					Height: int(ptr.Val(media.Dimensions.Height)),
				}
			}

			mediaType := payload.MediaType(ptr.Val(media.Type))
			switch mediaType {
			case payload.MediaTypeImage:
				info.Type = "photo"
				result.Photo = &info
			case payload.MediaTypeGif:
				info.Type = "animated_gif"
				result.AnimatedGif = &info
			case payload.MediaTypeVideo:
				info.Type = "video"
				result.Video = &info
			case payload.MediaTypeAudio:
				info.Type = "audio"
				info.AudioOnly = true
				result.Video = &info
			default:
				info.Type = "photo"
				result.Photo = &info
			}
		}

		if att.Post != nil {
			result.Tweet = &types.AttachmentTweet{
				ExpandedURL: ptr.Val(att.Post.PostUrl),
			}
		}

		if att.Url != nil {
			// URL attachments are typically cards
			result.Card = &types.AttachmentCard{
				BindingValues: types.AttachmentCardBinding{
					CardURL: types.AttachmentCardBindingValue{
						StringValue: ptr.Val(att.Url.Url),
					},
					Title: types.AttachmentCardBindingValue{
						StringValue: ptr.Val(att.Url.DisplayTitle),
					},
				},
			}
			// Store hash key for banner image download
			if att.Url.BannerImageMediaHashKey != nil {
				result.URLBannerMediaHashKey = ptr.Val(att.Url.BannerImageMediaHashKey.MediaHashKey)
			}
		}
	}

	return result
}

// convertXChatReplyPreview converts XChat ReplyingToPreview to types.ReplyData
func convertXChatReplyPreview(preview *payload.ReplyingToPreview) types.ReplyData {
	if preview == nil {
		return types.ReplyData{}
	}

	senderID := ""
	if preview.SenderId != nil {
		senderID = strconv.FormatInt(*preview.SenderId, 10)
	}

	return types.ReplyData{
		ID:       ptr.Val(preview.ReplyingToMessageId),
		SenderID: senderID,
		Text:     ptr.Val(preview.MessageText),
	}
}

// convertXChatReactionAdd converts XChat MessageReactionAdd to types.MessageReactionCreate
func convertXChatReactionAdd(evt *payload.MessageEvent, reaction *payload.MessageReactionAdd) *types.MessageReactionCreate {
	emoji := variationselector.FullyQualify(ptr.Val(reaction.Emoji))
	targetMsgID := ptr.Val(reaction.MessageSequenceId)
	if targetMsgID == "" {
		targetMsgID = ptr.Val(evt.MessageId)
	}
	if targetMsgID == "" {
		targetMsgID = ptr.Val(evt.SequenceId)
	}
	return (*types.MessageReactionCreate)(&types.MessageReaction{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(evt.ConversationId),
		MessageID:      targetMsgID,
		EmojiReaction:  emoji,
		ReactionKey:    emoji,
		SenderID:       ptr.Val(evt.SenderId),
	})
}

// convertXChatReactionRemove converts XChat MessageReactionRemove to types.MessageReactionDelete
func convertXChatReactionRemove(evt *payload.MessageEvent, reaction *payload.MessageReactionRemove) *types.MessageReactionDelete {
	emoji := variationselector.FullyQualify(ptr.Val(reaction.Emoji))
	targetMsgID := ptr.Val(reaction.MessageSequenceId)
	if targetMsgID == "" {
		targetMsgID = ptr.Val(evt.MessageId)
	}
	if targetMsgID == "" {
		targetMsgID = ptr.Val(evt.SequenceId)
	}
	return (*types.MessageReactionDelete)(&types.MessageReaction{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(evt.ConversationId),
		MessageID:      targetMsgID,
		EmojiReaction:  emoji,
		ReactionKey:    emoji,
		SenderID:       ptr.Val(evt.SenderId),
	})
}

// convertXChatMarkReadEvent converts XChat MarkConversationReadEvent to types.ConversationRead
func convertXChatMarkReadEvent(evt *payload.MessageEvent, read *payload.MarkConversationReadEvent) *types.ConversationRead {
	return &types.ConversationRead{
		ID:              ptr.Val(evt.SequenceId),
		Time:            ptr.Val(evt.CreatedAtMsec),
		ConversationID:  ptr.Val(evt.ConversationId),
		LastReadEventID: ptr.Val(read.SeenUntilSequenceId),
	}
}

// convertXChatMessageDelete converts XChat MessageDeleteEvent to types.MessageDelete
func convertXChatMessageDelete(evt *payload.MessageEvent, del *payload.MessageDeleteEvent) *types.MessageDelete {
	messages := make([]types.MessagesDeleted, len(del.SequenceIds))
	for i, seqID := range del.SequenceIds {
		messages[i] = types.MessagesDeleted{
			MessageID:            seqID,
			MessageCreateEventID: seqID,
		}
	}

	return &types.MessageDelete{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(evt.ConversationId),
		Messages:       messages,
	}
}

// convertXChatConversationDelete converts XChat ConversationDeleteEvent to types.ConversationDelete
func convertXChatConversationDelete(evt *payload.MessageEvent, del *payload.ConversationDeleteEvent) *types.ConversationDelete {
	return &types.ConversationDelete{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(del.ConversationId),
	}
}

// convertXChatGroupMemberAdd converts XChat GroupMemberAddChange to types.ParticipantsJoin
func convertXChatGroupMemberAdd(evt *payload.MessageEvent, add *payload.GroupMemberAddChange) *types.ParticipantsJoin {
	participants := make([]types.Participant, len(add.MemberIds))
	for i, memberID := range add.MemberIds {
		participants[i] = types.Participant{
			UserID: memberID,
		}
	}

	return &types.ParticipantsJoin{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(evt.ConversationId),
		SenderID:       ptr.Val(evt.SenderId),
		Participants:   participants,
	}
}

// convertXChatGroupMemberRemove converts XChat GroupMemberRemoveChange to types.ParticipantsLeave
func convertXChatGroupMemberRemove(evt *payload.MessageEvent, remove *payload.GroupMemberRemoveChange) *types.ParticipantsLeave {
	participants := make([]types.Participant, len(remove.MemberIds))
	for i, memberID := range remove.MemberIds {
		participants[i] = types.Participant{
			UserID: memberID,
		}
	}

	return &types.ParticipantsLeave{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(evt.ConversationId),
		Participants:   participants,
	}
}

// convertXChatGroupTitleChange converts XChat GroupTitleChange to types.ConversationNameUpdate
func convertXChatGroupTitleChange(evt *payload.MessageEvent, title *payload.GroupTitleChange) *types.ConversationNameUpdate {
	return &types.ConversationNameUpdate{
		ID:               ptr.Val(evt.SequenceId),
		Time:             ptr.Val(evt.CreatedAtMsec),
		ConversationID:   ptr.Val(evt.ConversationId),
		ConversationName: ptr.Val(title.CustomTitle),
		ByUserID:         ptr.Val(evt.SenderId),
	}
}

// convertXChatGroupAvatarChange converts XChat GroupAvatarUrlChange to types.ConversationAvatarUpdate
func convertXChatGroupAvatarChange(evt *payload.MessageEvent, avatar *payload.GroupAvatarUrlChange) *types.ConversationAvatarUpdate {
	return &types.ConversationAvatarUpdate{
		ID:                           ptr.Val(evt.SequenceId),
		Time:                         ptr.Val(evt.CreatedAtMsec),
		ConversationID:               ptr.Val(evt.ConversationId),
		ConversationAvatarImageHttps: ptr.Val(avatar.CustomAvatarUrl),
		ConversationKeyVersion:       ptr.Val(avatar.ConversationKeyVersion),
		ByUserID:                     ptr.Val(evt.SenderId),
	}
}

// convertXChatTypingEvent converts XChat MessageTypingEvent to a typing notification
func convertXChatTypingEvent(evt *payload.MessageEvent, typing *payload.MessageTypingEvent) (conversationID, senderID string) {
	if typing.ConversationId != nil {
		conversationID = *typing.ConversationId
	} else {
		conversationID = ptr.Val(evt.ConversationId)
	}
	senderID = ptr.Val(evt.SenderId)
	return
}

// ConvertXChatUserToUser converts an XChatUser from the initial inbox response
// to a types.User for cache compatibility with existing code.
func ConvertXChatUserToUser(xu *response.XChatUser) *types.User {
	if xu == nil {
		return nil
	}

	user := &types.User{
		IDStr: xu.RestID,
	}

	if xu.Core != nil {
		user.Name = xu.Core.Name
		user.ScreenName = xu.Core.ScreenName
	}

	if xu.Avatar != nil {
		user.ProfileImageURLHTTPS = xu.Avatar.ImageURL
	}

	if xu.Verification != nil {
		user.Verified = xu.Verification.Verified
		user.IsBlueVerified = xu.Verification.IsBlueVerified
	}

	if xu.Privacy != nil {
		user.Protected = xu.Privacy.Protected
	}

	return user
}
