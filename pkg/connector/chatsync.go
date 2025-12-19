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
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
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
		Str("portal_mxid", string(portal.MXID)).
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
	detail := item.ConversationDetail

	isGroup := len(detail.ParticipantsResults) > 2 || detail.GroupMetadata != nil

	memberMap := make(map[networkid.UserID]bridgev2.ChatMember, len(detail.ParticipantsResults))
	for _, p := range detail.ParticipantsResults {
		var userInfo *bridgev2.UserInfo
		// First try inline Result from participants_results, then fall back to users map or cache
		if p.Result != nil {
			user := twittermeow.ConvertXChatUserToUser(p.Result)
			userInfo = tc.connector.wrapUserInfo(tc.client, user)
		} else if users != nil {
			if user, ok := users[p.RestID]; ok {
				userInfo = tc.connector.wrapUserInfo(tc.client, user)
			}
		} else {
			// Fall back to user cache when users map is nil
			tc.userCacheLock.RLock()
			if user, ok := tc.userCache[p.RestID]; ok {
				userInfo = tc.connector.wrapUserInfo(tc.client, user)
			}
			tc.userCacheLock.RUnlock()
		}
		memberMap[networkid.UserID(p.RestID)] = bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(p.RestID),
			UserInfo:    userInfo,
		}
	}

	info := &bridgev2.ChatInfo{
		Members: &bridgev2.ChatMemberList{
			IsFull:           true,
			TotalMemberCount: len(detail.ParticipantsResults),
			MemberMap:        memberMap,
		},
		CanBackfill: true,
	}

	if isGroup {
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
