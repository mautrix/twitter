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
	"encoding/base64"
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
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

var (
	_ bridgev2.ReactionHandlingNetworkAPI        = (*TwitterClient)(nil)
	_ bridgev2.ReadReceiptHandlingNetworkAPI     = (*TwitterClient)(nil)
	_ bridgev2.EditHandlingNetworkAPI            = (*TwitterClient)(nil)
	_ bridgev2.TypingHandlingNetworkAPI          = (*TwitterClient)(nil)
	_ bridgev2.ChatViewingNetworkAPI             = (*TwitterClient)(nil)
	_ bridgev2.DeleteChatHandlingNetworkAPI      = (*TwitterClient)(nil)
	_ bridgev2.RedactionHandlingNetworkAPI       = (*TwitterClient)(nil)
	_ bridgev2.MembershipHandlingNetworkAPI      = (*TwitterClient)(nil)
	_ bridgev2.MessageRequestAcceptingNetworkAPI = (*TwitterClient)(nil)
	_ bridgev2.RoomAvatarHandlingNetworkAPI      = (*TwitterClient)(nil)
	_ bridgev2.RoomNameHandlingNetworkAPI        = (*TwitterClient)(nil)
	_ bridgev2.TagHandlingNetworkAPI             = (*TwitterClient)(nil)
	_ bridgev2.MuteHandlingNetworkAPI            = (*TwitterClient)(nil)
)

var _ bridgev2.TransactionIDGeneratingNetwork = (*TwitterConnector)(nil)

func (tc *TwitterClient) HandleMatrixTyping(ctx context.Context, msg *bridgev2.MatrixTyping) error {
	if !msg.IsTyping || msg.Type != bridgev2.TypingTypeText {
		return nil
	}

	conversationID := string(msg.Portal.ID)
	meta := ensurePortalMetadata(msg.Portal)

	// Use WebSocket for trusted conversations, GraphQL for untrusted
	if meta.IsTrusted() {
		return tc.client.SendXChatTypingNotification(ctx, conversationID)
	}
	return tc.client.SendTypingNotification(ctx, conversationID)
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

		// Get XChatClientMsgID from metadata (still stored for transaction ID matching)
		var metaCopy MessageMetadata
		if meta, ok := msg.ReplyTo.Metadata.(*MessageMetadata); ok && meta != nil {
			metaCopy = *meta
		}
		if extra := tc.lookupReplyMetadata(ctx, msg.Portal.PortalKey, msg.ReplyTo.ID); extra != nil {
			metaCopy.CopyFrom(extra)
		}
		if metaCopy.XChatClientMsgID != "" {
			replyMsgID = metaCopy.XChatClientMsgID
		}

		// Fetch text, display name, and attachments from Matrix event
		replyText, replyDisplayName, replyAttachments, ok := tc.fetchReplyInfoFromMatrix(ctx, msg.Portal, msg.ReplyTo)

		// Get sender ID
		var senderIDStr string
		if msg.ReplyTo.SenderID != "" {
			senderIDStr = string(msg.ReplyTo.SenderID)
		}

		// Fallback for display name if fetch failed
		if replyDisplayName == "" && senderIDStr != "" {
			replyDisplayName = tc.getDisplayNameForUser(ctx, senderIDStr)
			if replyDisplayName == "" {
				replyDisplayName = senderIDStr
			}
		}

		// If we couldn't fetch reply info, skip reply metadata
		if !ok {
			zerolog.Ctx(ctx).Debug().
				Str("reply_to_id", string(msg.ReplyTo.ID)).
				Msg("Could not fetch reply content, sending as standalone message")
		} else {
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
			zerolog.Ctx(ctx).Info().
				Str("conversation_id", conversationID).
				Str("reply_to_id", string(msg.ReplyTo.ID)).
				Int("reply_attachments", len(replyAttachments)).
				Str("reply_text", replyText).
				Msg("Preparing reply preview")
			opts.ReplyTo = &payload.ReplyingToPreview{
				ReplyingToMessageId:         &replyMsgID,
				ReplyingToMessageSequenceId: &replySeqID,
				MessageText:                 &replyText,
				SenderDisplayName:           ptr.Ptr(replyDisplayName),
				SenderId:                    senderIDPtr,
			}
			if len(replyAttachments) > 0 {
				opts.ReplyTo.Attachments = replyAttachments
			}
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

	if opts.Text != "" {
		urlAttachments, urlEntities := buildURLAttachments(opts.Text)
		if len(urlAttachments) > 0 {
			opts.Attachments = append(opts.Attachments, urlAttachments...)
		}
		if len(urlEntities) > 0 {
			opts.Entities = append(opts.Entities, urlEntities...)
		}
		mentionEntities := buildMentionEntities(opts.Text)
		if len(mentionEntities) > 0 {
			opts.Entities = append(opts.Entities, mentionEntities...)
		}
	}

	txnID := networkid.TransactionID(messageID)
	dbMsg := &database.Message{
		ID:        networkid.MessageID(messageID),
		SenderID:  networkid.UserID(tc.userLogin.ID),
		Timestamp: time.Now(),
		Metadata: &MessageMetadata{
			XChatClientMsgID: messageID,
		},
	}
	// Check portal metadata for trust status
	meta := ensurePortalMetadata(msg.Portal)
	if !meta.IsTrusted() {
		// Untrusted conversation - use REST API
		return tc.sendDirectMessageREST(ctx, msg, conversationID, messageID, text, opts, dbMsg, txnID)
	}

	// Trusted conversation - use XChat encrypted protocol, with REST fallback on key/token not found
	resp, err := tc.client.SendEncryptedMessage(ctx, opts)
	if err != nil {
		if errors.Is(err, crypto.ErrKeyNotFound) {
			zerolog.Ctx(ctx).Debug().
				Str("conversation_id", conversationID).
				Msg("Falling back to REST API for message send (key/token not found)")
			return tc.sendDirectMessageREST(ctx, msg, conversationID, messageID, text, opts, dbMsg, txnID)
		}
		return nil, err
	}

	// Extract sequence ID from the response
	if resp != nil && resp.Data.XChatSendCreateMessageEvent.EncodedMessageEvent != "" {
		if decoded, err := twittermeow.DecodeSendMessageEventResponse(resp.Data.XChatSendCreateMessageEvent.EncodedMessageEvent); err == nil && decoded.MessageEvent != nil {
			dbMsg.ID = networkid.MessageID(*decoded.MessageEvent)
		}
	}

	return &bridgev2.MatrixMessageResponse{
		DB: dbMsg,
	}, nil
}

// sendDirectMessageREST sends a message via the REST API for untrusted conversations
func (tc *TwitterClient) sendDirectMessageREST(
	ctx context.Context,
	msg *bridgev2.MatrixMessage,
	conversationID string,
	messageID string,
	text string,
	opts twittermeow.SendEncryptedMessageOpts,
	dbMsg *database.Message,
	_ networkid.TransactionID,
) (*bridgev2.MatrixMessageResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Debug().
		Str("conversation_id", conversationID).
		Str("message_id", messageID).
		Msg("Sending message via REST API (untrusted conversation)")

	pl := &payload.SendDirectMessagePayload{
		ConversationID: conversationID,
		RequestID:      messageID,
		Text:           text,
		CardsPlatform:  "Web-12",
		IncludeCards:   1,
	}

	// Handle replies for REST API
	if opts.ReplyTo != nil && opts.ReplyTo.ReplyingToMessageSequenceId != nil {
		pl.ReplyToDMID = *opts.ReplyTo.ReplyingToMessageSequenceId
	}

	// Handle media attachments for REST API
	content := msg.Content
	switch content.MsgType {
	case event.MsgVideo, event.MsgImage, event.MsgAudio:
		data, err := tc.connector.br.Bot.DownloadMedia(ctx, content.URL, content.File)
		if err != nil {
			return nil, err
		}

		mimeType := content.Info.MimeType
		// Convert audio to mp4 if needed
		if content.MsgType == event.MsgAudio && mimeType != "video/mp4" {
			converted, err := tc.client.ConvertAudioPayload(ctx, data, mimeType)
			if err != nil {
				return nil, err
			}
			data = converted
			mimeType = "video/mp4"
		}

		// Determine media category for REST API
		var mediaCategory payload.MediaCategory
		switch content.MsgType {
		case event.MsgVideo, event.MsgAudio:
			mediaCategory = payload.MEDIA_CATEGORY_DM_VIDEO
		default:
			mediaCategory = payload.MEDIA_CATEGORY_DM_IMAGE
		}

		// Upload media using non-encrypted flow for REST API
		uploadQuery := &payload.UploadMediaQuery{
			MediaType:     mimeType,
			MediaCategory: mediaCategory,
		}
		uploadResult, err := tc.client.UploadMedia(ctx, uploadQuery, data)
		if err != nil {
			return nil, fmt.Errorf("failed to upload media for REST API: %w", err)
		}

		pl.MediaID = uploadResult.MediaIDString
		log.Debug().
			Str("media_id", uploadResult.MediaIDString).
			Msg("Successfully uploaded media for REST API")
	}

	resp, err := tc.client.SendDirectMessage(ctx, pl)
	if err != nil {
		return nil, err
	}

	// Extract the message ID from the response
	if resp != nil && len(resp.Entries) > 0 {
		for _, entry := range resp.Entries {
			parsed := entry.ParseWithErrorLog(log)
			if msgEvt, ok := parsed.(*types.Message); ok && msgEvt.ConversationID == conversationID {
				dbMsg.ID = networkid.MessageID(msgEvt.ID)
				break
			}
		}
	}

	// Successfully sent - mark conversation as trusted
	meta := ensurePortalMetadata(msg.Portal)
	if !meta.Trusted {
		meta.Trusted = true
		if err := msg.Portal.Save(ctx); err != nil {
			log.Warn().Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to save portal metadata with Trusted=true after REST send")
		} else {
			log.Debug().
				Str("conversation_id", conversationID).
				Msg("Marked conversation as trusted after first REST message")
		}
	}

	return &bridgev2.MatrixMessageResponse{
		DB: dbMsg,
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

// fetchReplyInfoFromMatrix fetches reply text and sender display name from the Matrix event.
// Returns (messageText, senderDisplayName, attachments, ok).
// If GetEvent fails or the event can't be parsed, returns empty values with ok=false.
// The caller should send without reply metadata when ok=false.
func (tc *TwitterClient) fetchReplyInfoFromMatrix(ctx context.Context, portal *bridgev2.Portal, replyTo *database.Message) (text string, displayName string, attachments []*payload.MessageAttachment, ok bool) {
	if replyTo == nil || replyTo.MXID == "" {
		return "", "", nil, false
	}

	log := zerolog.Ctx(ctx)

	evt, err := tc.connector.br.Bot.GetEvent(ctx, portal.MXID, replyTo.MXID)
	if err != nil {
		log.Debug().
			Err(err).
			Stringer("event_id", replyTo.MXID).
			Msg("Failed to fetch Matrix event for reply, will send without reply metadata")
		return "", "", nil, false
	}

	// Parse message content (skip if already parsed by GetEvent)
	if evt.Content.Parsed == nil {
		if err := evt.Content.ParseRaw(event.EventMessage); err != nil {
			log.Debug().
				Err(err).
				Msg("Failed to parse Matrix event content for reply")
			return "", "", nil, false
		}
	}
	content := evt.Content.AsMessage()
	if content == nil {
		return "", "", nil, false
	}

	// Get message text
	messageText := content.Body

	// Get sender display name
	var senderDisplayName string
	// Try to get from ghost
	if networkUserID, ok := tc.connector.br.Matrix.ParseGhostMXID(evt.Sender); ok {
		ghost, err := tc.connector.br.GetGhostByID(ctx, networkUserID)
		if err == nil && ghost != nil && ghost.Name != "" {
			senderDisplayName = ghost.Name
		}
	}
	// Fallback to network user lookup
	if senderDisplayName == "" && replyTo.SenderID != "" {
		senderDisplayName = tc.getDisplayNameForUser(ctx, string(replyTo.SenderID))
	}

	// Build attachments from Matrix media (will extract MediaHashKey when direct media is enabled)
	if content.MsgType.IsMedia() {
		attachments = tc.buildReplyAttachmentsFromMatrixContent(ctx, portal, content)
	}

	return messageText, senderDisplayName, attachments, true
}

// buildReplyAttachmentsFromMatrixContent builds Twitter reply attachments from Matrix message content.
// When direct media is enabled and the content URL contains an encrypted media ID, extracts the MediaHashKey.
// Otherwise, just includes the media type for the reply preview.
func (tc *TwitterClient) buildReplyAttachmentsFromMatrixContent(ctx context.Context, portal *bridgev2.Portal, content *event.MessageEventContent) []*payload.MessageAttachment {
	if content == nil || content.URL == "" {
		return nil
	}

	// Determine media type
	var mediaType payload.MediaType
	switch content.MsgType {
	case event.MsgImage:
		mediaType = payload.MediaTypeImage
	case event.MsgVideo:
		mediaType = payload.MediaTypeVideo
	case event.MsgAudio:
		mediaType = payload.MediaTypeAudio
	default:
		return nil
	}

	mediaTypeInt := int32(mediaType)
	att := &payload.MessageAttachment{
		Media: &payload.MediaAttachment{
			Type: &mediaTypeInt,
		},
	}

	// Try to extract MediaHashKey from direct media URL if available
	if tc.connector.directMedia {
		if mediaHashKey := tc.extractMediaHashKeyFromContentURL(content.URL); mediaHashKey != "" {
			att.Media.MediaHashKey = &mediaHashKey
		}
	}

	// Include dimensions if available
	if content.Info != nil {
		if content.Info.Width > 0 || content.Info.Height > 0 {
			width := int64(content.Info.Width)
			height := int64(content.Info.Height)
			att.Media.Dimensions = &payload.MediaDimensions{
				Width:  &width,
				Height: &height,
			}
		}
		if content.Info.Size > 0 {
			size := int64(content.Info.Size)
			att.Media.FilesizeBytes = &size
		}
	}

	return []*payload.MessageAttachment{att}
}

// extractMediaHashKeyFromContentURL extracts the MediaHashKey from an encrypted media URL.
// Returns empty string if the URL is not an encrypted media URL or extraction fails.
func (tc *TwitterClient) extractMediaHashKeyFromContentURL(contentURL id.ContentURIString) string {
	// Parse the mxc:// URL
	uri, err := contentURL.Parse()
	if err != nil {
		return ""
	}

	// The FileID is base64url-encoded (without padding) mediaID bytes
	mediaIDBytes, err := base64.RawURLEncoding.DecodeString(uri.FileID)
	if err != nil {
		return ""
	}

	// Parse the media ID to check if it's encrypted media (version 2)
	parsed, err := ParseMediaID(networkid.MediaID(mediaIDBytes))
	if err != nil {
		return ""
	}

	if encInfo, ok := parsed.(*EncryptedMediaInfo); ok {
		return encInfo.MediaHashKey
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

	// Check portal metadata for trust status
	meta := ensurePortalMetadata(msg.Portal)
	if !meta.IsTrusted() {
		// Untrusted - only use REST API
		params := &payload.MarkConversationReadQuery{
			ConversationID:  conversationID,
			LastReadEventID: lastReadEventID,
		}
		return tc.client.MarkConversationRead(ctx, params)
	}

	// Trusted - use XChat, with REST fallback on key not found
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

func (tc *TwitterClient) HandleMatrixMembership(ctx context.Context, msg *bridgev2.MatrixMembershipChange) (*bridgev2.MatrixMembershipResult, error) {
	if msg.Type != bridgev2.Invite {
		return nil, errors.New("unsupported membership change type")
	}
	if msg.Portal.RoomType == database.RoomTypeDM {
		return nil, errors.New("cannot change members for DM")
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
		return nil, err
	}
	return &bridgev2.MatrixMembershipResult{}, nil
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

func (tc *TwitterClient) HandleMatrixAcceptMessageRequest(ctx context.Context, msg *bridgev2.MatrixAcceptMessageRequest) error {
	return tc.client.AcceptConversation(ctx, string(msg.Portal.ID))
}
