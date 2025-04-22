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
	response, err := tc.client.Search(ctx, payload.SearchQuery{
		Query:      identifier,
		ResultType: payload.SEARCH_RESULT_TYPE_USERS,
	})
	if err != nil {
		return nil, err
	}
	var resolvedUser types.User
	for _, user := range response.Users {
		if user.ScreenName == identifier {
			resolvedUser = user
		}
	}
	ghost, err := tc.connector.br.GetGhostByID(ctx, MakeUserID(resolvedUser.IDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to get ghost from Twitter User ID: %w", err)
	}

	var portalKey networkid.PortalKey
	if startChat {
		permissions, err := tc.client.GetDMPermissions(ctx, payload.GetDMPermissionsQuery{
			RecipientIDs: resolvedUser.IDStr,
			DMUsers:      true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get DM permissions for Twitter User: %w", err)
		}

		perms := permissions.Permissions.GetPermissionsForUser(resolvedUser.IDStr)

		if !perms.CanDM || perms.ErrorCode > 0 {
			return nil, fmt.Errorf("not allowed to DM this Twitter user: %v", resolvedUser.IDStr)
		}

		conversationID := methods.CreateConversationID([]string{resolvedUser.IDStr, ParseUserLoginID(tc.userLogin.ID)})
		portalKey = tc.MakePortalKeyFromID(conversationID)
	}

	return &bridgev2.ResolveIdentifierResponse{
		Ghost:  ghost,
		UserID: MakeUserID(resolvedUser.IDStr),
		Chat:   &bridgev2.CreateChatResponse{PortalKey: portalKey},
	}, nil
}
