package twitterfmt

import (
	"cmp"
	"context"
	"fmt"
	"html"
	"math"
	"slices"
	"strings"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func Parse(ctx context.Context, portal *bridgev2.Portal, msg *types.MessageData) *event.MessageEventContent {
	body := strings.Builder{}
	bodyHtml := strings.Builder{}
	charArr := []rune(msg.Text)
	cursor := 0
	sortedEntites := sortEntities(msg.Entities)

	for _, union := range sortedEntites {
		switch entity := union.(type) {
		case types.URLs:
			url := entity
			start, end := url.Indices[0], url.Indices[1]
			if cursor < start {
				body.WriteString(string(charArr[cursor:start]))
				bodyHtml.WriteString(string(charArr[cursor:start]))
			}
			body.WriteString(url.ExpandedURL)
			bodyHtml.WriteString(url.ExpandedURL)
			cursor = end
		case types.UserMention:
			mention := entity
			start, end := mention.Indices[0], mention.Indices[1]
			body.WriteString(string(charArr[cursor:end]))
			if cursor < start {
				bodyHtml.WriteString(string(charArr[cursor:start]))
			}

			uid := mention.IDStr
			ghost, err := portal.Bridge.GetGhostByID(ctx, networkid.UserID(uid))
			if err != nil {
				bodyHtml.WriteString("@" + mention.ScreenName)
				continue
			}

			fmt.Fprintf(&bodyHtml,
				`<a href="%s">%s</a>`,
				ghost.Intent.GetMXID().URI().MatrixToURL(),
				ghost.Name,
			)
			cursor = end
		}
	}

	body.WriteString(string(charArr[cursor:]))
	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    html.UnescapeString(body.String()),
	}

	if msg.Entities != nil {
		bodyHtml.WriteString(string(charArr[cursor:]))
		content.Format = event.FormatHTML
		content.FormattedBody = bodyHtml.String()
	}

	return content
}

type EntityUnion = any

func getStart(union any) int {
	switch entity := union.(type) {
	case types.URLs:
		return entity.Indices[0]
	case types.UserMention:
		return entity.Indices[0]
	default:
		return math.MaxInt32
	}
}

func sortEntities(entities *types.Entities) []EntityUnion {
	if entities == nil {
		return []EntityUnion{}
	}

	merged := make([]EntityUnion, 0)

	for _, url := range entities.URLs {
		merged = append(merged, url)
	}

	for _, mention := range entities.UserMentions {
		merged = append(merged, mention)
	}

	slices.SortFunc(merged, func(a EntityUnion, b EntityUnion) int {
		return cmp.Compare(getStart(a), getStart(b))
	})

	return merged
}
