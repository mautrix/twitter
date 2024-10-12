package connector

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (tc *TwitterClient) syncChannels(_ context.Context) {
	//log := zerolog.Ctx(ctx)

	reqQuery := ptr.Ptr(payload.DmRequestQuery{}.Default())
	initalInboxState, err := tc.client.GetInitialInboxState(reqQuery)
	if err != nil {
		panic(fmt.Sprintf("failed to fetch initial inbox state: %s", err.Error()))
	}

	inboxData := initalInboxState.InboxInitialState
	trustedInbox := inboxData.InboxTimelines.Trusted
	cursor := trustedInbox.MinEntryID
	paginationStatus := trustedInbox.Status

	// loop until no more threads can be found
	// add backfill configuration to limit this later
	for paginationStatus == types.HAS_MORE {
		reqQuery.MaxID = cursor
		nextInboxTimelineResponse, err := tc.client.FetchTrustedThreads(reqQuery)
		if err != nil {
			panic(fmt.Sprintf("failed to fetch threads in trusted inbox using cursor %s: %s", cursor, err.Error()))
		}

		methods.MergeMaps(inboxData.Conversations, nextInboxTimelineResponse.InboxTimeline.Conversations)
		methods.MergeMaps(inboxData.Users, nextInboxTimelineResponse.InboxTimeline.Users)
		inboxData.Entries = append(inboxData.Entries, nextInboxTimelineResponse.InboxTimeline.Entries...)

		cursor = nextInboxTimelineResponse.InboxTimeline.MinEntryID
		paginationStatus = nextInboxTimelineResponse.InboxTimeline.Status
	}

	methods.MergeMaps(tc.userCache, inboxData.Users)

	conversations, err := inboxData.Prettify()
	if err != nil {
		panic(fmt.Sprintf("failed to prettify inbox data after fetching conversations: %s", err.Error()))
	}

	for _, convInboxData := range conversations {
		conv := convInboxData.Conversation
		methods.SortMessagesByTime(convInboxData.Messages)
		messages := convInboxData.Messages
		latestMessage := messages[len(messages)-1]
		latestMessageTS, err := methods.UnixStringMilliToTime(latestMessage.MessageData.Time)
		if err != nil {
			panic(fmt.Sprintf("failed to convert latest message TS to time.Time: %s", err.Error()))
		}
		evt := &simplevent.ChatResync{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatResync,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("portal_key", conv.ConversationID)
				},
				PortalKey:    tc.MakePortalKey(conv),
				CreatePortal: true,
			},
			ChatInfo:        tc.ConversationToChatInfo(&conv),
			LatestMessageTS: latestMessageTS,
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, evt)
	}
}
