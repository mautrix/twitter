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
	"crypto/ecdh"
	"crypto/rand"
	"time"

	"go.mau.fi/util/exerrors"
	"go.mau.fi/util/random"
	"maunium.net/go/mautrix/bridgev2/database"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

func (tc *TwitterConnector) GetDBMetaTypes() database.MetaTypes {
	return database.MetaTypes{
		Reaction: nil,
		Portal: func() any {
			return &PortalMetadata{}
		},
		Message: func() any {
			return &MessageMetadata{}
		},
		Ghost: nil,
		UserLogin: func() any {
			return &UserLoginMetadata{}
		},
	}
}

// PortalMetadata stores per-conversation data in the portal.
type PortalMetadata struct {
	// Trusted indicates we can use XChat protocol (not a message request).
	// This is the sole routing flag - NOT derived from key/token presence.
	// Rules:
	//   - XChat sync → set true (never downgrade)
	//   - Untrusted REST sync → set false only if currently unset
	//   - First outbound message (REST) → set true after success
	//   - TrustConversation event → set true
	// TODO delete this and use the standard MessageRequest field in Portal
	Trusted bool `json:"trusted,omitempty"`

	// Encryption keys for XChat messages, keyed by keyVersion
	ConversationKeys map[string]*ConversationKeyData `json:"conversation_keys,omitempty"`

	// Server token for XChat API
	ConversationToken string `json:"conversation_token,omitempty"`
}

// IsTrusted returns true if this conversation is trusted (not a message request)
func (m *PortalMetadata) IsTrusted() bool {
	return m != nil && m.Trusted
}

type UserLoginMetadata struct {
	Cookies           string    `json:"cookies"`
	SecretKey         string    `json:"secret_key,omitempty"`
	SigningKey        string    `json:"signing_key,omitempty"`
	SigningKeyVersion string    `json:"signing_key_version,omitempty"`
	PushKeys          *PushKeys `json:"push_keys,omitempty"`

	Session            *twittermeow.CachedSession `json:"session,omitempty"`
	MaxUserSequenceID  string                     `json:"max_user_sequence_id,omitempty"` // Last processed sequence ID for incremental inbox fetching
	MessagePullVersion *int                       `json:"message_pull_version,omitempty"`

	// Migration tracking fields
	MigratedAt           *time.Time `json:"migrated_at,omitempty"`            // When encryption keys were first obtained via migration
	PendingEncryptedSync bool       `json:"pending_encrypted_sync,omitempty"` // True if encrypted rooms need full backfill after migration
}

// ConversationKeyData stores a conversation encryption key.
// Stored in PortalMetadata.ConversationKeys, keyed by keyVersion.
type ConversationKeyData struct {
	KeyVersion string     `json:"key_version"`
	Key        []byte     `json:"key"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type MessageMetadata struct {
	EditCount        int    `json:"edit_count,omitempty"`
	XChatClientMsgID string `json:"xchat_client_msg_id,omitempty"` // UUID/txn id for locally-sent messages
}

func (m *MessageMetadata) CopyFrom(other any) {
	o, ok := other.(*MessageMetadata)
	if !ok || o == nil || m == nil {
		return
	}
	if o.EditCount > m.EditCount {
		m.EditCount = o.EditCount
	}
	if m.XChatClientMsgID == "" {
		m.XChatClientMsgID = o.XChatClientMsgID
	}
}

type PushKeys struct {
	P256DH  []byte `json:"p256dh"`
	Auth    []byte `json:"auth"`
	Private []byte `json:"private"`
}

func (m *UserLoginMetadata) GeneratePushKeys() {
	privateKey := exerrors.Must(ecdh.P256().GenerateKey(rand.Reader))
	m.PushKeys = &PushKeys{
		P256DH:  privateKey.Public().(*ecdh.PublicKey).Bytes(),
		Auth:    random.Bytes(16),
		Private: privateKey.Bytes(),
	}
}
