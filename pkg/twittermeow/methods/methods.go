package methods

import (
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

const TwitterEpoch = 1288834974657

func ParseSnowflakeInt(msgID string) int64 {
	secs, err := strconv.ParseInt(msgID, 10, 64)
	if err != nil {
		return 0
	}
	return (secs >> 22) + TwitterEpoch
}

func ParseSnowflake(msgID string) time.Time {
	msec := ParseSnowflakeInt(msgID)
	if msec == 0 {
		return time.Time{}
	}
	return time.UnixMilli(msec)
}

func SortConversationsByTimestamp(conversations map[string]types.Conversation) []types.Conversation {
	conversationValues := maps.Values(conversations)
	slices.SortFunc(conversationValues, func(a, b types.Conversation) int {
		return strings.Compare(a.SortTimestamp, b.SortTimestamp)
	})

	return conversationValues
}

func SortMessagesByTime(messages []types.Message) {
	slices.SortFunc(messages, func(a, b types.Message) int {
		return strings.Compare(a.ID, b.ID)
	})
}

func CreateConversationID(conversationIDs []string) string {
	sort.Strings(conversationIDs)
	return strings.Join(conversationIDs, "-")
}
