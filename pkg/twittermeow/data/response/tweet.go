package response

// FetchPostQueryResponse represents the root response from the fetchPostQuery GraphQL endpoint
type FetchPostQueryResponse struct {
	Data FetchPostQueryData `json:"data"`
}

// FetchPostQueryData contains the top-level data from the GraphQL response
type FetchPostQueryData struct {
	PostResult PostResult `json:"tweet_result_by_rest_id"`
}

// PostResult contains the actual tweet result or error
type PostResult struct {
	Result TweetResult `json:"result"`
}

// TweetResult is the main tweet data structure
type TweetResult struct {
	Typename string       `json:"__typename"`
	ID       string       `json:"id"`
	RestID   string       `json:"rest_id"`
	Core     TweetCore    `json:"core"`
	Legacy   TweetLegacy  `json:"legacy"`
	EditInfo *EditInfo    `json:"edit_control,omitempty"`
}

// TweetCore contains user information for the tweet
type TweetCore struct {
	UserResults UserResults `json:"user_results"`
}

// UserResults wraps the user data
type UserResults struct {
	Result UserResult `json:"result"`
}

// UserResult contains the actual user data
type UserResult struct {
	Typename string     `json:"__typename"`
	ID       string     `json:"id"`
	RestID   string     `json:"rest_id"`
	Legacy   UserLegacy `json:"legacy"`
}

// UserLegacy contains user profile information
type UserLegacy struct {
	Name               string `json:"name"`
	ScreenName         string `json:"screen_name"`
	ProfileImageURL    string `json:"profile_image_url_https"`
	FollowersCount     int    `json:"followers_count,omitempty"`
	FollowingCount     int    `json:"following_count,omitempty"`
	Description        string `json:"description,omitempty"`
	VerifiedType       string `json:"verified_type,omitempty"`
	IsBlueVerified     bool   `json:"is_blue_verified,omitempty"`
}

// TweetLegacy contains the tweet content and metadata
type TweetLegacy struct {
	FullText          string             `json:"full_text"`
	DisplayTextRange  [2]int             `json:"display_text_range"`
	CreatedAt         string             `json:"created_at"`
	ConversationIDStr string             `json:"conversation_id_str,omitempty"`
	ReplyCount        int                `json:"reply_count,omitempty"`
	RetweetCount      int                `json:"retweet_count,omitempty"`
	FavoriteCount     int                `json:"favorite_count,omitempty"`
	PossiblySensitive bool               `json:"possibly_sensitive,omitempty"`
	ExtendedEntities  *ExtendedEntities  `json:"extended_entities,omitempty"`
	Entities          *Entities          `json:"entities,omitempty"`
	QuotedStatusIDStr string             `json:"quoted_status_id_str,omitempty"`
}

// ExtendedEntities contains media information for the tweet
type ExtendedEntities struct {
	Media []MediaEntity `json:"media"`
}

// Entities contains various entity types (URLs, mentions, hashtags)
type Entities struct {
	URLs     []URLEntity     `json:"urls,omitempty"`
	Hashtags []HashtagEntity `json:"hashtags,omitempty"`
	UserMentions []UserMentionEntity `json:"user_mentions,omitempty"`
	Media    []MediaEntity   `json:"media,omitempty"`
}

// MediaEntity represents a media attachment (photo, video, gif)
type MediaEntity struct {
	Type          string        `json:"type"` // "photo", "video", "animated_gif"
	MediaURLHTTPS string        `json:"media_url_https"`
	OriginalInfo  DimensionInfo `json:"original_info"`
	VideoInfo     *VideoInfo    `json:"video_info,omitempty"`
	IDStr         string        `json:"id_str,omitempty"`
	DisplayURL    string        `json:"display_url,omitempty"`
	ExpandedURL   string        `json:"expanded_url,omitempty"`
}

// DimensionInfo contains width and height information
type DimensionInfo struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// VideoInfo contains information about video variants
type VideoInfo struct {
	DurationMillis int            `json:"duration_millis,omitempty"`
	AspectRatio    [2]int         `json:"aspect_ratio,omitempty"`
	Variants       []VideoVariant `json:"variants,omitempty"`
}

// VideoVariant represents a specific video quality/format variant
type VideoVariant struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Bitrate     int    `json:"bitrate,omitempty"`
}

// URLEntity represents a URL in the tweet
type URLEntity struct {
	URL         string `json:"url"`
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     [2]int `json:"indices"`
}

// HashtagEntity represents a hashtag in the tweet
type HashtagEntity struct {
	Text    string `json:"text"`
	Indices [2]int `json:"indices"`
}

// UserMentionEntity represents a user mention in the tweet
type UserMentionEntity struct {
	ScreenName string `json:"screen_name"`
	Name       string `json:"name"`
	ID         int64  `json:"id"`
	Indices    [2]int `json:"indices"`
}

// EditInfo contains information about tweet edits
type EditInfo struct {
	EditTweetIDs           []string `json:"edit_tweet_ids"`
	EditableUntil          string   `json:"editable_until"`
	IsEditEligible         bool     `json:"is_edit_eligible"`
	EditsRemaining         int      `json:"edits_remaining"`
}
