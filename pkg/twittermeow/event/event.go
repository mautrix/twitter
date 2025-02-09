package event

import (
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type XEventType string

const (
	XMessageEvent                    XEventType = "message"
	XMessageDeleteEvent              XEventType = "message_delete"
	XMessageEditEvent                XEventType = "message_edit"
	XReactionCreatedEvent            XEventType = "reaction_create"
	XReactionDeletedEvent            XEventType = "reaction_delete"
	XConversationMetadataUpdateEvent XEventType = "conversation_metadata_update"
	XConversationNameUpdate          XEventType = "conversation_name_update"
	XConversationCreateEvent         XEventType = "conversation_create"
	XConversationDeleteEvent         XEventType = "remove_conversation"
	XParticipantsJoinedEvent         XEventType = "participants_join"
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
	Entities     *types.Entities
	Attachment   *types.Attachment
	ReplyData    types.ReplyData
	Reactions    []types.MessageReaction
	EditCount    int
}

type XEventConversationCreated struct {
	Conversation types.Conversation
	EventID      string
	CreatedAt    time.Time
	AffectsSort  bool
	RequestID    string
}

type XEventConversationDelete struct {
	ConversationID string
	EventID        string
	DeletedAt      time.Time
	AffectsSort    bool
	LastEventID    string
}

type XEventParticipantsJoined struct {
	EventID         string
	EventTime       time.Time
	AffectsSort     bool
	Conversation    types.Conversation
	Sender          types.User
	NewParticipants []types.User
}

type XEventConversationMetadataUpdate struct {
	Conversation types.Conversation
	EventID      string
	UpdatedAt    time.Time
	AffectsSort  bool
}

type XEventConversationNameUpdate struct {
	Conversation types.Conversation
	EventID      string
	UpdatedAt    time.Time
	Name         string
	Executor     types.User
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

type XEventReaction struct {
	Conversation  types.Conversation
	Action        types.MessageReactionAction
	ID            string
	Time          time.Time
	MessageID     string
	ReactionKey   string
	EmojiReaction string
	SenderID      string
	RecipientID   string // empty for group chats
	AffectsSort   bool
}
