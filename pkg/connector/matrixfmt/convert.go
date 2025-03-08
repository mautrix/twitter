package matrixfmt

import (
	"context"

	"maunium.net/go/mautrix/event"
)

func Parse(ctx context.Context, content *event.MessageEventContent) string {
	if content.FormattedBody == "" {
		return content.Body
	}

	parser := &HTMLParser{}
	parseCtx := NewContext(ctx)
	parseCtx.AllowedMentions = content.Mentions
	parsed := parser.Parse(content.FormattedBody, parseCtx)
	return string(parsed.String)
}
