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

	ACCOUNT_SETTINGS_URL               = API_BASE_URL + "/1.1/account/settings.json"
	CONVERSATION_MARK_READ_URL         = BASE_URL + "/i/api/1.1/dm/conversation/%s/mark_read.json"
	GRAPHQL_MESSAGE_DELETION_MUTATION  = BASE_URL + "/i/api/graphql/BJ6DtxA2llfjnRoRjaiIiw/DMMessageDeleteMutation"
	SEARCH_TYPEAHEAD_URL               = BASE_URL + "/i/api/1.1/search/typeahead.json"
	DM_PERMISSIONS_URL                 = BASE_URL + "/i/api/1.1/dm/permissions.json"
	DELETE_CONVERSATION_URL            = BASE_URL + "/i/api/1.1/dm/conversation/%s/delete.json"
	UPDATE_CONVERSATION_AVATAR_URL     = BASE_URL + "/i/api/1.1/dm/conversation/%s/update_avatar.json"
	UPDATE_CONVERSATION_NAME_URL       = BASE_URL + "/i/api/1.1/dm/conversation/%s/update_name.json"
	PIN_CONVERSATION_URL               = BASE_URL + "/i/api/graphql/o0aymgGiJY-53Y52YSUGVA/DMPinnedInboxAppend_Mutation"
	UNPIN_CONVERSATION_URL             = BASE_URL + "/i/api/graphql/_TQxP2Rb0expwVP9ktGrTQ/DMPinnedInboxDelete_Mutation"
	GET_PINNED_CONVERSATIONS_URL       = BASE_URL + "/i/api/graphql/_gBQBgClVuMQb8efxWkbbQ/DMPinnedInboxQuery"
	ADD_PARTICIPANTS_URL               = BASE_URL + "/i/api/graphql/oBwyQ0_xVbAQ8FAyG0pCRA/AddParticipantsMutation"
	ADD_REACTION_URL                   = BASE_URL + "/i/api/graphql/VyDyV9pC2oZEj6g52hgnhA/useDMReactionMutationAddMutation"
	REMOVE_REACTION_URL                = BASE_URL + "/i/api/graphql/bV_Nim3RYHsaJwMkTXJ6ew/useDMReactionMutationRemoveMutation"
	SEND_TYPING_NOTIFICATION           = BASE_URL + "/i/api/graphql/HL96-xZ3Y81IEzAdczDokg/useTypingNotifierMutation"
	GENERATE_XCHAT_TOKEN_MUTATION_URL  = API_BASE_URL + "/graphql/Qh3fZRjPPtPoHYR_2sCZsA/GenerateXChatTokenMutation"
	SEND_MESSAGE_MUTATION_URL          = API_BASE_URL + "/graphql/LkAIEchf8AGj-WgeLoTVcw/SendMessageMutation"
	DELETE_MESSAGE_MUTATION_URL        = API_BASE_URL + "/graphql/4gsDQKEmYkOtvsSIpHXdQA/DeleteMessageMutation"
	MUTE_CONVERSATION_URL              = API_BASE_URL + "/graphql/Dy7geJg7CL5dqhsl6QBteg/MuteConversation"
	UNMUTE_CONVERSATION_URL            = API_BASE_URL + "/graphql/LnNSeGu4vnbwqXAvh7OlGQ/UnmuteConversation"
	GET_INITIAL_XCHAT_PAGE_QUERY_URL   = API_BASE_URL + "/graphql/6gAgW7rM7oOZMq_um-zGyg/GetInitialXChatPageQuery"
	GET_INBOX_PAGE_REQUEST_QUERY_URL   = API_BASE_URL + "/graphql/wmieJEOHm6twV06EXwRdiA/GetInboxPageRequestQuery"
	GET_INBOX_PAGE_CONV_DATA_QUERY_URL = API_BASE_URL + "/graphql/uQEDp5FgdqNiG2jT5q07Jw/GetInboxPageConversationDataRequestQuery"
	GET_USERS_BY_IDS_FOR_XCHAT_URL     = API_BASE_URL + "/graphql/U5qYdpM3QOzlHh5K7ubs2A/GetUsersByIdsForXChat"
	GET_CONVERSATION_PAGE_QUERY_URL    = API_BASE_URL + "/graphql/IVlXls9JTnbgQ1gxsGAfJA/GetConversationPageQuery"
	GET_PUBLIC_KEYS_QUERY_URL          = API_BASE_URL + "/graphql/RQAjOoIX9dIsHoVjuVV0Iw/GetPublicKeys"

	XCHAT_WEBSOCKET_URL = "wss://chat-ws.x.com/ws"

	INITIALIZE_XCHAT_MEDIA_UPLOAD_URL = API_BASE_URL + "/graphql/g2n9PB_uaRYv_SFvQokEFw/InitializeXChatMediaUpload"
	FINALIZE_XCHAT_MEDIA_UPLOAD_URL   = API_BASE_URL + "/graphql/UK24H5vBa5MJspBmjZyFVQ/FinalizeXChatMediaUpload"
	TON_UPLOAD_BASE_URL               = "https://ton.x.com"

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
