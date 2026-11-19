package crypto

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"strconv"

	"go.mau.fi/util/ptr"
	"golang.org/x/crypto/hkdf"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

const (
	messageKeySalt    = "xchat-message-key"
	attachmentKeySalt = "xchat-attachment-key"
)

// UsesMessageKey reports whether a MessageCreateEvent's contents and media are
// encrypted under a key derived for that message instead of the conversation key.
func UsesMessageKey(mce *payload.MessageCreateEvent) bool {
	return mce != nil && ptr.Val(mce.ContentsEncryptionVersion) >= 1
}

// DeriveMessageKey derives the key of a single message from its conversation key.
func DeriveMessageKey(conversationKey []byte, senderID, messageID string) ([]byte, error) {
	sender, err := strconv.ParseInt(senderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("message key needs a numeric sender ID: %w", err)
	} else if messageID == "" {
		return nil, errors.New("message key needs a message ID")
	}
	return deriveKey(conversationKey, messageKeySalt, strconv.FormatInt(sender, 10)+":"+messageID)
}

// DeriveAttachmentKey derives the key of one attachment from its message key.
func DeriveAttachmentKey(messageKey []byte, attachmentID string) ([]byte, error) {
	return deriveKey(messageKey, attachmentKeySalt, attachmentID)
}

func deriveKey(secret []byte, salt, info string) ([]byte, error) {
	key := make([]byte, secretboxKeySize)
	_, err := io.ReadFull(hkdf.New(sha256.New, secret, []byte(salt), []byte(info)), key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// MessageContentsKey returns the key the contents of a MessageCreateEvent are encrypted under.
func MessageContentsKey(conversationKey []byte, evt *payload.MessageEvent, mce *payload.MessageCreateEvent) ([]byte, error) {
	if !UsesMessageKey(mce) {
		return conversationKey, nil
	}
	return DeriveMessageKey(conversationKey, ptr.Val(evt.SenderId), ptr.Val(evt.MessageId))
}

// AttachmentKeyInfo identifies the key of an attachment in a message that uses a message key.
type AttachmentKeyInfo struct {
	SenderID     string
	MessageID    string
	AttachmentID string
}

// NewAttachmentKeyInfo returns the key info of the attachment at index in a message,
// or nil if the message's media is encrypted under the conversation key.
// The index doesn't count payment receipts.
func NewAttachmentKeyInfo(evt *payload.MessageEvent, mce *payload.MessageCreateEvent, ownID *string, index int) *AttachmentKeyInfo {
	if !UsesMessageKey(mce) {
		return nil
	}
	messageID := ptr.Val(evt.MessageId)
	attachmentID := messageID
	if ownID != nil {
		attachmentID = *ownID
	} else if index > 0 {
		attachmentID = messageID + "#" + strconv.Itoa(index)
	}
	return &AttachmentKeyInfo{
		SenderID:     ptr.Val(evt.SenderId),
		MessageID:    messageID,
		AttachmentID: attachmentID,
	}
}

// Key returns the key of the attachment, or the conversation key if info is nil.
func (info *AttachmentKeyInfo) Key(conversationKey []byte) ([]byte, error) {
	if info == nil {
		return conversationKey, nil
	}
	messageKey, err := DeriveMessageKey(conversationKey, info.SenderID, info.MessageID)
	if err != nil {
		return nil, err
	}
	return DeriveAttachmentKey(messageKey, info.AttachmentID)
}
