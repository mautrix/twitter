package twittermeow

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"

	"github.com/google/uuid"
)

func (c *Client) GetInitialInboxState(params *payload.DmRequestQuery) (*response.InboxInitialStateResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.INBOX_INITIAL_STATE_URL, string(encodedQuery))

	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.InboxInitialStateResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) GetDMUserUpdates(params *payload.DmRequestQuery) (*response.GetDMUserUpdatesResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.DM_USER_UPDATES_URL, string(encodedQuery))

	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.GetDMUserUpdatesResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) MarkConversationRead(params *payload.MarkConversationReadQuery) error {
	encodedQueryBody, err := params.Encode()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(endpoints.CONVERSATION_MARK_READ_URL, params.ConversationID)
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           encodedQueryBody,
		ContentType:    types.FORM,
	}
	resp, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		c.Logger.Warn().Any("response_body", string(respBody)).Any("status_code", resp.StatusCode).Any("params", params).Msg("Failed to mark conversation as read")
		return ErrFailedMarkConversationRead
	}

	return nil
}

func (c *Client) FetchConversationContext(conversationId string, params *payload.DmRequestQuery, context payload.ContextInfo) (*response.ConversationDMResponse, error) {
	params.Context = context
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", fmt.Sprintf(endpoints.CONVERSATION_FETCH_MESSAGES, conversationId), string(encodedQuery))

	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.ConversationDMResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) FetchTrustedThreads(params *payload.DmRequestQuery) (*response.InboxTimelineResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.TRUSTED_INBOX_TIMELINE_URL, string(encodedQuery))

	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.InboxTimelineResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) SendDirectMessage(payload *payload.SendDirectMessagePayload) (*response.XInboxData, error) {
	if payload.RequestID == "" {
		payload.RequestID = uuid.NewString()
	}

	jsonBody, err := payload.Encode()
	if err != nil {
		return nil, err
	}

	url := endpoints.SEND_DM_URL
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Referer:        fmt.Sprintf("%s/%s", endpoints.BASE_MESSAGES_URL, payload.ConversationID),
		Origin:         endpoints.BASE_URL,
		ContentType:    types.JSON,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.XInboxData{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) EditDirectMessage(payload *payload.EditDirectMessagePayload) (*types.Message, error) {
	if payload.RequestID == "" {
		payload.RequestID = uuid.NewString()
	}

	encodedQuery, err := payload.Encode()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?%s", endpoints.EDIT_DM_URL, string(encodedQuery))
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Referer:        fmt.Sprintf("%s/%s", endpoints.BASE_MESSAGES_URL, payload.ConversationID),
		Origin:         endpoints.BASE_URL,
		ContentType:    types.FORM,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := types.Message{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) SendTypingNotification(conversationId string) error {
	variables := &payload.SendTypingNotificationVariables{
		ConversationID: conversationId,
	}

	GQLPayload := &payload.GraphQLPayload{
		Variables: variables,
		QueryID:   "HL96-xZ3Y81IEzAdczDokg",
	}

	jsonBody, err := GQLPayload.Encode()
	if err != nil {
		return err
	}

	apiRequestOpts := apiRequestOpts{
		Url:            endpoints.SEND_TYPING_NOTIFICATION,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Referer:        fmt.Sprintf("%s/%s", endpoints.BASE_MESSAGES_URL, conversationId),
		Origin:         endpoints.BASE_URL,
		ContentType:    types.JSON,
		Body:           jsonBody,
	}
	_, resp, err := c.makeAPIRequest(apiRequestOpts)
	fmt.Println(resp)
	return err
}

// keep in mind this only deletes the message for you
func (c *Client) DeleteMessageForMe(variables *payload.DMMessageDeleteMutationVariables) (*response.DMMessageDeleteMutationResponse, error) {
	if variables.RequestID == "" {
		variables.RequestID = uuid.NewString()
	}

	GQLPayload := &payload.GraphQLPayload{
		Variables: variables,
		QueryID:   "BJ6DtxA2llfjnRoRjaiIiw",
	}

	jsonBody, err := GQLPayload.Encode()
	if err != nil {
		return nil, err
	}

	url := endpoints.GRAPHQL_MESSAGE_DELETION_MUTATION
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.JSON,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.DMMessageDeleteMutationResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) DeleteConversation(conversationId string, payload *payload.DmRequestQuery) error {
	encodedQueryBody, err := payload.Encode()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(endpoints.DELETE_CONVERSATION_URL, conversationId)
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           encodedQueryBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.FORM,
	}
	resp, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("failed to delete conversation by id %s (status_code=%d, response_body=%s)", conversationId, resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) PinConversation(conversationId string) (*response.PinConversationResponse, error) {
	graphQlPayload := payload.GraphQLPayload{
		Variables: payload.PinAndUnpinConversationVariables{
			ConversationID: conversationId,
			Label:          payload.LABEL_TYPE_PINNED,
		},
		QueryID: "o0aymgGiJY-53Y52YSUGVA",
	}

	jsonBody, err := graphQlPayload.Encode()
	if err != nil {
		return nil, err
	}

	url := endpoints.PIN_CONVERSATION_URL
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.JSON,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.PinConversationResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) UnpinConversation(conversationId string) (*response.UnpinConversationResponse, error) {
	graphQlPayload := payload.GraphQLPayload{
		Variables: payload.PinAndUnpinConversationVariables{
			ConversationID: conversationId,
			LabelType:      payload.LABEL_TYPE_PINNED,
		},
		QueryID: "_TQxP2Rb0expwVP9ktGrTQ",
	}

	jsonBody, err := graphQlPayload.Encode()
	if err != nil {
		return nil, err
	}

	url := endpoints.UNPIN_CONVERSATION_URL
	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.JSON,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.UnpinConversationResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) React(reactionPayload *payload.ReactionActionPayload, remove bool) (*response.ReactionResponse, error) {
	graphQlPayload := payload.GraphQLPayload{
		Variables: reactionPayload,
		QueryID:   "VyDyV9pC2oZEj6g52hgnhA",
	}

	url := endpoints.ADD_REACTION_URL
	if remove {
		url = endpoints.REMOVE_REACTION_URL
		graphQlPayload.QueryID = "bV_Nim3RYHsaJwMkTXJ6ew"
	}

	jsonBody, err := graphQlPayload.Encode()
	if err != nil {
		return nil, err
	}

	apiRequestOpts := apiRequestOpts{
		Url:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.JSON,
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.ReactionResponse{}
	return &data, json.Unmarshal(respBody, &data)
}
