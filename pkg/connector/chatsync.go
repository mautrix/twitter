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

func (tc *TwitterClient) syncChannels(ctx context.Context, inbox *response.TwitterInboxData) {
	log := zerolog.Ctx(ctx)

	reqQuery := ptr.Ptr(payload.DMRequestQuery{}.Default())
	if inbox == nil {
		initialInboxState, err := tc.client.GetInitialInboxState(ctx, reqQuery)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch initial inbox state:")
			return
		}
		inbox = initialInboxState.InboxInitialState
	}

	trustedInbox := inbox.InboxTimelines.Trusted
	cursor := trustedInbox.MinEntryID
	paginationStatus := trustedInbox.Status

	for paginationStatus == types.PaginationStatusHasMore && (tc.connector.Config.ConversationSyncLimit == -1 || len(inbox.Entries) < tc.connector.Config.ConversationSyncLimit) {
		reqQuery.MaxID = cursor
		nextInboxTimelineResponse, err := tc.client.FetchTrustedThreads(ctx, reqQuery)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to fetch threads in trusted inbox using cursor %s:", cursor))
			return
		} else if len(nextInboxTimelineResponse.InboxTimeline.Entries) == 0 {
			break
		}

		if inbox.Conversations == nil {
			inbox.Conversations = map[string]*types.Conversation{}
		}
		if inbox.Users == nil {
			inbox.Users = map[string]*types.User{}
		}
		maps.Copy(inbox.Conversations, nextInboxTimelineResponse.InboxTimeline.Conversations)
		maps.Copy(inbox.Users, nextInboxTimelineResponse.InboxTimeline.Users)
		inbox.Entries = append(inbox.Entries, nextInboxTimelineResponse.InboxTimeline.Entries...)

		cursor = nextInboxTimelineResponse.InboxTimeline.MinEntryID
		paginationStatus = nextInboxTimelineResponse.InboxTimeline.Status
	}

	tc.userCacheLock.Lock()
	maps.Copy(tc.userCache, inbox.Users)
	tc.userCacheLock.Unlock()

	messages := inbox.SortedMessages(ctx)
	for _, conv := range inbox.SortedConversations() {
		convMessages := messages[conv.ConversationID]
		if len(convMessages) == 0 {
			continue
		}
		latestMessage := convMessages[len(convMessages)-1]
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
				CreatePortal: conv.Trusted,
			},
			ChatInfo: tc.conversationToChatInfo(conv, inbox),
			BundledBackfillData: &backfillDataBundle{
				Conv:     conv,
				Messages: convMessages,
				Inbox:    inbox,
			},
			LatestMessageTS: latestMessageTS,
		}
		tc.connector.br.QueueRemoteEvent(tc.userLogin, evt)
	}
}

type backfillDataBundle struct {
	Conv     *types.Conversation
	Messages []*types.Message
	Inbox    *response.TwitterInboxData
}
