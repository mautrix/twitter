package connector

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func TestOlderTrustedRESTInboxBudgetAndPagination(t *testing.T) {
	for _, test := range []struct {
		name                        string
		limit, synced, wantRequests int
	}{
		{"pending_xchat", 1, 0, 2},
		{"empty_pages", 1, 0, maxOlderTrustedRESTInboxPages},
		{"duplicate_pages", 1, 0, maxOlderTrustedRESTInboxPages},
		{"repeated_cursor", 1, 0, 1},
		{"missing_page", 1, 0, 1},
		{"cancelled", 1, 0, 0},
		{"budget_spent", 1, 1, 0},
		{"zero_first_page_only", 0, 0, 0},
		{"negative_first_page_only", -1, 0, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := twittermeow.NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
			tc := &TwitterClient{client: client, connector: &TwitterConnector{Config: Config{ConversationSyncLimit: test.limit}}}
			initial := &response.TwitterInboxData{Conversations: map[string]*types.Conversation{}}
			for i := 2; i < 22; i++ {
				id := fmt.Sprintf("1-%d", i)
				initial.Conversations[id] = &types.Conversation{ConversationID: id}
			}
			initial.InboxTimelines.Trusted.Status = types.PaginationStatusHasMore
			initial.InboxTimelines.Trusted.MinEntryID = "cursor1"
			client.SetXChatGapHandler(func(context.Context, string, string, string) error { return nil })
			client.GetXChatProcessor().MarkConversationGapUnresolved("1:22")
			requests := 0
			client.HTTP = &http.Client{Transport: connectorRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				requests++
				if requests > maxOlderTrustedRESTInboxPages {
					t.Fatal("startup pagination exceeded its request bound")
				}
				status, cursor := "HAS_MORE", fmt.Sprintf("cursor%d", requests+1)
				conversations := `{}`
				if test.name == "pending_xchat" {
					conversations = `{"pending":{"conversation_id":"1-22","trusted":true,"type":"ONE_TO_ONE"},"request":{"conversation_id":"1-23","trusted":false,"type":"ONE_TO_ONE"}}`
					if requests == 2 {
						status = "AT_END"
					}
				} else if test.name == "duplicate_pages" {
					conversations = `{"alias":{"conversation_id":"1:2","trusted":true,"type":"ONE_TO_ONE"}}`
				} else if test.name == "repeated_cursor" {
					cursor = "cursor1"
				}
				body := fmt.Sprintf(`{"inbox_timeline":{"status":"%s","min_entry_id":"%s","conversations":%s}}`, status, cursor, conversations)
				if test.name == "missing_page" {
					body = `{}`
				}
				return connectorTestHTTPResponse(body), nil
			})}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if test.name == "cancelled" {
				cancel()
			}
			query := payload.DMRequestQuery{}.Default()
			tc.syncOlderTrustedRESTChannels(ctx, initial, &query, test.synced)
			if requests != test.wantRequests {
				t.Fatalf("requests=%d want%d", requests, test.wantRequests)
			}
		})
	}
}

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
		if info.MessageRequest != nil || info.ExtraUpdates == nil {
			t.Fatalf("trust update was not deferred: MessageRequest = %v", info.MessageRequest)
		}
		meta := &PortalMetadata{}
		portal := &bridgev2.Portal{Portal: &database.Portal{Metadata: meta}}
		if !info.ExtraUpdates(t.Context(), portal) || meta.XChatTrusted == nil || *meta.XChatTrusted != *trusted || portal.MessageRequest == *trusted {
			t.Fatalf("XChatTrusted = %v, MessageRequest = %t", meta.XChatTrusted, portal.MessageRequest)
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
