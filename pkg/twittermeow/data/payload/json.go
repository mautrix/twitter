package payload

import "encoding/json"

type SendDirectMessagePayload struct {
	ConversationID           string `json:"conversation_id,omitempty"`
	MediaID                  string `json:"media_id,omitempty"`
	ReplyToDMID              string `json:"reply_to_dm_id,omitempty"`
	RecipientIDs             bool   `json:"recipient_ids"`
	RequestID                string `json:"request_id,omitempty"`
	Text                     string `json:"text"`
	CardsPlatform            string `json:"cards_platform,omitempty"`
	IncludeCards             int    `json:"include_cards,omitempty"`
	IncludeQuoteCount        bool   `json:"include_quote_count"`
	DMUsers                  bool   `json:"dm_users"`
	AudioOnlyMediaAttachment bool   `json:"audio_only_media_attachment"`
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

type SendTypingNotificationVariables struct {
	ConversationID string `json:"conversationId,omitempty"`
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

type AddParticipantsPayload struct {
	ConversationID    string   `json:"conversationId"`
	AddedParticipants []string `json:"addedParticipants"`
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

type PushNotificationSettings struct {
	Addressbook     string `json:"AddressbookSetting,omitempty"`
	Ads             string `json:"AdsSetting,omitempty"`
	DirectMessages  string `json:"DirectMessagesSetting,omitempty"`
	DMReaction      string `json:"DmReactionSetting,omitempty"`
	FollowersNonVit string `json:"FollowersNonVitSetting,omitempty"`
	FollowersVit    string `json:"FollowersVitSetting,omitempty"`
	LifelineAlerts  string `json:"LifelineAlertsSetting,omitempty"`
	LikesNonVit     string `json:"LikesNonVitSetting,omitempty"`
	LikesVit        string `json:"LikesVitSetting,omitempty"`
	LiveVideo       string `json:"LiveVideoSetting,omitempty"`
	Mentions        string `json:"MentionsSetting,omitempty"`
	Moments         string `json:"MomentsSetting,omitempty"`
	News            string `json:"NewsSetting,omitempty"`
	PhotoTags       string `json:"PhotoTagsSetting,omitempty"`
	Recommendations string `json:"RecommendationsSetting,omitempty"`
	Retweets        string `json:"RetweetsSetting,omitempty"`
	Spaces          string `json:"SpacesSetting,omitempty"`
	Topics          string `json:"TopicsSetting,omitempty"`
	Tweets          string `json:"TweetsSetting,omitempty"`
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

	Checksum string                    `json:"checksum,omitempty"`
	Settings *PushNotificationSettings `json:"settings,omitempty"`
}

type PushConfigPayloadWrapper struct {
	PushDeviceInfo *WebPushConfigPayload `json:"push_device_info"`
}
