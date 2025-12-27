package payload

import (
	"encoding/json"
	"net/url"

	"github.com/google/go-querystring/query"
)

type MigrationRequestPayload struct {
	Tok  string `url:"tok"`
	Data string `url:"data"`
}

type JotClientEventPayload struct {
	Category JotLoggingCategory `url:"category,omitempty"`
	Debug    bool               `url:"debug,omitempty"`
	Log      string             `url:"log"`
}

func (p *JotClientEventPayload) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type AccountSettingsQuery struct {
	IncludeExtSharingAudiospacesListeningDataWithFollowers bool   `url:"include_ext_sharing_audiospaces_listening_data_with_followers"`
	IncludeMentionFilter                                   bool   `url:"include_mention_filter"`
	IncludeNSFWUserFlag                                    bool   `url:"include_nsfw_user_flag"`
	IncludeNSFWAdminFlag                                   bool   `url:"include_nsfw_admin_flag"`
	IncludeRankedTimeline                                  bool   `url:"include_ranked_timeline"`
	IncludeAltTextCompose                                  bool   `url:"include_alt_text_compose"`
	Ext                                                    string `url:"ext"`
	IncludeCountryCode                                     bool   `url:"include_country_code"`
	IncludeExtDMNSFWMediaFilter                            bool   `url:"include_ext_dm_nsfw_media_filter"`
}

func (p *AccountSettingsQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type ContextInfo string

const (
	CONTEXT_FETCH_DM_CONVERSATION         ContextInfo = "FETCH_DM_CONVERSATION"
	CONTEXT_FETCH_DM_CONVERSATION_HISTORY ContextInfo = "FETCH_DM_CONVERSATION_HISTORY"
)

type DMRequestQuery struct {
	AvatarID                            string      `url:"avatar_id,omitempty"`
	ActiveConversationID                string      `url:"active_conversation_id,omitempty"`
	Cursor                              string      `url:"cursor,omitempty"`
	Count                               int         `url:"count,omitempty"`
	Context                             ContextInfo `url:"context,omitempty"`
	MaxID                               string      `url:"max_id,omitempty"` // when fetching messages, this is the message id
	MinID                               string      `url:"min_id,omitempty"`
	Name                                string      `url:"name,omitempty"`
	NSFWFilteringEnabled                bool        `url:"nsfw_filtering_enabled"`
	FilterLowQuality                    bool        `url:"filter_low_quality"`
	IncludeQuality                      string      `url:"include_quality"`
	IncludeProfileInterstitialType      int         `url:"include_profile_interstitial_type"`
	IncludeBlocking                     int         `url:"include_blocking"`
	IncludeConversationInfo             bool        `url:"include_conversation_info"`
	IncludeBlockedBy                    int         `url:"include_blocked_by"`
	IncludeFollowedBy                   int         `url:"include_followed_by"`
	IncludeWantRetweets                 int         `url:"include_want_retweets"`
	IncludeMuteEdge                     int         `url:"include_mute_edge"`
	IncludeCanDM                        int         `url:"include_can_dm"`
	IncludeCanMediaTag                  int         `url:"include_can_media_tag"`
	IncludeExtIsBlueVerified            int         `url:"include_ext_is_blue_verified"`
	IncludeExtVerifiedType              int         `url:"include_ext_verified_type"`
	IncludeExtProfileImageShape         int         `url:"include_ext_profile_image_shape"`
	SkipStatus                          int         `url:"skip_status"`
	DMSecretConversationsEnabled        bool        `url:"dm_secret_conversations_enabled"`
	KRSRegistrationEnabled              bool        `url:"krs_registration_enabled"`
	CardsPlatform                       string      `url:"cards_platform"`
	IncludeCards                        int         `url:"include_cards"`
	IncludeExtAltText                   bool        `url:"include_ext_alt_text"`
	IncludeExtLimitedActionResults      bool        `url:"include_ext_limited_action_results"`
	IncludeQuoteCount                   bool        `url:"include_quote_count"`
	IncludeReplyCount                   int         `url:"include_reply_count"`
	TweetMode                           string      `url:"tweet_mode"`
	IncludeExtViews                     bool        `url:"include_ext_views"`
	DMUsers                             bool        `url:"dm_users"`
	IncludeGroups                       bool        `url:"include_groups"`
	IncludeInboxTimelines               bool        `url:"include_inbox_timelines"`
	IncludeExtMediaColor                bool        `url:"include_ext_media_color"`
	SupportsReactions                   bool        `url:"supports_reactions"`
	SupportsEdit                        bool        `url:"supports_edit"`
	IncludeExtEditControl               bool        `url:"include_ext_edit_control"`
	IncludeExtBusinessAffiliationsLabel bool        `url:"include_ext_business_affiliations_label"`
	IncludeParodyCommentaryFanLabel     bool        `url:"include_ext_parody_commentary_fan_label"`
	Ext                                 string      `url:"ext"`
}

func (p DMRequestQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

func (p DMRequestQuery) Default() DMRequestQuery {
	return DMRequestQuery{
		NSFWFilteringEnabled:                false,
		FilterLowQuality:                    true,
		IncludeQuality:                      "all",
		IncludeProfileInterstitialType:      1,
		IncludeBlocking:                     1,
		IncludeBlockedBy:                    1,
		IncludeFollowedBy:                   1,
		IncludeWantRetweets:                 1,
		IncludeMuteEdge:                     1,
		IncludeCanDM:                        1,
		IncludeCanMediaTag:                  1,
		IncludeExtIsBlueVerified:            1,
		IncludeExtVerifiedType:              1,
		IncludeExtProfileImageShape:         1,
		SkipStatus:                          1,
		DMSecretConversationsEnabled:        false,
		KRSRegistrationEnabled:              false,
		CardsPlatform:                       "Web-12",
		IncludeCards:                        1,
		IncludeExtAltText:                   true,
		IncludeExtLimitedActionResults:      true,
		IncludeQuoteCount:                   true,
		IncludeReplyCount:                   1,
		TweetMode:                           "extended",
		IncludeExtViews:                     true,
		DMUsers:                             true,
		IncludeGroups:                       true,
		IncludeInboxTimelines:               true,
		IncludeExtMediaColor:                true,
		SupportsReactions:                   true,
		SupportsEdit:                        true,
		IncludeExtEditControl:               true,
		IncludeExtBusinessAffiliationsLabel: true,
		IncludeParodyCommentaryFanLabel:     true,
		Ext:                                 "mediaColor,altText,businessAffiliationsLabel,mediaStats,highlightedLabel,parodyCommentaryFanLabel,voiceInfo,birdwatchPivot,superFollowMetadata,unmentionInfo,editControl,article",
	}
}

type DMSendQuery struct {
	Ext                            string `url:"ext"`
	IncludeExtAltText              bool   `url:"include_ext_alt_text"`
	IncludeExtLimitedActionResults bool   `url:"include_ext_limited_action_results"`
	IncludeReplyCount              int    `url:"include_reply_count"`
	TweetMode                      string `url:"tweet_mode"`
	IncludeExtViews                bool   `url:"include_ext_views"`
	IncludeGroups                  bool   `url:"include_groups"`
	IncludeInboxTimelines          bool   `url:"include_inbox_timelines"`
	IncludeExtMediaColor           bool   `url:"include_ext_media_color"`
	SupportsReactions              bool   `url:"supports_reactions"`
	SupportsEdit                   bool   `url:"supports_edit"`
}

func (p DMSendQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

func (p DMSendQuery) Default() DMSendQuery {
	return DMSendQuery{
		IncludeExtAltText:              true,
		IncludeExtLimitedActionResults: true,
		IncludeReplyCount:              1,
		TweetMode:                      "extended",
		IncludeExtViews:                true,
		IncludeGroups:                  true,
		IncludeInboxTimelines:          true,
		IncludeExtMediaColor:           true,
		SupportsReactions:              true,
		SupportsEdit:                   true,
		Ext:                            "mediaColor,altText,mediaStats,highlightedLabel,voiceInfo,birdwatchPivot,superFollowMetadata,unmentionInfo,editControl,article",
	}
}

type MarkConversationReadQuery struct {
	ConversationID  string `url:"conversationId"`
	LastReadEventID string `url:"last_read_event_id"`
}

func (p *MarkConversationReadQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type MediaCategory string

const (
	MEDIA_CATEGORY_DM_IMAGE MediaCategory = "dm_image"
	MEDIA_CATEGORY_DM_VIDEO MediaCategory = "dm_video"
	MEDIA_CATEGORY_DM_GIF   MediaCategory = "dm_gif"
)

type UploadMediaQuery struct {
	Command         string        `url:"command,omitempty"`
	TotalBytes      int           `url:"total_bytes,omitempty"`
	SourceURL       string        `url:"source_url,omitempty"`
	MediaID         string        `url:"media_id,omitempty"`
	VideoDurationMS float32       `url:"video_duration_ms,omitempty"`
	OriginalMD5     string        `url:"original_md5,omitempty"`
	SegmentIndex    int           `url:"segment_index,omitempty"`
	MediaType       string        `url:"media_type,omitempty"`
	MediaCategory   MediaCategory `url:"media_category,omitempty"`
}

func (p *UploadMediaQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type SearchResultType string

const (
	SEARCH_RESULT_TYPE_USERS SearchResultType = "users"
)

type SearchQuery struct {
	IncludeExtIsBlueVerified    string           `url:"include_ext_is_blue_verified"`
	IncludeExtVerifiedType      string           `url:"include_ext_verified_type"`
	IncludeExtProfileImageShape string           `url:"include_ext_profile_image_shape"`
	Query                       string           `url:"q"`
	Src                         string           `url:"src"`
	ResultType                  SearchResultType `url:"result_type"`
}

func (p *SearchQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type GetDMPermissionsQuery struct {
	// seperated by commas: userid1,userid2,userid3
	RecipientIDs string `url:"recipient_ids"`
	DMUsers      bool   `url:"dm_users"`
}

func (p *GetDMPermissionsQuery) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type UpdateSubscriptionsPayload struct {
	SubTopics   string `url:"sub_topics"`
	UnsubTopics string `url:"unsub_topics"`
}

func (p *UpdateSubscriptionsPayload) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

var XCHAT_MESSAGE_PULL_VERSION int = 1761251295

type QuerySettings struct {
	InboxConversationEventLimit int `url:"inbox_conversation_event_limit"`
	InboxConversationLimit      int `url:"inbox_conversation_limit"`
	ConversationEventLimit      int `url:"conversation_event_limit"`
	UserEventLimit              int `url:"user_event_limit"`
}

func DefaultQuerySettings() *QuerySettings {
	return &QuerySettings{
		InboxConversationEventLimit: 5,
		InboxConversationLimit:      20,
		ConversationEventLimit:      200,
		UserEventLimit:              500,
	}
}

func (p *QuerySettings) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type XChatCursor struct {
	CursorId        string `url:"cursor_id,omitempty"`
	GraphSnapshotId string `url:"graph_snapshot_id,omitempty"`
}

func (p *XChatCursor) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type GetInitialXChatPageQueryVariables struct {
	MaxLocalSequenceId string         `url:"max_local_sequence_id,omitempty"`
	QuerySettings      *QuerySettings `url:"query_settings"`
	MessagePullVersion *int           `url:"message_pull_version,omitempty"`
	ContinueCursor     *XChatCursor   `url:"continue_cursor,omitempty"`
}

func NewInitialXChatPageQueryVariables(
	maxLocalSequenceId string,
) *GetInitialXChatPageQueryVariables {
	return &GetInitialXChatPageQueryVariables{
		MaxLocalSequenceId: maxLocalSequenceId,
		QuerySettings:      DefaultQuerySettings(),
		MessagePullVersion: &XCHAT_MESSAGE_PULL_VERSION,
	}
}

func (p *GetInitialXChatPageQueryVariables) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}

type GetInboxPageRequestQueryVariables struct {
	ContinueCursor *XChatCursor   `url:"continue_cursor,omitempty"`
	QuerySettings  *QuerySettings `url:"query_settings"`
}

func NewInboxPageRequestQueryVariables(
	cursor *XChatCursor,
) *GetInboxPageRequestQueryVariables {
	return &GetInboxPageRequestQueryVariables{
		ContinueCursor: cursor,
		QuerySettings:  DefaultQuerySettings(),
	}
}

func (p *GetInboxPageRequestQueryVariables) Encode() ([]byte, error) {
	encodedQuery, err := p.EncodeJSONQuery()
	if err != nil {
		return nil, err
	}

	return []byte(encodedQuery), nil
}

type GetInboxPageConversationDataQueryVariables struct {
	ConversationID        string `json:"conversation_id"`
	IncludeUserPublicKeys bool   `json:"include_user_public_keys"`
}

func NewInboxPageConversationDataQueryVariables(conversationID string, includeKeys bool) *GetInboxPageConversationDataQueryVariables {
	return &GetInboxPageConversationDataQueryVariables{
		ConversationID:        conversationID,
		IncludeUserPublicKeys: includeKeys,
	}
}

// Encode encodes the variables into a form body with a single "variables" JSON field,
// matching how other XChat GraphQL endpoints are called.
func (p *GetInboxPageConversationDataQueryVariables) Encode() ([]byte, error) {
	jsonVars, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("variables", string(jsonVars))
	return []byte(values.Encode()), nil
}

type GetUsersByIdsForXChatVariables struct {
	IDs []string `json:"ids"`
}

func NewGetUsersByIdsForXChatVariables(ids []string) *GetUsersByIdsForXChatVariables {
	return &GetUsersByIdsForXChatVariables{
		IDs: ids,
	}
}

func (p *GetUsersByIdsForXChatVariables) Encode() ([]byte, error) {
	jsonVars, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("variables", string(jsonVars))
	return []byte(values.Encode()), nil
}

// Encode encodes the variables into a form body with a single "variables" JSON field,
// matching how other XChat GraphQL endpoints are called.
func (p *GetConversationPageQueryVariables) Encode() ([]byte, error) {
	jsonVars, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("variables", string(jsonVars))
	return []byte(values.Encode()), nil
}

type GetPublicKeysQueryVariables struct {
	IDs                   []string `json:"ids"`
	IncludeJuiceboxTokens bool     `json:"include_juicebox_tokens"`
}

func NewGetPublicKeysQueryVariables(ids []string) *GetPublicKeysQueryVariables {
	return &GetPublicKeysQueryVariables{
		IDs:                   ids,
		IncludeJuiceboxTokens: true,
	}
}

func (p *GetPublicKeysQueryVariables) Encode() ([]byte, error) {
	jsonVars, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("variables", string(jsonVars))
	return []byte(values.Encode()), nil
}
