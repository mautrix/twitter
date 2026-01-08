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
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) MakePortalKey(conv *types.Conversation) networkid.PortalKey {
	return tc.MakePortalKeyFromID(conv.ConversationID)
}

func (tc *TwitterClient) MakePortalKeyFromID(conversationID string) networkid.PortalKey {
	return MakePortalKeyForConversation(conversationID, tc.userLogin.ID, tc.connector.br.Config.SplitPortals)
}

// MakePortalKeyForConversation creates a portal key using the same logic as MakePortalKeyFromID.
// This is used by the keystore which doesn't have direct access to TwitterClient.
func MakePortalKeyForConversation(conversationID string, loginID networkid.UserLoginID, splitPortals bool) networkid.PortalKey {
	var receiver networkid.UserLoginID
	// 1:1 DM conversation IDs use `:` as delimiter between user IDs
	if strings.Contains(conversationID, ":") || strings.HasPrefix(conversationID, "g") || splitPortals {
		receiver = loginID
	}
	return networkid.PortalKey{
		ID:       networkid.PortalID(conversationID),
		Receiver: receiver,
	}
}

func (tc *TwitterClient) MakeEventSender(userID string) bridgev2.EventSender {
	return bridgev2.EventSender{
		IsFromMe:    userID == tc.client.GetCurrentUserID(),
		SenderLogin: networkid.UserLoginID(userID),
		Sender:      networkid.UserID(userID),
	}
}

// MediaInfo stores legacy (non-encrypted) media info for direct media.
// Version 1 format.
type MediaInfo struct {
	UserID networkid.UserLoginID
	URL    string
}

// EncryptedMediaInfo stores encrypted XChat media info for direct media.
// Version 2 format.
type EncryptedMediaInfo struct {
	UserID         networkid.UserLoginID
	ConversationID string
	MediaHashKey   string
	KeyVersion     string
}

// MakeMediaID creates a version 1 media ID for legacy (non-encrypted) media.
func MakeMediaID(userID networkid.UserLoginID, URL string) networkid.MediaID {
	mediaID := []byte{1}
	uID, err := strconv.ParseUint(string(userID), 10, 64)
	if err != nil {
		panic(err)
	}
	mediaID = binary.AppendUvarint(mediaID, uID)

	bs := []byte(URL)
	mediaID = binary.AppendUvarint(mediaID, uint64(len(bs)))
	mediaID, err = binary.Append(mediaID, binary.BigEndian, bs)
	if err != nil {
		panic(err)
	}

	return mediaID
}

// MakeEncryptedMediaID creates a version 2 media ID for encrypted XChat media.
func MakeEncryptedMediaID(info EncryptedMediaInfo) networkid.MediaID {
	mediaID := []byte{2}
	uID, err := strconv.ParseUint(string(info.UserID), 10, 64)
	if err != nil {
		panic(err)
	}
	mediaID = binary.AppendUvarint(mediaID, uID)

	// Encode conversation ID
	convIDBytes := []byte(info.ConversationID)
	mediaID = binary.AppendUvarint(mediaID, uint64(len(convIDBytes)))
	mediaID = append(mediaID, convIDBytes...)

	// Encode media hash key
	hashKeyBytes := []byte(info.MediaHashKey)
	mediaID = binary.AppendUvarint(mediaID, uint64(len(hashKeyBytes)))
	mediaID = append(mediaID, hashKeyBytes...)

	// Encode key version
	keyVersionBytes := []byte(info.KeyVersion)
	mediaID = binary.AppendUvarint(mediaID, uint64(len(keyVersionBytes)))
	mediaID = append(mediaID, keyVersionBytes...)

	return mediaID
}

// ParseMediaID parses a media ID and returns either *MediaInfo (v1) or *EncryptedMediaInfo (v2).
func ParseMediaID(mediaID networkid.MediaID) (any, error) {
	buf := bytes.NewReader(mediaID)
	version := make([]byte, 1)
	_, err := io.ReadFull(buf, version)
	if err != nil {
		return nil, err
	}

	switch version[0] {
	case 1:
		return parseMediaIDV1(buf)
	case 2:
		return parseMediaIDV2(buf)
	default:
		return nil, fmt.Errorf("unknown mediaID version: %v", version[0])
	}
}

func parseMediaIDV1(buf *bytes.Reader) (*MediaInfo, error) {
	mediaInfo := &MediaInfo{}
	uID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	mediaInfo.UserID = networkid.UserLoginID(strconv.FormatUint(uID, 10))

	size, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	bs := make([]byte, size)
	_, err = io.ReadFull(buf, bs)
	if err != nil {
		return nil, err
	}
	mediaInfo.URL = string(bs)

	return mediaInfo, nil
}

func parseMediaIDV2(buf *bytes.Reader) (*EncryptedMediaInfo, error) {
	info := &EncryptedMediaInfo{}

	// Read user ID
	uID, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	info.UserID = networkid.UserLoginID(strconv.FormatUint(uID, 10))

	// Read conversation ID
	size, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	bs := make([]byte, size)
	if _, err := io.ReadFull(buf, bs); err != nil {
		return nil, err
	}
	info.ConversationID = string(bs)

	// Read media hash key
	size, err = binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	bs = make([]byte, size)
	if _, err := io.ReadFull(buf, bs); err != nil {
		return nil, err
	}
	info.MediaHashKey = string(bs)

	// Read key version
	size, err = binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	bs = make([]byte, size)
	if _, err := io.ReadFull(buf, bs); err != nil {
		return nil, err
	}
	info.KeyVersion = string(bs)

	return info, nil
}
