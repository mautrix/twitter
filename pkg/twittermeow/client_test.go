//lint:file-ignore U1000 -

package twittermeow_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/debug"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/event"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"

	"github.com/google/uuid"
)

var cli *twittermeow.Client

func TestXClientLogin(t *testing.T) {
	cookieStr, err := os.ReadFile("cookies.txt")
	if err != nil {
		log.Fatal(err)
	}
	cookieStruct := cookies.NewCookiesFromString(string(cookieStr))

	clientOptions := twittermeow.ClientOpts{
		Cookies:       cookieStruct,
		EventHandler:  eventHandler,
		WithJOTClient: true,
	}
	cli = twittermeow.NewClient(&clientOptions, debug.NewLogger())
	cli.SetEventHandler(eventHandler)

	_, _, err = cli.LoadMessagesPage()
	if err != nil {
		log.Fatal(err)
	}

	err = cli.Connect()
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan struct{})
	<-wait
}

func deleteConversationTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation

	payload := payload.DmRequestQuery{}.Default()
	err = cli.DeleteConversation(firstConversation.ConversationID, payload)
	if err != nil {
		log.Fatal(err)
	}

	cli.Logger.Info().Str("conversation_id", firstConversation.ConversationID).Msg("Successfully deleted the top conversation in my inbox")
}

func pinTopConversationTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation

	pinnedResponse, err := cli.PinConversation(firstConversation.ConversationID)
	if err != nil {
		log.Fatalf("failed to pin conversation with id: %s (%s)", firstConversation.ConversationID, err.Error())
	}

	cli.Logger.Info().Any("pinnedResponse", pinnedResponse).Str("conversation_id", firstConversation.ConversationID).Msg("Successfully pinned the conversation at the very top of my inbox")
}

func createConversationAndSendMessageTest() {
	searchQuery := payload.SearchQuery{
		IncludeExtIsBlueVerified:    "1",
		IncludeExtVerifiedType:      "1",
		IncludeExtProfileImageShape: "1",
		Query:                       "dest",
		Src:                         "compose_message",
		ResultType:                  payload.SEARCH_RESULT_TYPE_USERS,
	}
	searchResponse, err := cli.Search(searchQuery)
	if err != nil {
		log.Fatal(err)
	}

	var pickedUser *types.User
	for _, user := range searchResponse.Users {
		if user.IsDmAble {
			pickedUser = &user
			break
		}
	}

	if pickedUser == nil {
		log.Fatalf("failed to find a user that I can dm while searching for query %s", searchQuery.Query)
	}

	dmPermissionsQuery := payload.GetDMPermissionsQuery{
		RecipientIds: pickedUser.IDStr,
		DmUsers:      true,
	}
	dmPermissionsResponse, err := cli.GetDMPermissions(dmPermissionsQuery)
	if err != nil {
		log.Fatalf("failed to fetch dm permissions for recipients with ids %s", dmPermissionsQuery.RecipientIds)
	}

	pickedUserDMPermissions := dmPermissionsResponse.Permissions.GetPermissionsForUser(pickedUser.IDStr)
	if pickedUserDMPermissions == nil {
		log.Fatalf("failed to find permissions for user with id %s", pickedUser.IDStr)
	}

	if !pickedUserDMPermissions.CanDm {
		log.Fatalf("exiting because I do not have the correct permissions to dm user with id: %s (canDm=%v, errorCode=%d)", pickedUser.IDStr, pickedUserDMPermissions.CanDm, pickedUserDMPermissions.ErrorCode)
	}

	myUserID := cli.GetCurrentUserID()
	conversationId := fmt.Sprintf("%s-%s", pickedUser.IDStr, myUserID)

	contextQuery := (&payload.DmRequestQuery{}).Default()
	contextQuery.IncludeConversationInfo = true
	_, err = cli.FetchConversationContext(conversationId, contextQuery, payload.CONTEXT_FETCH_DM_CONVERSATION)
	if err != nil {
		log.Fatal(err)
	}

	sendDirectMessagePayload := &payload.SendDirectMessagePayload{
		ConversationID:    conversationId,
		CardsPlatform:     "Web-12",
		IncludeCards:      1,
		IncludeQuoteCount: true,
		DmUsers:           false,
		RecipientIds:      false,
		Text:              "testing creating a conversation by sending a message",
	}
	sentMessageResp, err := cli.SendDirectMessage(sendDirectMessagePayload)
	if err != nil {
		log.Fatalf("failed to initialize and send message to conversation by id %s", conversationId)
	}

	cli.Logger.Info().Any("response", sentMessageResp).Str("conversation_id", conversationId).Str("other_user_id", pickedUser.IDStr).Msg("Successfully initialized new conversation by sending a test message")
}

func deleteMessageForMeTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation
	mostRecentMessage := conversations[0].Messages[0]

	payload := &payload.DMMessageDeleteMutationVariables{
		MessageID: mostRecentMessage.MessageData.ID,
	}
	deleteMessageResp, err := cli.DeleteMessageForMe(payload)
	if err != nil {
		log.Fatalf("failed to delete message with id %s in conversation with id %s", payload.MessageID, firstConversation.ConversationID)
	}

	cli.Logger.Info().Any("deleteMessageResp", deleteMessageResp).Str("conversation_id", firstConversation.ConversationID).Str("message_id", payload.MessageID).Msg("Deleted most recent message in conversation for me")
}

func uploadAndSendImageTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation

	// Note: this file doesn't exist
	imgBytes, err := os.ReadFile("test_data/testimage1.jpg")
	if err != nil {
		log.Fatal(err)
	}

	uploadQuery := &payload.UploadMediaQuery{
		MediaType:     payload.MEDIA_TYPE_IMAGE_JPEG,
		MediaCategory: payload.MEDIA_CATEGORY_DM_IMAGE,
	}

	mediaResult, err := cli.UploadMedia(uploadQuery, imgBytes)
	if err != nil {
		log.Fatal(err)
	}

	payload := &payload.SendDirectMessagePayload{
		ConversationID:    firstConversation.ConversationID,
		RequestID:         uuid.NewString(),
		CardsPlatform:     "Web-12",
		IncludeCards:      1,
		Text:              "",
		MediaID:           mediaResult.MediaIDString,
		IncludeQuoteCount: true,
		RecipientIds:      false,
		DmUsers:           false,
	}

	sentMessageResponse, err := cli.SendDirectMessage(payload)
	if err != nil {
		log.Fatalf("failed to send image to conversation with id %s (%s)", firstConversation.ConversationID, err.Error())
	}

	cli.Logger.Info().Any("response", sentMessageResponse).Str("conversation_id", firstConversation.ConversationID).Msg("Sent test image to first conversation")
}

func testReplyToMessage(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation
	mostRecentMessage := conversations[0].Messages[0]

	payload := &payload.SendDirectMessagePayload{
		ConversationID:    firstConversation.ConversationID,
		RequestID:         uuid.NewString(),
		ReplyToDmID:       mostRecentMessage.MessageData.ID,
		CardsPlatform:     "Web-12",
		IncludeCards:      1,
		Text:              "this is a test reply",
		IncludeQuoteCount: true,
		RecipientIds:      false,
		DmUsers:           false,
	}

	sentReplyResponse, err := cli.SendDirectMessage(payload)
	if err != nil {
		log.Fatalf("failed to reply to message with id %s in conversation with id %s (%s)", mostRecentMessage.MessageData.ID, mostRecentMessage.ConversationID, err.Error())
	}

	cli.Logger.Info().Any("response", sentReplyResponse).Str("conversation_id", firstConversation.ConversationID).Msg("Sent test reply to most recent message in conversation")
}

func uploadAndSendGifTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation

	uploadQuery := &payload.UploadMediaQuery{
		SourceURL:     "https://media1.giphy.com/media/v1.Y2lkPWU4MjZjOWZjYzdkYTk5YWM3ODE2MjczYTlkYWFiZjY2MDkxYTIyZDJmMjVlMDAwYiZlcD12MV9naWZzX2NhdGVnb3JpZXNfY2F0ZWdvcnlfdGFnJmN0PWc/z3HFoEzXCMykr4L0TB/giphy.gif",
		MediaType:     payload.MEDIA_TYPE_IMAGE_GIF,
		MediaCategory: payload.MEDIA_CATEGORY_DM_GIF,
	}
	mediaResult, err := cli.UploadMedia(uploadQuery, nil)
	if err != nil {
		log.Fatal(err)
	}

	payload := &payload.SendDirectMessagePayload{
		ConversationID:    firstConversation.ConversationID,
		RequestID:         uuid.NewString(),
		CardsPlatform:     "Web-12",
		IncludeCards:      1,
		Text:              "",
		MediaID:           mediaResult.MediaIDString,
		IncludeQuoteCount: true,
		RecipientIds:      false,
		DmUsers:           false,
	}

	sentMessageResponse, err := cli.SendDirectMessage(payload)
	if err != nil {
		log.Fatalf("failed to send gif to conversation with id %s (%s)", firstConversation.ConversationID, err.Error())
	}

	cli.Logger.Info().Any("response", sentMessageResponse).Str("conversation_id", firstConversation.ConversationID).Msg("Sent test gif to first conversation")
}

func uploadAndSendVideoTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation

	// Note: this file doesn't exist
	videoBytes, err := os.ReadFile("test_data/testvideo1.mp4")
	if err != nil {
		log.Fatal(err)
	}

	uploadQuery := &payload.UploadMediaQuery{
		MediaType:     payload.MEDIA_TYPE_VIDEO_MP4,
		MediaCategory: payload.MEDIA_CATEGORY_DM_VIDEO,
	}

	mediaResult, err := cli.UploadMedia(uploadQuery, videoBytes)
	if err != nil {
		log.Fatal(err)
	}

	payload := &payload.SendDirectMessagePayload{
		ConversationID:    firstConversation.ConversationID,
		RequestID:         uuid.NewString(),
		CardsPlatform:     "Web-12",
		IncludeCards:      1,
		Text:              "",
		MediaID:           mediaResult.MediaIDString,
		IncludeQuoteCount: true,
		RecipientIds:      false,
		DmUsers:           false,
	}

	sentMessageResponse, err := cli.SendDirectMessage(payload)
	if err != nil {
		log.Fatalf("failed to send video to conversation with id %s (%s)", firstConversation.ConversationID, err.Error())
	}

	cli.Logger.Info().Any("response", sentMessageResponse).Str("conversation_id", firstConversation.ConversationID).Msg("Sent test video to first conversation")
}

func sendMessageTest(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}
	firstConversation := conversations[0].Conversation

	payload := &payload.SendDirectMessagePayload{
		ConversationID:    firstConversation.ConversationID,
		RequestID:         uuid.NewString(),
		Text:              "this is a test message",
		CardsPlatform:     "Web-12",
		IncludeCards:      1,
		IncludeQuoteCount: true,
		RecipientIds:      false,
		DmUsers:           false,
	}

	sentMessageResponse, err := cli.SendDirectMessage(payload)
	if err != nil {
		log.Fatalf("failed to send msg to conversation with id %s (%s)", firstConversation.ConversationID, err.Error())
	}

	cli.Logger.Info().Any("response", sentMessageResponse).Str("conversation_id", firstConversation.ConversationID).Msg("Sent test message to first conversation")
}

func logAllTrustedConversations(initialInboxData *response.XInboxData) {
	inboxTimelines := initialInboxData.InboxTimelines
	trustedInboxTimeline := inboxTimelines.Trusted

	paginationNextEntryID := trustedInboxTimeline.MinEntryID
	paginationStatus := trustedInboxTimeline.Status
	reqQuery := (&payload.DmRequestQuery{})

	for paginationStatus == types.HAS_MORE {
		reqQuery.MaxID = paginationNextEntryID
		nextInboxTimelineResponse, err := cli.FetchTrustedThreads(reqQuery)
		if err != nil {
			log.Fatal(err)
		}

		methods.MergeMaps(initialInboxData.Conversations, nextInboxTimelineResponse.InboxTimeline.Conversations)
		methods.MergeMaps(initialInboxData.Users, nextInboxTimelineResponse.InboxTimeline.Users)
		initialInboxData.Entries = append(initialInboxData.Entries, nextInboxTimelineResponse.InboxTimeline.Entries...)

		paginationNextEntryID = nextInboxTimelineResponse.InboxTimeline.MinEntryID
		paginationStatus = nextInboxTimelineResponse.InboxTimeline.Status
	}

	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}

	for i, c := range conversations {
		conv := c.Conversation
		mostRecentMessage := c.Messages[0]
		cli.Logger.Info().
			Int("conversation_inbox_position", i).
			Str("conversation_id", conv.ConversationID).
			Str("type", string(conv.Type)).
			Str("createdAt", conv.CreateTime).
			Str("createdByUserID", conv.CreatedByUserID).
			Any("participants", c.Participants).
			Any("most_recent_message", mostRecentMessage.MessageData.Text).
			Msg("Inbox Timeline Conversation")
	}
}

func logAllMessagesInConversation(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}

	firstConversation := conversations[0].Conversation

	conversationMessageHistoryStatus := firstConversation.Status
	if conversationMessageHistoryStatus == types.AT_END {
		log.Fatalf("conversation with id %s does not have any more messages to fetch", firstConversation.ConversationID)
	}

	totalMessages := len(conversations[0].Messages)
	paginationNextEntryID := firstConversation.MinEntryID
	reqQuery := (&payload.DmRequestQuery{}).Default()
	for conversationMessageHistoryStatus == types.HAS_MORE {
		reqQuery.MaxID = paginationNextEntryID
		fetchMessagesResponse, err := cli.FetchConversationContext(firstConversation.ConversationID, reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION_HISTORY)
		if err != nil {
			log.Fatal(err)
		}

		conversationTimeline := fetchMessagesResponse.ConversationTimeline
		messageBatch, err := conversationTimeline.PrettifyMessages(firstConversation.ConversationID)
		if err != nil {
			log.Fatalf("failed to prettify message batch for conversation with id %s (%s)", firstConversation.ConversationID, err.Error())
		}

		for _, msg := range messageBatch {
			cli.Logger.Info().
				Str("conversation_id", msg.ConversationID).
				Str("sender_name", msg.Sender.Name).
				Str("sender_screen_name", msg.Sender.ScreenName).
				Str("recipient_name", msg.Recipient.Name).
				Str("recipient_screen_name", msg.Recipient.ScreenName).
				Str("sent_at", msg.SentAt.String()).
				Str("text", msg.Text).
				Any("attachment", msg.Attachment).
				Any("entities", msg.Entities).
				Msg("Message")
		}

		totalMessages += len(messageBatch)
		conversationMessageHistoryStatus = conversationTimeline.Status
		paginationNextEntryID = conversationTimeline.MinEntryID
	}

	cli.Logger.Info().Int("total", totalMessages).Str("conversation_id", firstConversation.ConversationID).Msg("Successfully fetched all existing messages in conversation")
}

func logInitialDisplayFeed(initialInboxData *response.XInboxData) {
	conversations, err := initialInboxData.Prettify()
	if err != nil {
		log.Fatal(err)
	}

	for i, c := range conversations {
		conv := c.Conversation
		mostRecentMessage := c.Messages[0]
		cli.Logger.Info().
			Int("conversation_inbox_position", i).
			Str("conversation_id", conv.ConversationID).
			Str("type", string(conv.Type)).
			Str("createdAt", conv.CreateTime).
			Str("createdByUserID", conv.CreatedByUserID).
			Any("participants", c.Participants).
			Any("most_recent_message", mostRecentMessage.MessageData.Text).
			Msg("Initial Inbox Display")
	}
}

func eventHandler(evt interface{}) {
	switch evtData := evt.(type) {
	case event.XEventMessage:
		cli.Logger.Info().
			Str("conversation_id", evtData.Conversation.ConversationID).
			Str("sender_id", evtData.Sender.IDStr).
			Str("recipient_id", evtData.Recipient.IDStr).
			Str("message_id", evtData.MessageID).
			Str("createdAt", evtData.CreatedAt.String()).
			Str("text", evtData.Text).
			Any("entities", evtData.Entities).
			Any("attachment", evtData.Attachment).
			Msg("New message event!")
	case event.XEventConversationRead:
		cli.Logger.Info().
			Str("conversation_id", evtData.Conversation.ConversationID).
			Str("last_read_event_id", evtData.LastReadEventID).
			Str("read_at", evtData.ReadAt.String()).
			Msg("Conversation was read!")
	case event.XEventConversationCreated:
		cli.Logger.Info().
			Str("conversation_id", evtData.Conversation.ConversationID).
			Any("participants", evtData.Conversation.Participants).
			Str("type", string(evtData.Conversation.Type)).
			Str("created_at", evtData.CreatedAt.String()).
			Msg("New conversation was created!")
	case event.XEventMessageDeleted:
		cli.Logger.Info().
			Str("conversation_id", evtData.Conversation.ConversationID).
			Any("participants", evtData.Conversation.Participants).
			Any("messages_deleted", evtData.Messages).
			Str("type", string(evtData.Conversation.Type)).
			Str("deleted_at", evtData.DeletedAt.String()).
			Msg("Messages were deleted!")
	default:
		log.Println("unknown event:", evt)
	}
}
