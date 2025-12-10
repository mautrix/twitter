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
		Portal:   nil,
		Message: func() any {
			return &MessageMetadata{}
		},
		Ghost: nil,
		UserLogin: func() any {
			return &UserLoginMetadata{}
		},
	}
}

type UserLoginMetadata struct {
	Cookies           string    `json:"cookies"`
	SecretKey         string    `json:"secret_key,omitempty"`
	SigningKey        string    `json:"signing_key,omitempty"`
	SigningKeyVersion string    `json:"signing_key_version,omitempty"`
	PushKeys          *PushKeys `json:"push_keys,omitempty"`

	Session            *twittermeow.CachedSession      `json:"session,omitempty"`
	ConversationKeys   map[string]*ConversationKeyData `json:"conversation_keys,omitempty"`
	ConversationTokens map[string]string               `json:"conversation_tokens,omitempty"`  // conversationID -> server-provided token
	MaxUserSequenceID  string                          `json:"max_user_sequence_id,omitempty"` // Last processed sequence ID for incremental inbox fetching
	// Deprecated: kept for backward compatibility with older saves; prefer MaxUserSequenceID.
	MaxSequenceID      string `json:"max_sequence_id,omitempty"`
	MessagePullVersion *int   `json:"message_pull_version,omitempty"`
}

// ConversationKeyData stores a conversation encryption key.
// The map key in UserLoginMetadata.ConversationKeys is "conversationID:keyVersion".
type ConversationKeyData struct {
	KeyVersion string     `json:"key_version"`
	Key        []byte     `json:"key"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type MessageMetadata struct {
	EditCount int `json:"edit_count,omitempty"`
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
