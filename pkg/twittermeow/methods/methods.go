package methods

import (
	"sort"
	"strconv"
	"strings"
	"time"

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
	var conversationSlice []types.Conversation
	for _, conversation := range conversations {
		conversationSlice = append(conversationSlice, conversation)
	}

	sort.Slice(conversationSlice, func(j, i int) bool {
		timeJ, errJ := strconv.ParseInt(conversationSlice[j].SortTimestamp, 10, 64)
		timeI, errI := strconv.ParseInt(conversationSlice[i].SortTimestamp, 10, 64)

		if errI != nil || errJ != nil {
			return errI == nil
		}

		return timeI < timeJ
	})

	return conversationSlice
}

func SortMessagesByTime(messages []types.Message) {
	sort.Slice(messages, func(j, i int) bool {
		timeJ, errJ := strconv.ParseInt(messages[j].Time, 10, 64)
		timeI, errI := strconv.ParseInt(messages[i].Time, 10, 64)

		if errI != nil || errJ != nil {
			return errI == nil
		}

		return timeJ < timeI
	})
}

func CreateConversationID(conversationIDs []string) string {
	sort.Slice(conversationIDs, func(i, j int) bool {
		return conversationIDs[i] < conversationIDs[j]
	})

	return strings.Join(conversationIDs, "-")
}
