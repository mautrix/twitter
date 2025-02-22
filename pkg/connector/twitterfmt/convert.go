package twitterfmt

import (
	"context"
	"fmt"
	"strings"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func Parse(ctx context.Context, portal *bridgev2.Portal, msg *types.MessageData) *event.MessageEventContent {
	bodyHtml := strings.Builder{}
	charArr := []rune(msg.Text)
	cursor := 0
	sortedEntites := sortEntities(msg.Entities)

	for _, entity := range sortedEntites {
		if entity.URL != nil {
			url := entity.URL
			start, end := url.Indices[0], url.Indices[1]
			if cursor < start {
				bodyHtml.WriteString(string(charArr[cursor:start]))
			}
			bodyHtml.WriteString(url.ExpandedURL)
			cursor = end
		}
		if entity.UserMention != nil {
			mention := entity.UserMention
			start, end := mention.Indices[0], mention.Indices[1]
			if cursor < start {
				bodyHtml.WriteString(string(charArr[cursor:start]))
			}

			uid := mention.IDStr
			dmPortals, err := portal.Bridge.GetDMPortalsWith(ctx, networkid.UserID(uid))
			if err != nil {
				continue
			}

			if len(dmPortals) != 0 {
				fmt.Fprintf(&bodyHtml,
					`<a href="%s">%s</a>`,
					dmPortals[0].MXID.URI().MatrixToURL(),
					dmPortals[0].Name,
				)
			} else {
				userLogin := portal.Bridge.GetCachedUserLoginByID(networkid.UserLoginID(uid))
				text := "@" + mention.ScreenName
				if userLogin != nil {
					fmt.Fprintf(&bodyHtml,
						`<a href="%s">%s</a>`,
						userLogin.UserMXID.URI().MatrixToURL(),
						text,
					)

				} else {
					bodyHtml.WriteString(text)
				}
			}
			cursor = end
		}

	}

	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    msg.Text,
	}

	if msg.Entities != nil {
		bodyHtml.WriteString(string(charArr[cursor:]))
		content.Format = event.FormatHTML
		content.FormattedBody = bodyHtml.String()
	}

	return content
}

type Entity struct {
	UserMention *types.UserMention
	URL         *types.URLs
}

func sortEntities(entities *types.Entities) []Entity {
	if entities == nil {
		return []Entity{}
	}

	urls := entities.URLs
	mentions := entities.UserMentions
	urlIndex := 0
	mentionIndex := 0

	sorted := make([]Entity, 0)

	for urlIndex < len(urls) && mentionIndex < len(mentions) {
		urlStart := urls[urlIndex].Indices[0]
		mentionStart := mentions[mentionIndex].Indices[0]

		if urlStart < mentionStart {
			sorted = append(sorted, Entity{
				URL: &urls[urlIndex],
			})
			urlIndex++
		} else {
			sorted = append(sorted, Entity{
				UserMention: &mentions[mentionIndex],
			})
			mentionIndex++
		}
	}

	for ; urlIndex < len(urls); urlIndex++ {
		sorted = append(sorted, Entity{
			URL: &urls[urlIndex],
		})
	}

	for ; mentionIndex < len(mentions); mentionIndex++ {
		sorted = append(sorted, Entity{
			UserMention: &mentions[mentionIndex],
		})
	}

	return sorted
}
