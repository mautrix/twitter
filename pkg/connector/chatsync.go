// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2025 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// Conversation type constants for Twitter DM conversations.
const (
	ConversationTypeOneToOne = "ONE_TO_ONE"
	ConversationTypeGroupDM  = "GROUP_DM"
)

// syncXChatChannel syncs a single conversation from XChat inbox data.
// Creates the portal synchronously if it doesn't exist.
func (tc *TwitterClient) syncXChatChannel(ctx context.Context, item *response.XChatInboxItem, users map[string]*types.User) {
	log := zerolog.Ctx(ctx)

	conv := tc.xchatItemToConversation(ctx, item, users)
	if conv == nil {
		return
	}

	portalKey := tc.MakePortalKey(conv)

	// Get or create portal in database
	portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
	if err != nil {
		log.Warn().Err(err).
			Str("conversation_id", conv.ConversationID).
			Msg("Failed to get/create portal")
		return
	}

	// XChat conversations are socially trusted (not message requests).
	// Note: Independent of whether encryption keys are available yet.
	meta := portal.Metadata.(*PortalMetadata)
	if !meta.Trusted {
		meta.Trusted = true
		if err := portal.Save(ctx); err != nil {
			log.Warn().Err(err).
				Str("conversation_id", conv.ConversationID).
				Msg("Failed to save portal metadata with Trusted=true")
		}
	}

	// Ensure a backfill task exists even if we don't end up emitting a ChatInfoChange.
	// Beeper scrollback relies on the backfill task existing for the portal.
	if portal.MXID != "" {
		chatInfo := tc.xchatItemToChatInfo(ctx, item, users, conv)
		if chatInfo.CanBackfill {
			if err := tc.connector.br.DB.BackfillTask.EnsureExists(ctx, portal.PortalKey, tc.userLogin.ID); err != nil {
				log.Warn().Err(err).
					Str("conversation_id", conv.ConversationID).
					Msg("Failed to ensure backfill task exists")
			} else {
				tc.connector.br.WakeupBackfillQueue()
			}
		}
	}

	// Create Matrix room if it doesn't exist
	if portal.MXID == "" {
		chatInfo := tc.xchatItemToChatInfo(ctx, item, users, conv)
		err = portal.CreateMatrixRoom(ctx, tc.userLogin, chatInfo)
		if err != nil {
			log.Warn().Err(err).
				Str("conversation_id", conv.ConversationID).
				Msg("Failed to create Matrix room")
			return
		}
		// Register backfill task for the newly created room
		if chatInfo.CanBackfill {
			if err := tc.connector.br.DB.BackfillTask.EnsureExists(ctx, portal.PortalKey, tc.userLogin.ID); err != nil {
				log.Warn().Err(err).
					Str("conversation_id", conv.ConversationID).
					Msg("Failed to ensure backfill task exists for new room")
			} else {
				tc.connector.br.WakeupBackfillQueue()
			}
		}
	} else {
		chatInfo := tc.xchatItemToChatInfo(ctx, item, users, conv)
		if (chatInfo.Name != nil && *chatInfo.Name != "") || chatInfo.Avatar != nil {
			tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
				EventMeta: simplevent.EventMeta{
					Type:      bridgev2.RemoteEventChatInfoChange,
					PortalKey: portal.PortalKey,
					Timestamp: time.Now(),
				},
				ChatInfoChange: &bridgev2.ChatInfoChange{
					ChatInfo: chatInfo,
				},
			})
		}
	}

	log.Debug().
		Str("conversation_id", conv.ConversationID).
		Stringer("portal_mxid", portal.MXID).
		Msg("XChat channel synced")
}

// xchatItemToConversation converts an XChatInboxItem to a types.Conversation.
func (tc *TwitterClient) xchatItemToConversation(ctx context.Context, item *response.XChatInboxItem, users map[string]*types.User) *types.Conversation {
	detail := item.ConversationDetail

	conv := &types.Conversation{
		ConversationID: detail.ConversationID,
		Trusted:        true, // XChat conversations are always trusted
		Muted:          detail.IsMuted,
	}

	// Determine conversation type based on participants
	if len(detail.ParticipantsResults) > 2 || detail.GroupMetadata != nil {
		conv.Type = ConversationTypeGroupDM
	} else {
		conv.Type = ConversationTypeOneToOne
	}

	// Build participants list
	for _, p := range detail.ParticipantsResults {
		conv.Participants = append(conv.Participants, types.Participant{
			UserID: p.RestID,
		})
	}

	// Set sort timestamp from group metadata if available
	if detail.GroupMetadata != nil && detail.GroupMetadata.UpdatedAtMsec != "" {
		conv.SortTimestamp = detail.GroupMetadata.UpdatedAtMsec
		conv.AvatarImageHttps = detail.GroupMetadata.GroupAvatarURL
		if name := tc.decryptGroupName(ctx, detail.ConversationID, detail.GroupMetadata.GroupName); name != "" {
			conv.Name = name
		}
	}

	return conv
}

// xchatItemToChatInfo converts an XChatInboxItem to bridgev2 chat info.
func (tc *TwitterClient) xchatItemToChatInfo(ctx context.Context, item *response.XChatInboxItem, users map[string]*types.User, conv *types.Conversation) *bridgev2.ChatInfo {
	log := zerolog.Ctx(ctx)
	detail := item.ConversationDetail

	log.Debug().
		Str("conversation_id", detail.ConversationID).
		Int("participants_count", len(detail.ParticipantsResults)).
		Int("group_members_count", len(detail.GroupMembersResults)).
		Int("users_map_count", len(users)).
		Msg("xchatItemToChatInfo building member list")

	isGroup := len(detail.ParticipantsResults) > 2 || detail.GroupMetadata != nil

	memberMap := make(bridgev2.ChatMemberMap, len(detail.ParticipantsResults))
	for _, p := range detail.ParticipantsResults {
		var userInfo *bridgev2.UserInfo
		// First try inline Result from participants_results
		if p.Result != nil {
			user := twittermeow.ConvertXChatUserToUser(p.Result)
			userInfo = tc.connector.wrapUserInfo(tc.client, user)
		}
		// Then try users map if provided
		if userInfo == nil && users != nil {
			if user, ok := users[p.RestID]; ok {
				userInfo = tc.connector.wrapUserInfo(tc.client, user)
			}
		}
		// Finally fall back to user cache
		if userInfo == nil {
			tc.userCacheLock.RLock()
			if user, ok := tc.userCache[p.RestID]; ok {
				userInfo = tc.connector.wrapUserInfo(tc.client, user)
			}
			tc.userCacheLock.RUnlock()
		}
		log.Debug().
			Str("rest_id", p.RestID).
			Bool("has_result", p.Result != nil).
			Bool("has_user_info", userInfo != nil).
			Msg("xchatItemToChatInfo adding participant")
		memberMap.Set(bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(p.RestID),
			UserInfo:    userInfo,
		})
	}

	log.Debug().
		Int("final_member_count", len(memberMap)).
		Msg("xchatItemToChatInfo finished building members")

	// MessageRequest is true for untrusted conversations (message requests)
	var messageRequest *bool
	if conv != nil {
		messageRequest = ptr.Ptr(!conv.Trusted)
	}

	info := &bridgev2.ChatInfo{
		Members: &bridgev2.ChatMemberList{
			IsFull:           true,
			TotalMemberCount: len(detail.ParticipantsResults),
			MemberMap:        memberMap,
		},
		CanBackfill:    true,
		MessageRequest: messageRequest,
	}

	if isGroup {
		// TODO this is wrong, should be default
		info.Type = ptr.Ptr(database.RoomTypeGroupDM)
		if conv != nil && conv.AvatarImageHttps != "" {
			info.Avatar = tc.makeGroupAvatar(conv.ConversationID, conv.AvatarImageHttps, "")
		} else if detail.GroupMetadata != nil && detail.GroupMetadata.GroupAvatarURL != "" {
			info.Avatar = tc.makeGroupAvatar(conv.ConversationID, detail.GroupMetadata.GroupAvatarURL, "")
		}
		if conv != nil && conv.Name != "" {
			info.Name = &conv.Name
		} else if detail.GroupMetadata != nil {
			if name := tc.decryptGroupName(ctx, detail.ConversationID, detail.GroupMetadata.GroupName); name != "" {
				info.Name = &name
			}
		}
	} else {
		info.Type = ptr.Ptr(database.RoomTypeDM)
	}

	return info
}

func (tc *TwitterClient) decryptGroupName(ctx context.Context, conversationID, encName string) string {
	if encName == "" {
		return ""
	}

	keyVersion := ""
	if parts := strings.SplitN(encName, ":", 2); len(parts) == 2 && parts[0] != "" {
		keyVersion = parts[0]
		encName = parts[1]
	}

	km := tc.client.GetKeyManager()
	var convKey *crypto.ConversationKey
	var err error

	if keyVersion != "" {
		convKey, err = km.GetConversationKey(ctx, conversationID, keyVersion)
		if err != nil {
			zerolog.Ctx(ctx).Warn().
				Err(err).
				Str("conversation_id", conversationID).
				Str("key_version", keyVersion).
				Msg("Failed to get conversation key for group name by version")
			return ""
		}
	}

	if convKey == nil || len(convKey.Key) == 0 {
		if keyVersion != "" {
			zerolog.Ctx(ctx).Debug().
				Str("conversation_id", conversationID).
				Str("key_version", keyVersion).
				Msg("Conversation key with required version missing; cannot decrypt group name")
			return ""
		}
		convKey, err = km.GetLatestConversationKey(ctx, conversationID)
		if err != nil || convKey == nil || len(convKey.Key) == 0 {
			zerolog.Ctx(ctx).Debug().
				Str("conversation_id", conversationID).
				Str("key_version", keyVersion).
				Msg("No conversation key available to decrypt group name")
			return ""
		}
	}

	var ciphertext []byte
	if dec, err := base64.StdEncoding.DecodeString(encName); err == nil {
		ciphertext = dec
	} else if dec, err := base64.RawStdEncoding.DecodeString(encName); err == nil {
		ciphertext = dec
	} else {
		zerolog.Ctx(ctx).Warn().
			Err(err).
			Str("conversation_id", conversationID).
			Msg("Failed to base64 decode encrypted group name")
		return ""
	}

	plaintext, err := crypto.SecretboxDecrypt(ciphertext, convKey.Key)
	if err != nil {
		zerolog.Ctx(ctx).Warn().
			Err(err).
			Str("conversation_id", conversationID).
			Msg("Failed to decrypt group name with conversation key")
		return ""
	}
	return string(plaintext)
}

// syncUntrustedChannels fetches and syncs untrusted (message request) conversations via the REST API.
func (tc *TwitterClient) syncUntrustedChannels(ctx context.Context) {
	log := zerolog.Ctx(ctx)

	reqQuery := ptr.Ptr(payload.DMRequestQuery{}.Default())
	// Include low quality / untrusted conversations (message requests)
	reqQuery.FilterLowQuality = false
	initialInboxState, err := tc.client.GetInitialInboxState(ctx, reqQuery)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch initial inbox state for untrusted conversations")
		return
	}

	inbox := initialInboxState.InboxInitialState
	if inbox == nil {
		log.Debug().Msg("No inbox data in initial state response")
		return
	}

	// Set the polling cursor for REST API polling
	if inbox.Cursor != "" {
		session := tc.client.GetSession()
		if session.PollingCursor == "" {
			session.PollingCursor = inbox.Cursor
			log.Debug().Str("cursor", inbox.Cursor).Msg("Initialized polling cursor from inbox state")
		}
	}

	// Update ghost info for users (ensures profile pictures are visible)
	tc.updateTwitterUserInfo(ctx, inbox)

	// Cache users from inbox
	tc.userCacheLock.Lock()
	for userID, user := range inbox.Users {
		tc.userCache[userID] = user
	}
	tc.userCacheLock.Unlock()

	// Process only untrusted conversations (message requests)
	untrustedCount := 0
	trustedCount := 0
	for _, conv := range inbox.SortedConversations() {
		if conv.Trusted {
			trustedCount++
			continue // Skip trusted - handled by XChat
		}
		untrustedCount++
		log.Debug().
			Str("conversation_id", conv.ConversationID).
			Bool("trusted", conv.Trusted).
			Bool("low_quality", conv.LowQuality).
			Str("type", string(conv.Type)).
			Msg("Processing untrusted conversation")
		tc.syncUntrustedConversation(ctx, conv, inbox)
	}

	log.Info().
		Int("untrusted_conversations", untrustedCount).
		Int("trusted_conversations", trustedCount).
		Int("total_conversations", len(inbox.Conversations)).
		Msg("Finished syncing untrusted conversations")
}

// syncUntrustedConversation syncs a single untrusted conversation.
func (tc *TwitterClient) syncUntrustedConversation(ctx context.Context, conv *types.Conversation, inbox *response.TwitterInboxData) {
	log := zerolog.Ctx(ctx)

	portalKey := tc.MakePortalKey(conv)

	portal, err := tc.connector.br.GetPortalByKey(ctx, portalKey)
	if err != nil {
		log.Warn().Err(err).
			Str("conversation_id", conv.ConversationID).
			Msg("Failed to get/create portal for untrusted conversation")
		return
	}

	// Don't downgrade trust status - only set false if not already trusted
	meta := portal.Metadata.(*PortalMetadata)
	if !meta.Trusted {
		// Keep as untrusted (Trusted stays false)
		// No need to save since Trusted=false is the zero value
	}

	chatInfo := tc.conversationToChatInfo(conv, inbox)

	// Create Matrix room if it doesn't exist
	if portal.MXID == "" {
		err = portal.CreateMatrixRoom(ctx, tc.userLogin, chatInfo)
		if err != nil {
			log.Warn().Err(err).
				Str("conversation_id", conv.ConversationID).
				Msg("Failed to create Matrix room for untrusted conversation")
			return
		}
	}

	// Process messages for this conversation from inbox entries
	if inbox != nil {
		tc.processUntrustedMessages(ctx, conv.ConversationID, inbox)
	}

	log.Debug().
		Str("conversation_id", conv.ConversationID).
		Bool("trusted", conv.Trusted).
		Msg("Synced untrusted conversation")
}

// processUntrustedMessages processes message entries for an untrusted conversation.
func (tc *TwitterClient) processUntrustedMessages(ctx context.Context, conversationID string, inbox *response.TwitterInboxData) {
	log := zerolog.Ctx(ctx)

	for _, entry := range inbox.Entries {
		parsed := entry.ParseWithErrorLog(log)
		if parsed == nil {
			continue
		}

		// Only process messages for this conversation
		msg, ok := parsed.(*types.Message)
		if !ok {
			continue
		}
		if msg.ConversationID != conversationID {
			continue
		}

		// Queue the message event
		tc.HandlePollingEvent(msg, inbox)
	}
}

// conversationToChatInfo converts a REST API conversation to bridgev2 chat info.
func (tc *TwitterClient) conversationToChatInfo(conv *types.Conversation, inbox *response.TwitterInboxData) *bridgev2.ChatInfo {
	memberMap := make(bridgev2.ChatMemberMap, len(conv.Participants))
	for _, participant := range conv.Participants {
		var userInfo *bridgev2.UserInfo
		if inbox != nil {
			if user, ok := inbox.Users[participant.UserID]; ok {
				userInfo = tc.connector.wrapUserInfo(tc.client, user)
			}
		}
		if userInfo == nil {
			tc.userCacheLock.RLock()
			if user, ok := tc.userCache[participant.UserID]; ok {
				userInfo = tc.connector.wrapUserInfo(tc.client, user)
			}
			tc.userCacheLock.RUnlock()
		}
		memberMap.Set(bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(participant.UserID),
			UserInfo:    userInfo,
		})
	}

	messageRequest := !conv.Trusted

	info := &bridgev2.ChatInfo{
		Members: &bridgev2.ChatMemberList{
			IsFull:           true,
			TotalMemberCount: len(conv.Participants),
			MemberMap:        memberMap,
		},
		CanBackfill:    true,
		MessageRequest: &messageRequest,
	}

	isGroup := conv.Type == ConversationTypeGroupDM
	if isGroup {
		info.Type = ptr.Ptr(database.RoomTypeDefault)
		if conv.AvatarImageHttps != "" {
			info.Avatar = makeAvatar(tc.client, conv.AvatarImageHttps)
		}
		if conv.Name != "" {
			info.Name = &conv.Name
		}
	} else {
		info.Type = ptr.Ptr(database.RoomTypeDM)
	}

	return info
}
