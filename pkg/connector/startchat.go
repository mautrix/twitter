package connector

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

var (
	_ bridgev2.IdentifierResolvingNetworkAPI = (*TwitterClient)(nil)
)

func (tc *TwitterClient) ResolveIdentifier(ctx context.Context, identifier string, startChat bool) (*bridgev2.ResolveIdentifierResponse, error) {
	response, err := tc.client.Search(payload.SearchQuery{
		Query:      identifier,
		ResultType: payload.SEARCH_RESULT_TYPE_USERS,
	})
	var resolvedUser types.User
	for _, user := range response.Users {
		if user.ScreenName == identifier {
			resolvedUser = user
		}
	}
	if err != nil {
		return nil, err
	}
	ghost, err := tc.connector.br.GetGhostByID(ctx, networkid.UserID(resolvedUser.IDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to get ghost from Twitter User ID: %w", err)
	}

	var portalKey networkid.PortalKey
	if startChat {
		permissions, err := tc.client.GetDMPermissions(payload.GetDMPermissionsQuery{
			RecipientIds: resolvedUser.IDStr,
			DmUsers:      true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get DM permissions for Twitter User: %w", err)
		}

		perms := permissions.Permissions.GetPermissionsForUser(resolvedUser.IDStr)

		if perms.CanDm == false || perms.ErrorCode > 0 {
			return nil, fmt.Errorf("not allowed to DM this Twitter user: %v", resolvedUser.IDStr)
		}

		conversationId := methods.CreateConversationId([]string{resolvedUser.IDStr, tc.client.GetCurrentUserID()})
		portalKey = tc.MakePortalKeyFromID(conversationId)
	}

	return &bridgev2.ResolveIdentifierResponse{
		Ghost:  ghost,
		UserID: networkid.UserID(resolvedUser.ID),
		Chat:   &bridgev2.CreateChatResponse{PortalKey: portalKey},
	}, nil
}
