package connector

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

const (
	defaultXChatGapCatchupMessages = 50
	xchatGapCatchupRateLimitPause  = 15 * time.Minute
)

var errXChatGapCatchupPaused = errors.New("xchat gap catch-up paused after a rate limit")

type xchatGapCatchupState struct {
	sync.Mutex
	anchor *database.Message
}

func sendXChatGapBatches(messageCount, batchSize int, send func(start, end int, final bool) error) error {
	if messageCount <= 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = defaultXChatGapCatchupMessages
	}
	for start := 0; start < messageCount; start += batchSize {
		end := min(start+batchSize, messageCount)
		if err := send(start, end, end == messageCount); err != nil {
			return err
		}
	}
	return nil
}

func excludeXChatCurrentAndNewerMessages(messages []*bridgev2.BackfillMessage, currentSequenceID string) []*bridgev2.BackfillMessage {
	if currentSequenceID == "" {
		return messages
	}
	filtered := messages[:0]
	for _, message := range messages {
		if message == nil || compareIntStrings(ParseMessageID(message.ID), currentSequenceID) < 0 {
			filtered = append(filtered, message)
		}
	}
	return filtered
}

// getOrLoadAnchor retains the first pre-gap Matrix message until catch-up
// succeeds. The caller must hold the state's lock.
func (state *xchatGapCatchupState) getOrLoadAnchor(load func() (*database.Message, error)) (*database.Message, error) {
	if state.anchor != nil {
		return state.anchor, nil
	}
	anchor, err := load()
	if err != nil || anchor == nil {
		return anchor, err
	}
	anchorCopy := *anchor
	state.anchor = &anchorCopy
	return state.anchor, nil
}

func (state *xchatGapCatchupState) clearAnchor() {
	state.anchor = nil
}

type xchatGapAnchorQuery interface {
	GetLastNonFakePartAtOrBeforeTimeExcludingSender(ctx context.Context, portal networkid.PortalKey, maxTS time.Time, sender networkid.UserID) (*database.Message, error)
	GetLastNonFakePartAtOrBeforeTime(ctx context.Context, portal networkid.PortalKey, maxTS time.Time) (*database.Message, error)
}

// loadXChatGapAnchor skips the user's own sends, which can be stored after
// incoming messages that the gap is supposed to recover.
func loadXChatGapAnchor(ctx context.Context, query xchatGapAnchorQuery, portal networkid.PortalKey, ownUserID networkid.UserID) (*database.Message, error) {
	maxTS := time.Now().Add(10 * time.Second)
	anchor, err := query.GetLastNonFakePartAtOrBeforeTimeExcludingSender(ctx, portal, maxTS, ownUserID)
	if err != nil || anchor != nil {
		return anchor, err
	}
	return query.GetLastNonFakePartAtOrBeforeTime(ctx, portal, maxTS)
}

func (tc *TwitterClient) catchupXChatConversationGap(
	ctx context.Context,
	conversationID, previousSequenceID, currentSequenceID string,
) error {
	if time.Now().UnixMilli() < tc.xchatGapCatchupPausedUntil.Load() {
		return errXChatGapCatchupPaused
	}
	stateValue, _ := tc.xchatGapCatchupStates.LoadOrStore(conversationID, &xchatGapCatchupState{})
	state := stateValue.(*xchatGapCatchupState)
	state.Lock()
	defer state.Unlock()

	log := zerolog.Ctx(ctx).With().
		Str("component", "xchat_gap_catchup").
		Str("conversation_id", conversationID).
		Str("previous_sequence_id", previousSequenceID).
		Str("current_sequence_id", currentSequenceID).
		Logger()
	portal, err := tc.connector.br.GetExistingPortalByKey(ctx, tc.MakePortalKeyFromID(conversationID))
	if err != nil {
		return err
	}
	if portal == nil || portal.MXID == "" {
		state.clearAnchor()
		log.Debug().Msg("Skipping conversation catch-up because the portal does not exist yet")
		return nil
	}

	latestMessage, err := state.getOrLoadAnchor(func() (*database.Message, error) {
		return loadXChatGapAnchor(ctx, tc.connector.br.DB.Message, portal.PortalKey, MakeUserID(ParseUserLoginID(tc.userLogin.ID)))
	})
	if err != nil {
		return err
	}
	if latestMessage == nil {
		state.clearAnchor()
		log.Debug().Msg("Skipping conversation catch-up because the portal has no message anchor")
		return nil
	}

	count := tc.connector.br.Config.Backfill.MaxCatchupMessages
	if count <= 0 {
		count = defaultXChatGapCatchupMessages
	}
	resp, err := tc.fetchCompleteXChatForwardCatchup(ctx, conversationID, bridgev2.FetchMessagesParams{
		Portal:        portal,
		Forward:       true,
		AnchorMessage: latestMessage,
		Count:         count,
	})
	if errors.Is(err, twittermeow.ErrRatelimitExceeded) {
		tc.xchatGapCatchupPausedUntil.Store(time.Now().Add(xchatGapCatchupRateLimitPause).UnixMilli())
		log.Warn().Dur("pause", xchatGapCatchupRateLimitPause).Msg("X rate limited conversation catch-up, pausing catch-up requests")
		return err
	} else if err != nil {
		return err
	}
	if resp == nil || len(resp.Messages) == 0 {
		if resp != nil && resp.CompleteCallback != nil {
			resp.CompleteCallback()
		}
		state.clearAnchor()
		log.Debug().Msg("No missing conversation messages found")
		return nil
	}

	//lint:ignore SA1019 Gap repair needs synchronous backfill completion, which the public queue API does not provide.
	portalInternal := portal.Internal()
	resp.Messages = portalInternal.CutoffMessages(
		ctx,
		resp.Messages,
		resp.AggressiveDeduplication,
		true,
		latestMessage,
	)
	// The live event that triggered recovery, plus anything that arrived after
	// it, is handled by the websocket after this function returns. Excluding
	// those events preserves live ordering instead of backfilling newer messages
	// ahead of the current event.
	resp.Messages = excludeXChatCurrentAndNewerMessages(resp.Messages, currentSequenceID)
	if len(resp.Messages) == 0 {
		if resp.CompleteCallback != nil {
			resp.CompleteCallback()
		}
		state.clearAnchor()
		log.Debug().Msg("No conversation messages remained after deduplication")
		return nil
	}
	log.Info().Int("message_count", len(resp.Messages)).Msg("Bridging missed XChat conversation messages")
	err = sendXChatGapBatches(len(resp.Messages), count, func(start, end int, final bool) error {
		var completeCallback func()
		if final {
			completeCallback = resp.CompleteCallback
		}
		return portalInternal.SendBackfill(
			ctx,
			tc.userLogin,
			resp.Messages[start:end],
			true,
			final && resp.MarkRead,
			false,
			completeCallback,
		)
	})
	if err != nil {
		return err
	}
	state.clearAnchor()
	return nil
}
