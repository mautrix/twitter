package crypto

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	thrifter "github.com/thrift-iterator/go"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

// DecryptMessageContentsBytes decrypts raw ciphertext bytes to MessageContents.
func DecryptMessageContentsBytes(ciphertext []byte, conversationKey []byte) (*payload.MessageContents, error) {
	return DecryptMessageContentsBytesDebug(ciphertext, conversationKey, nil)
}

// DecryptMessageContentsBytesDebug decrypts raw ciphertext bytes to MessageContents with detailed logging.
func DecryptMessageContentsBytesDebug(ciphertext []byte, conversationKey []byte, log *zerolog.Logger) (*payload.MessageContents, error) {
	entry, err := DecryptMessageEntryContentsBytesDebug(ciphertext, conversationKey, log)
	if err != nil {
		return nil, err
	}
	if entry == nil || entry.Message == nil {
		if log != nil {
			log.Debug().
				Bool("contents_nil", entry == nil).
				Bool("message_nil", entry == nil || entry.Message == nil).
				Msg("No message contents in holder")
		}
		return nil, fmt.Errorf("no message contents in holder")
	}
	return entry.Message, nil
}

// DecryptMessageEntryContentsBytes decrypts ciphertext bytes to MessageEntryContents (message or reactions).
func DecryptMessageEntryContentsBytes(ciphertext []byte, conversationKey []byte) (*payload.MessageEntryContents, error) {
	return DecryptMessageEntryContentsBytesDebug(ciphertext, conversationKey, nil)
}

// DecryptMessageEntryContentsBytesDebug decrypts ciphertext bytes to MessageEntryContents (message/reaction/etc) with logging.
func DecryptMessageEntryContentsBytesDebug(ciphertext []byte, conversationKey []byte, log *zerolog.Logger) (*payload.MessageEntryContents, error) {
	if log != nil {
		log.Debug().
			Int("ciphertext_len", len(ciphertext)).
			Str("ciphertext_prefix", truncateHexBytes(ciphertext, 32)).
			Int("key_len", len(conversationKey)).
			Str("key_prefix", truncateHexBytes(conversationKey, 8)).
			Msg("Attempting secretbox decryption")
	}

	plaintext, err := SecretboxDecrypt(ciphertext, conversationKey)
	if err != nil {
		if log != nil {
			log.Debug().Err(err).Msg("Secretbox decryption failed")
		}
		return nil, fmt.Errorf("secretbox decrypt: %w", err)
	}

	if log != nil {
		log.Debug().
			Int("plaintext_len", len(plaintext)).
			Str("plaintext_hex", truncateHexBytes(plaintext, 64)).
			Msg("Secretbox decryption succeeded, attempting thrift decode")
	}

	entry, err := decodeMessageEntryHolder(plaintext, log)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("no message entry contents in holder")
	}
	if log != nil && entry.Message != nil {
		if msgJSON, err := json.Marshal(entry.Message); err == nil {
			evt := log.Info().
				Int("decrypted_json_len", len(msgJSON)).
				Int("plaintext_len", len(plaintext))
			if len(msgJSON) <= 4000 {
				evt = evt.RawJSON("decrypted_message", msgJSON)
			} else {
				evt = evt.Str("decrypted_message_prefix", string(msgJSON[:4000]))
			}
			evt.Msg("Decrypted XChat message contents")
		}
	}
	return entry, nil
}

// logHolderContents logs details about the decoded MessageEntryHolder structure.
func logHolderContents(log *zerolog.Logger, holder *payload.MessageEntryHolder, plaintext []byte) {
	evt := log.Debug().Bool("holder_contents_nil", holder.Contents == nil)

	if holder.Contents != nil {
		c := holder.Contents
		evt = evt.
			Bool("has_message", c.Message != nil).
			Bool("has_reaction_add", c.ReactionAdd != nil).
			Bool("has_reaction_remove", c.ReactionRemove != nil).
			Bool("has_message_edit", c.MessageEdit != nil).
			Bool("has_mark_read", c.MarkConversationRead != nil).
			Bool("has_mark_unread", c.MarkConversationUnread != nil).
			Bool("has_pin", c.PinConversation != nil).
			Bool("has_unpin", c.UnpinConversation != nil).
			Bool("has_screen_capture", c.ScreenCaptureDetected != nil).
			Bool("has_av_call_ended", c.AvCallEnded != nil).
			Bool("has_av_call_missed", c.AvCallMissed != nil).
			Bool("has_draft", c.DraftMessage != nil).
			Bool("has_accept_request", c.AcceptMessageRequest != nil).
			Bool("has_nickname", c.NicknameMessage != nil).
			Bool("has_verified_status", c.SetVerifiedStatus != nil).
			Bool("has_av_call_started", c.AvCallStarted != nil)

		if c.Message != nil {
			msg := c.Message
			hasText := msg.MessageText != nil && *msg.MessageText != ""
			evt = evt.
				Bool("message_has_text", hasText).
				Int("message_attachments", len(msg.Attachments)).
				Int("message_entities", len(msg.Entities)).
				Bool("message_has_reply", msg.ReplyingToPreview != nil).
				Bool("message_has_forward", msg.ForwardedMessage != nil)
		}
	}

	// Also try to dump the holder as JSON for full visibility
	holderJSON, jsonErr := json.Marshal(holder)
	if jsonErr == nil && len(holderJSON) < 2000 {
		evt = evt.RawJSON("holder_json", holderJSON)
	}

	evt = evt.Str("plaintext_full_hex", hex.EncodeToString(plaintext))
	evt.Msg("Thrift decode into MessageEntryHolder succeeded")
}

// truncateHexBytes returns hex representation of first n bytes.
func truncateHexBytes(b []byte, n int) string {
	if len(b) <= n {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:n]) + "..."
}

// DecryptMessageContents decrypts hex-encoded ciphertext to MessageContents.
func DecryptMessageContents(contentsHex string, conversationKey []byte) (*payload.MessageContents, error) {
	ciphertext, err := hex.DecodeString(contentsHex)
	if err != nil {
		return nil, fmt.Errorf("decode hex: %w", err)
	}
	return DecryptMessageContentsBytes(ciphertext, conversationKey)
}

// ParseMessageContentsBytes parses raw plaintext thrift bytes to MessageContents (unencrypted).
// Tries to decode as MessageContents first, falls back to MessageEntryHolder wrapper if that fails.
func ParseMessageContentsBytes(data []byte) (*payload.MessageContents, error) {
	// First try decoding directly as MessageContents
	var contents payload.MessageContents
	decoder := thrifter.NewDecoder(bytes.NewReader(data), nil)
	if err := decoder.Decode(&contents); err == nil {
		// Check if we got valid content
		if contents.MessageText != nil || len(contents.Attachments) > 0 {
			return &contents, nil
		}
	}

	// Fall back to MessageEntryHolder wrapper format
	entry, err := ParseMessageEntryContentsBytes(data)
	if err != nil {
		return nil, fmt.Errorf("thrift decode (tried both formats): %w", err)
	}

	if entry.Message != nil {
		return entry.Message, nil
	}

	// Check for other content types in the holder (reactions, edits, etc.)
	if entry != nil {
		// Return empty contents for non-message events that we don't handle
		return nil, fmt.Errorf("no message in holder (may be reaction/edit/other event type)")
	}

	return nil, fmt.Errorf("no message contents in holder")
}

// ParseMessageEntryContentsBytes parses raw plaintext thrift bytes into MessageEntryContents.
func ParseMessageEntryContentsBytes(data []byte) (*payload.MessageEntryContents, error) {
	return decodeMessageEntryHolder(data, nil)
}

func decodeMessageEntryHolder(data []byte, log *zerolog.Logger) (_ *payload.MessageEntryContents, err error) {
	defer func() {
		if r := recover(); r != nil {
			if log != nil {
				log.Warn().
					Interface("panic", r).
					Str("plaintext_full_hex", hex.EncodeToString(data)).
					Msg("Thrift decode panic in MessageEntryHolder")
			}
			err = fmt.Errorf("thrift decode panic: %v", r)
		}
	}()

	var holder payload.MessageEntryHolder
	decoder := thrifter.NewDecoder(bytes.NewReader(data), nil)
	if err := decoder.Decode(&holder); err != nil {
		if log != nil {
			log.Debug().
				Err(err).
				Str("plaintext_full_hex", hex.EncodeToString(data)).
				Msg("Thrift decode into MessageEntryHolder failed")
		}
		return nil, err
	}

	if log != nil {
		logHolderContents(log, &holder, data)
	}

	if holder.Contents == nil {
		return nil, fmt.Errorf("no message entry contents in holder")
	}

	return holder.Contents, nil
}

// ParseMessageContents parses hex-encoded plaintext thrift MessageContents (unencrypted).
func ParseMessageContents(contentsHex string) (*payload.MessageContents, error) {
	data, err := hex.DecodeString(contentsHex)
	if err != nil {
		return nil, fmt.Errorf("decode hex: %w", err)
	}
	return ParseMessageContentsBytes(data)
}

// EncryptMessageContents encrypts MessageContents to hex-encoded ciphertext.
func EncryptMessageContents(contents *payload.MessageContents, conversationKey []byte) (string, error) {
	ciphertext, err := EncryptMessageContentsRaw(contents, conversationKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}

// EncryptMessageContentsRaw encrypts MessageContents to raw ciphertext bytes.
func EncryptMessageContentsRaw(contents *payload.MessageContents, conversationKey []byte) ([]byte, error) {
	// Encode MessageContents to Thrift binary
	contentsBytes, err := EncodeMessageContents(contents)
	if err != nil {
		return nil, fmt.Errorf("thrift encode: %w", err)
	}

	ciphertext, err := SecretboxEncrypt(contentsBytes, conversationKey)
	if err != nil {
		return nil, fmt.Errorf("secretbox encrypt: %w", err)
	}

	return ciphertext, nil
}

// EncodeMessageContents encodes MessageContents to Thrift binary.
// The wire format is: MessageEntryHolder { MessageEntryContents.Message { MessageContents } }
func EncodeMessageContents(contents *payload.MessageContents) ([]byte, error) {
	// Ensure SentFrom is set (default to 1)
	if contents.SentFrom == nil {
		sentFrom := int32(1)
		contents.SentFrom = &sentFrom
	}

	return EncodeMessageEntryContents(&payload.MessageEntryContents{Message: contents})
}

// EncodeMessageEntryContents encodes MessageEntryContents to Thrift binary.
// The wire format is: MessageEntryHolder { MessageEntryContents }
func EncodeMessageEntryContents(contents *payload.MessageEntryContents) ([]byte, error) {
	if contents == nil {
		return nil, fmt.Errorf("message entry contents is nil")
	}

	holder := &payload.MessageEntryHolder{Contents: contents}

	var buf bytes.Buffer
	encoder := thrifter.NewEncoder(&buf)
	if err := encoder.Encode(holder); err != nil {
		return nil, fmt.Errorf("thrift encode: %w", err)
	}
	return buf.Bytes(), nil
}

// EncryptMessageEntryContentsRaw encrypts MessageEntryContents to raw ciphertext bytes.
func EncryptMessageEntryContentsRaw(contents *payload.MessageEntryContents, conversationKey []byte) ([]byte, error) {
	entryBytes, err := EncodeMessageEntryContents(contents)
	if err != nil {
		return nil, err
	}

	ciphertext, err := SecretboxEncrypt(entryBytes, conversationKey)
	if err != nil {
		return nil, fmt.Errorf("secretbox encrypt: %w", err)
	}

	return ciphertext, nil
}

// DecryptMessageEvent decrypts a MessageEvent's contents using the KeyManager.
func DecryptMessageEvent(ctx context.Context, km *KeyManager, event *payload.MessageEvent) (*payload.MessageContents, error) {
	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}
	if event.Detail == nil || event.Detail.MessageCreateEvent == nil {
		return nil, fmt.Errorf("not a message create event")
	}

	mce := event.Detail.MessageCreateEvent
	if len(mce.Contents) == 0 {
		return nil, fmt.Errorf("no contents in message")
	}

	conversationID := ""
	if event.ConversationId != nil {
		conversationID = *event.ConversationId
	}

	keyVersion := ""
	if mce.ConversationKeyVersion != nil {
		keyVersion = *mce.ConversationKeyVersion
	}

	convKey, err := km.GetConversationKey(ctx, conversationID, keyVersion)
	if err != nil {
		return nil, fmt.Errorf("get conversation key: %w", err)
	}

	return DecryptMessageContentsBytes(mce.Contents, convKey.Key)
}

// DecryptWithKey decrypts message contents using a provided key directly.
func DecryptWithKey(contentsHex string, key []byte) (*payload.MessageContents, error) {
	return DecryptMessageContents(contentsHex, key)
}

// MessageBuilder constructs MessageEvent structs for sending.
type MessageBuilder struct {
	km        *KeyManager
	ownUserID string

	// Message fields
	messageID      string
	conversationID string
	text           string
	attachments    []*payload.MessageAttachment
	replyTo        *payload.ReplyingToPreview
	entities       []*payload.RichTextEntity
	reactionAdd    *payload.MessageReactionAdd
	reactionRemove *payload.MessageReactionRemove
	messageEdit    *payload.MessageEdit

	// Pin/unpin fields
	pinConversation   *payload.PinConversation
	unpinConversation *payload.UnpinConversation

	// Crypto fields (per-call overrides)
	conversationKey     []byte
	keyVersion          string
	signingKey          *ecdsa.PrivateKey
	signatureKeyVersion string
}

// NewMessageBuilder creates a new builder.
// km can be nil for per-call key usage.
func NewMessageBuilder(km *KeyManager, ownUserID string) *MessageBuilder {
	return &MessageBuilder{
		km:        km,
		ownUserID: ownUserID,
	}
}

// SetMessageID sets the message ID (required).
func (b *MessageBuilder) SetMessageID(id string) *MessageBuilder {
	b.messageID = id
	return b
}

// SetConversationID sets the conversation ID (required).
func (b *MessageBuilder) SetConversationID(id string) *MessageBuilder {
	b.conversationID = id
	return b
}

// SetText sets the message text.
func (b *MessageBuilder) SetText(text string) *MessageBuilder {
	b.text = text
	return b
}

func (b *MessageBuilder) SetReactionAdd(targetMessageSequenceID, emoji string) *MessageBuilder {
	b.reactionAdd = &payload.MessageReactionAdd{
		MessageSequenceId: &targetMessageSequenceID,
		Emoji:             &emoji,
	}
	b.reactionRemove = nil
	return b
}

func (b *MessageBuilder) SetReactionRemove(targetMessageSequenceID, emoji string) *MessageBuilder {
	b.reactionRemove = &payload.MessageReactionRemove{
		MessageSequenceId: &targetMessageSequenceID,
		Emoji:             &emoji,
	}
	b.reactionAdd = nil
	return b
}

func (b *MessageBuilder) SetMessageEdit(targetMessageSequenceID, updatedText string, entities []*payload.RichTextEntity) *MessageBuilder {
	b.messageEdit = &payload.MessageEdit{
		MessageSequenceId: &targetMessageSequenceID,
		UpdatedText:       &updatedText,
		Entities:          entities,
	}
	b.reactionAdd = nil
	b.reactionRemove = nil
	return b
}

// SetPinConversation sets the conversation to pin.
func (b *MessageBuilder) SetPinConversation(targetConversationID string) *MessageBuilder {
	b.pinConversation = &payload.PinConversation{
		ConversationId: &targetConversationID,
	}
	b.unpinConversation = nil
	return b
}

// SetUnpinConversation sets the conversation to unpin.
func (b *MessageBuilder) SetUnpinConversation(targetConversationID string) *MessageBuilder {
	b.unpinConversation = &payload.UnpinConversation{
		ConversationId: &targetConversationID,
	}
	b.pinConversation = nil
	return b
}

// AddAttachment adds a media/url/post attachment.
func (b *MessageBuilder) AddAttachment(att *payload.MessageAttachment) *MessageBuilder {
	b.attachments = append(b.attachments, att)
	return b
}

// SetReplyTo sets the reply preview.
func (b *MessageBuilder) SetReplyTo(reply *payload.ReplyingToPreview) *MessageBuilder {
	b.replyTo = reply
	return b
}

// SetEntities sets rich text entities.
func (b *MessageBuilder) SetEntities(entities []*payload.RichTextEntity) *MessageBuilder {
	b.entities = entities
	return b
}

// WithConversationKey sets the conversation key directly (per-call approach).
func (b *MessageBuilder) WithConversationKey(key []byte, version string) *MessageBuilder {
	b.conversationKey = key
	b.keyVersion = version
	return b
}

// WithSigningKey sets the signing key directly (per-call approach).
func (b *MessageBuilder) WithSigningKey(key *ecdsa.PrivateKey, version string) *MessageBuilder {
	b.signingKey = key
	b.signatureKeyVersion = version
	return b
}

// Build constructs the MessageEvent, encrypting and signing as needed.
func (b *MessageBuilder) Build(ctx context.Context) (*payload.MessageEvent, error) {
	if b.messageID == "" {
		return nil, fmt.Errorf("message ID is required")
	}
	if b.conversationID == "" {
		return nil, fmt.Errorf("conversation ID is required")
	}

	// Get conversation key
	var convKey []byte
	var keyVersion string
	unencrypted := false

	if b.conversationKey != nil {
		convKey = b.conversationKey
		keyVersion = b.keyVersion
	} else if b.km != nil {
		key, err := b.km.GetLatestConversationKey(ctx, b.conversationID)
		if err != nil {
			if errors.Is(err, ErrKeyNotFound) {
				unencrypted = true
			} else {
				return nil, fmt.Errorf("get conversation key: %w", err)
			}
		} else if key == nil || len(key.Key) == 0 {
			unencrypted = true
		} else {
			convKey = key.Key
			keyVersion = key.KeyVersion
		}
	} else {
		// No key manager means we can't fetch a key; fall back to plaintext.
		unencrypted = true
	}

	// Build MessageEntryContents (message, reaction, pin/unpin, etc.).
	entryContents := &payload.MessageEntryContents{}
	if b.messageEdit != nil {
		entryContents.MessageEdit = b.messageEdit
	} else if b.reactionAdd != nil {
		entryContents.ReactionAdd = b.reactionAdd
	} else if b.reactionRemove != nil {
		entryContents.ReactionRemove = b.reactionRemove
	} else if b.pinConversation != nil {
		entryContents.PinConversation = b.pinConversation
	} else if b.unpinConversation != nil {
		entryContents.UnpinConversation = b.unpinConversation
	} else {
		contents := &payload.MessageContents{
			MessageText:       &b.text,
			Attachments:       b.attachments,
			Entities:          b.entities,
			ReplyingToPreview: b.replyTo,
		}

		// Ensure SentFrom is set (default to 1)
		if contents.SentFrom == nil {
			sentFrom := int32(1)
			contents.SentFrom = &sentFrom
		}

		entryContents.Message = contents
	}

	// Encrypt if we have a key, otherwise send plaintext contents.
	var contentsBytes []byte
	var err error
	if !unencrypted && len(convKey) > 0 {
		contentsBytes, err = EncryptMessageEntryContentsRaw(entryContents, convKey)
		if err != nil {
			return nil, fmt.Errorf("encrypt contents: %w", err)
		}
	} else {
		contentsBytes, err = EncodeMessageEntryContents(entryContents)
		if err != nil {
			return nil, fmt.Errorf("encode contents: %w", err)
		}
	}

	// Build MessageCreateEvent
	shouldNotify := true
	isPendingPublicKey := false
	priority := int32(1)
	mce := &payload.MessageCreateEvent{
		Contents:           contentsBytes,
		ShouldNotify:       &shouldNotify,
		IsPendingPublicKey: &isPendingPublicKey,
		Priority:           &priority,
	}
	if !unencrypted && keyVersion != "" {
		mce.ConversationKeyVersion = &keyVersion
	}

	// Get signing key
	var signingKey *ecdsa.PrivateKey
	var sigKeyVersion string

	if b.signingKey != nil {
		signingKey = b.signingKey
		sigKeyVersion = b.signatureKeyVersion
	} else if b.km != nil {
		ownKey, err := b.km.GetOwnSigningKey(ctx)
		if err == nil && ownKey != nil {
			signingKey = ownKey.SigningKey
			sigKeyVersion = ownKey.KeyVersion
		}
		// If no signing key, we'll just skip signing (not an error)
	}

	// Build MessageEvent
	event := &payload.MessageEvent{
		MessageId:      &b.messageID,
		SenderId:       &b.ownUserID,
		ConversationId: &b.conversationID,
		Detail: &payload.MessageEventDetail{
			MessageCreateEvent: mce,
		},
	}

	// Sign when we have a valid signing key+version. For plaintext (no conv key) the signature is still required.
	if signingKey != nil && sigKeyVersion != "" {
		signature, err := SignMessage(signingKey, b.messageID, b.ownUserID, b.conversationID, keyVersion, contentsBytes)
		if err != nil {
			return nil, fmt.Errorf("sign message: %w", err)
		}

		sigVersion := SignatureVersion3

		event.MessageEventSignature = &payload.MessageEventSignature{
			Signature:        &signature,
			PublicKeyVersion: &sigKeyVersion,
			SignatureVersion: &sigVersion,
		}
	}

	return event, nil
}

// BuildContentsOnly builds just the encrypted contents hex (for simpler use cases).
func (b *MessageBuilder) BuildContentsOnly(ctx context.Context) (contentsHex, keyVersion string, err error) {
	if b.conversationID == "" {
		return "", "", fmt.Errorf("conversation ID is required")
	}

	var convKey []byte

	if b.conversationKey != nil {
		convKey = b.conversationKey
		keyVersion = b.keyVersion
	} else if b.km != nil {
		key, err := b.km.GetLatestConversationKey(ctx, b.conversationID)
		if err != nil {
			return "", "", fmt.Errorf("get conversation key: %w", err)
		}
		convKey = key.Key
		keyVersion = key.KeyVersion
	} else {
		return "", "", fmt.Errorf("no conversation key provided")
	}

	contents := &payload.MessageContents{
		MessageText:       &b.text,
		Attachments:       b.attachments,
		Entities:          b.entities,
		ReplyingToPreview: b.replyTo,
	}

	contentsHex, err = EncryptMessageContents(contents, convKey)
	return contentsHex, keyVersion, err
}

// EncodeMessageCreateEvent encodes a MessageCreateEvent to base64 Thrift binary protocol.
func EncodeMessageCreateEvent(mce *payload.MessageCreateEvent) (string, error) {
	var buf bytes.Buffer
	encoder := thrifter.NewEncoder(&buf)
	if err := encoder.Encode(mce); err != nil {
		return "", fmt.Errorf("thrift encode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// EncodeMessageEventSignature encodes a MessageEventSignature to base64 Thrift binary protocol.
func EncodeMessageEventSignature(sig *payload.MessageEventSignature) (string, error) {
	var buf bytes.Buffer
	encoder := thrifter.NewEncoder(&buf)
	if err := encoder.Encode(sig); err != nil {
		return "", fmt.Errorf("thrift encode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// EncodeMessageEventDetail encodes a MessageEventDetail to base64 Thrift binary protocol.
// Used for delete message action signatures.
func EncodeMessageEventDetail(detail *payload.MessageEventDetail) (string, error) {
	var buf bytes.Buffer
	encoder := thrifter.NewEncoder(&buf)
	if err := encoder.Encode(detail); err != nil {
		return "", fmt.Errorf("thrift encode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// EncodeMuteConversation encodes a MuteConversation to base64 Thrift binary protocol.
func EncodeMuteConversation(detail *payload.MuteConversation) (string, error) {
	var buf bytes.Buffer
	encoder := thrifter.NewEncoder(&buf)
	if err := encoder.Encode(detail); err != nil {
		return "", fmt.Errorf("thrift encode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// EncodeUnmuteConversation encodes an UnmuteConversation to base64 Thrift binary protocol.
func EncodeUnmuteConversation(detail *payload.UnmuteConversation) (string, error) {
	var buf bytes.Buffer
	encoder := thrifter.NewEncoder(&buf)
	if err := encoder.Encode(detail); err != nil {
		return "", fmt.Errorf("thrift encode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// BuildForSend builds the encoded payloads needed for SendMessageMutation.
// Returns: encodedMessageCreateEvent, encodedSignature, error
func (b *MessageBuilder) BuildForSend(ctx context.Context) (string, string, error) {
	event, err := b.Build(ctx)
	if err != nil {
		return "", "", err
	}

	if event.Detail == nil || event.Detail.MessageCreateEvent == nil {
		return "", "", fmt.Errorf("no message create event in built message")
	}

	encodedMCE, err := EncodeMessageCreateEvent(event.Detail.MessageCreateEvent)
	if err != nil {
		return "", "", fmt.Errorf("encode message create event: %w", err)
	}

	encodedSig := ""
	if event.MessageEventSignature != nil {
		encodedSig, err = EncodeMessageEventSignature(event.MessageEventSignature)
		if err != nil {
			return "", "", fmt.Errorf("encode signature: %w", err)
		}
	}

	return encodedMCE, encodedSig, nil
}
