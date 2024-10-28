package connector

import (
	"context"

	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

var _ bridgev2.BackfillingNetworkAPI = (*TwitterClient)(nil)

func (tc *TwitterClient) FetchMessages(ctx context.Context, params bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	conversationId := string(params.Portal.PortalKey.ID)

	reqQuery := ptr.Ptr(payload.DmRequestQuery{}.Default())
	reqQuery.Count = params.Count
	reqQuery.MaxID = string(params.Cursor)
	messageResp, err := tc.client.FetchConversationContext(conversationId, reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
	if err != nil {
		return nil, err
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
