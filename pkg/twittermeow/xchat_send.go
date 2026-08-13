package twittermeow

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (c *Client) SendXChatReadReceipt(ctx context.Context, conversationID, lastReadEventID string, readAt time.Time) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}
	if lastReadEventID == "" {
		return errors.New("last read event ID is required")
	}

	seenAt := time.Now()
	if readAt.IsZero() {
		readAt = seenAt
	}

	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return errors.New("sender ID is required")
	}

	conversationToken, err := c.ensureConversationToken(ctx, conversationID)
	if err != nil {
		return err
	}

	keyPair, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	}
	event, err := buildXChatReadReceiptEvent(
		uuid.NewString(), senderID, conversationID, conversationToken,
		lastReadEventID, readAt, seenAt, keyPair,
	)
	if err != nil {
		return err
	}

	return c.SendXChatPayload(ctx, &payload.Message{MessageEvent: event})
}

func buildXChatReadReceiptEvent(
	messageID, senderID, conversationID, conversationToken, lastReadEventID string,
	createdAt, seenAt time.Time,
	keyPair *crypto.SigningKeyPair,
) (*payload.MessageEvent, error) {
	if keyPair == nil || keyPair.SigningKey == nil || keyPair.KeyVersion == "" {
		return nil, errors.New("complete signing key is required")
	}

	seenAtMillis := seenAt.UnixMilli()
	createdAtMsec := strconv.FormatInt(createdAt.UnixMilli(), 10)
	signature, err := crypto.SignMarkConversationReadEvent(
		keyPair.SigningKey,
		messageID,
		senderID,
		conversationID,
		lastReadEventID,
		seenAtMillis,
	)
	if err != nil {
		return nil, fmt.Errorf("sign read receipt: %w", err)
	}
	signingPublicKey, err := crypto.EncodePublicKeySPKI(&keyPair.SigningKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("encode signing public key: %w", err)
	}

	relaySource := int32(0)
	isGrok := false
	sigVersion := crypto.SignatureVersion7
	return &payload.MessageEvent{
		MessageId:          &messageID,
		SenderId:           &senderID,
		ConversationId:     &conversationID,
		ConversationToken:  &conversationToken,
		CreatedAtMsec:      &createdAtMsec,
		RelaySource:        &relaySource,
		PreviousSequenceId: &lastReadEventID,
		Detail: &payload.MessageEventDetail{
			MarkConversationReadEvent: &payload.MarkConversationReadEvent{
				SeenUntilSequenceId: &lastReadEventID,
				SeenAtMillis:        &seenAtMillis,
				IsGrok:              &isGrok,
			},
		},
		MessageEventSignature: &payload.MessageEventSignature{
			Signature:        &signature,
			PublicKeyVersion: &keyPair.KeyVersion,
			SignatureVersion: &sigVersion,
			SigningPublicKey: &signingPublicKey,
		},
	}, nil
}

func (c *Client) ensureConversationToken(ctx context.Context, conversationID string) (string, error) {
	token, err := c.keyManager.GetConversationToken(ctx, conversationID)
	if err == nil {
		return token, nil
	}
	if err != nil && !errors.Is(err, crypto.ErrKeyNotFound) {
		return "", fmt.Errorf("get conversation token: %w", err)
	}

	if err := c.refreshConversationToken(ctx, conversationID); err != nil && !errors.Is(err, crypto.ErrKeyNotFound) {
		return "", err
	}

	token, err = c.keyManager.GetConversationToken(ctx, conversationID)
	if err != nil {
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
	for _, readEvt := range item.LatestReadEventsPerParticipant {
		encoded = append(encoded, readEvt.LatestMarkConversationReadEvent)
	}

	for _, encodedEvt := range encoded {
		err := c.putConversationTokenFromEncodedEvent(ctx, conversationID, encodedEvt)
		if err == nil {
			return nil
		}
		if !errors.Is(err, crypto.ErrKeyNotFound) {
			return err
		}
	}

	return crypto.ErrKeyNotFound
}

func (c *Client) putConversationTokenFromEncodedEvent(ctx context.Context, conversationID, encoded string) error {
	if encoded == "" {
		return crypto.ErrKeyNotFound
	}
	evt, err := DecodeMessageEvent(encoded)
	if err != nil {
		return crypto.ErrKeyNotFound
	}
	if evt == nil || evt.ConversationToken == nil || *evt.ConversationToken == "" {
		return crypto.ErrKeyNotFound
	}
	if evt.ConversationId != nil && *evt.ConversationId != "" {
		conversationID = *evt.ConversationId
	}
	return c.keyManager.PutConversationToken(ctx, conversationID, *evt.ConversationToken)
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

	_, err = c.sendMessageMutation(ctx, pl)
	return err
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

	_, err = c.sendMessageMutation(ctx, pl)
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

	// Sync room data even if one key event was malformed: the response may still
	// contain complete member/profile data and other usable keys.
	keyErr := c.processKeyChangeEventsFromItem(ctx, conversationID, &item)
	c.notifyConversationDataRefresh(ctx, conversationID, item)
	return keyErr
}

func (c *Client) notifyConversationDataRefresh(ctx context.Context, conversationID string, item response.XChatInboxItem) {
	callback := c.onConversationDataRefresh
	if callback == nil {
		return
	}
	// Room refresh may queue on the portal currently converting this message, so never wait inline.
	go callback(ctx, conversationID, &item)
}

func (c *Client) processKeyChangeEventsFromItem(ctx context.Context, conversationID string, item *response.XChatInboxItem) error {
	if conversationID == "" {
		return errors.New("conversation ID is empty")
	}
	if item == nil || len(item.LatestConversationKeyChangeEvents) == 0 {
		return nil
	}

	signingKey, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	}
	if signingKey == nil || signingKey.DecryptKeyB64 == "" {
		return errors.New("own XChat decryption key is missing")
	}
	ownUserID := c.GetCurrentUserID()
	if ownUserID == "" {
		return fmt.Errorf("current user ID is empty")
	}

	var eventErrs []error
	for index, encoded := range item.LatestConversationKeyChangeEvents {
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			eventErrs = append(eventErrs, fmt.Errorf("decode conversation key event %d: %w", index, err))
			continue
		}
		var evt payload.MessageEvent
		if err := payload.Decode(data, &evt); err != nil {
			eventErrs = append(eventErrs, fmt.Errorf("decode conversation key event %d thrift: %w", index, err))
			continue
		}
		if evt.Detail == nil {
			eventErrs = append(eventErrs, fmt.Errorf("conversation key event %d has no detail", index))
			continue
		}
		eventConversationID := ptr.Val(evt.ConversationId)
		if eventConversationID != "" && eventConversationID != conversationID {
			eventErrs = append(eventErrs, fmt.Errorf("conversation key event %d belongs to a different conversation", index))
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
			eventErrs = append(eventErrs, fmt.Errorf("unwrap conversation key event %d: %w", index, err))
			continue
		}
		keyVersion := ptr.Val(ckce.ConversationKeyVersion)
		if keyVersion == "" {
			eventErrs = append(eventErrs, fmt.Errorf("conversation key event %d has no key version", index))
			continue
		}
		keyCreatedAt := methods.ParseMsecTimestamp(ptr.Val(evt.CreatedAtMsec))
		if keyCreatedAt.IsZero() {
			c.Logger.Warn().
				Str("conversation_id", conversationID).
				Str("key_version", keyVersion).
				Str("created_at_msec", ptr.Val(evt.CreatedAtMsec)).
				Msg("Skipping conversation key update without valid XChat timestamp")
			eventErrs = append(eventErrs, fmt.Errorf("conversation key event %d has no valid timestamp", index))
			continue
		}

		err = c.keyManager.PutConversationKey(ctx, &crypto.ConversationKey{
			ConversationID: conversationID,
			KeyVersion:     keyVersion,
			Key:            convKeyBytes,
			CreatedAt:      keyCreatedAt,
		})
		if err != nil {
			eventErrs = append(eventErrs, fmt.Errorf("store conversation key event %d: %w", index, err))
		}
	}
	return errors.Join(eventErrs...)
}
