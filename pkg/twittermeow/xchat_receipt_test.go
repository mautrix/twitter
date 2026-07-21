package twittermeow

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func TestBuildXChatReadReceiptEventUsesCurrentWireFormat(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate signing key: %v", err)
	}
	keyPair := &crypto.SigningKeyPair{
		KeyVersion: "1768512266300",
		SigningKey: privateKey,
	}
	createdAt := time.UnixMilli(1784086879919)
	seenAt := time.UnixMilli(1784086876864)

	evt, err := buildXChatReadReceiptEvent(
		"a674feb1-3a80-4f26-8a24-003c28edccf2",
		"1892329226424225792",
		"1155463061127467008:1892329226424225792",
		"conversation-token",
		"2077236977481687040",
		createdAt,
		seenAt,
		keyPair,
	)
	if err != nil {
		t.Fatalf("build read receipt: %v", err)
	}

	if got := ptr.Val(evt.CreatedAtMsec); got != "1784086879919" {
		t.Fatalf("CreatedAtMsec = %q", got)
	}
	if got := ptr.Val(evt.PreviousSequenceId); got != "2077236977481687040" {
		t.Fatalf("PreviousSequenceId = %q", got)
	}
	if got := ptr.Val(evt.RelaySource); got != 0 {
		t.Fatalf("RelaySource = %d", got)
	}
	read := evt.Detail.MarkConversationReadEvent
	if got := ptr.Val(read.SeenAtMillis); got != 1784086876864 {
		t.Fatalf("SeenAtMillis = %d", got)
	}
	if read.IsGrok == nil || *read.IsGrok {
		t.Fatalf("IsGrok = %v, want false", read.IsGrok)
	}
	sig := evt.MessageEventSignature
	if got := ptr.Val(sig.SignatureVersion); got != crypto.SignatureVersion7 {
		t.Fatalf("SignatureVersion = %q", got)
	}
	if ptr.Val(sig.SigningPublicKey) == "" {
		t.Fatal("SigningPublicKey is empty")
	}
	preimage := crypto.SignaturePreimageMarkConversationReadEvent(
		ptr.Val(evt.MessageId),
		ptr.Val(evt.SenderId),
		ptr.Val(evt.ConversationId),
		ptr.Val(read.SeenUntilSequenceId),
		ptr.Val(read.SeenAtMillis),
	)
	wantPreimage := "MarkConversationReadEvent,a674feb1-3a80-4f26-8a24-003c28edccf2,1892329226424225792,1155463061127467008:1892329226424225792,2077236977481687040,1784086876864"
	if string(preimage) != wantPreimage {
		t.Fatalf("signature preimage = %q", preimage)
	}
	if err := crypto.Verify(&privateKey.PublicKey, preimage, ptr.Val(sig.Signature)); err != nil {
		t.Fatalf("verify signature: %v", err)
	}
}

func TestConvertXChatMarkReadEventPreservesSenderAndSeenTime(t *testing.T) {
	seenAt := int64(1784081234567)
	evt := &payload.MessageEvent{
		SequenceId:     ptr.Ptr("2077229000000000000"),
		SenderId:       ptr.Ptr("1892329226424225792"),
		ConversationId: ptr.Ptr("1155463061127467008:1892329226424225792"),
		CreatedAtMsec:  ptr.Ptr("1784081234000"),
	}
	read := &payload.MarkConversationReadEvent{
		SeenUntilSequenceId: ptr.Ptr("2077228999999999999"),
		SeenAtMillis:        &seenAt,
	}

	got := convertXChatMarkReadEvent(evt, read)
	if got.SenderID != "1892329226424225792" {
		t.Fatalf("SenderID = %q, want remote participant", got.SenderID)
	}
	if got.Time != "1784081234567" {
		t.Fatalf("Time = %q, want seen_at_millis", got.Time)
	}
	if got.LastReadEventID != "2077228999999999999" {
		t.Fatalf("LastReadEventID = %q", got.LastReadEventID)
	}
}

func TestProcessMessageAndReadEventsUsesParticipantSenderFallback(t *testing.T) {
	evt := &payload.MessageEvent{
		SequenceId:     ptr.Ptr("2077229000000000000"),
		ConversationId: ptr.Ptr("1155463061127467008:1892329226424225792"),
		CreatedAtMsec:  ptr.Ptr("1784081234000"),
		Detail: &payload.MessageEventDetail{
			MarkConversationReadEvent: &payload.MarkConversationReadEvent{
				SeenUntilSequenceId: ptr.Ptr("2077228999999999999"),
			},
		},
	}
	encoded, err := payload.Encode(evt)
	if err != nil {
		t.Fatalf("encode read event: %v", err)
	}

	processor := newXChatEventProcessor(&Client{Logger: zerolog.Nop()})
	var got *types.ConversationRead
	processor.SetEventHandler(func(_ context.Context, evt types.TwitterEvent) bool {
		got, _ = evt.(*types.ConversationRead)
		return true
	})
	err = processor.ProcessMessageAndReadEvents(context.Background(), &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{
			ConversationID: "1155463061127467008:1892329226424225792",
		},
		LatestReadEventsPerParticipant: []response.XChatParticipantReadEvent{{
			ParticipantID:                   response.XChatParticipantID{RestID: "1892329226424225792"},
			LatestMarkConversationReadEvent: base64.StdEncoding.EncodeToString(encoded),
		}},
	})
	if err != nil {
		t.Fatalf("ProcessMessageAndReadEvents: %v", err)
	}
	if got == nil {
		t.Fatal("read event was not emitted")
	}
	if got.SenderID != "1892329226424225792" {
		t.Fatalf("SenderID = %q, want participant fallback", got.SenderID)
	}
}
