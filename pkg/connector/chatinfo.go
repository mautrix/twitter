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
	"maunium.net/go/mautrix/bridgev2/database"
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

func (tc *TwitterClient) getDMChatInfoCompleteness(info *bridgev2.ChatInfo) (isDM bool, complete bool, memberCount int, missingUserInfo int) {
	if info == nil {
		return false, false, 0, 0
	}
	if info.Type == nil || *info.Type != database.RoomTypeDM {
		// Non-DM conversations don't use member-derived room naming.
		return false, true, 0, 0
	}
	if info.Members == nil || !info.Members.IsFull {
		return true, false, 0, 0
	}

	memberCount = info.Members.TotalMemberCount
	if memberCount == 0 && len(info.Members.MemberMap) > 0 {
		memberCount = len(info.Members.MemberMap)
	}
	for _, member := range info.Members.MemberMap {
		if member.UserInfo == nil {
			missingUserInfo++
			continue
		}
		if member.UserInfo.Name == nil || strings.TrimSpace(*member.UserInfo.Name) == "" {
			missingUserInfo++
		}
	}

	complete = memberCount >= 2 && missingUserInfo == 0
	return true, complete, memberCount, missingUserInfo
}

func (tc *TwitterClient) isDMChatInfoComplete(info *bridgev2.ChatInfo) bool {
	_, complete, _, _ := tc.getDMChatInfoCompleteness(info)
	return complete
}

func (tc *TwitterClient) currentUserID() string {
	currentUserID := strings.TrimSpace(tc.client.GetCurrentUserID())
	if currentUserID == "" {
		currentUserID = ParseUserLoginID(tc.userLogin.ID)
	}
	return currentUserID
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

	return tc.fetchUsersByIDAndCache(ctx, missing)
}

func (tc *TwitterClient) forceRefreshUserInCacheByID(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}
	return tc.fetchUsersByIDAndCache(ctx, []string{userID})
}

func (tc *TwitterClient) fetchUsersByIDAndCache(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// Keep request sizes reasonable (GraphQL variables can get large quickly).
	const batchSize = 100
	for i := 0; i < len(ids); i += batchSize {
		batch := ids[i:min(i+batchSize, len(ids))]
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

func parseDMParticipantIDs(conversationID string) (string, []string, bool) {
	normalizedConversationID := NormalizeConversationID(conversationID)
	if strings.HasPrefix(normalizedConversationID, "g") {
		return normalizedConversationID, nil, false
	}

	parts := strings.Split(normalizedConversationID, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return normalizedConversationID, nil, false
	}
	return normalizedConversationID, parts, true
}

func isSelfConversationID(conversationID, selfUserID string) bool {
	if selfUserID == "" {
		return false
	}
	_, parts, ok := parseDMParticipantIDs(conversationID)
	return ok && parts[0] == selfUserID && parts[1] == selfUserID
}

func (tc *TwitterClient) cacheUserForID(userID string, user *types.User) *types.User {
	if userID == "" || user == nil {
		return nil
	}
	tc.userCacheLock.Lock()
	tc.userCache[userID] = user
	tc.userCacheLock.Unlock()
	return user
}

func (tc *TwitterClient) resolveSelfDMUser(ctx context.Context, selfUserID string, users map[string]*types.User, inbox *response.TwitterInboxData) *types.User {
	if selfUserID == "" {
		return nil
	}

	if users != nil {
		if user := tc.cacheUserForID(selfUserID, users[selfUserID]); user != nil {
			return user
		}
	}

	if inbox != nil && inbox.Users != nil {
		if user := tc.cacheUserForID(selfUserID, inbox.Users[selfUserID]); user != nil {
			return user
		}
	}

	if err := tc.forceRefreshUserInCacheByID(ctx, selfUserID); err != nil {
		zerolog.Ctx(ctx).Warn().
			Err(err).
			Str("user_id", selfUserID).
			Msg("Failed to refresh current user profile for self DM metadata")
	}

	tc.userCacheLock.RLock()
	user := tc.userCache[selfUserID]
	tc.userCacheLock.RUnlock()
	return user
}

func (tc *TwitterClient) applySelfDMNameAndAvatar(ctx context.Context, conversationID string, info *bridgev2.ChatInfo, users map[string]*types.User, inbox *response.TwitterInboxData) {
	if info == nil {
		return
	}

	selfUserID := tc.currentUserID()
	if !isSelfConversationID(conversationID, selfUserID) {
		return
	}

	user := tc.resolveSelfDMUser(ctx, selfUserID, users, inbox)
	if user == nil {
		zerolog.Ctx(ctx).Warn().
			Str("conversation_id", conversationID).
			Str("user_id", selfUserID).
			Msg("Could not resolve current user profile for self DM metadata")
		return
	}

	userInfo := tc.connector.wrapUserInfo(tc.client, user)
	info.Name = userInfo.Name
	info.Avatar = userInfo.Avatar
}

func (tc *TwitterClient) buildChatMembersFromUserIDs(ctx context.Context, conversationID string, userIDs []string, inbox *response.TwitterInboxData) *bridgev2.ChatMemberList {
	log := zerolog.Ctx(ctx).With().
		Str("conversation_id", conversationID).
		Logger()

	uniq := make(map[string]struct{}, len(userIDs))
	ids := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		if id == "" {
			continue
		}
		if _, seen := uniq[id]; seen {
			continue
		}
		uniq[id] = struct{}{}
		ids = append(ids, id)
	}

	memberMap := make(bridgev2.ChatMemberMap, len(ids))
	if len(ids) == 0 {
		return &bridgev2.ChatMemberList{
			IsFull:           true,
			TotalMemberCount: 0,
			MemberMap:        memberMap,
		}
	}

	// Find which participant user IDs we can't resolve from the inbox or cache, then fetch those.
	missing := make([]string, 0, len(ids))
	tc.userCacheLock.RLock()
	for _, id := range ids {
		if inbox != nil && inbox.Users != nil {
			if u := inbox.Users[id]; u != nil {
				continue
			}
		}
		if tc.userCache[id] != nil {
			continue
		}
		missing = append(missing, id)
	}
	tc.userCacheLock.RUnlock()

	if len(missing) > 0 {
		if err := tc.ensureUsersInCacheByID(ctx, missing); err != nil {
			log.Warn().
				Err(err).
				Int("missing_users", len(missing)).
				Msg("Failed to fetch missing users for chat member list")
		}
	}

	missingAfter := 0
	tc.userCacheLock.RLock()
	for _, id := range ids {
		var userInfo *bridgev2.UserInfo
		if inbox != nil && inbox.Users != nil {
			if u := inbox.Users[id]; u != nil {
				userInfo = tc.connector.wrapUserInfo(tc.client, u)
			}
		}
		if userInfo == nil {
			if u := tc.userCache[id]; u != nil {
				userInfo = tc.connector.wrapUserInfo(tc.client, u)
			}
		}
		if userInfo == nil {
			missingAfter++
		}
		memberMap.Set(bridgev2.ChatMember{
			EventSender: tc.MakeEventSender(id),
			UserInfo:    userInfo,
		})
	}
	tc.userCacheLock.RUnlock()

	log.Debug().
		Int("participants", len(ids)).
		Int("missing_before", len(missing)).
		Int("missing_after", missingAfter).
		Msg("Built chat member list")

	return &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(ids),
		MemberMap:        memberMap,
	}
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

	normalizedConversationID, parts, ok := parseDMParticipantIDs(conversationID)
	if !ok {
		// Group chats use format "g[snowflake]".
		if strings.HasPrefix(normalizedConversationID, "g") {
			// Group chat - we can't easily determine members from the conversation ID.
			log.Debug().
				Str("conversation_id", normalizedConversationID).
				Msg("Group chat trust event - returning minimal ChatInfo")
		} else {
			log.Warn().
				Str("conversation_id", normalizedConversationID).
				Msg("Could not parse 1:1 conversation ID - returning minimal ChatInfo")
		}
		return &bridgev2.ChatInfo{
			MessageRequest: &messageRequest,
		}
	}

	log.Debug().
		Str("conversation_id", normalizedConversationID).
		Int("member_count", len(parts)).
		Msg("Built ChatInfo from conversation ID for trust event")

	info := &bridgev2.ChatInfo{
		Members:        tc.buildChatMembersFromUserIDs(ctx, normalizedConversationID, parts, nil),
		MessageRequest: &messageRequest,
	}
	tc.applySelfDMNameAndAvatar(ctx, normalizedConversationID, info, nil, nil)
	return info
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
			chatInfo := tc.conversationToChatInfo(ctx, conv, inbox)
			isDM, complete, memberCount, missingUserInfo := tc.getDMChatInfoCompleteness(chatInfo)
			if complete {
				return chatInfo
			}
			log.Info().
				Bool("is_dm", isDM).
				Int("member_count", memberCount).
				Int("missing_userinfo", missingUserInfo).
				Msg("Polling inbox chat info incomplete, falling back to REST fetch")
		}
	}

	// Try 2: Fetch via REST API
	log.Debug().Msg("Fetching conversation via REST API")
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
			chatInfo := tc.conversationToChatInfo(ctx, conv, resp.ConversationTimeline)
			isDM, complete, memberCount, missingUserInfo := tc.getDMChatInfoCompleteness(chatInfo)
			if complete {
				return chatInfo
			}
			log.Warn().
				Bool("is_dm", isDM).
				Int("member_count", memberCount).
				Int("missing_userinfo", missingUserInfo).
				Msg("REST chat info incomplete, falling back to minimal conversation-ID chat info")
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
