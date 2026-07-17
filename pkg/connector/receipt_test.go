package connector

import (
	"testing"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func TestConversationReadSenderID(t *testing.T) {
	remote := &types.ConversationRead{SenderID: "remote-user"}
	if got := conversationReadSenderID(remote, "self-user"); got != "remote-user" {
		t.Fatalf("remote sender = %q", got)
	}

	legacy := &types.ConversationRead{}
	if got := conversationReadSenderID(legacy, "self-user"); got != "self-user" {
		t.Fatalf("legacy sender fallback = %q", got)
	}
}

func TestMatrixReadReceiptTimesSeparateTargetAndReceipt(t *testing.T) {
	targetTime := time.UnixMilli(1784081200000)
	receiptTime := time.UnixMilli(1784081234567)
	readUpTo, readAt := matrixReadReceiptTimes(&bridgev2.MatrixReadReceipt{
		ReadUpTo: targetTime,
		Receipt:  event.ReadReceipt{Timestamp: receiptTime},
	})
	if !readUpTo.Equal(targetTime) {
		t.Fatalf("readUpTo = %v, want %v", readUpTo, targetTime)
	}
	if !readAt.Equal(receiptTime) {
		t.Fatalf("readAt = %v, want %v", readAt, receiptTime)
	}
}

func TestMatrixReadReceiptTimesFallsBackToReceiptForTarget(t *testing.T) {
	receiptTime := time.UnixMilli(1784081234567)
	readUpTo, readAt := matrixReadReceiptTimes(&bridgev2.MatrixReadReceipt{
		Receipt: event.ReadReceipt{Timestamp: receiptTime},
	})
	if !readUpTo.Equal(receiptTime) || !readAt.Equal(receiptTime) {
		t.Fatalf("times = (%v, %v), want receipt time %v", readUpTo, readAt, receiptTime)
	}
}
