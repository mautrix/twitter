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

// Re-exported errors from types package
var (
	ErrInvalidAuth       = types.ErrInvalidAuth
	ErrUpgradeRequired   = types.ErrUpgradeRequired
	ErrRateLimitExceeded = types.ErrRateLimitExceeded
	ErrAssertion         = types.ErrAssertion
	ErrTransient         = types.ErrTransient
	ErrNotRegistered     = types.ErrNotRegistered
)

// RecoverError type alias
type RecoverError = types.RecoverError

// ErrInvalidPin creates an error for invalid PIN with guesses remaining.
var ErrInvalidPin = types.ErrInvalidPin
