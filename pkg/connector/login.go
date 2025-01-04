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

	"maunium.net/go/mautrix/bridge/status"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

type TwitterLogin struct {
	User    *bridgev2.User
	Cookies string
	tc      *TwitterConnector
}

var (
	LoginStepIDCookies  = "fi.mau.twitter.login.enter_cookies"
	LoginStepIDComplete = "fi.mau.twitter.login.complete"
)

var _ bridgev2.LoginProcessCookies = (*TwitterLogin)(nil)

func (tc *TwitterConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{
		{
			Name:        "Cookies",
			Description: "Log in with your Twitter account using your cookies",
			ID:          "cookies",
		},
	}
}

func (tc *TwitterConnector) CreateLogin(_ context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != "cookies" {
		return nil, fmt.Errorf("unknown login flow ID: %s", flowID)
	}
	return &TwitterLogin{User: user, tc: tc}, nil
}

func (t *TwitterLogin) Start(_ context.Context) (*bridgev2.LoginStep, error) {
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       LoginStepIDCookies,
		Instructions: "Open the Login URL in an Incognito/Private browsing mode. Then, extract the cookies as a JSON object/cURL command copied from the Network tab of your browser's DevTools. After that, close the browser **before** pasting the cookies.\n\nFor example: `{\"ct0\":\"123466-...\",\"auth_token\":\"abcde-...\"}`",
		CookiesParams: &bridgev2.LoginCookiesParams{
			URL:       "https://x.com",
			UserAgent: "",
			Fields: []bridgev2.LoginCookieField{
				{
					ID:       "ct0",
					Required: true,
					Sources: []bridgev2.LoginCookieFieldSource{
						{Type: bridgev2.LoginCookieTypeCookie, Name: "ct0"},
					},
				},
				{
					ID:       "auth_token",
					Required: true,
					Sources: []bridgev2.LoginCookieFieldSource{
						{Type: bridgev2.LoginCookieTypeCookie, Name: "auth_token"},
					},
				},
			},
		},
	}, nil
}

func (t *TwitterLogin) Cancel() {}

func (t *TwitterLogin) SubmitCookies(ctx context.Context, cookies map[string]string) (*bridgev2.LoginStep, error) {
	cookieStruct := twitCookies.NewCookies(cookies)
	meta := &UserLoginMetadata{
		Cookies: cookieStruct.String(),
	}

	clientOpts := &twittermeow.ClientOpts{
		Cookies:       cookieStruct,
		WithJOTClient: true,
	}
	client := twittermeow.NewClient(clientOpts, t.User.Log.With().Str("component", "login_twitter_client").Logger())

	inboxState, settings, err := client.LoadMessagesPage()
	if err != nil {
		return nil, fmt.Errorf("failed to load messages page after submitting cookies: %w", err)
	}
	selfUser := inboxState.InboxInitialState.GetUserByID(client.GetCurrentUserID())

	id := networkid.UserLoginID(client.GetCurrentUserID())
	ul, err := t.User.NewLogin(
		ctx,
		&database.UserLogin{
			ID:         id,
			Metadata:   meta,
			RemoteName: settings.ScreenName,
			RemoteProfile: status.RemoteProfile{
				Username: settings.ScreenName,
				Name:     selfUser.Name,
			},
		},
		&bridgev2.NewLoginParams{
			DeleteOnConflict:  true,
			DontReuseExisting: false,
			LoadUserLogin: func(ctx context.Context, login *bridgev2.UserLogin) error {
				client.Logger = login.Log.With().Str("component", "twitter_client").Logger()
				login.Client = NewTwitterClient(login, t.tc, client)
				return nil
			},
		},
	)
	if err != nil {
		return nil, err
	}

	go func(ctx context.Context, client *TwitterClient, inboxState *response.InboxInitialStateResponse) {
		client.syncChannels(ctx, inboxState)
		client.startPolling(ctx)
	}(context.WithoutCancel(ctx), ul.Client.(*TwitterClient), inboxState)

	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeComplete,
		StepID:       LoginStepIDComplete,
		Instructions: fmt.Sprintf("Successfully logged into @%s", ul.UserLogin.RemoteName),
		CompleteParams: &bridgev2.LoginCompleteParams{
			UserLoginID: ul.ID,
			UserLogin:   ul,
		},
	}, nil
}
