package connector

import (
	"context"
	"log"
	"os"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

var _ bridgev2.BackfillingNetworkAPI = (*TwitterClient)(nil)

func (tc *TwitterClient) FetchMessages(ctx context.Context, params bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	conversationId := string(params.Portal.PortalKey.ID)
	cursor := params.Cursor
	//count := params.Count

	reqQuery := payload.DmRequestQuery{}.Default()
	reqQuery.Count = 25
	messageResp, err := tc.client.FetchConversationContext(conversationId, reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
	if err != nil {
		return nil, err
	}

	if cursor != "" {
		log.Println("found cursor:", params)
		os.Exit(1)
	}

	messages, err := messageResp.ConversationTimeline.GetMessageEntriesByConversationID(conversationId, true)
	if err != nil {
		return nil, err
	}

	backfilledMessages, err := tc.MessagesToBackfillMessages(ctx, messages, messageResp.ConversationTimeline.GetConversationByID(conversationId))
	if err != nil {
		return nil, err
	}

	fetchMessagesResp := &bridgev2.FetchMessagesResponse{
		Messages: backfilledMessages,
		Cursor:   networkid.PaginationCursor(messageResp.ConversationTimeline.MinEntryID),
		HasMore:  messageResp.ConversationTimeline.Status == types.HAS_MORE,
		Forward:  params.Forward,
	}

	return fetchMessagesResp, nil
}
