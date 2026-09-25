package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
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
		result.Items += len(page.Items) + len(page.EncodedMessageEvents)
		result.CheckpointBlocked = processed.CheckpointBlocked
		return nil
	}

	checkpoint := func() error {
		if result.CheckpointBlocked {
			// Keep the cursor as well as the sequence behind the gap so a restart
			// retries the unresolved page instead of skipping its conversation.
			return nil
		}
		result.MaxSequenceID = pendingMaxSequenceID
		result.MessagePullVersion = cloneXChatInt(pendingMessagePullVersion)
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
		if nextCursor != nil && ((nextCursor.CursorId != "" && nextCursor.CursorId == requestCursor.CursorId) ||
			(nextCursor.MaxLocalSequenceID != "" && compareIntStrings(nextCursor.MaxLocalSequenceID, requestCursor.MaxLocalSequenceID) <= 0)) {
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
	if page.MessageEventsCursor != nil {
		return nil, tc.client.GetXChatProcessor().ProcessEncodedMessageEvents(ctx, page.EncodedMessageEvents)
	}
	log := zerolog.Ctx(ctx)
	var pageMissing []string
	unavailableCount := 0
	for i := range page.Items {
		item := &page.Items[i]
		if item.ConversationUnavailable {
			unavailableCount++
			continue
		}
		if item.ConversationDetail.ConversationID == "" {
			return nil, fmt.Errorf("XChat inbox item %d has no conversation ID", i)
		}
		if totalItems != nil {
			totalItems.Add(1)
		}
		pageMissing = append(pageMissing, tc.cacheUsersFromItem(item)...)
	}
	if unavailableCount > 0 {
		log.Warn().Int("unavailable_conversations", unavailableCount).Msg("Skipping unavailable XChat conversations")
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
		if item.ConversationUnavailable {
			continue
		}
		g.Go(func() error {
			conversationID := item.ConversationDetail.ConversationID
			keyErr := processor.ProcessKeyChangeEvents(pageCtx, item)
			if keyErr != nil {
				if errors.Is(keyErr, twittermeow.ErrXChatFailurePersistence) {
					return keyErr
				}
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
				if errors.Is(messageErr, twittermeow.ErrXChatFailurePersistence) {
					return messageErr
				}
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
	return pageMissing, errors.Join(g.Wait(), processor.FailedEventPersistenceError())
}

func (tc *TwitterClient) syncXChatInboxAfterConnect(
	ctx context.Context,
	getMaxSequenceID func() string,
	setMaxSequenceID func(string),
	getMessagePullVersion func() *int,
	setMessagePullVersion func(*int),
	drain twittermeow.XChatLiveDrain,
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
	publishedSequence := maxSequenceID
	processor.SetSequenceIDCallback(nil)
	defer func() {
		processor.CapHandledSequenceID(publishedSequence)
		processor.SetSequenceIDCallback(setMaxSequenceID)
	}()

	var totalItems atomic.Int32
	result, err := runXChatPriorityCatchup(ctx, xchatInboxCatchupState{
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
			return xchatInboxPageProcessResult{
				CheckpointBlocked: processor.SequenceCheckpointBlocked(),
			}, err
		},
	}, tc.prepareXChatSnapshotProfiles, processor.ConversationRecoveryPending, drain)
	if err != nil {
		return err
	}

	if !processor.SequenceCheckpointBlocked() {
		if err = tc.saveXChatInboxCheckpoint(ctx, nil, result.MaxSequenceID, result.MessagePullVersion); err != nil {
			return err
		}
		publishedSequence = result.MaxSequenceID
	}
	setMaxSequenceID(publishedSequence)
	if !processor.SequenceCheckpointBlocked() {
		setMessagePullVersion(result.MessagePullVersion)
	}
	result.CheckpointBlocked = processor.SequenceCheckpointBlocked()
	completionLog := log.Info()
	if result.CheckpointBlocked {
		completionLog = log.Warn()
	}
	completionLog.
		Int("pages", result.Pages).
		Int("items", result.Items).
		Str("max_sequence_id", publishedSequence).
		Bool("checkpoint_blocked", result.CheckpointBlocked).
		Msg("XChat reconnect catch-up completed")
	return nil
}

// runXChatPriorityCatchup uses runXChatInboxCatchup to collect inbox pages, then
// processes conversations referenced by queued live messages before dispatching
// those messages. Remaining conversations are processed in batches between drains.
// If the snapshot exceeds the buffer limits or cannot be safely reordered,
// pages are processed in fetch order before draining live messages.
func runXChatPriorityCatchup(ctx context.Context, state xchatInboxCatchupState, ops xchatInboxCatchupOps, prepare func(context.Context, []response.XChatInboxPage) error, recoveryPending func(string) bool, drain twittermeow.XChatLiveDrain) (xchatInboxCatchupResult, error) {
	var pages []response.XChatInboxPage
	bufferedBytes, bufferedItems := 0, 0
	fallback := drain == nil
	collect := ops
	collect.Checkpoint = nil
	seen := map[payload.XChatCursor]bool{}
	collect.FetchNext = func(ctx context.Context, vars *payload.GetInboxPageRequestQueryVariables) (response.XChatInboxPage, error) {
		if !fallback {
			key := *vars.ContinueCursor
			if seen[key] {
				return response.XChatInboxPage{}, fmt.Errorf("XChat inbox cursor cycle")
			}
			seen[key] = true
		}
		return ops.FetchNext(ctx, vars)
	}
	collect.ProcessPage = func(ctx context.Context, page response.XChatInboxPage) (xchatInboxPageProcessResult, error) {
		if !fallback {
			encoded, err := json.Marshal(page)
			if err != nil {
				return xchatInboxPageProcessResult{}, err
			}
			bufferedBytes += len(encoded)
			bufferedItems += len(page.Items)
			fallback = len(pages) >= 128 || bufferedBytes > 16<<20 || bufferedItems > 4096 || page.MessageEventsCursor != nil
			for i := range page.Items {
				fallback = fallback || !twittermeow.XChatInboxItemIsConversationScoped(&page.Items[i])
			}
		}
		if fallback {
			seen = nil
			for _, buffered := range pages {
				if _, err := ops.ProcessPage(ctx, buffered); err != nil {
					return xchatInboxPageProcessResult{}, err
				}
			}
			pages = nil
			processed, err := ops.ProcessPage(ctx, page)
			// Ignore sequence progress from on-demand history.
			processed.MaxSequenceID = ""
			if page.MessageEventsCursor != nil {
				processed.MaxSequenceID = page.MessageEventsCursor.MaxLocalSequenceID
				for _, encoded := range page.EncodedMessageEvents {
					if event, decodeErr := twittermeow.DecodeMessageEvent(encoded); decodeErr == nil && event.SequenceId != nil {
						processed.MaxSequenceID = maxXChatSequenceID(processed.MaxSequenceID, *event.SequenceId)
					}
				}
			}
			return processed, err
		}
		pages = append(pages, page)
		return xchatInboxPageProcessResult{}, nil
	}
	result, err := runXChatInboxCatchup(ctx, state, collect)
	if err != nil {
		return result, err
	}
	if !fallback {
		if err = prepare(ctx, pages); err != nil {
			return result, err
		}
		pending := map[string][]response.XChatInboxItem{}
		known := map[string]bool{}
		order := []string{}
		for _, page := range pages {
			for _, item := range page.Items {
				id := item.ConversationDetail.ConversationID
				known[id] = true
				if _, ok := pending[id]; !ok {
					order = append(order, id)
				}
				pending[id] = append(pending[id], item)
			}
		}
		apply := func(ctx context.Context, items []response.XChatInboxItem) error {
			for _, item := range items {
				if _, err := ops.ProcessPage(ctx, response.XChatInboxPage{Items: []response.XChatInboxItem{item}}); err != nil {
					return err
				}
			}
			return nil
		}
		applyIDs := func(ids []string) error {
			for _, id := range ids {
				if err := apply(ctx, pending[id]); err != nil {
					return err
				}
				delete(pending, id)
			}
			return nil
		}
		before := func(message *payload.Message) error {
			ids := twittermeow.XChatMessageConversations(message, "")
			if ids == nil {
				return applyIDs(order)
			}
			if err := applyIDs(ids); err != nil {
				return err
			}
			for _, id := range ids {
				if !known[id] || recoveryPending(id) {
					return applyIDs(order)
				}
			}
			return nil
		}
		for offset := 0; offset < len(order); offset += 10 {
			if err = drain(before); err != nil {
				return result, err
			}
			group, groupCtx := errgroup.WithContext(ctx)
			for _, id := range order[offset:min(offset+10, len(order))] {
				items := pending[id]
				delete(pending, id)
				group.Go(func() error { return apply(groupCtx, items) })
			}
			if err = group.Wait(); err != nil {
				return result, err
			}
		}
	}
	if drain != nil {
		err = drain(nil)
	}
	return result, err
}

// Deferred rooms must use the newest profiles from the snapshot.
func (tc *TwitterClient) prepareXChatSnapshotProfiles(ctx context.Context, pages []response.XChatInboxPage) error {
	var missing []string
	strip := func(input []response.XChatUserResult) []response.XChatUserResult {
		users := append([]response.XChatUserResult(nil), input...)
		for i := range users {
			if id, _ := xchatUserFromResult(users[i]); id != "" {
				users[i].RestID = id
				users[i].Result = nil
			}
		}
		return users
	}
	for p := range pages {
		if err := ctx.Err(); err != nil {
			return err
		}
		pages[p].Items = append([]response.XChatInboxItem(nil), pages[p].Items...)
		for i := range pages[p].Items {
			missing = append(missing, tc.cacheUsersFromItem(&pages[p].Items[i])...)
			detail := &pages[p].Items[i].ConversationDetail
			detail.ParticipantsResults = strip(detail.ParticipantsResults)
			detail.GroupMembersResults = strip(detail.GroupMembersResults)
			detail.GroupAdminsResults = strip(detail.GroupAdminsResults)
		}
	}
	if err := tc.ensureUsersInCacheByID(ctx, missing); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Failed to prefetch snapshot users")
	}
	return ctx.Err()
}
