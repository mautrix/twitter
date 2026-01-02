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

// Package crypto implements cryptographic operations for Juicebox.
package crypto

import (
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"hash"

	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20poly1305"

	"github.com/bwesterb/go-ristretto"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/oprf"
)

// ErrDecryptionFailed is returned when secret decryption fails.
var ErrDecryptionFailed = errors.New("decryption failed")

// MaxUserSecretLength is the maximum length of user secrets.
const MaxUserSecretLength = 128

// DeriveUnlockKeyAndCommitment derives the unlock key and commitment from OPRF output.
// Returns (unlockKey [32]byte, unlockKeyCommitment [32]byte).
func DeriveUnlockKeyAndCommitment(oprfOutput oprf.Output) ([32]byte, [32]byte) {
	digest := sha512.Sum512(oprfOutput[:])
	var commitment [32]byte
	var key [32]byte
	copy(commitment[:], digest[:32])
	copy(key[:], digest[32:])
	return key, commitment
}

// DeriveUnlockKeyTag derives the unlock key tag for a specific realm.
// Uses Blake2s-128 MAC: MAC(unlock_key, "Unlock Key Tag" || realm_id)
func DeriveUnlockKeyTag(unlockKey [32]byte, realmID [16]byte) [16]byte {
	label := []byte("Unlock Key Tag")

	mac, _ := blake2s.New128(unlockKey[:])
	writeLengthPrefixed(mac, label)
	writeLengthPrefixed(mac, realmID[:])

	var tag [16]byte
	copy(tag[:], mac.Sum(nil))
	return tag
}

// DeriveEncryptionKey derives the encryption key from the seed and scalar.
// Uses Blake2s-256 MAC: MAC(seed, "User Secret Encryption Key" || scalar)
func DeriveEncryptionKey(seed [32]byte, scalar *ristretto.Scalar) [32]byte {
	label := []byte("User Secret Encryption Key")
	scalarBytes := scalar.Bytes()

	mac, _ := blake2s.New256(seed[:])
	writeLengthPrefixed(mac, label)
	writeLengthPrefixed(mac, scalarBytes[:])

	var key [32]byte
	copy(key[:], mac.Sum(nil))
	return key
}

// DecryptSecret decrypts the user secret using ChaCha20-Poly1305.
// The encrypted secret uses a fixed all-zeros nonce since the key is unique per encryption.
func DecryptSecret(encryptedSecret []byte, encryptionKey [32]byte) ([]byte, error) {
	cipher, err := chacha20poly1305.New(encryptionKey[:])
	if err != nil {
		return nil, err
	}

	// Fixed nonce: all zeros (safe because key is unique per encryption)
	nonce := make([]byte, chacha20poly1305.NonceSize)

	// Decrypt
	padded, err := cipher.Open(nil, nonce, encryptedSecret, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	// Remove padding: first byte is length, followed by data, then zeros
	if len(padded) < 1 {
		return nil, ErrDecryptionFailed
	}
	length := int(padded[0])
	if length > len(padded)-1 || length > MaxUserSecretLength {
		return nil, ErrDecryptionFailed
	}

	return padded[1 : 1+length], nil
}

// DeriveEncryptedUserSecretCommitment derives the commitment for verification.
// Uses Blake2s-128 MAC.
func DeriveEncryptedUserSecretCommitment(
	unlockKey [32]byte,
	realmID [16]byte,
	scalarShare *ristretto.Scalar,
	encryptedSecret []byte,
) [16]byte {
	label := []byte("Encrypted User Secret Commitment")
	scalarBytes := scalarShare.Bytes()

	mac, _ := blake2s.New128(unlockKey[:])
	writeLengthPrefixed(mac, label)
	writeLengthPrefixed(mac, realmID[:])
	writeLengthPrefixed(mac, scalarBytes[:])
	writeLengthPrefixed(mac, encryptedSecret)

	var commitment [16]byte
	copy(commitment[:], mac.Sum(nil))
	return commitment
}

// writeLengthPrefixed writes BE4(len) || data to the hash.
func writeLengthPrefixed(h hash.Hash, data []byte) {
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))
	h.Write(lenBuf[:])
	h.Write(data)
}
