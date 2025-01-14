package connector

import (
	"context"
	"fmt"
	"maps"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (tc *TwitterClient) syncChannels(ctx context.Context, initialInboxState *response.InboxInitialStateResponse) {
	log := zerolog.Ctx(ctx)

	reqQuery := ptr.Ptr(payload.DMRequestQuery{}.Default())
	if initialInboxState == nil {
		var err error
		initialInboxState, err = tc.client.GetInitialInboxState(reqQuery)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch initial inbox state:")
			return
		}
	}

	inboxData := initialInboxState.InboxInitialState
	trustedInbox := inboxData.InboxTimelines.Trusted
	cursor := trustedInbox.MinEntryID
	paginationStatus := trustedInbox.Status

	for paginationStatus == types.HAS_MORE && (tc.connector.Config.ConversationSyncLimit == -1 || len(inboxData.Entries) < tc.connector.Config.ConversationSyncLimit) {
		reqQuery.MaxID = cursor
		nextInboxTimelineResponse, err := tc.client.FetchTrustedThreads(reqQuery)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to fetch threads in trusted inbox using cursor %s:", cursor))
			return
		} else if len(nextInboxTimelineResponse.InboxTimeline.Entries) == 0 {
			break
		}

		if inboxData.Conversations == nil {
			inboxData.Conversations = map[string]types.Conversation{}
		}
		maps.Copy(inboxData.Conversations, nextInboxTimelineResponse.InboxTimeline.Conversations)
		if inboxData.Users == nil {
			inboxData.Users = map[string]types.User{}
		}
		maps.Copy(inboxData.Users, nextInboxTimelineResponse.InboxTimeline.Users)
		inboxData.Entries = append(inboxData.Entries, nextInboxTimelineResponse.InboxTimeline.Entries...)

		cursor = nextInboxTimelineResponse.InboxTimeline.MinEntryID
		paginationStatus = nextInboxTimelineResponse.InboxTimeline.Status
	}

	maps.Copy(tc.userCache, inboxData.Users)

	conversations, err := inboxData.Prettify()
	if err != nil {
		log.Error().Err(err).Msg("failed to prettify inbox data after fetching conversations:")
		return
	}

	for _, convInboxData := range conversations {
		conv := convInboxData.Conversation
		methods.SortMessagesByTime(convInboxData.Messages)
		messages := convInboxData.Messages
		if len(messages) == 0 {
			continue
		}
		latestMessage := messages[len(messages)-1]
		latestMessageTS := methods.ParseSnowflake(latestMessage.MessageData.ID)
		evt := &simplevent.ChatResync{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatResync,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", conv.ConversationID).
						Bool("conv_low_quality", conv.LowQuality).
						Bool("conv_trusted", conv.Trusted)
				},
				PortalKey:    tc.MakePortalKey(conv),
				CreatePortal: !conv.LowQuality,
			},
			ChatInfo:        tc.ConversationToChatInfo(&conv),
			LatestMessageTS: latestMessageTS,
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, evt)
	}
}
