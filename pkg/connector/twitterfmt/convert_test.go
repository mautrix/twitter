package twitterfmt_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mau.fi/mautrix-twitter/pkg/connector/twitterfmt"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		ins  string
		ine  *types.Entities
		body string
		html string
	}{
		{
			name: "plain",
			ins:  "Hello world!",
			body: "Hello world!",
		},
		{
			name: "emoji before url",
			ins:  "🚀 https://t.co/WCPQgzfcO4 abc",
			ine: &types.Entities{
				URLs: []types.URLs{
					{ExpandedURL: "https://x.com",
						Indices: []int{2, 25},
					},
				},
			},
			body: "🚀 https://x.com abc",
			html: "🚀 <a href=\"https://x.com\">https://x.com</a> abc",
		},
		{
			name: "xchat url without expanded",
			ins:  "Check this https://example.com/path",
			ine: &types.Entities{
				URLs: []types.URLs{
					{Indices: []int{11, 35}},
				},
			},
			body: "Check this https://example.com/path",
			html: "Check this <a href=\"https://example.com/path\">https://example.com/path</a>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := &types.MessageData{
				Text:     test.ins,
				Entities: test.ine,
			}
			parsed := twitterfmt.Parse(context.TODO(), nil, msg)
			assert.Equal(t, test.body, parsed.Body)
			assert.Equal(t, test.html, parsed.FormattedBody)
		})
	}

}
