package twittermeow

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"
)

type JotClient struct {
	client *Client
}

func (c *Client) newJotClient() *JotClient {
	return &JotClient{
		client: c,
	}
}

func (jc *JotClient) sendClientLoggingEvent(category payload.JotLoggingCategory, debug bool, body []interface{}) error {
	logPayloadBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	clientLogPayload := &payload.JotClientEventPayload{
		Category: category,
		Debug:    debug,
		Log:      string(logPayloadBytes),
	}

	clientLogPayloadBytes, err := clientLogPayload.Encode()
	if err != nil {
		return err
	}

	clientTransactionId, err := crypto.SignTransaction(jc.client.session.verificationToken, endpoints.JOT_CLIENT_EVENT_URL, http.MethodPost)
	if err != nil {
		return err
	}

	extraHeaders := map[string]string{
		"accept":                  "*/*",
		"sec-fetch-site":          "same-site",
		"sec-fetch-mode":          "cors",
		"sec-fetch-dest":          "empty",
		"x-client-transaction-id": clientTransactionId,
	}

	headerOpts := HeaderOpts{
		WithAuthBearer:      true,
		WithCookies:         true,
		WithXGuestToken:     true,
		WithXTwitterHeaders: true,
		Origin:              endpoints.BASE_URL,
		Referer:             endpoints.BASE_URL + "/",
		Extra:               extraHeaders,
	}

	clientLogResponse, _, err := jc.client.MakeRequest(endpoints.JOT_CLIENT_EVENT_URL, http.MethodPost, jc.client.buildHeaders(headerOpts), clientLogPayloadBytes, types.FORM)
	if err != nil {
		return err
	}

	if clientLogResponse.StatusCode > 204 {
		return fmt.Errorf("failed to send jot client event, status code: %d", clientLogResponse.StatusCode)
	}

	return nil
}
