package twittermeow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (c *Client) SendXChatReadReceipt(ctx context.Context, conversationID, lastReadEventID string, readAt time.Time) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}
	if lastReadEventID == "" {
		return errors.New("last read event ID is required")
	}

	if readAt.IsZero() {
		readAt = time.Now()
	}
	readAtMillis := readAt.UnixMilli()
	createdAtMsec := strconv.FormatInt(readAtMillis, 10)

	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return errors.New("sender ID is required")
	}

	conversationToken, err := c.ensureConversationToken(ctx, conversationID)
	if err != nil {
		return err
	}

	messageID := uuid.NewString()
	readEvent := &payload.MarkConversationReadEvent{
		SeenUntilSequenceId: &lastReadEventID,
		SeenAtMillis:        &readAtMillis,
	}
	event := &payload.MessageEvent{
		MessageId:         &messageID,
		SenderId:          &senderID,
		ConversationId:    &conversationID,
		ConversationToken: &conversationToken,
		CreatedAtMsec:     &createdAtMsec,
		Detail: &payload.MessageEventDetail{
			MarkConversationReadEvent: readEvent,
		},
	}

	if keyPair, err := c.keyManager.GetOwnSigningKey(ctx); err == nil && keyPair != nil && keyPair.SigningKey != nil && keyPair.KeyVersion != "" {
		signature, err := crypto.SignMarkConversationReadEvent(
			keyPair.SigningKey,
			messageID,
			senderID,
			conversationID,
			conversationToken,
			createdAtMsec,
			lastReadEventID,
			readAtMillis,
		)
		if err != nil {
			return fmt.Errorf("sign read receipt: %w", err)
		}
		sigVersion := crypto.SignatureVersion4
		event.MessageEventSignature = &payload.MessageEventSignature{
			Signature:        &signature,
			PublicKeyVersion: &keyPair.KeyVersion,
			SignatureVersion: &sigVersion,
		}
	}

	return c.SendXChatPayload(ctx, &payload.Message{MessageEvent: event})
}

func (c *Client) ensureConversationToken(ctx context.Context, conversationID string) (string, error) {
	token, err := c.keyManager.GetConversationToken(ctx, conversationID)
	if err == nil && token != "" {
		return token, nil
	}
	if err != nil && !errors.Is(err, crypto.ErrKeyNotFound) {
		return "", fmt.Errorf("get conversation token: %w", err)
	}

	if err := c.refreshConversationToken(ctx, conversationID); err != nil && !errors.Is(err, crypto.ErrKeyNotFound) {
		return "", err
	}

	token, err = c.keyManager.GetConversationToken(ctx, conversationID)
	if err != nil || token == "" {
		if err == nil {
			err = crypto.ErrKeyNotFound
		}
		return "", fmt.Errorf("get conversation token: %w", err)
	}

	return token, nil
}

func (c *Client) refreshConversationToken(ctx context.Context, conversationID string) error {
	vars := payload.NewInboxPageConversationDataQueryVariables(conversationID, false)
	resp, err := c.GetConversationData(ctx, vars)
	if err != nil {
		return fmt.Errorf("fetch conversation data: %w", err)
	}

	item := resp.Data.GetInboxPageConversationData.Data
	if err := c.storeConversationTokenFromInboxItem(ctx, conversationID, &item); err != nil {
		if errors.Is(err, crypto.ErrKeyNotFound) {
			return err
		}
		return fmt.Errorf("store conversation token: %w", err)
	}

	return nil
}

func (c *Client) storeConversationTokenFromInboxItem(ctx context.Context, conversationID string, item *response.XChatInboxItem) error {
	if item == nil {
		return crypto.ErrKeyNotFound
	}

	encoded := make([]string, 0, len(item.LatestMessageEvents)+len(item.EncodedMessageEvents)+len(item.LatestConversationKeyChangeEvents)+2)
	encoded = append(encoded, item.LatestMessageEvents...)
	encoded = append(encoded, item.EncodedMessageEvents...)
	encoded = append(encoded, item.LatestConversationKeyChangeEvents...)
	if item.LatestNotifiableMessageCreateEvent != "" {
		encoded = append(encoded, item.LatestNotifiableMessageCreateEvent)
	}
	if item.ConversationDetail.LatestGroupTitleChangeMessageEvent != "" {
		encoded = append(encoded, item.ConversationDetail.LatestGroupTitleChangeMessageEvent)
	}

	for _, encodedEvt := range encoded {
		ok, err := c.storeConversationTokenFromEncodedEvent(ctx, conversationID, encodedEvt)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	for _, readEvt := range item.LatestReadEventsPerParticipant {
		ok, err := c.storeConversationTokenFromEncodedEvent(ctx, conversationID, readEvt.LatestMarkConversationReadEvent)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	return crypto.ErrKeyNotFound
}

func (c *Client) storeConversationTokenFromEncodedEvent(ctx context.Context, fallbackConversationID, encoded string) (bool, error) {
	if encoded == "" {
		return false, nil
	}
	evt, err := DecodeMessageEvent(encoded)
	if err != nil {
		return false, nil
	}
	if evt == nil || evt.ConversationToken == nil || *evt.ConversationToken == "" {
		return false, nil
	}

	conversationID := fallbackConversationID
	if evt.ConversationId != nil && *evt.ConversationId != "" {
		conversationID = *evt.ConversationId
	}

	if err := c.keyManager.PutConversationToken(ctx, conversationID, *evt.ConversationToken); err != nil {
		return false, err
	}
	return true, nil
}

// getSelfConversationID returns the user's self-conversation ID (user_id:user_id format).
func (c *Client) getSelfConversationID() string {
	userID := c.GetCurrentUserID()
	return userID + ":" + userID
}

// SendXChatPinConversation pins a conversation via XChat.
func (c *Client) SendXChatPinConversation(ctx context.Context, targetConversationID string) error {
	if targetConversationID == "" {
		return errors.New("target conversation ID is required")
	}

	selfConvID := c.getSelfConversationID()

	token, err := c.ensureConversationToken(ctx, selfConvID)
	if err != nil {
		return fmt.Errorf("get self conversation token: %w", err)
	}

	messageID := uuid.NewString()

	builder := crypto.NewMessageBuilder(c.keyManager, c.GetCurrentUserID()).
		SetMessageID(messageID).
		SetConversationID(selfConvID).
		SetPinConversation(targetConversationID)

	encodedMCE, encodedSig, err := builder.BuildForSend(ctx)
	if err != nil {
		return fmt.Errorf("build pin message: %w", err)
	}

	var sigPtr *string
	if encodedSig != "" {
		sigPtr = &encodedSig
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               selfConvID,
		MessageID:                    messageID,
		ConversationToken:            token,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: sigPtr,
	})

	return c.sendMessageMutation(ctx, pl)
}

// SendXChatUnpinConversation unpins a conversation via XChat.
func (c *Client) SendXChatUnpinConversation(ctx context.Context, targetConversationID string) error {
	if targetConversationID == "" {
		return errors.New("target conversation ID is required")
	}

	selfConvID := c.getSelfConversationID()

	token, err := c.ensureConversationToken(ctx, selfConvID)
	if err != nil {
		return fmt.Errorf("get self conversation token: %w", err)
	}

	messageID := uuid.NewString()

	builder := crypto.NewMessageBuilder(c.keyManager, c.GetCurrentUserID()).
		SetMessageID(messageID).
		SetConversationID(selfConvID).
		SetUnpinConversation(targetConversationID)

	encodedMCE, encodedSig, err := builder.BuildForSend(ctx)
	if err != nil {
		return fmt.Errorf("build unpin message: %w", err)
	}

	var sigPtr *string
	if encodedSig != "" {
		sigPtr = &encodedSig
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               selfConvID,
		MessageID:                    messageID,
		ConversationToken:            token,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: sigPtr,
	})

	return c.sendMessageMutation(ctx, pl)
}

// sendMessageMutation sends a SendMessageMutation request.
func (c *Client) sendMessageMutation(ctx context.Context, pl *payload.SendMessageMutationPayload) error {
	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("SendMessageMutation payload")

	_, _, err = c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.SEND_MESSAGE_MUTATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	return err
}

// SendXChatTypingNotification sends a typing indicator via XChat WebSocket.
func (c *Client) SendXChatTypingNotification(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}

	conversationToken, err := c.ensureConversationToken(ctx, conversationID)
	if err != nil {
		return err
	}

	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return errors.New("sender ID is required")
	}

	messageID := uuid.NewString()
	createdAtMsec := strconv.FormatInt(time.Now().UnixMilli(), 10)

	event := &payload.MessageEvent{
		MessageId:         &messageID,
		SenderId:          &senderID,
		ConversationId:    &conversationID,
		ConversationToken: &conversationToken,
		CreatedAtMsec:     &createdAtMsec,
		Detail: &payload.MessageEventDetail{
			MessageTypingEvent: &payload.MessageTypingEvent{
				ConversationId: &conversationID,
			},
		},
	}

	return c.SendXChatPayload(ctx, &payload.Message{MessageEvent: event})
}

// RefreshConversationKeys fetches conversation data and processes key change events.
// Called when message decryption fails due to missing keys.
// Also invokes the conversation data callback (if set) to sync room data.
func (c *Client) RefreshConversationKeys(ctx context.Context, conversationID string) error {
	vars := payload.NewInboxPageConversationDataQueryVariables(conversationID, true)
	resp, err := c.GetConversationData(ctx, vars)
	if err != nil {
		return fmt.Errorf("fetch conversation data: %w", err)
	}

	item := resp.Data.GetInboxPageConversationData.Data

	// Process key change events to store conversation keys
	if err := c.processKeyChangeEventsFromItem(ctx, conversationID, &item); err != nil {
		return err
	}

	// Notify callback to sync room data (members, name, avatar, etc.)
	if c.onConversationDataRefresh != nil {
		c.onConversationDataRefresh(ctx, conversationID, &item)
	}

	return nil
}

func (c *Client) processKeyChangeEventsFromItem(ctx context.Context, conversationID string, item *response.XChatInboxItem) error {
	if item == nil || len(item.LatestConversationKeyChangeEvents) == 0 {
		return nil
	}

	signingKey, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	}
	ownUserID := c.GetCurrentUserID()
	if ownUserID == "" {
		return fmt.Errorf("current user ID is empty")
	}

	for _, encoded := range item.LatestConversationKeyChangeEvents {
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			continue
		}
		var evt payload.MessageEvent
		if err := payload.Decode(data, &evt); err != nil {
			continue
		}
		ckce := evt.Detail.ConversationKeyChangeEvent
		if ckce == nil {
			continue
		}

		var ourEncryptedKey string
		for _, pk := range ckce.ConversationParticipantKeys {
			if ptr.Val(pk.UserId) == ownUserID {
				ourEncryptedKey = ptr.Val(pk.EncryptedConversationKey)
				break
			}
		}
		if ourEncryptedKey == "" {
			continue
		}

		convKeyBytes, err := crypto.UnwrapConversationKey(ourEncryptedKey, signingKey.DecryptKeyB64)
		if err != nil {
			continue
		}
		c.keyManager.PutConversationKey(ctx, &crypto.ConversationKey{
			ConversationID: conversationID,
			KeyVersion:     ptr.Val(ckce.ConversationKeyVersion),
			Key:            convKeyBytes,
			CreatedAt:      time.Now(),
		})
	}
	return nil
}
