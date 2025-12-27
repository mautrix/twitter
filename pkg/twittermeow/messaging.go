package twittermeow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// FormEncoder is implemented by types that can encode themselves to form data.
type FormEncoder interface {
	Encode() ([]byte, error)
}

// makeXChatFormRequest is a generic helper for XChat API form-encoded POST requests.
func makeXChatFormRequest[T any](c *Client, ctx context.Context, url string, variables FormEncoder, opName string) (*T, error) {
	formBody, err := variables.Encode()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		Str("url", url).
		Str("form_body", string(formBody)).
		Msg(opName + " payload")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            url,
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
			Msg(opName + " failed")
		return nil, err
	}

	c.Logger.Trace().
		Str("response_body", string(respBody)).
		Msg(opName + " response")

	var resp T
	return &resp, json.Unmarshal(respBody, &resp)
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
	var sigPtr *string
	if encodedSig != "" {
		sigPtr = &encodedSig
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               conversationID,
		MessageID:                    messageID,
		ConversationToken:            conversationToken,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: sigPtr,
	})

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

// SendEncryptedEditOpts contains options for sending an encrypted message edit.
type SendEncryptedEditOpts struct {
	ConversationID          string
	MessageID               string // optional, generates UUID if empty
	TargetMessageSequenceID string
	UpdatedText             string
	Entities                []*payload.RichTextEntity
}

// SendEncryptedReaction sends a reaction add/remove via the XChat protocol.
// targetMessageSequenceID must be the XChat message sequence ID of the message being reacted to.
func (c *Client) SendEncryptedReaction(ctx context.Context, conversationID, targetMessageSequenceID, emoji string, remove bool) (*response.SendMessageMutationResponse, error) {
	token, err := c.keyManager.GetConversationToken(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("get conversation token: %w", err)
	}

	messageID := uuid.NewString()

	builder := crypto.NewMessageBuilder(c.keyManager, c.GetCurrentUserID()).
		SetMessageID(messageID).
		SetConversationID(conversationID)

	if remove {
		builder.SetReactionRemove(targetMessageSequenceID, emoji)
	} else {
		builder.SetReactionAdd(targetMessageSequenceID, emoji)
	}

	encodedMCE, encodedSig, err := builder.BuildForSend(ctx)
	if err != nil {
		return nil, fmt.Errorf("build reaction: %w", err)
	}

	var sigPtr *string
	if encodedSig != "" {
		sigPtr = &encodedSig
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               conversationID,
		MessageID:                    messageID,
		ConversationToken:            token,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: sigPtr,
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("SendMessageMutation reaction payload")

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

// SendEncryptedEdit sends a message edit via the XChat protocol.
// targetMessageSequenceID must be the XChat message sequence ID of the message being edited.
func (c *Client) SendEncryptedEdit(ctx context.Context, opts SendEncryptedEditOpts) (*response.SendMessageMutationResponse, error) {
	if opts.ConversationID == "" {
		return nil, fmt.Errorf("conversation ID is required")
	}
	if opts.TargetMessageSequenceID == "" {
		return nil, fmt.Errorf("target message sequence ID is required")
	}

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
		SetMessageEdit(opts.TargetMessageSequenceID, opts.UpdatedText, opts.Entities)

	encodedMCE, encodedSig, err := builder.BuildForSend(ctx)
	if err != nil {
		return nil, fmt.Errorf("build edit: %w", err)
	}

	var sigPtr *string
	if encodedSig != "" {
		sigPtr = &encodedSig
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               opts.ConversationID,
		MessageID:                    messageID,
		ConversationToken:            token,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: sigPtr,
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("SendMessageMutation edit payload")

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
		if replyJSON, err := json.Marshal(opts.ReplyTo); err != nil {
			c.Logger.Warn().
				Err(err).
				Interface("reply_preview", opts.ReplyTo).
				Str("conversation_id", opts.ConversationID).
				Str("message_id", messageID).
				Msg("Failed to marshal ReplyingToPreview for logging")
		} else {
			c.Logger.Debug().
				Str("conversation_id", opts.ConversationID).
				Str("message_id", messageID).
				RawJSON("reply_preview", replyJSON).
				Msg("Sending reply with ReplyingToPreview")
		}
		builder.SetReplyTo(opts.ReplyTo)
	}
	if len(opts.Entities) > 0 {
		builder.SetEntities(opts.Entities)
	}

	encodedMCE, encodedSig, err := builder.BuildForSend(ctx)
	if err != nil {
		return nil, fmt.Errorf("build message: %w", err)
	}

	var sigPtr *string
	if encodedSig != "" {
		sigPtr = &encodedSig
	}

	pl := payload.NewSendMessageMutationPayload(payload.SendMessageMutationVariables{
		ConversationID:               opts.ConversationID,
		MessageID:                    messageID,
		ConversationToken:            token,
		EncodedMessageCreateEvent:    encodedMCE,
		EncodedMessageEventSignature: sigPtr,
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("SendMessageMutation payload")

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
	return makeXChatFormRequest[response.GetInitialXChatPageQueryResponse](c, ctx, endpoints.GET_INITIAL_XCHAT_PAGE_QUERY_URL, variables, "GetInitialXChatPage")
}

func (c *Client) GetInboxPageRequest(ctx context.Context, variables *payload.GetInboxPageRequestQueryVariables) (*response.GetInboxPageRequestQueryResponse, error) {
	return makeXChatFormRequest[response.GetInboxPageRequestQueryResponse](c, ctx, endpoints.GET_INBOX_PAGE_REQUEST_QUERY_URL, variables, "GetInboxPageRequest")
}

func (c *Client) GetConversationData(ctx context.Context, variables *payload.GetInboxPageConversationDataQueryVariables) (*response.GetInboxPageConversationDataResponse, error) {
	return makeXChatFormRequest[response.GetInboxPageConversationDataResponse](c, ctx, endpoints.GET_INBOX_PAGE_CONV_DATA_QUERY_URL, variables, "GetConversationData")
}

func (c *Client) GetUsersByIdsForXChat(ctx context.Context, variables *payload.GetUsersByIdsForXChatVariables) (*response.GetUsersByIdsForXChatResponse, error) {
	return makeXChatFormRequest[response.GetUsersByIdsForXChatResponse](c, ctx, endpoints.GET_USERS_BY_IDS_FOR_XCHAT_URL, variables, "GetUsersByIdsForXChat")
}

func (c *Client) GetConversationPage(ctx context.Context, variables *payload.GetConversationPageQueryVariables) (*response.GetConversationPageQueryResponse, error) {
	return makeXChatFormRequest[response.GetConversationPageQueryResponse](c, ctx, endpoints.GET_CONVERSATION_PAGE_QUERY_URL, variables, "GetConversationPage")
}

// DeleteXChatMessageOpts contains options for deleting messages via XChat.
type DeleteXChatMessageOpts struct {
	ConversationID string
	SequenceIDs    []string // The sequence IDs of the messages to delete
	DeleteForAll   bool     // If true, delete for everyone; if false, delete only for self
}

// DeleteXChatMessage deletes messages via the XChat DeleteMessageMutation GraphQL endpoint.
func (c *Client) DeleteXChatMessage(ctx context.Context, opts DeleteXChatMessageOpts) error {
	if opts.ConversationID == "" {
		return fmt.Errorf("conversation ID is required")
	}
	if len(opts.SequenceIDs) == 0 {
		return fmt.Errorf("at least one sequence ID is required")
	}

	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return fmt.Errorf("sender ID is required")
	}

	token, err := c.keyManager.GetConversationToken(ctx, opts.ConversationID)
	if err != nil {
		return fmt.Errorf("get conversation token: %w", err)
	}

	keyPair, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	}

	deleteAction := payload.DeleteMessageActionTypeForSelf
	thriftAction := int32(payload.DeleteMessageActionDeleteForSelf)
	if opts.DeleteForAll {
		deleteAction = payload.DeleteMessageActionTypeForAll
		thriftAction = int32(payload.DeleteMessageActionDeleteForAll)
	}

	createdAtMsec := fmt.Sprintf("%d", time.Now().UnixMilli())

	actionSignatures := make([]payload.DeleteMessageActionSignature, 0, len(opts.SequenceIDs))
	for _, seqID := range opts.SequenceIDs {
		messageID := uuid.NewString()

		detail := &payload.MessageEventDetail{
			MessageDeleteEvent: &payload.MessageDeleteEvent{
				SequenceIds:         []string{seqID},
				DeleteMessageAction: &thriftAction,
			},
		}

		encodedDetail, err := crypto.EncodeMessageEventDetail(detail)
		if err != nil {
			return fmt.Errorf("encode message event detail for %s: %w", seqID, err)
		}

		actionSig := payload.DeleteMessageActionSignature{
			MessageID:                 messageID,
			EncodedMessageEventDetail: encodedDetail,
		}

		if keyPair != nil && keyPair.SigningKey != nil && keyPair.KeyVersion != "" {
			signature, err := crypto.SignMessageDeleteEvent(
				keyPair.SigningKey,
				messageID,
				senderID,
				opts.ConversationID,
				token,
				createdAtMsec,
				encodedDetail,
			)
			if err != nil {
				return fmt.Errorf("sign delete event for %s: %w", seqID, err)
			}

			sigVersion := crypto.SignatureVersion4
			actionSig.MessageEventSignature = &payload.DeleteMessageEventSignatureJSON{
				Signature:        signature,
				PublicKeyVersion: keyPair.KeyVersion,
				SignatureVersion: sigVersion,
			}
		}

		actionSignatures = append(actionSignatures, actionSig)
	}

	pl := payload.NewDeleteMessageMutationPayload(payload.DeleteMessageMutationVariables{
		SequenceIDs:         opts.SequenceIDs,
		ConversationID:      opts.ConversationID,
		DeleteMessageAction: deleteAction,
		ActionSignatures:    actionSignatures,
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("DeleteMessageMutation payload")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.DELETE_MESSAGE_MUTATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("response", respBody).
		Msg("DeleteMessageMutation response")

	return nil
}

// MuteConversation mutes a conversation via the XChat MuteConversation GraphQL endpoint.
func (c *Client) MuteConversation(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return fmt.Errorf("conversation ID is required")
	}

	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return fmt.Errorf("sender ID is required")
	}

	token, err := c.keyManager.GetConversationToken(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation token: %w", err)
	}

	keyPair, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	}

	createdAtMsec := fmt.Sprintf("%d", time.Now().UnixMilli())
	messageID := uuid.NewString()

	detail := &payload.MuteConversation{
		MutedConversationIds: []string{conversationID},
	}

	encodedDetail, err := crypto.EncodeMuteConversation(detail)
	if err != nil {
		return fmt.Errorf("encode mute conversation: %w", err)
	}

	actionSig := payload.DeleteMessageActionSignature{
		MessageID:                 messageID,
		EncodedMessageEventDetail: encodedDetail,
	}

	if keyPair != nil && keyPair.SigningKey != nil && keyPair.KeyVersion != "" {
		signature, err := crypto.SignMuteConversation(
			keyPair.SigningKey,
			messageID,
			senderID,
			conversationID,
			token,
			createdAtMsec,
			encodedDetail,
		)
		if err != nil {
			return fmt.Errorf("sign mute event: %w", err)
		}

		sigVersion := crypto.SignatureVersion4
		actionSig.MessageEventSignature = &payload.DeleteMessageEventSignatureJSON{
			Signature:        signature,
			PublicKeyVersion: keyPair.KeyVersion,
			SignatureVersion: sigVersion,
		}
	}

	pl := payload.NewMuteConversationMutationPayload(payload.MuteConversationMutationVariables{
		ConversationIDs:  []string{conversationID},
		ActionSignatures: []payload.DeleteMessageActionSignature{actionSig},
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("MuteConversation payload")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.MUTE_CONVERSATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("response", respBody).
		Msg("MuteConversation response")

	return nil
}

// UnmuteConversation unmutes a conversation via the XChat UnmuteConversation GraphQL endpoint.
func (c *Client) UnmuteConversation(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return fmt.Errorf("conversation ID is required")
	}

	senderID := c.GetCurrentUserID()
	if senderID == "" {
		return fmt.Errorf("sender ID is required")
	}

	token, err := c.keyManager.GetConversationToken(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("get conversation token: %w", err)
	}

	keyPair, err := c.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		return fmt.Errorf("get signing key: %w", err)
	}

	createdAtMsec := fmt.Sprintf("%d", time.Now().UnixMilli())
	messageID := uuid.NewString()

	detail := &payload.UnmuteConversation{
		UnmutedConversationIds: []string{conversationID},
	}

	encodedDetail, err := crypto.EncodeUnmuteConversation(detail)
	if err != nil {
		return fmt.Errorf("encode unmute conversation: %w", err)
	}

	actionSig := payload.DeleteMessageActionSignature{
		MessageID:                 messageID,
		EncodedMessageEventDetail: encodedDetail,
	}

	if keyPair != nil && keyPair.SigningKey != nil && keyPair.KeyVersion != "" {
		signature, err := crypto.SignUnmuteConversation(
			keyPair.SigningKey,
			messageID,
			senderID,
			conversationID,
			token,
			createdAtMsec,
			encodedDetail,
		)
		if err != nil {
			return fmt.Errorf("sign unmute event: %w", err)
		}

		sigVersion := crypto.SignatureVersion4
		actionSig.MessageEventSignature = &payload.DeleteMessageEventSignatureJSON{
			Signature:        signature,
			PublicKeyVersion: keyPair.KeyVersion,
			SignatureVersion: sigVersion,
		}
	}

	pl := payload.NewUnmuteConversationMutationPayload(payload.UnmuteConversationMutationVariables{
		ConversationIDs:  []string{conversationID},
		ActionSignatures: []payload.DeleteMessageActionSignature{actionSig},
	})

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("UnmuteConversation payload")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.UNMUTE_CONVERSATION_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return err
	}

	c.Logger.Debug().
		RawJSON("response", respBody).
		Msg("UnmuteConversation response")

	return nil
}

func (c *Client) GetPublicKeys(ctx context.Context, userIDs []string) (*response.GetPublicKeysResponse, error) {
	variables := payload.NewGetPublicKeysQueryVariables(userIDs)

	formBody, err := variables.Encode()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		Str("url", endpoints.GET_PUBLIC_KEYS_QUERY_URL).
		Str("form_body", string(formBody)).
		Msg("GetPublicKeys payload")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.GET_PUBLIC_KEYS_QUERY_URL,
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
			Msg("GetPublicKeys failed")
		return nil, err
	}

	c.Logger.Trace().
		Str("response_body", string(respBody)).
		Msg("GetPublicKeys response")

	var resp response.GetPublicKeysResponse
	return &resp, json.Unmarshal(respBody, &resp)
}
