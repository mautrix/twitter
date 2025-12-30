package response

// XChatGetAuthTokenResponse models the GraphQL response for fetching an XChat auth token.
type XChatGetAuthTokenResponse struct {
	Data struct {
		UserGetXChatAuthToken struct {
			Typename string `json:"__typename,omitempty"`
			Token    string `json:"token,omitempty"`
		} `json:"user_get_x_chat_auth_token"`
	} `json:"data"`
}

// GetUsersByIdsForXChatResponse models the GraphQL response for fetching XChat member/user data.
type GetUsersByIdsForXChatResponse struct {
	Data struct {
		GetMemberResults struct {
			Results []XChatMemberResult `json:"results,omitempty"`
		} `json:"get_member_results"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}

type XChatMemberResult struct {
	Typename      string           `json:"__typename,omitempty"`
	ID            string           `json:"id,omitempty"`
	Status        string           `json:"status,omitempty"`
	MemberResults *XChatUserResult `json:"member_results,omitempty"`
}

// GetConversationPageQueryResponse models the GraphQL response for fetching message history in a conversation.
type GetConversationPageQueryResponse struct {
	Data struct {
		GetConversationPage XChatConversationPage `json:"get_conversation_page"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}

type XChatConversationPage struct {
	Typename                           string   `json:"__typename,omitempty"`
	EncodedMessageEvents               []string `json:"encoded_message_events,omitempty"`
	MissingConversationKeyChangeEvents []string `json:"missing_conversation_key_change_events,omitempty"`
	HasMore                            bool     `json:"has_more,omitempty"`
}

// GetInitialXChatPageQueryResponse models the GraphQL response for fetching the
// initial XChat inbox page.
type GetInitialXChatPageQueryResponse struct {
	Data struct {
		GetInboxPage XChatInboxPage `json:"get_initial_chat_page"`
	} `json:"data"`
}

// GetInboxPageRequestQueryResponse models the GraphQL response for fetching subsequent inbox pages.
type GetInboxPageRequestQueryResponse struct {
	Data struct {
		GetInboxPage XChatInboxPage `json:"get_inbox_page"`
	} `json:"data"`
}

// GetInboxPageConversationDataResponse models fetching a single conversation's data by ID.
type GetInboxPageConversationDataResponse struct {
	Data struct {
		GetInboxPageConversationData struct {
			Typename string         `json:"__typename,omitempty"`
			Data     XChatInboxItem `json:"data"`
		} `json:"get_inbox_page_conversation_data"`
	} `json:"data"`
}

type XChatInboxPage struct {
	Typename           string           `json:"__typename,omitempty"`
	InboxCursor        XChatInboxCursor `json:"inboxCursor"`
	Items              []XChatInboxItem `json:"items,omitempty"`
	HasMessageRequests bool             `json:"has_message_requests,omitempty"`
	MaxUserSequenceID  *string          `json:"max_user_sequence_id,omitempty"`
	MessagePullVersion *int             `json:"message_pull_version,omitempty"`
	Errors             []map[string]any `json:"errors,omitempty"`
	Extensions         map[string]any   `json:"extensions,omitempty"`
}

type XChatInboxCursor struct {
	Typename        string `json:"__typename,omitempty"`
	CursorID        string `json:"cursor_id,omitempty"`
	GraphSnapshotID string `json:"graph_snapshot_id,omitempty"`
	PullFinished    bool   `json:"pull_finished,omitempty"`
}

type XChatInboxItem struct {
	Typename                           string                      `json:"__typename,omitempty"`
	LatestMessageEvents                []string                    `json:"latest_message_events,omitempty"`
	EncodedMessageEvents               []string                    `json:"encoded_message_events,omitempty"`
	ConversationDetail                 XChatConversationDetail     `json:"conversation_detail"`
	LatestConversationKeyChangeEvents  []string                    `json:"latest_conversation_key_change_events,omitempty"`
	LatestNotifiableMessageCreateEvent string                      `json:"latest_notifiable_message_create_event,omitempty"`
	LatestReadEventsPerParticipant     []XChatParticipantReadEvent `json:"latest_read_events_per_participant,omitempty"`
	HasMore                            bool                        `json:"has_more,omitempty"`
}

type XChatConversationDetail struct {
	Typename                           string              `json:"__typename,omitempty"`
	IsMuted                            bool                `json:"is_muted,omitempty"`
	ConversationID                     string              `json:"conversation_id,omitempty"`
	ParticipantsResults                []XChatUserResult   `json:"participants_results,omitempty"`
	GroupMetadata                      *XChatGroupMetadata `json:"group_metadata,omitempty"`
	GroupAdminsResults                 []XChatUserResult   `json:"group_admins_results,omitempty"`
	GroupMembersResults                []XChatUserResult   `json:"group_members_results,omitempty"`
	GroupRemovedUsers                  []XChatUserResult   `json:"group_removed_users,omitempty"`
	GroupJoinRequestedUsers            []XChatUserResult   `json:"group_join_requested_users,omitempty"`
	LatestGroupTitleChangeMessageEvent string              `json:"latest_group_title_change_message_event,omitempty"`
}

type XChatGroupMetadata struct {
	Typename                      string                   `json:"__typename,omitempty"`
	GroupName                     string                   `json:"group_name,omitempty"`
	CreatedAtMsec                 string                   `json:"created_at_msec,omitempty"`
	UpdatedAtMsec                 string                   `json:"updated_at_msec,omitempty"`
	GroupAvatarURL                string                   `json:"group_avatar_url,omitempty"`
	GroupInviteDetails            *XChatGroupInviteDetails `json:"group_invite_details,omitempty"`
	ScreenCaptureDetectionEnabled bool                     `json:"screen_capture_detection_enabled,omitempty"`
}

type XChatGroupInviteDetails struct {
	Typename       string `json:"__typename,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	InviteURL      string `json:"invite_url,omitempty"`
	Token          string `json:"token,omitempty"`
}

type XChatParticipantReadEvent struct {
	Typename                        string             `json:"__typename,omitempty"`
	ParticipantID                   XChatParticipantID `json:"participant_id,omitempty"`
	LatestMarkConversationReadEvent string             `json:"latest_mark_conversation_read_event,omitempty"`
}

type XChatParticipantID struct {
	Typename string `json:"__typename,omitempty"`
	RestID   string `json:"rest_id,omitempty"`
}

type XChatUserResult struct {
	Typename string     `json:"__typename,omitempty"`
	RestID   string     `json:"rest_id,omitempty"`
	Result   *XChatUser `json:"result,omitempty"`
}

type XChatUser struct {
	Typename                   string                     `json:"__typename,omitempty"`
	RestID                     string                     `json:"rest_id,omitempty"`
	Avatar                     *XChatUserAvatar           `json:"avatar,omitempty"`
	ChatPermissions            *XChatUserDMPermissions    `json:"chat_permissions,omitempty"`
	Core                       *XChatUserCore             `json:"core,omitempty"`
	Privacy                    *XChatUserPrivacy          `json:"privacy,omitempty"`
	AffiliatesHighlightedLabel *XChatHighlightedUserLabel `json:"affiliates_highlighted_label,omitempty"`
	Verification               *XChatUserVerification     `json:"verification,omitempty"`
	ProfileImageShape          string                     `json:"profile_image_shape,omitempty"`
}

type XChatUserAvatar struct {
	Typename string `json:"__typename,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type XChatUserDMPermissions struct {
	Typename          string `json:"__typename,omitempty"`
	CanDM             bool   `json:"can_dm,omitempty"`
	CanDMOnXChat      bool   `json:"can_dm_on_xchat,omitempty"`
	HasPublicKey      bool   `json:"has_public_key,omitempty"`
	CanBeAddedToGroup bool   `json:"can_be_added_to_group,omitempty"`
}

type XChatUserCore struct {
	Typename    string `json:"__typename,omitempty"`
	Name        string `json:"name,omitempty"`
	ScreenName  string `json:"screen_name,omitempty"`
	CreatedAtMS int64  `json:"created_at_ms,omitempty"`
}

type XChatUserPrivacy struct {
	Typename  string `json:"__typename,omitempty"`
	Protected bool   `json:"protected,omitempty"`
	Suspended bool   `json:"suspended,omitempty"`
}

type XChatHighlightedUserLabel struct {
	Typename string          `json:"__typename,omitempty"`
	Label    *XChatUserLabel `json:"label,omitempty"`
}

type XChatUserLabel struct {
	Typename             string            `json:"__typename,omitempty"`
	Badge                *XChatBadgeInfo   `json:"badge,omitempty"`
	Description          string            `json:"description,omitempty"`
	URL                  *XChatTimelineURL `json:"url,omitempty"`
	UserLabelDisplayType string            `json:"user_label_display_type,omitempty"`
	UserLabelType        string            `json:"user_label_type,omitempty"`
}

type XChatBadgeInfo struct {
	Typename string `json:"__typename,omitempty"`
	URL      string `json:"url,omitempty"`
}

type XChatTimelineURL struct {
	Typename string `json:"__typename,omitempty"`
	URL      string `json:"url,omitempty"`
	URLType  string `json:"url_type,omitempty"`
}

type XChatUserVerification struct {
	Typename                        string `json:"__typename,omitempty"`
	IsBlueVerified                  bool   `json:"is_blue_verified,omitempty"`
	IsVerifiedOrganizationAffiliate bool   `json:"is_verified_organization_affiliate,omitempty"`
	Verified                        bool   `json:"verified,omitempty"`
}

// GetPublicKeysResponse models the GraphQL response for fetching public keys and Juicebox tokens.
type GetPublicKeysResponse struct {
	Data struct {
		UserResultsByRestIDs []UserResultsWithPublicKeys `json:"user_results_by_rest_ids"`
	} `json:"data"`
}

type UserResultsWithPublicKeys struct {
	Typename string `json:"__typename,omitempty"`
	RestID   string `json:"rest_id,omitempty"`
	Result   struct {
		Typename      string              `json:"__typename,omitempty"`
		GetPublicKeys GetPublicKeysResult `json:"get_public_keys"`
	} `json:"result,omitempty"`
}

type GetPublicKeysResult struct {
	Typename               string                  `json:"__typename,omitempty"`
	PublicKeysWithTokenMap []PublicKeyWithTokenMap `json:"public_keys_with_token_map,omitempty"`
	IsManagedPinUser       bool                    `json:"is_managed_pin_user,omitempty"`
}

type PublicKeyWithTokenMap struct {
	Typename              string                 `json:"__typename,omitempty"`
	PublicKeyWithMetadata XChatPublicKeyWithMeta `json:"public_key_with_metadata"`
	TokenMap              KeyStoreTokenMap       `json:"token_map"`
}

type XChatPublicKeyWithMeta struct {
	Typename  string         `json:"__typename,omitempty"`
	PublicKey XChatPublicKey `json:"public_key"`
	Version   string         `json:"version,omitempty"` // This is the signing_key_version
}

type XChatPublicKey struct {
	Typename                   string `json:"__typename,omitempty"`
	PublicKey                  string `json:"public_key,omitempty"`
	SigningPublicKey           string `json:"signing_public_key,omitempty"`
	IdentityPublicKeySignature string `json:"identity_public_key_signature,omitempty"`
}

type KeyStoreTokenMap struct {
	Typename             string               `json:"__typename,omitempty"`
	RealmState           string               `json:"realm_state,omitempty"`
	MaxGuessCount        int                  `json:"max_guess_count,omitempty"`
	RecoverThreshold     int                  `json:"recover_threshold,omitempty"`
	RegisterThreshold    int                  `json:"register_threshold,omitempty"`
	TokenMap             []KeyStoreTokenEntry `json:"token_map,omitempty"`
	KeyStoreTokenMapJSON string               `json:"key_store_token_map_json,omitempty"` // Ready-to-use Juicebox config JSON
}

type KeyStoreTokenEntry struct {
	Typename string        `json:"__typename,omitempty"`
	Key      string        `json:"key,omitempty"` // Realm ID (hex)
	Value    KeyStoreToken `json:"value,omitempty"`
}

type KeyStoreToken struct {
	Typename  string `json:"__typename,omitempty"`
	Token     string `json:"token,omitempty"`      // Auth token for this realm
	Address   string `json:"address,omitempty"`    // Realm URL
	PublicKey string `json:"public_key,omitempty"` // Optional realm public key
}
