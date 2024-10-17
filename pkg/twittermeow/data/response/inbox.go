package response

import (
	"encoding/json"
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

type XInboxData struct {
	Status                   types.PaginationStatus             `json:"status,omitempty"`
	MinEntryID               string                             `json:"min_entry_id,omitempty"`
	MaxEntryID               string                             `json:"max_entry_id,omitempty"`
	LastSeenEventID          string                             `json:"last_seen_event_id,omitempty"`
	TrustedLastSeenEventID   string                             `json:"trusted_last_seen_event_id,omitempty"`
	UntrustedLastSeenEventID string                             `json:"untrusted_last_seen_event_id,omitempty"`
	Cursor                   string                             `json:"cursor,omitempty"`
	InboxTimelines           InboxTimelines                     `json:"inbox_timelines,omitempty"`
	Entries                  []map[event.XEventType]interface{} `json:"entries,omitempty"`
	Users                    map[string]types.User              `json:"users,omitempty"`
	Conversations            map[string]types.Conversation      `json:"conversations,omitempty"`
	KeyRegistryState         KeyRegistryState                   `json:"key_registry_state,omitempty"`
}

type XInboxConversationFeedData struct {
	Conversation types.Conversation
	Participants []types.User
	// sorted by timestamp, first index is the most recent message
	Messages []types.Message
}

func (data *XInboxData) Prettify() ([]XInboxConversationFeedData, error) {
	conversationFeeds := make([]XInboxConversationFeedData, 0)
	sortedConversations := methods.SortConversationsByTimestamp(data.Conversations)

	for _, conv := range sortedConversations {
		messages, err := data.GetMessageEntriesByConversationID(conv.ConversationID, true)
		if err != nil {
			return nil, err
		}

		feedData := XInboxConversationFeedData{
			Conversation: conv,
			Participants: data.GetParticipantUsers(conv.Participants),
			Messages:     messages,
		}
		conversationFeeds = append(conversationFeeds, feedData)
	}

	return conversationFeeds, nil
}

type PrettifiedMessage struct {
	EventID        string
	ConversationID string
	MessageID      string
	Recipient      types.User
	Sender         types.User
	SentAt         time.Time
	AffectsSort    bool
	Text           string
	Attachment     *types.Attachment
	Entities       types.Entities
	Reactions      []types.MessageReaction
}

func (data *XInboxData) PrettifyMessages(conversationId string) ([]PrettifiedMessage, error) {
	messages, err := data.GetMessageEntriesByConversationID(conversationId, true)
	if err != nil {
		return nil, err
	}

	prettifiedMessages := make([]PrettifiedMessage, 0)
	for _, msg := range messages {
		sentAt, err := methods.UnixStringMilliToTime(msg.Time)
		if err != nil {
			return nil, err
		}

		prettifiedMessage := PrettifiedMessage{
			EventID:        msg.ID,
			ConversationID: msg.ConversationID,
			MessageID:      msg.MessageData.ID,
			Sender:         data.GetUserByID(msg.MessageData.SenderID),
			Recipient:      data.GetUserByID(msg.MessageData.RecipientID),
			SentAt:         sentAt,
			Text:           msg.MessageData.Text,
			Attachment:     msg.MessageData.Attachment,
			Entities:       msg.MessageData.Entities,
			AffectsSort:    msg.AffectsSort,
			Reactions:      msg.MessageReactions,
		}
		prettifiedMessages = append(prettifiedMessages, prettifiedMessage)
	}

	return prettifiedMessages, nil
}

func (data *XInboxData) GetParticipantUsers(participants []types.Participant) []types.User {
	result := make([]types.User, 0)
	for _, participant := range participants {
		result = append(result, data.GetUserByID(participant.UserID))
	}
	return result
}

func (data *XInboxData) GetMessageEntriesByConversationID(conversationId string, sortByTimestamp bool) ([]types.Message, error) {
	messages := make([]types.Message, 0)
	for _, entry := range data.Entries {
		for entryType, entryData := range entry {
			if entryType == event.XMessageEvent {
				jsonEvData, err := json.Marshal(entryData)
				if err != nil {
					return nil, err
				}
				var message types.Message
				err = json.Unmarshal(jsonEvData, &message)
				if err != nil {
					return nil, err
				}

				if message.ConversationID == conversationId {
					messages = append(messages, message)
				}
			} else if entryType == event.XMessageEditEvent {
				// todo: add edits and return them alongside messages
			}
		}
	}

	if sortByTimestamp {
		methods.SortMessagesByTime(messages)
	}

	return messages, nil
}
