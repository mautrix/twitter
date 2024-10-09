// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
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
	"strings"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridge/status"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	bridgeEvt "maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

type TwitterClient struct {
	connector *TwitterConnector
	client    *twittermeow.Client

	userLogin *bridgev2.UserLogin

	userCache map[string]types.User
}

var (
	_ bridgev2.NetworkAPI                    = (*TwitterClient)(nil)
	_ bridgev2.ReactionHandlingNetworkAPI    = (*TwitterClient)(nil)
	_ bridgev2.ReadReceiptHandlingNetworkAPI = (*TwitterClient)(nil)
)

func NewTwitterClient(ctx context.Context, tc *TwitterConnector, login *bridgev2.UserLogin) (*TwitterClient, error) {
	log := zerolog.Ctx(ctx).With().
		Str("component", "twitter_client").
		Str("user_login_id", string(login.ID)).
		Logger()

	meta := login.Metadata.(*UserLoginMetadata)
	clientOpts := &twittermeow.ClientOpts{
		Cookies:       cookies.NewCookiesFromString(meta.Cookies),
		WithJOTClient: true,
	}
	twitClient := &TwitterClient{
		client:    twittermeow.NewClient(clientOpts, log),
		userLogin: login,
		userCache: make(map[string]types.User),
	}

	twitClient.client.SetEventHandler(twitClient.HandleTwitterEvent)

	return twitClient, nil
}

func (tc *TwitterClient) Connect(ctx context.Context) error {
	if tc.client == nil {
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateBadCredentials,
			Error:      "twitter-not-logged-in",
		})
		return nil
	}

	go tc.syncChannels(ctx)
	return tc.client.Connect()
}

func (tc *TwitterClient) Disconnect() {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) IsLoggedIn() bool {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) LogoutRemote(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	conversationId := string(portal.PortalKey.ID)
	queryConversationPayload := payload.DmRequestQuery{}.Default()
	queryConversationPayload.IncludeConversationInfo = true
	conversationData, err := tc.client.FetchConversationContext(conversationId, queryConversationPayload, payload.CONTEXT_FETCH_DM_CONVERSATION)
	if err != nil {
		return nil, err
	}

	conversations := conversationData.ConversationTimeline.Conversations
	if len(conversations) <= 0 {
		return nil, fmt.Errorf("failed to find conversation by id %s", string(conversationId))
	}

	conversation := conversations[conversationId]
	users := conversationData.ConversationTimeline.Users

	methods.MergeMaps(tc.userCache, users)

	return tc.ConversationToChatInfo(&conversation), nil
}

func (tc *TwitterClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	userInfo := tc.GetUserInfoBridge(string(ghost.ID))
	if userInfo == nil {
		return nil, fmt.Errorf("failed to find user info in cache by id: %s", ghost.ID)
	}
	return userInfo, nil
}

func (tc *TwitterClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *bridgev2.NetworkRoomCapabilities {
	return &bridgev2.NetworkRoomCapabilities{
		FormattedText: false,
		UserMentions:  true,
		RoomMentions:  false,

		Captions:      true,
		Replies:       true,
		Reactions:     true,
		ReactionCount: 1,
	}
}

func (tc *TwitterClient) convertToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *event.XEventMessage) (*bridgev2.ConvertedMessage, error) {
	partId := networkid.PartID("")
	var MessageOptionalPartID *networkid.MessageOptionalPartID
	if msg.ReplyData.ID != "" {
		MessageOptionalPartID = &networkid.MessageOptionalPartID{
			MessageID: networkid.MessageID(msg.ReplyData.ID),
			PartID:    &partId,
		}
	}

	textPart := &bridgev2.ConvertedMessagePart{
		ID:   partId,
		Type: bridgeEvt.EventMessage,
		Content: &bridgeEvt.MessageEventContent{
			MsgType: bridgeEvt.MsgText,
			Body:    msg.Text,
		},
	}

	parts := make([]*bridgev2.ConvertedMessagePart, 0)

	if msg.Attachment != nil {
		convertedAttachmentPart, indices, err := tc.TwitterAttachmentToMatrix(ctx, portal, intent, msg.Attachment)
		if err != nil {
			return nil, err
		}
		parts = append(parts, convertedAttachmentPart)

		RemoveEntityLinkFromText(textPart, indices)
	}

	if len(textPart.Content.Body) > 0 {
		parts = append(parts, textPart)
	}

	cm := &bridgev2.ConvertedMessage{
		ReplyTo: MessageOptionalPartID,
		Parts:   parts,
	}

	return cm, nil
}

func (tc *TwitterClient) MakePortalKey(conv types.Conversation) networkid.PortalKey {
	var receiver networkid.UserLoginID
	if conv.Type == types.ONE_TO_ONE {
		receiver = tc.userLogin.ID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conv.ConversationID),
		Receiver: receiver,
	}
}

func (tc *TwitterClient) MakePortalKeyFromID(conversationId string) networkid.PortalKey {
	var receiver networkid.UserLoginID
	if strings.Contains(conversationId, "-") {
		receiver = tc.userLogin.ID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conversationId),
		Receiver: receiver,
	}
}
