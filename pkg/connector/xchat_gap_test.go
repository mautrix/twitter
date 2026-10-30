package connector

import (
	"context"
	"errors"
	"testing"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

func TestXChatGapCatchupStateRetainsAnchorUntilSuccess(t *testing.T) {
	state := &xchatGapCatchupState{}
	original := &database.Message{ID: "original", Timestamp: time.Unix(100, 0)}
	newest := &database.Message{ID: "newest", Timestamp: time.Unix(200, 0)}
	loadCalls := 0

	anchor, err := state.getOrLoadAnchor(func() (*database.Message, error) {
		loadCalls++
		return original, nil
	})
	if err != nil || anchor.ID != original.ID {
		t.Fatalf("first anchor = %#v, %v", anchor, err)
	}
	original.ID = "mutated"
	anchor, err = state.getOrLoadAnchor(func() (*database.Message, error) {
		loadCalls++
		return newest, nil
	})
	if err != nil || anchor.ID != "original" || loadCalls != 1 {
		t.Fatalf("retained anchor/load calls = %#v/%d, want original/1", anchor, loadCalls)
	}

	state.clearAnchor()
	anchor, err = state.getOrLoadAnchor(func() (*database.Message, error) {
		loadCalls++
		return newest, nil
	})
	if err != nil || anchor.ID != newest.ID || loadCalls != 2 {
		t.Fatalf("reloaded anchor/load calls = %#v/%d, want newest/2", anchor, loadCalls)
	}
}

func TestSendXChatGapBatchesChunksAndMarksOnlyFinalBatch(t *testing.T) {
	type call struct {
		start int
		end   int
		final bool
	}
	var calls []call
	err := sendXChatGapBatches(125, 50, func(start, end int, final bool) error {
		calls = append(calls, call{start: start, end: end, final: final})
		return nil
	})
	if err != nil {
		t.Fatalf("sendXChatGapBatches() error = %v", err)
	}
	want := []call{{0, 50, false}, {50, 100, false}, {100, 125, true}}
	if len(calls) != len(want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("calls = %#v, want %#v", calls, want)
		}
	}
}

func TestSendXChatGapBatchesStopsAfterFailure(t *testing.T) {
	calls := 0
	err := sendXChatGapBatches(125, 50, func(int, int, bool) error {
		calls++
		if calls == 2 {
			return errors.New("send failed")
		}
		return nil
	})
	if err == nil || calls != 2 {
		t.Fatalf("sendXChatGapBatches() error/calls = %v/%d, want failure/2", err, calls)
	}
}

func TestExcludeXChatCurrentAndNewerMessages(t *testing.T) {
	messages := []*bridgev2.BackfillMessage{
		{ID: MakeMessageID("100")},
		{ID: MakeMessageID("101")},
		{ID: MakeMessageID("102")},
		{ID: MakeMessageID("103")},
	}
	filtered := excludeXChatCurrentAndNewerMessages(messages, "102")
	if len(filtered) != 2 || filtered[0].ID != MakeMessageID("100") || filtered[1].ID != MakeMessageID("101") {
		t.Fatalf("filtered messages = %#v", filtered)
	}
}

type fakeXChatGapAnchorQuery struct {
	latestFromOthers *database.Message
	oldest           *database.Message
	err              error
	excludedSender   networkid.UserID
	oldestCalls      int
}

func (q *fakeXChatGapAnchorQuery) GetLastNonFakePartAtOrBeforeTimeExcludingSender(_ context.Context, _ networkid.PortalKey, _ time.Time, sender networkid.UserID) (*database.Message, error) {
	q.excludedSender = sender
	return q.latestFromOthers, q.err
}

func (q *fakeXChatGapAnchorQuery) GetLastNonFakePartAtOrBeforeTime(context.Context, networkid.PortalKey, time.Time) (*database.Message, error) {
	q.oldestCalls++
	return q.oldest, nil
}

func TestLoadXChatGapAnchorExcludesOwnMessages(t *testing.T) {
	ownUserID := MakeUserID("1")
	incoming := &database.Message{ID: MakeMessageID("100"), SenderID: MakeUserID("2"), Timestamp: time.Unix(100, 0)}
	oldestOwn := &database.Message{ID: MakeMessageID("50"), SenderID: ownUserID, Timestamp: time.Unix(50, 0)}
	tests := []struct {
		name            string
		query           *fakeXChatGapAnchorQuery
		want            *database.Message
		wantOldestCalls int
	}{
		{name: "latest message from others", query: &fakeXChatGapAnchorQuery{latestFromOthers: incoming, oldest: oldestOwn}, want: incoming},
		{name: "only own messages", query: &fakeXChatGapAnchorQuery{oldest: oldestOwn}, want: oldestOwn, wantOldestCalls: 1},
		{name: "empty portal", query: &fakeXChatGapAnchorQuery{}, wantOldestCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anchor, err := loadXChatGapAnchor(context.Background(), tt.query, networkid.PortalKey{ID: "portal"}, ownUserID)
			if err != nil || anchor != tt.want {
				t.Fatalf("loadXChatGapAnchor() = %#v, %v, want %#v", anchor, err, tt.want)
			}
			if tt.query.excludedSender != ownUserID {
				t.Fatalf("excluded sender = %q, want %q", tt.query.excludedSender, ownUserID)
			}
			if tt.query.oldestCalls != tt.wantOldestCalls {
				t.Fatalf("oldest message lookups = %d, want %d", tt.query.oldestCalls, tt.wantOldestCalls)
			}
		})
	}
}

func TestLoadXChatGapAnchorReturnsQueryError(t *testing.T) {
	query := &fakeXChatGapAnchorQuery{err: errors.New("db failed")}
	anchor, err := loadXChatGapAnchor(context.Background(), query, networkid.PortalKey{ID: "portal"}, MakeUserID("1"))
	if err == nil || anchor != nil || query.oldestCalls != 0 {
		t.Fatalf("loadXChatGapAnchor() = %#v, %v, oldest lookups %d, want error without fallback", anchor, err, query.oldestCalls)
	}
}
