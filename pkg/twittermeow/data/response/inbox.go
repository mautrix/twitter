package response

import (
	"context"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

type GetDMUserUpdatesResponse struct {
	InboxInitialState *TwitterInboxData `json:"inbox_initial_state,omitempty"`
	UserEvents        *TwitterInboxData `json:"user_events,omitempty"`
}

type TwitterInboxData struct {
	Status                   types.PaginationStatus         `json:"status,omitempty"`
	MinEntryID               string                         `json:"min_entry_id,omitempty"`
	MaxEntryID               string                         `json:"max_entry_id,omitempty"`
	LastSeenEventID          string                         `json:"last_seen_event_id,omitempty"`
	TrustedLastSeenEventID   string                         `json:"trusted_last_seen_event_id,omitempty"`
	UntrustedLastSeenEventID string                         `json:"untrusted_last_seen_event_id,omitempty"`
	Cursor                   string                         `json:"cursor,omitempty"`
	InboxTimelines           InboxTimelines                 `json:"inbox_timelines,omitempty"`
	Entries                  []types.RawTwitterEvent        `json:"entries,omitempty"`
	Users                    map[string]*types.User         `json:"users,omitempty"`
	Conversations            map[string]*types.Conversation `json:"conversations,omitempty"`
	KeyRegistryState         KeyRegistryState               `json:"key_registry_state,omitempty"`
}

func (data *TwitterInboxData) GetUserByID(userID string) *types.User {
	if data == nil {
		return nil
	}
	if user, ok := data.Users[userID]; ok {
		return user
	}
	return nil
}

func (data *TwitterInboxData) GetConversationByID(conversationID string) *types.Conversation {
	if data == nil {
		return nil
	}
	if conv, ok := data.Conversations[conversationID]; ok {
		return conv
	}
	return nil
}

func (data *TwitterInboxData) SortedConversations() []*types.Conversation {
	return methods.SortConversationsByTimestamp(data.Conversations)
}

func (data *TwitterInboxData) SortedMessages(ctx context.Context) map[string][]*types.Message {
	conversations := make(map[string][]*types.Message)
	log := zerolog.Ctx(ctx)
	for _, entry := range data.Entries {
		switch evt := entry.ParseWithErrorLog(log).(type) {
		case *types.Message:
			conversations[evt.ConversationID] = append(conversations[evt.ConversationID], evt)
		}
	}
	for _, conv := range conversations {
		methods.SortMessagesByTime(conv)
	}
	return conversations
}
