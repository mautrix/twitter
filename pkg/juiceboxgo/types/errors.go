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

package types

import (
	"errors"
	"fmt"
)

// Standard errors that can occur during Juicebox operations.
var (
	// ErrInvalidAuth indicates the realm rejected the auth token.
	ErrInvalidAuth = errors.New("realm rejected auth token")
	// ErrUpgradeRequired indicates the SDK version is too old.
	ErrUpgradeRequired = errors.New("SDK upgrade required")
	// ErrRateLimitExceeded indicates too many requests.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	// ErrAssertion indicates an unexpected protocol error.
	ErrAssertion = errors.New("assertion error")
	// ErrTransient indicates a temporary error that may succeed on retry.
	ErrTransient = errors.New("transient error")
	// ErrNotRegistered indicates the user has not registered a secret.
	ErrNotRegistered = errors.New("secret not registered")
)

// RecoverError provides details about a recovery failure.
type RecoverError struct {
	Reason           error
	GuessesRemaining *uint16
}

func (e *RecoverError) Error() string {
	if e.GuessesRemaining != nil {
		return fmt.Sprintf("%s (guesses remaining: %d)", e.Reason.Error(), *e.GuessesRemaining)
	}
	return e.Reason.Error()
}

func (e *RecoverError) Unwrap() error {
	return e.Reason
}

// ErrInvalidPin creates an error for an invalid PIN with the remaining guesses.
func ErrInvalidPin(guessesRemaining uint16) error {
	return &RecoverError{
		Reason:           errors.New("invalid PIN"),
		GuessesRemaining: &guessesRemaining,
	}
}
