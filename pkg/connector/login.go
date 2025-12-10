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

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	twitCookies "go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

type TwitterLogin struct {
	User       *bridgev2.User
	Cookies    string
	SecretKey  string
	SigningKey string
	tc         *TwitterConnector

	client   *twittermeow.Client
	settings *response.AccountSettingsResponse
}

var (
	LoginStepIDCookies  = "fi.mau.twitter.login.enter_cookies"
	LoginStepSecretKey  = "fi.mau.twitter.login.secret_key"
	LoginStepIDComplete = "fi.mau.twitter.login.complete"
)

var _ bridgev2.LoginProcessCookies = (*TwitterLogin)(nil)
var _ bridgev2.LoginProcessUserInput = (*TwitterLogin)(nil)

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
			URL:       "https://x.com/i/flow/login",
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
	t.Cookies = cookieStruct.String()

	client := twittermeow.NewClient(cookieStruct, t.User.Log.With().Str("component", "login_twitter_client").Logger())

	settings, err := client.LoadMessagesPage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load messages page after submitting cookies: %w", err)
	}
	t.client = client
	t.settings = settings

	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       LoginStepSecretKey,
		Instructions: "Enter the secret key (for decrypting conversation keys) and signing key (for signing messages).",
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type:        bridgev2.LoginInputFieldTypePassword,
					ID:          "secret_key",
					Name:        "Secret key",
					Description: "The secret key for decrypting conversation keys (base64 32-byte scalar).",
					Validate: func(input string) (string, error) {
						input = strings.TrimSpace(input)
						if input == "" {
							return "", fmt.Errorf("secret key cannot be empty")
						}
						if _, err := crypto.ParsePrivateKeyScalar(input); err != nil {
							return "", fmt.Errorf("invalid secret key: %w", err)
						}
						return input, nil
					},
				},
				{
					Type:        bridgev2.LoginInputFieldTypePassword,
					ID:          "signing_key",
					Name:        "Signing key",
					Description: "The signing key for signing messages (base64 32-byte scalar).",
					Validate: func(input string) (string, error) {
						input = strings.TrimSpace(input)
						if input == "" {
							return "", nil
						}
						if _, err := crypto.ParsePrivateKeyScalar(input); err != nil {
							return "", fmt.Errorf("invalid signing key: %w", err)
						}
						return input, nil
					},
				},
			},
		},
	}, nil
}

func (t *TwitterLogin) SubmitUserInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	secretKey, ok := input["secret_key"]
	if !ok {
		return nil, fmt.Errorf("secret_key input is required")
	}
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" {
		return nil, fmt.Errorf("secret key cannot be empty")
	}
	if _, err := crypto.ParsePrivateKeyScalar(secretKey); err != nil {
		return nil, fmt.Errorf("invalid secret key: %w", err)
	}
	t.SecretKey = secretKey

	signingKey := strings.TrimSpace(input["signing_key"])
	if signingKey != "" {
		if _, err := crypto.ParsePrivateKeyScalar(signingKey); err != nil {
			return nil, fmt.Errorf("invalid signing key: %w", err)
		}
	}
	t.SigningKey = signingKey

	if t.client == nil {
		if t.Cookies == "" {
			return nil, fmt.Errorf("cookies must be submitted before secret key")
		}
		cookieStruct := twitCookies.NewCookiesFromString(t.Cookies)
		t.client = twittermeow.NewClient(cookieStruct, t.User.Log.With().Str("component", "login_twitter_client").Logger())
		settings, err := t.client.LoadMessagesPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load messages page after submitting secret key: %w", err)
		}
		t.settings = settings
	}

	meta := &UserLoginMetadata{
		Cookies:    t.Cookies,
		SecretKey:  t.SecretKey,
		SigningKey: t.SigningKey,
	}

	remoteProfile := &status.RemoteProfile{
		Username: t.settings.ScreenName,
	}
	id := MakeUserLoginID(t.client.GetCurrentUserID())
	ul, err := t.User.NewLogin(
		ctx,
		&database.UserLogin{
			ID:            id,
			Metadata:      meta,
			RemoteName:    remoteProfile.Username,
			RemoteProfile: *remoteProfile,
		},
		&bridgev2.NewLoginParams{
			DeleteOnConflict:  true,
			DontReuseExisting: false,
			LoadUserLogin: func(ctx context.Context, login *bridgev2.UserLogin) error {
				ensureUserLoginMetadata(login)
				if t.client != nil {
					t.client.SetKeyStore(newUserLoginKeyStore(login))
				}
				t.client.Logger = login.Log.With().Str("component", "twitter_client").Logger()
				login.Client = NewTwitterClient(login, t.tc, t.client)
				return nil
			},
		},
	)
	if err != nil {
		return nil, err
	}
	ul.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})

	go func(ctx context.Context, client *TwitterClient) {
		client.DoConnect(ctx)
	}(context.WithoutCancel(ctx), ul.Client.(*TwitterClient))

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
