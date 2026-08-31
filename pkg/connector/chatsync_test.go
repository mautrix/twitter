package connector

import (
	"context"
	"encoding/base64"
	"testing"

	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func TestXChatItemTrustControlsMessageRequest(t *testing.T) {
	for _, trusted := range []*bool{ptr.Ptr(false), ptr.Ptr(true), nil} {
		client := &TwitterClient{}
		item := &response.XChatInboxItem{
			ConversationDetail:  response.XChatConversationDetail{ConversationID: "g123"},
			LatestMessageEvents: []string{"AgABMA=="},
		}
		if trusted != nil {
			encoded, err := payload.Encode(&payload.MessageEvent{IsTrusted: trusted})
			if err != nil {
				t.Fatal(err)
			}
			item.LatestMessageEvents = []string{base64.StdEncoding.EncodeToString(encoded)}
		}
		conv := client.xchatItemToConversation(t.Context(), item, nil)
		if conv.Trusted != (trusted == nil || *trusted) {
			t.Fatalf("Trusted = %t", conv.Trusted)
		}
		info := client.xchatItemToChatInfo(t.Context(), item, nil, conv)
		if trusted == nil {
			if info.MessageRequest != nil || info.ExtraUpdates != nil {
				t.Fatal("missing trust changed message-request state")
			}
			return
		}
		if info.MessageRequest == nil || *info.MessageRequest == *trusted {
			t.Fatalf("MessageRequest = %v", info.MessageRequest)
		}
		meta := &PortalMetadata{}
		portal := &bridgev2.Portal{Portal: &database.Portal{Metadata: meta}}
		if !info.ExtraUpdates(t.Context(), portal) || meta.XChatTrusted == nil || *meta.XChatTrusted != *trusted {
			t.Fatalf("XChatTrusted = %v", meta.XChatTrusted)
		}
	}
}

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
