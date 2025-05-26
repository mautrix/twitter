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
	"sync"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/status"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type TwitterClient struct {
	connector *TwitterConnector
	client    *twittermeow.Client

	userLogin *bridgev2.UserLogin

	userCache     map[string]*types.User
	userCacheLock sync.RWMutex

	participantCache map[string][]types.Participant

	matrixParser *format.HTMLParser
}

var _ bridgev2.NetworkAPI = (*TwitterClient)(nil)

func NewTwitterClient(login *bridgev2.UserLogin, connector *TwitterConnector, client *twittermeow.Client) *TwitterClient {
	tc := &TwitterClient{
		connector:        connector,
		client:           client,
		userLogin:        login,
		userCache:        make(map[string]*types.User),
		participantCache: make(map[string][]types.Participant),
	}
	client.SetEventHandler(tc.HandleTwitterEvent, tc.HandleStreamEvent)
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
			return "@" + identifierToScreenName(id)
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

func (tc *TwitterClient) Connect(ctx context.Context) {
	if tc.client == nil {
		tc.userLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateBadCredentials,
			Error:      "twitter-not-logged-in",
		})
		return
	}

	inboxState, _, err := tc.client.LoadMessagesPage(ctx)
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

	currentUserID := tc.client.GetCurrentUserID()
	if MakeUserLoginID(currentUserID) != tc.userLogin.ID {
		zerolog.Ctx(ctx).Warn().
			Str("user_login_id", string(tc.userLogin.ID)).
			Str("current_user_id", currentUserID).
			Msg("User login ID mismatch")
	}

	remoteProfile := tc.connector.makeRemoteProfile(ctx, tc.client, currentUserID, inboxState.InboxInitialState)
	if remoteProfile != nil && (tc.userLogin.RemoteName != remoteProfile.Username ||
		tc.userLogin.RemoteProfile != *remoteProfile) {
		tc.userLogin.RemoteName = remoteProfile.Username
		tc.userLogin.RemoteProfile = *remoteProfile
		err = tc.userLogin.Save(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to save user login after updating remote profile")
		}
	}

	go tc.syncChannels(ctx, inboxState.InboxInitialState)
	tc.startPolling(ctx)
}

func (tc *TwitterConnector) makeRemoteProfile(ctx context.Context, cli *twittermeow.Client, currentUserID string, inbox *response.TwitterInboxData) *status.RemoteProfile {
	selfUser := inbox.GetUserByID(currentUserID)
	if selfUser == nil {
		zerolog.Ctx(ctx).Warn().Msg("Own user info not found in inbox state")
		return nil
	}
	var avatarMXC id.ContentURIString
	ownGhost, err := tc.br.GetGhostByID(ctx, MakeUserID(currentUserID))
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get own ghost by ID")
	} else {
		ownGhost.UpdateInfo(ctx, tc.wrapUserInfo(cli, selfUser))
		avatarMXC = ownGhost.AvatarMXC
	}
	return &status.RemoteProfile{
		// TODO fetch from /1.1/users/email_phone_info.json?
		Phone:    "",
		Email:    "",
		Username: selfUser.ScreenName,
		Name:     selfUser.Name,
		Avatar:   avatarMXC,
	}
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
	tc.client.Disconnect()
}

func (tc *TwitterClient) IsLoggedIn() bool {
	return tc.client.IsLoggedIn()
}

func (tc *TwitterClient) LogoutRemote(ctx context.Context) {
	log := zerolog.Ctx(ctx)
	_, err := tc.client.Logout(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error logging out")
	}
}

func (tc *TwitterClient) IsThisUser(_ context.Context, userID networkid.UserID) bool {
	return UserLoginIDToUserID(tc.userLogin.ID) == userID
}
