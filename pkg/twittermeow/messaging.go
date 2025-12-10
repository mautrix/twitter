package twittermeow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (c *Client) GetInitialInboxState(ctx context.Context, params *payload.DMRequestQuery) (*response.InboxInitialStateResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.INBOX_INITIAL_STATE_URL, string(encodedQuery))

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	})
	if err != nil {
		return nil, err
	}

	data := response.InboxInitialStateResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) GetDMUserUpdates(ctx context.Context, params *payload.DMRequestQuery) (*response.GetDMUserUpdatesResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.DM_USER_UPDATES_URL, string(encodedQuery))

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	})
	if err != nil {
		return nil, err
	}

	data := response.GetDMUserUpdatesResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) MarkConversationRead(ctx context.Context, params *payload.MarkConversationReadQuery) error {
	encodedQueryBody, err := params.Encode()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(endpoints.CONVERSATION_MARK_READ_URL, params.ConversationID)
	_, _, err = c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           encodedQueryBody,
		ContentType:    types.ContentTypeForm,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) FetchConversationContext(ctx context.Context, conversationID string, params *payload.DMRequestQuery, context payload.ContextInfo) (*response.ConversationDMResponse, error) {
	params.Context = context
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", fmt.Sprintf(endpoints.CONVERSATION_FETCH_MESSAGES, conversationID), string(encodedQuery))

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	})
	if err != nil {
		return nil, err
	}

	data := response.ConversationDMResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) FetchTrustedThreads(ctx context.Context, params *payload.DMRequestQuery) (*response.InboxTimelineResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            fmt.Sprintf("%s?%s", endpoints.TRUSTED_INBOX_TIMELINE_URL, string(encodedQuery)),
		Method:         http.MethodGet,
		WithClientUUID: true,
	})
	if err != nil {
		return nil, err
	}

	data := response.InboxTimelineResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) SendDirectMessage(ctx context.Context, pl *payload.SendDirectMessagePayload) (*response.TwitterInboxData, error) {
	if pl.RequestID == "" {
		pl.RequestID = uuid.NewString()
	}

	jsonBody, err := pl.Encode()
	if err != nil {
		return nil, err
	}

	query, _ := (payload.DMSendQuery{}).Default().Encode()
	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.SEND_DM_URL + "?" + string(query),
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Referer:        fmt.Sprintf("%s/%s", endpoints.BASE_MESSAGES_URL, pl.ConversationID),
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	data := response.TwitterInboxData{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) EditDirectMessage(ctx context.Context, payload *payload.EditDirectMessagePayload) (*types.Message, error) {
	if payload.RequestID == "" {
		payload.RequestID = uuid.NewString()
	}

	encodedQuery, err := payload.Encode()
	if err != nil {
		return nil, err
	}

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            fmt.Sprintf("%s?%s", endpoints.EDIT_DM_URL, string(encodedQuery)),
		Method:         http.MethodPost,
		WithClientUUID: true,
		Referer:        fmt.Sprintf("%s/%s", endpoints.BASE_MESSAGES_URL, payload.ConversationID),
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeForm,
	})
	if err != nil {
		return nil, err
	}

	data := types.Message{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) SendTypingNotification(ctx context.Context, conversationID string) error {
	variables := &payload.SendTypingNotificationVariables{
		ConversationID: conversationID,
	}

	GQLPayload := &payload.GraphQLPayload{
		Variables: variables,
		QueryID:   "HL96-xZ3Y81IEzAdczDokg",
	}

	jsonBody, err := GQLPayload.Encode()
	if err != nil {
		return err
	}

	_, _, err = c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.SEND_TYPING_NOTIFICATION,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Referer:        fmt.Sprintf("%s/%s", endpoints.BASE_MESSAGES_URL, conversationID),
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	return err
}

// keep in mind this only deletes the message for you
func (c *Client) DeleteMessageForMe(ctx context.Context, variables *payload.DMMessageDeleteMutationVariables) (*response.DMMessageDeleteMutationResponse, error) {
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

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.GRAPHQL_MESSAGE_DELETION_MUTATION,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	data := response.DMMessageDeleteMutationResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) DeleteConversation(ctx context.Context, conversationID string, payload *payload.DMRequestQuery) error {
	encodedQueryBody, err := payload.Encode()
	if err != nil {
		return err
	}

	resp, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            fmt.Sprintf(endpoints.DELETE_CONVERSATION_URL, conversationID),
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           encodedQueryBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeForm,
	})
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("failed to delete conversation by id %s (status_code=%d, response_body=%s)", conversationID, resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) PinConversation(ctx context.Context, conversationID string) (*response.PinConversationResponse, error) {
	graphQlPayload := payload.GraphQLPayload{
		Variables: payload.PinAndUnpinConversationVariables{
			ConversationID: conversationID,
			Label:          payload.LABEL_TYPE_PINNED,
		},
		QueryID: "o0aymgGiJY-53Y52YSUGVA",
	}

	jsonBody, err := graphQlPayload.Encode()
	if err != nil {
		return nil, err
	}

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.PIN_CONVERSATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	data := response.PinConversationResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) UnpinConversation(ctx context.Context, conversationID string) (*response.UnpinConversationResponse, error) {
	graphQlPayload := payload.GraphQLPayload{
		Variables: payload.PinAndUnpinConversationVariables{
			ConversationID: conversationID,
			LabelType:      payload.LABEL_TYPE_PINNED,
		},
		QueryID: "_TQxP2Rb0expwVP9ktGrTQ",
	}

	jsonBody, err := graphQlPayload.Encode()
	if err != nil {
		return nil, err
	}

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.UNPIN_CONVERSATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	data := response.UnpinConversationResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) React(ctx context.Context, reactionPayload *payload.ReactionActionPayload, remove bool) (*response.ReactionResponse, error) {
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

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	data := response.ReactionResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

func (c *Client) UpdateConversationAvatar(ctx context.Context, conversationID string, payload *payload.DMRequestQuery) error {
	encodedQueryBody, err := payload.Encode()
	if err != nil {
		return err
	}

	resp, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            fmt.Sprintf(endpoints.UPDATE_CONVERSATION_AVATAR_URL, conversationID),
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           encodedQueryBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeForm,
	})
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("failed to update conversation avatar id=%s (status_code=%d, response_body=%s)", conversationID, resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) UpdateConversationName(ctx context.Context, conversationID string, payload *payload.DMRequestQuery) error {
	encodedQueryBody, err := payload.Encode()
	if err != nil {
		return err
	}

	resp, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            fmt.Sprintf(endpoints.UPDATE_CONVERSATION_NAME_URL, conversationID),
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           encodedQueryBody,
		Referer:        endpoints.BASE_MESSAGES_URL,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeForm,
	})
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("failed to update conversation name id=%s (status_code=%d, response_body=%s)", conversationID, resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) AddParticipants(ctx context.Context, variables *payload.AddParticipantsPayload) (*response.AddParticipantsResponse, error) {
	graphQlPayload := payload.GraphQLPayload{
		Variables: variables,
		QueryID:   "oBwyQ0_xVbAQ8FAyG0pCRA",
	}

	url := endpoints.ADD_PARTICIPANTS_URL

	jsonBody, err := graphQlPayload.Encode()
	if err != nil {
		return nil, err
	}

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Body:           jsonBody,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
	})
	if err != nil {
		return nil, err
	}

	data := response.AddParticipantsResponse{}
	return &data, json.Unmarshal(respBody, &data)
}

// SendRawEncryptedMessage sends an already-encoded encrypted message payload.
// This is primarily for testing/debugging purposes.
func (c *Client) SendRawEncryptedMessage(ctx context.Context, conversationID, messageID, conversationToken, encodedMCE, encodedSig string) (*response.SendMessageMutationResponse, error) {
	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               conversationID,
		MessageID:                    messageID,
		ConversationToken:            conversationToken,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: encodedSig,
	},
	)

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().RawJSON("payload", jsonBody).Msg("Sending raw encrypted message")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.SEND_MESSAGE_MUTATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().RawJSON("response", respBody).Msg("Got response for encrypted message")

	var resp response.SendMessageMutationResponse
	return &resp, json.Unmarshal(respBody, &resp)
}

// SendEncryptedMessageOpts contains options for sending an encrypted message.
type SendEncryptedMessageOpts struct {
	ConversationID string
	MessageID      string // optional, generates UUID if empty
	Text           string
	Attachments    []*payload.MessageAttachment
	ReplyTo        *payload.ReplyingToPreview
	Entities       []*payload.RichTextEntity
}

// SendEncryptedMessage sends an encrypted message via the XChat protocol.
func (c *Client) SendEncryptedMessage(ctx context.Context, opts SendEncryptedMessageOpts) (*response.SendMessageMutationResponse, error) {
	// Get the server-provided conversation token for this conversation
	token, err := c.keyManager.GetConversationToken(ctx, opts.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("get conversation token: %w", err)
	}

	messageID := opts.MessageID
	if messageID == "" {
		messageID = uuid.NewString()
	}

	builder := crypto.NewMessageBuilder(c.keyManager, c.GetCurrentUserID()).
		SetMessageID(messageID).
		SetConversationID(opts.ConversationID).
		SetText(opts.Text)

	for _, att := range opts.Attachments {
		builder.AddAttachment(att)
	}
	if opts.ReplyTo != nil {
		builder.SetReplyTo(opts.ReplyTo)
	}
	if len(opts.Entities) > 0 {
		builder.SetEntities(opts.Entities)
	}

	encodedMCE, encodedSig, err := builder.BuildForSend(ctx)
	if err != nil {
		return nil, fmt.Errorf("build message: %w", err)
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               opts.ConversationID,
		MessageID:                    messageID,
		ConversationToken:            token,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: encodedSig,
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.SEND_MESSAGE_MUTATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return nil, err
	}

	var resp response.SendMessageMutationResponse
	return &resp, json.Unmarshal(respBody, &resp)
}

func (c *Client) GetInitialXChatPage(ctx context.Context, variables *payload.GetInitialXChatPageQueryVariables) (*response.GetInitialXChatPageQueryResponse, error) {
	formBody, err := variables.Encode()

	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		Str("form_body", string(formBody)).
		Str("url", endpoints.GET_INITIAL_XCHAT_PAGE_QUERY_URL).
		Msg("GetInitialXChatPage request payload")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.GET_INITIAL_XCHAT_PAGE_QUERY_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeForm,
		Body:           formBody,
	})

	if err != nil {
		c.Logger.Debug().
			Str("response_body", string(respBody)).
			Err(err).
			Msg("GetInitialXChatPage request failed")
		return nil, err
	}

	c.Logger.Trace().
		Str("response_body", string(respBody)).
		Msg("GetInitialXChatPage response")

	var resp response.GetInitialXChatPageQueryResponse
	return &resp, json.Unmarshal(respBody, &resp)

}

func (c *Client) GetInboxPageRequest(ctx context.Context, variables *payload.GetInboxPageRequestQueryVariables) (*response.GetInboxPageRequestQueryResponse, error) {
	queryString, err := variables.EncodeJSONQuery()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		Str("url", endpoints.GET_INBOX_PAGE_REQUEST_QUERY_URL).
		Str("query", queryString).
		Msg("GetInboxPageRequest payload")

	url := endpoints.GET_INBOX_PAGE_REQUEST_QUERY_URL + "?" + queryString

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
	})
	if err != nil {
		c.Logger.Debug().
			Str("response_body", string(respBody)).
			Err(err).
			Msg("GetInboxPageRequest failed")
		return nil, err
	}

	c.Logger.Trace().
		Str("response_body", string(respBody)).
		Msg("GetInboxPageRequest response")

	var resp response.GetInboxPageRequestQueryResponse
	return &resp, json.Unmarshal(respBody, &resp)
}
