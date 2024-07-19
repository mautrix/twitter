package twittermeow

import (
	"log"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func (c *Client) processEventEntries(resp *response.GetDMUserUpdatesResponse) {
	entries, err := resp.UserEvents.ToEventEntries()
	if err != nil {
		log.Fatal(err) // send event handler error
	}

	for _, entry := range entries {
		c.eventHandler(entry)
	}
}
