package twittermeow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	thrifter "github.com/thrift-iterator/go"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// XChatEventHandler processes XChat events.
// Returns true if the event was successfully handled.
type XChatEventHandler func(ctx context.Context, evt types.TwitterEvent) bool

// SequenceIDCallback is called with the sequence ID of each processed event.
type SequenceIDCallback func(seqID string)

// XChatEventProcessor processes XChat websocket events and converts them
// to TwitterEvent types for the bridge.
type XChatEventProcessor struct {
	client             *Client
	eventHandler       XChatEventHandler
	sequenceIDCallback SequenceIDCallback
	log                zerolog.Logger
}

func newXChatEventProcessor(client *Client) *XChatEventProcessor {
	return &XChatEventProcessor{
		client: client,
		log:    client.Logger.With().Str("component", "xchat_processor").Logger(),
	}
}

// SetEventHandler sets the handler for processed XChat events.
func (p *XChatEventProcessor) SetEventHandler(handler XChatEventHandler) {
	p.eventHandler = handler
}

// SetSequenceIDCallback sets a callback that will be called with each processed event's sequence ID.
// This can be used to track the max sequence ID for incremental inbox fetching.
func (p *XChatEventProcessor) SetSequenceIDCallback(callback SequenceIDCallback) {
	p.sequenceIDCallback = callback
}

// trackSequenceID reports the sequence ID to the callback if set.
func (p *XChatEventProcessor) trackSequenceID(seqID string) {
	if seqID != "" && p.sequenceIDCallback != nil {
		p.sequenceIDCallback(seqID)
	}
}

// ProcessMessage handles a decoded payload.Message from the websocket.
// It may emit zero or more TwitterEvent objects via the event handler.
func (p *XChatEventProcessor) ProcessMessage(ctx context.Context, msg *payload.Message) error {
	if msg == nil {
		return nil
	}

	// Handle single MessageEvent
	if msg.MessageEvent != nil {
		if err := p.processMessageEvent(ctx, msg.MessageEvent); err != nil {
			p.log.Err(err).Msg("Failed to process MessageEvent")
		}
	}

	// Handle MessageInstruction
	if msg.MessageInstruction != nil {
		if err := p.processInstruction(ctx, msg.MessageInstruction); err != nil {
			p.log.Err(err).Msg("Failed to process MessageInstruction")
		}
	}

	// Handle batched events
	if msg.BatchedMessageEvents != nil {
		for _, evt := range msg.BatchedMessageEvents.MessageEvents {
			if err := p.processMessageEvent(ctx, evt); err != nil {
				p.log.Err(err).Msg("Failed to process batched MessageEvent")
			}
		}
	}

	return nil
}

// processMessageEvent processes an individual MessageEvent and dispatches the appropriate TwitterEvent.
func (p *XChatEventProcessor) processMessageEvent(ctx context.Context, evt *payload.MessageEvent) error {
	if evt == nil {
		return nil
	}

	// Track the sequence ID for incremental fetching
	p.trackSequenceID(ptr.Val(evt.SequenceId))

	// Store conversation token if present
	if err := p.storeConversationToken(ctx, evt); err != nil {
		p.log.Warn().Err(err).Msg("Failed to store conversation token")
	}

	detail := evt.Detail
	if detail == nil {
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("conversation_id", ptr.Val(evt.ConversationId)).
			Msg("MessageEvent has no detail")
		return nil
	}

	// Process based on event type
	switch {
	case detail.MessageCreateEvent != nil:
		return p.processMessageCreateEvent(ctx, evt, detail.MessageCreateEvent)

	case detail.MessageDeleteEvent != nil:
		return p.emitEvent(ctx, convertXChatMessageDelete(evt, detail.MessageDeleteEvent))

	case detail.MessageTypingEvent != nil:
		convID, senderID := convertXChatTypingEvent(evt, detail.MessageTypingEvent)
		return p.emitEvent(ctx, &types.XChatTyping{
			ConversationID: convID,
			SenderID:       senderID,
			Timestamp:      time.Now(),
		})

	case detail.GroupChangeEvent != nil:
		return p.processGroupChangeEvent(ctx, evt, detail.GroupChangeEvent)

	case detail.ConversationKeyChangeEvent != nil:
		return p.processConversationKeyChange(ctx, evt, detail.ConversationKeyChangeEvent)

	case detail.ConversationDeleteEvent != nil:
		return p.emitEvent(ctx, convertXChatConversationDelete(evt, detail.ConversationDeleteEvent))

	case detail.MarkConversationReadEvent != nil:
		return p.emitEvent(ctx, convertXChatMarkReadEvent(evt, detail.MarkConversationReadEvent))

	case detail.MessageFailureEvent != nil:
		return p.processMessageFailure(ctx, evt, detail.MessageFailureEvent)

	case detail.ConversationMetadataChangeEvent != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Msg("Ignoring ConversationMetadataChangeEvent")
		return nil

	case detail.RequestForEncryptedResendEvent != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("min_seq", ptr.Val(detail.RequestForEncryptedResendEvent.MinSequenceId)).
			Str("max_seq", ptr.Val(detail.RequestForEncryptedResendEvent.MaxSequenceId)).
			Msg("Received RequestForEncryptedResendEvent")
		return nil

	case detail.GrokSearchResponseEvent != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Msg("Ignoring GrokSearchResponseEvent")
		return nil

	case detail.MemberAccountDeleteEvent != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("member_id", ptr.Val(detail.MemberAccountDeleteEvent.MemberId)).
			Msg("Received MemberAccountDeleteEvent")
		return nil

	case detail.MarkConversationUnreadEvent != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Msg("Ignoring MarkConversationUnreadEvent")
		return nil

	default:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("conversation_id", ptr.Val(evt.ConversationId)).
			Msg("Unknown MessageEventDetail type")
		return nil
	}
}

// processMessageCreateEvent handles a MessageCreateEvent, decrypting or parsing contents.
// If the event has a signature, the contents are encrypted and need decryption.
// If no signature, the contents are plaintext thrift and can be parsed directly.
func (p *XChatEventProcessor) processMessageCreateEvent(ctx context.Context, evt *payload.MessageEvent, mce *payload.MessageCreateEvent) error {
	conversationID := ptr.Val(evt.ConversationId)
	contentsBytes := mce.Contents
	keyVersion := ptr.Val(mce.ConversationKeyVersion)

	if len(contentsBytes) == 0 {
		p.log.Warn().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("conversation_id", conversationID).
			Msg("MessageCreateEvent has no contents")
		return nil
	}

	var contents *payload.MessageEntryContents
	var err error

	// Check if message has a signature - if so, it's encrypted
	if evt.MessageEventSignature != nil && evt.MessageEventSignature.Signature != nil {
		// Encrypted message - decrypt using conversation key
		convKey, err := p.client.keyManager.GetConversationKey(ctx, conversationID, keyVersion)
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("sequence_id", ptr.Val(evt.SequenceId)).
				Str("conversation_id", conversationID).
				Str("key_version", keyVersion).
				Int("contents_len", len(contentsBytes)).
				Msg("Failed to get conversation key, skipping message")
			return nil
		}

		debugLog := p.log.With().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("conversation_id", conversationID).
			Str("key_version", keyVersion).
			Logger()
		contents, err = crypto.DecryptMessageEntryContentsBytesDebug(contentsBytes, convKey.Key, &debugLog)
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("sequence_id", ptr.Val(evt.SequenceId)).
				Str("conversation_id", conversationID).
				Int("contents_len", len(contentsBytes)).
				Str("contents_hex_prefix", truncateBytes(contentsBytes, 32)).
				Msg("Failed to decrypt message contents, skipping")
			return nil
		}
		} else {
			// No signature - parse as plaintext thrift
			contents, err = crypto.ParseMessageEntryContentsBytes(contentsBytes)
			if err != nil {
				p.log.Warn().
				Err(err).
				Str("sequence_id", ptr.Val(evt.SequenceId)).
				Str("conversation_id", conversationID).
				Int("contents_len", len(contentsBytes)).
				Str("contents_hex_prefix", truncateBytes(contentsBytes, 32)).
					Msg("Failed to parse message contents, skipping")
				return nil
			}
		}

		p.logDecodedMessageContents(evt, conversationID, contents)

		// MessageContents directly contains message data (MessageText, Attachments, etc.)
		// Check if it has actual message content
		if contents.Message != nil && (contents.Message.MessageText != nil || len(contents.Message.Attachments) > 0) {
			msg := convertXChatMessageToTwitterMessage(evt, contents.Message, keyVersion)
			return p.emitEvent(ctx, msg)
	}

	if contents.ReactionAdd != nil {
		return p.emitEvent(ctx, convertXChatReactionAdd(evt, contents.ReactionAdd))
	}

	if contents.ReactionRemove != nil {
		return p.emitEvent(ctx, convertXChatReactionRemove(evt, contents.ReactionRemove))
	}

	p.log.Debug().
		Str("sequence_id", ptr.Val(evt.SequenceId)).
		Str("conversation_id", conversationID).
		Msg("MessageCreateEvent parsed but no message text, attachments, or reactions")

	return nil
}

// processGroupChangeEvent handles group-related changes.
func (p *XChatEventProcessor) processGroupChangeEvent(ctx context.Context, evt *payload.MessageEvent, gce *payload.GroupChangeEvent) error {
	gc := gce.GroupChange
	if gc == nil {
		return nil
	}

	switch {
	case gc.GroupMemberAdd != nil:
		return p.emitEvent(ctx, convertXChatGroupMemberAdd(evt, gc.GroupMemberAdd))
	case gc.GroupMemberRemove != nil:
		return p.emitEvent(ctx, convertXChatGroupMemberRemove(evt, gc.GroupMemberRemove))
	case gc.GroupTitleChange != nil:
		return p.emitEvent(ctx, convertXChatGroupTitleChange(evt, gc.GroupTitleChange))
	case gc.GroupAvatarChange != nil:
		return p.emitEvent(ctx, convertXChatGroupAvatarChange(evt, gc.GroupAvatarChange))
	case gc.GroupCreate != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Int("member_count", len(gc.GroupCreate.MemberIds)).
			Msg("Received GroupCreate event")
		return nil
	case gc.GroupAdminAdd != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Strs("admin_ids", gc.GroupAdminAdd.AdminIds).
			Msg("Received GroupAdminAdd event")
		return nil
	case gc.GroupAdminRemove != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Strs("admin_ids", gc.GroupAdminRemove.AdminIds).
			Msg("Received GroupAdminRemove event")
		return nil
	case gc.GroupInviteEnable != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("invite_url", ptr.Val(gc.GroupInviteEnable.InviteUrl)).
			Msg("Received GroupInviteEnable event")
		return nil
	case gc.GroupInviteDisable != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Msg("Received GroupInviteDisable event")
		return nil
	case gc.GroupJoinRequest != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("requesting_user", ptr.Val(gc.GroupJoinRequest.RequestingUserId)).
			Msg("Received GroupJoinRequest event")
		return nil
	case gc.GroupJoinReject != nil:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Strs("rejected_users", gc.GroupJoinReject.RejectedUserIds).
			Msg("Received GroupJoinReject event")
		return nil
	default:
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Msg("Unknown GroupChange type")
		return nil
	}
}

// processConversationKeyChange handles key rotation events.
func (p *XChatEventProcessor) processConversationKeyChange(ctx context.Context, evt *payload.MessageEvent, ckce *payload.ConversationKeyChangeEvent) error {
	conversationID := ptr.Val(evt.ConversationId)
	newKeyVersion := ptr.Val(ckce.ConversationKeyVersion)

	p.log.Info().
		Str("sequence_id", ptr.Val(evt.SequenceId)).
		Str("conversation_id", conversationID).
		Str("new_key_version", newKeyVersion).
		Int("participant_keys", len(ckce.ConversationParticipantKeys)).
		Msg("Processing ConversationKeyChangeEvent")

	signingKey, err := p.client.keyManager.GetOwnSigningKey(ctx)
	if err != nil {
		p.log.Err(err).
			Str("conversation_id", conversationID).
			Msg("Failed to get own signing key for key unwrap")
		return err
	}

	ownUserID := p.client.GetCurrentUserID()
	if ownUserID == "" {
		p.log.Warn().
			Str("conversation_id", conversationID).
			Msg("Current user ID is empty while handling key change; cannot unwrap key")
		return nil
	}

	var ourEncryptedKey string

	for _, pk := range ckce.ConversationParticipantKeys {
		if ptr.Val(pk.UserId) == ownUserID {
			ourEncryptedKey = ptr.Val(pk.EncryptedConversationKey)
			p.log.Info().Msgf("Our Encrypted Key Is Found")
			break
		}
	}

	if ourEncryptedKey == "" {
		p.log.Warn().
			Str("conversation_id", conversationID).
			Str("own_user_id", ownUserID).
			Msg("No encrypted key for own user in key change event")
		return nil
	}

	convKeyBytes, err := crypto.UnwrapConversationKey(ourEncryptedKey, signingKey.DecryptKeyB64)
	if err != nil {
		p.log.Err(err).
			Str("conversation_id", conversationID).
			Str("key_version", newKeyVersion).
			Msg("Failed to unwrap conversation key")
		return err
	}

	// Store the new key
	if err := p.client.keyManager.PutConversationKey(ctx, &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     newKeyVersion,
		Key:            convKeyBytes,
		CreatedAt:      time.Now(),
	}); err != nil {
		p.log.Err(err).
			Str("conversation_id", conversationID).
			Str("key_version", newKeyVersion).
			Msg("Failed to store new conversation key")
		return err
	}

	p.log.Info().
		Str("conversation_id", conversationID).
		Str("key_version", newKeyVersion).
		Msg("Successfully stored new conversation key")

	// Emit informational event
	return p.emitEvent(ctx, &types.XChatKeyChange{
		ID:             ptr.Val(evt.SequenceId),
		ConversationID: conversationID,
		SenderID:       ptr.Val(evt.SenderId),
		NewKeyVersion:  newKeyVersion,
		Timestamp:      time.Now(),
	})
}

// processMessageFailure handles message send failure events.
func (p *XChatEventProcessor) processMessageFailure(ctx context.Context, evt *payload.MessageEvent, failure *payload.MessageFailureEvent) error {
	failureType := payload.FailureType(ptr.Val(failure.FailureType))

	p.log.Warn().
		Str("sequence_id", ptr.Val(evt.SequenceId)).
		Str("conversation_id", ptr.Val(evt.ConversationId)).
		Int32("failure_type", ptr.Val(failure.FailureType)).
		Msg("Received MessageFailureEvent")

	return p.emitEvent(ctx, &types.XChatMessageFailure{
		ConversationID: ptr.Val(evt.ConversationId),
		MessageID:      ptr.Val(evt.MessageId),
		FailureType:    failureType,
		Timestamp:      time.Now(),
	})
}

// processInstruction handles MessageInstruction types.
func (p *XChatEventProcessor) processInstruction(ctx context.Context, inst *payload.MessageInstruction) error {
	if inst == nil {
		return nil
	}

	switch {
	case inst.PullMessagesInstruction != nil:
		p.log.Debug().
			Str("sequence_start", ptr.Val(inst.PullMessagesInstruction.SequenceStart)).
			Str("sender_id", ptr.Val(inst.PullMessagesInstruction.SenderId)).
			Bool("is_batched", ptr.Val(inst.PullMessagesInstruction.IsBatchedPull)).
			Msg("Received PullMessagesInstruction")
		return nil

	case inst.PullMessagesFinishedInstruction != nil:
		p.log.Debug().
			Bool("finished", ptr.Val(inst.PullMessagesFinishedInstruction.FinishedPull)).
			Str("sequence_continue", ptr.Val(inst.PullMessagesFinishedInstruction.SequenceContinue)).
			Msg("Received PullMessagesFinishedInstruction")
		return nil

	case inst.KeepAliveInstruction != nil:
		p.log.Trace().Msg("Received KeepAliveInstruction")
		return nil

	case inst.PinReminderInstruction != nil:
		p.log.Debug().
			Bool("should_register", ptr.Val(inst.PinReminderInstruction.ShouldRegister)).
			Bool("should_generate", ptr.Val(inst.PinReminderInstruction.ShouldGenerate)).
			Msg("Received PinReminderInstruction")
		return nil

	case inst.SwitchToHybridPullInstruction != nil:
		p.log.Debug().
			Str("user_agent", ptr.Val(inst.SwitchToHybridPullInstruction.RequestingUserAgent)).
			Msg("Received SwitchToHybridPullInstruction")
		return nil

	case inst.DisplayTemporaryPasscodeInstruction != nil:
		p.log.Debug().
			Str("token", ptr.Val(inst.DisplayTemporaryPasscodeInstruction.Token)).
			Str("public_key_version", ptr.Val(inst.DisplayTemporaryPasscodeInstruction.LatestPublicKeyVersion)).
			Msg("Received DisplayTemporaryPasscodeInstruction")
		return nil

	default:
		p.log.Debug().Msg("Unknown MessageInstruction type")
		return nil
	}
}

// storeConversationToken stores the conversation token from a MessageEvent if present.
func (p *XChatEventProcessor) storeConversationToken(ctx context.Context, evt *payload.MessageEvent) error {
	token := ptr.Val(evt.ConversationToken)
	conversationID := ptr.Val(evt.ConversationId)

	if token == "" || conversationID == "" {
		return nil
	}

	return p.client.keyManager.PutConversationToken(ctx, conversationID, token)
}

// emitEvent dispatches a TwitterEvent to the handler.
func (p *XChatEventProcessor) emitEvent(ctx context.Context, evt types.TwitterEvent) error {
	if p.eventHandler == nil {
		p.log.Warn().Type("event_type", evt).Msg("No event handler set, dropping event")
		return nil
	}

	p.eventHandler(ctx, evt)
	return nil
}

// truncateBytes returns a hex representation of the first n bytes for logging.
func truncateBytes(b []byte, n int) string {
	if len(b) <= n {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:n]) + "..."
}

func (p *XChatEventProcessor) logDecodedMessageContents(evt *payload.MessageEvent, conversationID string, contents *payload.MessageEntryContents) {
	if contents == nil {
		return
	}
	seqID := ptr.Val(evt.SequenceId)
	raw, err := json.Marshal(contents)
	if err != nil {
		p.log.Warn().
			Err(err).
			Str("sequence_id", seqID).
			Str("conversation_id", conversationID).
			Msg("Failed to marshal decrypted message contents")
		return
	}

	logEvt := p.log.Info().
		Str("sequence_id", seqID).
		Str("conversation_id", conversationID).
		Int("decoded_json_len", len(raw))
	if len(raw) <= 4000 {
		logEvt = logEvt.RawJSON("decrypted_contents", raw)
	} else {
		logEvt = logEvt.Str("decrypted_contents_prefix", string(raw[:4000]))
	}
	logEvt.Msg("Decrypted MessageCreateEvent contents")
}

// DecodeMessageEvent decodes a base64-encoded thrift MessageEvent string.
func DecodeMessageEvent(encoded string) (*payload.MessageEvent, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	decoder := thrifter.NewDecoder(bytes.NewReader(data), nil)
	var evt payload.MessageEvent
	if err := decoder.Decode(&evt); err != nil {
		return nil, fmt.Errorf("thrift decode: %w", err)
	}
	return &evt, nil
}

type decodedInboxMessageEvent struct {
	seq int64
	evt *payload.MessageEvent
}

func (p *XChatEventProcessor) decodeAndSortInboxEvents(conversationID string, encodedEvents []string) []decodedInboxMessageEvent {
	if len(encodedEvents) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(encodedEvents))
	out := make([]decodedInboxMessageEvent, 0, len(encodedEvents))

	for _, encodedEvt := range encodedEvents {
		if encodedEvt == "" {
			continue
		}

		evt, err := DecodeMessageEvent(encodedEvt)
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to decode inbox message event")
			continue
		}

		seqID := ptr.Val(evt.SequenceId)
		seenKey := seqID
		if seenKey == "" {
			seenKey = encodedEvt
		}
		if _, ok := seen[seenKey]; ok {
			continue
		}
		seen[seenKey] = struct{}{}

		var seq int64
		if seqID != "" {
			seq, _ = strconv.ParseInt(seqID, 10, 64)
		}
		out = append(out, decodedInboxMessageEvent{seq: seq, evt: evt})
	}

	slices.SortStableFunc(out, func(a, b decodedInboxMessageEvent) int {
		switch {
		case a.seq == 0 && b.seq != 0:
			return 1
		case a.seq != 0 && b.seq == 0:
			return -1
		case a.seq < b.seq:
			return -1
		case a.seq > b.seq:
			return 1
		default:
			return 0
		}
	})

	return out
}

// ProcessKeyChangeEvents processes key change events from an XChatInboxItem.
// This should be called BEFORE syncing the channel, as keys are needed for decryption.
func (p *XChatEventProcessor) ProcessKeyChangeEvents(ctx context.Context, item *response.XChatInboxItem) error {
	conversationID := item.ConversationDetail.ConversationID

	encodedEvents := make([]string, 0, len(item.LatestConversationKeyChangeEvents)+len(item.EncodedMessageEvents))
	encodedEvents = append(encodedEvents, item.LatestConversationKeyChangeEvents...)
	encodedEvents = append(encodedEvents, item.EncodedMessageEvents...)

	for _, decoded := range p.decodeAndSortInboxEvents(conversationID, encodedEvents) {
		detail := decoded.evt.Detail
		if detail == nil || detail.ConversationKeyChangeEvent == nil {
			continue
		}
		if err := p.processMessageEvent(ctx, decoded.evt); err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to process key change event from inbox")
		}
	}

	return nil
}

// ProcessMessageAndReadEvents processes message and read events from an XChatInboxItem.
// This should be called AFTER syncing the channel, as portals must exist for message handling.
func (p *XChatEventProcessor) ProcessMessageAndReadEvents(ctx context.Context, item *response.XChatInboxItem) error {
	conversationID := item.ConversationDetail.ConversationID

	processedSeqIDs := make(map[string]struct{})

	encodedEvents := make([]string, 0, len(item.LatestMessageEvents)+len(item.EncodedMessageEvents))
	encodedEvents = append(encodedEvents, item.LatestMessageEvents...)
	encodedEvents = append(encodedEvents, item.EncodedMessageEvents...)

	for _, decoded := range p.decodeAndSortInboxEvents(conversationID, encodedEvents) {
		seqID := ptr.Val(decoded.evt.SequenceId)
		if seqID != "" {
			if _, ok := processedSeqIDs[seqID]; ok {
				continue
			}
			processedSeqIDs[seqID] = struct{}{}
		}

		detail := decoded.evt.Detail
		if detail != nil && detail.ConversationKeyChangeEvent != nil {
			continue
		}

		if err := p.processMessageEvent(ctx, decoded.evt); err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to process message event from inbox")
		}
	}

	// Process read events per participant
	for _, readEvt := range item.LatestReadEventsPerParticipant {
		if readEvt.LatestMarkConversationReadEvent == "" {
			continue
		}
		evt, err := DecodeMessageEvent(readEvt.LatestMarkConversationReadEvent)
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Str("participant_id", readEvt.ParticipantID.RestID).
				Msg("Failed to decode read event from initial inbox")
			continue
		}
		if seqID := ptr.Val(evt.SequenceId); seqID != "" {
			if _, ok := processedSeqIDs[seqID]; ok {
				continue
			}
			processedSeqIDs[seqID] = struct{}{}
		}
		if err := p.processMessageEvent(ctx, evt); err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to process read event from initial inbox")
		}
	}

	return nil
}
