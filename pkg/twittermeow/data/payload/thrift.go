package payload

// Code generated from reverse-xchat/thrift/dmv2.thriftjava; manual edits may be overwritten.

type AdditionalAction int32

const (
	AdditionalActionFetchConvIfMissingCkey AdditionalAction = 0
)

type DeleteMessageAction int32

const (
	DeleteMessageActionDeleteForSelf DeleteMessageAction = 1
	DeleteMessageActionDeleteForAll  DeleteMessageAction = 2
)

type FailureType int32

const (
	FailureTypeEmptyDetail                        FailureType = 1
	FailureTypeInternalError                      FailureType = 2
	FailureTypeContentsTooLarge                   FailureType = 3
	FailureTypeTooManyMessages                    FailureType = 4
	FailureTypeInvalidSenderSignature             FailureType = 5
	FailureTypeNonLatestCkeyVersion               FailureType = 6
	FailureTypeRecipientHasNotTrustedConversation FailureType = 7
	FailureTypeRecipientKeyHasChanged             FailureType = 8
)

type MediaType int32

const (
	MediaTypeImage MediaType = 1
	MediaTypeGif   MediaType = 2
	MediaTypeVideo MediaType = 3
	MediaTypeAudio MediaType = 4
	MediaTypeFile  MediaType = 5
	MediaTypeSvg   MediaType = 6
)

type ScreenCaptureType int32

const (
	ScreenCaptureTypeScreenshot ScreenCaptureType = 1
	ScreenCaptureTypeRecording  ScreenCaptureType = 2
)

type SentFromSurface int32

const (
	SentFromSurfaceConversationScreenComposer SentFromSurface = 1
	SentFromSurfaceNotificationReply          SentFromSurface = 2
	SentFromSurfaceShareSheet                 SentFromSurface = 3
	SentFromSurfacePaymentsSupportComposer    SentFromSurface = 4
	SentFromSurfaceMessageForwardSheet        SentFromSurface = 5
)

type AVCallEnded struct {
	SentAtMillis    *int64  `thrift:"sent_at_millis,1" json:"sent_at_millis,omitempty"`
	DurationSeconds *int64  `thrift:"duration_seconds,2" json:"duration_seconds,omitempty"`
	IsAudioOnly     *bool   `thrift:"is_audio_only,3" json:"is_audio_only,omitempty"`
	BroadcastId     *string `thrift:"broadcast_id,5" json:"broadcast_id,omitempty"`
}

type AVCallMissed struct {
	SentAtMillis *int64 `thrift:"sent_at_millis,1" json:"sent_at_millis,omitempty"`
	IsAudioOnly  *bool  `thrift:"is_audio_only,2" json:"is_audio_only,omitempty"`
}

type AVCallStarted struct {
	IsAudioOnly *bool   `thrift:"is_audio_only,1" json:"is_audio_only,omitempty"`
	BroadcastId *string `thrift:"broadcast_id,3" json:"broadcast_id,omitempty"`
}

type AcceptMessageRequest struct {
}

type AddressRichTextContent struct {
}

type BatchedMessageEvents struct {
	MessageEvents []*MessageEvent `thrift:"message_events,1" json:"message_events,omitempty"`
}

type CallToAction struct {
	Label *string `thrift:"label,1" json:"label,omitempty"`
	Url   *string `thrift:"url,2" json:"url,omitempty"`
}

type CashtagRichTextContent struct {
}

type ConversationDeleteEvent struct {
	ConversationId *string `thrift:"conversation_id,1" json:"conversation_id,omitempty"`
}

type ConversationKeyChangeEvent struct {
	ConversationKeyVersion      *string                       `thrift:"conversation_key_version,1" json:"conversation_key_version,omitempty"`
	ConversationParticipantKeys []*ConversationParticipantKey `thrift:"conversation_participant_keys,2" json:"conversation_participant_keys,omitempty"`
	RatchetTree                 *KeyRotation                  `thrift:"ratchet_tree,3" json:"ratchet_tree,omitempty"`
}

type ConversationMetadataChange struct {
	MessageDurationChange         *MessageDurationChange         `thrift:"message_duration_change,1" json:"message_duration_change,omitempty"`
	MessageDurationRemove         *MessageDurationRemove         `thrift:"message_duration_remove,2" json:"message_duration_remove,omitempty"`
	MuteConversation              *MuteConversation              `thrift:"mute_conversation,3" json:"mute_conversation,omitempty"`
	UnmuteConversation            *UnmuteConversation            `thrift:"unmute_conversation,4" json:"unmute_conversation,omitempty"`
	EnableScreenCaptureDetection  *EnableScreenCaptureDetection  `thrift:"enable_screen_capture_detection,5" json:"enable_screen_capture_detection,omitempty"`
	DisableScreenCaptureDetection *DisableScreenCaptureDetection `thrift:"disable_screen_capture_detection,6" json:"disable_screen_capture_detection,omitempty"`
	EnableScreenCaptureBlocking   *EnableScreenCaptureBlocking   `thrift:"enable_screen_capture_blocking,7" json:"enable_screen_capture_blocking,omitempty"`
	DisableScreenCaptureBlocking  *DisableScreenCaptureBlocking  `thrift:"disable_screen_capture_blocking,8" json:"disable_screen_capture_blocking,omitempty"`
}

type ConversationMetadataChangeEvent struct {
	ConversationMetadataChange *ConversationMetadataChange `thrift:"conversation_metadata_change,1" json:"conversation_metadata_change,omitempty"`
}

type ConversationParticipantKey struct {
	UserId                   *string `thrift:"user_id,1" json:"user_id,omitempty"`
	EncryptedConversationKey *string `thrift:"encrypted_conversation_key,2" json:"encrypted_conversation_key,omitempty"`
	PublicKeyVersion         *string `thrift:"public_key_version,3" json:"public_key_version,omitempty"`
}

type DisableScreenCaptureBlocking struct {
	Placeholder *string `thrift:"placeholder,1" json:"placeholder,omitempty"`
}

type DisableScreenCaptureDetection struct {
	Placeholder *string `thrift:"placeholder,1" json:"placeholder,omitempty"`
}

type DisplayTemporaryPasscodeInstruction struct {
	Token                  *string `thrift:"token,1" json:"token,omitempty"`
	LatestPublicKeyVersion *string `thrift:"latest_public_key_version,2" json:"latest_public_key_version,omitempty"`
}

type DraftMessage struct {
	ConversationId *string `thrift:"conversation_id,1" json:"conversation_id,omitempty"`
	DraftText      *string `thrift:"draft_text,2" json:"draft_text,omitempty"`
}

type EmailRichTextContent struct {
}

type EmptyNode struct {
	Description *string `thrift:"description,1" json:"description,omitempty"`
}

type EnableScreenCaptureBlocking struct {
	Placeholder *string `thrift:"placeholder,1" json:"placeholder,omitempty"`
}

type EnableScreenCaptureDetection struct {
	Placeholder *string `thrift:"placeholder,1" json:"placeholder,omitempty"`
}

type EventQueuePriority struct {
}

type ForwardedMessage struct {
	MessageText *string           `thrift:"message_text,1" json:"message_text,omitempty"`
	Entities    []*RichTextEntity `thrift:"entities,2" json:"entities,omitempty"`
}

type GrokSearchResponseEvent struct {
	SearchResponseId *string `thrift:"search_response_id,1" json:"search_response_id,omitempty"`
}

type GroupAdminAddChange struct {
	AdminIds []string `thrift:"admin_ids,1" json:"admin_ids,omitempty"`
}

type GroupAdminRemoveChange struct {
	AdminIds []string `thrift:"admin_ids,1" json:"admin_ids,omitempty"`
}

type GroupAvatarUrlChange struct {
	CustomAvatarUrl        *string `thrift:"custom_avatar_url,1" json:"custom_avatar_url,omitempty"`
	ConversationKeyVersion *string `thrift:"conversation_key_version,2" json:"conversation_key_version,omitempty"`
}

type GroupChange struct {
	GroupCreate        *GroupCreate             `thrift:"group_create,1" json:"group_create,omitempty"`
	GroupTitleChange   *GroupTitleChange        `thrift:"group_title_change,2" json:"group_title_change,omitempty"`
	GroupAvatarChange  *GroupAvatarUrlChange    `thrift:"group_avatar_change,3" json:"group_avatar_change,omitempty"`
	GroupAdminAdd      *GroupAdminAddChange     `thrift:"group_admin_add,4" json:"group_admin_add,omitempty"`
	GroupMemberAdd     *GroupMemberAddChange    `thrift:"group_member_add,5" json:"group_member_add,omitempty"`
	GroupAdminRemove   *GroupAdminRemoveChange  `thrift:"group_admin_remove,6" json:"group_admin_remove,omitempty"`
	GroupMemberRemove  *GroupMemberRemoveChange `thrift:"group_member_remove,7" json:"group_member_remove,omitempty"`
	GroupInviteEnable  *GroupInviteEnable       `thrift:"group_invite_enable,8" json:"group_invite_enable,omitempty"`
	GroupInviteDisable *GroupInviteDisable      `thrift:"group_invite_disable,9" json:"group_invite_disable,omitempty"`
	GroupJoinRequest   *GroupJoinRequest        `thrift:"group_join_request,10" json:"group_join_request,omitempty"`
	GroupJoinReject    *GroupJoinReject         `thrift:"group_join_reject,11" json:"group_join_reject,omitempty"`
}

type GroupChangeEvent struct {
	GroupChange *GroupChange `thrift:"group_change,1" json:"group_change,omitempty"`
}

type GroupCreate struct {
	MemberIds              []string `thrift:"member_ids,1" json:"member_ids,omitempty"`
	AdminIds               []string `thrift:"admin_ids,2" json:"admin_ids,omitempty"`
	Title                  *string  `thrift:"title,3" json:"title,omitempty"`
	AvatarUrl              *string  `thrift:"avatar_url,4" json:"avatar_url,omitempty"`
	ConversationKeyVersion *string  `thrift:"conversation_key_version,5" json:"conversation_key_version,omitempty"`
}

type GroupInviteDisable struct {
	DisabledByMemberId *string `thrift:"disabled_by_member_id,1" json:"disabled_by_member_id,omitempty"`
}

type GroupInviteEnable struct {
	ExpiresAtMsec *int64  `thrift:"expires_at_msec,1" json:"expires_at_msec,omitempty"`
	InviteUrl     *string `thrift:"invite_url,2" json:"invite_url,omitempty"`
	AffiliateId   *string `thrift:"affiliate_id,3" json:"affiliate_id,omitempty"`
}

type GroupJoinReject struct {
	RejectedUserIds []string `thrift:"rejected_user_ids,1" json:"rejected_user_ids,omitempty"`
}

type GroupJoinRequest struct {
	RequestingUserId *string `thrift:"requesting_user_id,1" json:"requesting_user_id,omitempty"`
}

type GroupMemberAddChange struct {
	MemberIds               []string `thrift:"member_ids,1" json:"member_ids,omitempty"`
	CurrentMemberIds        []string `thrift:"current_member_ids,2" json:"current_member_ids,omitempty"`
	CurrentAdminIds         []string `thrift:"current_admin_ids,3" json:"current_admin_ids,omitempty"`
	CurrentTitle            *string  `thrift:"current_title,4" json:"current_title,omitempty"`
	CurrentAvatarUrl        *string  `thrift:"current_avatar_url,5" json:"current_avatar_url,omitempty"`
	ConversationKeyVersion  *string  `thrift:"conversation_key_version,6" json:"conversation_key_version,omitempty"`
	CurrentTtlMsec          *int64   `thrift:"current_ttl_msec,7" json:"current_ttl_msec,omitempty"`
	CurrentPendingMemberIds []string `thrift:"current_pending_member_ids,8" json:"current_pending_member_ids,omitempty"`
}

type GroupMemberRemoveChange struct {
	MemberIds []string `thrift:"member_ids,1" json:"member_ids,omitempty"`
}

type GroupTitleChange struct {
	CustomTitle            *string `thrift:"custom_title,1" json:"custom_title,omitempty"`
	ConversationKeyVersion *string `thrift:"conversation_key_version,2" json:"conversation_key_version,omitempty"`
}

type HashtagRichTextContent struct {
}

type KeepAliveInstruction struct {
}

type KeyRotation struct {
	PreviousVersion     *string           `thrift:"previous_version,1" json:"previous_version,omitempty"`
	RatchetTree         *RatchetTree      `thrift:"ratchet_tree,2" json:"ratchet_tree,omitempty"`
	Nodes               []*UpdatePathNode `thrift:"nodes,3" json:"nodes,omitempty"`
	EncryptedPrivateKey *string           `thrift:"encrypted_private_key,4" json:"encrypted_private_key,omitempty"`
}

type LeafNode struct {
	SubtreeEncryptionPublicKey  *string `thrift:"subtree_encryption_public_key,1" json:"subtree_encryption_public_key,omitempty"`
	SignaturePublicKey          *string `thrift:"signature_public_key,2" json:"signature_public_key,omitempty"`
	KeypairId                   *string `thrift:"keypair_id,3" json:"keypair_id,omitempty"`
	MaxSupportedProtocolVersion *int32  `thrift:"max_supported_protocol_version,4" json:"max_supported_protocol_version,omitempty"`
	ParentHash                  *string `thrift:"parent_hash,5" json:"parent_hash,omitempty"`
	Signature                   *string `thrift:"signature,6" json:"signature,omitempty"`
}

type MarkConversationRead struct {
	SeenUntilSequenceId *string `thrift:"seen_until_sequence_id,1" json:"seen_until_sequence_id,omitempty"`
	SeenAtMillis        *int64  `thrift:"seen_at_millis,2" json:"seen_at_millis,omitempty"`
}

type MarkConversationReadEvent struct {
	SeenUntilSequenceId *string `thrift:"seen_until_sequence_id,1" json:"seen_until_sequence_id,omitempty"`
	SeenAtMillis        *int64  `thrift:"seen_at_millis,2" json:"seen_at_millis,omitempty"`
}

type MarkConversationUnread struct {
	SeenUntilSequenceId *string `thrift:"seen_until_sequence_id,1" json:"seen_until_sequence_id,omitempty"`
}

type MarkConversationUnreadEvent struct {
	SeenUntilSequenceId *string `thrift:"seen_until_sequence_id,1" json:"seen_until_sequence_id,omitempty"`
}

type MaybeKeypair struct {
	Empty   *string        `thrift:"empty,1" json:"empty,omitempty"`
	Keypair *StoredKeypair `thrift:"keypair,2" json:"keypair,omitempty"`
}

type MediaAttachment struct {
	MediaHashKey          *string          `thrift:"media_hash_key,1" json:"media_hash_key,omitempty"`
	Dimensions            *MediaDimensions `thrift:"dimensions,2" json:"dimensions,omitempty"`
	Type                  *int32           `thrift:"type,3" json:"type,omitempty"`
	DurationMillis        *int64           `thrift:"duration_millis,4" json:"duration_millis,omitempty"`
	FilesizeBytes         *int64           `thrift:"filesize_bytes,5" json:"filesize_bytes,omitempty"`
	Filename              *string          `thrift:"filename,6" json:"filename,omitempty"`
	AttachmentId          *string          `thrift:"attachment_id,7" json:"attachment_id,omitempty"`
	LegacyMediaUrlHttps   *string          `thrift:"legacy_media_url_https,8" json:"legacy_media_url_https,omitempty"`
	LegacyMediaPreviewUrl *string          `thrift:"legacy_media_preview_url,9" json:"legacy_media_preview_url,omitempty"`
}

type MediaDimensions struct {
	Width  *int64 `thrift:"width,1" json:"width,omitempty"`
	Height *int64 `thrift:"height,2" json:"height,omitempty"`
}

type MemberAccountDeleteEvent struct {
	MemberId *string `thrift:"member_id,1" json:"member_id,omitempty"`
}

type MentionRichTextContent struct {
}

type Message struct {
	MessageEvent         *MessageEvent         `thrift:"messageEvent,1" json:"messageEvent,omitempty"`
	MessageInstruction   *MessageInstruction   `thrift:"messageInstruction,2" json:"messageInstruction,omitempty"`
	BatchedMessageEvents *BatchedMessageEvents `thrift:"batchedMessageEvents,3" json:"batchedMessageEvents,omitempty"`
}

type MessageAttachment struct {
	Media       *MediaAttachment       `thrift:"media,1" json:"media,omitempty"`
	Post        *PostAttachment        `thrift:"post,2" json:"post,omitempty"`
	Url         *UrlAttachment         `thrift:"url,3" json:"url,omitempty"`
	UnifiedCard *UnifiedCardAttachment `thrift:"unified_card,4" json:"unified_card,omitempty"`
	Money       *MoneyAttachment       `thrift:"money,5" json:"money,omitempty"`
}

type MessageContents struct {
	MessageText       *string              `thrift:"message_text,1" json:"message_text,omitempty"`
	Entities          []*RichTextEntity    `thrift:"entities,2" json:"entities,omitempty"`
	Attachments       []*MessageAttachment `thrift:"attachments,3" json:"attachments,omitempty"`
	ReplyingToPreview *ReplyingToPreview   `thrift:"replying_to_preview,4" json:"replying_to_preview,omitempty"`
	ForwardedMessage  *ForwardedMessage    `thrift:"forwarded_message,6" json:"forwarded_message,omitempty"`
	SentFrom          *int32               `thrift:"sent_from,7" json:"sent_from,omitempty"`
	QuickReply        *QuickReply          `thrift:"quick_reply,8" json:"quick_reply,omitempty"`
	Ctas              []*CallToAction      `thrift:"ctas,9" json:"ctas,omitempty"`
}

type MessageCreateEvent struct {
	Contents               []byte  `thrift:"contents,100" json:"contents,omitempty"`
	ConversationKeyVersion *string `thrift:"conversation_key_version,101" json:"conversation_key_version,omitempty"`
	ShouldNotify           *bool   `thrift:"should_notify,102" json:"should_notify,omitempty"`
	TtlMsec                *int64  `thrift:"ttl_msec,103" json:"ttl_msec,omitempty"`
	DeliveredAtMsec        *int64  `thrift:"delivered_at_msec,104" json:"delivered_at_msec,omitempty"`
	IsPendingPublicKey     *bool   `thrift:"is_pending_public_key,105" json:"is_pending_public_key,omitempty"`
	Priority               *int32  `thrift:"priority,106" json:"priority,omitempty"`
	AdditionalActionList   []int32 `thrift:"additional_action_list,107" json:"additional_action_list,omitempty"`
}

type MessageDeleteEvent struct {
	SequenceIds         []string `thrift:"sequence_ids,1" json:"sequence_ids,omitempty"`
	DeleteMessageAction *int32   `thrift:"delete_message_action,2" json:"delete_message_action,omitempty"`
}

type MessageDurationChange struct {
	TtlMsec *int64 `thrift:"ttl_msec,1" json:"ttl_msec,omitempty"`
}

type MessageDurationRemove struct {
	CurrentTtlMsec *int64 `thrift:"current_ttl_msec,1" json:"current_ttl_msec,omitempty"`
}

type MessageEdit struct {
	MessageSequenceId *string           `thrift:"message_sequence_id,1" json:"message_sequence_id,omitempty"`
	UpdatedText       *string           `thrift:"updated_text,2" json:"updated_text,omitempty"`
	Entities          []*RichTextEntity `thrift:"entities,3" json:"entities,omitempty"`
}

type MessageEntryContents struct {
	Message                *MessageContents        `thrift:"message,1" json:"message,omitempty"`
	ReactionAdd            *MessageReactionAdd     `thrift:"reaction_add,2" json:"reaction_add,omitempty"`
	ReactionRemove         *MessageReactionRemove  `thrift:"reaction_remove,3" json:"reaction_remove,omitempty"`
	MessageEdit            *MessageEdit            `thrift:"message_edit,4" json:"message_edit,omitempty"`
	MarkConversationRead   *MarkConversationRead   `thrift:"mark_conversation_read,5" json:"mark_conversation_read,omitempty"`
	MarkConversationUnread *MarkConversationUnread `thrift:"mark_conversation_unread,6" json:"mark_conversation_unread,omitempty"`
	PinConversation        *PinConversation        `thrift:"pin_conversation,7" json:"pin_conversation,omitempty"`
	UnpinConversation      *UnpinConversation      `thrift:"unpin_conversation,8" json:"unpin_conversation,omitempty"`
	ScreenCaptureDetected  *ScreenCaptureDetected  `thrift:"screen_capture_detected,9" json:"screen_capture_detected,omitempty"`
	AvCallEnded            *AVCallEnded            `thrift:"av_call_ended,10" json:"av_call_ended,omitempty"`
	AvCallMissed           *AVCallMissed           `thrift:"av_call_missed,11" json:"av_call_missed,omitempty"`
	DraftMessage           *DraftMessage           `thrift:"draft_message,12" json:"draft_message,omitempty"`
	AcceptMessageRequest   *AcceptMessageRequest   `thrift:"accept_message_request,13" json:"accept_message_request,omitempty"`
	NicknameMessage        *NicknameMessage        `thrift:"nickname_message,14" json:"nickname_message,omitempty"`
	SetVerifiedStatus      *SetVerifiedStatus      `thrift:"set_verified_status,15" json:"set_verified_status,omitempty"`
	AvCallStarted          *AVCallStarted          `thrift:"av_call_started,16" json:"av_call_started,omitempty"`
}

type MessageEntryHolder struct {
	Contents *MessageEntryContents `thrift:"contents,1" json:"contents,omitempty"`
}

type MessageEvent struct {
	SequenceId            *string                `thrift:"sequence_id,1" json:"sequence_id,omitempty"`
	MessageId             *string                `thrift:"message_id,2" json:"message_id,omitempty"`
	SenderId              *string                `thrift:"sender_id,3" json:"sender_id,omitempty"`
	ConversationId        *string                `thrift:"conversation_id,4" json:"conversation_id,omitempty"`
	ConversationToken     *string                `thrift:"conversation_token,5" json:"conversation_token,omitempty"`
	CreatedAtMsec         *string                `thrift:"created_at_msec,6" json:"created_at_msec,omitempty"`
	Detail                *MessageEventDetail    `thrift:"detail,7" json:"detail,omitempty"`
	RelaySource           *int32                 `thrift:"relay_source,8" json:"relay_source,omitempty"`
	MessageEventSignature *MessageEventSignature `thrift:"message_event_signature,9" json:"message_event_signature,omitempty"`
	PreviousSequenceId    *string                `thrift:"previous_sequence_id,10" json:"previous_sequence_id,omitempty"`
	IsTrusted             *bool                  `thrift:"is_trusted,11" json:"is_trusted,omitempty"`
}

type MessageEventDetail struct {
	MessageCreateEvent              *MessageCreateEvent              `thrift:"messageCreateEvent,1" json:"messageCreateEvent,omitempty"`
	ConversationKeyChangeEvent      *ConversationKeyChangeEvent      `thrift:"conversationKeyChangeEvent,3" json:"conversationKeyChangeEvent,omitempty"`
	GroupChangeEvent                *GroupChangeEvent                `thrift:"groupChangeEvent,4" json:"groupChangeEvent,omitempty"`
	MessageFailureEvent             *MessageFailureEvent             `thrift:"messageFailureEvent,5" json:"messageFailureEvent,omitempty"`
	MessageTypingEvent              *MessageTypingEvent              `thrift:"messageTypingEvent,6" json:"messageTypingEvent,omitempty"`
	MessageDeleteEvent              *MessageDeleteEvent              `thrift:"messageDeleteEvent,7" json:"messageDeleteEvent,omitempty"`
	ConversationDeleteEvent         *ConversationDeleteEvent         `thrift:"conversationDeleteEvent,8" json:"conversationDeleteEvent,omitempty"`
	ConversationMetadataChangeEvent *ConversationMetadataChangeEvent `thrift:"conversationMetadataChangeEvent,9" json:"conversationMetadataChangeEvent,omitempty"`
	GrokSearchResponseEvent         *GrokSearchResponseEvent         `thrift:"grokSearchResponseEvent,10" json:"grokSearchResponseEvent,omitempty"`
	RequestForEncryptedResendEvent  *RequestForEncryptedResendEvent  `thrift:"requestForEncryptedResendEvent,11" json:"requestForEncryptedResendEvent,omitempty"`
	MarkConversationReadEvent       *MarkConversationReadEvent       `thrift:"markConversationReadEvent,12" json:"markConversationReadEvent,omitempty"`
	MarkConversationUnreadEvent     *MarkConversationUnreadEvent     `thrift:"markConversationUnreadEvent,13" json:"markConversationUnreadEvent,omitempty"`
	MemberAccountDeleteEvent        *MemberAccountDeleteEvent        `thrift:"memberAccountDeleteEvent,14" json:"memberAccountDeleteEvent,omitempty"`
}

type MessageEventRelaySource struct {
}

type MessageEventSignature struct {
	Signature        *string `thrift:"signature,1" json:"signature,omitempty"`
	PublicKeyVersion *string `thrift:"public_key_version,2" json:"public_key_version,omitempty"`
	SignatureVersion *string `thrift:"signature_version,3" json:"signature_version,omitempty"`
	SigningPublicKey *string `thrift:"signing_public_key,4" json:"signing_public_key,omitempty"`
}

type MessageFailureEvent struct {
	FailureType *int32 `thrift:"failure_type,1" json:"failure_type,omitempty"`
}

type MessageInstruction struct {
	PullMessagesInstruction             *PullMessagesInstruction             `thrift:"pullMessagesInstruction,1" json:"pullMessagesInstruction,omitempty"`
	KeepAliveInstruction                *KeepAliveInstruction                `thrift:"keepAliveInstruction,2" json:"keepAliveInstruction,omitempty"`
	PullMessagesFinishedInstruction     *PullMessagesFinishedInstruction     `thrift:"pullMessagesFinishedInstruction,3" json:"pullMessagesFinishedInstruction,omitempty"`
	PinReminderInstruction              *PinReminderInstruction              `thrift:"pinReminderInstruction,4" json:"pinReminderInstruction,omitempty"`
	SwitchToHybridPullInstruction       *SwitchToHybridPullInstruction       `thrift:"switchToHybridPullInstruction,5" json:"switchToHybridPullInstruction,omitempty"`
	DisplayTemporaryPasscodeInstruction *DisplayTemporaryPasscodeInstruction `thrift:"displayTemporaryPasscodeInstruction,6" json:"displayTemporaryPasscodeInstruction,omitempty"`
}

type MessageReactionAdd struct {
	MessageSequenceId *string `thrift:"message_sequence_id,1" json:"message_sequence_id,omitempty"`
	Emoji             *string `thrift:"emoji,2" json:"emoji,omitempty"`
}

type MessageReactionRemove struct {
	MessageSequenceId *string `thrift:"message_sequence_id,1" json:"message_sequence_id,omitempty"`
	Emoji             *string `thrift:"emoji,2" json:"emoji,omitempty"`
}

type MessageTypingEvent struct {
	ConversationId *string `thrift:"conversation_id,1" json:"conversation_id,omitempty"`
}

type MoneyAttachment struct {
	FallbackText *string `thrift:"fallbackText,1" json:"fallbackText,omitempty"`
	Payload      *string `thrift:"payload,2" json:"payload,omitempty"`
}

type MuteConversation struct {
	MutedConversationIds []string `thrift:"muted_conversation_ids,1" json:"muted_conversation_ids,omitempty"`
}

type NicknameMessage struct {
	UserId       *int64  `thrift:"user_id,1" json:"user_id,omitempty"`
	NicknameText *string `thrift:"nickname_text,2" json:"nickname_text,omitempty"`
}

type ParentNode struct {
	SubtreeEncryptionPublicKey *string `thrift:"subtree_encryption_public_key,1" json:"subtree_encryption_public_key,omitempty"`
	ParentHash                 *string `thrift:"parent_hash,2" json:"parent_hash,omitempty"`
}

type PhoneNumberRichTextContent struct {
}

type PinConversation struct {
	ConversationId *string `thrift:"conversation_id,1" json:"conversation_id,omitempty"`
}

type PinReminderInstruction struct {
	ShouldRegister *bool `thrift:"should_register,1" json:"should_register,omitempty"`
	ShouldGenerate *bool `thrift:"should_generate,2" json:"should_generate,omitempty"`
}

type PostAttachment struct {
	RestId       *string `thrift:"rest_id,1" json:"rest_id,omitempty"`
	PostUrl      *string `thrift:"post_url,2" json:"post_url,omitempty"`
	AttachmentId *string `thrift:"attachment_id,3" json:"attachment_id,omitempty"`
}

type PullMessagePageDetails struct {
	MinSequenceId *string `thrift:"min_sequence_id,3" json:"min_sequence_id,omitempty"`
	MaxSequenceId *string `thrift:"max_sequence_id,4" json:"max_sequence_id,omitempty"`
	IsBatchedPull *bool   `thrift:"is_batched_pull,7" json:"is_batched_pull,omitempty"`
}

type PullMessagesFinishedInstruction struct {
	FinishedPull           *bool                   `thrift:"finished_pull,1" json:"finished_pull,omitempty"`
	SequenceContinue       *string                 `thrift:"sequence_continue,2" json:"sequence_continue,omitempty"`
	PullMessagePageDetails *PullMessagePageDetails `thrift:"pull_message_page_details,3" json:"pull_message_page_details,omitempty"`
}

type PullMessagesInstruction struct {
	SequenceStart *string `thrift:"sequence_start,1" json:"sequence_start,omitempty"`
	SenderId      *string `thrift:"sender_id,2" json:"sender_id,omitempty"`
	IsBatchedPull *bool   `thrift:"is_batched_pull,6" json:"is_batched_pull,omitempty"`
}

type QuickReply struct {
	Request  *QuickReplyRequest  `thrift:"request,1" json:"request,omitempty"`
	Response *QuickReplyResponse `thrift:"response,2" json:"response,omitempty"`
}

type QuickReplyOption struct {
	Id          *string `thrift:"id,1" json:"id,omitempty"`
	Label       *string `thrift:"label,2" json:"label,omitempty"`
	Metadata    *string `thrift:"metadata,3" json:"metadata,omitempty"`
	Description *string `thrift:"description,4" json:"description,omitempty"`
}

type QuickReplyOptionsRequest struct {
	Id      *string             `thrift:"id,1" json:"id,omitempty"`
	Options []*QuickReplyOption `thrift:"options,2" json:"options,omitempty"`
}

type QuickReplyOptionsResponse struct {
	RequestId        *string `thrift:"request_id,1" json:"request_id,omitempty"`
	Metadata         *string `thrift:"metadata,2" json:"metadata,omitempty"`
	SelectedOptionId *string `thrift:"selected_option_id,3" json:"selected_option_id,omitempty"`
}

type QuickReplyRequest struct {
	Options *QuickReplyOptionsRequest `thrift:"options,1" json:"options,omitempty"`
}

type QuickReplyResponse struct {
	Options *QuickReplyOptionsResponse `thrift:"options,1" json:"options,omitempty"`
}

type RatchetTree struct {
	Leaves  []*RatchetTreeLeaf   `thrift:"leaves,1" json:"leaves,omitempty"`
	Parents []*RatchetTreeParent `thrift:"parents,2" json:"parents,omitempty"`
}

type RatchetTreeLeaf struct {
	Empty *EmptyNode `thrift:"empty,1" json:"empty,omitempty"`
	Leaf  *LeafNode  `thrift:"leaf,2" json:"leaf,omitempty"`
}

type RatchetTreeParent struct {
	Empty  *EmptyNode  `thrift:"empty,1" json:"empty,omitempty"`
	Parent *ParentNode `thrift:"parent,2" json:"parent,omitempty"`
}

type ReplyingToPreview struct {
	SenderId                    *int64               `thrift:"sender_id,1" json:"sender_id,omitempty"`
	MessageText                 *string              `thrift:"message_text,2" json:"message_text,omitempty"`
	Entities                    []*RichTextEntity    `thrift:"entities,3" json:"entities,omitempty"`
	Attachments                 []*MessageAttachment `thrift:"attachments,4" json:"attachments,omitempty"`
	SenderDisplayName           *string              `thrift:"sender_display_name,5" json:"sender_display_name,omitempty"`
	ReplyingToMessageSequenceId *string              `thrift:"replying_to_message_sequence_id,6" json:"replying_to_message_sequence_id,omitempty"`
	ReplyingToMessageId         *string              `thrift:"replying_to_message_id,7" json:"replying_to_message_id,omitempty"`
}

type RequestForEncryptedResendEvent struct {
	MinSequenceId *string `thrift:"min_sequence_id,1" json:"min_sequence_id,omitempty"`
	MaxSequenceId *string `thrift:"max_sequence_id,2" json:"max_sequence_id,omitempty"`
}

type RichTextContent struct {
	Hashtag     *HashtagRichTextContent     `thrift:"hashtag,1" json:"hashtag,omitempty"`
	Cashtag     *CashtagRichTextContent     `thrift:"cashtag,2" json:"cashtag,omitempty"`
	Mention     *MentionRichTextContent     `thrift:"mention,3" json:"mention,omitempty"`
	Url         *UrlRichTextContent         `thrift:"url,4" json:"url,omitempty"`
	Email       *EmailRichTextContent       `thrift:"email,5" json:"email,omitempty"`
	PhoneNumber *PhoneNumberRichTextContent `thrift:"phoneNumber,7" json:"phoneNumber,omitempty"`
}

type RichTextEntity struct {
	StartIndex *int32           `thrift:"start_index,1" json:"start_index,omitempty"`
	EndIndex   *int32           `thrift:"end_index,2" json:"end_index,omitempty"`
	Content    *RichTextContent `thrift:"content,3" json:"content,omitempty"`
}

type ScreenCaptureDetected struct {
	Type *int32 `thrift:"type,1" json:"type,omitempty"`
}

type SetVerifiedStatus struct {
	UserId         *int64 `thrift:"user_id,1" json:"user_id,omitempty"`
	VerifiedStatus *bool  `thrift:"verified_status,2" json:"verified_status,omitempty"`
}

type StoredGroupState struct {
	Keypairs    []*MaybeKeypair `thrift:"keypairs,1" json:"keypairs,omitempty"`
	RatchetTree *RatchetTree    `thrift:"ratchet_tree,2" json:"ratchet_tree,omitempty"`
}

type StoredKeypair struct {
	PublicKey  *string `thrift:"public_key,1" json:"public_key,omitempty"`
	PrivateKey *string `thrift:"private_key,2" json:"private_key,omitempty"`
}

type SwitchToHybridPullInstruction struct {
	RequestingUserAgent *string `thrift:"requesting_user_agent,1" json:"requesting_user_agent,omitempty"`
}

type UnifiedCardAttachment struct {
	Url          *string `thrift:"url,1" json:"url,omitempty"`
	AttachmentId *string `thrift:"attachment_id,2" json:"attachment_id,omitempty"`
}

type UnmuteConversation struct {
	UnmutedConversationIds []string `thrift:"unmuted_conversation_ids,1" json:"unmuted_conversation_ids,omitempty"`
}

type UnpinConversation struct {
	ConversationId *string `thrift:"conversation_id,1" json:"conversation_id,omitempty"`
}

type UpdatePathNode struct {
	EncryptedSecrets    []string `thrift:"encrypted_secrets,1" json:"encrypted_secrets,omitempty"`
	EncryptedPrivateKey *string  `thrift:"encrypted_private_key,2" json:"encrypted_private_key,omitempty"`
}

type UrlAttachment struct {
	Url                      *string             `thrift:"url,1" json:"url,omitempty"`
	BannerImageMediaHashKey  *UrlAttachmentImage `thrift:"banner_image_media_hash_key,2" json:"banner_image_media_hash_key,omitempty"`
	FaviconImageMediaHashKey *UrlAttachmentImage `thrift:"favicon_image_media_hash_key,3" json:"favicon_image_media_hash_key,omitempty"`
	DisplayTitle             *string             `thrift:"display_title,4" json:"display_title,omitempty"`
	AttachmentId             *string             `thrift:"attachment_id,5" json:"attachment_id,omitempty"`
}

type UrlAttachmentImage struct {
	MediaHashKey  *string          `thrift:"media_hash_key,1" json:"media_hash_key,omitempty"`
	FilesizeBytes *int64           `thrift:"filesize_bytes,2" json:"filesize_bytes,omitempty"`
	Filename      *string          `thrift:"filename,3" json:"filename,omitempty"`
	Dimensions    *MediaDimensions `thrift:"dimensions,4" json:"dimensions,omitempty"`
}

type UrlRichTextContent struct {
}
