package methods

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

var Charset = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

// retrieved from main page resp, its a 2 year old timestamp I don't think this is changing
const fetchedTime = 1661971138705

func GenerateEventValue() int64 {
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	return ts - fetchedTime
}

func GetTimestampMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

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

func CreateConversationId(conversationIds []string) string {
	sort.Slice(conversationIds, func(i, j int) bool {
		return conversationIds[i] < conversationIds[j]
	})

	return strings.Join(conversationIds, "-")
}
