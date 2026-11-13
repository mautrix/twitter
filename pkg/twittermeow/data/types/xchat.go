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

// XChatConversationTTLChange is emitted when the disappearing message timer of a conversation changes.
type XChatConversationTTLChange struct {
	ID             string
	ConversationID string
	SenderID       string
	TTL            time.Duration
	Timestamp      time.Time
}

func (*XChatConversationTTLChange) isTwitterEvent() {}

// XChatTyping represents a typing indicator from XChat.
type XChatTyping struct {
	ConversationID string
	SenderID       string
	Timestamp      time.Time
}

func (*XChatTyping) isTwitterEvent() {}
