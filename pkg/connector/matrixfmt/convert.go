package matrixfmt

import (
	"context"

	"maunium.net/go/mautrix/event"
)

func Parse(ctx context.Context, parser *HTMLParser, content *event.MessageEventContent) string {
	if content.FormattedBody == "" {
		return content.Body
	}

	parseCtx := NewContext(ctx)
	parseCtx.AllowedMentions = content.Mentions
	parsed := parser.Parse(content.FormattedBody, parseCtx)
	return string(parsed.String)
}
