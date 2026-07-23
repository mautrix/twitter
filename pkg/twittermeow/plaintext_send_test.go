package twittermeow

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

func TestBuildUnencryptedMessageMutationPayload(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate signing key: %v", err)
	}
	signingKey := &crypto.SigningKeyPair{
		KeyVersion: "1781660395506",
		SigningKey: privateKey,
	}
	opts := SendEncryptedMessageOpts{
		ConversationID: "g1709621683324379335",
		MessageID:      "16eec3e2-a7cc-479a-b9af-67b317ba03a8",
		Text:           "Our depth is pretty insane",
	}

	pl, err := buildUnencryptedMessageMutationPayload(
		context.Background(),
		opts,
		"1739875480390750208",
		signingKey,
	)
	if err != nil {
		t.Fatalf("build payload: %v", err)
	}
	if got := pl.Variables.ConversationID; got != opts.ConversationID {
		t.Fatalf("conversation ID = %q, want %q", got, opts.ConversationID)
	}
	if got := pl.Variables.MessageID; got != opts.MessageID {
		t.Fatalf("message ID = %q, want %q", got, opts.MessageID)
	}

	encodedMCE, err := base64.StdEncoding.DecodeString(pl.Variables.EncodedMessageCreateEvent)
	if err != nil {
		t.Fatalf("decode message create event: %v", err)
	}
	var mce payload.MessageCreateEvent
	if err = payload.Decode(encodedMCE, &mce); err != nil {
		t.Fatalf("decode thrift message create event: %v", err)
	}
	if mce.ConversationKeyVersion != nil {
		t.Fatalf("plaintext message has conversation key version %q", ptr.Val(mce.ConversationKeyVersion))
	}
	entry, err := crypto.ParseMessageEntryContentsBytes(mce.Contents)
	if err != nil {
		t.Fatalf("parse plaintext message contents: %v", err)
	}
	if got := ptr.Val(entry.Message.MessageText); got != opts.Text {
		t.Fatalf("message text = %q, want %q", got, opts.Text)
	}

	if pl.Variables.EncodedMessageEventSignature == nil {
		t.Fatal("encoded message signature is nil")
	}
	encodedSig, err := base64.StdEncoding.DecodeString(*pl.Variables.EncodedMessageEventSignature)
	if err != nil {
		t.Fatalf("decode message signature: %v", err)
	}
	var sig payload.MessageEventSignature
	if err = payload.Decode(encodedSig, &sig); err != nil {
		t.Fatalf("decode thrift message signature: %v", err)
	}
	if got := ptr.Val(sig.PublicKeyVersion); got != signingKey.KeyVersion {
		t.Fatalf("public key version = %q, want %q", got, signingKey.KeyVersion)
	}
	if err = crypto.VerifyMessage(
		&privateKey.PublicKey,
		opts.MessageID,
		"1739875480390750208",
		opts.ConversationID,
		"",
		mce.Contents,
		ptr.Val(sig.Signature),
	); err != nil {
		t.Fatalf("verify plaintext signature: %v", err)
	}
}

func TestMessageCreateEventEncryptionDetection(t *testing.T) {
	keyVersion := "1781660395506"
	tests := []struct {
		name string
		mce  *payload.MessageCreateEvent
		want bool
	}{
		{
			name: "signed plaintext",
			mce:  &payload.MessageCreateEvent{},
			want: false,
		},
		{
			name: "encrypted with key version",
			mce:  &payload.MessageCreateEvent{ConversationKeyVersion: &keyVersion},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isEncryptedMessageCreateEvent(test.mce); got != test.want {
				t.Fatalf("isEncryptedMessageCreateEvent() = %t, want %t", got, test.want)
			}
		})
	}
}
