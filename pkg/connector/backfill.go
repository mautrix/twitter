// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2025 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
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

	var inbox *response.TwitterInboxData
	var conv *types.Conversation
	var messages []*types.Message
	bundle, ok := params.BundledData.(*backfillDataBundle)
	if ok && params.Forward && len(bundle.Messages) > 0 {
		inbox = bundle.Inbox
		conv = bundle.Conv
		messages = bundle.Messages
		// TODO support for fetching more messages
	} else {
		messageResp, err := tc.client.FetchConversationContext(conversationID, &reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
		if err != nil {
			return nil, err
		}
		inbox = messageResp.ConversationTimeline
		conv = inbox.GetConversationByID(conversationID)
		messages = messageResp.ConversationTimeline.SortedMessages(ctx)[conversationID]
	}

	if len(messages) == 0 {
		log.Debug().
			Str("timeline_status", string(inbox.Status)).
			Msg("No messages in backfill response")
		return &bridgev2.FetchMessagesResponse{
			Messages: make([]*bridgev2.BackfillMessage, 0),
			HasMore:  false,
			Forward:  params.Forward,
		}, nil
	}

	converted := make([]*bridgev2.BackfillMessage, 0, len(messages))
	log.Debug().
		Bool("is_bundled_data", bundle != nil).
		Int("message_count", len(messages)).
		Str("oldest_raw_ts", messages[0].Time).
		Str("newest_raw_ts", messages[len(messages)-1].Time).
		Str("oldest_id", messages[0].ID).
		Str("newest_id", messages[len(messages)-1].ID).
		Str("last_read_id", conv.LastReadEventID).
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
			Sender:           tc.MakeEventSender(msg.MessageData.SenderID),
			ID:               networkid.MessageID(msg.MessageData.ID),
			Timestamp:        messageTS,
			Reactions:        tc.convertBackfillReactions(msg.MessageReactions),
			StreamOrder:      methods.ParseSnowflakeInt(msg.MessageData.ID),
		}
		converted = append(converted, convertedMsg)
	}

	fetchMessagesResp := &bridgev2.FetchMessagesResponse{
		Messages: converted,
		HasMore:  bundle != nil || inbox.Status == types.PaginationStatusHasMore,
		Forward:  params.Forward,
		MarkRead: conv.LastReadEventID == messages[len(messages)-1].ID,
	}
	if !params.Forward {
		fetchMessagesResp.Cursor = networkid.PaginationCursor(inbox.MinEntryID)
	}

	return fetchMessagesResp, nil
}

func (tc *TwitterClient) convertBackfillReactions(reactions []types.MessageReaction) []*bridgev2.BackfillReaction {
	backfillReactions := make([]*bridgev2.BackfillReaction, 0)
	for _, reaction := range reactions {
		backfillReaction := &bridgev2.BackfillReaction{
			Timestamp: methods.ParseSnowflake(reaction.ID),
			Sender:    tc.MakeEventSender(reaction.SenderID),
			EmojiID:   "",
			Emoji:     reaction.EmojiReaction,
			// StreamOrder: methods.ParseSnowflakeInt(reaction.ID),
		}
		backfillReactions = append(backfillReactions, backfillReaction)
	}
	return backfillReactions
}
