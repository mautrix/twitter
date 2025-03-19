package matrixfmt

import (
	"context"

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

func Parse(ctx context.Context, parser *format.HTMLParser, content *event.MessageEventContent) string {
	if content.FormattedBody == "" {
		return content.Body
	}

	parseCtx := format.NewContext(ctx)
	return parser.Parse(content.FormattedBody, parseCtx)
}
