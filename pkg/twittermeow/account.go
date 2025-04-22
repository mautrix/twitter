package twittermeow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (c *Client) Login(ctx context.Context) error {
	err := c.session.LoadPage(ctx, endpoints.BASE_LOGIN_URL)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccountSettings(ctx context.Context, params payload.AccountSettingsQuery) (*response.AccountSettingsResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.ACCOUNT_SETTINGS_URL, string(encodedQuery))
	apiRequestOpts := apiRequestOpts{
		URL:    url,
		Method: http.MethodGet,
	}
	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.AccountSettingsResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) GetDMPermissions(ctx context.Context, params payload.GetDMPermissionsQuery) (*response.GetDMPermissionsResponse, error) {
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
	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts)
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

	Settings *payload.PushNotificationSettings
}

type PushNotificationConfigAction int

const (
	PushRegister PushNotificationConfigAction = iota
	PushUnregister
	PushCheckin
	PushSave
)

func (c *Client) SetPushNotificationConfig(ctx context.Context, action PushNotificationConfigAction, config WebPushConfig) error {
	var url string
	switch action {
	case PushRegister:
		url = endpoints.NOTIFICATION_LOGIN_URL
	case PushUnregister:
		url = endpoints.NOTIFICATION_LOGOUT_URL
	case PushCheckin:
		url = endpoints.NOTIFICATION_CHECKIN_URL
	case PushSave:
		url = endpoints.NOTIFICATION_SAVE_URL
	default:
		return fmt.Errorf("unknown push notification setting: %d", action)
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

		Settings: config.Settings,
	}

	var wrappedPayload any
	if action != PushUnregister {
		wrappedPayload = &payload.PushConfigPayloadWrapper{
			PushDeviceInfo: &webPushPayload,
		}
	} else {
		wrappedPayload = &webPushPayload
	}

	encodedBody, err := json.Marshal(wrappedPayload)
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
		ContentType:    types.ContentTypeJSON,
	}
	_, _, err = c.makeAPIRequest(ctx, apiRequestOpts)
	return err
}
