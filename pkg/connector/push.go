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
	"fmt"

	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var _ bridgev2.PushableNetworkAPI = (*TwitterClient)(nil)

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
