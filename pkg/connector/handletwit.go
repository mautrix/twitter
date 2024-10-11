package connector

import (
	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
)

func (tc *TwitterClient) HandleTwitterEvent(rawEvt any) {
	switch evtData := rawEvt.(type) {
	case event.XEventMessage:
		sender := evtData.Sender
		isFromMe := sender.IDStr == string(tc.userLogin.ID)
		msgType := bridgev2.RemoteEventMessage

		if evtData.EditCount > 0 {
			msgType = bridgev2.RemoteEventEdit
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, &simplevent.Message[*event.XEventMessage]{
			EventMeta: simplevent.EventMeta{
				Type: msgType,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("message_id", evtData.MessageID).
						Str("sender", sender.IDStr).
						Str("sender_login", sender.ScreenName).
						Bool("is_from_me", isFromMe)
				},
				PortalKey:    tc.MakePortalKey(evtData.Conversation),
				CreatePortal: true,
				Sender: bridgev2.EventSender{
					IsFromMe:    isFromMe,
					SenderLogin: networkid.UserLoginID(sender.IDStr),
					Sender:      networkid.UserID(sender.IDStr),
				},
				Timestamp: evtData.CreatedAt,
			},
			ID:                 networkid.MessageID(evtData.MessageID),
			TargetMessage:      networkid.MessageID(evtData.MessageID),
			Data:               &evtData,
			ConvertMessageFunc: tc.convertToMatrix,
			ConvertEditFunc:    tc.convertEditToMatrix,
		})
	case event.XEventReaction:
		reactionRemoteEvent := tc.wrapReaction(evtData)
		tc.connector.br.QueueRemoteEvent(tc.userLogin, reactionRemoteEvent)
	case event.XEventConversationRead:
		// conversation read events are only fired by yourself??
		// if another user reads your message this is never fired
		// they use user_updates and last_read_event_id to figure that out with the polling client
		tc.client.Logger.Info().
			Str("conversation_id", evtData.Conversation.ConversationID).
			Str("last_read_event_id", evtData.LastReadEventID).
			Str("read_at", evtData.ReadAt.String()).
			Msg("Conversation was read!")
	case event.XEventConversationCreated:
		// honestly not sure when this is ever called... ? might be when they initialize the conversation with me?
		tc.client.Logger.Warn().Any("data", evtData).Msg("XEventConversationCreated")
	case event.XEventMessageDeleted:
		for _, deletedMsg := range evtData.Messages {
			messageDeleteRemoteEvent := &simplevent.MessageRemove{
				EventMeta: simplevent.EventMeta{
					Type:      bridgev2.RemoteEventMessageRemove,
					PortalKey: tc.MakePortalKey(evtData.Conversation),
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("message_id", deletedMsg.MessageID).
							Str("message_create_event_id", deletedMsg.MessageCreateEventID)
					},
					Timestamp: evtData.DeletedAt,
				},

				TargetMessage: networkid.MessageID(deletedMsg.MessageID),
			}
			tc.connector.br.QueueRemoteEvent(tc.userLogin, messageDeleteRemoteEvent)
		}
	case event.XEventConversationNameUpdate:
		portalUpdateRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				Sender: bridgev2.EventSender{
					IsFromMe:    evtData.Executor.IDStr == string(tc.userLogin.ID),
					SenderLogin: networkid.UserLoginID(evtData.Executor.IDStr),
					Sender:      networkid.UserID(evtData.Executor.IDStr),
				},
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evtData.Conversation.ConversationID).
						Str("new_name", evtData.Name).
						Str("changed_by_user_id", evtData.Executor.IDStr)
				},
				PortalKey: tc.MakePortalKey(evtData.Conversation),
				Timestamp: evtData.UpdatedAt,
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{
					Name: &evtData.Name,
				},
			},
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, portalUpdateRemoteEvent)
	case event.XEventParticipantsJoined:
		portalMembersAddedRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evtData.Conversation.ConversationID).
						Int("total_new_members", len(evtData.NewParticipants))
				},
				PortalKey: tc.MakePortalKey(evtData.Conversation),
				Timestamp: evtData.EventTime,
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				MemberChanges: tc.UsersToMemberList(evtData.NewParticipants),
			},
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, portalMembersAddedRemoteEvent)
	case event.XEventConversationDelete:
		portalDeleteRemoteEvent := &simplevent.ChatDelete{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventChatDelete,
				PortalKey: tc.MakePortalKeyFromID(evtData.ConversationID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evtData.ConversationID)
				},
				Timestamp: evtData.DeletedAt,
			},
			OnlyForMe: true,
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, portalDeleteRemoteEvent)
		tc.client.Logger.Info().Any("data", evtData).Msg("Deleted conversation")
	default:
		tc.client.Logger.Warn().Any("event_data", evtData).Msg("Received unhandled event case from twitter library")
	}
}

func (tc *TwitterClient) wrapReaction(data event.XEventReaction) *simplevent.Reaction {
	var eventType bridgev2.RemoteEventType
	if data.Action == types.MessageReactionRemove {
		eventType = bridgev2.RemoteEventReactionRemove
	} else {
		eventType = bridgev2.RemoteEventReaction
	}

	var receiver networkid.UserLoginID
	if data.Conversation.Type == types.ONE_TO_ONE {
		receiver = networkid.UserLoginID(tc.userLogin.ID)
	}
	return &simplevent.Reaction{
		EventMeta: simplevent.EventMeta{
			Type: eventType,
			LogContext: func(c zerolog.Context) zerolog.Context {
				return c.
					Str("message_id", data.MessageID).
					Str("sender", data.SenderID).
					Str("reaction_key", data.ReactionKey).
					Str("emoji_reaction", data.EmojiReaction)
			},
			PortalKey: networkid.PortalKey{
				ID:       networkid.PortalID(data.Conversation.ConversationID),
				Receiver: receiver,
			},
			Timestamp: data.Time,
			Sender: bridgev2.EventSender{
				IsFromMe: data.SenderID == string(tc.userLogin.ID),
				Sender:   networkid.UserID(data.SenderID),
			},
		},
		EmojiID:       "",
		Emoji:         data.EmojiReaction,
		TargetMessage: networkid.MessageID(data.MessageID),
	}
}
