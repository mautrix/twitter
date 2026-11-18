package connector

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

var (
	_ bridgev2.IdentifierResolvingNetworkAPI = (*TwitterClient)(nil)
	_ bridgev2.GhostDMCreatingNetworkAPI     = (*TwitterClient)(nil)
)

func (tc *TwitterClient) ResolveIdentifier(ctx context.Context, identifier string, startChat bool) (*bridgev2.ResolveIdentifierResponse, error) {
	identifier = strings.TrimPrefix(identifier, "@")
	response, err := tc.client.Search(ctx, payload.SearchQuery{
		Query:      identifier,
		ResultType: payload.SEARCH_RESULT_TYPE_USERS,
	})
	if errors.Is(err, twittermeow.ErrRatelimitExceeded) {
		return nil, bridgev2.RespError(mautrix.MLimitExceeded.WithMessage("X is rate limiting requests, try again later").WithStatus(http.StatusBadRequest))
	} else if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to search for X users")
		return nil, bridgev2.RespError(mautrix.MUnknown.WithMessage("Failed to search for X users")).WithInternalError(err)
	}
	var resolvedUser types.User
	for _, user := range response.Users {
		if strings.EqualFold(user.ScreenName, identifier) {
			resolvedUser = user
		}
	}
	if resolvedUser.IDStr == "" {
		return nil, bridgev2.RespError(mautrix.MNotFound.WithMessage("X user not found").WithStatus(http.StatusNotFound))
	}
	ghost, err := tc.connector.br.GetGhostByID(ctx, MakeUserID(resolvedUser.IDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to get ghost from Twitter User ID: %w", err)
	}

	resp := &bridgev2.ResolveIdentifierResponse{
		Ghost:  ghost,
		UserID: MakeUserID(resolvedUser.IDStr),
	}
	if startChat {
		resp.Chat, err = tc.createDMWithUser(ctx, resolvedUser.IDStr)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (tc *TwitterClient) CreateChatWithGhost(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.CreateChatResponse, error) {
	return tc.createDMWithUser(ctx, ParseUserID(ghost.ID))
}

func (tc *TwitterClient) createDMWithUser(ctx context.Context, userID string) (*bridgev2.CreateChatResponse, error) {
	permissions, err := tc.client.GetDMPermissions(ctx, payload.GetDMPermissionsQuery{
		RecipientIDs: userID,
		DMUsers:      true,
	})
	if errors.Is(err, twittermeow.ErrRatelimitExceeded) {
		return nil, bridgev2.RespError(mautrix.MLimitExceeded.WithMessage("X is rate limiting requests, try again later").WithStatus(http.StatusBadRequest))
	} else if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get DM permissions for X user")
		return nil, bridgev2.RespError(mautrix.MUnknown.WithMessage("Failed to check if you can message this X user")).WithInternalError(err)
	}

	perms := permissions.Permissions.GetPermissionsForUser(userID)
	if perms == nil || !perms.CanDM || perms.ErrorCode > 0 {
		return nil, bridgev2.RespError(mautrix.MForbidden.WithMessage("You can't send direct messages to this X user"))
	}

	conversationID := methods.CreateConversationID([]string{userID, ParseUserLoginID(tc.userLogin.ID)})
	return &bridgev2.CreateChatResponse{PortalKey: tc.MakePortalKeyFromID(conversationID)}, nil
}
