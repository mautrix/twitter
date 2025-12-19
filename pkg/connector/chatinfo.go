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

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	bridgeEvt "maunium.net/go/mautrix/event"

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

	conversationID := string(portal.PortalKey.ID)
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
	userID := string(ghost.ID)
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

func (tc *TwitterClient) conversationToChatInfo(conv *types.Conversation, inbox *response.TwitterInboxData) *bridgev2.ChatInfo {
	memberList := tc.participantsToMemberList(conv.Participants, inbox)
	var userLocal bridgev2.UserLocalPortalInfo
	if conv.Muted {
		userLocal.MutedUntil = ptr.Ptr(bridgeEvt.MutedForever)
	} else {
		userLocal.MutedUntil = ptr.Ptr(bridgev2.Unmuted)
	}
	chatInfo := &bridgev2.ChatInfo{
		Members:     memberList,
		Type:        tc.conversationTypeToRoomType(conv.Type),
		UserLocal:   &userLocal,
		CanBackfill: true,
	}

	if *chatInfo.Type != database.RoomTypeDM {
		chatInfo.Name = &conv.Name
		chatInfo.Avatar = tc.makeGroupAvatar(conv.ConversationID, conv.AvatarImageHttps, "")
	} else {
		chatInfo.Name = bridgev2.DefaultChatName
	}

	return chatInfo
}

func (tc *TwitterClient) conversationTypeToRoomType(convType types.ConversationType) *database.RoomType {
	var roomType database.RoomType
	switch convType {
	case types.ConversationTypeOneToOne:
		roomType = database.RoomTypeDM
	case types.ConversationTypeGroupDM:
		roomType = database.RoomTypeGroupDM
	}

	return &roomType
}

func (tc *TwitterClient) participantsToMemberList(participants []types.Participant, inbox *response.TwitterInboxData) *bridgev2.ChatMemberList {
	memberMap := make(map[networkid.UserID]bridgev2.ChatMember, len(participants))
	for _, participant := range participants {
		memberMap[networkid.UserID(participant.UserID)] = tc.participantToChatMember(participant, inbox)
	}
	return &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(participants),
		MemberMap:        memberMap,
	}
}

func (tc *TwitterClient) participantToChatMember(participant types.Participant, inbox *response.TwitterInboxData) bridgev2.ChatMember {
	var userInfo *bridgev2.UserInfo
	if user := inbox.GetUserByID(participant.UserID); user != nil {
		userInfo = tc.connector.wrapUserInfo(tc.client, user)
	} else {
		userInfo = tc.getCachedUserInfo(participant.UserID)
	}
	return bridgev2.ChatMember{
		EventSender: tc.MakeEventSender(participant.UserID),
		UserInfo:    userInfo,
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
			data, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return data, err
		},
		Remove: avatarURL == "",
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
