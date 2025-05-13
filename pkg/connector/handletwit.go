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
	"maps"
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (tc *TwitterClient) HandleTwitterEvent(rawEvt types.TwitterEvent, inbox *response.TwitterInboxData) {
	if rawEvt == nil {
		tc.updateTwitterUserInfo(inbox)
		tc.updateTwitterReadReceipt(inbox)
		tc.userCacheLock.Lock()
		maps.Copy(tc.userCache, inbox.Users)
		tc.userCacheLock.Unlock()
		return
	}
	isEdit := false
	if edit, ok := rawEvt.(*types.MessageEdit); ok {
		rawEvt = (*types.Message)(edit)
		isEdit = true
	}
	switch evt := rawEvt.(type) {
	case *types.Message:
		isFromMe := evt.MessageData.SenderID == string(tc.userLogin.ID)
		msgType := bridgev2.RemoteEventMessage
		if isEdit {
			msgType = bridgev2.RemoteEventEdit
		}
		conversation := inbox.GetConversationByID(evt.ConversationID)
		tc.connector.br.QueueRemoteEvent(tc.userLogin, &simplevent.Message[*types.MessageData]{
			EventMeta: simplevent.EventMeta{
				Type: msgType,
				LogContext: func(c zerolog.Context) zerolog.Context {
					if conversation != nil {
						c = c.
							Bool("conv_low_quality", conversation.LowQuality).
							Bool("conv_trusted", conversation.Trusted)
					} else {
						c = c.Bool("conversation_nil", true)
					}
					return c.
						Str("message_id", evt.MessageData.ID).
						Str("sender", evt.MessageData.SenderID).
						Bool("is_from_me", isFromMe)
				},
				PortalKey:    tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				CreatePortal: isFromMe || (conversation != nil && !conversation.LowQuality && conversation.Trusted),
				Sender:       tc.MakeEventSender(evt.MessageData.SenderID),
				StreamOrder:  methods.ParseSnowflakeInt(evt.MessageData.ID),
				Timestamp:    methods.ParseSnowflake(evt.MessageData.ID),
			},
			ID:            networkid.MessageID(evt.MessageData.ID),
			TransactionID: networkid.TransactionID(evt.RequestID),
			TargetMessage: networkid.MessageID(evt.MessageData.ID),
			Data:          &evt.MessageData,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
				return tc.convertToMatrix(ctx, portal, intent, data), nil
			},
			ConvertEditFunc: tc.convertEditToMatrix,
		})
	case *types.ConversationRead:
		tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventReadReceipt,
				PortalKey: tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				Sender:    bridgev2.EventSender{IsFromMe: true},
				Timestamp: methods.ParseSnowflake(evt.ID),
			},
			LastTarget: networkid.MessageID(evt.LastReadEventID),
		})
	case *types.MessageReactionCreate:
		portalKey := tc.makePortalKeyFromInbox(evt.ConversationID, inbox)
		wrappedEvt := tc.wrapReaction((*types.MessageReaction)(evt), portalKey, bridgev2.RemoteEventReaction)
		tc.connector.br.QueueRemoteEvent(tc.userLogin, wrappedEvt)
	case *types.MessageReactionDelete:
		portalKey := tc.makePortalKeyFromInbox(evt.ConversationID, inbox)
		wrappedEvt := tc.wrapReaction((*types.MessageReaction)(evt), portalKey, bridgev2.RemoteEventReactionRemove)
		tc.connector.br.QueueRemoteEvent(tc.userLogin, wrappedEvt)
	case *types.ConversationCreate:
		// honestly not sure when this is ever called... ? might be when they initialize the conversation with me?
		tc.client.Logger.Warn().Any("data", evt).Msg("Unhandled conversation create event")
	case *types.MessageDelete:
		for _, deletedMsg := range evt.Messages {
			messageDeleteRemoteEvent := &simplevent.MessageRemove{
				EventMeta: simplevent.EventMeta{
					Type:      bridgev2.RemoteEventMessageRemove,
					PortalKey: tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("message_id", deletedMsg.MessageID).
							Str("message_create_event_id", deletedMsg.MessageCreateEventID)
					},
					Timestamp:   methods.ParseSnowflake(evt.ID),
					StreamOrder: methods.ParseSnowflakeInt(evt.ID),
				},
				TargetMessage: networkid.MessageID(deletedMsg.MessageID),
			}
			tc.connector.br.QueueRemoteEvent(tc.userLogin, messageDeleteRemoteEvent)
		}
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
				PortalKey:   tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				Timestamp:   methods.ParseSnowflake(evt.ID),
				StreamOrder: methods.ParseSnowflakeInt(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{
					Name: &evt.ConversationName,
				},
			},
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, portalUpdateRemoteEvent)
	case *types.ConversationMetadataUpdate:
		tc.client.Logger.Warn().Any("data", evt).Msg("Unhandled conversation metadata update event")
	case *types.ParticipantsJoin:
		conversation := inbox.GetConversationByID(evt.ConversationID)
		portalMembersAddedRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				LogContext: func(c zerolog.Context) zerolog.Context {
					if conversation != nil {
						c = c.
							Bool("conv_low_quality", conversation.LowQuality).
							Bool("conv_trusted", conversation.Trusted)
					} else {
						c = c.Bool("conversation_nil", true)
					}
					return c.
						Str("conversation_id", evt.ConversationID).
						Int("total_new_members", len(evt.Participants))
				},
				PortalKey:    tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				CreatePortal: conversation != nil && !conversation.LowQuality && conversation.Trusted,
				StreamOrder:  methods.ParseSnowflakeInt(evt.ID),
				Timestamp:    methods.ParseSnowflake(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				MemberChanges: tc.participantsToMemberList(evt.Participants, inbox),
			},
		}
		portalMembersAddedRemoteEvent.ChatInfoChange.MemberChanges.IsFull = false
		portalMembersAddedRemoteEvent.ChatInfoChange.MemberChanges.TotalMemberCount = len(conversation.Participants)
		tc.connector.br.QueueRemoteEvent(tc.userLogin, portalMembersAddedRemoteEvent)
	case *types.ConversationDelete:
		portalDeleteRemoteEvent := &simplevent.ChatDelete{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventChatDelete,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID)
				},
				StreamOrder: methods.ParseSnowflakeInt(evt.ID),
				Timestamp:   methods.ParseSnowflake(evt.ID),
			},
			OnlyForMe: true,
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, portalDeleteRemoteEvent)
		tc.client.Logger.Info().Any("data", evt).Msg("Deleted conversation")
	default:
		tc.client.Logger.Warn().
			Type("event_data_type", rawEvt).
			Any("event_data", rawEvt).
			Msg("Received unhandled event case from twitter library")
	}
}

func (tc *TwitterClient) wrapReaction(data *types.MessageReaction, portalKey networkid.PortalKey, evtType bridgev2.RemoteEventType) *simplevent.Reaction {
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
			Sender:      tc.MakeEventSender(data.SenderID),
			Timestamp:   methods.ParseSnowflake(data.ID),
			StreamOrder: methods.ParseSnowflakeInt(data.ID),
		},
		EmojiID:       "",
		Emoji:         data.EmojiReaction,
		TargetMessage: networkid.MessageID(data.MessageID),
	}
}

func (tc *TwitterClient) updateTwitterReadReceipt(inbox *response.TwitterInboxData) {
	for conversationID, conversation := range inbox.Conversations {
		cache := tc.participantCache[conversationID]
		for _, participant := range conversation.Participants {
			if participant.UserID == ParseUserLoginID(tc.userLogin.ID) {
				continue
			}
			var cachedParticipant *types.Participant
			for _, p := range cache {
				if p.UserID == participant.UserID {
					cachedParticipant = &p
					break
				}
			}
			if cachedParticipant == nil || cachedParticipant.LastReadEventID < participant.LastReadEventID {
				tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
					EventMeta: simplevent.EventMeta{
						Type:      bridgev2.RemoteEventReadReceipt,
						PortalKey: tc.makePortalKeyFromInbox(conversationID, inbox),
						Sender:    tc.MakeEventSender(participant.UserID),
						Timestamp: methods.ParseSnowflake(participant.LastReadEventID),
					},
					LastTarget: networkid.MessageID(participant.LastReadEventID),
				})
			}
		}
		tc.participantCache[conversationID] = conversation.Participants
	}
}

func (tc *TwitterClient) updateTwitterUserInfo(inbox *response.TwitterInboxData) {
	ctx := tc.userLogin.Log.With().Str("action", "update user info").Logger().WithContext(context.Background())
	for userID, user := range inbox.Users {
		cached := tc.userCache[userID]
		if cached == nil || cached.Name != user.Name || cached.ScreenName != user.ScreenName || cached.ProfileImageURLHTTPS != user.ProfileImageURLHTTPS {
			ghost, err := tc.connector.br.GetGhostByID(ctx, MakeUserID(userID))
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("Failed to get ghost by ID")
			} else {
				ghost.UpdateInfo(ctx, tc.connector.wrapUserInfo(tc.client, user))
			}
		}
	}
}

func (tc *TwitterClient) HandleStreamEvent(evt response.StreamEvent) {
	updateData := evt.Payload.DmUpdate
	typingData := evt.Payload.DmTyping

	if updateData != nil {
		tc.client.PollConversation(updateData.ConversationID)
	}

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
