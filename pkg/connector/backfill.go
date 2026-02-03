package connector

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"go.mau.fi/util/variationselector"
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

func (tc *TwitterClient) FetchMessages(ctx context.Context, fetchParams bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	conversationID := ParsePortalID(fetchParams.Portal.PortalKey.ID)
	meta := fetchParams.Portal.Metadata.(*PortalMetadata)

	// Use REST API for conversations without encryption keys.
	// XChat GraphQL API requires encryption keys to decrypt messages.
	if !meta.CanUseXChat() {
		return tc.fetchRESTMessages(ctx, conversationID, fetchParams)
	}

	if fetchParams.Forward {
		// Forward backfill for XChat: initial messages are already processed during room creation
		// via ProcessMessageAndReadEvents(). No additional forward fetch needed.
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

	return &bridgev2.FetchMessagesResponse{
		Messages: msgs,
		Cursor:   networkid.PaginationCursor(nextCursor),
		HasMore:  hasMore,
	}, nil
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

	return tc.client.GetKeyManager().PutConversationKey(ctx, &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     newKeyVersion,
		Key:            convKeyBytes,
		CreatedAt:      time.Now(),
	})
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

// fetchRESTMessages fetches messages via the REST API for unencrypted conversations.
func (tc *TwitterClient) fetchRESTMessages(ctx context.Context, conversationID string, fetchParams bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	log := zerolog.Ctx(ctx)

	reqQuery := payload.DMRequestQuery{}.Default()

	// Use REST conversation ID format (hyphen-separated)
	restConvID := ConvertConversationIDToREST(conversationID)

	if fetchParams.Forward {
		// Forward backfill: fetch recent messages.
		// For initial sync, no anchor means fetch the most recent messages.
		// With anchor, fetch messages newer than the anchor (using min_id).
		if fetchParams.Cursor != "" {
			reqQuery.MinID = string(fetchParams.Cursor)
		} else if fetchParams.AnchorMessage != nil {
			reqQuery.MinID = ParseMessageID(fetchParams.AnchorMessage.ID)
		}
		// No min_id means fetch the most recent messages (initial forward backfill)
	} else {
		// Backward backfill: fetch older messages using max_id.
		if fetchParams.Cursor != "" {
			reqQuery.MaxID = string(fetchParams.Cursor)
		} else if fetchParams.AnchorMessage != nil {
			reqQuery.MaxID = ParseMessageID(fetchParams.AnchorMessage.ID)
		} else {
			// No cursor or anchor for backward backfill - this shouldn't happen
			// after forward backfill has run, but log it for debugging.
			log.Debug().
				Str("conversation_id", conversationID).
				Msg("REST backward backfill: no cursor or anchor message")
			return &bridgev2.FetchMessagesResponse{HasMore: false, Forward: false}, nil
		}
	}

	log.Debug().
		Str("conversation_id", restConvID).
		Str("max_id", reqQuery.MaxID).
		Str("min_id", reqQuery.MinID).
		Bool("forward", fetchParams.Forward).
		Msg("REST API backfill request")

	resp, err := tc.client.FetchConversationContext(ctx, restConvID, &reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
	if err != nil {
		return nil, fmt.Errorf("fetch conversation context: %w", err)
	}

	inbox := resp.ConversationTimeline
	if inbox == nil {
		return &bridgev2.FetchMessagesResponse{HasMore: false}, nil
	}

	// Get sorted messages for this conversation
	messages := inbox.SortedMessages(ctx)[restConvID]

	if len(messages) == 0 {
		return &bridgev2.FetchMessagesResponse{
			Messages: []*bridgev2.BackfillMessage{},
			HasMore:  false,
			Forward:  fetchParams.Forward,
		}, nil
	}

	// Convert messages to BackfillMessage format
	backfillMessages := make([]*bridgev2.BackfillMessage, 0, len(messages))
	for _, msg := range messages {
		ts := methods.ParseSnowflake(msg.ID)
		streamOrder := methods.ParseInt64(msg.ID)

		converted := tc.convertToMatrix(ctx, fetchParams.Portal, tc.connector.br.Bot, &msg.MessageData)
		if converted == nil || len(converted.Parts) == 0 {
			log.Warn().Str("message_id", msg.ID).Msg("Failed to convert message for backfill")
			continue
		}

		backfillMessages = append(backfillMessages, &bridgev2.BackfillMessage{
			ConvertedMessage: converted,
			Sender:           tc.MakeEventSender(msg.MessageData.SenderID),
			ID:               MakeMessageID(msg.ID),
			TxnID:            networkid.TransactionID(msg.ID),
			Timestamp:        ts,
			StreamOrder:      streamOrder,
			Reactions:        tc.convertBackfillReactions(msg.MessageReactions),
		})
	}

	sortBackfillMessages(backfillMessages)

	hasMore := inbox.Status == types.PaginationStatusHasMore

	result := &bridgev2.FetchMessagesResponse{
		Messages: backfillMessages,
		HasMore:  hasMore,
		Forward:  fetchParams.Forward,
	}

	// Set cursor for next backfill request
	if fetchParams.Forward {
		// Forward backfill: use MaxEntryID to get newer messages next time
		if hasMore && inbox.MaxEntryID != "" {
			result.Cursor = networkid.PaginationCursor(inbox.MaxEntryID)
		}
	} else {
		// Backward backfill: use MinEntryID to get older messages next time
		if hasMore && inbox.MinEntryID != "" {
			result.Cursor = networkid.PaginationCursor(inbox.MinEntryID)
		}
	}

	log.Debug().
		Int("message_count", len(backfillMessages)).
		Bool("has_more", hasMore).
		Str("min_entry_id", inbox.MinEntryID).
		Str("max_entry_id", inbox.MaxEntryID).
		Str("cursor", string(result.Cursor)).
		Bool("forward", fetchParams.Forward).
		Msg("REST API backfill response")

	return result, nil
}

func (tc *TwitterClient) convertBackfillReactions(reactions []types.MessageReaction) []*bridgev2.BackfillReaction {
	if len(reactions) == 0 {
		return nil
	}
	backfillReactions := make([]*bridgev2.BackfillReaction, 0, len(reactions))
	for _, reaction := range reactions {
		emoji := variationselector.FullyQualify(reaction.EmojiReaction)
		backfillReactions = append(backfillReactions, &bridgev2.BackfillReaction{
			Timestamp: methods.ParseSnowflake(reaction.ID),
			Sender:    tc.MakeEventSender(reaction.SenderID),
			Emoji:     emoji,
		})
	}
	return backfillReactions
}
