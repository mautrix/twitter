package crypto

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

const (
	// SignatureVersion3 is the current signature version for MessageCreateEvent.
	SignatureVersion3 = "3"
	// SignatureVersion4 is the current signature version for MarkConversationReadEvent.
	SignatureVersion4 = "4"

	// SignatureSize is the size of a raw ECDSA P-256 signature (r || s).
	SignatureSize = 64
)

// SignaturePreimage builds the preimage for signature version 3.
// Format: "MessageCreateEvent,{message_id},{sender_id},{conversation_id},{key_version},{base64_nopad(contents)}"
func SignaturePreimage(messageID, senderID, conversationID, keyVersion string, contents []byte) []byte {
	// Use standard base64 without padding (matches Twitter's format)
	contentsB64 := base64.RawStdEncoding.EncodeToString(contents)
	preimage := fmt.Sprintf("MessageCreateEvent,%s,%s,%s,%s,%s",
		messageID, senderID, conversationID, keyVersion, contentsB64)
	return []byte(preimage)
}

// SignaturePreimageMarkConversationReadEvent builds the preimage for signature version 4.
func SignaturePreimageMarkConversationReadEvent(messageID, senderID, conversationID, conversationToken, createdAtMsec, seenUntilSequenceID string, seenAtMillis int64) []byte {
	preimage := fmt.Sprintf("MarkConversationReadEvent,%s,%s,%s,%s,%s,%s,%d",
		messageID, senderID, conversationID, conversationToken, createdAtMsec, seenUntilSequenceID, seenAtMillis)
	return []byte(preimage)
}

// SignaturePreimageMessageDeleteEvent builds the preimage for signature version 4 for delete events.
// Format: "MessageDeleteEvent,{message_id},{sender_id},{conversation_id},{conversation_token},{created_at_msec},{encoded_message_event_detail}"
func SignaturePreimageMessageDeleteEvent(messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail string) []byte {
	preimage := fmt.Sprintf("MessageDeleteEvent,%s,%s,%s,%s,%s,%s",
		messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail)
	return []byte(preimage)
}

// Sign creates an ECDSA P-256 signature over the given preimage.
// Returns the raw 64-byte signature (r || s), base64-encoded.
func Sign(privateKey *ecdsa.PrivateKey, preimage []byte) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is nil")
	}

	hash := sha256.Sum256(preimage)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("ecdsa sign: %w", err)
	}

	// Encode as raw 64 bytes: r (32) || s (32), both padded to 32 bytes
	sig := make([]byte, SignatureSize)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)

	// Use StdEncoding (with padding) to match Twitter's format
	return base64.StdEncoding.EncodeToString(sig), nil
}

// SignMessage creates a signature for a MessageCreateEvent.
// contents should be the raw ciphertext bytes (not hex-encoded).
func SignMessage(privateKey *ecdsa.PrivateKey, messageID, senderID, conversationID, keyVersion string, contents []byte) (string, error) {
	preimage := SignaturePreimage(messageID, senderID, conversationID, keyVersion, contents)
	return Sign(privateKey, preimage)
}

// SignMarkConversationReadEvent creates a signature for a MarkConversationReadEvent.
func SignMarkConversationReadEvent(privateKey *ecdsa.PrivateKey, messageID, senderID, conversationID, conversationToken, createdAtMsec, seenUntilSequenceID string, seenAtMillis int64) (string, error) {
	preimage := SignaturePreimageMarkConversationReadEvent(messageID, senderID, conversationID, conversationToken, createdAtMsec, seenUntilSequenceID, seenAtMillis)
	return Sign(privateKey, preimage)
}

// SignMessageDeleteEvent creates a signature for a MessageDeleteEvent.
func SignMessageDeleteEvent(privateKey *ecdsa.PrivateKey, messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail string) (string, error) {
	preimage := SignaturePreimageMessageDeleteEvent(messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail)
	return Sign(privateKey, preimage)
}

// SignaturePreimageMuteConversation builds the preimage for signature version 4 for mute events.
func SignaturePreimageMuteConversation(messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail string) []byte {
	preimage := fmt.Sprintf("MuteConversation,%s,%s,%s,%s,%s,%s",
		messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail)
	return []byte(preimage)
}

// SignMuteConversation creates a signature for a MuteConversation event.
func SignMuteConversation(privateKey *ecdsa.PrivateKey, messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail string) (string, error) {
	preimage := SignaturePreimageMuteConversation(messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail)
	return Sign(privateKey, preimage)
}

// SignaturePreimageUnmuteConversation builds the preimage for signature version 4 for unmute events.
func SignaturePreimageUnmuteConversation(messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail string) []byte {
	preimage := fmt.Sprintf("UnmuteConversation,%s,%s,%s,%s,%s,%s",
		messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail)
	return []byte(preimage)
}

// SignUnmuteConversation creates a signature for an UnmuteConversation event.
func SignUnmuteConversation(privateKey *ecdsa.PrivateKey, messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail string) (string, error) {
	preimage := SignaturePreimageUnmuteConversation(messageID, senderID, conversationID, conversationToken, createdAtMsec, encodedMessageEventDetail)
	return Sign(privateKey, preimage)
}

// Verify verifies an ECDSA P-256 signature.
// signatureB64 is the raw 64-byte signature, base64-encoded.
func Verify(publicKey *ecdsa.PublicKey, preimage []byte, signatureB64 string) error {
	if publicKey == nil {
		return errors.New("public key is nil")
	}

	sigBytes, err := decodeBase64Flexible(signatureB64)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	if len(sigBytes) != SignatureSize {
		return fmt.Errorf("signature must be %d bytes, got %d", SignatureSize, len(sigBytes))
	}

	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	hash := sha256.Sum256(preimage)

	if !ecdsa.Verify(publicKey, hash[:], r, s) {
		return ErrSignatureInvalid
	}

	return nil
}

// VerifyMessage verifies a signature for a MessageCreateEvent.
// contents should be the raw ciphertext bytes (not hex-encoded).
func VerifyMessage(publicKey *ecdsa.PublicKey, messageID, senderID, conversationID, keyVersion string, contents []byte, signatureB64 string) error {
	preimage := SignaturePreimage(messageID, senderID, conversationID, keyVersion, contents)
	return Verify(publicKey, preimage, signatureB64)
}

// VerifyMessageEvent verifies the signature on a MessageEvent.
// If km is provided, it will be used to look up/cache the public key.
// Returns nil if verification succeeds.
func VerifyMessageEvent(ctx context.Context, km *KeyManager, event *payload.MessageEvent) error {
	if event == nil {
		return errors.New("event is nil")
	}
	if event.MessageEventSignature == nil {
		return errors.New("no signature present")
	}

	sig := event.MessageEventSignature
	if sig.SignatureVersion == nil || *sig.SignatureVersion != SignatureVersion3 {
		return fmt.Errorf("unsupported signature version: %v", sig.SignatureVersion)
	}

	if sig.Signature == nil {
		return errors.New("signature is nil")
	}

	if event.Detail == nil || event.Detail.MessageCreateEvent == nil {
		return errors.New("not a message create event")
	}

	mce := event.Detail.MessageCreateEvent
	if len(mce.Contents) == 0 || mce.ConversationKeyVersion == nil {
		return errors.New("missing contents or key version")
	}

	// Contents is already raw bytes
	contents := mce.Contents

	// Get required fields
	messageID := ""
	if event.MessageId != nil {
		messageID = *event.MessageId
	}

	senderID := ""
	if event.SenderId != nil {
		senderID = *event.SenderId
	}

	conversationID := ""
	if event.ConversationId != nil {
		conversationID = *event.ConversationId
	}

	keyVersion := ""
	if sig.PublicKeyVersion != nil {
		keyVersion = *sig.PublicKeyVersion
	}

	spki := ""
	if sig.SigningPublicKey != nil {
		spki = *sig.SigningPublicKey
	}

	// Get the public key
	var pubKey *ecdsa.PublicKey
	var err error
	if km != nil {
		pubKey, err = km.GetPublicKeyForVerification(ctx, senderID, keyVersion, spki)
	} else if spki != "" {
		pubKey, err = ParsePublicKeySPKI(spki)
	} else {
		return errors.New("no public key available for verification")
	}
	if err != nil {
		return fmt.Errorf("get public key: %w", err)
	}

	return VerifyMessage(pubKey, messageID, senderID, conversationID, *mce.ConversationKeyVersion, contents, *sig.Signature)
}
