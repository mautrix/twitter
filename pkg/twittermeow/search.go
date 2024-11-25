package twittermeow

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func (c *Client) Search(params payload.SearchQuery) (*response.SearchResponse, error) {
	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s?%s", endpoints.SEARCH_TYPEAHEAD_URL, string(encodedQuery))

	apiRequestOpts := apiRequestOpts{
		URL:            url,
		Method:         http.MethodGet,
		WithClientUUID: true,
		Referer:        endpoints.BASE_MESSAGES_URL + "/compose",
	}
	_, respBody, err := c.makeAPIRequest(apiRequestOpts)
	if err != nil {
		return nil, err
	}

	data := response.SearchResponse{}
	return &data, json.Unmarshal(respBody, &data)
}
