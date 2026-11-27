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
	ConversationIDs []string
	Archived        bool
	Timestamp       time.Time
}

func (*XChatConversationArchiveChange) isTwitterEvent() {}

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
	ActivityType   payload.MessageActivityType
}

func (*XChatTyping) isTwitterEvent() {}

// XChatUnsupportedMessage is emitted for a message whose contents use an encryption version
// the bridge can't read.
type XChatUnsupportedMessage struct {
	SequenceID     string
	ConversationID string
	SenderID       string
	Time           string
	Version        int16
}

func (*XChatUnsupportedMessage) isTwitterEvent() {}
