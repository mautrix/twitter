package connector

import (
	"errors"
	"testing"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
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

func TestXChatGapCatchupStateDoesNotRememberLoadFailure(t *testing.T) {
	state := &xchatGapCatchupState{}
	_, err := state.getOrLoadAnchor(func() (*database.Message, error) {
		return nil, errors.New("database unavailable")
	})
	if err == nil || state.anchor != nil {
		t.Fatalf("failed load error/anchor = %v/%#v", err, state.anchor)
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
