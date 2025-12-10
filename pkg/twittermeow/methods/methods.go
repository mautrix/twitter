package methods

import (
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

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

// ParseMsecTimestamp parses a milliseconds-since-epoch string to time.Time.
func ParseMsecTimestamp(msec string) time.Time {
	ms, err := strconv.ParseInt(msec, 10, 64)
	if err != nil || ms == 0 {
		return time.Time{}
	}
	return time.UnixMilli(ms)
}

// ParseInt64 parses a string to int64 for sequence IDs.
func ParseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func CompareSnowflake(a, b string) int {
	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return strings.Compare(a, b)
}

func SortConversationsByTimestamp(conversations map[string]*types.Conversation) []*types.Conversation {
	return slices.SortedFunc(maps.Values(conversations), func(a, b *types.Conversation) int {
		return CompareSnowflake(a.SortTimestamp, b.SortTimestamp)
	})
}

func SortMessagesByTime(messages []*types.Message) {
	slices.SortFunc(messages, func(a, b *types.Message) int {
		return CompareSnowflake(a.ID, b.ID)
	})
}

func CreateConversationID(conversationIDs []string) string {
	sort.Strings(conversationIDs)
	return strings.Join(conversationIDs, ":")
}
