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
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// previewString returns a truncated string preview for logging.
func previewString(b []byte, max int) string {
	if len(b) == 0 {
		return ""
	}
	if len(b) > max {
		b = b[:max]
	}
	return string(b)
}

// previewBase64 returns a truncated base64 preview for logging.
func previewBase64(b []byte, maxChars int) string {
	if len(b) == 0 {
		return ""
	}
	enc := base64.StdEncoding.EncodeToString(b)
	if len(enc) > maxChars {
		return enc[:maxChars]
	}
	return enc
}

// previewHex returns a truncated hex preview for logging.
func previewHex(b []byte, maxBytes int) string {
	if len(b) > maxBytes {
		b = b[:maxBytes]
	}
	return hex.EncodeToString(b)
}

// getLatestConversationKey retrieves the latest conversation key for a conversation.
func (tc *TwitterClient) getLatestConversationKey(ctx context.Context, conversationID string) *crypto.ConversationKey {
	km := tc.client.GetKeyManager()
	if km == nil {
		return nil
	}
	key, err := km.GetLatestConversationKey(ctx, conversationID)
	if err != nil || key == nil {
		return nil
	}
	return key
}

func (tc *TwitterClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	if portal == nil {
		return nil, fmt.Errorf("portal is nil")
	}

	conversationID := ParsePortalID(portal.PortalKey.ID)
	item, users, err := tc.fetchConversationData(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	conv := tc.xchatItemToConversation(ctx, item, users)
	if conv == nil {
		return nil, fmt.Errorf("failed to build conversation from fetched data")
	}

	return tc.xchatItemToChatInfo(ctx, item, users, conv), nil
}

func (tc *TwitterClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	userID := ParseUserID(ghost.ID)
	userInfo := tc.getCachedUserInfo(userID)
	if userInfo == nil {
		if err := tc.ensureUsersInCacheByID(ctx, []string{userID}); err != nil {
			return nil, err
		}
		userInfo = tc.getCachedUserInfo(userID)
		if userInfo == nil {
			return nil, fmt.Errorf("failed to find user info in cache by id: %s", ghost.ID)
		}
	}
	return userInfo, nil
}

func (tc *TwitterClient) getCachedUserInfo(userID string) *bridgev2.UserInfo {
	tc.userCacheLock.RLock()
	defer tc.userCacheLock.RUnlock()
	var userinfo *bridgev2.UserInfo
	if userCacheEntry, ok := tc.userCache[userID]; ok {
		userinfo = tc.connector.wrapUserInfo(tc.client, userCacheEntry)
	}
	return userinfo
}

func (tc *TwitterClient) ensureUsersInCacheByID(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	uniq := make(map[string]struct{}, len(ids))
	missing := make([]string, 0, len(ids))

	tc.userCacheLock.RLock()
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, seen := uniq[id]; seen {
			continue
		}
		uniq[id] = struct{}{}
		if _, ok := tc.userCache[id]; !ok {
			missing = append(missing, id)
		}
	}
	tc.userCacheLock.RUnlock()

	if len(missing) == 0 {
		return nil
	}

	// Keep request sizes reasonable (GraphQL variables can get large quickly).
	const batchSize = 100
	for i := 0; i < len(missing); i += batchSize {
		batch := missing[i:min(i+batchSize, len(missing))]
		resp, err := tc.client.GetUsersByIdsForXChat(ctx, payload.NewGetUsersByIdsForXChatVariables(batch))
		if err != nil {
			return err
		}
		if len(resp.Errors) > 0 && resp.Errors[0].Message != "" {
			return fmt.Errorf("GetUsersByIdsForXChat error: %s", resp.Errors[0].Message)
		}

		tc.userCacheLock.Lock()
		for _, r := range resp.Data.GetMemberResults.Results {
			if r.MemberResults == nil || r.MemberResults.Result == nil {
				continue
			}
			user := twittermeow.ConvertXChatUserToUser(r.MemberResults.Result)
			if user == nil || user.IDStr == "" {
				continue
			}
			tc.userCache[user.IDStr] = user
		}
		tc.userCacheLock.Unlock()
	}

	return nil
}

func (tc *TwitterConnector) wrapUserInfo(cli *twittermeow.Client, user *types.User) *bridgev2.UserInfo {
	avatarURL := user.ProfileImageURL
	if avatarURL == "" {
		avatarURL = user.ProfileImageURLHTTPS
	}
	return &bridgev2.UserInfo{
		Name:        ptr.Ptr(tc.Config.FormatDisplayname(user.ScreenName, user.Name)),
		Avatar:      makeAvatar(cli, avatarURL),
		Identifiers: []string{fmt.Sprintf("twitter:%s", user.ScreenName)},
	}
}

func makeAvatar(cli *twittermeow.Client, avatarURL string) *bridgev2.Avatar {
	return &bridgev2.Avatar{
		ID: networkid.AvatarID(avatarURL),
		Get: func(ctx context.Context) ([]byte, error) {
			resp, err := downloadFile(ctx, cli, avatarURL)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 400 {
				return nil, fmt.Errorf("failed to download avatar: HTTP %d", resp.StatusCode)
			}
			return io.ReadAll(resp.Body)
		},
		Remove: avatarURL == "",
	}
}

// getTrustedChatInfo returns a ChatInfo with MessageRequest: false for a conversation that became trusted.
func (tc *TwitterClient) getTrustedChatInfo(ctx context.Context, conversationID string) *bridgev2.ChatInfo {
	log := zerolog.Ctx(ctx)

	// Try to fetch conversation data via XChat API first
	item, users, err := tc.fetchConversationData(ctx, conversationID)
	if err == nil && item != nil {
		log.Debug().
			Str("conversation_id", conversationID).
			Int("participants_count", len(item.ConversationDetail.ParticipantsResults)).
			Int("users_count", len(users)).
			Msg("getTrustedChatInfo fetched conversation data via XChat API")

		// If XChat API returns empty participants, fall back to building from conversation ID
		if len(item.ConversationDetail.ParticipantsResults) == 0 {
			log.Warn().
				Str("conversation_id", conversationID).
				Msg("XChat API returned empty participants, falling back to conversation ID parsing")
		} else {
			conv := tc.xchatItemToConversation(ctx, item, users)
			if conv != nil {
				// Override Trusted to true since this conversation was just accepted
				conv.Trusted = true
				return tc.xchatItemToChatInfo(ctx, item, users, conv)
			}
			log.Warn().
				Str("conversation_id", conversationID).
				Msg("xchatItemToConversation returned nil")
		}
	}

	log.Debug().
		Str("conversation_id", conversationID).
		Err(err).
		Bool("item_nil", item == nil).
		Msg("Could not fetch conversation data via XChat API for trust event, building from conversation ID")

	// XChat API failed or returned empty participants - build ChatInfo from conversation ID and user cache
	// This ensures we still have member info even if the conversation hasn't fully migrated to XChat
	return tc.buildChatInfoFromConversationID(ctx, conversationID)
}

// buildChatInfoFromConversationID builds a ChatInfo from a conversation ID by parsing user IDs
// and fetching user info from the cache. This is used as a fallback when XChat API fails.
func (tc *TwitterClient) buildChatInfoFromConversationID(ctx context.Context, conversationID string) *bridgev2.ChatInfo {
	log := zerolog.Ctx(ctx)
	messageRequest := false

	// Normalize to colon format for parsing 1:1 DM conversation IDs
	conversationID = NormalizeConversationID(conversationID)

	// Parse user IDs from conversation ID
	// 1:1 DMs use format "userID1:userID2" (lower ID first)
	// Group chats use format "g[snowflake]"
	if strings.HasPrefix(conversationID, "g") {
		// Group chat - we can't easily determine members from the conversation ID
		log.Debug().
			Str("conversation_id", conversationID).
			Msg("Group chat trust event - returning minimal ChatInfo")
		return &bridgev2.ChatInfo{
			MessageRequest: &messageRequest,
		}
	}

	// Parse 1:1 DM conversation ID
	parts := strings.Split(conversationID, ":")
	if len(parts) != 2 {
		log.Warn().
			Str("conversation_id", conversationID).
			Msg("Could not parse 1:1 conversation ID - returning minimal ChatInfo")
		return &bridgev2.ChatInfo{
			MessageRequest: &messageRequest,
		}
	}

	// Try to ensure we have both users in the cache
	if err := tc.ensureUsersInCacheByID(ctx, parts); err != nil {
		log.Warn().
			Err(err).
			Str("conversation_id", conversationID).
			Msg("Failed to fetch users for trust event")
	}

	// Build member map from parsed user IDs
	memberMap := make(bridgev2.ChatMemberMap, len(parts))
	for _, userID := range parts {
		var userInfo *bridgev2.UserInfo
		tc.userCacheLock.RLock()
		if user, ok := tc.userCache[userID]; ok {
			userInfo = tc.connector.wrapUserInfo(tc.client, user)
		}
		tc.userCacheLock.RUnlock()

		memberMap.Set(bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(userID),
			UserInfo:    userInfo,
		})
	}

	log.Debug().
		Str("conversation_id", conversationID).
		Int("member_count", len(parts)).
		Msg("Built ChatInfo from conversation ID for trust event")

	return &bridgev2.ChatInfo{
		Members: &bridgev2.ChatMemberList{
			IsFull:           true,
			TotalMemberCount: len(parts),
			MemberMap:        memberMap,
		},
		MessageRequest: &messageRequest,
	}
}

// makeGroupAvatar downloads and attempts to decrypt a group avatar using the conversation key.
// If decryption fails or no key is found, it falls back to returning the raw bytes.
func (tc *TwitterClient) makeGroupAvatar(conversationID, avatarURL, keyVersion string) *bridgev2.Avatar {
	if avatarURL == "" {
		return &bridgev2.Avatar{Remove: true}
	}

	return &bridgev2.Avatar{
		ID: networkid.AvatarID(avatarURL),
		Get: func(ctx context.Context) ([]byte, error) {
			logger := zerolog.Ctx(ctx).With().Str("conversation_id", conversationID).Str("url", avatarURL).Logger()

			resp, body, err := tc.client.FetchRaw(ctx, avatarURL)
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to download group avatar")
				return nil, err
			}

			if resp != nil {
				logger.Info().
					Int("status_code", resp.StatusCode).
					Int("body_len", len(body)).
					Str("content_type", resp.Header.Get("content-type")).
					Str("body_prefix_str", previewString(body, 200)).
					Str("body_prefix_b64", previewBase64(body, 50)).
					Msg("Fetched group avatar")
			}

			if len(body) == 0 {
				return body, nil
			}

			convKey := tc.getLatestConversationKey(ctx, conversationID)
			if convKey == nil || len(convKey.Key) == 0 {
				logger.Warn().Msg("No conversation key available for group avatar; returning raw bytes")
				return body, nil
			}

			pt, err := crypto.SecretstreamDecrypt(body, convKey.Key)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("key_version", convKey.KeyVersion).
					Str("body_prefix_hex", previewHex(body, 16)).
					Msg("Failed to decrypt group avatar; returning raw bytes")
				return body, nil
			}
			logger.Info().Msg("Successfully decrypted group avatar")
			return pt, nil
		},
		Remove: false,
	}
}

// getOrFetchChatInfoForPolling gets ChatInfo for a conversation, trying multiple sources:
// 1. From the polling inbox data (if conversation is present)
// 2. By fetching via REST API FetchConversationContext
// 3. By building minimal info from conversation ID (for 1:1 DMs)
func (tc *TwitterClient) getOrFetchChatInfoForPolling(ctx context.Context, conversationID string, inbox *response.TwitterInboxData) *bridgev2.ChatInfo {
	log := zerolog.Ctx(ctx).With().
		Str("conversation_id", conversationID).
		Str("handler", "polling").
		Logger()

	// Try 1: Get from inbox if present
	if inbox != nil {
		if conv := inbox.GetConversationByID(conversationID); conv != nil {
			log.Debug().Msg("Found conversation in polling inbox")
			return tc.conversationToChatInfo(conv, inbox)
		}
	}

	// Try 2: Fetch via REST API
	log.Debug().Msg("Conversation not in inbox, fetching via REST API")
	restConvID := ConvertConversationIDToREST(conversationID)
	reqQuery := payload.DMRequestQuery{}.Default()
	resp, err := tc.client.FetchConversationContext(ctx, restConvID, &reqQuery, payload.CONTEXT_FETCH_DM_CONVERSATION)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to fetch conversation via REST API")
	} else if resp.ConversationTimeline != nil {
		// Cache users from the fetched data
		if resp.ConversationTimeline.Users != nil {
			tc.userCacheLock.Lock()
			for userID, user := range resp.ConversationTimeline.Users {
				tc.userCache[userID] = user
			}
			tc.userCacheLock.Unlock()
		}

		// Try to get conversation from fetched data
		if conv := resp.ConversationTimeline.GetConversationByID(restConvID); conv != nil {
			log.Debug().Msg("Got conversation from REST API fetch")
			return tc.conversationToChatInfo(conv, resp.ConversationTimeline)
		}
		log.Warn().Msg("REST API returned data but conversation not found in response")
	}

	// Try 3: Build minimal ChatInfo from conversation ID (works for 1:1 DMs)
	// For group chats, we can't determine members from the ID alone, so return nil
	if strings.HasPrefix(conversationID, "g") {
		log.Warn().Msg("Cannot build ChatInfo for group chat without conversation data")
		return nil
	}
	log.Debug().Msg("Building minimal ChatInfo from conversation ID")
	return tc.buildChatInfoFromConversationID(ctx, conversationID)
}
