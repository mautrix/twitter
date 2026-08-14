package connector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

type xchatInboxCatchupState struct {
	MaxSequenceID      string
	MessagePullVersion *int
	Cursor             *payload.XChatCursor
}

type xchatInboxCatchupResult struct {
	MaxSequenceID      string
	MessagePullVersion *int
	Pages              int
	Items              int
	CheckpointBlocked  bool
}

type xchatInboxPageProcessResult struct {
	MaxSequenceID     string
	CheckpointBlocked bool
}

type xchatInboxCatchupOps struct {
	FetchInitial func(context.Context, *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error)
	FetchNext    func(context.Context, *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error)
	ProcessPage  func(context.Context, response.XChatInboxPage) (xchatInboxPageProcessResult, error)
	Checkpoint   func(context.Context, *payload.XChatCursor, string, *int) error
}

func cloneXChatInt(value *int) *int {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneXChatCursor(cursor *payload.XChatCursor) *payload.XChatCursor {
	if cursor == nil {
		return nil
	}
	copy := *cursor
	return &copy
}

func maxXChatSequenceID(current, candidate string) string {
	if compareIntStrings(candidate, current) > 0 {
		return candidate
	}
	return current
}

func runXChatInboxCatchup(
	ctx context.Context,
	state xchatInboxCatchupState,
	ops xchatInboxCatchupOps,
) (xchatInboxCatchupResult, error) {
	result := xchatInboxCatchupResult{
		MaxSequenceID:      state.MaxSequenceID,
		MessagePullVersion: cloneXChatInt(state.MessagePullVersion),
	}
	cursor := cloneXChatCursor(state.Cursor)
	pendingMaxSequenceID := result.MaxSequenceID
	pendingMessagePullVersion := cloneXChatInt(result.MessagePullVersion)
	if ops.ProcessPage == nil {
		return result, errors.New("XChat inbox page processor is nil")
	}

	processPage := func(page response.XChatInboxPage) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if len(page.Errors) > 0 {
			return fmt.Errorf("XChat inbox page returned %d errors", len(page.Errors))
		}
		processed, err := ops.ProcessPage(ctx, page)
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		pendingMaxSequenceID = maxXChatSequenceID(pendingMaxSequenceID, processed.MaxSequenceID)
		if page.MaxUserSequenceID != nil {
			pendingMaxSequenceID = maxXChatSequenceID(pendingMaxSequenceID, *page.MaxUserSequenceID)
		}
		if page.MessagePullVersion != nil {
			pendingMessagePullVersion = cloneXChatInt(page.MessagePullVersion)
		}
		result.Pages++
		result.Items += len(page.Items)
		result.CheckpointBlocked = processed.CheckpointBlocked
		return nil
	}

	checkpoint := func() error {
		if !result.CheckpointBlocked {
			result.MaxSequenceID = pendingMaxSequenceID
			result.MessagePullVersion = cloneXChatInt(pendingMessagePullVersion)
		}
		if ops.Checkpoint != nil {
			return ops.Checkpoint(ctx, cursor, result.MaxSequenceID, result.MessagePullVersion)
		}
		return nil
	}

	if cursor == nil {
		if ops.FetchInitial == nil {
			return result, errors.New("initial XChat inbox fetcher is nil")
		}
		variables := payload.NewInitialXChatPageQueryVariables(result.MaxSequenceID)
		if result.MessagePullVersion != nil {
			variables.MessagePullVersion = cloneXChatInt(result.MessagePullVersion)
		}
		page, err := ops.FetchInitial(ctx, variables)
		if err != nil {
			return result, fmt.Errorf("fetch initial XChat catch-up page: %w", err)
		}
		if err = processPage(page); err != nil {
			return result, fmt.Errorf("process initial XChat catch-up page: %w", err)
		}
		cursor, err = validatedNextXChatInboxCursor(page)
		if err != nil {
			return result, fmt.Errorf("read initial XChat catch-up cursor: %w", err)
		}
		if err = checkpoint(); err != nil {
			return result, fmt.Errorf("save initial XChat catch-up checkpoint: %w", err)
		}
	}

	for cursor != nil {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		if ops.FetchNext == nil {
			return result, errors.New("continuation XChat inbox fetcher is nil")
		}
		requestCursor := cloneXChatCursor(cursor)
		page, err := ops.FetchNext(ctx, payload.NewInboxPageRequestQueryVariables(requestCursor))
		if err != nil {
			return result, fmt.Errorf("fetch XChat catch-up page: %w", err)
		}
		if err = processPage(page); err != nil {
			return result, fmt.Errorf("process XChat catch-up page: %w", err)
		}

		nextCursor, err := validatedNextXChatInboxCursor(page)
		if err != nil {
			return result, fmt.Errorf("read XChat catch-up cursor: %w", err)
		}
		if nextCursor != nil && nextCursor.CursorId == requestCursor.CursorId {
			return result, fmt.Errorf("xchat inbox cursor did not advance from %q", requestCursor.CursorId)
		}
		cursor = nextCursor
		if err = checkpoint(); err != nil {
			return result, fmt.Errorf("save XChat catch-up checkpoint: %w", err)
		}
	}

	return result, nil
}

func (tc *TwitterClient) processXChatInboxPage(
	ctx context.Context,
	page response.XChatInboxPage,
	totalItems *atomic.Int32,
	repairTruncatedItems bool,
) ([]string, error) {
	log := zerolog.Ctx(ctx)
	var pageMissing []string
	for i := range page.Items {
		item := &page.Items[i]
		if item.ConversationDetail.ConversationID == "" {
			return nil, fmt.Errorf("XChat inbox item %d has no conversation ID", i)
		}
		if totalItems != nil {
			totalItems.Add(1)
		}
		pageMissing = append(pageMissing, tc.cacheUsersFromItem(item)...)
	}

	if len(pageMissing) > 0 {
		if err := tc.ensureUsersInCacheByID(ctx, pageMissing); err != nil {
			log.Warn().
				Err(err).
				Int("missing_users", len(pageMissing)).
				Msg("Failed to prefetch missing users for inbox page")
		}
	}

	processor := tc.client.GetXChatProcessor()
	g, pageCtx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	for i := range page.Items {
		item := &page.Items[i]
		g.Go(func() error {
			conversationID := item.ConversationDetail.ConversationID
			keyErr := processor.ProcessKeyChangeEvents(pageCtx, item)
			if keyErr != nil {
				log.Warn().
					Err(keyErr).
					Str("conversation_id", conversationID).
					Msg("Failed to process key change events")
				processor.MarkConversationGapUnresolved(conversationID)
				if pageCtx.Err() != nil {
					return pageCtx.Err()
				}
			}

			syncErr := tc.syncXChatChannel(pageCtx, item, nil)
			if syncErr != nil {
				log.Warn().
					Err(syncErr).
					Str("conversation_id", conversationID).
					Msg("Failed to sync XChat channel")
				processor.MarkConversationGapUnresolved(conversationID)
				if pageCtx.Err() != nil {
					return pageCtx.Err()
				}
			}
			// Message handlers can often recover/create the portal themselves. Keep
			// delivering the latest events after a metadata sync failure, while
			// retaining the unresolved marker so the room sync is retried.
			repairConversationGap := keyErr == nil && syncErr == nil && ((repairTruncatedItems && item.HasMore) ||
				processor.ConversationGapUnresolved(conversationID))
			if repairConversationGap {
				if gapErr := tc.catchupXChatConversationGap(pageCtx, conversationID, "", ""); gapErr != nil {
					log.Warn().
						Err(gapErr).
						Str("conversation_id", conversationID).
						Msg("Failed to repair XChat inbox conversation gap; continuing with latest events")
					processor.MarkConversationGapUnresolved(conversationID)
				} else {
					processor.MarkConversationCaughtUp(conversationID)
				}
			}

			messageErr := processor.ProcessMessageAndReadEvents(pageCtx, item)
			if messageErr != nil {
				log.Warn().
					Err(messageErr).
					Str("conversation_id", conversationID).
					Msg("Failed to process message/read events")
				if pageCtx.Err() != nil {
					return pageCtx.Err()
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return pageMissing, err
	}
	return pageMissing, nil
}

func (tc *TwitterClient) syncXChatInboxAfterConnect(
	ctx context.Context,
	getMaxSequenceID func() string,
	setMaxSequenceID func(string),
	getMessagePullVersion func() *int,
	setMessagePullVersion func(*int),
) error {
	tc.xchatInboxSyncLock.Lock()
	defer tc.xchatInboxSyncLock.Unlock()
	if ctx.Err() != nil {
		return ctx.Err()
	}

	log := zerolog.Ctx(ctx).With().Str("component", "xchat_socket_catchup").Logger()
	meta := tc.userLogin.Metadata.(*UserLoginMetadata)
	maxSequenceID := maxXChatSequenceID(getMaxSequenceID(), meta.MaxUserSequenceID)
	messagePullVersion := getMessagePullVersion()
	if meta.MessagePullVersion != nil {
		messagePullVersion = cloneXChatInt(meta.MessagePullVersion)
	}

	processor := tc.client.GetXChatProcessor()
	var stagedMaxSequenceID string
	var stagedMaxLock sync.Mutex
	processor.SetSequenceIDCallback(func(sequenceID string) {
		stagedMaxLock.Lock()
		stagedMaxSequenceID = maxXChatSequenceID(stagedMaxSequenceID, sequenceID)
		stagedMaxLock.Unlock()
	})
	defer processor.SetSequenceIDCallback(setMaxSequenceID)

	var totalItems atomic.Int32
	result, err := runXChatInboxCatchup(ctx, xchatInboxCatchupState{
		MaxSequenceID:      maxSequenceID,
		MessagePullVersion: messagePullVersion,
		Cursor:             xchatInboxCursorFromMetadata(meta.XChatInboxCursor),
	}, xchatInboxCatchupOps{
		FetchInitial: func(ctx context.Context, variables *payload.GetInitialXChatPageQueryVariables) (response.XChatInboxPage, error) {
			log.Info().Str("max_sequence_id", variables.MaxLocalSequenceId).Msg("Fetching XChat reconnect catch-up page")
			resp, err := tc.client.GetInitialXChatPage(ctx, variables)
			if err != nil {
				return response.XChatInboxPage{}, err
			}
			return resp.Data.GetInboxPage, nil
		},
		FetchNext: func(ctx context.Context, variables *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
			resp, err := tc.client.GetInboxPageRequest(ctx, variables)
			if err != nil {
				return response.XChatInboxPage{}, err
			}
			return resp.Data.GetInboxPage, nil
		},
		ProcessPage: func(ctx context.Context, page response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
			_, err := tc.processXChatInboxPage(ctx, page, &totalItems, true)
			stagedMaxLock.Lock()
			observedMax := stagedMaxSequenceID
			stagedMaxLock.Unlock()
			observedMax = maxXChatSequenceID(observedMax, processor.MaxHandledSequenceID())
			return xchatInboxPageProcessResult{
				MaxSequenceID:     observedMax,
				CheckpointBlocked: processor.SequenceCheckpointBlocked(),
			}, err
		},
		Checkpoint: tc.saveXChatInboxCheckpoint,
	})
	if err != nil {
		return err
	}

	setMaxSequenceID(result.MaxSequenceID)
	setMessagePullVersion(result.MessagePullVersion)
	completionLog := log.Info()
	if result.CheckpointBlocked {
		completionLog = log.Warn()
	}
	completionLog.
		Int("pages", result.Pages).
		Int("items", result.Items).
		Str("max_sequence_id", result.MaxSequenceID).
		Bool("checkpoint_blocked", result.CheckpointBlocked).
		Msg("XChat reconnect catch-up completed")
	return nil
}
