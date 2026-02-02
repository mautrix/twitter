package twittermeow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// FetchTweet retrieves tweet data by ID using GraphQL fetchPostQuery endpoint
func (c *Client) FetchTweet(ctx context.Context, tweetID string) (*response.FetchPostQueryResponse, error) {
	// Build GraphQL variables
	variables := map[string]any{
		"postId":        tweetID,
		"withCommunity": false,
	}

	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal variables: %w", err)
	}

	// Build query URL with encoded variables
	queryURL := fmt.Sprintf("%s?variables=%s",
		endpoints.FETCH_POST_QUERY_URL,
		url.QueryEscape(string(variablesJSON)),
	)

	c.Logger.Info().
		Str("tweet_id", tweetID).
		Str("url", queryURL).
		Msg("[FetchTweet] Fetching tweet via GraphQL")

	// Make authenticated request
	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            queryURL,
		Method:         http.MethodGet,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeNone,
	})
	if err != nil {
		c.Logger.Error().Err(err).Str("tweet_id", tweetID).Msg("[FetchTweet] Request failed")
		return nil, fmt.Errorf("fetchPostQuery request failed: %w", err)
	}

	c.Logger.Info().
		Int("response_length", len(respBody)).
		Msg("[FetchTweet] Got response")

	// Parse response
	var resp response.FetchPostQueryResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		c.Logger.Error().Err(err).Str("response", string(respBody)).Msg("[FetchTweet] Failed to unmarshal")
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Validate response structure
	if resp.Data.PostResult.Result.RestID == "" {
		c.Logger.Warn().
			Str("tweet_id", tweetID).
			Str("response", string(respBody)).
			Msg("[FetchTweet] Invalid response - missing tweet data")
		return nil, fmt.Errorf("invalid response: missing tweet data for ID %s", tweetID)
	}

	c.Logger.Info().
		Str("tweet_id", tweetID).
		Str("rest_id", resp.Data.PostResult.Result.RestID).
		Str("author", resp.Data.PostResult.Result.Core.UserResults.Result.Legacy.Name).
		Msg("[FetchTweet] Successfully parsed tweet")

	return &resp, nil
}
