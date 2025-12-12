package connector

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

var _ bridgev2.BackfillingNetworkAPI = (*TwitterClient)(nil)
var _ bridgev2.BackfillingNetworkAPIWithLimits = (*TwitterClient)(nil)

const xchatBackfillMaxInt = "9223372036854775807"

func (tc *TwitterClient) GetBackfillMaxBatchCount(_ context.Context, portal *bridgev2.Portal, _ *database.BackfillTask) int {
	key := "channel"
	if portal != nil {
		switch portal.RoomType {
		case database.RoomTypeDM:
			key = "dm"
		case database.RoomTypeGroupDM:
			key = "group_dm"
		}
	}
	return tc.connector.br.Config.Backfill.Queue.GetOverride(key)
}

func (tc *TwitterClient) FetchMessages(ctx context.Context, fetchParams bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	if fetchParams.Forward {
		return &bridgev2.FetchMessagesResponse{
			Forward: true,
		}, nil
	}
	if fetchParams.Portal == nil {
		return nil, fmt.Errorf("portal is nil")
	}

	conversationID := string(fetchParams.Portal.PortalKey.ID)
	minSeqID := string(fetchParams.AnchorMessage.ID)

	zerolog.Ctx(ctx).Info().
		Str("conversation_id", conversationID).
		Bool("forward", fetchParams.Forward).
		Str("cursor", minSeqID).
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
			ID:               networkid.MessageID(msgID),
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
	var err error

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

		entry, err = crypto.DecryptMessageEntryContentsBytesDebug(contentsBytes, convKey.Key, &debugLog)
	} else {
		entry, err = crypto.ParseMessageEntryContentsBytes(contentsBytes)
	}
	if err != nil || entry == nil || entry.Message == nil {
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
