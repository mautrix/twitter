package connector

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

var _ bridgev2.BackfillingNetworkAPI = (*TwitterClient)(nil)

func (tc *TwitterClient) FetchMessages(ctx context.Context, params bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	conversationID := string(params.Portal.PortalKey.ID)

	reqQuery := payload.DMRequestQuery{}.Default()
	reqQuery.Count = params.Count
	log := zerolog.Ctx(ctx)
	log.Debug().
		Bool("forward", params.Forward).
		Str("cursor", string(params.Cursor)).
		Int("count", params.Count).
		Msg("Backfill params")
	if params.AnchorMessage != nil {
		log.Debug().
			Time("anchor_ts", params.AnchorMessage.Timestamp).
			Str("anchor_id", string(params.AnchorMessage.ID)).
			Msg("Backfill anchor message")
	}
	if !params.Forward {
		if params.Cursor != "" {
			reqQuery.MaxID = string(params.Cursor)
			log.Debug().Msg("Using cursor as max ID")
		} else if params.AnchorMessage != nil {
			reqQuery.MaxID = string(params.AnchorMessage.ID)
			log.Debug().Msg("Using anchor as max ID")
		} else {
			return nil, fmt.Errorf("no cursor or anchor message provided for backward backfill")
		}
	} else if params.AnchorMessage != nil {
		reqQuery.MinID = string(params.AnchorMessage.ID)
		log.Debug().Msg("Using anchor as min ID")
	} else {
		log.Debug().Msg("No anchor for forward backfill, fetching latest messages")
	}

	messageResp, err := tc.client.FetchConversationContext(conversationID, &reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
	if err != nil {
		return nil, err
	}

	messages, err := messageResp.ConversationTimeline.GetMessageEntriesByConversationID(conversationID, true)
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		log.Debug().
			Str("timeline_status", string(messageResp.ConversationTimeline.Status)).
			Msg("No messages in backfill response")
		return &bridgev2.FetchMessagesResponse{
			Messages: make([]*bridgev2.BackfillMessage, 0),
			HasMore:  false,
			Forward:  params.Forward,
		}, nil
	}

	converted := make([]*bridgev2.BackfillMessage, 0, len(messages))
	log.Debug().
		Int("message_count", len(messages)).
		Str("oldest_raw_ts", messages[0].Time).
		Str("newest_raw_ts", messages[len(messages)-1].Time).
		Msg("Fetched messages")
	for _, msg := range messages {
		messageTS := methods.ParseSnowflake(msg.ID)
		log := log.With().
			Str("message_id", msg.MessageData.ID).
			Str("message_raw_ts", msg.Time).
			Time("message_ts", messageTS).
			Logger()
		if params.AnchorMessage != nil {
			if string(params.AnchorMessage.ID) == msg.ID {
				log.Warn().Msg("Skipping anchor message")
				continue
			} else if params.Forward && messageTS.Before(params.AnchorMessage.Timestamp) {
				log.Warn().Msg("Skipping too old message in forwards backfill")
				continue
			} else if !params.Forward && messageTS.After(params.AnchorMessage.Timestamp) {
				log.Warn().Msg("Skipping too new message in backwards backfill")
				continue
			}
		}
		log.Trace().Msg("Converting message")
		// TODO get correct intent
		intent := tc.userLogin.Bridge.Matrix.BotIntent()
		convertedMsg := &bridgev2.BackfillMessage{
			ConvertedMessage: tc.convertToMatrix(ctx, params.Portal, intent, &msg.MessageData),
			Sender: bridgev2.EventSender{
				IsFromMe: msg.MessageData.SenderID == string(tc.userLogin.ID),
				Sender:   networkid.UserID(msg.MessageData.SenderID),
			},
			ID:        networkid.MessageID(msg.MessageData.ID),
			Timestamp: messageTS,
			Reactions: tc.convertBackfillReactions(msg.MessageReactions),
		}
		converted = append(converted, convertedMsg)
	}

	fetchMessagesResp := &bridgev2.FetchMessagesResponse{
		Messages: converted,
		HasMore:  messageResp.ConversationTimeline.Status == types.HAS_MORE,
		Forward:  params.Forward,
	}
	if !params.Forward {
		fetchMessagesResp.Cursor = networkid.PaginationCursor(messageResp.ConversationTimeline.MinEntryID)
	}

	return fetchMessagesResp, nil
}

func (tc *TwitterClient) convertBackfillReactions(reactions []types.MessageReaction) []*bridgev2.BackfillReaction {
	backfillReactions := make([]*bridgev2.BackfillReaction, 0)
	for _, reaction := range reactions {
		backfillReaction := &bridgev2.BackfillReaction{
			Timestamp: methods.ParseSnowflake(reaction.ID),
			Sender: bridgev2.EventSender{
				IsFromMe: reaction.SenderID == string(tc.userLogin.ID),
				Sender:   networkid.UserID(reaction.SenderID),
			},
			EmojiID: "",
			Emoji:   reaction.EmojiReaction,
		}
		backfillReactions = append(backfillReactions, backfillReaction)
	}
	return backfillReactions
}
