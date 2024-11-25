package payload

import "github.com/google/go-querystring/query"

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
	ActiveConversationID                string      `url:"active_conversation_id,omitempty"`
	Cursor                              string      `url:"cursor,omitempty"`
	Count                               int         `url:"count,omitempty"`
	Context                             ContextInfo `url:"context,omitempty"`
	MaxID                               string      `url:"max_id,omitempty"` // when fetching messages, this is the message id
	MinID                               string      `url:"min_id,omitempty"`
	NSFWFilteringEnabled                bool        `url:"nsfw_filtering_enabled"`
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
	IncludeExtEditControl               bool        `url:"include_ext_edit_control"`
	IncludeExtBusinessAffiliationsLabel bool        `url:"include_ext_business_affiliations_label"`
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
		KRSRegistrationEnabled:              true,
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
		IncludeExtEditControl:               true,
		IncludeExtBusinessAffiliationsLabel: true,
		Ext:                                 "mediaColor,altText,mediaStats,highlightedLabel,voiceInfo,birdwatchPivot,superFollowMetadata,unmentionInfo,editControl,article",
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

type MediaType string

const (
	MEDIA_TYPE_IMAGE_JPEG MediaType = "image/jpeg"
	MEDIA_TYPE_IMAGE_GIF  MediaType = "image/gif"
	MEDIA_TYPE_VIDEO_MP4  MediaType = "video/mp4"
)

type UploadMediaQuery struct {
	Command         string        `url:"command,omitempty"`
	TotalBytes      int           `url:"total_bytes,omitempty"`
	SourceURL       string        `url:"source_url,omitempty"`
	MediaID         string        `url:"media_id,omitempty"`
	VideoDurationMS float32       `url:"video_duration_ms,omitempty"`
	OriginalMD5     string        `url:"original_md5,omitempty"`
	SegmentIndex    int           `url:"segment_index,omitempty"`
	MediaType       MediaType     `url:"media_type,omitempty"`
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

type EditDirectMessagePayload struct {
	ConversationID string `url:"conversation_id,omitempty"`
	RequestID      string `url:"request_id,omitempty"`
	DMID           string `url:"dm_id,omitempty"` // used to specify a message, specifically for editing
	Text           string `url:"text"`
}

func (p *EditDirectMessagePayload) Encode() ([]byte, error) {
	values, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	return []byte(values.Encode()), nil
}
