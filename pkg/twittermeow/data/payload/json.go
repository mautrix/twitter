package payload

import (
	"encoding/json"
	"fmt"
	"net/url"
)

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

type SendMessageMutationVariables struct {
	ConversationID               string  `json:"conversation_id"`
	MessageID                    string  `json:"message_id"`
	ConversationToken            string  `json:"conversation_token"`
	EncodedMessageCreateEvent    string  `json:"encoded_message_create_event"`
	EncodedMessageEventSignature *string `json:"encoded_message_event_signature,omitempty"`
}

type SendMessageMutationPayload struct {
	Variables SendMessageMutationVariables `json:"variables"`
	QueryId   string                       `json:"queryId,omitempty"`
}

func NewSendMessageMutationPayload(vars SendMessageMutationVariables) *SendMessageMutationPayload {
	return &SendMessageMutationPayload{
		Variables: vars,
	}
}

type GenerateXChatTokenMutationPayload struct {
	Variables  map[string]any `json:"variables"`
	Features   map[string]any `json:"features,omitempty"`
	Extensions struct {
		PersistedQuery struct {
			Version    int    `json:"version"`
			Sha256Hash string `json:"sha256Hash"`
		} `json:"persistedQuery"`
	} `json:"extensions"`
}

func (p *GenerateXChatTokenMutationPayload) Default() *GenerateXChatTokenMutationPayload {
	p.Variables = map[string]any{}
	p.Extensions.PersistedQuery.Version = 1
	return p
}

// ApolloExtensions - reusable struct for Apollo-style GraphQL persisted queries
type ApolloExtensions struct {
	PersistedQuery struct {
		Version    int    `json:"version"`
		Sha256Hash string `json:"sha256Hash"`
	} `json:"persistedQuery"`
	ClientLibrary struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"clientLibrary"`
}

// DefaultApolloExtensions returns ApolloExtensions with default values
func DefaultApolloExtensions() ApolloExtensions {
	return ApolloExtensions{
		PersistedQuery: struct {
			Version    int    `json:"version"`
			Sha256Hash string `json:"sha256Hash"`
		}{Version: 1},
		ClientLibrary: struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{Name: "apollo-kotlin", Version: "4.3.3"},
	}
}

// InitializeXChatMediaUploadPayload for the InitializeXChatMediaUpload GraphQL request
type InitializeXChatMediaUploadPayload struct {
	OperationName string                              `json:"operationName"`
	Variables     InitializeXChatMediaUploadVariables `json:"variables"`
	Extensions    ApolloExtensions                    `json:"extensions"`
}

type InitializeXChatMediaUploadVariables struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	TotalBytes     string `json:"total_bytes"`
}

func (p *InitializeXChatMediaUploadPayload) Default() *InitializeXChatMediaUploadPayload {
	p.OperationName = "InitializeXChatMediaUpload"
	p.Extensions = DefaultApolloExtensions()
	return p
}

// FinalizeXChatMediaUploadPayload for the FinalizeXChatMediaUpload GraphQL request
type FinalizeXChatMediaUploadPayload struct {
	OperationName string                            `json:"operationName"`
	Variables     FinalizeXChatMediaUploadVariables `json:"variables"`
	Extensions    ApolloExtensions                  `json:"extensions"`
}

type FinalizeXChatMediaUploadVariables struct {
	ConversationID string  `json:"conversation_id"`
	MessageID      string  `json:"message_id"`
	MediaHashKey   string  `json:"media_hash_key"`
	ResumeID       string  `json:"resume_id"`
	NumParts       string  `json:"num_parts"`
	TTLMsec        *string `json:"ttl_msec"`
}

func (p *FinalizeXChatMediaUploadPayload) Default() *FinalizeXChatMediaUploadPayload {
	p.OperationName = "FinalizeXChatMediaUpload"
	p.Extensions = DefaultApolloExtensions()
	return p
}

type GetConversationPageQuerySettings struct {
	InboxConversationEventLimit int `json:"inbox_conversation_event_limit"`
	InboxConversationLimit      int `json:"inbox_conversation_limit"`
	ConversationEventLimit      int `json:"conversation_event_limit"`
	UserEventLimit              int `json:"user_event_limit"`
}

func DefaultGetConversationPageQuerySettings() *GetConversationPageQuerySettings {
	return &GetConversationPageQuerySettings{
		InboxConversationEventLimit: 5,
		InboxConversationLimit:      20,
		ConversationEventLimit:      200,
		UserEventLimit:              500,
	}
}

// GetConversationPageQueryVariables is the JSON body for the GetConversationPageQuery endpoint.
type GetConversationPageQueryVariables struct {
	ConversationID            string                            `json:"conversation_id"`
	MinLocalSequenceID        string                            `json:"min_local_sequence_id,omitempty"`
	MinConversationKeyVersion string                            `json:"min_conversation_key_version,omitempty"`
	QuerySettings             *GetConversationPageQuerySettings `json:"query_settings,omitempty"`
}

func NewGetConversationPageQueryVariables(conversationID, minLocalSequenceID, minConversationKeyVersion string, settings *GetConversationPageQuerySettings) *GetConversationPageQueryVariables {
	if settings == nil {
		settings = DefaultGetConversationPageQuerySettings()
	}
	return &GetConversationPageQueryVariables{
		ConversationID:            conversationID,
		MinLocalSequenceID:        minLocalSequenceID,
		MinConversationKeyVersion: minConversationKeyVersion,
		QuerySettings:             settings,
	}
}

// EncodeJSONQuery serializes the inbox page request variables to JSON and wraps
// them in a query string suitable for ?variables=... on GraphQL GET requests.
func (p *GetInboxPageRequestQueryVariables) EncodeJSONQuery() (string, error) {
	if p.ContinueCursor == nil {
		return "", fmt.Errorf("continue_cursor is required")
	}

	type gqlCursor struct {
		CursorID        string `json:"cursor_id,omitempty"`
		GraphSnapshotID string `json:"graph_snapshot_id,omitempty"`
	}
	type gqlSettings struct {
		InboxConversationEventLimit int `json:"inbox_conversation_event_limit"`
		InboxConversationLimit      int `json:"inbox_conversation_limit"`
		ConversationEventLimit      int `json:"conversation_event_limit"`
		UserEventLimit              int `json:"user_event_limit"`
	}
	type gqlVars struct {
		ContinueCursor *gqlCursor   `json:"continue_cursor"`
		QuerySettings  *gqlSettings `json:"query_settings"`
	}

	vars := gqlVars{
		ContinueCursor: &gqlCursor{
			CursorID:        p.ContinueCursor.CursorId,
			GraphSnapshotID: p.ContinueCursor.GraphSnapshotId,
		},
		QuerySettings: &gqlSettings{
			InboxConversationEventLimit: p.QuerySettings.InboxConversationEventLimit,
			InboxConversationLimit:      p.QuerySettings.InboxConversationLimit,
			ConversationEventLimit:      p.QuerySettings.ConversationEventLimit,
			UserEventLimit:              p.QuerySettings.UserEventLimit,
		},
	}

	jsonVars, err := json.Marshal(vars)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("variables", string(jsonVars))
	return values.Encode(), nil
}
