package connector

import (
	"context"
	"testing"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/bridgeconfig"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
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
