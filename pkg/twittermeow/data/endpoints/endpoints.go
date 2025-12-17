package endpoints

// TODO fix variable name casing

const (
	TWITTER_BASE_HOST = "twitter.com"
	TWITTER_BASE_URL  = "https://" + TWITTER_BASE_HOST

	BASE_HOST                      = "x.com"
	BASE_URL                       = "https://" + BASE_HOST
	BASE_LOGIN_URL                 = BASE_URL + "/login"
	BASE_MESSAGES_URL              = BASE_URL + "/messages"
	BASE_LOGOUT_URL                = BASE_URL + "/logout"
	BASE_NOTIFICATION_SETTINGS_URL = BASE_URL + "/settings/push_notifications"

	API_BASE_HOST = "api.x.com"
	API_BASE_URL  = "https://" + API_BASE_HOST

	ACCOUNT_SETTINGS_URL              = API_BASE_URL + "/1.1/account/settings.json"
	INBOX_INITIAL_STATE_URL           = BASE_URL + "/i/api/1.1/dm/inbox_initial_state.json"
	DM_USER_UPDATES_URL               = BASE_URL + "/i/api/1.1/dm/user_updates.json"
	CONVERSATION_MARK_READ_URL        = BASE_URL + "/i/api/1.1/dm/conversation/%s/mark_read.json"
	CONVERSATION_FETCH_MESSAGES       = BASE_URL + "/i/api/1.1/dm/conversation/%s.json"
	UPDATE_LAST_SEEN_EVENT_ID_URL     = BASE_URL + "/i/api/1.1/dm/update_last_seen_event_id.json"
	TRUSTED_INBOX_TIMELINE_URL        = BASE_URL + "/i/api/1.1/dm/inbox_timeline/trusted.json"
	SEND_DM_URL                       = BASE_URL + "/i/api/1.1/dm/new2.json"
	EDIT_DM_URL                       = BASE_URL + "/i/api/1.1/dm/edit.json"
	GRAPHQL_MESSAGE_DELETION_MUTATION = BASE_URL + "/i/api/graphql/BJ6DtxA2llfjnRoRjaiIiw/DMMessageDeleteMutation"
	SEARCH_TYPEAHEAD_URL              = BASE_URL + "/i/api/1.1/search/typeahead.json"
	DM_PERMISSIONS_URL                = BASE_URL + "/i/api/1.1/dm/permissions.json"
	ACCEPT_CONVERSATION_URL           = BASE_URL + "/i/api/1.1/dm/conversation/%s/accept.json"
	DELETE_CONVERSATION_URL           = BASE_URL + "/i/api/1.1/dm/conversation/%s/delete.json"
	UPDATE_CONVERSATION_AVATAR_URL    = BASE_URL + "/i/api/1.1/dm/conversation/%s/update_avatar.json"
	UPDATE_CONVERSATION_NAME_URL      = BASE_URL + "/i/api/1.1/dm/conversation/%s/update_name.json"
	PIN_CONVERSATION_URL              = BASE_URL + "/i/api/graphql/o0aymgGiJY-53Y52YSUGVA/DMPinnedInboxAppend_Mutation"
	UNPIN_CONVERSATION_URL            = BASE_URL + "/i/api/graphql/_TQxP2Rb0expwVP9ktGrTQ/DMPinnedInboxDelete_Mutation"
	GET_PINNED_CONVERSATIONS_URL      = BASE_URL + "/i/api/graphql/_gBQBgClVuMQb8efxWkbbQ/DMPinnedInboxQuery"
	ADD_PARTICIPANTS_URL              = BASE_URL + "/i/api/graphql/oBwyQ0_xVbAQ8FAyG0pCRA/AddParticipantsMutation"
	ADD_REACTION_URL                  = BASE_URL + "/i/api/graphql/VyDyV9pC2oZEj6g52hgnhA/useDMReactionMutationAddMutation"
	REMOVE_REACTION_URL               = BASE_URL + "/i/api/graphql/bV_Nim3RYHsaJwMkTXJ6ew/useDMReactionMutationRemoveMutation"
	SEND_TYPING_NOTIFICATION          = BASE_URL + "/i/api/graphql/HL96-xZ3Y81IEzAdczDokg/useTypingNotifierMutation"

	JOT_CLIENT_EVENT_URL = API_BASE_URL + "/1.1/jot/client_event.json"
	JOT_CES_P2_URL       = API_BASE_URL + "/1.1/jot/ces/p2"

	PIPELINE_EVENTS_URL = API_BASE_URL + "/live_pipeline/events"
	PIPELINE_UPDATE_URL = API_BASE_URL + "/1.1/live_pipeline/update_subscriptions"

	UPLOAD_BASE_HOST = "upload.x.com"
	UPLOAD_BASE_URL  = "https://" + UPLOAD_BASE_HOST
	UPLOAD_MEDIA_URL = UPLOAD_BASE_URL + "/i/media/upload.json"

	NOTIFICATION_SETTINGS_URL = BASE_URL + "/i/api/1.1/notifications/settings"
	NOTIFICATION_LOGIN_URL    = NOTIFICATION_SETTINGS_URL + "/login.json"
	NOTIFICATION_LOGOUT_URL   = NOTIFICATION_SETTINGS_URL + "/logout.json"
	NOTIFICATION_CHECKIN_URL  = NOTIFICATION_SETTINGS_URL + "/checkin.json"
	NOTIFICATION_SAVE_URL     = NOTIFICATION_SETTINGS_URL + "/save.json"
)
