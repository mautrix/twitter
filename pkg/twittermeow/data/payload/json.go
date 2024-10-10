package payload

import "encoding/json"

type SendDirectMessagePayload struct {
	ConversationID    string `json:"conversation_id,omitempty"`
	MediaID           string `json:"media_id,omitempty"`
	ReplyToDmID       string `json:"reply_to_dm_id,omitempty"`
	RecipientIds      bool   `json:"recipient_ids"`
	RequestID         string `json:"request_id,omitempty"`
	Text              string `json:"text"`
	CardsPlatform     string `json:"cards_platform,omitempty"`
	IncludeCards      int    `json:"include_cards,omitempty"`
	IncludeQuoteCount bool   `json:"include_quote_count"`
	DmUsers           bool   `json:"dm_users"`
}

func (p *SendDirectMessagePayload) Encode() ([]byte, error) {
	return json.Marshal(p)
}

type GraphQLPayload struct {
	Variables interface{} `json:"variables,omitempty"`
	QueryID   string      `json:"queryId,omitempty"`
}

func (p *GraphQLPayload) Encode() ([]byte, error) {
	return json.Marshal(p)
}

type DMMessageDeleteMutationVariables struct {
	MessageID string `json:"messageId,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}

type LabelType string

const (
	LABEL_TYPE_PINNED LabelType = "Pinned"
)

type PinAndUnpinConversationVariables struct {
	ConversationID string    `json:"conversation_id,omitempty"`
	LabelType      LabelType `json:"label_type,omitempty"`
	Label          LabelType `json:"label,omitempty"`
}

type ReactionActionPayload struct {
	ConversationID string   `json:"conversationId"`
	MessageID      string   `json:"messageId"`
	ReactionTypes  []string `json:"reactionTypes"`
	EmojiReactions []string `json:"emojiReactions"`
}

func (p *ReactionActionPayload) Encode() ([]byte, error) {
	return json.Marshal(p)
}

type WebPushConfigPayload struct {
	Token  string `json:"token"`
	P256DH string `json:"encryption_key1"`
	Auth   string `json:"encryption_key2"`

	OSVersion string `json:"os_version"`
	UDID      string `json:"udid"`
	Locale    string `json:"locale"`

	Env             int `json:"env"`
	ProtocolVersion int `json:"protocol_version"`
}
