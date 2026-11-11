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
	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return errors.New("sender ID is required")
	}
	convKey, err := c.keyManager.GetLatestConversationKey(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation key: %w", err)
	} else if convKey == nil || len(convKey.Key) == 0 || convKey.KeyVersion == "" {
		return crypto.ErrKeyNotFound
	}
	keyPair, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	} else if keyPair == nil || keyPair.SigningKey == nil || keyPair.KeyVersion == "" {
		return errors.New("missing signing key")
	}

	encryptedTitle, err := crypto.SecretboxEncrypt([]byte(title), convKey.Key)
	if err != nil {
		return fmt.Errorf("encrypt group title: %w", err)
	}
	encodedTitle := base64.StdEncoding.EncodeToString(encryptedTitle)

	messageID := uuid.NewString()
	encodedDetail, err := crypto.EncodeMessageEventDetail(&payload.MessageEventDetail{
		GroupChangeEvent: &payload.GroupChangeEvent{
			GroupChange: &payload.GroupChange{
				GroupTitleChange: &payload.GroupTitleChange{
					CustomTitle:            &encodedTitle,
					ConversationKeyVersion: &convKey.KeyVersion,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("encode group title change: %w", err)
	}
	signaturePayload, err := crypto.SignaturePayloadGroupChange("GroupTitleChange", messageID, senderID, conversationID, convKey.KeyVersion, encodedTitle)
	if err != nil {
		return err
	}
	signature, err := crypto.Sign(keyPair.SigningKey, []byte(signaturePayload))
	if err != nil {
		return fmt.Errorf("sign group title change: %w", err)
	}

	pl := payload.NewEditGroupConversationMutationPayload(payload.EditGroupConversationMutationVariables{
		ConversationID:         conversationID,
		ConversationKeyVersion: convKey.KeyVersion,
		GroupName:              &encodedTitle,
		ModifiedFields:         []payload.EditGroupFieldName{payload.EditGroupFieldTitle},
		ActionSignatures: []payload.DeleteMessageActionSignature{{
			MessageID:                 messageID,
			EncodedMessageEventDetail: encodedDetail,
			SignaturePayload:          signaturePayload,
			MessageEventSignature: &payload.DeleteMessageEventSignatureJSON{
				Signature:        signature,
				PublicKeyVersion: keyPair.KeyVersion,
				SignatureVersion: crypto.SignatureVersion7,
			},
		}},
	})
	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return err
	}
	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.EDIT_GROUP_CONVERSATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return err
	}
	return parseEditGroupConversationResponse(respBody)
}

func parseEditGroupConversationResponse(respBody []byte) error {
	var resp struct {
		Data struct {
			EditGroup *struct {
				Typename  string `json:"__typename"`
				ErrorCode string `json:"error_code"`
			} `json:"xchat_edit_group"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("decode edit group response: %w", err)
	} else if len(resp.Errors) > 0 {
		return fmt.Errorf("edit group error: %s", resp.Errors[0].Message)
	} else if resp.Data.EditGroup == nil {
		return errors.New("edit group returned no result")
	} else if resp.Data.EditGroup.ErrorCode != "" {
		return fmt.Errorf("edit group failed: %s", resp.Data.EditGroup.ErrorCode)
	}
	return nil
}
