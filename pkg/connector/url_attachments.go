package connector

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var urlAttachmentRegex = regexp.MustCompile(`https?://[^\s<>()]+`)

const urlAttachmentTrim = ".,;:!?)]}>\"'"

func buildURLAttachments(text string) ([]*payload.MessageAttachment, []*payload.RichTextEntity) {
	if text == "" {
		return nil, nil
	}

	matches := urlAttachmentRegex.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	var attachments []*payload.MessageAttachment
	var entities []*payload.RichTextEntity
	seen := make(map[string]struct{}, len(matches))

	for _, match := range matches {
		raw := text[match[0]:match[1]]
		url := strings.TrimRight(raw, urlAttachmentTrim)
		if url == "" {
			continue
		}

		endByte := match[0] + len(url)
		if endByte <= match[0] || endByte > len(text) {
			continue
		}

		startRune := utf8.RuneCountInString(text[:match[0]])
		endRune := startRune + utf8.RuneCountInString(text[match[0]:endByte])
		if endRune <= startRune {
			continue
		}

		startIdx := int32(startRune)
		endIdx := int32(endRune)
		entities = append(entities, &payload.RichTextEntity{
			StartIndex: &startIdx,
			EndIndex:   &endIdx,
			Content: &payload.RichTextContent{
				Url: &payload.UrlRichTextContent{},
			},
		})

		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}

		urlCopy := url
		display := url
		attachments = append(attachments, &payload.MessageAttachment{
			Url: &payload.UrlAttachment{
				Url:          &urlCopy,
				DisplayTitle: &display,
			},
		})
	}

	return attachments, entities
}
