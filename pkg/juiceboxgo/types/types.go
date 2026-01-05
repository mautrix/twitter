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

// Package types defines shared types for the Juicebox SDK.
package types

import (
	"encoding/hex"
	"fmt"
)

// RealmID is a unique 16-byte identifier for a realm.
type RealmID [16]byte

func (r RealmID) String() string {
	return hex.EncodeToString(r[:])
}

// ParseRealmID parses a hex-encoded realm ID.
func ParseRealmID(s string) (RealmID, error) {
	var id RealmID
	b, err := hex.DecodeString(s)
	if err != nil {
		return id, fmt.Errorf("invalid realm ID hex: %w", err)
	}
	if len(b) != 16 {
		return id, fmt.Errorf("realm ID must be 16 bytes, got %d", len(b))
	}
	copy(id[:], b)
	return id, nil
}

// Pin is a user-chosen password (max 128 bytes).
type Pin []byte

// Secret is the recovered secret data (max 128 bytes).
type Secret []byte

// UserInfo is additional data added to the salt for PIN hashing.
type UserInfo []byte

// RegistrationVersion is a 16-byte version identifier.
type RegistrationVersion [16]byte

// Realm represents a remote service that the client interacts with.
type Realm struct {
	// ID is a unique identifier specified by the realm.
	ID RealmID `json:"id"`
	// Address is the network address to connect to the service.
	Address string `json:"address"`
	// PublicKey is a long-lived public key for hardware-backed realms (optional).
	PublicKey []byte `json:"public_key,omitempty"`
}

// PinHashingMode specifies the strategy for hashing the user's PIN.
type PinHashingMode string

const (
	// PinHashingModeStandard2019 is a tuned hash secure for use on modern devices.
	PinHashingModeStandard2019 PinHashingMode = "Standard2019"
	// PinHashingModeFastInsecure is a fast hash for testing only.
	PinHashingModeFastInsecure PinHashingMode = "FastInsecure"
)

// UnlockKey is a 32-byte key derived from the OPRF output.
type UnlockKey [32]byte

// UnlockKeyCommitment is a 32-byte commitment to the unlock key.
type UnlockKeyCommitment [32]byte

// UnlockKeyTag is a 16-byte tag derived from unlock key and realm ID.
type UnlockKeyTag [16]byte

// EncryptedUserSecret is the encrypted secret stored on realms (145 bytes).
type EncryptedUserSecret [145]byte

// EncryptedUserSecretCommitment is a 16-byte commitment.
type EncryptedUserSecretCommitment [16]byte

// SessionID is a 32-bit session identifier.
type SessionID uint32

// AuthToken is a bearer token for realm authentication.
type AuthToken string

// MaxUserSecretLength is the maximum allowed bytes for a user secret.
const MaxUserSecretLength = 128
