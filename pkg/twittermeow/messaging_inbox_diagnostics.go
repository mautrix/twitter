package twittermeow

import (
	"strconv"

	"github.com/tidwall/gjson"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func (c *Client) logXChatInboxItemShape(rawPage gjson.Result, page *response.XChatInboxPage, decoderPath string) {
	for i, item := range page.Items {
		if item.ConversationUnavailable || item.ConversationDetail.ConversationID != "" {
			continue
		}
		// Typed decoding collapses absent/null objects and empty IDs. Inspect only
		// the first item that the connector will reject, without retaining values.
		rawItem := rawPage.Get("items." + strconv.Itoa(i))
		detail := rawItem.Get("conversation_detail")
		c.Logger.Warn().
			Str("decoder_path", decoderPath).
			Str("page_kind", inboxJSONKind(rawPage)).
			Int("item_index", i).Int("item_count", len(page.Items)).
			Str("item_kind", inboxJSONKind(rawItem)).
			Int("item_field_count", len(rawItem.Map())).
			Bool("item_typename_present", rawItem.Get("__typename").Exists()).
			Str("detail_kind", inboxJSONKind(detail)).
			Int("detail_field_count", len(detail.Map())).
			Bool("detail_typename_present", detail.Get("__typename").Exists()).
			Str("conversation_id_kind", inboxJSONKind(detail.Get("conversation_id"))).
			Int("page_error_count", len(page.Errors)).
			Int("latest_message_count", len(item.LatestMessageEvents)).
			Int("encoded_message_count", len(item.EncodedMessageEvents)).
			Bool("notifiable_message_present", item.LatestNotifiableMessageCreateEvent != "").
			Int("key_change_count", len(item.LatestConversationKeyChangeEvents)).
			Int("read_event_count", len(item.LatestReadEventsPerParticipant)).
			Bool("has_more", item.HasMore).
			Msg("XChat inbox item has no conversation ID")
		return
	}
}

func inboxJSONKind(value gjson.Result) string {
	switch {
	case !value.Exists():
		return "missing"
	case value.Type == gjson.Null:
		return "null"
	case value.IsObject():
		return "object"
	case value.IsArray():
		return "array"
	case value.Type == gjson.String && value.Str == "":
		return "empty_string"
	default:
		return value.Type.String()
	}
}
