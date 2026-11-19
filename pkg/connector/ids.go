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

	"github.com/google/uuid"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func ConvertConversationIDToREST(conversationID string) string {
	if strings.HasPrefix(conversationID, "g") {
		return conversationID
	}
	return strings.ReplaceAll(conversationID, ":", "-")
}

func ConvertConversationIDToRESTAPI(conversationID string) string {
	if restKey, ok := restGroupPortalAliasKey(networkid.PortalKey{ID: networkid.PortalID(conversationID)}); ok {
		return string(restKey.ID)
	}
	return ConvertConversationIDToREST(conversationID)
}

func MakePortalID(conversationID string) networkid.PortalID {
	// Backwards compat for old pre-xchat portal IDs
	return networkid.PortalID(ConvertConversationIDToREST(conversationID))
}

func ParsePortalID(portalID networkid.PortalID) string {
	if strings.HasPrefix(string(portalID), "g") {
		return string(portalID)
	}
	return strings.ReplaceAll(string(portalID), "-", ":")
}

func xchatGroupPortalAliasKey(restKey networkid.PortalKey) (networkid.PortalKey, bool) {
	restID := string(restKey.ID)
	if _, err := strconv.ParseUint(restID, 10, 64); err != nil {
		return networkid.PortalKey{}, false
	}
	xchatKey := restKey
	xchatKey.ID = networkid.PortalID("g" + restID)
	return xchatKey, true
}

func restGroupPortalAliasKey(xchatKey networkid.PortalKey) (networkid.PortalKey, bool) {
	xchatID := string(xchatKey.ID)
	if len(xchatID) < 2 || xchatID[0] != 'g' {
		return networkid.PortalKey{}, false
	}
	restID := xchatID[1:]
	if _, err := strconv.ParseUint(restID, 10, 64); err != nil {
		return networkid.PortalKey{}, false
	}
	restKey := xchatKey
	restKey.ID = networkid.PortalID(restID)
	return restKey, true
}

func MakeUserID(userID string) networkid.UserID {
	return networkid.UserID(userID)
}

func MakeUserLoginID(userID string) networkid.UserLoginID {
	return networkid.UserLoginID(userID)
}

func ParseUserLoginID(userLoginID networkid.UserLoginID) string {
	return string(userLoginID)
}

func ParseUserID(userID networkid.UserID) string {
	return string(userID)
}

func MakeMessageID(seqID string) networkid.MessageID {
	return networkid.MessageID(seqID)
}

func ParseMessageID(messageID networkid.MessageID) string {
	return string(messageID)
}

func (tc *TwitterClient) MakePortalKey(conv *types.Conversation) networkid.PortalKey {
	return tc.MakePortalKeyFromID(conv.ConversationID)
}

func (tc *TwitterClient) MakePortalKeyFromID(conversationID string) networkid.PortalKey {
	return MakePortalKeyForConversation(conversationID, tc.userLogin.ID, tc.connector.br.Config.SplitPortals)
}

// MakePortalKeyForConversation creates a portal key using the same logic as MakePortalKeyFromID.
// This is used by the keystore which doesn't have direct access to TwitterClient.
func MakePortalKeyForConversation(conversationID string, loginID networkid.UserLoginID, splitPortals bool) (key networkid.PortalKey) {
	key.ID = MakePortalID(conversationID)
	// 1:1 DM conversation IDs use `:` as delimiter between user IDs
	if strings.Contains(string(key.ID), "-") || splitPortals {
		key.Receiver = loginID
	}
	return
}

func (tc *TwitterClient) MakeEventSender(userID string) bridgev2.EventSender {
	return bridgev2.EventSender{
		IsFromMe:    userID == tc.client.GetCurrentUserID(),
		SenderLogin: MakeUserLoginID(userID),
		Sender:      MakeUserID(userID),
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
	KeyInfo        *crypto.AttachmentKeyInfo
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

	if info.KeyInfo != nil {
		mediaID = appendKeyInfoField(mediaID, info.KeyInfo.SenderID)
		mediaID = appendKeyInfoField(mediaID, info.KeyInfo.MessageID)
		mediaID = appendAttachmentIDField(mediaID, info.KeyInfo.MessageID, info.KeyInfo.AttachmentID)
	}

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
	mediaInfo.UserID = MakeUserLoginID(strconv.FormatUint(uID, 10))

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
	info.UserID = MakeUserLoginID(strconv.FormatUint(uID, 10))

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

	if buf.Len() > 0 {
		keyInfo := &crypto.AttachmentKeyInfo{}
		if keyInfo.SenderID, err = readKeyInfoField(buf, nil); err != nil {
			return nil, fmt.Errorf("read sender ID: %w", err)
		} else if keyInfo.MessageID, err = readKeyInfoField(buf, nil); err != nil {
			return nil, fmt.Errorf("read message ID: %w", err)
		} else if keyInfo.AttachmentID, err = readKeyInfoField(buf, &keyInfo.MessageID); err != nil {
			return nil, fmt.Errorf("read attachment ID: %w", err)
		}
		info.KeyInfo = keyInfo
	}

	return info, nil
}

const (
	keyInfoFieldString byte = iota
	keyInfoFieldUint
	keyInfoFieldUUID
	keyInfoFieldMessageID
	keyInfoFieldMessageIDIndex
)

func parseCanonicalUint(val string) (uint64, bool) {
	u, err := strconv.ParseUint(val, 10, 64)
	return u, err == nil && strconv.FormatUint(u, 10) == val
}

func appendKeyInfoField(buf []byte, val string) []byte {
	if u, ok := parseCanonicalUint(val); ok {
		return binary.AppendUvarint(append(buf, keyInfoFieldUint), u)
	} else if u, err := uuid.Parse(val); err == nil && u.String() == val {
		return append(append(buf, keyInfoFieldUUID), u[:]...)
	}
	buf = binary.AppendUvarint(append(buf, keyInfoFieldString), uint64(len(val)))
	return append(buf, val...)
}

func appendAttachmentIDField(buf []byte, messageID, attachmentID string) []byte {
	if attachmentID == messageID {
		return append(buf, keyInfoFieldMessageID)
	} else if suffix, ok := strings.CutPrefix(attachmentID, messageID+"#"); ok {
		if index, ok := parseCanonicalUint(suffix); ok {
			return binary.AppendUvarint(append(buf, keyInfoFieldMessageIDIndex), index)
		}
	}
	return appendKeyInfoField(buf, attachmentID)
}

func readKeyInfoField(buf *bytes.Reader, messageID *string) (string, error) {
	fieldType, err := buf.ReadByte()
	if err != nil {
		return "", err
	}
	switch fieldType {
	case keyInfoFieldString:
		size, err := binary.ReadUvarint(buf)
		if err != nil {
			return "", err
		} else if size > uint64(buf.Len()) {
			return "", io.ErrUnexpectedEOF
		}
		bs := make([]byte, size)
		if _, err = io.ReadFull(buf, bs); err != nil {
			return "", err
		}
		return string(bs), nil
	case keyInfoFieldUint:
		u, err := binary.ReadUvarint(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatUint(u, 10), nil
	case keyInfoFieldUUID:
		var u uuid.UUID
		if _, err = io.ReadFull(buf, u[:]); err != nil {
			return "", err
		}
		return u.String(), nil
	case keyInfoFieldMessageID:
		if messageID != nil {
			return *messageID, nil
		}
	case keyInfoFieldMessageIDIndex:
		if messageID != nil {
			index, err := binary.ReadUvarint(buf)
			if err != nil {
				return "", err
			}
			return *messageID + "#" + strconv.FormatUint(index, 10), nil
		}
	}
	return "", fmt.Errorf("unknown field type %d", fieldType)
}
