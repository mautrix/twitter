package connector

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

type initialChatListContextKey struct{}

func (tc *TwitterClient) deferInitialHistory(ctx context.Context, portal *bridgev2.Portal, info *bridgev2.ChatInfo) *bridgev2.ChatInfo {
	if ctx.Value(initialChatListContextKey{}) != nil || portal.Metadata.(*PortalMetadata).InitialHistoryPending != "" {
		info.CanBackfill = false
		info.ExtraUpdates = bridgev2.MergeExtraUpdaters(info.ExtraUpdates, func(_ context.Context, portal *bridgev2.Portal) bool {
			meta := portal.Metadata.(*PortalMetadata)
			if meta.InitialHistoryPending != "" {
				return false
			}
			meta.InitialHistoryPending = string(tc.userLogin.ID)
			return true
		})
	}
	return info
}

func (tc *TwitterClient) remainingInitialHistory(ctx context.Context, portal *bridgev2.Portal, task *database.BackfillTask) (int, error) {
	limit := tc.connector.br.Config.Backfill.MaxInitialMessages
	pageSize := payload.DefaultGetConversationPageQuerySettings().ConversationEventLimit
	if batchSize := tc.connector.br.Config.Backfill.Queue.BatchSize; batchSize > 0 {
		pageSize = min(pageSize, batchSize)
	}
	if limit < 0 {
		limit = pageSize
	}
	budget := limit
	if task != nil {
		// Count control-only pages to prevent unbounded crawling in sparse chats.
		budget -= max(0, task.BatchCount) * pageSize
	}
	if budget <= 0 {
		return 0, nil
	}
	for count := limit; ; count *= 2 {
		messages, err := tc.connector.br.DB.Message.GetLastNInPortal(ctx, portal.PortalKey, count)
		if err != nil {
			return 0, err
		}
		ids := make(map[networkid.MessageID]struct{}, len(messages))
		for _, msg := range messages {
			if !strings.HasPrefix(string(msg.MXID), "~fake:") {
				ids[msg.ID] = struct{}{}
			}
		}
		if len(ids) >= limit || len(messages) < count {
			return min(budget, max(0, limit-len(ids))), nil
		}
	}
}

func (tc *TwitterClient) backfillInitialHistory(ctx context.Context) {
	tc.initialHistoryLock.Lock()
	defer tc.initialHistoryLock.Unlock()
	for ctx.Err() == nil {
		if err := tc.runInitialHistory(ctx); err == nil {
			zerolog.Ctx(ctx).Info().Msg("Finished importing initial chat history")
			return
		} else if ctx.Err() == nil {
			zerolog.Ctx(ctx).Err(err).Msg("Initial history import interrupted, retrying")
		}
		if !waitForXChatInboxRetry(ctx, time.Minute) {
			return
		}
	}
}

func (tc *TwitterClient) runInitialHistory(ctx context.Context) error {
	portals, err := tc.connector.br.GetAllPortals(ctx)
	if err != nil {
		return err
	}
	var g errgroup.Group
	g.SetLimit(5)
	for _, portal := range portals {
		if portal.MXID == "" || portal.Metadata.(*PortalMetadata).InitialHistoryPending != string(tc.userLogin.ID) {
			continue
		}
		g.Go(func() error { return tc.backfillInitialPortal(ctx, portal) })
	}
	return g.Wait()
}

func (tc *TwitterClient) backfillInitialPortal(ctx context.Context, portal *bridgev2.Portal) error {
	task, err := tc.connector.br.DB.BackfillTask.GetNextForPortal(ctx, portal.PortalKey, true)
	if err != nil {
		return err
	}
	if task == nil {
		task = &database.BackfillTask{PortalKey: portal.PortalKey, UserLoginID: tc.userLogin.ID}
	}
	task.UserLoginID = tc.userLogin.ID
	task.IsDone = false
	// Keep initial tasks out of the automatic queue; scrollback can resume them.
	task.QueueDone = true
	task.NextDispatchMinTS = database.BackfillNextDispatchNever
	if err = tc.connector.br.DB.BackfillTask.Upsert(ctx, task); err != nil {
		return err
	}
	for tc.connector.br.Config.Backfill.Enabled && !task.IsDone {
		remaining, err := tc.remainingInitialHistory(ctx, portal, task)
		if err != nil {
			return err
		} else if remaining <= 0 {
			break
		} else if err = ctx.Err(); err != nil {
			return err
		}
		previousCursor := task.Cursor
		task.FromQueue = true
		tc.connector.br.DoBackfillTask(ctx, task)
		task, err = tc.connector.br.DB.BackfillTask.GetNextForPortal(ctx, portal.PortalKey, true)
		if err != nil {
			return err
		} else if task == nil || task.CompletedAt.IsZero() {
			return fmt.Errorf("initial history batch did not complete")
		}
		if !task.IsDone && task.Cursor == previousCursor {
			nextRemaining, err := tc.remainingInitialHistory(ctx, portal, task)
			if err != nil {
				return err
			} else if nextRemaining >= remaining {
				return fmt.Errorf("initial history cursor and message count did not advance")
			}
		}
	}
	if reads := portal.Metadata.(*PortalMetadata).InitialReadEvents; len(reads) > 0 {
		if err := tc.client.GetXChatProcessor().ProcessMessageAndReadEvents(ctx, &response.XChatInboxItem{
			ConversationDetail:             response.XChatConversationDetail{ConversationID: NormalizeConversationID(ParsePortalID(portal.ID))},
			LatestReadEventsPerParticipant: reads,
		}); err != nil {
			return err
		}
	}
	result := tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
		EventMeta: simplevent.EventMeta{Type: bridgev2.RemoteEventChatInfoChange, PortalKey: portal.PortalKey},
		ChatInfoChange: &bridgev2.ChatInfoChange{ChatInfo: &bridgev2.ChatInfo{
			ExtraUpdates: func(_ context.Context, portal *bridgev2.Portal) bool {
				portal.Metadata.(*PortalMetadata).InitialHistoryPending = ""
				portal.Metadata.(*PortalMetadata).InitialReadEvents = nil
				return true
			},
		}},
	})
	if !xchatRemoteEventHandled(result) {
		return fmt.Errorf("finish initial history: %v", result.Error)
	}
	return portal.Save(ctx)
}
