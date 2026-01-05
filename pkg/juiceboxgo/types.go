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

package juiceboxgo

import "go.mau.fi/mautrix-twitter/pkg/juiceboxgo/types"

// Type aliases for convenience - users can import just the main package.
type (
	RealmID                       = types.RealmID
	Pin                           = types.Pin
	Secret                        = types.Secret
	UserInfo                      = types.UserInfo
	RegistrationVersion           = types.RegistrationVersion
	Realm                         = types.Realm
	PinHashingMode                = types.PinHashingMode
	UnlockKey                     = types.UnlockKey
	UnlockKeyCommitment           = types.UnlockKeyCommitment
	UnlockKeyTag                  = types.UnlockKeyTag
	EncryptedUserSecret           = types.EncryptedUserSecret
	EncryptedUserSecretCommitment = types.EncryptedUserSecretCommitment
	SessionID                     = types.SessionID
	AuthToken                     = types.AuthToken
)

// Re-exported constants
const (
	PinHashingModeStandard2019 = types.PinHashingModeStandard2019
	PinHashingModeFastInsecure = types.PinHashingModeFastInsecure
	MaxUserSecretLength        = types.MaxUserSecretLength
)

// ParseRealmID parses a hex-encoded realm ID.
var ParseRealmID = types.ParseRealmID
