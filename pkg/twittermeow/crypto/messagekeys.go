package crypto

import (
	"crypto/hkdf"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"

	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

const (
	messageKeySalt    = "xchat-message-key"
	attachmentKeySalt = "xchat-attachment-key"
)

func usesMessageKey(mce *payload.MessageCreateEvent) bool {
	return mce != nil && ptr.Val(mce.ContentsEncryptionVersion) >= 1
}

func deriveMessageKey(conversationKey []byte, senderID, messageID string) ([]byte, error) {
	sender, err := strconv.ParseInt(senderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("message key needs a numeric sender ID: %w", err)
	} else if messageID == "" {
		return nil, errors.New("message key needs a message ID")
	}
	return hkdf.Key(sha256.New, conversationKey, []byte(messageKeySalt), strconv.FormatInt(sender, 10)+":"+messageID, secretboxKeySize)
}

func MessageContentsKey(conversationKey []byte, evt *payload.MessageEvent, mce *payload.MessageCreateEvent) ([]byte, error) {
	if !usesMessageKey(mce) {
		return conversationKey, nil
	}
	return deriveMessageKey(conversationKey, ptr.Val(evt.SenderId), ptr.Val(evt.MessageId))
}

type AttachmentKeyInfo struct {
	SenderID     string
	MessageID    string
	AttachmentID string
}

// NewAttachmentKeyInfo returns nil for conversation-key media; index excludes payment attachments.
func NewAttachmentKeyInfo(evt *payload.MessageEvent, mce *payload.MessageCreateEvent, ownID *string, index int) *AttachmentKeyInfo {
	if !usesMessageKey(mce) {
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

func (info *AttachmentKeyInfo) Key(conversationKey []byte) ([]byte, error) {
	if info == nil {
		return conversationKey, nil
	}
	messageKey, err := deriveMessageKey(conversationKey, info.SenderID, info.MessageID)
	if err != nil {
		return nil, err
	}
	return hkdf.Key(sha256.New, messageKey, []byte(attachmentKeySalt), info.AttachmentID, secretboxKeySize)
}
