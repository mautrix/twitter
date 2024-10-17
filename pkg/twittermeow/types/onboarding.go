package types

type TaskResponse struct {
	FlowToken string     `json:"flow_token,omitempty"`
	Status    string     `json:"status,omitempty"`
	Subtasks  []Subtasks `json:"subtasks,omitempty"`
}
type ScribeConfig struct {
	Action    string `json:"action,omitempty"`
	Component string `json:"component,omitempty"`
	Page      string `json:"page,omitempty"`
	Section   string `json:"section,omitempty"`
}
type Callbacks struct {
	Endpoint     string       `json:"endpoint,omitempty"`
	Trigger      string       `json:"trigger,omitempty"`
	ScribeConfig ScribeConfig `json:"scribe_config,omitempty"`
}

type Entities struct {
	FromIndex      int            `json:"from_index,omitempty"`
	NavigationLink NavigationLink `json:"navigation_link,omitempty"`
	ToIndex        int            `json:"to_index,omitempty"`
}
type DetailText struct {
	Entities []Entities `json:"entities,omitempty"`
	Text     string     `json:"text,omitempty"`
}
type PrimaryText struct {
	Entities []any  `json:"entities,omitempty"`
	Text     string `json:"text,omitempty"`
}
type Header struct {
	PrimaryText PrimaryText `json:"primary_text,omitempty"`
}
type SecondaryText struct {
	Entities []Entities `json:"entities,omitempty"`
	Text     string     `json:"text,omitempty"`
}
type Icon struct {
	Icon string `json:"icon,omitempty"`
}
type NavigationLink struct {
	Label     string `json:"label,omitempty"`
	LinkID    string `json:"link_id,omitempty"`
	LinkType  string `json:"link_type,omitempty"`
	SubtaskID string `json:"subtask_id,omitempty"`
}
type Button struct {
	Icon           Icon           `json:"icon,omitempty"`
	NavigationLink NavigationLink `json:"navigation_link,omitempty"`
	PreferredSize  string         `json:"preferred_size,omitempty"`
	Style          string         `json:"style,omitempty"`
}
type Message struct {
	Entities []any  `json:"entities,omitempty"`
	Text     string `json:"text,omitempty"`
}
type InAppNotification struct {
	Message    Message  `json:"message,omitempty"`
	NextLink   NextLink `json:"next_link,omitempty"`
	WaitTimeMs int      `json:"wait_time_ms,omitempty"`
}
type Link struct {
	LinkID    string `json:"link_id,omitempty"`
	LinkType  string `json:"link_type,omitempty"`
	SubtaskID string `json:"subtask_id,omitempty"`
}
type OpenLink struct {
	Link                   Link   `json:"link,omitempty"`
	OnboardingCallbackPath string `json:"onboarding_callback_path,omitempty"`
}
type Text struct {
	Entities []any  `json:"entities,omitempty"`
	Text     string `json:"text,omitempty"`
}
type ProgressIndication struct {
	Text Text `json:"text,omitempty"`
}
type EmailNextLink struct {
	Label     string `json:"label,omitempty"`
	LinkID    string `json:"link_id,omitempty"`
	LinkType  string `json:"link_type,omitempty"`
	SubtaskID string `json:"subtask_id,omitempty"`
}
type JsInstrumentation struct {
	URL string `json:"url,omitempty"`
}
type SignUp struct {
	AllowedIdentifiers  string            `json:"allowed_identifiers,omitempty"`
	BirthdayExplanation string            `json:"birthday_explanation,omitempty"`
	BirthdayHint        string            `json:"birthday_hint,omitempty"`
	BirthdayType        string            `json:"birthday_type,omitempty"`
	EmailHint           string            `json:"email_hint,omitempty"`
	EmailNextLink       EmailNextLink     `json:"email_next_link,omitempty"`
	IgnoreBirthday      bool              `json:"ignore_birthday,omitempty"`
	JsInstrumentation   JsInstrumentation `json:"js_instrumentation,omitempty"`
	NameHint            string            `json:"name_hint,omitempty"`
	NextLink            NextLink          `json:"next_link,omitempty"`
	PasswordHint        string            `json:"password_hint,omitempty"`
	PhoneEmailHint      string            `json:"phone_email_hint,omitempty"`
	PhoneHint           string            `json:"phone_hint,omitempty"`
	PrimaryText         string            `json:"primary_text,omitempty"`
	UseDevicePrefill    bool              `json:"use_device_prefill,omitempty"`
	UseEmailText        string            `json:"use_email_text,omitempty"`
	UsePhoneText        string            `json:"use_phone_text,omitempty"`
}
type BooleanData struct {
	InitialValue bool `json:"initial_value,omitempty"`
}
type ValueData struct {
	BooleanData BooleanData `json:"boolean_data,omitempty"`
}
type Settings struct {
	PrimaryText     PrimaryText `json:"primary_text,omitempty"`
	ValueType       string      `json:"value_type,omitempty"`
	ValueData       ValueData   `json:"value_data,omitempty"`
	ValueIdentifier string      `json:"value_identifier,omitempty"`
}
type SettingsList0 struct {
	DetailText  DetailText  `json:"detail_text,omitempty"`
	NextLink    NextLink    `json:"next_link,omitempty"`
	PrimaryText PrimaryText `json:"primary_text,omitempty"`
	Settings    []Settings  `json:"settings,omitempty"`
}
type SubtaskDataReference struct {
	Key       string `json:"key,omitempty"`
	SubtaskID string `json:"subtask_id,omitempty"`
}
type Birthday struct {
	SubtaskDataReference SubtaskDataReference `json:"subtask_data_reference,omitempty"`
}
type SubtaskNavigationContext struct {
	Action string `json:"action,omitempty"`
}
type BirthdayEditLink struct {
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
}
type BirthdayRequirement struct {
	Day   int `json:"day,omitempty"`
	Month int `json:"month,omitempty"`
	Year  int `json:"year,omitempty"`
}
type Email struct {
	SubtaskDataReference SubtaskDataReference `json:"subtask_data_reference,omitempty"`
}
type EmailEditLink struct {
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
}
type InvalidBirthdayLink struct {
	Label                string `json:"label,omitempty"`
	LinkID               string `json:"link_id,omitempty"`
	LinkType             string `json:"link_type,omitempty"`
	SubtaskID            string `json:"subtask_id,omitempty"`
	SuppressClientEvents bool   `json:"suppress_client_events,omitempty"`
}
type Name struct {
	SubtaskDataReference SubtaskDataReference `json:"subtask_data_reference,omitempty"`
}
type NameEditLink struct {
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
}
type PhoneEditLink struct {
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
}
type PhoneNextLink struct {
	Label     string `json:"label,omitempty"`
	LinkID    string `json:"link_id,omitempty"`
	LinkType  string `json:"link_type,omitempty"`
	SubtaskID string `json:"subtask_id,omitempty"`
}
type PhoneNumber struct {
	SubtaskDataReference SubtaskDataReference `json:"subtask_data_reference,omitempty"`
}
type SignUpReview struct {
	Birthday            Birthday            `json:"birthday,omitempty"`
	BirthdayEditLink    BirthdayEditLink    `json:"birthday_edit_link,omitempty"`
	BirthdayHint        string              `json:"birthday_hint,omitempty"`
	BirthdayRequirement BirthdayRequirement `json:"birthday_requirement,omitempty"`
	BirthdayType        string              `json:"birthday_type,omitempty"`
	DetailText          DetailText          `json:"detail_text,omitempty"`
	Email               Email               `json:"email,omitempty"`
	EmailEditLink       EmailEditLink       `json:"email_edit_link,omitempty"`
	EmailHint           string              `json:"email_hint,omitempty"`
	EmailNextLink       EmailNextLink       `json:"email_next_link,omitempty"`
	IgnoreBirthday      bool                `json:"ignore_birthday,omitempty"`
	InvalidBirthdayLink InvalidBirthdayLink `json:"invalid_birthday_link,omitempty"`
	Name                Name                `json:"name,omitempty"`
	NameEditLink        NameEditLink        `json:"name_edit_link,omitempty"`
	NameHint            string              `json:"name_hint,omitempty"`
	PhoneEditLink       PhoneEditLink       `json:"phone_edit_link,omitempty"`
	PhoneHint           string              `json:"phone_hint,omitempty"`
	PhoneNextLink       PhoneNextLink       `json:"phone_next_link,omitempty"`
	PhoneNumber         PhoneNumber         `json:"phone_number,omitempty"`
	PrimaryText         string              `json:"primary_text,omitempty"`
}
type DiscoverableByEmailDetailText struct {
	Entities []any  `json:"entities,omitempty"`
	Text     string `json:"text,omitempty"`
}
type DiscoverableByPhoneDetailText struct {
	Entities []any  `json:"entities,omitempty"`
	Text     string `json:"text,omitempty"`
}
type PrivacyOptions struct {
	DiscoverableByEmail           bool                          `json:"discoverable_by_email,omitempty"`
	DiscoverableByEmailDetailText DiscoverableByEmailDetailText `json:"discoverable_by_email_detail_text,omitempty"`
	DiscoverableByEmailLabel      string                        `json:"discoverable_by_email_label,omitempty"`
	DiscoverableByPhone           bool                          `json:"discoverable_by_phone,omitempty"`
	DiscoverableByPhoneDetailText DiscoverableByPhoneDetailText `json:"discoverable_by_phone_detail_text,omitempty"`
	DiscoverableByPhoneLabel      string                        `json:"discoverable_by_phone_label,omitempty"`
	NextLink                      NextLink                      `json:"next_link,omitempty"`
	PrimaryText                   string                        `json:"primary_text,omitempty"`
	SecondaryText                 string                        `json:"secondary_text,omitempty"`
}
type CancelLink struct {
	Label                    string                   `json:"label,omitempty"`
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
}
type AlertDialog struct {
	CancelLink    CancelLink    `json:"cancel_link,omitempty"`
	NextLink      NextLink      `json:"next_link,omitempty"`
	PrimaryText   PrimaryText   `json:"primary_text,omitempty"`
	SecondaryText SecondaryText `json:"secondary_text,omitempty"`
}
type PrimaryActionLinks struct {
	Label                    string                   `json:"label,omitempty"`
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
	IsDestructive            bool                     `json:"is_destructive,omitempty"`
}
type MenuDialog struct {
	CancelLink         CancelLink           `json:"cancel_link,omitempty"`
	DismissLink        DismissLink          `json:"dismiss_link,omitempty"`
	PrimaryActionLinks []PrimaryActionLinks `json:"primary_action_links,omitempty"`
	PrimaryText        PrimaryText          `json:"primary_text,omitempty"`
}
type Subtext struct {
	Entities []any  `json:"entities,omitempty"`
	Text     string `json:"text,omitempty"`
}
type Choices struct {
	Icon    Icon    `json:"icon,omitempty"`
	ID      string  `json:"id,omitempty"`
	Subtext Subtext `json:"subtext,omitempty"`
	Text    Text    `json:"text,omitempty"`
}
type NextLinkOptions struct {
	MinEnableCount int `json:"min_enable_count,omitempty"`
}
type SkipLink struct {
	Label                    string                   `json:"label,omitempty"`
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
}
type ChoiceSelection struct {
	Choices         []Choices       `json:"choices,omitempty"`
	DetailText      DetailText      `json:"detail_text,omitempty"`
	NextLink        NextLink        `json:"next_link,omitempty"`
	NextLinkOptions NextLinkOptions `json:"next_link_options,omitempty"`
	PrimaryText     PrimaryText     `json:"primary_text,omitempty"`
	SecondaryText   SecondaryText   `json:"secondary_text,omitempty"`
	Sections        []any           `json:"sections,omitempty"`
	SelectedChoices []string        `json:"selected_choices,omitempty"`
	SelectionType   string          `json:"selection_type,omitempty"`
	SkipLink        SkipLink        `json:"skip_link,omitempty"`
	Style           string          `json:"style,omitempty"`
}
type PhoneVerification struct {
	AutoVerifyHintText string        `json:"auto_verify_hint_text,omitempty"`
	CancelLink         CancelLink    `json:"cancel_link,omitempty"`
	DetailText         DetailText    `json:"detail_text,omitempty"`
	HintText           string        `json:"hint_text,omitempty"`
	NextLink           NextLink      `json:"next_link,omitempty"`
	PhoneNumber        PhoneNumber   `json:"phone_number,omitempty"`
	PrimaryText        PrimaryText   `json:"primary_text,omitempty"`
	SecondaryText      SecondaryText `json:"secondary_text,omitempty"`
	SendViaVoice       bool          `json:"send_via_voice,omitempty"`
}
type EmailVerification struct {
	CancelLink                       CancelLink    `json:"cancel_link,omitempty"`
	DetailText                       DetailText    `json:"detail_text,omitempty"`
	Email                            Email         `json:"email,omitempty"`
	HintText                         string        `json:"hint_text,omitempty"`
	Name                             Name          `json:"name,omitempty"`
	NextLink                         NextLink      `json:"next_link,omitempty"`
	PrimaryText                      PrimaryText   `json:"primary_text,omitempty"`
	SecondaryText                    SecondaryText `json:"secondary_text,omitempty"`
	VerificationStatusPollingEnabled bool          `json:"verification_status_polling_enabled,omitempty"`
}
type DismissLink struct {
	IsDestructive            bool                     `json:"is_destructive,omitempty"`
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
	SuppressClientEvents     bool                     `json:"suppress_client_events,omitempty"`
}
type NextLink struct {
	IsDestructive            bool                     `json:"is_destructive,omitempty"`
	Label                    string                   `json:"label,omitempty"`
	LinkID                   string                   `json:"link_id,omitempty"`
	LinkType                 string                   `json:"link_type,omitempty"`
	SubtaskID                string                   `json:"subtask_id,omitempty"`
	SubtaskNavigationContext SubtaskNavigationContext `json:"subtask_navigation_context,omitempty"`
	SuppressClientEvents     bool                     `json:"suppress_client_events,omitempty"`
}
type AlertDialogSuppressClientEvents struct {
	DismissLink DismissLink `json:"dismiss_link,omitempty"`
	NextLink    NextLink    `json:"next_link,omitempty"`
	PrimaryText PrimaryText `json:"primary_text,omitempty"`
}
type Subtasks struct {
	Callbacks             []Callbacks        `json:"callbacks,omitempty"`
	SubtaskBackNavigation string             `json:"subtask_back_navigation,omitempty"`
	SubtaskID             string             `json:"subtask_id,omitempty"`
	InAppNotification     InAppNotification  `json:"in_app_notification,omitempty"`
	OpenLink              OpenLink           `json:"open_link,omitempty"`
	ProgressIndication    ProgressIndication `json:"progress_indication,omitempty"`
	SignUp                SignUp             `json:"sign_up,omitempty"`
	SettingsList          SettingsList0      `json:"settings_list,omitempty"`
	//SettingsList1                   SettingsList1                   `json:"settings_list,omitempty"`
	SignUpReview                    SignUpReview                    `json:"sign_up_review,omitempty"`
	PrivacyOptions                  PrivacyOptions                  `json:"privacy_options,omitempty"`
	AlertDialog                     AlertDialog                     `json:"alert_dialog,omitempty"`
	MenuDialog                      MenuDialog                      `json:"menu_dialog,omitempty"`
	ChoiceSelection                 ChoiceSelection                 `json:"choice_selection,omitempty"`
	PhoneVerification               PhoneVerification               `json:"phone_verification,omitempty"`
	EmailVerification               EmailVerification               `json:"email_verification,omitempty"`
	AlertDialogSuppressClientEvents AlertDialogSuppressClientEvents `json:"alert_dialog_suppress_client_events,omitempty"`
	SelectAvatar                    *SelectAvatar                   `json:"select_avatar,omitempty"`
}

type CallbackPayload struct {
	Product     string `json:"product,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
	Params      string `json:"params,omitempty"`
	TimestampMs int64  `json:"timestampMs,omitempty"`
}

type UploadMediaStartResponse struct {
	MediaID          int64  `json:"media_id,omitempty"`
	MediaIDString    string `json:"media_id_string,omitempty"`
	Size             int    `json:"size,omitempty"`
	ExpiresAfterSecs int    `json:"expires_after_secs,omitempty"`
	Image            struct {
		ImageType string `json:"image_type,omitempty"`
		W         int    `json:"w,omitempty"`
		H         int    `json:"h,omitempty"`
	} `json:"image,omitempty"`
}

type SelectAvatar struct {
	NextLink struct {
		Callbacks []struct {
			Endpoint     string `json:"endpoint,omitempty"`
			Trigger      string `json:"trigger,omitempty"`
			ScribeConfig struct {
				Action    string `json:"action,omitempty"`
				Component string `json:"component,omitempty"`
				Element   string `json:"element,omitempty"`
				Page      string `json:"page,omitempty"`
				Section   string `json:"section,omitempty"`
			} `json:"scribe_config,omitempty"`
		} `json:"callbacks,omitempty"`
		Label     string `json:"label,omitempty"`
		LinkID    string `json:"link_id,omitempty"`
		LinkType  string `json:"link_type,omitempty"`
		SubtaskID string `json:"subtask_id,omitempty"`
	} `json:"next_link,omitempty"`
	PrimaryText struct {
		Entities []any  `json:"entities,omitempty"`
		Text     string `json:"text,omitempty"`
	} `json:"primary_text,omitempty"`
	SecondaryText struct {
		Entities []any  `json:"entities,omitempty"`
		Text     string `json:"text,omitempty"`
	} `json:"secondary_text,omitempty"`
	SkipLink struct {
		Label     string `json:"label,omitempty"`
		LinkID    string `json:"link_id,omitempty"`
		LinkType  string `json:"link_type,omitempty"`
		SubtaskID string `json:"subtask_id,omitempty"`
	} `json:"skip_link,omitempty"`
}

type UserRecommendationsOnboardingResponse struct {
	GlobalObjects struct {
		Tweets struct {
		} `json:"tweets,omitempty"`
		Users   map[string]*RecommendedUser `json:"users,omitempty"`
		Moments struct {
		} `json:"moments,omitempty"`
		Cards struct {
		} `json:"cards,omitempty"`
		Places struct {
		} `json:"places,omitempty"`
		Media struct {
		} `json:"media,omitempty"`
		Broadcasts struct {
		} `json:"broadcasts,omitempty"`
		Topics struct {
		} `json:"topics,omitempty"`
		Lists struct {
		} `json:"lists,omitempty"`
	} `json:"globalObjects,omitempty"`
	Timeline struct {
		ID           string `json:"id,omitempty"`
		Instructions []struct {
			AddEntries struct {
				Entries []struct {
					EntryID   string `json:"entryId,omitempty"`
					SortIndex string `json:"sortIndex,omitempty"`
					Content   struct {
						TimelineModule struct {
							Items []struct {
								EntryID string `json:"entryId,omitempty"`
								Item    struct {
									Content struct {
										User struct {
											ID                     string `json:"id,omitempty"`
											DisplayType            string `json:"displayType,omitempty"`
											EnableReactiveBlending bool   `json:"enableReactiveBlending,omitempty"`
										} `json:"user,omitempty"`
									} `json:"content,omitempty"`
									ClientEventInfo struct {
										Component string `json:"component,omitempty"`
										Element   string `json:"element,omitempty"`
										Details   struct {
											TimelinesDetails struct {
												SourceData string `json:"sourceData,omitempty"`
											} `json:"timelinesDetails,omitempty"`
										} `json:"details,omitempty"`
										Action string `json:"action,omitempty"`
									} `json:"clientEventInfo,omitempty"`
								} `json:"item,omitempty"`
							} `json:"items,omitempty"`
							DisplayType string `json:"displayType,omitempty"`
							Header      struct {
								Text        string `json:"text,omitempty"`
								DisplayType string `json:"displayType,omitempty"`
							} `json:"header,omitempty"`
						} `json:"timelineModule,omitempty"`
					} `json:"content,omitempty"`
				} `json:"entries,omitempty"`
			} `json:"addEntries,omitempty"`
		} `json:"instructions,omitempty"`
	} `json:"timeline,omitempty"`
}

type RecommendedUser struct {
	ID          int    `json:"id,omitempty"`
	IDStr       string `json:"id_str,omitempty"`
	Name        string `json:"name,omitempty"`
	ScreenName  string `json:"screen_name,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Entities    struct {
		URL struct {
			Urls []struct {
				URL         string `json:"url,omitempty"`
				ExpandedURL string `json:"expanded_url,omitempty"`
				DisplayURL  string `json:"display_url,omitempty"`
				Indices     []int  `json:"indices,omitempty"`
			} `json:"urls,omitempty"`
		} `json:"url,omitempty"`
		Description struct {
			Urls []any `json:"urls,omitempty"`
		} `json:"description,omitempty"`
	} `json:"entities,omitempty"`
	Protected                      bool   `json:"protected,omitempty"`
	FollowersCount                 int    `json:"followers_count,omitempty"`
	FriendsCount                   int    `json:"friends_count,omitempty"`
	ListedCount                    int    `json:"listed_count,omitempty"`
	CreatedAt                      string `json:"created_at,omitempty"`
	FavouritesCount                int    `json:"favourites_count,omitempty"`
	UtcOffset                      any    `json:"utc_offset,omitempty"`
	TimeZone                       any    `json:"time_zone,omitempty"`
	GeoEnabled                     bool   `json:"geo_enabled,omitempty"`
	Verified                       bool   `json:"verified,omitempty"`
	StatusesCount                  int    `json:"statuses_count,omitempty"`
	Lang                           any    `json:"lang,omitempty"`
	ContributorsEnabled            bool   `json:"contributors_enabled,omitempty"`
	IsTranslator                   bool   `json:"is_translator,omitempty"`
	IsTranslationEnabled           bool   `json:"is_translation_enabled,omitempty"`
	ProfileBackgroundColor         string `json:"profile_background_color,omitempty"`
	ProfileBackgroundImageURL      string `json:"profile_background_image_url,omitempty"`
	ProfileBackgroundImageURLHTTPS string `json:"profile_background_image_url_https,omitempty"`
	ProfileBackgroundTile          bool   `json:"profile_background_tile,omitempty"`
	ProfileImageURL                string `json:"profile_image_url,omitempty"`
	ProfileImageURLHTTPS           string `json:"profile_image_url_https,omitempty"`
	ProfileBannerURL               string `json:"profile_banner_url,omitempty"`
	ProfileLinkColor               string `json:"profile_link_color,omitempty"`
	ProfileSidebarBorderColor      string `json:"profile_sidebar_border_color,omitempty"`
	ProfileSidebarFillColor        string `json:"profile_sidebar_fill_color,omitempty"`
	ProfileTextColor               string `json:"profile_text_color,omitempty"`
	ProfileUseBackgroundImage      bool   `json:"profile_use_background_image,omitempty"`
	HasExtendedProfile             bool   `json:"has_extended_profile,omitempty"`
	DefaultProfile                 bool   `json:"default_profile,omitempty"`
	DefaultProfileImage            bool   `json:"default_profile_image,omitempty"`
	CanMediaTag                    any    `json:"can_media_tag,omitempty"`
	Following                      bool   `json:"following,omitempty"`
	FollowRequestSent              bool   `json:"follow_request_sent,omitempty"`
	Notifications                  bool   `json:"notifications,omitempty"`
	Muting                         bool   `json:"muting,omitempty"`
	Blocking                       bool   `json:"blocking,omitempty"`
	BlockedBy                      bool   `json:"blocked_by,omitempty"`
	TranslatorType                 string `json:"translator_type,omitempty"`
	WithheldInCountries            []any  `json:"withheld_in_countries,omitempty"`
	FollowedBy                     bool   `json:"followed_by,omitempty"`
	ExtIsBlueVerified              bool   `json:"ext_is_blue_verified,omitempty"`
}
