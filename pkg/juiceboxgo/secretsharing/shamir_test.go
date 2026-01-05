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

package secretsharing

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/bwesterb/go-ristretto"
)

// TestIndexToScalar verifies that index-to-scalar conversion matches Rust's Scalar::from(u32).
func TestIndexToScalar(t *testing.T) {
	testCases := []struct {
		index    Index
		expected string // hex of expected scalar bytes (little-endian)
	}{
		{1, "0100000000000000000000000000000000000000000000000000000000000000"},
		{2, "0200000000000000000000000000000000000000000000000000000000000000"},
		{3, "0300000000000000000000000000000000000000000000000000000000000000"},
		{255, "ff00000000000000000000000000000000000000000000000000000000000000"},
		{256, "0001000000000000000000000000000000000000000000000000000000000000"},
		{65536, "0000010000000000000000000000000000000000000000000000000000000000"},
	}

	for _, tc := range testCases {
		t.Run("index_"+string(rune(tc.index)), func(t *testing.T) {
			scalar := indexToScalar(tc.index)
			expected, _ := hex.DecodeString(tc.expected)

			if !bytes.Equal(scalar.Bytes(), expected) {
				t.Errorf("Index %d: got %s, want %s",
					tc.index,
					hex.EncodeToString(scalar.Bytes()),
					tc.expected)
			}
		})
	}
}

// TestLagrangeInterpolation tests that we can recover a secret from shares.
func TestLagrangeInterpolation(t *testing.T) {
	// Create a known secret
	var secret ristretto.Scalar
	var secretBytes [64]byte
	rand.Read(secretBytes[:])
	secret.SetReduced(&secretBytes)

	// Manually create shares for a degree-2 polynomial (threshold 3)
	// f(x) = a0 + a1*x + a2*x^2
	// where a0 = secret

	var a1, a2 ristretto.Scalar
	var a1Bytes, a2Bytes [64]byte
	rand.Read(a1Bytes[:])
	rand.Read(a2Bytes[:])
	a1.SetReduced(&a1Bytes)
	a2.SetReduced(&a2Bytes)

	// Evaluate f(x) at x = 1, 2, 3
	shares := make([]ScalarShare, 3)
	for i := 0; i < 3; i++ {
		x := indexToScalar(Index(i + 1))
		var x2 ristretto.Scalar
		x2.Mul(&x, &x) // x^2

		// f(x) = a0 + a1*x + a2*x^2
		var term1, term2, share ristretto.Scalar
		term1.Mul(&a1, &x)
		term2.Mul(&a2, &x2)
		share.Add(&secret, &term1)
		share.Add(&share, &term2)

		shares[i] = ScalarShare{
			Index:  Index(i + 1),
			Secret: share,
		}
	}

	// Recover secret
	recovered, err := RecoverScalar(shares)
	if err != nil {
		t.Fatalf("RecoverScalar failed: %v", err)
	}

	if !bytes.Equal(recovered.Bytes(), secret.Bytes()) {
		t.Errorf("Recovered secret mismatch:\ngot:  %s\nwant: %s",
			hex.EncodeToString(recovered.Bytes()),
			hex.EncodeToString(secret.Bytes()))
	} else {
		t.Logf("✓ Secret recovery works correctly")
	}
}

// TestPointLagrangeInterpolation tests that we can recover a point from shares.
func TestPointLagrangeInterpolation(t *testing.T) {
	// Create a known secret point
	var secretScalar ristretto.Scalar
	var secretBytes [64]byte
	rand.Read(secretBytes[:])
	secretScalar.SetReduced(&secretBytes)

	var secret ristretto.Point
	secret.ScalarMultBase(&secretScalar)

	// Create random polynomial coefficients (as points)
	var a1Scalar, a2Scalar ristretto.Scalar
	var a1Bytes, a2Bytes [64]byte
	rand.Read(a1Bytes[:])
	rand.Read(a2Bytes[:])
	a1Scalar.SetReduced(&a1Bytes)
	a2Scalar.SetReduced(&a2Bytes)

	var a1, a2 ristretto.Point
	a1.ScalarMultBase(&a1Scalar)
	a2.ScalarMultBase(&a2Scalar)

	// Evaluate f(x) = a0 + a1*x + a2*x^2 at x = 1, 2, 3
	// For points: f(x) = a0 + x*a1 + x^2*a2
	shares := make([]PointShare, 3)
	for i := 0; i < 3; i++ {
		x := indexToScalar(Index(i + 1))
		var x2 ristretto.Scalar
		x2.Mul(&x, &x) // x^2

		// f(x) = a0 + x*a1 + x^2*a2
		var term1, term2, share ristretto.Point
		term1.ScalarMult(&a1, &x)
		term2.ScalarMult(&a2, &x2)
		share.Add(&secret, &term1)
		share.Add(&share, &term2)

		shares[i] = PointShare{
			Index:  Index(i + 1),
			Secret: share,
		}
	}

	// Recover secret
	recovered, err := RecoverPoint(shares)
	if err != nil {
		t.Fatalf("RecoverPoint failed: %v", err)
	}

	if !bytes.Equal(recovered.Bytes(), secret.Bytes()) {
		t.Errorf("Recovered point mismatch:\ngot:  %s\nwant: %s",
			hex.EncodeToString(recovered.Bytes()),
			hex.EncodeToString(secret.Bytes()))
	} else {
		t.Logf("✓ Point recovery works correctly")
	}
}

// TestSubsetRecovery tests that any threshold subset can recover the secret.
func TestSubsetRecovery(t *testing.T) {
	// Create a known secret
	var secret ristretto.Scalar
	var secretBytes [64]byte
	rand.Read(secretBytes[:])
	secret.SetReduced(&secretBytes)

	// Create shares with threshold 2, count 5
	var a1 ristretto.Scalar
	var a1Bytes [64]byte
	rand.Read(a1Bytes[:])
	a1.SetReduced(&a1Bytes)

	// Create 5 shares of degree-1 polynomial
	allShares := make([]ScalarShare, 5)
	for i := 0; i < 5; i++ {
		x := indexToScalar(Index(i + 1))
		var term, share ristretto.Scalar
		term.Mul(&a1, &x)
		share.Add(&secret, &term)

		allShares[i] = ScalarShare{
			Index:  Index(i + 1),
			Secret: share,
		}
	}

	// Test recovery with different subsets of 2 shares
	subsets := [][2]int{
		{0, 1}, {0, 2}, {0, 3}, {0, 4},
		{1, 2}, {1, 3}, {1, 4},
		{2, 3}, {2, 4},
		{3, 4},
	}

	for _, subset := range subsets {
		shares := []ScalarShare{allShares[subset[0]], allShares[subset[1]]}
		recovered, err := RecoverScalar(shares)
		if err != nil {
			t.Errorf("Subset %v: RecoverScalar failed: %v", subset, err)
			continue
		}
		if !bytes.Equal(recovered.Bytes(), secret.Bytes()) {
			t.Errorf("Subset %v: recovered secret mismatch", subset)
		}
	}
	t.Logf("✓ All subsets recover the secret correctly")
}
