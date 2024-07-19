package types

type MessageData struct {
	ID          string     `json:"id,omitempty"`
	Time        string     `json:"time,omitempty"`
	RecipientID string     `json:"recipient_id,omitempty"`
	SenderID    string     `json:"sender_id,omitempty"`
	Text        string     `json:"text,omitempty"`
	Entities    Entities   `json:"entities,omitempty"`
	Attachment  Attachment `json:"attachment,omitempty"`
}

type Message struct {
	ID             string      `json:"id,omitempty"`
	Time           string      `json:"time,omitempty"`
	AffectsSort    bool        `json:"affects_sort,omitempty"`
	RequestID      string      `json:"request_id,omitempty"`
	ConversationID string      `json:"conversation_id,omitempty"`
	MessageData    MessageData `json:"message_data,omitempty"`
}

type ConversationRead struct {
	ID              string `json:"id,omitempty"`
	Time            string `json:"time,omitempty"`
	AffectsSort     bool   `json:"affects_sort,omitempty"`
	ConversationID  string `json:"conversation_id,omitempty"`
	LastReadEventID string `json:"last_read_event_id,omitempty"`
}

type ConversationMetadataUpdate struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
}

type ConversationCreatedData struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	RequestID      string `json:"request_id,omitempty"`
}

type ConversationType string

const (
	ONE_TO_ONE ConversationType = "ONE_TO_ONE"
	GROUP_DM   ConversationType = "GROUP_DM"
)

type PaginationStatus string

const (
	AT_END   PaginationStatus = "AT_END"
	HAS_MORE PaginationStatus = "HAS_MORE"
)

type Conversation struct {
	ConversationID               string           `json:"conversation_id,omitempty"`
	Type                         ConversationType `json:"type,omitempty"`
	SortEventID                  string           `json:"sort_event_id,omitempty"`
	SortTimestamp                string           `json:"sort_timestamp,omitempty"`
	CreateTime                   string           `json:"create_time,omitempty"`
	CreatedByUserID              string           `json:"created_by_user_id,omitempty"`
	Participants                 []Participants   `json:"participants,omitempty"`
	Nsfw                         bool             `json:"nsfw,omitempty"`
	NotificationsDisabled        bool             `json:"notifications_disabled,omitempty"`
	MentionNotificationsDisabled bool             `json:"mention_notifications_disabled,omitempty"`
	LastReadEventID              string           `json:"last_read_event_id,omitempty"`
	ReadOnly                     bool             `json:"read_only,omitempty"`
	Trusted                      bool             `json:"trusted,omitempty"`
	Muted                        bool             `json:"muted,omitempty"`
	Status                       PaginationStatus `json:"status,omitempty"`
	MinEntryID                   string           `json:"min_entry_id,omitempty"`
	MaxEntryID                   string           `json:"max_entry_id,omitempty"`
}

type Participants struct {
	UserID                  string `json:"user_id,omitempty"`
	LastReadEventID         string `json:"last_read_event_id,omitempty"`
	IsAdmin                 bool   `json:"is_admin,omitempty"`
	JoinTime                string `json:"join_time,omitempty"`
	JoinConversationEventId string `json:"join_conversation_event_id,omitempty"`
}

type MessageDeleted struct {
	ID             string            `json:"id,omitempty"`
	Time           string            `json:"time,omitempty"`
	AffectsSort    bool              `json:"affects_sort,omitempty"`
	RequestID      string            `json:"request_id,omitempty"`
	ConversationID string            `json:"conversation_id,omitempty"`
	Messages       []MessagesDeleted `json:"messages,omitempty"`
}

type MessagesDeleted struct {
	MessageID            string `json:"message_id,omitempty"`
	MessageCreateEventID string `json:"message_create_event_id,omitempty"`
}
