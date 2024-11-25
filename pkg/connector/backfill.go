package connector

import (
	"context"

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
	if !params.Forward {
		if params.Cursor != "" {
			reqQuery.MaxID = string(params.Cursor)
		} else if params.AnchorMessage != nil {
			reqQuery.MaxID = string(params.AnchorMessage.ID)
		}
	} else if params.AnchorMessage != nil {
		reqQuery.MinID = string(params.AnchorMessage.ID)
	}

	messageResp, err := tc.client.FetchConversationContext(conversationID, &reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
	if err != nil {
		return nil, err
	}

	messages, err := messageResp.ConversationTimeline.GetMessageEntriesByConversationID(conversationID, true)
	if err != nil {
		return nil, err
	}

	converted := make([]*bridgev2.BackfillMessage, len(messages))
	for i, msg := range messages {
		converted[i] = tc.convertBackfillMessage(ctx, params.Portal, msg)
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

func (tc *TwitterClient) convertBackfillMessage(ctx context.Context, portal *bridgev2.Portal, message types.Message) *bridgev2.BackfillMessage {
	sentAt, _ := methods.UnixStringMilliToTime(message.MessageData.Time)
	// TODO get correct intent
	intent := tc.userLogin.Bridge.Matrix.BotIntent()
	return &bridgev2.BackfillMessage{
		ConvertedMessage: tc.convertToMatrix(ctx, portal, intent, &message.MessageData),
		Sender: bridgev2.EventSender{
			IsFromMe: message.MessageData.SenderID == string(tc.userLogin.ID),
			Sender:   networkid.UserID(message.MessageData.SenderID),
		},
		ID:        networkid.MessageID(message.MessageData.ID),
		Timestamp: sentAt,
		Reactions: tc.convertBackfillReactions(message.MessageReactions),
	}
}

func (tc *TwitterClient) convertBackfillReactions(reactions []types.MessageReaction) []*bridgev2.BackfillReaction {
	backfillReactions := make([]*bridgev2.BackfillReaction, 0)
	for _, reaction := range reactions {
		reactionTime, _ := methods.UnixStringMilliToTime(reaction.Time)
		backfillReaction := &bridgev2.BackfillReaction{
			Timestamp: reactionTime,
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
