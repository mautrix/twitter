package methods

import (
	"cmp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func UnixStringMilliToTime(input string) (time.Time, error) {
	secs, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(secs), nil
}

func SortConversationsByTimestamp(conversations map[string]types.Conversation) []types.Conversation {
	conversationValues := maps.Values(conversations)
	slices.SortFunc(conversationValues, func(a, b types.Conversation) int {
		timeA, errA := strconv.ParseInt(a.SortTimestamp, 10, 64)
		timeB, errB := strconv.ParseInt(b.SortTimestamp, 10, 64)
		if errB != nil || errA != nil {
			return 0
		}

		return cmp.Compare(timeA, timeB)
	})

	return conversationValues
}

func SortMessagesByTime(messages []types.Message) {
	slices.SortFunc(messages, func(a, b types.Message) int {
		timeA, errA := strconv.ParseInt(a.Time, 10, 64)
		timeB, errB := strconv.ParseInt(b.Time, 10, 64)
		if errB != nil || errA != nil {
			return 0
		}

		return cmp.Compare(timeA, timeB)
	})
}

func CreateConversationID(conversationIDs []string) string {
	sort.Strings(conversationIDs)
	return strings.Join(conversationIDs, "-")
}
