package connector

import (
	"context"
	"errors"
	"strings"
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

type xchatCheckpointRecord struct {
	cursor      *payload.XChatCursor
	maxSequence string
	pullVersion *int
}

func xchatTestInt(value int) *int {
	return &value
}

func xchatTestString(value string) *string {
	return &value
}

func TestRunXChatInboxCatchupUsesSequenceAndCheckpointsPages(t *testing.T) {
	var checkpoints []xchatCheckpointRecord
	processCalls := 0
	result, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{
		MaxSequenceID:      "100",
		MessagePullVersion: xchatTestInt(1),
	}, xchatInboxCatchupOps{
		FetchInitial: func(_ context.Context, variables *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			if variables.MaxLocalSequenceId != "100" {
				t.Fatalf("initial max sequence = %q, want 100", variables.MaxLocalSequenceId)
			}
			if variables.MessagePullVersion == nil || *variables.MessagePullVersion != 1 {
				t.Fatalf("initial pull version = %v, want 1", variables.MessagePullVersion)
			}
			return response.XChatInboxPage{
				Items:              make([]response.XChatInboxItem, 2),
				MaxUserSequenceID:  xchatTestString("110"),
				MessagePullVersion: xchatTestInt(2),
				InboxCursor: response.XChatInboxCursor{
					CursorID:        "cursor-1",
					GraphSnapshotID: "graph-1",
				},
			}, nil
		},
		FetchNext: func(_ context.Context, variables *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
			if variables.ContinueCursor == nil || variables.ContinueCursor.CursorId != "cursor-1" {
				t.Fatalf("continuation cursor = %#v, want cursor-1", variables.ContinueCursor)
			}
			return response.XChatInboxPage{
				Items:              make([]response.XChatInboxItem, 3),
				MaxUserSequenceID:  xchatTestString("120"),
				MessagePullVersion: xchatTestInt(3),
				InboxCursor:        response.XChatInboxCursor{PullFinished: true},
			}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			processCalls++
			if processCalls == 1 {
				return xchatInboxPageProcessResult{MaxSequenceID: "115"}, nil
			}
			return xchatInboxPageProcessResult{MaxSequenceID: "125"}, nil
		},
		Checkpoint: func(_ context.Context, cursor *payload.XChatCursor, maxSequence string, pullVersion *int) error {
			checkpoints = append(checkpoints, xchatCheckpointRecord{
				cursor:      cloneXChatCursor(cursor),
				maxSequence: maxSequence,
				pullVersion: cloneXChatInt(pullVersion),
			})
			return nil
		},
	})
	if err != nil {
		t.Fatalf("runXChatInboxCatchup() error = %v", err)
	}
	if result.Pages != 2 || result.Items != 5 {
		t.Fatalf("result pages/items = %d/%d, want 2/5", result.Pages, result.Items)
	}
	if result.MaxSequenceID != "125" {
		t.Fatalf("result max sequence = %q, want 125", result.MaxSequenceID)
	}
	if result.MessagePullVersion == nil || *result.MessagePullVersion != 3 {
		t.Fatalf("result pull version = %v, want 3", result.MessagePullVersion)
	}
	if len(checkpoints) != 2 {
		t.Fatalf("checkpoint count = %d, want 2", len(checkpoints))
	}
	if checkpoints[0].cursor == nil || checkpoints[0].cursor.CursorId != "cursor-1" || checkpoints[0].maxSequence != "115" {
		t.Fatalf("first checkpoint = %#v", checkpoints[0])
	}
	if checkpoints[1].cursor != nil || checkpoints[1].maxSequence != "125" {
		t.Fatalf("final checkpoint = %#v", checkpoints[1])
	}
}

func TestRunXChatInboxCatchupResumesSavedCursor(t *testing.T) {
	initialCalled := false
	result, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{
		MaxSequenceID: "200",
		Cursor: &payload.XChatCursor{
			CursorId:        "resume",
			GraphSnapshotId: "graph",
		},
	}, xchatInboxCatchupOps{
		FetchInitial: func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			initialCalled = true
			return response.XChatInboxPage{}, nil
		},
		FetchNext: func(_ context.Context, variables *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
			if variables.ContinueCursor == nil || variables.ContinueCursor.CursorId != "resume" {
				t.Fatalf("continuation cursor = %#v, want resume", variables.ContinueCursor)
			}
			return response.XChatInboxPage{
				MaxUserSequenceID: xchatTestString("210"),
				InboxCursor:       response.XChatInboxCursor{PullFinished: true},
			}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			return xchatInboxPageProcessResult{}, nil
		},
	})
	if err != nil {
		t.Fatalf("runXChatInboxCatchup() error = %v", err)
	}
	if initialCalled {
		t.Fatal("saved cursor path unexpectedly fetched an initial page")
	}
	if result.MaxSequenceID != "210" || result.Pages != 1 {
		t.Fatalf("result = %#v, want max 210 and one page", result)
	}
}

func TestRunXChatInboxCatchupDoesNotCheckpointFailedPage(t *testing.T) {
	checkpointed := false
	_, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{
		MaxSequenceID: "300",
	}, xchatInboxCatchupOps{
		FetchInitial: func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{MaxUserSequenceID: xchatTestString("310")}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			return xchatInboxPageProcessResult{MaxSequenceID: "309"}, errors.New("queue failed")
		},
		Checkpoint: func(context.Context, *payload.XChatCursor, string, *int) error {
			checkpointed = true
			return nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "queue failed") {
		t.Fatalf("runXChatInboxCatchup() error = %v, want queue failure", err)
	}
	if checkpointed {
		t.Fatal("failed page advanced the reconnect checkpoint")
	}
}

func TestRunXChatInboxCatchupRejectsPageErrors(t *testing.T) {
	processed := false
	_, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{}, xchatInboxCatchupOps{
		FetchInitial: func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{
				Errors:      []map[string]any{{"message": "redacted"}},
				InboxCursor: response.XChatInboxCursor{PullFinished: true},
			}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			processed = true
			return xchatInboxPageProcessResult{}, nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "returned 1 errors") {
		t.Fatalf("runXChatInboxCatchup() error = %v, want page error", err)
	}
	if processed {
		t.Fatal("page with API errors was processed")
	}
}

func TestRunXChatInboxCatchupRejectsStalledCursor(t *testing.T) {
	_, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{
		Cursor: &payload.XChatCursor{CursorId: "stalled", GraphSnapshotId: "graph"},
	}, xchatInboxCatchupOps{
		FetchNext: func(context.Context, *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{
				InboxCursor: response.XChatInboxCursor{
					CursorID:        "stalled",
					GraphSnapshotID: "graph-2",
				},
			}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			return xchatInboxPageProcessResult{}, nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "did not advance") {
		t.Fatalf("runXChatInboxCatchup() error = %v, want stalled cursor error", err)
	}
}

func TestRunXChatInboxCatchupAcceptsTerminalPageWithoutCursor(t *testing.T) {
	continuationCalled := false
	checkpointCalls := 0
	result, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{}, xchatInboxCatchupOps{
		FetchInitial: func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{
				MaxUserSequenceID: xchatTestString("610"),
			}, nil
		},
		FetchNext: func(context.Context, *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
			continuationCalled = true
			return response.XChatInboxPage{}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			return xchatInboxPageProcessResult{MaxSequenceID: "609"}, nil
		},
		Checkpoint: func(_ context.Context, cursor *payload.XChatCursor, maxSequence string, _ *int) error {
			checkpointCalls++
			if cursor != nil {
				t.Fatalf("terminal page checkpoint cursor = %#v, want nil", cursor)
			}
			if maxSequence != "610" {
				t.Fatalf("terminal page checkpoint sequence = %q, want 610", maxSequence)
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("runXChatInboxCatchup() error = %v", err)
	}
	if continuationCalled {
		t.Fatal("terminal page without cursor unexpectedly fetched a continuation")
	}
	if checkpointCalls != 1 {
		t.Fatalf("terminal page checkpoint count = %d, want 1", checkpointCalls)
	}
	if result.Pages != 1 || result.MaxSequenceID != "610" {
		t.Fatalf("terminal page result = %#v, want one page at sequence 610", result)
	}
}

func TestProcessXChatInboxPageRejectsMissingConversationID(t *testing.T) {
	client := &TwitterClient{}
	_, err := client.processXChatInboxPage(context.Background(), response.XChatInboxPage{
		Items: []response.XChatInboxItem{{}},
	}, nil, false)
	if err == nil || !strings.Contains(err.Error(), "no conversation ID") {
		t.Fatalf("processXChatInboxPage() error = %v, want missing conversation ID error", err)
	}
}

func TestRunXChatInboxCatchupDoesNotAdvancePastBlockedGap(t *testing.T) {
	checkpointed := false
	pageCalls := 0
	result, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{
		MaxSequenceID:      "400",
		MessagePullVersion: xchatTestInt(4),
	}, xchatInboxCatchupOps{
		FetchInitial: func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{
				MaxUserSequenceID:  xchatTestString("450"),
				MessagePullVersion: xchatTestInt(5),
				InboxCursor: response.XChatInboxCursor{
					CursorID:        "cursor-2",
					GraphSnapshotID: "graph-2",
				},
			}, nil
		},
		FetchNext: func(context.Context, *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{
				MaxUserSequenceID:  xchatTestString("500"),
				MessagePullVersion: xchatTestInt(6),
				InboxCursor:        response.XChatInboxCursor{PullFinished: true},
			}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			pageCalls++
			maxSequenceID := "449"
			if pageCalls == 2 {
				maxSequenceID = "499"
			}
			return xchatInboxPageProcessResult{
				MaxSequenceID:     maxSequenceID,
				CheckpointBlocked: true,
			}, nil
		},
		Checkpoint: func(context.Context, *payload.XChatCursor, string, *int) error {
			checkpointed = true
			return nil
		},
	})
	if err != nil {
		t.Fatalf("runXChatInboxCatchup() error = %v", err)
	}
	if checkpointed {
		t.Fatal("blocked page advanced the persisted checkpoint")
	}
	if result.Pages != 2 {
		t.Fatalf("processed pages = %d, want 2 despite blocked checkpoint", result.Pages)
	}
	if !result.CheckpointBlocked || result.MaxSequenceID != "400" {
		t.Fatalf("result = %#v, want blocked at sequence 400", result)
	}
	if result.MessagePullVersion == nil || *result.MessagePullVersion != 4 {
		t.Fatalf("result pull version = %v, want 4", result.MessagePullVersion)
	}
}

func TestRunXChatInboxCatchupPropagatesCheckpointFailure(t *testing.T) {
	_, err := runXChatInboxCatchup(context.Background(), xchatInboxCatchupState{}, xchatInboxCatchupOps{
		FetchInitial: func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			return response.XChatInboxPage{InboxCursor: response.XChatInboxCursor{PullFinished: true}}, nil
		},
		ProcessPage: func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			return xchatInboxPageProcessResult{}, nil
		},
		Checkpoint: func(context.Context, *payload.XChatCursor, string, *int) error {
			return errors.New("database unavailable")
		},
	})
	if err == nil || !strings.Contains(err.Error(), "database unavailable") {
		t.Fatalf("runXChatInboxCatchup() error = %v, want checkpoint failure", err)
	}
}
