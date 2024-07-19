package response

import "go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"

type AccountSettingsResponse struct {
	Protected                                       bool             `json:"protected,omitempty"`
	ScreenName                                      string           `json:"screen_name,omitempty"`
	AlwaysUseHTTPS                                  bool             `json:"always_use_https,omitempty"`
	UseCookiePersonalization                        bool             `json:"use_cookie_personalization,omitempty"`
	SleepTime                                       SleepTime        `json:"sleep_time,omitempty"`
	GeoEnabled                                      bool             `json:"geo_enabled,omitempty"`
	Language                                        string           `json:"language,omitempty"`
	DiscoverableByEmail                             bool             `json:"discoverable_by_email,omitempty"`
	DiscoverableByMobilePhone                       bool             `json:"discoverable_by_mobile_phone,omitempty"`
	DisplaySensitiveMedia                           bool             `json:"display_sensitive_media,omitempty"`
	PersonalizedTrends                              bool             `json:"personalized_trends,omitempty"`
	AllowMediaTagging                               string           `json:"allow_media_tagging,omitempty"`
	AllowContributorRequest                         string           `json:"allow_contributor_request,omitempty"`
	AllowAdsPersonalization                         bool             `json:"allow_ads_personalization,omitempty"`
	AllowLoggedOutDevicePersonalization             bool             `json:"allow_logged_out_device_personalization,omitempty"`
	AllowLocationHistoryPersonalization             bool             `json:"allow_location_history_personalization,omitempty"`
	AllowSharingDataForThirdPartyPersonalization    bool             `json:"allow_sharing_data_for_third_party_personalization,omitempty"`
	AllowDmsFrom                                    string           `json:"allow_dms_from,omitempty"`
	AlwaysAllowDmsFromSubscribers                   any              `json:"always_allow_dms_from_subscribers,omitempty"`
	AllowDmGroupsFrom                               string           `json:"allow_dm_groups_from,omitempty"`
	TranslatorType                                  string           `json:"translator_type,omitempty"`
	CountryCode                                     string           `json:"country_code,omitempty"`
	NsfwUser                                        bool             `json:"nsfw_user,omitempty"`
	NsfwAdmin                                       bool             `json:"nsfw_admin,omitempty"`
	RankedTimelineSetting                           any              `json:"ranked_timeline_setting,omitempty"`
	RankedTimelineEligible                          any              `json:"ranked_timeline_eligible,omitempty"`
	AddressBookLiveSyncEnabled                      bool             `json:"address_book_live_sync_enabled,omitempty"`
	UniversalQualityFilteringEnabled                string           `json:"universal_quality_filtering_enabled,omitempty"`
	DmReceiptSetting                                string           `json:"dm_receipt_setting,omitempty"`
	AltTextComposeEnabled                           any              `json:"alt_text_compose_enabled,omitempty"`
	MentionFilter                                   string           `json:"mention_filter,omitempty"`
	AllowAuthenticatedPeriscopeRequests             bool             `json:"allow_authenticated_periscope_requests,omitempty"`
	ProtectPasswordReset                            bool             `json:"protect_password_reset,omitempty"`
	RequirePasswordLogin                            bool             `json:"require_password_login,omitempty"`
	RequiresLoginVerification                       bool             `json:"requires_login_verification,omitempty"`
	ExtSharingAudiospacesListeningDataWithFollowers bool             `json:"ext_sharing_audiospaces_listening_data_with_followers,omitempty"`
	Ext                                             Ext              `json:"ext,omitempty"`
	DmQualityFilter                                 string           `json:"dm_quality_filter,omitempty"`
	AutoplayDisabled                                bool             `json:"autoplay_disabled,omitempty"`
	SettingsMetadata                                SettingsMetadata `json:"settings_metadata,omitempty"`
}
type SleepTime struct {
	Enabled   bool `json:"enabled,omitempty"`
	EndTime   any  `json:"end_time,omitempty"`
	StartTime any  `json:"start_time,omitempty"`
}
type Ok struct {
	SsoIDHash   string `json:"ssoIdHash,omitempty"`
	SsoProvider string `json:"ssoProvider,omitempty"`
}
type R struct {
	Ok []Ok `json:"ok,omitempty"`
}
type SsoConnections struct {
	R   R   `json:"r,omitempty"`
	TTL int `json:"ttl,omitempty"`
}
type Ext struct {
	SsoConnections SsoConnections `json:"ssoConnections,omitempty"`
}
type SettingsMetadata struct {
	IsEu string `json:"is_eu,omitempty"`
}

type GetDMPermissionsResponse struct {
	Permissions Permissions           `json:"permissions,omitempty"`
	Users       map[string]types.User `json:"users,omitempty"`
}

type PermissionDetails struct {
	CanDm     bool `json:"can_dm,omitempty"`
	ErrorCode int  `json:"error_code,omitempty"`
}

type Permissions struct {
	IDKeys map[string]PermissionDetails `json:"id_keys,omitempty"`
}

func (perms Permissions) GetPermissionsForUser(userId string) *PermissionDetails {
	if user, ok := perms.IDKeys[userId]; ok {
		return &user
	}

	return nil
}
