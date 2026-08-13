package twittermeow

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

// XChatEventHandler processes XChat events.
// Returns true if the event was successfully handled.
type XChatEventHandler func(ctx context.Context, evt types.TwitterEvent) bool

// SequenceIDCallback is called with the sequence ID of each processed event.
type SequenceIDCallback func(seqID string)

// XChatGapHandler backfills a conversation before the current live event is
// emitted when the processor detects a reconnect or a sequence discontinuity.
type XChatGapHandler func(ctx context.Context, conversationID, previousSequenceID, currentSequenceID string) error

// XChatEventProcessor processes XChat websocket events and converts them
// to TwitterEvent types for the bridge.
type XChatEventProcessor struct {
	client             *Client
	eventHandler       XChatEventHandler
	sequenceIDCallback SequenceIDCallback
	gapHandler         XChatGapHandler
	log                zerolog.Logger

	sequenceStateLock          sync.Mutex
	lastConversationSequence   map[string]string
	reconnectGeneration        uint64
	catchupGeneration          map[string]uint64
	unresolvedConversationGaps map[string]struct{}
	maxHandledSequenceID       string
	checkpointPublicationHolds int
}

func newXChatEventProcessor(client *Client) *XChatEventProcessor {
	return &XChatEventProcessor{
		client:                     client,
		log:                        client.Logger.With().Str("component", "xchat_processor").Logger(),
		lastConversationSequence:   make(map[string]string),
		catchupGeneration:          make(map[string]uint64),
		unresolvedConversationGaps: make(map[string]struct{}),
	}
}

// SetEventHandler sets the handler for processed XChat events.
func (p *XChatEventProcessor) SetEventHandler(handler XChatEventHandler) {
	p.eventHandler = handler
}

// SetSequenceIDCallback sets a callback that will be called with each processed event's sequence ID.
// This can be used to track the max sequence ID for incremental inbox fetching.
func (p *XChatEventProcessor) SetSequenceIDCallback(callback SequenceIDCallback) {
	p.sequenceStateLock.Lock()
	p.sequenceIDCallback = callback
	p.sequenceStateLock.Unlock()
}

// MaxHandledSequenceID returns the highest event sequence that completed
// handling during this connection, even while a gap blocks publication.
func (p *XChatEventProcessor) MaxHandledSequenceID() string {
	p.sequenceStateLock.Lock()
	defer p.sequenceStateLock.Unlock()
	return p.maxHandledSequenceID
}

// SetGapHandler sets the synchronous conversation catch-up hook.
func (p *XChatEventProcessor) SetGapHandler(handler XChatGapHandler) {
	p.gapHandler = handler
}

// MarkReconnected makes the next live event in each active conversation verify
// that conversation's forward history before emitting the event.
func (p *XChatEventProcessor) MarkReconnected() {
	p.sequenceStateLock.Lock()
	p.reconnectGeneration++
	p.sequenceStateLock.Unlock()
}

func (p *XChatEventProcessor) needsConversationCatchup(conversationID, previousSequenceID string) bool {
	if conversationID == "" {
		return false
	}
	p.sequenceStateLock.Lock()
	defer p.sequenceStateLock.Unlock()
	if _, unresolved := p.unresolvedConversationGaps[conversationID]; unresolved {
		return true
	}
	if p.catchupGeneration[conversationID] < p.reconnectGeneration {
		return true
	}
	lastSequenceID := p.lastConversationSequence[conversationID]
	return previousSequenceID != "" && lastSequenceID != "" &&
		compareXChatSequenceIDs(previousSequenceID, lastSequenceID) > 0
}

func (p *XChatEventProcessor) markConversationCaughtUp(conversationID string) {
	if conversationID == "" {
		return
	}
	p.sequenceStateLock.Lock()
	p.catchupGeneration[conversationID] = p.reconnectGeneration
	delete(p.unresolvedConversationGaps, conversationID)
	p.sequenceStateLock.Unlock()
}

// MarkConversationCaughtUp clears a previously blocked per-conversation gap
// after an inbox retry or explicit history repair succeeds.
func (p *XChatEventProcessor) MarkConversationCaughtUp(conversationID string) {
	p.markConversationCaughtUp(conversationID)
}

func (p *XChatEventProcessor) markConversationGapUnresolved(conversationID string) {
	if conversationID == "" || p.gapHandler == nil {
		return
	}
	p.sequenceStateLock.Lock()
	p.unresolvedConversationGaps[conversationID] = struct{}{}
	p.sequenceStateLock.Unlock()
}

// MarkConversationGapUnresolved keeps the global inbox checkpoint behind a
// conversation that could not be imported from an inbox page.
func (p *XChatEventProcessor) MarkConversationGapUnresolved(conversationID string) {
	p.markConversationGapUnresolved(conversationID)
}

// SequenceCheckpointBlocked reports whether at least one conversation still
// has history that must be repaired before the global inbox checkpoint moves.
func (p *XChatEventProcessor) SequenceCheckpointBlocked() bool {
	p.sequenceStateLock.Lock()
	defer p.sequenceStateLock.Unlock()
	return len(p.unresolvedConversationGaps) > 0
}

// ConversationGapUnresolved reports whether this conversation is preventing
// publication of the global inbox checkpoint.
func (p *XChatEventProcessor) ConversationGapUnresolved(conversationID string) bool {
	p.sequenceStateLock.Lock()
	defer p.sequenceStateLock.Unlock()
	_, unresolved := p.unresolvedConversationGaps[conversationID]
	return unresolved
}

// ResetSequenceState clears connection-local ordering state before a full
// reconnect. The persisted inbox checkpoint remains the source of truth.
func (p *XChatEventProcessor) ResetSequenceState() {
	p.sequenceStateLock.Lock()
	p.lastConversationSequence = make(map[string]string)
	p.catchupGeneration = make(map[string]uint64)
	p.unresolvedConversationGaps = make(map[string]struct{})
	p.reconnectGeneration = 0
	p.maxHandledSequenceID = ""
	p.checkpointPublicationHolds = 0
	p.sequenceStateLock.Unlock()
}

func (p *XChatEventProcessor) beginCheckpointBatch() {
	p.sequenceStateLock.Lock()
	p.checkpointPublicationHolds++
	p.sequenceStateLock.Unlock()
}

func (p *XChatEventProcessor) finishCheckpointBatch(success bool) {
	p.sequenceStateLock.Lock()
	if p.checkpointPublicationHolds > 0 {
		p.checkpointPublicationHolds--
	}
	checkpointSequenceID := ""
	callback := p.sequenceIDCallback
	if success && p.checkpointPublicationHolds == 0 && len(p.unresolvedConversationGaps) == 0 && callback != nil {
		checkpointSequenceID = p.maxHandledSequenceID
	}
	p.sequenceStateLock.Unlock()
	if checkpointSequenceID != "" {
		callback(checkpointSequenceID)
	}
}

// recordHandledSequence records successful per-conversation progress, but only
// publishes the global maximum after all known conversation gaps are repaired.
func (p *XChatEventProcessor) recordHandledSequence(conversationID, sequenceID string) {
	p.sequenceStateLock.Lock()
	if conversationID != "" && sequenceID != "" &&
		compareXChatSequenceIDs(sequenceID, p.lastConversationSequence[conversationID]) > 0 {
		p.lastConversationSequence[conversationID] = sequenceID
	}
	if sequenceID != "" && compareXChatSequenceIDs(sequenceID, p.maxHandledSequenceID) > 0 {
		p.maxHandledSequenceID = sequenceID
	}
	checkpointSequenceID := ""
	callback := p.sequenceIDCallback
	if p.checkpointPublicationHolds == 0 && len(p.unresolvedConversationGaps) == 0 && callback != nil {
		checkpointSequenceID = p.maxHandledSequenceID
	}
	p.sequenceStateLock.Unlock()
	if checkpointSequenceID != "" {
		callback(checkpointSequenceID)
	}
}

func compareXChatSequenceIDs(a, b string) int {
	a, aValid := normalizeXChatSequenceID(a)
	b, bValid := normalizeXChatSequenceID(b)
	if !aValid {
		if !bValid {
			return 0
		}
		return -1
	} else if !bValid {
		return 1
	}
	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return strings.Compare(a, b)
}

func normalizeXChatSequenceID(sequenceID string) (string, bool) {
	if sequenceID == "" {
		return "0", true
	}
	for _, char := range sequenceID {
		if char < '0' || char > '9' {
			return "", false
		}
	}
	sequenceID = strings.TrimLeft(sequenceID, "0")
	if sequenceID == "" {
		return "0", true
	}
	return sequenceID, true
}

// ProcessMessage handles a decoded payload.Message from the websocket.
// It may emit zero or more TwitterEvent objects via the event handler.
func (p *XChatEventProcessor) ProcessMessage(ctx context.Context, msg *payload.Message) (err error) {
	if msg == nil {
		return nil
	}
	p.beginCheckpointBatch()
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic processing XChat message (%T)", recovered)
		}
		p.finishCheckpointBatch(err == nil)
	}()
	// Handle single MessageEvent
	if msg.MessageEvent != nil {
		if err := p.processMessageEvent(ctx, msg.MessageEvent); err != nil {
			p.log.Err(err).Msg("Failed to process MessageEvent")
			return err
		}
	}

	// Handle MessageInstruction
	if msg.MessageInstruction != nil {
		if err := p.processInstruction(ctx, msg.MessageInstruction); err != nil {
			p.log.Err(err).Msg("Failed to process MessageInstruction")
			return err
		}
	}

	// Handle batched events
	if msg.BatchedMessageEvents != nil {
		for _, evt := range msg.BatchedMessageEvents.MessageEvents {
			if err := p.processMessageEvent(ctx, evt); err != nil {
				p.log.Err(err).Msg("Failed to process batched MessageEvent")
				return err
			}
		}
	}

	return nil
}

// processMessageEvent processes an individual MessageEvent and dispatches the appropriate TwitterEvent.
func (p *XChatEventProcessor) processMessageEvent(ctx context.Context, evt *payload.MessageEvent) error {
	return p.processMessageEventWithGapCatchup(ctx, evt, true)
}

func (p *XChatEventProcessor) processMessageEventWithGapCatchup(
	ctx context.Context,
	evt *payload.MessageEvent,
	allowGapCatchup bool,
) (err error) {
	if evt == nil {
		return nil
	}

	// Advance the incremental catch-up checkpoint only after the event has been
	// processed successfully. Advancing before decryption/queueing can make a
	// failed event permanently invisible to the next reconnect catch-up.
	sequenceID := ptr.Val(evt.SequenceId)
	conversationID := ptr.Val(evt.ConversationId)
	previousSequenceID := ptr.Val(evt.PreviousSequenceId)
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic processing XChat message event (%T)", recovered)
		}
		if err == nil {
			p.recordHandledSequence(conversationID, sequenceID)
		} else {
			p.markConversationGapUnresolved(conversationID)
		}
	}()
	if allowGapCatchup && p.needsConversationCatchup(conversationID, previousSequenceID) && p.gapHandler != nil {
		p.log.Info().
			Str("conversation_id", conversationID).
			Str("previous_sequence_id", previousSequenceID).
			Str("current_sequence_id", sequenceID).
			Msg("Detected XChat conversation gap, running forward catch-up")
		if catchupErr := p.gapHandler(ctx, conversationID, previousSequenceID, sequenceID); catchupErr != nil {
			// A historical gap must remain retryable, but it must not block the
			// current websocket event. Otherwise a transient Matrix/backfill
			// failure turns into a continuous live-message outage.
			p.markConversationGapUnresolved(conversationID)
			p.log.Warn().Err(catchupErr).
				Msg("XChat conversation catch-up failed; continuing with live event")
		} else {
			p.markConversationCaughtUp(conversationID)
		}
	}
	// Store conversation token if present
	if err := p.storeConversationToken(ctx, evt); err != nil {
		p.log.Warn().Err(err).Msg("Failed to store conversation token")
		return fmt.Errorf("store XChat conversation token: %w", err)
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

func isEncryptedMessageCreateEvent(mce *payload.MessageCreateEvent) bool {
	return mce != nil && ptr.Val(mce.ConversationKeyVersion) != ""
}

// processMessageCreateEvent handles a MessageCreateEvent, decrypting or parsing contents.
// Encrypted events declare a conversation key version. Signed plaintext XChat
// events omit it and contain directly encoded Thrift.
func (p *XChatEventProcessor) processMessageCreateEvent(ctx context.Context, evt *payload.MessageEvent, mce *payload.MessageCreateEvent) error {
	conversationID := ptr.Val(evt.ConversationId)
	contentsBytes := mce.Contents
	keyVersion := ptr.Val(mce.ConversationKeyVersion)

	if len(contentsBytes) == 0 {
		p.log.Debug().
			Str("sequence_id", ptr.Val(evt.SequenceId)).
			Str("conversation_id", conversationID).
			Msg("MessageCreateEvent has no contents, emitting ConversationCreate")
		return p.emitConversationCreate(ctx, evt)
	}

	var contents *payload.MessageEntryContents
	var err error

	if isEncryptedMessageCreateEvent(mce) {
		// Encrypted message - decrypt using conversation key
		convKey, err := p.client.keyManager.GetConversationKey(ctx, conversationID, keyVersion)
		if err == nil && (convKey == nil || len(convKey.Key) == 0) {
			err = crypto.ErrKeyNotFound
		}
		if errors.Is(err, crypto.ErrKeyNotFound) {
			// Try to fetch missing keys
			refreshErr := p.client.RefreshConversationKeys(ctx, conversationID)
			if refreshErr != nil {
				p.log.Warn().Err(refreshErr).
					Str("conversation_id", conversationID).
					Msg("Failed to refresh conversation keys")
			}
			// A refresh can store the requested key even if another key-change event
			// on the same response was malformed, so always retry the local lookup.
			convKey, err = p.client.keyManager.GetConversationKey(ctx, conversationID, keyVersion)
			if err == nil && (convKey == nil || len(convKey.Key) == 0) {
				err = crypto.ErrKeyNotFound
			}
			if err != nil && refreshErr != nil {
				err = errors.Join(err, refreshErr)
			}
		}
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("sequence_id", ptr.Val(evt.SequenceId)).
				Str("conversation_id", conversationID).
				Str("key_version", keyVersion).
				Int("contents_len", len(contentsBytes)).
				Msg("Failed to get conversation key, skipping message")
			return fmt.Errorf("get XChat conversation key %s: %w", keyVersion, err)
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
				Msg("Failed to decrypt message contents, skipping")
			return fmt.Errorf("decrypt XChat message contents: %w", err)
		}
	} else {
		// No conversation key version - parse as plaintext Thrift.
		contents, err = crypto.ParseMessageEntryContentsBytes(contentsBytes)
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("sequence_id", ptr.Val(evt.SequenceId)).
				Str("conversation_id", conversationID).
				Int("contents_len", len(contentsBytes)).
				Msg("Failed to parse message contents, skipping")
			return fmt.Errorf("parse plaintext XChat message contents: %w", err)
		}
	}

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

	if contents.MessageEdit != nil {
		return p.emitEvent(ctx, convertXChatMessageEdit(evt, contents.MessageEdit, keyVersion))
	}

	// Empty MessageCreateEvent - this often happens when a message request is accepted.
	// Emit a ConversationCreate event so the connector can create the room and backfill.
	p.log.Debug().
		Str("sequence_id", ptr.Val(evt.SequenceId)).
		Str("conversation_id", conversationID).
		Msg("MessageCreateEvent has no message content, emitting ConversationCreate")

	return p.emitConversationCreate(ctx, evt)
}

func (p *XChatEventProcessor) emitConversationCreate(ctx context.Context, evt *payload.MessageEvent) error {
	return p.emitEvent(ctx, &types.ConversationCreate{
		ID:             ptr.Val(evt.SequenceId),
		Time:           ptr.Val(evt.CreatedAtMsec),
		ConversationID: ptr.Val(evt.ConversationId),
		RequestID:      ptr.Val(evt.MessageId),
	})
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
	if newKeyVersion == "" {
		return errors.New("XChat conversation key change has no key version")
	}

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
	if signingKey == nil || signingKey.DecryptKeyB64 == "" {
		return errors.New("own XChat decryption key is missing")
	}

	ownUserID := p.client.GetCurrentUserID()
	if ownUserID == "" {
		p.log.Warn().
			Str("conversation_id", conversationID).
			Msg("Current user ID is empty while handling key change; cannot unwrap key")
		return errors.New("current user ID is empty while handling XChat key change")
	}

	var ourEncryptedKey string

	for _, pk := range ckce.ConversationParticipantKeys {
		if ptr.Val(pk.UserId) == ownUserID {
			ourEncryptedKey = ptr.Val(pk.EncryptedConversationKey)
			p.log.Info().Msg("Found encrypted key for current user")
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

	keyCreatedAt := methods.ParseMsecTimestamp(ptr.Val(evt.CreatedAtMsec))
	if keyCreatedAt.IsZero() {
		return fmt.Errorf("missing valid XChat key timestamp for conversation %s key %s", conversationID, newKeyVersion)
	}

	// Store the new key
	if err := p.client.keyManager.PutConversationKey(ctx, &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     newKeyVersion,
		Key:            convKeyBytes,
		CreatedAt:      keyCreatedAt,
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
			Bool("has_token", ptr.Val(inst.DisplayTemporaryPasscodeInstruction.Token) != "").
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
		return errors.New("xchat event handler is not configured")
	}

	if !p.eventHandler(ctx, evt) {
		return fmt.Errorf("xchat event handler rejected %T", evt)
	}
	return nil
}

// DecodeMessageEvent decodes a base64-encoded thrift MessageEvent string.
func DecodeMessageEvent(encoded string) (*payload.MessageEvent, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	var evt payload.MessageEvent
	if err := payload.Decode(data, &evt); err != nil {
		return nil, fmt.Errorf("thrift decode: %w", err)
	}
	return &evt, nil
}

// DecodeSendMessageEventResponse decodes the base64-encoded thrift response from SendMessageMutation.
func DecodeSendMessageEventResponse(encoded string) (*payload.SendMessageEventResponse, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	var resp payload.SendMessageEventResponse
	if err := payload.Decode(data, &resp); err != nil {
		return nil, fmt.Errorf("thrift decode: %w", err)
	}
	return &resp, nil
}

// DecodeSendMessageMutationMessageEvent decodes SendMessageMutation responses.
// The endpoint has returned both direct MessageEvent payloads and wrapper payloads.
func DecodeSendMessageMutationMessageEvent(encoded string) (*payload.MessageEvent, error) {
	evt, err := DecodeMessageEvent(encoded)
	if err == nil && evt != nil && evt.Detail != nil {
		return evt, nil
	}

	resp, wrapperErr := DecodeSendMessageEventResponse(encoded)
	if wrapperErr != nil {
		if err != nil {
			return nil, err
		}
		return nil, wrapperErr
	}
	if resp == nil || resp.MessageEvent == nil || *resp.MessageEvent == "" {
		return nil, fmt.Errorf("send message response did not contain a message event")
	}
	return DecodeMessageEvent(*resp.MessageEvent)
}

type decodedInboxMessageEvent struct {
	sequenceID string
	evt        *payload.MessageEvent
}

func (p *XChatEventProcessor) decodeAndSortInboxEvents(conversationID string, encodedEvents []string) ([]decodedInboxMessageEvent, error) {
	if len(encodedEvents) == 0 {
		return nil, nil
	}

	seen := make(map[string]struct{}, len(encodedEvents))
	out := make([]decodedInboxMessageEvent, 0, len(encodedEvents))
	var decodeErrs []error

	for index, encodedEvt := range encodedEvents {
		if encodedEvt == "" {
			continue
		}

		evt, err := DecodeMessageEvent(encodedEvt)
		if err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to decode inbox message event")
			decodeErrs = append(decodeErrs, fmt.Errorf("decode inbox message event %d: %w", index, err))
			continue
		}
		eventConversationID := ptr.Val(evt.ConversationId)
		if eventConversationID == "" && conversationID != "" {
			evt.ConversationId = ptr.Ptr(conversationID)
		} else if conversationID != "" && eventConversationID != conversationID {
			decodeErrs = append(decodeErrs, fmt.Errorf("inbox message event %d belongs to a different conversation", index))
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

		out = append(out, decodedInboxMessageEvent{sequenceID: seqID, evt: evt})
	}

	slices.SortStableFunc(out, func(a, b decodedInboxMessageEvent) int {
		switch {
		case a.sequenceID == "" && b.sequenceID != "":
			return 1
		case a.sequenceID != "" && b.sequenceID == "":
			return -1
		default:
			return compareXChatSequenceIDs(a.sequenceID, b.sequenceID)
		}
	})

	return out, errors.Join(decodeErrs...)
}

// ProcessKeyChangeEvents processes key change events from an XChatInboxItem.
// This should be called BEFORE syncing the channel, as keys are needed for decryption.
func (p *XChatEventProcessor) ProcessKeyChangeEvents(ctx context.Context, item *response.XChatInboxItem) (err error) {
	conversationID := item.ConversationDetail.ConversationID
	p.beginCheckpointBatch()
	defer func() {
		if err != nil {
			p.markConversationGapUnresolved(conversationID)
		}
		p.finishCheckpointBatch(err == nil)
	}()

	encodedEvents := make([]string, 0, len(item.LatestConversationKeyChangeEvents)+len(item.EncodedMessageEvents))
	encodedEvents = append(encodedEvents, item.LatestConversationKeyChangeEvents...)
	encodedEvents = append(encodedEvents, item.EncodedMessageEvents...)

	decodedEvents, decodeErr := p.decodeAndSortInboxEvents(conversationID, encodedEvents)
	for _, decoded := range decodedEvents {
		detail := decoded.evt.Detail
		if detail == nil || detail.ConversationKeyChangeEvent == nil {
			continue
		}
		if err := p.processMessageEventWithGapCatchup(ctx, decoded.evt, false); err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to process key change event from inbox")
			return errors.Join(decodeErr, err)
		}
	}

	return decodeErr
}

// ProcessMessageAndReadEvents processes message and read events from an XChatInboxItem.
// This should be called AFTER syncing the channel, as portals must exist for message handling.
func (p *XChatEventProcessor) ProcessMessageAndReadEvents(ctx context.Context, item *response.XChatInboxItem) (err error) {
	conversationID := item.ConversationDetail.ConversationID
	p.beginCheckpointBatch()
	defer func() {
		if err != nil {
			p.markConversationGapUnresolved(conversationID)
		}
		p.finishCheckpointBatch(err == nil)
	}()

	processedSeqIDs := make(map[string]struct{})

	encodedEvents := make([]string, 0, len(item.LatestMessageEvents)+len(item.EncodedMessageEvents)+1)
	encodedEvents = append(encodedEvents, item.LatestMessageEvents...)
	encodedEvents = append(encodedEvents, item.EncodedMessageEvents...)
	if item.LatestNotifiableMessageCreateEvent != "" {
		encodedEvents = append(encodedEvents, item.LatestNotifiableMessageCreateEvent)
	}

	decodedEvents, decodeErr := p.decodeAndSortInboxEvents(conversationID, encodedEvents)
	for _, decoded := range decodedEvents {
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

		if err := p.processMessageEventWithGapCatchup(ctx, decoded.evt, false); err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to process message event from inbox")
			return errors.Join(decodeErr, err)
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
		if ptr.Val(evt.SenderId) == "" && readEvt.ParticipantID.RestID != "" {
			evt.SenderId = ptr.Ptr(readEvt.ParticipantID.RestID)
		}
		if ptr.Val(evt.ConversationId) == "" {
			evt.ConversationId = ptr.Ptr(conversationID)
		} else if ptr.Val(evt.ConversationId) != conversationID {
			p.log.Warn().
				Str("conversation_id", conversationID).
				Str("participant_id", readEvt.ParticipantID.RestID).
				Msg("Ignoring inbox read event for a different conversation")
			continue
		}
		if seqID := ptr.Val(evt.SequenceId); seqID != "" {
			if _, ok := processedSeqIDs[seqID]; ok {
				continue
			}
			processedSeqIDs[seqID] = struct{}{}
		}
		if err := p.processMessageEventWithGapCatchup(ctx, evt, false); err != nil {
			p.log.Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Msg("Failed to process read event from initial inbox")
			return errors.Join(decodeErr, err)
		}
	}

	return decodeErr
}
