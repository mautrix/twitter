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
	"fmt"
	"io"

	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	bridgeEvt "maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) GetChatInfo(_ context.Context, _ *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	// Chat info is populated during initial sync via GetInitialXChatPage.
	// We don't support on-demand fetching via REST API as it doesn't work with XChat/group chats.
	return nil, nil
}

func (tc *TwitterClient) GetUserInfo(_ context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	userInfo := tc.getCachedUserInfo(string(ghost.ID))
	if userInfo == nil {
		return nil, fmt.Errorf("failed to find user info in cache by id: %s", ghost.ID)
	}
	return userInfo, nil
}

func (tc *TwitterClient) getCachedUserInfo(userID string) *bridgev2.UserInfo {
	tc.userCacheLock.Lock()
	defer tc.userCacheLock.Unlock()
	var userinfo *bridgev2.UserInfo
	if userCacheEntry, ok := tc.userCache[userID]; ok {
		userinfo = tc.connector.wrapUserInfo(tc.client, userCacheEntry)
	}
	return userinfo
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
		chatInfo.Avatar = makeAvatar(tc.client, conv.AvatarImageHttps)
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
		memberMap[MakeUserID(participant.UserID)] = tc.participantToChatMember(participant, inbox)
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
