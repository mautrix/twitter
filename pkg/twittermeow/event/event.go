package event

import (
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type XEventType string

const (
	XMessageEvent                    XEventType = "message"
	XMessageDeleteEvent              XEventType = "message_delete"
	XConversationReadEvent           XEventType = "conversation_read"
	XConversationMetadataUpdateEvent XEventType = "conversation_metadata_update"
	XConversationCreateEvent         XEventType = "conversation_create"
	XDisableNotificationsEvent       XEventType = "disable_notifications"
)

type XEventMessage struct {
	Conversation types.Conversation
	Sender       types.User
	Recipient    types.User
	MessageID    string
	Text         string
	CreatedAt    time.Time
	AffectsSort  bool
	Entities     types.Entities
	Attachment   types.Attachment
}

type XEventConversationRead struct {
	EventID         string
	Conversation    types.Conversation
	ReadAt          time.Time
	AffectsSort     bool
	LastReadEventID string
}

type XEventConversationCreated struct {
	EventID      string
	Conversation types.Conversation
	CreatedAt    time.Time
	AffectsSort  bool
	RequestID    string
}

type XEventConversationMetadataUpdate struct {
	Conversation types.Conversation
	EventID      string
	UpdatedAt    time.Time
	AffectsSort  bool
}

type XEventMessageDeleted struct {
	Conversation types.Conversation
	EventID      string
	RequestID    string
	DeletedAt    time.Time
	AffectsSort  bool
	Messages     []types.MessagesDeleted
}
