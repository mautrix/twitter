package twittermeow

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func sequenceTestDeleteEvent(sequenceID string) *payload.MessageEvent {
	return sequenceTestDeleteEventForConversation(sequenceID, "conversation")
}

func sequenceTestDeleteEventForConversation(sequenceID, conversationID string) *payload.MessageEvent {
	return &payload.MessageEvent{
		SequenceId:     ptr.Ptr(sequenceID),
		ConversationId: ptr.Ptr(conversationID),
		Detail: &payload.MessageEventDetail{
			MessageDeleteEvent: &payload.MessageDeleteEvent{SequenceIds: []string{"message"}},
		},
	}
}

func encodeSequenceTestEvent(t *testing.T, evt *payload.MessageEvent) string {
	t.Helper()
	encoded, err := payload.Encode(evt)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	return base64.StdEncoding.EncodeToString(encoded)
}

func TestCompareXChatSequenceIDsSupportsLargeDecimalValues(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{a: "18446744073709551615", b: "9223372036854775807", want: 1},
		{a: "9223372036854775808", b: "18446744073709551615", want: -1},
		{a: "00042", b: "42", want: 0},
		{a: "invalid", b: "42", want: -1},
	}
	for _, test := range tests {
		if got := compareXChatSequenceIDs(test.a, test.b); got != test.want {
			t.Errorf("compareXChatSequenceIDs(%q, %q) = %d, want %d", test.a, test.b, got, test.want)
		}
	}
}

func TestXChatProcessMessageStopsAfterFirstHandlerFailure(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	handlerCalls := 0
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		handlerCalls++
		return handlerCalls != 1
	})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})

	err := processor.ProcessMessage(context.Background(), &payload.Message{
		BatchedMessageEvents: &payload.BatchedMessageEvents{
			MessageEvents: []*payload.MessageEvent{
				sequenceTestDeleteEvent("211"),
				sequenceTestDeleteEvent("212"),
			},
		},
	})
	if err == nil {
		t.Fatal("ProcessMessage() error = nil, want rejected handler error")
	}
	if handlerCalls != 1 {
		t.Fatalf("handler calls = %d, want 1", handlerCalls)
	}
	if len(tracked) != 0 {
		t.Fatalf("tracked sequences = %v, want none", tracked)
	}
}

func TestXChatProcessMessageDoesNotPublishOutOfOrderSuccessBeforeFailure(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	handlerCalls := 0
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		handlerCalls++
		return handlerCalls != 2
	})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})

	err := processor.ProcessMessage(context.Background(), &payload.Message{
		BatchedMessageEvents: &payload.BatchedMessageEvents{
			MessageEvents: []*payload.MessageEvent{
				sequenceTestDeleteEvent("900"),
				sequenceTestDeleteEvent("800"),
			},
		},
	})
	if err == nil {
		t.Fatal("ProcessMessage() error = nil, want rejected handler error")
	}
	if len(tracked) != 0 {
		t.Fatalf("tracked sequences = %v, want none", tracked)
	}

	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool { return true })
	if err := processor.ProcessMessage(context.Background(), &payload.Message{
		MessageEvent: sequenceTestDeleteEvent("800"),
	}); err != nil {
		t.Fatalf("retry ProcessMessage() error = %v", err)
	}
	if len(tracked) != 1 || tracked[0] != "900" {
		t.Fatalf("tracked sequences after retry = %v, want [900]", tracked)
	}
}

func TestXChatSequenceAdvancesOnlyAfterHandlerSuccess(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		return false
	})

	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("101")); err == nil {
		t.Fatal("processMessageEvent() error = nil, want rejected handler error")
	}
	if len(tracked) != 0 {
		t.Fatalf("tracked sequences after rejected event = %v, want none", tracked)
	}

	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		return true
	})
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("101")); err != nil {
		t.Fatalf("processMessageEvent() error = %v", err)
	}
	if len(tracked) != 1 || tracked[0] != "101" {
		t.Fatalf("tracked sequences = %v, want [101]", tracked)
	}
}

func TestXChatProcessMessagePropagatesHandlerFailure(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		return false
	})
	var tracked bool
	processor.SetSequenceIDCallback(func(string) {
		tracked = true
	})

	err := processor.ProcessMessage(context.Background(), &payload.Message{
		MessageEvent: sequenceTestDeleteEvent("202"),
	})
	if err == nil {
		t.Fatal("ProcessMessage() error = nil, want rejected handler error")
	}
	if tracked {
		t.Fatal("failed event advanced the sequence checkpoint")
	}
}

func TestXChatProcessMessagePanicDoesNotAdvanceCheckpoint(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetGapHandler(func(context.Context, string, string, string) error { return nil })
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		panic("simulated handler panic")
	})
	var tracked bool
	processor.SetSequenceIDCallback(func(string) {
		tracked = true
	})

	err := processor.ProcessMessage(context.Background(), &payload.Message{
		MessageEvent: sequenceTestDeleteEvent("203"),
	})
	if err == nil {
		t.Fatal("ProcessMessage() error = nil, want recovered panic error")
	}
	if tracked {
		t.Fatal("panicked event advanced the sequence checkpoint")
	}
	if !processor.ConversationGapUnresolved("conversation") {
		t.Fatal("panicked event did not leave its conversation retryable")
	}
}

func TestXChatIgnoredEventStillAdvancesSequence(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var tracked string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = sequenceID
	})

	if err := processor.processMessageEvent(context.Background(), &payload.MessageEvent{
		SequenceId: ptr.Ptr("303"),
	}); err != nil {
		t.Fatalf("processMessageEvent() error = %v", err)
	}
	if tracked != "303" {
		t.Fatalf("tracked sequence = %q, want 303", tracked)
	}
}

func TestXChatMalformedMessageDoesNotAdvanceSequence(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})
	evt := &payload.MessageEvent{
		SequenceId:     ptr.Ptr("304"),
		ConversationId: ptr.Ptr("conversation"),
		Detail: &payload.MessageEventDetail{
			MessageCreateEvent: &payload.MessageCreateEvent{Contents: []byte{0xff}},
		},
	}

	if err := processor.processMessageEvent(context.Background(), evt); err == nil {
		t.Fatal("processMessageEvent() error = nil, want malformed message error")
	}
	if len(tracked) != 0 {
		t.Fatalf("tracked sequences = %v, want none", tracked)
	}
}

func TestXChatEmptyMessageCreateEmitsConversationCreate(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var emitted types.TwitterEvent
	processor.SetEventHandler(func(_ context.Context, evt types.TwitterEvent) bool {
		emitted = evt
		return true
	})
	var tracked string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = sequenceID
	})

	evt := &payload.MessageEvent{
		SequenceId:     ptr.Ptr("305"),
		ConversationId: ptr.Ptr("conversation"),
		Detail: &payload.MessageEventDetail{
			MessageCreateEvent: &payload.MessageCreateEvent{},
		},
	}
	if err := processor.processMessageEvent(context.Background(), evt); err != nil {
		t.Fatalf("processMessageEvent() error = %v", err)
	}
	create, ok := emitted.(*types.ConversationCreate)
	if !ok || create.ConversationID != "conversation" || create.ID != "305" {
		t.Fatalf("emitted event = %#v, want ConversationCreate for sequence 305", emitted)
	}
	if tracked != "305" {
		t.Fatalf("tracked sequence = %q, want 305", tracked)
	}
}

func TestXChatInboxDecodeFailureDoesNotAdvanceSequence(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetGapHandler(func(context.Context, string, string, string) error { return nil })
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})

	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail:  response.XChatConversationDetail{ConversationID: "conversation"},
		LatestMessageEvents: []string{"not-base64"},
	})
	if err == nil {
		t.Fatal("ProcessMessageAndReadEvents() error = nil, want decode error")
	}
	if len(tracked) != 0 {
		t.Fatalf("tracked sequences = %v, want none", tracked)
	}
	if !processor.SequenceCheckpointBlocked() {
		t.Fatal("decode failure did not block the conversation checkpoint")
	}
}

func TestXChatInboxDecodeFailureStillProcessesValidEvent(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetGapHandler(func(context.Context, string, string, string) error { return nil })
	handlerCalls := 0
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		handlerCalls++
		return true
	})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})

	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{ConversationID: "conversation"},
		LatestMessageEvents: []string{
			"not-base64",
			encodeSequenceTestEvent(t, sequenceTestDeleteEvent("307")),
		},
	})
	if err == nil {
		t.Fatal("ProcessMessageAndReadEvents() error = nil, want decode error")
	}
	if handlerCalls != 1 {
		t.Fatalf("handler calls = %d, want 1 valid event delivered", handlerCalls)
	}
	if len(tracked) != 0 || !processor.SequenceCheckpointBlocked() {
		t.Fatalf("tracked/blocked = %v/%t, want no checkpoint and blocked retry", tracked, processor.SequenceCheckpointBlocked())
	}
}

func TestXChatInboxValidEventDoesNotClearPriorUnresolvedState(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	gapCalls := 0
	processor.SetGapHandler(func(context.Context, string, string, string) error {
		gapCalls++
		return nil
	})
	processor.MarkConversationGapUnresolved("conversation")
	handlerCalls := 0
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		handlerCalls++
		return true
	})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})

	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{ConversationID: "conversation"},
		LatestMessageEvents: []string{
			encodeSequenceTestEvent(t, sequenceTestDeleteEvent("308")),
		},
	})
	if err != nil {
		t.Fatalf("ProcessMessageAndReadEvents() error = %v", err)
	}
	if handlerCalls != 1 || gapCalls != 0 {
		t.Fatalf("handler/gap calls = %d/%d, want 1/0", handlerCalls, gapCalls)
	}
	if len(tracked) != 0 || !processor.ConversationGapUnresolved("conversation") {
		t.Fatalf("tracked/unresolved = %v/%t, want no checkpoint and unresolved", tracked, processor.ConversationGapUnresolved("conversation"))
	}
}

func TestXChatInboxMalformedReadReceiptDoesNotBlockMessages(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetGapHandler(func(context.Context, string, string, string) error { return nil })

	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{ConversationID: "conversation"},
		LatestReadEventsPerParticipant: []response.XChatParticipantReadEvent{{
			LatestMarkConversationReadEvent: "not-base64",
		}},
	})
	if err != nil {
		t.Fatalf("ProcessMessageAndReadEvents() error = %v", err)
	}
	if processor.SequenceCheckpointBlocked() {
		t.Fatal("malformed ephemeral read receipt blocked the message checkpoint")
	}
}

func TestXChatInboxProcessesLatestNotifiableMessage(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetGapHandler(func(context.Context, string, string, string) error { return nil })
	handlerCalls := 0
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		handlerCalls++
		return true
	})

	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail:                 response.XChatConversationDetail{ConversationID: "conversation"},
		LatestNotifiableMessageCreateEvent: encodeSequenceTestEvent(t, sequenceTestDeleteEvent("306")),
	})
	if err != nil {
		t.Fatalf("ProcessMessageAndReadEvents() error = %v", err)
	}
	if handlerCalls != 1 {
		t.Fatalf("handler calls = %d, want 1", handlerCalls)
	}
}

func TestXChatInboxFillsMissingEventConversationID(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetEventHandler(func(_ context.Context, evt types.TwitterEvent) bool {
		deleted, ok := evt.(*types.MessageDelete)
		return ok && deleted.ConversationID == "conversation"
	})
	evt := sequenceTestDeleteEvent("309")
	evt.ConversationId = nil
	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail:  response.XChatConversationDetail{ConversationID: "conversation"},
		LatestMessageEvents: []string{encodeSequenceTestEvent(t, evt)},
	})
	if err != nil {
		t.Fatalf("ProcessMessageAndReadEvents() error = %v", err)
	}
}

func TestXChatInboxStopsConversationAfterFirstHandlerFailure(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	processor.SetGapHandler(func(context.Context, string, string, string) error { return nil })
	handlerCalls := 0
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		handlerCalls++
		return false
	})

	err := processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{ConversationID: "conversation"},
		LatestMessageEvents: []string{
			encodeSequenceTestEvent(t, sequenceTestDeleteEvent("310")),
			encodeSequenceTestEvent(t, sequenceTestDeleteEvent("311")),
		},
	})
	if err == nil {
		t.Fatal("ProcessMessageAndReadEvents() error = nil, want handler error")
	}
	if handlerCalls != 1 {
		t.Fatalf("handler calls = %d, want 1", handlerCalls)
	}
	if !processor.SequenceCheckpointBlocked() {
		t.Fatal("handler failure did not block the conversation checkpoint")
	}
}

func TestXChatReconnectCatchupRunsBeforeFirstConversationEvent(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var order []string
	processor.SetGapHandler(func(_ context.Context, conversationID, _, currentSequenceID string) error {
		order = append(order, "catchup:"+conversationID+":"+currentSequenceID)
		return nil
	})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		order = append(order, "event")
		return true
	})
	processor.MarkReconnected()

	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("401")); err != nil {
		t.Fatalf("first event error = %v", err)
	}
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("402")); err != nil {
		t.Fatalf("second event error = %v", err)
	}
	want := []string{"catchup:conversation:401", "event", "event"}
	if len(order) != len(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order = %v, want %v", order, want)
		}
	}
}

func TestXChatReconnectCatchupFailureDoesNotDropLiveEvent(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var catchupCalls, eventCalls int
	var tracked []string
	catchupAvailable := false
	processor.SetGapHandler(func(context.Context, string, string, string) error {
		catchupCalls++
		if !catchupAvailable {
			return errors.New("catch-up unavailable")
		}
		return nil
	})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		eventCalls++
		return true
	})
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})
	processor.MarkReconnected()

	for i := 0; i < 2; i++ {
		if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("501")); err != nil {
			t.Fatalf("event error = %v, want live event to continue", err)
		}
	}
	if catchupCalls != 2 {
		t.Fatalf("catch-up calls = %d, want 2", catchupCalls)
	}
	if eventCalls != 2 || len(tracked) != 0 {
		t.Fatalf("event/tracked calls = %d/%d, want 2/0", eventCalls, len(tracked))
	}

	catchupAvailable = true
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("502")); err != nil {
		t.Fatalf("recovery event error = %v", err)
	}
	if catchupCalls != 3 || eventCalls != 3 {
		t.Fatalf("catch-up/event calls after recovery = %d/%d, want 3/3", catchupCalls, eventCalls)
	}
	if len(tracked) != 1 || tracked[0] != "502" {
		t.Fatalf("tracked sequences after recovery = %v, want [502]", tracked)
	}

	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("503")); err != nil {
		t.Fatalf("post-recovery event error = %v", err)
	}
	if catchupCalls != 3 || eventCalls != 4 {
		t.Fatalf("catch-up/event calls after completed recovery = %d/%d, want 3/4", catchupCalls, eventCalls)
	}
	if len(tracked) != 2 || tracked[1] != "503" {
		t.Fatalf("tracked sequences after completed recovery = %v, want [502 503]", tracked)
	}
}

func TestXChatUnresolvedGapBlocksOtherConversationCheckpoint(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	catchupAvailable := false
	processor.SetGapHandler(func(_ context.Context, conversationID, _, _ string) error {
		if conversationID == "conversation-a" && !catchupAvailable {
			return errors.New("catch-up unavailable")
		}
		return nil
	})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		return true
	})
	var tracked []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		tracked = append(tracked, sequenceID)
	})
	processor.MarkReconnected()

	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEventForConversation("701", "conversation-a")); err != nil {
		t.Fatalf("gapped event error = %v", err)
	}
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEventForConversation("900", "conversation-b")); err != nil {
		t.Fatalf("other conversation event error = %v", err)
	}
	if len(tracked) != 0 {
		t.Fatalf("tracked sequences while gap unresolved = %v, want none", tracked)
	}

	catchupAvailable = true
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEventForConversation("702", "conversation-a")); err != nil {
		t.Fatalf("recovery event error = %v", err)
	}
	if len(tracked) != 1 || tracked[0] != "900" {
		t.Fatalf("tracked sequences after recovery = %v, want [900]", tracked)
	}
}

func TestXChatSequenceCallbackCanBeReplacedDuringBlockedGap(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	catchupAvailable := false
	processor.SetGapHandler(func(context.Context, string, string, string) error {
		if !catchupAvailable {
			return errors.New("catch-up unavailable")
		}
		return nil
	})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool { return true })
	processor.MarkReconnected()

	var firstCallback, secondCallback []string
	processor.SetSequenceIDCallback(func(sequenceID string) {
		firstCallback = append(firstCallback, sequenceID)
	})
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("801")); err != nil {
		t.Fatalf("blocked event error = %v", err)
	}
	processor.SetSequenceIDCallback(func(sequenceID string) {
		secondCallback = append(secondCallback, sequenceID)
	})
	catchupAvailable = true
	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("802")); err != nil {
		t.Fatalf("recovery event error = %v", err)
	}

	if len(firstCallback) != 0 {
		t.Fatalf("first callback sequences = %v, want none", firstCallback)
	}
	if len(secondCallback) != 1 || secondCallback[0] != "802" {
		t.Fatalf("second callback sequences = %v, want [802]", secondCallback)
	}
}

func TestXChatSequenceDiscontinuityTriggersConversationCatchup(t *testing.T) {
	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	catchupCalls := 0
	processor.SetGapHandler(func(context.Context, string, string, string) error {
		catchupCalls++
		return nil
	})
	processor.SetEventHandler(func(context.Context, types.TwitterEvent) bool {
		return true
	})

	if err := processor.processMessageEvent(context.Background(), sequenceTestDeleteEvent("601")); err != nil {
		t.Fatalf("initial event error = %v", err)
	}
	gapped := sequenceTestDeleteEvent("603")
	gapped.PreviousSequenceId = ptr.Ptr("602")
	if err := processor.processMessageEvent(context.Background(), gapped); err != nil {
		t.Fatalf("gapped event error = %v", err)
	}
	if catchupCalls != 1 {
		t.Fatalf("catch-up calls = %d, want 1", catchupCalls)
	}
}
