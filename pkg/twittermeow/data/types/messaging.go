package types

type MessageData struct {
	ID          string      `json:"id,omitempty"`
	Time        string      `json:"time,omitempty"`
	RecipientID string      `json:"recipient_id,omitempty"`
	SenderID    string      `json:"sender_id,omitempty"`
	Text        string      `json:"text,omitempty"`
	Entities    Entities    `json:"entities,omitempty"`
	Attachment  *Attachment `json:"attachment,omitempty"`
	ReplyData   ReplyData   `json:"reply_data,omitempty"`
	EditCount   int         `json:"edit_count,omitempty"`
}

type Message struct {
	ID               string            `json:"id,omitempty"`
	Time             string            `json:"time,omitempty"`
	AffectsSort      bool              `json:"affects_sort,omitempty"`
	RequestID        string            `json:"request_id,omitempty"`
	ConversationID   string            `json:"conversation_id,omitempty"`
	MessageData      MessageData       `json:"message_data,omitempty"`
	MessageReactions []MessageReaction `json:"message_reactions,omitempty"`
}

type ReplyData struct {
	ID          string `json:"id,omitempty"`
	Time        string `json:"time,omitempty"`
	RecipientID string `json:"recipient_id,omitempty"`
	SenderID    string `json:"sender_id,omitempty"`
	Text        string `json:"text,omitempty"`
}

type MessageReactionAction string

const (
	MessageReactionAdd    MessageReactionAction = "reaction_add"
	MessageReactionRemove MessageReactionAction = "reaction_remove"
)

type MessageReaction struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	MessageID      string `json:"message_id,omitempty"`
	ReactionKey    string `json:"reaction_key,omitempty"`
	EmojiReaction  string `json:"emoji_reaction,omitempty"`
	SenderID       string `json:"sender_id,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
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

type ConversationDeletedData struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	LastEventID    string `json:"last_event_id,omitempty"`
}

type ConversationNameUpdateData struct {
	ID               string `json:"id,omitempty"`
	Time             string `json:"time,omitempty"`
	ConversationID   string `json:"conversation_id,omitempty"`
	ConversationName string `json:"conversation_name,omitempty"`
	ByUserID         string `json:"by_user_id,omitempty"`
	AffectsSort      bool   `json:"affects_sort,omitempty"`
}

type ParticipantsJoinedData struct {
	ID             string        `json:"id,omitempty"`
	Time           string        `json:"time,omitempty"`
	AffectsSort    bool          `json:"affects_sort,omitempty"`
	ConversationID string        `json:"conversation_id,omitempty"`
	SenderID       string        `json:"sender_id,omitempty"`
	Participants   []Participant `json:"participants,omitempty"`
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
	Name                         string           `json:"name,omitempty"`
	AvatarImageHttps             string           `json:"avatar_image_https,omitempty"`
	Avatar                       Avatar           `json:"avatar,omitempty"`
	SortEventID                  string           `json:"sort_event_id,omitempty"`
	SortTimestamp                string           `json:"sort_timestamp,omitempty"`
	CreateTime                   string           `json:"create_time,omitempty"`
	CreatedByUserID              string           `json:"created_by_user_id,omitempty"`
	Participants                 []Participant    `json:"participants,omitempty"`
	Nsfw                         bool             `json:"nsfw,omitempty"`
	NotificationsDisabled        bool             `json:"notifications_disabled,omitempty"`
	MentionNotificationsDisabled bool             `json:"mention_notifications_disabled,omitempty"`
	LastReadEventID              string           `json:"last_read_event_id,omitempty"`
	ReadOnly                     bool             `json:"read_only,omitempty"`
	Trusted                      bool             `json:"trusted,omitempty"`
	Muted                        bool             `json:"muted,omitempty"`
	LowQuality                   bool             `json:"low_quality,omitempty"`
	Status                       PaginationStatus `json:"status,omitempty"`
	MinEntryID                   string           `json:"min_entry_id,omitempty"`
	MaxEntryID                   string           `json:"max_entry_id,omitempty"`
}

type Image struct {
	OriginalInfo OriginalInfo `json:"original_info,omitempty"`
}

type Avatar struct {
	Image Image `json:"image,omitempty"`
}

type Participant struct {
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
