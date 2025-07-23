package types

type MessageData struct {
	ID          string      `json:"id,omitempty"`
	Time        string      `json:"time,omitempty"`
	RecipientID string      `json:"recipient_id,omitempty"`
	SenderID    string      `json:"sender_id,omitempty"`
	Text        string      `json:"text,omitempty"`
	Entities    *Entities   `json:"entities,omitempty"`
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

type MessageEdit Message

type ReplyData struct {
	ID          string `json:"id,omitempty"`
	Time        string `json:"time,omitempty"`
	RecipientID string `json:"recipient_id,omitempty"`
	SenderID    string `json:"sender_id,omitempty"`
	Text        string `json:"text,omitempty"`
}

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

type MessageReactionCreate MessageReaction
type MessageReactionDelete MessageReaction

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

type ConversationCreate struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	RequestID      string `json:"request_id,omitempty"`
}

type ConversationDelete struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	LastEventID    string `json:"last_event_id,omitempty"`
}

type ConversationNameUpdate struct {
	ID               string `json:"id,omitempty"`
	Time             string `json:"time,omitempty"`
	ConversationID   string `json:"conversation_id,omitempty"`
	ConversationName string `json:"conversation_name,omitempty"`
	ByUserID         string `json:"by_user_id,omitempty"`
	AffectsSort      bool   `json:"affects_sort,omitempty"`
}

type ParticipantsJoin struct {
	ID             string        `json:"id,omitempty"`
	Time           string        `json:"time,omitempty"`
	AffectsSort    bool          `json:"affects_sort,omitempty"`
	ConversationID string        `json:"conversation_id,omitempty"`
	SenderID       string        `json:"sender_id,omitempty"`
	Participants   []Participant `json:"participants,omitempty"`
}

type ParticipantsLeave struct {
	ID             string        `json:"id,omitempty"`
	Time           string        `json:"time,omitempty"`
	AffectsSort    bool          `json:"affects_sort,omitempty"`
	ConversationID string        `json:"conversation_id,omitempty"`
	Participants   []Participant `json:"participants,omitempty"`
}

type ConversationJoin ParticipantsJoin

type ConversationType string

const (
	ConversationTypeOneToOne ConversationType = "ONE_TO_ONE"
	ConversationTypeGroupDM  ConversationType = "GROUP_DM"
)

type PaginationStatus string

const (
	PaginationStatusAtEnd   PaginationStatus = "AT_END"
	PaginationStatusHasMore PaginationStatus = "HAS_MORE"
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
	NSFW                         bool             `json:"nsfw,omitempty"`
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
	JoinConversationEventID string `json:"join_conversation_event_id,omitempty"`
}

type MessageDelete struct {
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

type TrustConversation struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

type EndAVBroadcast struct {
	ID             string `json:"id,omitempty"`
	Time           string `json:"time,omitempty"`
	AffectsSort    bool   `json:"affects_sort,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	IsCaller       bool   `json:"is_caller,omitempty"`
	StartedAtMs    string `json:"started_at_ms,omitempty"`
	EndedAtMs      string `json:"ended_at_ms,omitempty"`
	EndReason      string `json:"end_reason,omitempty"`
	CallType       string `json:"call_type,omitempty"`
}
