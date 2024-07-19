package response

import (
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
)

type InboxTimelineResponse struct {
	InboxTimeline XInboxData `json:"inbox_timeline"`
}

type ConversationDMResponse struct {
	ConversationTimeline XInboxData `json:"conversation_timeline"`
}

type InboxInitialStateResponse struct {
	InboxInitialState XInboxData `json:"inbox_initial_state,omitempty"`
}
type Trusted struct {
	Status     types.PaginationStatus `json:"status,omitempty"`
	MinEntryID string                 `json:"min_entry_id,omitempty"`
}
type Untrusted struct {
	Status types.PaginationStatus `json:"status,omitempty"`
}
type InboxTimelines struct {
	Trusted   Trusted   `json:"trusted,omitempty"`
	Untrusted Untrusted `json:"untrusted,omitempty"`
}

type KeyRegistryState struct {
	Status types.PaginationStatus `json:"status,omitempty"`
}
type InboxInitialState struct {
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

type DMMessageDeleteMutationResponse struct {
	Data struct {
		DmMessageHideDelete string `json:"dm_message_hide_delete,omitempty"`
	} `json:"data,omitempty"`
}

type PinConversationResponse struct {
	Data struct {
		AddDmConversationLabelV3 struct {
			Typename  string `json:"__typename,omitempty"`
			LabelType string `json:"label_type,omitempty"`
			Timestamp int64  `json:"timestamp,omitempty"`
		} `json:"add_dm_conversation_label_v3,omitempty"`
	}
}

type UnpinConversationResponse struct {
	Data struct {
		DmConversationLabelDelete string `json:"dm_conversation_label_delete,omitempty"`
	} `json:"data,omitempty"`
}
