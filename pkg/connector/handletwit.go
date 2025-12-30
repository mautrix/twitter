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
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (tc *TwitterClient) wrapReaction(data *types.MessageReaction, portalKey networkid.PortalKey, evtType bridgev2.RemoteEventType) *simplevent.Reaction {
	senderID := data.SenderID
	if senderID == "" {
		senderID = tc.client.GetCurrentUserID()
	}
	reactionKey := data.ReactionKey
	if reactionKey == "" {
		reactionKey = data.EmojiReaction
	}
	emojiID := networkid.EmojiID(reactionKey)
	if emojiID == "" {
		emojiID = networkid.EmojiID(data.EmojiReaction)
	}
	return &simplevent.Reaction{
		EventMeta: simplevent.EventMeta{
			Type: evtType,
			LogContext: func(c zerolog.Context) zerolog.Context {
				return c.
					Str("message_id", data.MessageID).
					Str("sender", data.SenderID).
					Str("reaction_key", data.ReactionKey).
					Str("emoji_reaction", data.EmojiReaction)
			},
			PortalKey:   portalKey,
			Sender:      tc.MakeEventSender(senderID),
			Timestamp:   methods.ParseSnowflake(data.ID),
			StreamOrder: methods.ParseSnowflakeInt(data.ID),
		},
		EmojiID:       emojiID,
		Emoji:         data.EmojiReaction,
		TargetMessage: networkid.MessageID(data.MessageID),
	}
}

func (tc *TwitterClient) HandleStreamEvent(evt response.StreamEvent) {
	typingData := evt.Payload.DmTyping

	if typingData != nil {
		tc.userLogin.QueueRemoteEvent(&simplevent.Typing{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventTyping,
				PortalKey: tc.MakePortalKeyFromID(typingData.ConversationID),
				Sender:    tc.MakeEventSender(typingData.UserID),
			},
			Timeout: 3 * time.Second,
		})
	}
}

// buildMemberChangeEvent creates a ChatInfoChange event for participant joins/leaves.
func (tc *TwitterClient) buildMemberChangeEvent(
	conversationID, eventID, eventTime string,
	participants []types.Participant,
	membership event.Membership,
) *simplevent.ChatInfoChange {
	memberChanges := &bridgev2.ChatMemberList{
		IsFull:    false,
		MemberMap: make(map[networkid.UserID]bridgev2.ChatMember, len(participants)),
	}
	for _, p := range participants {
		memberChanges.MemberMap[networkid.UserID(p.UserID)] = bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(p.UserID),
			Membership:  membership,
		}
	}
	return &simplevent.ChatInfoChange{
		EventMeta: simplevent.EventMeta{
			Type: bridgev2.RemoteEventChatInfoChange,
			LogContext: func(c zerolog.Context) zerolog.Context {
				return c.
					Str("conversation_id", conversationID).
					Int("member_count", len(participants))
			},
			PortalKey:    tc.MakePortalKeyFromID(conversationID),
			CreatePortal: false,
			StreamOrder:  methods.ParseInt64(eventID),
			Timestamp:    methods.ParseMsecTimestamp(eventTime),
		},
		ChatInfoChange: &bridgev2.ChatInfoChange{
			MemberChanges: memberChanges,
		},
	}
}

// HandleXChatEvent handles events from the XChat websocket processor.
func (tc *TwitterClient) HandleXChatEvent(ctx context.Context, rawEvt types.TwitterEvent) bool {
	if rawEvt == nil {
		return true
	}

	log := tc.userLogin.Log.With().
		Str("handler", "xchat").
		Type("event_type", rawEvt).
		Logger()

	switch evt := rawEvt.(type) {
	case *types.MessageEdit:
		isFromMe := evt.MessageData.SenderID == string(tc.userLogin.ID)
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		requiredKeyVersion := evt.ConversationKeyVersion
		eventID := evt.SequenceID
		if eventID == "" {
			eventID = evt.ID
		}
		streamOrder := methods.ParseInt64(eventID)
		if streamOrder == 0 {
			streamOrder = methods.ParseInt64(evt.Time)
		}
		targetMessageID := evt.MessageData.ID
		if targetMessageID == "" {
			targetMessageID = eventID
		}

		if ctx == nil || ctx.Value(ensurePortalContextKey{}) == nil {
			if _, err := tc.ensurePortalForConversation(ctx, evt.ConversationID, requiredKeyVersion); err != nil {
				log.Warn().
					Err(err).
					Str("conversation_id", evt.ConversationID).
					Str("required_key_version", requiredKeyVersion).
					Msg("Failed to ensure portal and key exist before handling XChat edit")
				return false
			}
		}

		txnID := evt.RequestID
		if txnID == "" {
			txnID = eventID
		}

		return tc.userLogin.QueueRemoteEvent(&simplevent.Message[*types.MessageData]{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventEdit,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("message_id", targetMessageID).
						Str("sender", evt.MessageData.SenderID).
						Bool("is_from_me", isFromMe)
				},
				PortalKey:    portalKey,
				CreatePortal: false,
				Sender:       tc.MakeEventSender(evt.MessageData.SenderID),
				StreamOrder:  streamOrder,
				Timestamp:    methods.ParseMsecTimestamp(evt.Time),
			},
			ID:            networkid.MessageID(targetMessageID),
			TransactionID: networkid.TransactionID(txnID),
			TargetMessage: networkid.MessageID(targetMessageID),
			Data:          &evt.MessageData,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
				return tc.convertToMatrix(ctx, portal, intent, data), nil
			},
			ConvertEditFunc: tc.convertEditToMatrix,
		}).Success

	case *types.Message:
		isFromMe := evt.MessageData.SenderID == string(tc.userLogin.ID)
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		requiredKeyVersion := evt.ConversationKeyVersion
		msgID := evt.SequenceID
		if msgID == "" {
			msgID = evt.ID
		}
		streamOrder := methods.ParseInt64(msgID)
		if streamOrder == 0 {
			streamOrder = methods.ParseInt64(evt.Time)
		}

		if ctx == nil || ctx.Value(ensurePortalContextKey{}) == nil {
			if _, err := tc.ensurePortalForConversation(ctx, evt.ConversationID, requiredKeyVersion); err != nil {
				log.Warn().
					Err(err).
					Str("conversation_id", evt.ConversationID).
					Str("required_key_version", requiredKeyVersion).
					Msg("Failed to ensure portal and key exist before handling XChat message")
				return false
			}
		}

		txnID := evt.RequestID
		if txnID == "" {
			txnID = msgID
		}
		clientMsgID := evt.RequestID

		return tc.userLogin.QueueRemoteEvent(&simplevent.Message[*types.MessageData]{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventMessage,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("message_id", msgID).
						Str("sender", evt.MessageData.SenderID).
						Bool("is_from_me", isFromMe)
				},
				PortalKey:    portalKey,
				CreatePortal: false, // Portal should already exist from initial sync
				Sender:       tc.MakeEventSender(evt.MessageData.SenderID),
				StreamOrder:  streamOrder,
				Timestamp:    methods.ParseMsecTimestamp(evt.Time),
			},
			ID:            networkid.MessageID(msgID),
			TransactionID: networkid.TransactionID(txnID),
			TargetMessage: networkid.MessageID(msgID),
			Data:          &evt.MessageData,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
				converted := tc.convertToMatrix(ctx, portal, intent, data)
				if converted == nil {
					return nil, nil
				}
				for _, part := range converted.Parts {
					if part == nil {
						continue
					}
					meta, ok := part.DBMetadata.(*MessageMetadata)
					if !ok || meta == nil {
						meta = &MessageMetadata{}
						part.DBMetadata = meta
					}
					if meta.XChatSequenceID == "" {
						meta.XChatSequenceID = msgID
					}
					if meta.XChatCreatedAtMS == "" {
						meta.XChatCreatedAtMS = evt.Time
					}
					if clientMsgID != "" && meta.XChatClientMsgID == "" {
						meta.XChatClientMsgID = clientMsgID
					}
				}
				return converted, nil
			},
			ConvertEditFunc: tc.convertEditToMatrix,
			HandleExistingFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, existing []*database.Message, data *types.MessageData) (bridgev2.UpsertResult, error) {
				for _, part := range existing {
					meta, ok := part.Metadata.(*MessageMetadata)
					if !ok || meta == nil {
						meta = &MessageMetadata{}
						part.Metadata = meta
					}
					if meta.XChatSequenceID == "" {
						meta.XChatSequenceID = msgID
					}
					if meta.XChatCreatedAtMS == "" {
						meta.XChatCreatedAtMS = evt.Time
					}
					if clientMsgID != "" && meta.XChatClientMsgID == "" {
						meta.XChatClientMsgID = clientMsgID
					}
					if data != nil {
						if meta.MessageText == "" {
							meta.MessageText = data.Text
						}
						if meta.SenderID == "" {
							meta.SenderID = data.SenderID
						}
						if meta.SenderDisplayName == "" {
							meta.SenderDisplayName = tc.getDisplayNameForUser(ctx, data.SenderID)
							if meta.SenderDisplayName == "" {
								meta.SenderDisplayName = data.SenderID
							}
						}
						if len(meta.ReplyAttachments) == 0 && len(data.OriginalAttachments) > 0 {
							meta.ReplyAttachments = filterReplyPreviewAttachments(data.OriginalAttachments)
						}
					}
					if data != nil && data.EditCount > meta.EditCount {
						meta.EditCount = data.EditCount
					}
				}
				// This is typically the remote echo for a pending outgoing message; don't bridge again.
				return bridgev2.UpsertResult{SaveParts: true, ContinueMessageHandling: false}, nil
			},
		}).Success

	case *types.MessageReactionCreate:
		reaction := (*types.MessageReaction)(evt)
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		wrappedEvt := tc.wrapReaction(reaction, portalKey, bridgev2.RemoteEventReaction)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success

	case *types.MessageReactionDelete:
		reaction := (*types.MessageReaction)(evt)
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		wrappedEvt := tc.wrapReaction(reaction, portalKey, bridgev2.RemoteEventReactionRemove)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success

	case *types.ConversationRead:
		lastTarget := networkid.MessageID(evt.LastReadEventID)
		readUpTo := methods.ParseMsecTimestamp(evt.Time)
		readUpToStreamOrder := methods.ParseInt64(evt.LastReadEventID)
		if readUpToStreamOrder == 0 {
			readUpToStreamOrder = methods.ParseInt64(evt.ID)
		}
		var targets []networkid.MessageID
		if lastTarget != "" {
			targets = []networkid.MessageID{lastTarget}
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventReadReceipt,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				Sender:    tc.MakeEventSender(string(tc.userLogin.ID)),
				Timestamp: readUpTo,
				PreHandleFunc: func(ctx context.Context, portal *bridgev2.Portal) {
					if intent := tc.userLogin.User.DoublePuppet(ctx); intent != nil {
						_ = intent.EnsureJoined(ctx, portal.MXID)
					}
				},
			},
			LastTarget:          lastTarget,
			Targets:             targets,
			ReadUpTo:            readUpTo,
			ReadUpToStreamOrder: readUpToStreamOrder,
		}).Success

	case *types.MessageDelete:
		allSuccess := true
		for _, deletedMsg := range evt.Messages {
			messageDeleteRemoteEvent := &simplevent.MessageRemove{
				EventMeta: simplevent.EventMeta{
					Type:      bridgev2.RemoteEventMessageRemove,
					PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("message_id", deletedMsg.MessageID).
							Str("message_create_event_id", deletedMsg.MessageCreateEventID)
					},
					Timestamp:   methods.ParseMsecTimestamp(evt.Time),
					StreamOrder: methods.ParseInt64(evt.ID),
				},
				TargetMessage: networkid.MessageID(deletedMsg.MessageID),
			}
			allSuccess = tc.userLogin.QueueRemoteEvent(messageDeleteRemoteEvent).Success && allSuccess
		}
		return allSuccess

	case *types.ConversationDelete:
		portalDeleteRemoteEvent := &simplevent.ChatDelete{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventChatDelete,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.Str("conversation_id", evt.ConversationID)
				},
				StreamOrder: methods.ParseInt64(evt.ID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
			},
			OnlyForMe: true,
		}
		log.Info().Any("data", evt).Msg("Deleted conversation")
		return tc.userLogin.QueueRemoteEvent(portalDeleteRemoteEvent).Success

	case *types.ConversationNameUpdate:
		portalUpdateRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:   bridgev2.RemoteEventChatInfoChange,
				Sender: tc.MakeEventSender(evt.ByUserID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Str("new_name", evt.ConversationName).
						Str("changed_by_user_id", evt.ByUserID)
				},
				PortalKey:   tc.MakePortalKeyFromID(evt.ConversationID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
				StreamOrder: methods.ParseInt64(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{
					Name: &evt.ConversationName,
				},
			},
		}
		return tc.userLogin.QueueRemoteEvent(portalUpdateRemoteEvent).Success

	case *types.ConversationAvatarUpdate:
		chatInfo := &bridgev2.ChatInfo{
			Avatar: tc.makeGroupAvatar(evt.ConversationID, evt.ConversationAvatarImageHttps, evt.ConversationKeyVersion),
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:        bridgev2.RemoteEventChatInfoChange,
				Sender:      tc.MakeEventSender(evt.ByUserID),
				PortalKey:   tc.MakePortalKeyFromID(evt.ConversationID),
				StreamOrder: methods.ParseInt64(evt.ID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: chatInfo,
			},
		}).Success

	case *types.ParticipantsJoin:
		changeEvt := tc.buildMemberChangeEvent(evt.ConversationID, evt.ID, evt.Time, evt.Participants, event.MembershipJoin)
		return tc.userLogin.QueueRemoteEvent(changeEvt).Success

	case *types.ParticipantsLeave:
		changeEvt := tc.buildMemberChangeEvent(evt.ConversationID, evt.ID, evt.Time, evt.Participants, event.MembershipLeave)
		return tc.userLogin.QueueRemoteEvent(changeEvt).Success

	case *types.XChatTyping:
		tc.userLogin.QueueRemoteEvent(&simplevent.Typing{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventTyping,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				Sender:    tc.MakeEventSender(evt.SenderID),
			},
			Timeout: 3 * time.Second,
		})
		return true

	case *types.XChatKeyChange:
		log.Info().
			Str("conversation_id", evt.ConversationID).
			Str("new_key_version", evt.NewKeyVersion).
			Msg("Conversation key changed")
		return true

	case *types.XChatMessageFailure:
		log.Warn().
			Str("conversation_id", evt.ConversationID).
			Str("message_id", evt.MessageID).
			Int32("failure_type", int32(evt.FailureType)).
			Msg("Message failure event received")
		return true

	default:
		log.Debug().
			Type("event_data_type", rawEvt).
			Any("event_data", rawEvt).
			Msg("Received unhandled XChat event")
		return true
	}
}

// ignorePayloadImport is used to prevent the import from being removed
var _ = payload.FailureType(0)
