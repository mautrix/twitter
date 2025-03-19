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
	"maps"
	"strings"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/status"
	bridgeEvt "maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/connector/twitterfmt"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type TwitterClient struct {
	connector *TwitterConnector
	client    *twittermeow.Client

	userLogin *bridgev2.UserLogin

	userCache map[string]types.User

	matrixParser *format.HTMLParser
}

var (
	_ bridgev2.NetworkAPI         = (*TwitterClient)(nil)
	_ bridgev2.PushableNetworkAPI = (*TwitterClient)(nil)
)

func NewTwitterClient(login *bridgev2.UserLogin, connector *TwitterConnector, client *twittermeow.Client) *TwitterClient {
	tc := &TwitterClient{
		connector: connector,
		client:    client,
		userLogin: login,
		userCache: make(map[string]types.User),
	}
	client.SetEventHandler(tc.HandleTwitterEvent)
	tc.matrixParser = &format.HTMLParser{
		TabsToSpaces:   4,
		Newline:        "\n",
		HorizontalLine: "\n---\n",
		PillConverter: func(displayname, mxid, eventID string, ctx format.Context) string {
			userID, ok := tc.connector.br.Matrix.ParseGhostMXID(id.UserID(mxid))
			if !ok {
				return displayname
			}
			ghost, err := tc.connector.br.GetGhostByID(context.TODO(), userID)
			if err != nil || len(ghost.Identifiers) < 1 {
				return displayname
			}
			id := ghost.Identifiers[0]
			return "@" + strings.TrimPrefix(id, "twitter:")
		},
	}
	return tc
}

func (tc *TwitterConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	login.Client = NewTwitterClient(login, tc, twittermeow.NewClient(&twittermeow.ClientOpts{
		Cookies:       cookies.NewCookiesFromString(login.Metadata.(*UserLoginMetadata).Cookies),
		WithJOTClient: true,
	}, login.Log.With().Str("component", "twitter_client").Logger()))
	return nil
}

var pushCfg = &bridgev2.PushConfig{
	Web: &bridgev2.WebPushConfig{VapidKey: "BF5oEo0xDUpgylKDTlsd8pZmxQA1leYINiY-rSscWYK_3tWAkz4VMbtf1MLE_Yyd6iII6o-e3Q9TCN5vZMzVMEs"},
}

var pushSettings = &payload.PushNotificationSettings{
	Addressbook:     "off",
	Ads:             "off",
	DirectMessages:  "on",
	DMReaction:      "reaction_your_own",
	FollowersNonVit: "off",
	FollowersVit:    "off",
	LifelineAlerts:  "off",
	LikesNonVit:     "off",
	LikesVit:        "off",
	LiveVideo:       "off",
	Mentions:        "off",
	Moments:         "off",
	News:            "off",
	PhotoTags:       "off",
	Recommendations: "off",
	Retweets:        "off",
	Spaces:          "off",
	Topics:          "off",
	Tweets:          "off",
}

func (tc *TwitterClient) GetPushConfigs() *bridgev2.PushConfig {
	return pushCfg
}

func (tc *TwitterClient) RegisterPushNotifications(ctx context.Context, pushType bridgev2.PushType, token string) error {
	if tc.client == nil {
		return bridgev2.ErrNotLoggedIn
	}
	switch pushType {
	case bridgev2.PushTypeWeb:
		meta := tc.userLogin.Metadata.(*UserLoginMetadata)
		if meta.PushKeys == nil {
			meta.GeneratePushKeys()
			err := tc.userLogin.Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to save push key: %w", err)
			}
		}
		pc := twittermeow.WebPushConfig{
			Endpoint: token,
			Auth:     meta.PushKeys.Auth,
			P256DH:   meta.PushKeys.P256DH,
		}
		err := tc.client.SetPushNotificationConfig(twittermeow.PushRegister, pc)
		if err != nil {
			return fmt.Errorf("failed to set push notification config: %w", err)
		}
		pc.Settings = pushSettings
		err = tc.client.SetPushNotificationConfig(twittermeow.PushSave, pc)
		if err != nil {
			return fmt.Errorf("failed to set push notification preferences: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported push type: %v", pushType)
	}
}

func (tc *TwitterClient) Connect(ctx context.Context) {
	if tc.client == nil {
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateBadCredentials,
			Error:      "twitter-not-logged-in",
		})
		return
	}

	inboxState, currentUser, err := tc.client.LoadMessagesPage()
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to load messages page")
		if twittermeow.IsAuthError(err) {
			tc.userLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateBadCredentials,
				Error:      "twitter-invalid-credentials",
				Message:    err.Error(),
			})
		} else {
			tc.userLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateUnknownError,
				Error:      "twitter-load-error",
			})
		}
		return
	}

	selfUser := inboxState.InboxInitialState.GetUserByID(tc.client.GetCurrentUserID())
	if tc.userLogin.RemoteName != currentUser.ScreenName || tc.userLogin.RemoteProfile.Name != selfUser.Name {
		tc.userLogin.RemoteName = currentUser.ScreenName
		tc.userLogin.RemoteProfile = status.RemoteProfile{
			// TODO fetch from /1.1/users/email_phone_info.json?
			Phone:    "",
			Email:    "",
			Username: currentUser.ScreenName,
			Name:     selfUser.Name,
			// TODO set on ghost and reuse same mxc
			Avatar: "",
		}
		err = tc.userLogin.Save(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to save user login after updating remote profile")
		}
	}

	go tc.syncChannels(ctx, inboxState)
	tc.startPolling(ctx)
}

func (tc *TwitterClient) startPolling(ctx context.Context) {
	err := tc.client.Connect()
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to start polling")
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateUnknownError,
			Error:      "twitter-connect-error",
		})
	} else {
		tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})
	}
}

func (tc *TwitterClient) Disconnect() {
	err := tc.client.Disconnect()
	if err != nil {
		tc.userLogin.Log.Error().Err(err).Msg("failed to disconnect, err:")
	}
}

func (tc *TwitterClient) IsLoggedIn() bool {
	return tc.client.IsLoggedIn()
}

func (tc *TwitterClient) LogoutRemote(ctx context.Context) {
	log := zerolog.Ctx(ctx)
	_, err := tc.client.Logout()
	if err != nil {
		log.Error().Err(err).Msg("error logging out")
	}
}

func (tc *TwitterClient) IsThisUser(_ context.Context, userID networkid.UserID) bool {
	return networkid.UserID(tc.client.GetCurrentUserID()) == userID
}

func (tc *TwitterClient) GetCurrentUser() (user *types.User, err error) {
	// TODO wtf is this
	_, settings, err := tc.client.LoadMessagesPage()
	if err != nil {
		return nil, err
	}
	searchResponse, err := tc.client.Search(payload.SearchQuery{
		Query:      settings.ScreenName,
		ResultType: payload.SEARCH_RESULT_TYPE_USERS,
	})
	if err != nil {
		return nil, err
	}
	user = &searchResponse.Users[0]
	return
}

func (tc *TwitterClient) GetChatInfo(_ context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	conversationID := string(portal.PortalKey.ID)
	queryConversationPayload := payload.DMRequestQuery{}.Default()
	queryConversationPayload.IncludeConversationInfo = true
	conversationData, err := tc.client.FetchConversationContext(conversationID, &queryConversationPayload, payload.CONTEXT_FETCH_DM_CONVERSATION)
	if err != nil {
		return nil, err
	}

	conversations := conversationData.ConversationTimeline.Conversations
	if len(conversations) <= 0 {
		return nil, fmt.Errorf("failed to find conversation by id %s", string(conversationID))
	}

	conversation := conversations[conversationID]
	users := conversationData.ConversationTimeline.Users

	maps.Copy(tc.userCache, users)

	return tc.ConversationToChatInfo(&conversation), nil
}

func (tc *TwitterClient) GetUserInfo(_ context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	userInfo := tc.GetUserInfoBridge(string(ghost.ID))
	if userInfo == nil {
		return nil, fmt.Errorf("failed to find user info in cache by id: %s", ghost.ID)
	}
	return userInfo, nil
}

func (tc *TwitterClient) convertEditToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, existing []*database.Message, data *types.MessageData) (*bridgev2.ConvertedEdit, error) {
	data.Text = strings.TrimPrefix(data.Text, "Edited: ")
	return &bridgev2.ConvertedEdit{
		ModifiedParts: []*bridgev2.ConvertedEditPart{tc.convertToMatrix(ctx, portal, intent, data).Parts[0].ToEditPart(existing[0])},
	}, nil
}

func (tc *TwitterClient) convertToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *types.MessageData) *bridgev2.ConvertedMessage {
	var replyTo *networkid.MessageOptionalPartID
	if msg.ReplyData.ID != "" {
		replyTo = &networkid.MessageOptionalPartID{
			MessageID: networkid.MessageID(msg.ReplyData.ID),
		}
	}

	textPart := &bridgev2.ConvertedMessagePart{
		ID:      "",
		Type:    bridgeEvt.EventMessage,
		Content: twitterfmt.Parse(ctx, portal, msg),
	}

	parts := make([]*bridgev2.ConvertedMessagePart, 0)

	if msg.Attachment != nil {
		convertedAttachmentPart, indices, err := tc.TwitterAttachmentToMatrix(ctx, portal, intent, msg)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to convert attachment")
			parts = append(parts, &bridgev2.ConvertedMessagePart{
				ID:   "",
				Type: bridgeEvt.EventMessage,
				Content: &bridgeEvt.MessageEventContent{
					MsgType: bridgeEvt.MsgNotice,
					Body:    "Failed to convert attachment from Twitter",
				},
			})
		} else {
			if msg.Attachment.Card != nil || msg.Attachment.Tweet != nil {
				textPart.Content.BeeperLinkPreviews = convertedAttachmentPart.Content.BeeperLinkPreviews
			} else {
				parts = append(parts, convertedAttachmentPart)
				RemoveEntityLinkFromText(textPart, indices)
			}
		}
	}

	if len(textPart.Content.Body) > 0 {
		parts = append(parts, textPart)
	}

	cm := &bridgev2.ConvertedMessage{
		ReplyTo: replyTo,
		Parts:   parts,
	}
	cm.MergeCaption()

	return cm
}

func (tc *TwitterClient) MakePortalKey(conv types.Conversation) networkid.PortalKey {
	var receiver networkid.UserLoginID
	if conv.Type == types.ONE_TO_ONE || tc.connector.br.Config.SplitPortals {
		receiver = tc.userLogin.ID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conv.ConversationID),
		Receiver: receiver,
	}
}

func (tc *TwitterClient) MakePortalKeyFromID(conversationID string) networkid.PortalKey {
	var receiver networkid.UserLoginID
	if strings.Contains(conversationID, "-") || tc.connector.br.Config.SplitPortals {
		receiver = tc.userLogin.ID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conversationID),
		Receiver: receiver,
	}
}
