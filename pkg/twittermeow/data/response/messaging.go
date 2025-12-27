package response

import (
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type Trusted struct {
	Status     types.PaginationStatus `json:"status,omitempty"`
	MinEntryID string                 `json:"min_entry_id,omitempty"`
}

type Untrusted struct {
	Status types.PaginationStatus `json:"status,omitempty"`
}

type InboxTimelines struct {
	Trusted   Trusted   `json:"trusted,omitempty"`
	Untrusted Untrusted `json:"untrusted,omitempty"`
}

type KeyRegistryState struct {
	Status types.PaginationStatus `json:"status,omitempty"`
}

type DMMessageDeleteMutationResponse struct {
	Data struct {
		DmMessageHideDelete string `json:"dm_message_hide_delete,omitempty"`
	} `json:"data,omitempty"`
}

type PinConversationResponse struct {
	Data struct {
		AddDmConversationLabelV3 struct {
			Typename  string `json:"__typename,omitempty"`
			LabelType string `json:"label_type,omitempty"`
			Timestamp int64  `json:"timestamp,omitempty"`
		} `json:"add_dm_conversation_label_v3,omitempty"`
	}
}

type UnpinConversationResponse struct {
	Data struct {
		DmConversationLabelDelete string `json:"dm_conversation_label_delete,omitempty"`
	} `json:"data,omitempty"`
}

type ReactionResponse struct {
	Data struct {
		DeleteDmReaction struct {
			Typename string `json:"__typename,omitempty"`
		} `json:"delete_dm_reaction,omitempty"`
		CreateDmReaction struct {
			Typename string `json:"__typename,omitempty"`
		} `json:"create_dm_reaction,omitempty"`
	} `json:"data,omitempty"`
}

type AddParticipantsResponse struct {
	Data struct {
		AddParticipants struct {
			Typename   string   `json:"__typename,omitempty"`
			AddedUsers []string `json:"added_users,omitempty"`
		} `json:"add_participants,omitempty"`
	} `json:"data,omitempty"`
}

type SendMessageMutationResponse struct {
	Data struct {
		SendMessage struct {
			Typename string `json:"__typename,omitempty"`
		} `json:"send_message,omitempty"`
	} `json:"data,omitempty"`
	Errors []struct {
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}
