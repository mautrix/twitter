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
	"go.mau.fi/util/variationselector"
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
	reactionKey := data.EmojiReaction

	reactionKey = variationselector.FullyQualify(reactionKey)
	emojiID := networkid.EmojiID(reactionKey)

	tc.userLogin.Log.Debug().
		Str("portal_id", string(portalKey.ID)).
		Str("portal_receiver", string(portalKey.Receiver)).
		Str("target_message_id", string(MakeMessageID(data.MessageID))).
		Str("sender_id", string(MakeUserID(senderID))).
		Str("emoji_id", string(emojiID)).
		Str("emoji", reactionKey).
		Msg("Wrapping reaction for remote event")

	return &simplevent.Reaction{
		EventMeta: simplevent.EventMeta{
			Type: evtType,
			LogContext: func(c zerolog.Context) zerolog.Context {
				return c.
					Str("message_id", data.MessageID).
					Str("sender", data.SenderID).
					Str("emoji_reaction", data.EmojiReaction)
			},
			PortalKey:   portalKey,
			Sender:      tc.MakeEventSender(senderID),
			Timestamp:   methods.ParseSnowflake(data.ID),
			StreamOrder: methods.ParseSnowflakeInt(data.ID),
		},
		EmojiID:       emojiID,
		Emoji:         reactionKey,
		TargetMessage: MakeMessageID(data.MessageID),
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
		memberChanges.MemberMap.Set(bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(p.UserID),
			Membership:  membership,
		})
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
		isFromMe := MakeUserLoginID(evt.MessageData.SenderID) == tc.userLogin.ID
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
			ID:            MakeMessageID(targetMessageID),
			TransactionID: networkid.TransactionID(txnID),
			TargetMessage: MakeMessageID(targetMessageID),
			Data:          &evt.MessageData,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
				return tc.convertToMatrix(ctx, portal, intent, data), nil
			},
			ConvertEditFunc: tc.convertEditToMatrix,
		}).Success

	case *types.Message:
		isFromMe := MakeUserLoginID(evt.MessageData.SenderID) == tc.userLogin.ID
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
			ID:            MakeMessageID(msgID),
			TransactionID: networkid.TransactionID(txnID),
			TargetMessage: MakeMessageID(msgID),
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
					if clientMsgID != "" && meta.XChatClientMsgID == "" {
						meta.XChatClientMsgID = clientMsgID
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
		lastTarget := MakeMessageID(evt.LastReadEventID)
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
				Sender:    tc.MakeEventSender(ParseUserLoginID(tc.userLogin.ID)),
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
				TargetMessage: MakeMessageID(deletedMsg.MessageID),
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
		if ctx == nil {
			ctx = context.Background()
		}

		// XChat group titles are encrypted. Decrypt before forwarding to Matrix so
		// we don't set the room name to ciphertext.
		newName := evt.ConversationName
		if decrypted := tc.decryptGroupName(ctx, evt.ConversationID, newName); decrypted != "" {
			newName = decrypted
		} else if isProbablyEncryptedGroupName(newName) {
			log.Debug().
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to decrypt XChat group title change, skipping name update")
			return true
		}

		portalUpdateRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:   bridgev2.RemoteEventChatInfoChange,
				Sender: tc.MakeEventSender(evt.ByUserID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Str("new_name", newName).
						Str("changed_by_user_id", evt.ByUserID)
				},
				PortalKey:   tc.MakePortalKeyFromID(evt.ConversationID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
				StreamOrder: methods.ParseInt64(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{
					Name: &newName,
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

	case *types.ConversationCreate:
		// Conversation created or became active - ensure portal exists and trigger backfill.
		// This event is emitted when we receive an empty MessageCreateEvent, which often
		// happens when a message request is accepted.
		log.Info().
			Str("conversation_id", evt.ConversationID).
			Str("event_id", evt.ID).
			Msg("Conversation create event received, ensuring portal exists")

		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		portal, err := tc.ensurePortalForConversation(ctx, evt.ConversationID, "")
		if err != nil {
			log.Warn().Err(err).
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to ensure portal for ConversationCreate event")
			return false
		}

		// If the portal was just created or doesn't have a room yet, trigger a resync
		// to ensure it gets created and backfilled
		if portal.MXID == "" {
			chatInfo := tc.getTrustedChatInfo(ctx, evt.ConversationID)
			if chatInfo == nil {
				log.Warn().
					Str("conversation_id", evt.ConversationID).
					Msg("Failed to get chat info for ConversationCreate event")
				return false
			}

			return tc.userLogin.QueueRemoteEvent(&simplevent.ChatResync{
				EventMeta: simplevent.EventMeta{
					Type:         bridgev2.RemoteEventChatResync,
					PortalKey:    portalKey,
					CreatePortal: true,
					Timestamp:    methods.ParseMsecTimestamp(evt.Time),
					StreamOrder:  methods.ParseInt64(evt.ID),
				},
				ChatInfo: chatInfo,
				CheckNeedsBackfillFunc: func(ctx context.Context, latestMessage *database.Message) (bool, error) {
					return true, nil
				},
			}).Success
		}
		return true

	case *types.TrustConversation:
		// Conversation was accepted (became trusted) - emit ChatResync to update MessageRequest status
		log.Info().
			Str("conversation_id", evt.ConversationID).
			Str("reason", evt.Reason).
			Msg("Conversation became trusted (message request accepted)")

		// Update portal metadata to mark as trusted
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
		if err != nil {
			log.Warn().Err(err).
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to get portal for TrustConversation event")
		} else {
			meta := portal.Metadata.(*PortalMetadata)
			if !meta.Trusted {
				meta.Trusted = true
				if err := portal.Save(ctx); err != nil {
					log.Warn().Err(err).
						Str("conversation_id", evt.ConversationID).
						Msg("Failed to save portal metadata with Trusted=true")
				}
			}
		}

		// Fetch updated conversation data and create ChatInfo with MessageRequest: false
		chatInfo := tc.getTrustedChatInfo(ctx, evt.ConversationID)
		if chatInfo == nil {
			log.Warn().
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to get chat info for trusted conversation")
			return false
		}

		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatResync{
			EventMeta: simplevent.EventMeta{
				Type:         bridgev2.RemoteEventChatResync,
				PortalKey:    tc.MakePortalKeyFromID(evt.ConversationID),
				CreatePortal: true,
				Timestamp:    methods.ParseMsecTimestamp(evt.Time),
				StreamOrder:  methods.ParseInt64(evt.ID),
			},
			ChatInfo: chatInfo,
		}).Success

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

// HandlePollingEvent handles events from REST API polling.
// This is used for untrusted (message request) conversations that don't
// receive real-time updates via XChat WebSocket.
// Returns true to continue polling, false to stop.
func (tc *TwitterClient) HandlePollingEvent(evt types.TwitterEvent, inbox *response.TwitterInboxData) bool {
	// Always cache users from inbox when available - needed for portal creation
	if inbox != nil {
		tc.updateTwitterUserInfo(context.Background(), inbox)
		tc.userCacheLock.Lock()
		for userID, user := range inbox.Users {
			tc.userCache[userID] = user
		}
		tc.userCacheLock.Unlock()
	}

	if evt == nil {
		// nil event with inbox means initial poll - nothing more to do
		return true
	}

	log := tc.userLogin.Log.With().
		Str("handler", "polling").
		Type("event_type", evt).
		Logger()

	// Get conversation ID from the event
	var conversationID string
	switch e := evt.(type) {
	case *types.Message:
		conversationID = e.ConversationID
	case *types.MessageReactionCreate:
		conversationID = e.ConversationID
	case *types.MessageReactionDelete:
		conversationID = e.ConversationID
	case *types.ConversationRead:
		conversationID = e.ConversationID
	case *types.ConversationDelete:
		conversationID = e.ConversationID
	case *types.TrustConversation:
		conversationID = e.ConversationID
	case *types.PollingError:
		// Handle polling errors
		if e.IsAuth {
			log.Warn().Err(e.Error).Msg("Authentication error during polling")
			return false
		}
		if e.Error != nil {
			log.Warn().Err(e.Error).Msg("Polling error")
		}
		return true
	default:
		// For other event types, let them through
		log.Debug().Msg("Received unhandled polling event type")
		return true
	}

	if conversationID == "" {
		return true
	}

	// Skip if conversation can use XChat (has encryption keys) - XChat WebSocket handles those
	portalKey := tc.MakePortalKeyFromID(conversationID)
	ctx := context.Background()
	if portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey); err == nil && portal != nil {
		meta := portal.Metadata.(*PortalMetadata)
		if meta.CanUseXChat() {
			log.Debug().
				Str("conversation_id", conversationID).
				Msg("Skipping conversation event (handled by XChat WebSocket)")
			return true
		}
	}

	// Dispatch to the appropriate handler based on event type
	switch e := evt.(type) {
	case *types.Message:
		return tc.handlePollingMessage(e, inbox)
	case *types.MessageReactionCreate:
		reaction := (*types.MessageReaction)(e)
		portalKey := tc.MakePortalKeyFromID(conversationID)
		wrappedEvt := tc.wrapReaction(reaction, portalKey, bridgev2.RemoteEventReaction)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success
	case *types.MessageReactionDelete:
		reaction := (*types.MessageReaction)(e)
		portalKey := tc.MakePortalKeyFromID(conversationID)
		wrappedEvt := tc.wrapReaction(reaction, portalKey, bridgev2.RemoteEventReactionRemove)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success
	case *types.ConversationRead:
		lastTarget := MakeMessageID(e.LastReadEventID)
		readUpTo := methods.ParseMsecTimestamp(e.Time)
		readUpToStreamOrder := methods.ParseInt64(e.LastReadEventID)
		if readUpToStreamOrder == 0 {
			readUpToStreamOrder = methods.ParseInt64(e.ID)
		}
		var targets []networkid.MessageID
		if lastTarget != "" {
			targets = []networkid.MessageID{lastTarget}
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventReadReceipt,
				PortalKey: tc.MakePortalKeyFromID(conversationID),
				Sender:    tc.MakeEventSender(ParseUserLoginID(tc.userLogin.ID)),
				Timestamp: readUpTo,
			},
			LastTarget:          lastTarget,
			Targets:             targets,
			ReadUpTo:            readUpTo,
			ReadUpToStreamOrder: readUpToStreamOrder,
		}).Success
	case *types.ConversationDelete:
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatDelete{
			EventMeta: simplevent.EventMeta{
				Type:        bridgev2.RemoteEventChatDelete,
				PortalKey:   tc.MakePortalKeyFromID(conversationID),
				StreamOrder: methods.ParseInt64(e.ID),
				Timestamp:   methods.ParseMsecTimestamp(e.Time),
			},
			OnlyForMe: true,
		}).Success
	case *types.TrustConversation:
		// Conversation became trusted - update the portal
		log.Info().
			Str("conversation_id", conversationID).
			Str("reason", e.Reason).
			Msg("Conversation became trusted via polling")

		// Update portal metadata to mark as trusted (same as XChat path)
		ctx := context.Background()
		portalKey := tc.MakePortalKeyFromID(conversationID)
		portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
		if err != nil {
			log.Warn().Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to get portal for TrustConversation event")
		} else {
			meta := portal.Metadata.(*PortalMetadata)
			if !meta.Trusted {
				meta.Trusted = true
				if err := portal.Save(ctx); err != nil {
					log.Warn().Err(err).
						Str("conversation_id", conversationID).
						Msg("Failed to save portal metadata with Trusted=true")
				}
			}
		}

		chatInfo := tc.getTrustedChatInfo(ctx, conversationID)
		if chatInfo == nil {
			return false
		}

		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatResync{
			EventMeta: simplevent.EventMeta{
				Type:         bridgev2.RemoteEventChatResync,
				PortalKey:    tc.MakePortalKeyFromID(conversationID),
				CreatePortal: true,
				Timestamp:    methods.ParseMsecTimestamp(e.Time),
				StreamOrder:  methods.ParseInt64(e.ID),
			},
			ChatInfo: chatInfo,
		}).Success
	}

	return true
}

// handlePollingMessage handles a message event from REST API polling.
func (tc *TwitterClient) handlePollingMessage(evt *types.Message, inbox *response.TwitterInboxData) bool {
	isFromMe := MakeUserLoginID(evt.MessageData.SenderID) == tc.userLogin.ID
	portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
	msgID := evt.ID

	// For polling messages, ensure the portal exists
	ctx := context.Background()
	portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
	if err != nil {
		tc.userLogin.Log.Warn().
			Err(err).
			Str("conversation_id", evt.ConversationID).
			Msg("Failed to get portal for polling message")
		return false
	}

	// Create portal if it doesn't exist
	if portal.MXID == "" {
		chatInfo := tc.getOrFetchChatInfoForPolling(ctx, evt.ConversationID, inbox)
		if chatInfo == nil {
			tc.userLogin.Log.Warn().
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to get chat info for polling message")
			return false
		}
		if err := portal.CreateMatrixRoom(ctx, tc.userLogin, chatInfo); err != nil {
			tc.userLogin.Log.Warn().
				Err(err).
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to create Matrix room for polling message")
			return false
		}
		// Register backfill task for the newly created room
		if err := tc.connector.br.DB.BackfillTask.EnsureExists(ctx, portal.PortalKey, tc.userLogin.ID); err != nil {
			tc.userLogin.Log.Warn().Err(err).
				Str("conversation_id", evt.ConversationID).
				Msg("Failed to ensure backfill task exists for new polling room")
		} else {
			tc.connector.br.WakeupBackfillQueue()
		}
	}

	return tc.userLogin.QueueRemoteEvent(&simplevent.Message[*types.MessageData]{
		EventMeta: simplevent.EventMeta{
			Type: bridgev2.RemoteEventMessage,
			LogContext: func(c zerolog.Context) zerolog.Context {
				return c.
					Str("message_id", msgID).
					Str("sender", evt.MessageData.SenderID).
					Bool("is_from_me", isFromMe).
					Str("handler", "polling")
			},
			PortalKey:    portalKey,
			CreatePortal: true,
			Sender:       tc.MakeEventSender(evt.MessageData.SenderID),
			StreamOrder:  methods.ParseInt64(msgID),
			Timestamp:    methods.ParseMsecTimestamp(evt.Time),
		},
		ID:            MakeMessageID(msgID),
		TransactionID: networkid.TransactionID(evt.RequestID),
		TargetMessage: MakeMessageID(msgID),
		Data:          &evt.MessageData,
		ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
			return tc.convertToMatrix(ctx, portal, intent, data), nil
		},
		ConvertEditFunc: tc.convertEditToMatrix,
	}).Success
}

// updateTwitterUserInfo updates ghost info when user data changes.
// This ensures profile pictures are visible for users in message request conversations.
func (tc *TwitterClient) updateTwitterUserInfo(ctx context.Context, inbox *response.TwitterInboxData) {
	if inbox == nil || inbox.Users == nil {
		return
	}
	log := zerolog.Ctx(ctx)
	tc.userCacheLock.RLock()
	defer tc.userCacheLock.RUnlock()

	for userID, user := range inbox.Users {
		cached := tc.userCache[userID]
		if cached == nil || cached.Name != user.Name || cached.ScreenName != user.ScreenName || cached.ProfileImageURLHTTPS != user.ProfileImageURLHTTPS {
			ghost, err := tc.connector.br.GetGhostByID(ctx, MakeUserID(userID))
			if err != nil {
				log.Debug().Err(err).Str("user_id", userID).Msg("Failed to get ghost by ID for user info update")
				continue
			}
			ghost.UpdateInfo(ctx, tc.connector.wrapUserInfo(tc.client, user))
		}
	}
}
