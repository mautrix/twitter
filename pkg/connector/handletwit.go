package connector

import (
	"github.com/rs/zerolog"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

func (tc *TwitterClient) HandleTwitterEvent(rawEvt any) {
	switch evtData := rawEvt.(type) {
		case event.XEventMessage:
			sender := evtData.Sender
			isFromMe := sender.IDStr == string(tc.userLogin.ID)
			tc.connector.br.QueueRemoteEvent(tc.userLogin, &bridgev2.SimpleRemoteEvent[*event.XEventMessage]{
				Type: bridgev2.RemoteEventMessage,
				ID: networkid.MessageID(evtData.MessageID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("message_id", evtData.MessageID).
						Str("sender", sender.IDStr).
						Str("sender_login", sender.ScreenName).
						Bool("is_from_me", isFromMe)
				},
				Sender: bridgev2.EventSender{
					IsFromMe: isFromMe,
					SenderLogin: networkid.UserLoginID(sender.IDStr),
					Sender: networkid.UserID(sender.IDStr),
				},
				PortalKey: tc.MakePortalKey(evtData.Conversation),
				Data: &evtData,
				ConvertMessageFunc: tc.convertToMatrix,
				CreatePortal: true,
				Timestamp: evtData.CreatedAt,
			})
		case event.XEventReaction:
			reactionRemoteEvent := tc.wrapReaction(evtData)
			tc.connector.br.QueueRemoteEvent(tc.userLogin, reactionRemoteEvent)
		case event.XEventConversationRead:
			/*
			eventData := &simplevent.Receipt{
				EventMeta: simplevent.EventMeta{
					Type: bridgev2.RemoteEventReadReceipt,
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("conversation_id", evtData.Conversation.ConversationID).
							Str("last_read_event_id", evtData.LastReadEventID).
							Str("read_at", evtData.ReadAt.String())
					},
					PortalKey: tc.MakePortalKey(evtData.Conversation),
				},
				LastTarget: networkid.MessageID(evtData.LastReadEventID),
				Targets: []networkid.MessageID{networkid.MessageID(evtData.LastReadEventID)},
			}
			*/
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
				messageDeleteRemoteEvent := &bridgev2.SimpleRemoteEvent[*types.MessagesDeleted]{
					Type: bridgev2.RemoteEventMessageRemove,
					PortalKey: tc.MakePortalKey(evtData.Conversation),
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("message_id", deletedMsg.MessageID).
							Str("message_create_event_id", deletedMsg.MessageCreateEventID)
					},
					TargetMessage: networkid.MessageID(deletedMsg.MessageID),
					Timestamp: evtData.DeletedAt,
					Data: &deletedMsg,
				}
				tc.connector.br.QueueRemoteEvent(tc.userLogin, messageDeleteRemoteEvent)
			}
		default:
			tc.client.Logger.Warn().Any("event_data", evtData).Msg("Received unhandled event case from twitter library")
		}
}

func (tc *TwitterClient) wrapReaction(data event.XEventReaction) *bridgev2.SimpleRemoteEvent[*event.XEventReaction] {
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
	return &bridgev2.SimpleRemoteEvent[*event.XEventReaction]{
		Type: eventType,
		Data: &data,
		LogContext: func(c zerolog.Context) zerolog.Context {
			return c.
				Str("message_id", data.MessageID).
				Str("sender", data.SenderID).
				Str("reaction_key", data.ReactionKey).
				Str("emoji_reaction", data.EmojiReaction)
		},
		PortalKey: networkid.PortalKey{
			ID: networkid.PortalID(data.Conversation.ConversationID),
			Receiver: receiver,
		},
		EmojiID: "",
		Emoji: data.EmojiReaction,
		TargetMessage: networkid.MessageID(data.MessageID),
		Timestamp: data.Time,
		Sender: bridgev2.EventSender{
			IsFromMe: data.SenderID == string(tc.userLogin.ID),
			Sender: networkid.UserID(data.SenderID),
		},
	}
}