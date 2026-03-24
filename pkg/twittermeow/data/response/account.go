package response

import "go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"

type GetDMPermissionsResponse struct {
	Permissions Permissions           `json:"permissions,omitempty"`
	Users       map[string]types.User `json:"users,omitempty"`
}

type PermissionDetails struct {
	CanDM     bool `json:"can_dm,omitempty"`
	ErrorCode int  `json:"error_code,omitempty"`
}

type Permissions struct {
	IDKeys map[string]PermissionDetails `json:"id_keys,omitempty"`
}

func (perms Permissions) GetPermissionsForUser(userID string) *PermissionDetails {
	if user, ok := perms.IDKeys[userID]; ok {
		return &user
	}

	return nil
}
