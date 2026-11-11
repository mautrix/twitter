package twittermeow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// EditXChatGroupTitle renames an X Chat group. The title is encrypted with the latest conversation key.
func (c *Client) EditXChatGroupTitle(ctx context.Context, conversationID, title string) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}
	convKey, err := c.keyManager.GetLatestConversationKey(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation key: %w", err)
	} else if convKey == nil || len(convKey.Key) == 0 || convKey.KeyVersion == "" {
		return crypto.ErrKeyNotFound
	}
	encryptedTitle, err := crypto.SecretboxEncrypt([]byte(title), convKey.Key)
	if err != nil {
		return fmt.Errorf("encrypt group title: %w", err)
	}
	encodedTitle := base64.StdEncoding.EncodeToString(encryptedTitle)
	return c.editXChatGroup(ctx, payload.EditGroupConversationMutationVariables{
		ConversationID:         conversationID,
		ConversationKeyVersion: convKey.KeyVersion,
		GroupName:              &encodedTitle,
		ModifiedFields:         []payload.EditGroupFieldName{payload.EditGroupFieldTitle},
	}, "GroupTitleChange", &payload.GroupChange{
		GroupTitleChange: &payload.GroupTitleChange{
			CustomTitle:            &encodedTitle,
			ConversationKeyVersion: &convKey.KeyVersion,
		},
	}, convKey.KeyVersion, encodedTitle)
}

// EditXChatGroupAvatar uploads a new X Chat group avatar, encrypted with the latest conversation key, and sets it.
func (c *Client) EditXChatGroupAvatar(ctx context.Context, conversationID string, avatar []byte) (string, error) {
	if conversationID == "" {
		return "", errors.New("conversation ID is required")
	}
	uploaded, err := c.UploadXChatMedia(ctx, conversationID, fmt.Sprintf("avatar-%s-%s", conversationID, uuid.NewString()), avatar)
	if err != nil {
		return "", fmt.Errorf("upload group avatar: %w", err)
	}
	avatarURL := fmt.Sprintf("%s/%s/%s", endpoints.XCHAT_MEDIA_BASE_URL, conversationID, uploaded.MediaHashKey)
	return avatarURL, c.editXChatGroup(ctx, payload.EditGroupConversationMutationVariables{
		ConversationID:         conversationID,
		ConversationKeyVersion: uploaded.KeyVersion,
		AvatarURL:              &avatarURL,
		ModifiedFields:         []payload.EditGroupFieldName{payload.EditGroupFieldAvatar},
	}, "GroupAvatarUrlChange", &payload.GroupChange{
		GroupAvatarChange: &payload.GroupAvatarUrlChange{
			CustomAvatarUrl:        &avatarURL,
			ConversationKeyVersion: &uploaded.KeyVersion,
		},
	}, uploaded.KeyVersion, avatarURL)
}

func (c *Client) editXChatGroup(
	ctx context.Context,
	vars payload.EditGroupConversationMutationVariables,
	changeName string,
	change *payload.GroupChange,
	signatureFields ...string,
) error {
	sig, err := c.signXChatGroupChange(ctx, vars.ConversationID, changeName, change, signatureFields...)
	if err != nil {
		return err
	}
	vars.ActionSignatures = []payload.DeleteMessageActionSignature{*sig}
	respBody, err := c.postXChatGroupMutation(ctx, endpoints.EDIT_GROUP_CONVERSATION_URL, payload.NewEditGroupConversationMutationPayload(vars))
	if err != nil {
		return err
	}
	return parseEditGroupConversationResponse(respBody)
}

// SetXChatGroupAdmin makes a group member an admin or removes their admin status.
func (c *Client) SetXChatGroupAdmin(ctx context.Context, conversationID, userID string, admin bool) error {
	if conversationID == "" || userID == "" {
		return errors.New("conversation ID and user ID are required")
	}
	var (
		changeName  string
		change      payload.GroupChange
		url         string
		pl          *payload.GroupAdminMutationPayload
		resultField string
	)
	vars := payload.GroupAdminMutationVariables{ConversationID: conversationID, UserID: userID}
	if admin {
		changeName = "GroupAdminAddChange"
		change.GroupAdminAdd = &payload.GroupAdminAddChange{AdminIds: []string{userID}}
		url = endpoints.ADD_GROUP_ADMIN_URL
		pl = payload.NewGroupAdminMutationPayload("AddAsAdminMutation", "1RrduaHSFpenZuOQzgSZqA", vars)
		resultField = "xchat_add_group_admins"
	} else {
		changeName = "GroupAdminRemoveChange"
		change.GroupAdminRemove = &payload.GroupAdminRemoveChange{AdminIds: []string{userID}}
		url = endpoints.REMOVE_GROUP_ADMIN_URL
		pl = payload.NewGroupAdminMutationPayload("RemoveAsAdminMutation", "i3RDeN7PE4lDpu7dpPjbZw", vars)
		resultField = "xchat_remove_group_admins"
	}
	sig, err := c.signXChatGroupChange(ctx, conversationID, changeName, &change, userID)
	if err != nil {
		return err
	}
	pl.Variables.ActionSignatures = []payload.DeleteMessageActionSignature{*sig}
	respBody, err := c.postXChatGroupMutation(ctx, url, pl)
	if err != nil {
		return err
	}
	return parseXChatGroupMutationResponse(respBody, resultField)
}

func (c *Client) signXChatGroupChange(
	ctx context.Context,
	conversationID, changeName string,
	change *payload.GroupChange,
	signatureFields ...string,
) (*payload.DeleteMessageActionSignature, error) {
	return c.signXChatEvent(ctx, &payload.MessageEventDetail{
		GroupChangeEvent: &payload.GroupChangeEvent{GroupChange: change},
	}, "GroupChangeEvent."+changeName, conversationID, signatureFields...)
}

func (c *Client) signXChatEvent(
	ctx context.Context,
	detail *payload.MessageEventDetail,
	eventName, conversationID string,
	signatureFields ...string,
) (*payload.DeleteMessageActionSignature, error) {
	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return nil, errors.New("sender ID is required")
	}
	keyPair, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("get signing key: %w", err)
	} else if keyPair == nil || keyPair.SigningKey == nil || keyPair.KeyVersion == "" {
		return nil, errors.New("missing signing key")
	}

	messageID := uuid.NewString()
	encodedDetail, err := crypto.EncodeMessageEventDetail(detail)
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", eventName, err)
	}
	signaturePayload, err := crypto.SignaturePayloadV7(eventName, messageID, senderID, conversationID, signatureFields...)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(keyPair.SigningKey, []byte(signaturePayload))
	if err != nil {
		return nil, fmt.Errorf("sign %s: %w", eventName, err)
	}
	// The payload isn't sent: key change payloads contain the new conversation key.
	return &payload.DeleteMessageActionSignature{
		MessageID:                 messageID,
		EncodedMessageEventDetail: encodedDetail,
		MessageEventSignature: &payload.DeleteMessageEventSignatureJSON{
			Signature:        signature,
			PublicKeyVersion: keyPair.KeyVersion,
			SignatureVersion: crypto.SignatureVersion7,
		},
	}, nil
}

func (c *Client) postXChatGroupMutation(ctx context.Context, url string, pl any) ([]byte, error) {
	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}
	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
		// A server error doesn't mean the change wasn't applied, and resending it would sign a
		// second change against state that may already be out of date.
		NoRetry: true,
	})
	return respBody, err
}

func parseEditGroupConversationResponse(respBody []byte) error {
	return parseXChatGroupMutationResponse(respBody, "xchat_edit_group")
}

func parseXChatGroupMutationResponse(respBody []byte, resultField string) error {
	var resp struct {
		Data map[string]*struct {
			Typename  string `json:"__typename"`
			ErrorCode string `json:"error_code"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("decode %s response: %w", resultField, err)
	} else if len(resp.Errors) > 0 {
		return fmt.Errorf("%s error: %s", resultField, resp.Errors[0].Message)
	} else if result := resp.Data[resultField]; result == nil {
		return fmt.Errorf("%s returned no result", resultField)
	} else if result.ErrorCode != "" {
		return fmt.Errorf("%s failed: %s", resultField, result.ErrorCode)
	}
	return nil
}
