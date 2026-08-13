package connector

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/bridgeconfig"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

func TestGetBackfillMaxBatchCount(t *testing.T) {
	tests := []struct {
		name       string
		portalID   networkid.PortalID
		overrides  map[string]int
		maxBatches int
		want       int
	}{
		{
			name:       "direct message override",
			portalID:   "123-456",
			overrides:  map[string]int{"dm": -1, "group_dm": 4},
			maxBatches: 0,
			want:       -1,
		},
		{
			name:       "group message override",
			portalID:   "g123",
			overrides:  map[string]int{"dm": -1, "group_dm": 4},
			maxBatches: 0,
			want:       4,
		},
		{
			name:       "fallback to global limit",
			portalID:   "123-456",
			overrides:  map[string]int{"group_dm": 4},
			maxBatches: 2,
			want:       2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &TwitterClient{}
			portal := &bridgev2.Portal{
				Portal: &database.Portal{PortalKey: networkid.PortalKey{ID: test.portalID}},
				Bridge: &bridgev2.Bridge{Config: &bridgeconfig.BridgeConfig{
					Backfill: bridgeconfig.BackfillConfig{Queue: bridgeconfig.BackfillQueueConfig{
						MaxBatches:         test.maxBatches,
						MaxBatchesOverride: test.overrides,
					}},
				}},
			}
			if got := client.GetBackfillMaxBatchCount(context.Background(), portal, nil); got != test.want {
				t.Fatalf("GetBackfillMaxBatchCount() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestXChatBackfillEncryptionDetectionUsesKeyVersion(t *testing.T) {
	if !isEncryptedXChatBackfillMessage(&payload.MessageCreateEvent{
		ConversationKeyVersion: ptr.Ptr("7"),
	}) {
		t.Fatal("message with a conversation key version was treated as plaintext")
	}
	if isEncryptedXChatBackfillMessage(&payload.MessageCreateEvent{}) {
		t.Fatal("message without a conversation key version was treated as encrypted")
	}
}

func TestXChatBackfillEmptyConversationMarkerIsNotCritical(t *testing.T) {
	if isCriticalXChatBackfillSkip(xchatSkipReasonEmptyContents) {
		t.Fatal("empty conversation-create marker was treated as a missing message")
	}
}

func TestNormalizeXChatBackfillEventConversation(t *testing.T) {
	evt := &payload.MessageEvent{}
	if !normalizeXChatBackfillEventConversation(evt, "conversation") || ptr.Val(evt.ConversationId) != "conversation" {
		t.Fatalf("missing conversation ID was not filled: %#v", evt.ConversationId)
	}
	evt.ConversationId = ptr.Ptr("different")
	if normalizeXChatBackfillEventConversation(evt, "conversation") {
		t.Fatal("mismatched backfill event conversation was accepted")
	}
}

func TestCompleteXChatForwardCatchupIsNotLimitedToFiftyMessages(t *testing.T) {
	anchorTimestamp := time.Unix(125, 0)
	pageCursors := []string{xchatBackfillMaxInt, "176", "151", "126"}
	pageRanges := [][2]int{{176, 200}, {151, 175}, {126, 150}, {100, 125}}
	fetchCalls := 0

	result, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "message-125", Timestamp: anchorTimestamp},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(_ context.Context, cursor string, pageSize int) (*parsedXChatPage, error) {
			if pageSize != 50 {
				t.Fatalf("page size = %d, want 50", pageSize)
			}
			if fetchCalls >= len(pageCursors) || cursor != pageCursors[fetchCalls] {
				t.Fatalf("fetch %d cursor = %q", fetchCalls, cursor)
			}
			bounds := pageRanges[fetchCalls]
			messages := make([]*bridgev2.BackfillMessage, 0, bounds[1]-bounds[0]+1)
			for timestamp := bounds[1]; timestamp >= bounds[0]; timestamp-- {
				messages = append(messages, &bridgev2.BackfillMessage{
					ID:        networkid.MessageID(fmt.Sprintf("message-%d", timestamp)),
					Timestamp: time.Unix(int64(timestamp), 0),
				})
			}
			page := &parsedXChatPage{
				messages:      messages,
				oldestEventTS: time.Unix(int64(bounds[0]), 0),
				pageHasMore:   fetchCalls < len(pageCursors)-1,
			}
			fetchCalls++
			if page.pageHasMore {
				page.nextCursor = pageCursors[fetchCalls]
			}
			return page, nil
		},
	)
	if err != nil {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v", err)
	}
	if fetchCalls != 4 {
		t.Fatalf("fetch calls = %d, want 4", fetchCalls)
	}
	if len(result.Messages) != 75 {
		t.Fatalf("message count = %d, want 75", len(result.Messages))
	}
	if !result.Messages[0].Timestamp.Equal(time.Unix(126, 0)) ||
		!result.Messages[len(result.Messages)-1].Timestamp.Equal(time.Unix(200, 0)) {
		t.Fatalf("message range = %s..%s, want 126..200", result.Messages[0].Timestamp, result.Messages[len(result.Messages)-1].Timestamp)
	}
}

func TestCompleteXChatForwardCatchupRejectsStalledCursor(t *testing.T) {
	_, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{
				nextCursor:    xchatBackfillMaxInt,
				oldestEventTS: time.Unix(150, 0),
				pageHasMore:   true,
			}, nil
		},
	)
	if err == nil || !strings.Contains(err.Error(), "did not advance") {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v, want stalled cursor error", err)
	}
}

func TestCompleteXChatForwardCatchupKeepsDifferentMessageAtAnchorTimestamp(t *testing.T) {
	result, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{
				messages: []*bridgev2.BackfillMessage{{
					ID:        "same-timestamp-peer",
					Timestamp: time.Unix(100, 0),
				}},
				oldestEventTS: time.Unix(99, 0),
			}, nil
		},
	)
	if err != nil {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v", err)
	}
	if len(result.Messages) != 1 || result.Messages[0].ID != "same-timestamp-peer" {
		t.Fatalf("messages = %#v, want same-timestamp peer", result.Messages)
	}
	if !result.AggressiveDeduplication {
		t.Fatal("automatic catch-up did not enable database deduplication")
	}
}

func TestCompleteXChatForwardCatchupRejectsTimestampLessAnchor(t *testing.T) {
	_, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor"},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			t.Fatal("fetcher called with unusable anchor")
			return nil, nil
		},
	)
	if err == nil || !strings.Contains(err.Error(), "anchor has no timestamp") {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v, want anchor timestamp error", err)
	}
}

func TestCompleteXChatForwardCatchupRejectsUnconvertedNewerMessage(t *testing.T) {
	_, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{
				oldestEventTS: time.Unix(90, 0),
				unconvertedMessages: []xchatUnconvertedMessage{{
					timestamp: time.Unix(110, 0),
				}},
			}, nil
		},
	)
	if err == nil || !strings.Contains(err.Error(), "unconverted newer messages") {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v, want incomplete page error", err)
	}
}

func TestCompleteXChatForwardCatchupRejectsUndecodableEvent(t *testing.T) {
	_, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{decodeFailedCount: 1}, nil
		},
	)
	if err == nil || !strings.Contains(err.Error(), "decode failures") {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v, want incomplete page error", err)
	}
}

func TestCompleteXChatForwardCatchupRejectsKeyStoreFailure(t *testing.T) {
	_, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{keyStoreFailedCount: 1}, nil
		},
	)
	if err == nil || !strings.Contains(err.Error(), "key-store failures") {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v, want key-store error", err)
	}
}

func TestCompleteXChatForwardCatchupRejectsMessageWithoutTimestamp(t *testing.T) {
	_, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{
				messages: []*bridgev2.BackfillMessage{{ID: "missing-timestamp"}},
			}, nil
		},
	)
	if err == nil || !strings.Contains(err.Error(), "messages without timestamps") {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v, want timestamp error", err)
	}
}

func TestCompleteXChatForwardCatchupIgnoresUnconvertedOlderMessage(t *testing.T) {
	result, err := fetchXChatForwardCatchupPages(
		context.Background(),
		"conversation",
		&database.Message{ID: "anchor", Timestamp: time.Unix(100, 0)},
		xchatForwardCatchupOptions{PageSize: 50, RequireComplete: true},
		func(context.Context, string, int) (*parsedXChatPage, error) {
			return &parsedXChatPage{
				oldestEventTS: time.Unix(90, 0),
				unconvertedMessages: []xchatUnconvertedMessage{{
					timestamp: time.Unix(90, 0),
				}},
			}, nil
		},
	)
	if err != nil {
		t.Fatalf("fetchXChatForwardCatchupPages() error = %v", err)
	}
	if len(result.Messages) != 0 {
		t.Fatalf("message count = %d, want 0", len(result.Messages))
	}
}

func TestPrepareRESTFallbackPreservesAnchorForFiltering(t *testing.T) {
	anchor := &database.Message{ID: "anchor"}
	params := bridgev2.FetchMessagesParams{
		Cursor:        networkid.PaginationCursor(restFallbackCursorPrefix + "xchat-cursor"),
		AnchorMessage: anchor,
	}

	restParams, opts := prepareRESTBackfillFetchParams(params, restBackfillCursorModeFallback)
	if restParams.Cursor != "" {
		t.Fatalf("REST fallback cursor = %q, want empty", restParams.Cursor)
	}
	if !opts.IgnoreAnchorForQuery {
		t.Fatal("REST fallback did not ignore the XChat anchor for the REST query")
	}
	if restParams.AnchorMessage != anchor {
		t.Fatal("REST fallback removed the anchor needed for response filtering")
	}
}

func TestCompareIntStringsSupportsLargeDecimalValues(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{a: "18446744073709551615", b: "9223372036854775807", want: 1},
		{a: "9223372036854775808", b: "18446744073709551615", want: -1},
		{a: "00042", b: "42", want: 0},
		{a: "invalid", b: "42", want: -1},
	}
	for _, test := range tests {
		if got := compareIntStrings(test.a, test.b); got != test.want {
			t.Errorf("compareIntStrings(%q, %q) = %d, want %d", test.a, test.b, got, test.want)
		}
	}
}
