package twitterfmt

import (
	"context"
	"strings"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func Parse(ctx context.Context, portal *bridgev2.Portal, msg *types.MessageData) *event.MessageEventContent {
	bodyHtml := strings.Builder{}

	charArr := []rune(msg.Text)
	cursor := 0

	if msg.Entities != nil {
		for _, url := range msg.Entities.URLs {
			start, end := url.Indices[0], url.Indices[1]
			if cursor < start {
				bodyHtml.WriteString(string(charArr[cursor:start]))
			}
			bodyHtml.WriteString(url.ExpandedURL)
			cursor = end
		}
		bodyHtml.WriteString(string(charArr[cursor:]))
	}

	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    msg.Text,
	}

	bodyHtmlString := bodyHtml.String()
	if bodyHtmlString != "" {
		content.Format = event.FormatHTML
		content.FormattedBody = bodyHtmlString
	}

	return content
}
