package twittermeow

import (
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func (c *Client) processEventEntries(resp *response.GetDMUserUpdatesResponse) error {
	entries, err := resp.UserEvents.ToEventEntries()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		c.eventHandler(entry)
	}
	return nil
}
