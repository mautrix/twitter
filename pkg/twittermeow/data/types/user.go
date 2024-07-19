package types

type User struct {
	ID                   int64    `json:"id,omitempty"`
	IDStr                string   `json:"id_str,omitempty"`
	Name                 string   `json:"name,omitempty"`
	ScreenName           string   `json:"screen_name,omitempty"`
	ProfileImageURL      string   `json:"profile_image_url,omitempty"`
	ProfileImageURLHTTPS string   `json:"profile_image_url_https,omitempty"`
	Following            bool     `json:"following,omitempty"`
	FollowRequestSent    bool     `json:"follow_request_sent,omitempty"`
	Description          string   `json:"description,omitempty"`
	Entities             Entities `json:"entities,omitempty"`
	Verified             bool     `json:"verified,omitempty"`
	IsBlueVerified       bool     `json:"is_blue_verified,omitempty"`
	ExtIsBlueVerified    bool     `json:"ext_is_blue_verified,omitempty"`
	Protected            bool     `json:"protected,omitempty"`
	IsProtected          bool     `json:"is_protected,omitempty"`
	Blocking             bool     `json:"blocking,omitempty"`
	IsBlocked            bool     `json:"is_blocked,omitempty"`
	IsSecretDmAble       bool     `json:"is_secret_dm_able,omitempty"`
	IsDmAble             bool     `json:"is_dm_able,omitempty"`
	SubscribedBy         bool     `json:"subscribed_by,omitempty"`
	CanMediaTag          bool     `json:"can_media_tag,omitempty"`
	CreatedAt            string   `json:"created_at,omitempty"`
	Location             string   `json:"location,omitempty"`
	FriendsCount         int      `json:"friends_count,omitempty"`
	SocialProof          int      `json:"social_proof,omitempty"`
	RoundedScore         int      `json:"rounded_score,omitempty"`
	FollowersCount       int      `json:"followers_count,omitempty"`
	ConnectingUserCount  int      `json:"connecting_user_count,omitempty"`
	ConnectingUserIds    []any    `json:"connecting_user_ids,omitempty"`
	SocialProofsOrdered  []any    `json:"social_proofs_ordered,omitempty"`
	Tokens               []any    `json:"tokens,omitempty"`
	Inline               bool     `json:"inline,omitempty"`
}

type UserEntities struct {
	URL         URL         `json:"url,omitempty"`
	Description Description `json:"description,omitempty"`
}

type URL struct {
	Urls []Urls `json:"urls,omitempty"`
}

type Description struct {
	Urls []Urls `json:"urls,omitempty"`
}
