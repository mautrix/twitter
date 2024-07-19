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
)

type TwitterLogin struct {
}

var _ bridgev2.LoginProcessCookies = (*TwitterLogin)(nil)

func (tc *TwitterConnector) GetLoginFlows() []bridgev2.LoginFlow {
	//TODO implement me
	panic("implement me")
}

func (tc *TwitterConnector) CreateLogin(ctx context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TwitterLogin) Start(ctx context.Context) (*bridgev2.LoginStep, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TwitterLogin) Cancel() {
	//TODO implement me
	panic("implement me")
}

func (t *TwitterLogin) SubmitCookies(ctx context.Context, cookies map[string]string) (*bridgev2.LoginStep, error) {
	//TODO implement me
	panic("implement me")
}
