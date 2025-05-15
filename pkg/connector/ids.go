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
	"encoding/json"
	"log"
	"strings"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) makePortalKeyFromInbox(conversationID string, inbox *response.TwitterInboxData) networkid.PortalKey {
	conv := inbox.GetConversationByID(conversationID)
	if conv != nil {
		return tc.MakePortalKey(conv)
	} else {
		return tc.MakePortalKeyFromID(conversationID)
	}
}

func (tc *TwitterClient) MakePortalKey(conv *types.Conversation) networkid.PortalKey {
	var receiver networkid.UserLoginID
	if conv.Type == types.ConversationTypeOneToOne || tc.connector.br.Config.SplitPortals {
		receiver = tc.userLogin.ID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conv.ConversationID),
		Receiver: receiver,
	}
}

func (tc *TwitterClient) MakePortalKeyFromID(conversationID string) networkid.PortalKey {
	var receiver networkid.UserLoginID
	if strings.Contains(conversationID, "-") || tc.connector.br.Config.SplitPortals {
		receiver = tc.userLogin.ID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conversationID),
		Receiver: receiver,
	}
}

func MakeUserID(userID string) networkid.UserID {
	return networkid.UserID(userID)
}

func ParseUserID(userID networkid.UserID) string {
	return string(userID)
}

func UserIDToUserLoginID(userID networkid.UserID) networkid.UserLoginID {
	return networkid.UserLoginID(userID)
}

func UserLoginIDToUserID(userID networkid.UserLoginID) networkid.UserID {
	return networkid.UserID(userID)
}

func MakeUserLoginID(userID string) networkid.UserLoginID {
	return networkid.UserLoginID(userID)
}

func ParseUserLoginID(userID networkid.UserLoginID) string {
	return string(userID)
}

func (tc *TwitterClient) MakeEventSender(userID string) bridgev2.EventSender {
	return bridgev2.EventSender{
		IsFromMe:    userID == string(tc.userLogin.ID),
		SenderLogin: MakeUserLoginID(userID),
		Sender:      MakeUserID(userID),
	}
}

type MediaInfo struct {
	UserID networkid.UserLoginID
	URL    string
}

func MakeMediaID(userID networkid.UserLoginID, URL string) networkid.MediaID {
	info := &MediaInfo{
		UserID: userID,
		URL:    URL,
	}
	id, err := json.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func ParseMediaID(mediaID networkid.MediaID) (MediaInfo, error) {
	var info MediaInfo
	err := json.Unmarshal(mediaID, &info)
	return info, err
}
