package twittermeow

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"
)

func (c *Client) Login() error {
	err := c.session.LoadPage(endpoints.BASE_LOGIN_URL)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccountSettings(params payload.AccountSettingsQuery) (*response.AccountSettingsResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.ACCOUNT_SETTINGS_URL, string(encodedQuery))
	apiRequestOpts := apiRequestOpts{
		URL:    url,
		Method: http.MethodGet,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.AccountSettingsResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) GetDMPermissions(params payload.GetDMPermissionsQuery) (*response.GetDMPermissionsResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.DM_PERMISSIONS_URL, string(encodedQuery))
	apiRequestOpts := apiRequestOpts{
		URL:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.GetDMPermissionsResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

type WebPushConfig struct {
	Endpoint string
	P256DH   []byte
	Auth     []byte
}

type PushNotificationSetting int

const (
	REGISTER_PUSH   PushNotificationSetting = 0
	UNREGISTER_PUSH PushNotificationSetting = 1
)

func (c *Client) SetPushNotificationConfig(setting PushNotificationSetting, config WebPushConfig) error {
	var url string
	switch setting {
	case REGISTER_PUSH:
		url = endpoints.NOTIFICATION_LOGIN_URL
	case UNREGISTER_PUSH:
		url = endpoints.NOTIFICATION_LOGOUT_URL
	default:
		return fmt.Errorf("unknown push notification setting: %d", setting)
	}

	webPushPayload := payload.WebPushConfigPayload{
		Env:             3,
		ProtocolVersion: 1,

		Locale:    "en",
		OSVersion: UDID,
		UDID:      UDID,

		Token:  config.Endpoint,
		P256DH: base64.RawURLEncoding.EncodeToString(config.P256DH),
		Auth:   base64.RawURLEncoding.EncodeToString(config.Auth),
	}

	encodedBody, err := json.Marshal(webPushPayload)

	if err != nil {
		return err
	}

	apiRequestOpts := apiRequestOpts{
		URL:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Referer:        endpoints.BASE_NOTIFICATION_SETTINGS_URL,
		Origin:         endpoints.BASE_URL,
		Body:           encodedBody,
		ContentType:    types.JSON,
	}
	_, _, err = c.makeAPIRequest(apiRequestOpts)
	return err
}
