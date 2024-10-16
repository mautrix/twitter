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

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
)

type TwitterLogin struct {
	User    *bridgev2.User
	Cookies string
}

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
	return &TwitterLogin{User: user}, nil
}

func (t *TwitterLogin) Start(_ context.Context) (*bridgev2.LoginStep, error) {
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeCookies,
		StepID:       "fi.mau.twitter.login.enter_cookies",
		Instructions: "Enter a JSON object with your cookies, or a cURL command copied from browser devtools.\n\nFor example: `{\"ct0\":\"123466-...\",\"auth_token\":\"abcde-...\"}`",
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
				/*
						{
							ID: "guest_id",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "guest_id" },
							},
						},
						{
							ID: "twid",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "twid" },
							},
						},
						{
							ID: "kdt",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "kdt" },
							},
						},
						{
							ID: "night_mode",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "night_mode" },
							},
						},
						{
							ID: "personalization_id",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "personalization_id" },
							},
						},
						{
							ID: "guest_id_marketing",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "guest_id_marketing" },
							},
						},
						{
							ID: "guest_id_ads",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "guest_id_ads" },
							},
						},
						{
							ID: "d_prefs",
							Required: false,
							Sources: []bridgev2.LoginCookieFieldSource{
								{ Type: bridgev2.LoginCookieTypeCookie, Name: "d_prefs" },
							},
						},
					},
				*/
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
	client := twittermeow.NewClient(clientOpts, t.User.Log)

	_, _, err := client.LoadMessagesPage()
	if err != nil {
		return nil, fmt.Errorf("failed to load messages page after submitting cookies")
	}

	id := networkid.UserLoginID(client.GetCurrentUserID())
	ul, err := t.User.NewLogin(
		ctx,
		&database.UserLogin{
			ID:       id,
			Metadata: meta,
		},
		&bridgev2.NewLoginParams{
			DeleteOnConflict:  true,
			DontReuseExisting: false,
		},
	)
	if err != nil {
		return nil, err
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeComplete,
		StepID:       "fi.mau.twitter.login.complete",
		Instructions: fmt.Sprintf("Successfully logged into @%s", ul.UserLogin.RemoteName),
		CompleteParams: &bridgev2.LoginCompleteParams{
			UserLoginID: ul.ID,
			UserLogin:   ul,
		},
	}, nil
}
