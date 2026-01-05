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

package pin

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// Test vector from Rust SDK: rust/sdk/src/pin.rs:119-141
func TestHashPIN_RustVector(t *testing.T) {
	// Input
	pinBytes := []byte("1234")
	var version [16]byte
	for i := range version {
		version[i] = 0x05
	}
	userInfo := []byte("artemis")

	// Expected output from Rust
	expectedAccessKey := [32]byte{
		41, 53, 218, 132, 201, 116, 35, 179, 127, 52, 87, 35, 27, 135, 124, 230,
		172, 32, 147, 80, 29, 114, 85, 219, 238, 29, 235, 9, 165, 216, 130, 27,
	}
	expectedEncryptionKeySeed := [32]byte{
		135, 200, 201, 181, 211, 234, 159, 234, 131, 182, 172, 106, 100, 226, 91, 151,
		196, 114, 44, 164, 228, 11, 234, 37, 35, 239, 234, 38, 33, 37, 226, 42,
	}

	// Run hash
	result := HashPIN(pinBytes, HashingModeStandard2019, version, userInfo)

	// Verify
	if !bytes.Equal(result.AccessKey[:], expectedAccessKey[:]) {
		t.Errorf("AccessKey mismatch:\ngot:  %s\nwant: %s",
			hex.EncodeToString(result.AccessKey[:]),
			hex.EncodeToString(expectedAccessKey[:]))
	}
	if !bytes.Equal(result.EncryptionKeySeed[:], expectedEncryptionKeySeed[:]) {
		t.Errorf("EncryptionKeySeed mismatch:\ngot:  %s\nwant: %s",
			hex.EncodeToString(result.EncryptionKeySeed[:]),
			hex.EncodeToString(expectedEncryptionKeySeed[:]))
	}

	t.Logf("PIN hash test passed")
	t.Logf("AccessKey: %s", hex.EncodeToString(result.AccessKey[:]))
	t.Logf("EncryptionKeySeed: %s", hex.EncodeToString(result.EncryptionKeySeed[:]))
}
