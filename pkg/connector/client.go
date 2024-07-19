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

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

type TwitterClient struct {
}

var _ bridgev2.NetworkAPI = (*TwitterClient)(nil)

func (tc *TwitterClient) Connect(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *bridgev2.NetworkRoomCapabilities {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	//TODO implement me
	panic("implement me")
}
