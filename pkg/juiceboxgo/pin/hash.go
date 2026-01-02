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

// Package pin implements PIN hashing for Juicebox.
package pin

import (
	"encoding/binary"

	"golang.org/x/crypto/argon2"
)

// HashingMode specifies the PIN hashing parameters.
type HashingMode string

const (
	// HashingModeStandard2019 uses secure parameters for production.
	HashingModeStandard2019 HashingMode = "Standard2019"
	// HashingModeFastInsecure uses fast parameters for testing.
	HashingModeFastInsecure HashingMode = "FastInsecure"
)

// HashResult contains the two keys derived from PIN hashing.
type HashResult struct {
	// AccessKey is used as input to the OPRF (first 32 bytes).
	AccessKey [32]byte
	// EncryptionKeySeed is used to derive the encryption key (last 32 bytes).
	EncryptionKeySeed [32]byte
}

// HashPIN hashes the PIN using Argon2id with the specified parameters.
// The salt is constructed as: BE4(len(version)) || version || BE4(len(userInfo)) || userInfo
func HashPIN(pinBytes []byte, mode HashingMode, version [16]byte, userInfo []byte) HashResult {
	var time, memory uint32
	var threads uint8

	switch mode {
	case HashingModeStandard2019:
		// Tuned for security on modern devices
		memory = 16 * 1024 // 16 MiB
		time = 32
		threads = 1
	case HashingModeFastInsecure:
		// Fast for testing
		memory = 8  // Minimum
		time = 1    // Minimum
		threads = 1 // Minimum
	default:
		// Default to standard
		memory = 16 * 1024
		time = 32
		threads = 1
	}

	// Construct salt: BE4(len(version)) || version || BE4(len(userInfo)) || userInfo
	salt := make([]byte, 0, 4+len(version)+4+len(userInfo))
	salt = appendBE4(salt, uint32(len(version)))
	salt = append(salt, version[:]...)
	salt = appendBE4(salt, uint32(len(userInfo)))
	salt = append(salt, userInfo...)

	// Hash using Argon2id, output 64 bytes
	hash := argon2.IDKey(pinBytes, salt, time, memory, threads, 64)

	var result HashResult
	copy(result.AccessKey[:], hash[:32])
	copy(result.EncryptionKeySeed[:], hash[32:])
	return result
}

// appendBE4 appends a uint32 in big-endian format.
func appendBE4(buf []byte, v uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], v)
	return append(buf, b[:]...)
}
