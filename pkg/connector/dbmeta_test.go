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

func TestPortalMetadataTokenDoesNotEnableInteractiveXChat(t *testing.T) {
	meta := &PortalMetadata{ConversationToken: "token"}
	if meta.CanUseXChat() {
		t.Fatal("CanUseXChat() = true for token-only metadata; interactive actions require a conversation key")
	}
}
