package response

import (
	"encoding/json"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

type GetDMUserUpdatesResponse struct {
	InboxInitialState XInboxData `json:"inbox_initial_state,omitempty"`
	UserEvents        XInboxData `json:"user_events,omitempty"`
}

type UserEvents struct {
	MinEntryID               string                             `json:"min_entry_id,omitempty"`
	MaxEntryID               string                             `json:"max_entry_id,omitempty"`
	Cursor                   string                             `json:"cursor,omitempty"`
	LastSeenEventID          string                             `json:"last_seen_event_id,omitempty"`
	TrustedLastSeenEventID   string                             `json:"trusted_last_seen_event_id,omitempty"`
	UntrustedLastSeenEventID string                             `json:"untrusted_last_seen_event_id,omitempty"`
	Entries                  []map[event.XEventType]interface{} `json:"entries,omitempty"`
	Users                    map[string]types.User              `json:"users,omitempty"`
	Conversations            map[string]types.Conversation      `json:"conversations,omitempty"`
}

func (data *XInboxData) GetUserByID(userId string) types.User {
	if user, ok := data.Users[userId]; ok {
		return user
	}
	return types.User{}
}

func (data *XInboxData) GetConversationByID(conversationId string) types.Conversation {
	if conv, ok := data.Conversations[conversationId]; ok {
		return conv
	}
	return types.Conversation{}
}

func (data *XInboxData) ToEventEntries() ([]interface{}, error) {
	entries := make([]interface{}, 0)
	if len(data.Entries) <= 0 {
		return entries, nil
	}

	for _, entry := range data.Entries {
		for entryType, entryData := range entry {
			var updatedEntry interface{}
			jsonEvData, err := json.Marshal(entryData)
			if err != nil {
				return nil, err
			}
			switch entryType {
			case event.XReactionCreatedEvent, event.XReactionDeletedEvent:
				var reactionEventData types.MessageReaction
				err = json.Unmarshal(jsonEvData, &reactionEventData)
				if err != nil {
					return nil, err
				}

				reactionActionAt, err := methods.UnixStringMilliToTime(reactionEventData.Time)
				if err != nil {
					return nil, err
				}
				
				updatedReactionEventData := event.XEventReaction{
					Conversation: data.GetConversationByID(reactionEventData.ConversationID),
					Time: reactionActionAt,
					ID: reactionEventData.ID,
					ReactionKey: reactionEventData.ReactionKey,
					EmojiReaction: reactionEventData.EmojiReaction,
					AffectsSort: reactionEventData.AffectsSort,
					SenderID: reactionEventData.SenderID,
					MessageID: reactionEventData.MessageID,
				}
				switch entryType {
				case event.XReactionCreatedEvent:
					updatedReactionEventData.Action = types.MessageReactionAdd
				case event.XReactionDeletedEvent:
					updatedReactionEventData.Action = types.MessageReactionRemove
				default:
					break
				}
				updatedEntry = updatedReactionEventData
			case event.XMessageEvent:
				var messageEventData types.Message
				err = json.Unmarshal(jsonEvData, &messageEventData)
				if err != nil {
					return nil, err
				}

				createdAt, err := methods.UnixStringMilliToTime(messageEventData.MessageData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventMessage{
					Conversation: data.GetConversationByID(messageEventData.ConversationID),
					Sender:       data.GetUserByID(messageEventData.MessageData.SenderID),
					Recipient:    data.GetUserByID(messageEventData.MessageData.RecipientID),
					MessageID:    messageEventData.MessageData.ID,
					CreatedAt:    createdAt,
					Text:         messageEventData.MessageData.Text,
					Entities:     &messageEventData.MessageData.Entities,
					Attachment:   messageEventData.MessageData.Attachment,
					ReplyData: 	  messageEventData.MessageData.ReplyData,
					AffectsSort:  messageEventData.AffectsSort,
					Reactions: 	  messageEventData.MessageReactions,
				}
			case event.XConversationNameUpdate:
				var convNameUpdateEventData types.ConversationNameUpdateData
				err = json.Unmarshal(jsonEvData, &convNameUpdateEventData)
				if err != nil {
					return nil, err
				}

				updatedAt, err := methods.UnixStringMilliToTime(convNameUpdateEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventConversationNameUpdate{
					Conversation: data.GetConversationByID(convNameUpdateEventData.ConversationID),
					Name: convNameUpdateEventData.ConversationName,
					Executor: data.GetUserByID(convNameUpdateEventData.ByUserID),
					EventID: convNameUpdateEventData.ID,
					AffectsSort: convNameUpdateEventData.AffectsSort,
					UpdatedAt: updatedAt,
				}
			case event.XParticipantsJoinedEvent:
				var participantsJoinedEventData types.ParticipantsJoinedData
				err = json.Unmarshal(jsonEvData, &participantsJoinedEventData)
				if err != nil {
					return nil, err
				}

				eventTime, err := methods.UnixStringMilliToTime(participantsJoinedEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventParticipantsJoined{
					EventID: participantsJoinedEventData.ID,
					EventTime: eventTime,
					AffectsSort: participantsJoinedEventData.AffectsSort,
					Conversation: data.GetConversationByID(participantsJoinedEventData.ConversationID),
					Sender: data.GetUserByID(participantsJoinedEventData.SenderID),
					NewParticipants: data.GetParticipantUsers(participantsJoinedEventData.Participants),
				}
			case event.XMessageDeleteEvent:
				var messageDeletedEventData types.MessageDeleted
				err = json.Unmarshal(jsonEvData, &messageDeletedEventData)
				if err != nil {
					return nil, err
				}

				deletedAt, err := methods.UnixStringMilliToTime(messageDeletedEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventMessageDeleted{
					Conversation: data.GetConversationByID(messageDeletedEventData.ConversationID),
					DeletedAt:    deletedAt,
					EventID:      messageDeletedEventData.ID,
					RequestID:    messageDeletedEventData.RequestID,
					AffectsSort:  messageDeletedEventData.AffectsSort,
					Messages:     messageDeletedEventData.Messages,
				}
			case event.XConversationReadEvent:
				var convReadEventData types.ConversationRead
				err = json.Unmarshal(jsonEvData, &convReadEventData)
				if err != nil {
					return nil, err
				}

				readAt, err := methods.UnixStringMilliToTime(convReadEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventConversationRead{
					EventID:         convReadEventData.ID,
					Conversation:    data.GetConversationByID(convReadEventData.ConversationID),
					ReadAt:          readAt,
					AffectsSort:     convReadEventData.AffectsSort,
					LastReadEventID: convReadEventData.LastReadEventID,
				}
			case event.XConversationDeleteEvent:
				var convDeletedEventData types.ConversationDeletedData
				err = json.Unmarshal(jsonEvData, &convDeletedEventData)
				if err != nil {
					return nil, err
				}

				deletedAt, err := methods.UnixStringMilliToTime(convDeletedEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventConversationDelete{
					ConversationID:  convDeletedEventData.ConversationID,
					EventID:         convDeletedEventData.ID,
					DeletedAt: 		 deletedAt,
					AffectsSort:     convDeletedEventData.AffectsSort,
					LastEventID: 	 convDeletedEventData.LastEventID,
				}
			case event.XConversationCreateEvent:
				var convCreatedEventData types.ConversationCreatedData
				err = json.Unmarshal(jsonEvData, &convCreatedEventData)
				if err != nil {
					return nil, err
				}

				createdAt, err := methods.UnixStringMilliToTime(convCreatedEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventConversationCreated{
					EventID:      convCreatedEventData.ID,
					Conversation: data.GetConversationByID(convCreatedEventData.ConversationID),
					CreatedAt:    createdAt,
					AffectsSort:  convCreatedEventData.AffectsSort,
					RequestID:    convCreatedEventData.RequestID,
				}
				
			case event.XConversationMetadataUpdateEvent:
				var convMetadataUpdateEventData types.ConversationMetadataUpdate
				err = json.Unmarshal(jsonEvData, &convMetadataUpdateEventData)
				if err != nil {
					return nil, err
				}

				updatedAt, err := methods.UnixStringMilliToTime(convMetadataUpdateEventData.Time)
				if err != nil {
					return nil, err
				}

				updatedEntry = event.XEventConversationMetadataUpdate{
					EventID:      convMetadataUpdateEventData.ID,
					Conversation: data.GetConversationByID(convMetadataUpdateEventData.ConversationID),
					UpdatedAt:    updatedAt,
					AffectsSort:  convMetadataUpdateEventData.AffectsSort,
				}
			}
			entries = append(entries, updatedEntry)
		}
	}

	return entries, nil
}
