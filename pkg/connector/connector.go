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

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/id"
)

type TwitterConnector struct {
	br *bridgev2.Bridge

	Config      Config
	directMedia bool
}

var _ bridgev2.NetworkConnector = (*TwitterConnector)(nil)

func (tc *TwitterConnector) Init(bridge *bridgev2.Bridge) {
	tc.br = bridge
	// XChat checkpoints must reflect completed Matrix handling, not admission to
	// an in-memory portal queue. Synchronous portal handling preserves per-room
	// ordering while making QueueRemoteEvent's result authoritative.
	bridgev2.PortalEventBuffer = 0
}

func xchatRemoteEventHandled(result bridgev2.EventHandlingResult) bool {
	return result.Success && !result.Queued
}

func (tc *TwitterConnector) Start(_ context.Context) error {
	return nil
}

func (tc *TwitterConnector) GetName() bridgev2.BridgeName {
	networkIcon := "mxc://maunium.net/HVHcnusJkQcpVcsVGZRELLCn"
	displayName := "Twitter"
	if tc.Config.X {
		displayName = "X"
		networkIcon = "mxc://maunium.net/kbYtMeRPvXehbHXZtvOnUcyu"
	}
	return bridgev2.BridgeName{
		DisplayName:      displayName,
		NetworkURL:       "https://x.com",
		NetworkIcon:      id.ContentURIString(networkIcon),
		NetworkID:        "twitter",
		BeeperBridgeType: "twitter",
		DefaultPort:      29327,
	}
}
