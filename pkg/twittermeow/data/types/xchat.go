package types

import (
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

// XChatKeyChange is emitted when a conversation key changes.
// This is informational - the key has already been stored by the processor.
type XChatKeyChange struct {
	ID             string
	ConversationID string
	SenderID       string
	NewKeyVersion  string
	Timestamp      time.Time
}

func (*XChatKeyChange) isTwitterEvent() {}

// XChatMessageFailure is emitted when a message send fails.
type XChatMessageFailure struct {
	ConversationID string
	MessageID      string
	FailureType    payload.FailureType
	Timestamp      time.Time
}

func (*XChatMessageFailure) isTwitterEvent() {}

// XChatConversationMuteChange is emitted when conversations are muted or unmuted on another device.
type XChatConversationMuteChange struct {
	ConversationIDs []string
	Muted           bool
	Timestamp       time.Time
}

func (*XChatConversationMuteChange) isTwitterEvent() {}

// XChatPinnedConversations is emitted with the full list of pinned conversations from the user metadata.
type XChatPinnedConversations struct {
	ConversationIDs []string
	Timestamp       time.Time
}

func (*XChatPinnedConversations) isTwitterEvent() {}

// XChatTyping represents a typing indicator from XChat.
type XChatTyping struct {
	ConversationID string
	SenderID       string
	Timestamp      time.Time
}

func (*XChatTyping) isTwitterEvent() {}
