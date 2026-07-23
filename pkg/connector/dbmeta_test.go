package connector

import "testing"

func TestPortalMetadataCanBackfillXChat(t *testing.T) {
	tests := []struct {
		name string
		meta *PortalMetadata
		want bool
	}{
		{name: "nil metadata", meta: nil, want: false},
		{name: "legacy conversation", meta: &PortalMetadata{}, want: false},
		{name: "stored conversation key", meta: &PortalMetadata{
			ConversationKeys: map[string]*ConversationKeyData{"123": {}},
		}, want: true},
		{name: "token can bootstrap missing key", meta: &PortalMetadata{
			ConversationToken: "token",
		}, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.meta.CanBackfillXChat(); got != test.want {
				t.Fatalf("CanBackfillXChat() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestPortalMetadataTokenDoesNotEnableEncryptedXChat(t *testing.T) {
	meta := &PortalMetadata{ConversationToken: "token"}
	if meta.CanUseXChat() {
		t.Fatal("CanUseXChat() = true for token-only metadata; encrypted actions require a conversation key")
	}
	if !meta.IsXChatConversation() {
		t.Fatal("IsXChatConversation() = false for token-only metadata")
	}
}

func TestSelectMessageSendMode(t *testing.T) {
	tests := []struct {
		name string
		meta *PortalMetadata
		want messageSendMode
	}{
		{name: "legacy REST conversation", meta: &PortalMetadata{}, want: messageSendREST},
		{name: "plaintext XChat conversation", meta: &PortalMetadata{
			ConversationToken: "token",
		}, want: messageSendXChatPlaintext},
		{name: "encrypted XChat conversation", meta: &PortalMetadata{
			ConversationToken: "token",
			ConversationKeys:  map[string]*ConversationKeyData{"123": {}},
		}, want: messageSendXChatEncrypted},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := selectMessageSendMode(test.meta); got != test.want {
				t.Fatalf("selectMessageSendMode() = %d, want %d", got, test.want)
			}
		})
	}
}
