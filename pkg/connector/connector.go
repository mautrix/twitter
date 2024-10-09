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
	"log"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

type TwitterConnector struct {
	br *bridgev2.Bridge

	Config *TwitterConfig
}

var _ bridgev2.NetworkConnector = (*TwitterConnector)(nil)

func NewConnector() *TwitterConnector {
	return &TwitterConnector{}
}

func (tc *TwitterConnector) Init(bridge *bridgev2.Bridge) {
	tc.br = bridge
}

func (tc *TwitterConnector) Start(_ context.Context) error {

	log.Println("starting....")
	return nil
}

func (tc *TwitterConnector) GetName() bridgev2.BridgeName {
	return bridgev2.BridgeName{
		DisplayName:      "Twitter",
		NetworkURL:       "https://twitter.com",
		NetworkIcon:      "mxc://maunium.net/HVHcnusJkQcpVcsVGZRELLCn",
		NetworkID:        "twitter",
		BeeperBridgeType: "twitter",
		DefaultPort:      29327,
	}
}

func (tc *TwitterConnector) GetDBMetaTypes() database.MetaTypes {
	return database.MetaTypes{
		Reaction: nil,
		Portal:   nil,
		Message:  nil,
		Ghost:    nil,
		UserLogin: func() any {
			return &UserLoginMetadata{}
		},
	}
}

func (tc *TwitterConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	return &bridgev2.NetworkGeneralCapabilities{}
}

type UserLoginMetadata struct {
	Cookies string
}

func (tc *TwitterConnector) LoadUserLogin(_ context.Context, login *bridgev2.UserLogin) error {
	meta := login.Metadata.(*UserLoginMetadata)
	clientOpts := &twittermeow.ClientOpts{
		Cookies:       twitCookies.NewCookiesFromString(meta.Cookies),
		WithJOTClient: true,
	}
	twitClient := &TwitterClient{
		connector: tc,
		userLogin: login,
		client:    twittermeow.NewClient(clientOpts, login.Log),
		userCache: make(map[string]types.User),
	}
	twitClient.client.SetEventHandler(twitClient.HandleTwitterEvent)

	_, currentUser, err := twitClient.client.LoadMessagesPage()
	if err != nil {
		return fmt.Errorf("failed to load messages page")
	}

	login.RemoteName = currentUser.ScreenName
	login.Client = twitClient

	return nil
}
