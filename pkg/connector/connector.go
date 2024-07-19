// mautrix-twitter - A Matrix-Slack puppeting bridge.
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

	"go.mau.fi/util/configupgrade"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
)

type TwitterConnector struct{}

var _ bridgev2.NetworkConnector = (*TwitterConnector)(nil)

func (tc *TwitterConnector) Init(bridge *bridgev2.Bridge) {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterConnector) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterConnector) GetConfig() (example string, data any, upgrader configupgrade.Upgrader) {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	//TODO implement me
	panic("implement me")
}
