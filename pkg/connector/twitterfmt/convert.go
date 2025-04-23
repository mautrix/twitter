package twitterfmt

import (
	"cmp"
	"context"
	"fmt"
	"html"
	"math"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func Parse(ctx context.Context, portal *bridgev2.Portal, msg *types.MessageData) *event.MessageEventContent {
	body := strings.Builder{}
	bodyHTML := strings.Builder{}
	charArr := []rune(msg.Text)
	cursor := 0
	sortedEntites := sortEntities(msg.Entities)
	var mentions event.Mentions

	for _, union := range sortedEntites {
		switch entity := union.(type) {
		case types.URLs:
			url := entity
			start, end := url.Indices[0], url.Indices[1]
			if cursor < start {
				body.WriteString(string(charArr[cursor:start]))
				bodyHTML.WriteString(string(charArr[cursor:start]))
			}
			body.WriteString(url.ExpandedURL)
			bodyHTML.WriteString(url.ExpandedURL)
			cursor = end
		case types.UserMention:
			mention := entity
			start, end := mention.Indices[0], mention.Indices[1]
			body.WriteString(string(charArr[cursor:end]))
			if cursor < start {
				bodyHTML.WriteString(string(charArr[cursor:start]))
			}

			uid := mention.IDStr
			ghost, err := portal.Bridge.GetGhostByID(ctx, networkid.UserID(uid)) // TODO use MakeUserID
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("Failed to get ghost")
				bodyHTML.WriteString(string(charArr[start:end]))
				continue
			}
			targetMXID := ghost.Intent.GetMXID()
			login := portal.Bridge.GetCachedUserLoginByID(networkid.UserLoginID(uid)) // TODO use MakeUserLoginID
			if login != nil {
				targetMXID = login.UserMXID
			}
			_, _ = fmt.Fprintf(&bodyHTML,
				`<a href="%s">%s</a>`,
				targetMXID.URI().MatrixToURL(),
				string(charArr[start:end]),
			)
			mentions.Add(targetMXID)
			cursor = end
		}
	}

	body.WriteString(string(charArr[cursor:]))
	content := &event.MessageEventContent{
		MsgType:  event.MsgText,
		Body:     html.UnescapeString(body.String()),
		Mentions: &mentions,
	}

	if msg.Entities != nil {
		bodyHTML.WriteString(string(charArr[cursor:]))
		content.Format = event.FormatHTML
		content.FormattedBody = strings.ReplaceAll(bodyHTML.String(), "\n", "<br>")
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
