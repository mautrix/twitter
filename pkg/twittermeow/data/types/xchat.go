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

// XChatConversationArchiveChange is emitted when conversations are archived or unarchived on another device.
type XChatConversationArchiveChange struct {
	SequenceID      string
	ConversationIDs []string
	Archived        bool
	Timestamp       time.Time
}

func (*XChatConversationArchiveChange) isTwitterEvent() {}

// XChatConversationMarkedUnread is emitted when the user marks a conversation as unread on another device.
type XChatConversationMarkedUnread struct {
	// FromSnapshot is set when the event was replayed from an inbox snapshot rather than received live.
	FromSnapshot   bool
	SequenceID     string
	ConversationID string
	Timestamp      time.Time
}

func (*XChatConversationMarkedUnread) isTwitterEvent() {}

// XChatPinnedConversations is emitted with the full list of pinned conversations from the user metadata.
type XChatPinnedConversations struct {
	SequenceID      string
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
