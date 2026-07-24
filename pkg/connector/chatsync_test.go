package connector

import (
	"context"
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func TestXChatItemToConversationPreservesPlaintextGroupName(t *testing.T) {
	tc := &TwitterClient{}
	item := &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{
			ConversationID: "g1709621683324379335",
			GroupMetadata: &response.XChatGroupMetadata{
				GroupName:     "Outlaws of CSU",
				UpdatedAtMsec: "1784671554078",
			},
		},
	}

	conv := tc.xchatItemToConversation(context.Background(), item, nil)
	if conv.Name != "Outlaws of CSU" {
		t.Fatalf("conversation name = %q, want %q", conv.Name, "Outlaws of CSU")
	}
}

func TestDecryptGroupNamePreservesPlaintextWithColon(t *testing.T) {
	tc := &TwitterClient{}
	const name = "2026: Outlaws of CSU"
	if got := tc.decryptGroupName(context.Background(), "g1709621683324379335", name); got != name {
		t.Fatalf("decryptGroupName() = %q, want %q", got, name)
	}
}
