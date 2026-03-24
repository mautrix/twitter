package connector

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

var _ bridgev2.BackfillingNetworkAPI = (*TwitterClient)(nil)

const xchatBackfillMaxInt = "9223372036854775807"
const restFallbackCursorPrefix = "rest_fallback:"
const restPaginationCursorPrefix = "rest_cursor:"

type restBackfillOptions struct {
	// IgnoreAnchorForQuery skips using AnchorMessage.ID as max_id/min_id.
	// This is needed when transitioning from XChat to REST where the anchor ID
	// may be an XChat sequence ID that REST can't paginate with.
	IgnoreAnchorForQuery bool
	// IgnoreAnchorForFiltering skips message-side anchor filtering after a REST
	// page is fetched. This is needed in cross-protocol fallback, where the
	// anchor originates from XChat and doesn't map cleanly to REST history.
	IgnoreAnchorForFiltering bool
}

type restBackfillCursorMode int

const (
	restBackfillCursorModeRaw restBackfillCursorMode = iota
	restBackfillCursorModeFallback
	restBackfillCursorModePagination
)

func parseRESTBackfillCursorMode(cursor networkid.PaginationCursor) restBackfillCursorMode {
	cursorStr := string(cursor)
	switch {
	case strings.HasPrefix(cursorStr, restFallbackCursorPrefix):
		return restBackfillCursorModeFallback
	case strings.HasPrefix(cursorStr, restPaginationCursorPrefix):
		return restBackfillCursorModePagination
	default:
		return restBackfillCursorModeRaw
	}
}

func prepareRESTBackfillFetchParams(
	fetchParams bridgev2.FetchMessagesParams,
	mode restBackfillCursorMode,
	crossProtocol bool,
) (bridgev2.FetchMessagesParams, restBackfillOptions) {
	restParams := fetchParams
	opts := restBackfillOptions{}
	switch mode {
	case restBackfillCursorModeFallback:
		restParams.Cursor = ""
		opts.IgnoreAnchorForQuery = true
		opts.IgnoreAnchorForFiltering = true
	case restBackfillCursorModePagination:
		restParams.Cursor = networkid.PaginationCursor(strings.TrimPrefix(string(fetchParams.Cursor), restPaginationCursorPrefix))
		if crossProtocol {
			opts.IgnoreAnchorForFiltering = true
		}
	}
	return restParams, opts
}

func (tc *TwitterClient) fetchRESTMessagesForCursorMode(
	ctx context.Context,
	conversationID string,
	fetchParams bridgev2.FetchMessagesParams,
	mode restBackfillCursorMode,
	crossProtocol bool,
) (*bridgev2.FetchMessagesResponse, error) {
	restParams, opts := prepareRESTBackfillFetchParams(fetchParams, mode, crossProtocol)
	restResp, err := tc.fetchRESTMessagesWithOptions(ctx, conversationID, restParams, opts)
	if err != nil {
		return nil, err
	}
	return ensureRESTFallbackCursor(restResp), nil
}

func (tc *TwitterClient) FetchMessages(ctx context.Context, fetchParams bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	conversationID := ParsePortalID(fetchParams.Portal.PortalKey.ID)
	meta := fetchParams.Portal.Metadata.(*PortalMetadata)
	cursorMode := parseRESTBackfillCursorMode(fetchParams.Cursor)

	// Use REST API backfill for conversations without encryption keys.
	if !meta.CanUseXChat() {
		return tc.fetchRESTMessagesForCursorMode(ctx, conversationID, fetchParams, cursorMode, false)
	}

	// Prefixed cursors indicate an in-progress XChat -> REST transition.
	if cursorMode != restBackfillCursorModeRaw {
		return tc.fetchRESTMessagesForCursorMode(ctx, conversationID, fetchParams, cursorMode, true)
	}

	if fetchParams.Forward {
		// Forward backfill for XChat: initial messages are already processed during room creation
		// via ProcessMessageAndReadEvents(). No additional forward fetch needed.
		// FIXME the above is likely false and will fail to fill gaps correctly
		return &bridgev2.FetchMessagesResponse{
			Forward: true,
			HasMore: false,
		}, nil
	}

	// Priority: Cursor > AnchorMessage > xchatBackfillMaxInt
	var minSeqID string
	var cursorSource string
	if fetchParams.Cursor != "" {
		minSeqID = string(fetchParams.Cursor)
		cursorSource = "cursor"
	} else if fetchParams.AnchorMessage != nil {
		minSeqID = ParseMessageID(fetchParams.AnchorMessage.ID)
		cursorSource = "anchor_message"
	} else {
		// No cursor or anchor - fetch from the beginning (max int = oldest first)
		minSeqID = xchatBackfillMaxInt
		cursorSource = "max_int"
	}

	zerolog.Ctx(ctx).Info().
		Str("conversation_id", conversationID).
		Bool("forward", fetchParams.Forward).
		Str("cursor", minSeqID).
		Str("cursor_source", cursorSource).
		Int("count", fetchParams.Count).
		Msg("XChat backfill requested")

	if minSeqID == "" {
		return &bridgev2.FetchMessagesResponse{HasMore: false}, nil
	}
	if _, err := strconv.ParseUint(minSeqID, 10, 64); err != nil {
		return nil, fmt.Errorf("oldest message id is not numeric (cannot use for XChat backfill): %q", minSeqID)
	}

	count := fetchParams.Count
	if count <= 0 {
		count = 50
	}

	settings := payload.DefaultGetConversationPageQuerySettings()
	if count < settings.ConversationEventLimit {
		settings.ConversationEventLimit = count
	}

	minKeyVersion := xchatBackfillMaxInt
	if km := tc.client.GetKeyManager(); km != nil {
		if ck, err := km.GetLatestConversationKey(ctx, conversationID); err == nil && ck != nil && ck.KeyVersion != "" {
			minKeyVersion = ck.KeyVersion
		}
	}

	vars := payload.NewGetConversationPageQueryVariables(conversationID, minSeqID, minKeyVersion, settings)
	resp, err := tc.client.GetConversationPage(ctx, vars)
	if err != nil {
		return nil, err
	}
	if len(resp.Errors) > 0 && resp.Errors[0].Message != "" {
		return nil, fmt.Errorf("GetConversationPageQuery error: %s", resp.Errors[0].Message)
	}

	page := resp.Data.GetConversationPage

	// Process missing key change events first to maximize decryption success.
	for _, enc := range page.MissingConversationKeyChangeEvents {
		evt, err := twittermeow.DecodeMessageEvent(enc)
		if err != nil || evt == nil || evt.Detail == nil || evt.Detail.ConversationKeyChangeEvent == nil {
			continue
		}
		_ = tc.storeConversationKeyFromChangeEvent(ctx, evt, evt.Detail.ConversationKeyChangeEvent)
	}

	msgs := make([]*bridgev2.BackfillMessage, 0, len(page.EncodedMessageEvents))
	nextCursor := ""

	for _, enc := range page.EncodedMessageEvents {
		evt, err := twittermeow.DecodeMessageEvent(enc)
		if err != nil || evt == nil || evt.Detail == nil {
			continue
		}

		if evt.Detail.ConversationKeyChangeEvent != nil {
			_ = tc.storeConversationKeyFromChangeEvent(ctx, evt, evt.Detail.ConversationKeyChangeEvent)
			continue
		}
		if evt.Detail.MessageCreateEvent == nil {
			continue
		}

		msg, ts, streamOrder := tc.decodeXChatMessageCreateForBackfill(ctx, conversationID, evt)
		if msg == nil {
			continue
		}

		converted := tc.convertToMatrix(ctx, fetchParams.Portal, tc.connector.br.Bot, &msg.MessageData)
		if converted == nil || len(converted.Parts) == 0 {
			continue
		}

		msgID := msg.SequenceID
		if msgID == "" {
			msgID = msg.ID
		}

		if ts.IsZero() && msg.Time != "" {
			ts = methods.ParseMsecTimestamp(msg.Time)
		}
		if ts.IsZero() && msgID != "" {
			ts = methods.ParseSnowflake(msgID)
		}

		if streamOrder == 0 && msgID != "" {
			streamOrder = methods.ParseInt64(msgID)
		}

		msgs = append(msgs, &bridgev2.BackfillMessage{
			ConvertedMessage: converted,
			Sender:           tc.MakeEventSender(msg.MessageData.SenderID),
			ID:               MakeMessageID(msgID),
			TxnID:            networkid.TransactionID(msgID),
			Timestamp:        ts,
			StreamOrder:      streamOrder,
		})

		if seq := ptr.Val(evt.SequenceId); seq != "" {
			if nextCursor == "" || compareIntStrings(seq, nextCursor) < 0 {
				nextCursor = seq
			}
		} else if msgID != "" {
			if nextCursor == "" || compareIntStrings(msgID, nextCursor) < 0 {
				nextCursor = msgID
			}
		}
	}

	sortBackfillMessages(msgs)

	hasMore := page.HasMore
	if nextCursor == "" {
		hasMore = false
	}

	xchatResp := &bridgev2.FetchMessagesResponse{
		Messages: msgs,
		Cursor:   networkid.PaginationCursor(nextCursor),
		HasMore:  hasMore,
	}
	if hasMore {
		return xchatResp, nil
	}

	// If XChat is exhausted but returned a final page, schedule one more batch at the
	// terminal cursor so the next call can attempt REST fallback for legacy history.
	// Guard against loops by not reusing the same cursor repeatedly.
	if len(msgs) > 0 && nextCursor != "" && string(fetchParams.Cursor) != nextCursor {
		zerolog.Ctx(ctx).Debug().
			Str("conversation_id", conversationID).
			Str("xchat_terminal_cursor", nextCursor).
			Int("xchat_message_count", len(msgs)).
			Msg("XChat backfill exhausted, scheduling REST fallback on next batch")
		xchatResp.Cursor = networkid.PaginationCursor(restFallbackCursorPrefix + nextCursor)
		xchatResp.HasMore = true
		return xchatResp, nil
	}

	zerolog.Ctx(ctx).Debug().
		Str("conversation_id", conversationID).
		Str("cursor", string(fetchParams.Cursor)).
		Bool("had_xchat_messages", len(msgs) > 0).
		Msg("XChat backfill exhausted, attempting REST fallback")
	restResp, err := tc.fetchRESTMessagesForCursorMode(ctx, conversationID, fetchParams, restBackfillCursorModeFallback, true)
	if err != nil {
		return nil, err
	}
	return restResp, nil
}

func ensureRESTFallbackCursor(resp *bridgev2.FetchMessagesResponse) *bridgev2.FetchMessagesResponse {
	if resp == nil || !resp.HasMore || resp.Cursor == "" || strings.HasPrefix(string(resp.Cursor), restPaginationCursorPrefix) {
		return resp
	}
	resp.Cursor = networkid.PaginationCursor(restPaginationCursorPrefix + string(resp.Cursor))
	return resp
}

func (tc *TwitterClient) decodeXChatMessageCreateForBackfill(ctx context.Context, conversationID string, evt *payload.MessageEvent) (*types.Message, time.Time, int64) {
	mce := evt.Detail.MessageCreateEvent
	keyVersion := ptr.Val(mce.ConversationKeyVersion)
	contentsBytes := mce.Contents

	if len(contentsBytes) == 0 {
		return nil, time.Time{}, 0
	}

	var entry *payload.MessageEntryContents

	if evt.MessageEventSignature != nil && evt.MessageEventSignature.Signature != nil {
		km := tc.client.GetKeyManager()
		convKey, err := km.GetConversationKey(ctx, conversationID, keyVersion)
		if err != nil || convKey == nil || len(convKey.Key) == 0 {
			return nil, time.Time{}, 0
		}

		debugLog := zerolog.Ctx(ctx).With().
			Str("component", "xchat_backfill_decrypt").
			Str("conversation_id", conversationID).
			Str("key_version", keyVersion).
			Logger()

		decrypted, err := crypto.DecryptMessageEntryContentsBytesDebug(contentsBytes, convKey.Key, &debugLog)
		if err != nil {
			return nil, time.Time{}, 0
		}
		entry = decrypted
	} else {
		parsed, err := crypto.ParseMessageEntryContentsBytes(contentsBytes)
		if err != nil {
			return nil, time.Time{}, 0
		}
		entry = parsed
	}
	if entry == nil || entry.Message == nil {
		return nil, time.Time{}, 0
	}

	if entry.Message.MessageText == nil && len(entry.Message.Attachments) == 0 {
		return nil, time.Time{}, 0
	}

	msg := twittermeow.ConvertXChatMessageContentsToMessage(evt, entry.Message, keyVersion)
	if msg == nil {
		return nil, time.Time{}, 0
	}

	msgID := msg.SequenceID
	if msgID == "" {
		msgID = msg.ID
	}

	ts := methods.ParseMsecTimestamp(msg.Time)
	streamOrder := methods.ParseInt64(msgID)
	return msg, ts, streamOrder
}

func (tc *TwitterClient) storeConversationKeyFromChangeEvent(ctx context.Context, evt *payload.MessageEvent, ckce *payload.ConversationKeyChangeEvent) error {
	conversationID := ptr.Val(evt.ConversationId)
	newKeyVersion := ptr.Val(ckce.ConversationKeyVersion)

	signingKey, err := tc.client.GetKeyManager().GetOwnSigningKey(ctx)
	if err != nil {
		return err
	}

	ownUserID := tc.client.GetCurrentUserID()
	if ownUserID == "" {
		return nil
	}

	var ourEncryptedKey string
	for _, pk := range ckce.ConversationParticipantKeys {
		if ptr.Val(pk.UserId) == ownUserID {
			ourEncryptedKey = ptr.Val(pk.EncryptedConversationKey)
			break
		}
	}
	if ourEncryptedKey == "" {
		return nil
	}

	convKeyBytes, err := crypto.UnwrapConversationKey(ourEncryptedKey, signingKey.DecryptKeyB64)
	if err != nil {
		return err
	}

	keyCreatedAt := methods.ParseMsecTimestamp(ptr.Val(evt.CreatedAtMsec))
	if keyCreatedAt.IsZero() {
		return fmt.Errorf("missing valid XChat key timestamp for conversation %s key %s", conversationID, newKeyVersion)
	}

	return tc.client.GetKeyManager().PutConversationKey(ctx, &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     newKeyVersion,
		Key:            convKeyBytes,
		CreatedAt:      keyCreatedAt,
	})
}

func (tc *TwitterClient) fetchRESTMessagesWithOptions(
	ctx context.Context,
	conversationID string,
	fetchParams bridgev2.FetchMessagesParams,
	opts restBackfillOptions,
) (*bridgev2.FetchMessagesResponse, error) {
	log := zerolog.Ctx(ctx)
	reqQuery := payload.DMRequestQuery{}.Default()
	if fetchParams.Count > 0 {
		reqQuery.Count = fetchParams.Count
	}

	// REST APIs require dash-form conversation IDs.
	restConvID := ConvertConversationIDToREST(conversationID)
	if restConvID == "" {
		log.Warn().
			Str("conversation_id", conversationID).
			Msg("Skipping REST backfill: empty conversation ID")
		return &bridgev2.FetchMessagesResponse{
			Messages: []*bridgev2.BackfillMessage{},
			HasMore:  false,
			Forward:  fetchParams.Forward,
		}, nil
	}

	// When explicitly ignoring anchor IDs, also ignore cursor IDs to avoid
	// passing XChat sequence IDs into REST max_id/min_id pagination.
	cursorForQuery := fetchParams.Cursor
	if opts.IgnoreAnchorForQuery && cursorForQuery != "" {
		log.Debug().
			Str("conversation_id", restConvID).
			Str("original_cursor", string(cursorForQuery)).
			Bool("cursor_cleared_for_rest", true).
			Msg("Ignoring cursor for REST backfill query due to cross-protocol fallback")
		cursorForQuery = ""
	}

	if fetchParams.Forward {
		if cursorForQuery != "" {
			reqQuery.MinID = string(cursorForQuery)
		} else if fetchParams.AnchorMessage != nil && !opts.IgnoreAnchorForQuery {
			reqQuery.MinID = ParseMessageID(fetchParams.AnchorMessage.ID)
		}
	} else if cursorForQuery != "" {
		reqQuery.MaxID = string(cursorForQuery)
	} else if fetchParams.AnchorMessage != nil && !opts.IgnoreAnchorForQuery {
		reqQuery.MaxID = ParseMessageID(fetchParams.AnchorMessage.ID)
	} else {
		log.Debug().
			Str("conversation_id", restConvID).
			Msg("REST backward backfill: no cursor or anchor message, bootstrapping from latest page")
	}

	log.Debug().
		Str("conversation_id", restConvID).
		Str("max_id", reqQuery.MaxID).
		Str("min_id", reqQuery.MinID).
		Int("count", reqQuery.Count).
		Bool("forward", fetchParams.Forward).
		Msg("REST API backfill request")

	resp, err := tc.client.FetchConversationContext(ctx, restConvID, &reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
	if err != nil {
		var twitterErr twittermeow.TwitterError
		// Treat "conversation doesn't exist" as terminal for this backfill task.
		if errors.As(err, &twitterErr) && twitterErr.Code == 279 {
			log.Warn().
				Str("conversation_id", restConvID).
				Int("twitter_error_code", twitterErr.Code).
				Msg("REST backfill conversation not found, marking task complete")
			return &bridgev2.FetchMessagesResponse{
				Messages: []*bridgev2.BackfillMessage{},
				HasMore:  false,
				Forward:  fetchParams.Forward,
			}, nil
		}
		return nil, fmt.Errorf("fetch conversation context: %w", err)
	}

	inbox := resp.ConversationTimeline
	if inbox == nil {
		return &bridgev2.FetchMessagesResponse{
			HasMore: false,
			Forward: fetchParams.Forward,
		}, nil
	}

	messageMap := inbox.SortedMessages(ctx)
	messages := messageMap[restConvID]
	selectedMessageConversationID := restConvID
	if len(messages) == 0 {
		// Older pages can sometimes use colon-format conversation IDs in entry payloads.
		altConvID := NormalizeConversationID(restConvID)
		if altConvID != "" && altConvID != restConvID {
			if altMessages := messageMap[altConvID]; len(altMessages) > 0 {
				messages = altMessages
				selectedMessageConversationID = altConvID
			}
		}
	}
	if len(messages) == 0 && conversationID != "" && conversationID != selectedMessageConversationID {
		if altMessages := messageMap[conversationID]; len(altMessages) > 0 {
			messages = altMessages
			selectedMessageConversationID = conversationID
		}
	}

	rawMessageCount := len(messages)
	filteredAnchorCount := 0
	filteredAnchorIDCount := 0
	filteredAnchorTimeCount := 0
	canonicalIDFromMessageDataCount := 0
	canonicalIDFromEntryCount := 0
	canonicalIDMissingCount := 0
	canonicalTSFromSnowflakeMessageDataIDCount := 0
	canonicalTSFromSnowflakeEntryIDCount := 0
	canonicalTSFromMessageDataTimeCount := 0
	canonicalTSFromEntryTimeCount := 0
	canonicalTSMissingCount := 0

	backfillMessages := make([]*bridgev2.BackfillMessage, 0, len(messages))
	for _, msg := range messages {
		canonicalID, canonicalIDSource := restMessageCanonicalID(msg)
		switch canonicalIDSource {
		case restCanonicalIDSourceMessageData:
			canonicalIDFromMessageDataCount++
		case restCanonicalIDSourceEntry:
			canonicalIDFromEntryCount++
		default:
			canonicalIDMissingCount++
		}

		timestamp, canonicalTSSource := restMessageCanonicalTimestamp(msg, canonicalID, canonicalIDSource)
		switch canonicalTSSource {
		case restCanonicalTSSourceSnowflakeMessageDataID:
			canonicalTSFromSnowflakeMessageDataIDCount++
		case restCanonicalTSSourceSnowflakeEntryID:
			canonicalTSFromSnowflakeEntryIDCount++
		case restCanonicalTSSourceMessageDataTime:
			canonicalTSFromMessageDataTimeCount++
		case restCanonicalTSSourceEntryTime:
			canonicalTSFromEntryTimeCount++
		default:
			canonicalTSMissingCount++
		}

		converted := tc.convertToMatrix(ctx, fetchParams.Portal, tc.connector.br.Bot, &msg.MessageData)
		if converted == nil || len(converted.Parts) == 0 {
			log.Warn().
				Str("conversation_id", restConvID).
				Str("entry_id", msg.ID).
				Str("message_id", canonicalID).
				Msg("Failed to convert REST backfill message")
			continue
		}

		// Keep only messages strictly on the requested side of the anchor.
		// This avoids poisoning the queue when transitioning from XChat IDs to REST IDs.
		if fetchParams.AnchorMessage != nil && !opts.IgnoreAnchorForFiltering {
			anchorID := ParseMessageID(fetchParams.AnchorMessage.ID)
			if canonicalID != "" && canonicalID == anchorID {
				filteredAnchorCount++
				filteredAnchorIDCount++
				continue
			}
			// Compatibility fallback for previously bridged rows keyed by entry ID.
			if msg.ID == anchorID {
				filteredAnchorCount++
				filteredAnchorIDCount++
				continue
			}
			if fetchParams.Forward && !timestamp.IsZero() && timestamp.Before(fetchParams.AnchorMessage.Timestamp) {
				filteredAnchorCount++
				filteredAnchorTimeCount++
				continue
			}
			if !fetchParams.Forward && !timestamp.IsZero() && timestamp.After(fetchParams.AnchorMessage.Timestamp) {
				filteredAnchorCount++
				filteredAnchorTimeCount++
				continue
			}
		}

		messageIDForBridge := canonicalID
		if messageIDForBridge == "" {
			messageIDForBridge = msg.ID
		}
		if timestamp.IsZero() && messageIDForBridge != "" {
			timestamp = methods.ParseSnowflake(messageIDForBridge)
		}

		streamOrder := methods.ParseInt64(messageIDForBridge)
		if streamOrder == 0 && msg.ID != "" {
			streamOrder = methods.ParseInt64(msg.ID)
		}

		backfillMessages = append(backfillMessages, &bridgev2.BackfillMessage{
			ConvertedMessage: converted,
			Sender:           tc.MakeEventSender(msg.MessageData.SenderID),
			ID:               MakeMessageID(messageIDForBridge),
			TxnID:            networkid.TransactionID(messageIDForBridge),
			Timestamp:        timestamp,
			StreamOrder:      streamOrder,
		})
	}

	sortBackfillMessages(backfillMessages)

	hasMore := inbox.Status == types.PaginationStatusHasMore
	forcedOlderPage := false
	result := &bridgev2.FetchMessagesResponse{
		Messages: backfillMessages,
		HasMore:  hasMore,
		Forward:  fetchParams.Forward,
	}
	if fetchParams.Forward {
		if hasMore && inbox.MaxEntryID != "" {
			result.Cursor = networkid.PaginationCursor(inbox.MaxEntryID)
		}
	} else if hasMore && inbox.MinEntryID != "" {
		result.Cursor = networkid.PaginationCursor(inbox.MinEntryID)
	}
	// When transitioning from XChat to REST, the first REST page may contain only
	// newer messages than the anchor. Keep paginating backward using min_entry_id
	// even if status doesn't advertise has_more yet.
	if len(backfillMessages) == 0 &&
		!fetchParams.Forward &&
		fetchParams.AnchorMessage != nil &&
		!opts.IgnoreAnchorForFiltering &&
		inbox.MinEntryID != "" &&
		reqQuery.MaxID != inbox.MinEntryID {
		result.HasMore = true
		result.Cursor = networkid.PaginationCursor(inbox.MinEntryID)
		forcedOlderPage = true
	}
	if len(backfillMessages) == 0 && result.HasMore && result.Cursor == "" {
		result.HasMore = false
	}
	if rawMessageCount == 0 {
		messageMapKeys := make([]string, 0, len(messageMap))
		for key := range messageMap {
			messageMapKeys = append(messageMapKeys, key)
		}
		slices.Sort(messageMapKeys)
		log.Debug().
			Str("conversation_id", restConvID).
			Str("selected_message_conversation_id", selectedMessageConversationID).
			Strs("message_map_keys", messageMapKeys).
			Int("entry_count", len(inbox.Entries)).
			Str("timeline_status", string(inbox.Status)).
			Str("min_entry_id", inbox.MinEntryID).
			Str("max_entry_id", inbox.MaxEntryID).
			Msg("REST backfill parsed zero messages from timeline entries")
	}

	log.Debug().
		Int("raw_message_count", rawMessageCount).
		Int("filtered_anchor_count", filteredAnchorCount).
		Int("filtered_anchor_id_count", filteredAnchorIDCount).
		Int("filtered_anchor_time_count", filteredAnchorTimeCount).
		Int("canonical_id_from_message_data_count", canonicalIDFromMessageDataCount).
		Int("canonical_id_from_entry_count", canonicalIDFromEntryCount).
		Int("canonical_id_missing_count", canonicalIDMissingCount).
		Int("canonical_ts_from_snowflake_message_data_id_count", canonicalTSFromSnowflakeMessageDataIDCount).
		Int("canonical_ts_from_snowflake_entry_id_count", canonicalTSFromSnowflakeEntryIDCount).
		Int("canonical_ts_from_message_data_time_count", canonicalTSFromMessageDataTimeCount).
		Int("canonical_ts_from_entry_time_count", canonicalTSFromEntryTimeCount).
		Int("canonical_ts_missing_count", canonicalTSMissingCount).
		Int("converted_message_count", len(backfillMessages)).
		Int("message_count", len(backfillMessages)).
		Bool("has_more", hasMore).
		Str("min_entry_id", inbox.MinEntryID).
		Str("max_entry_id", inbox.MaxEntryID).
		Str("selected_message_conversation_id", selectedMessageConversationID).
		Str("cursor", string(result.Cursor)).
		Bool("forward", fetchParams.Forward).
		Bool("ignore_anchor_for_query", opts.IgnoreAnchorForQuery).
		Bool("ignore_anchor_for_filtering", opts.IgnoreAnchorForFiltering).
		Bool("forced_older_page", forcedOlderPage).
		Msg("REST API backfill response")

	return result, nil
}

const (
	restCanonicalIDSourceMessageData = "message_data_id"
	restCanonicalIDSourceEntry       = "entry_id"
	restCanonicalIDSourceNone        = "none"
)

func restMessageCanonicalID(msg *types.Message) (id string, source string) {
	if msg == nil {
		return "", restCanonicalIDSourceNone
	}
	if msg.MessageData.ID != "" {
		return msg.MessageData.ID, restCanonicalIDSourceMessageData
	}
	if msg.ID != "" {
		return msg.ID, restCanonicalIDSourceEntry
	}
	return "", restCanonicalIDSourceNone
}

const (
	restCanonicalTSSourceSnowflakeMessageDataID = "snowflake_message_data_id"
	restCanonicalTSSourceSnowflakeEntryID       = "snowflake_entry_id"
	restCanonicalTSSourceMessageDataTime        = "message_data_time"
	restCanonicalTSSourceEntryTime              = "entry_time"
	restCanonicalTSSourceNone                   = "none"
)

func restMessageCanonicalTimestamp(msg *types.Message, canonicalID, canonicalIDSource string) (ts time.Time, source string) {
	if canonicalID != "" {
		if ts = methods.ParseSnowflake(canonicalID); !ts.IsZero() {
			if canonicalIDSource == restCanonicalIDSourceMessageData {
				return ts, restCanonicalTSSourceSnowflakeMessageDataID
			}
			return ts, restCanonicalTSSourceSnowflakeEntryID
		}
	}
	if msg != nil && msg.MessageData.Time != "" {
		if ts = methods.ParseMsecTimestamp(msg.MessageData.Time); !ts.IsZero() {
			return ts, restCanonicalTSSourceMessageDataTime
		}
	}
	if msg != nil && msg.Time != "" {
		if ts = methods.ParseMsecTimestamp(msg.Time); !ts.IsZero() {
			return ts, restCanonicalTSSourceEntryTime
		}
	}
	return time.Time{}, restCanonicalTSSourceNone
}

func compareIntStrings(a, b string) int {
	ai, errA := strconv.ParseInt(a, 10, 64)
	bi, errB := strconv.ParseInt(b, 10, 64)
	if errA == nil && errB == nil {
		switch {
		case ai < bi:
			return -1
		case ai > bi:
			return 1
		default:
			return 0
		}
	}
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func sortBackfillMessages(msgs []*bridgev2.BackfillMessage) {
	slices.SortFunc(msgs, func(a, b *bridgev2.BackfillMessage) int {
		if a == nil || b == nil {
			return 0
		}
		switch {
		case a.StreamOrder != 0 && b.StreamOrder != 0:
			switch {
			case a.StreamOrder < b.StreamOrder:
				return -1
			case a.StreamOrder > b.StreamOrder:
				return 1
			default:
				return 0
			}
		case !a.Timestamp.IsZero() && !b.Timestamp.IsZero():
			if a.Timestamp.Before(b.Timestamp) {
				return -1
			} else if a.Timestamp.After(b.Timestamp) {
				return 1
			}
			return 0
		default:
			return 0
		}
	})
}
