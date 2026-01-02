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

package oprf

import (
	"crypto/sha512"
	"crypto/subtle"
	"errors"

	"github.com/bwesterb/go-ristretto"
)

// Proof is a DLEQ proof that the server evaluated the OPRF correctly.
type Proof struct {
	C      ristretto.Scalar // Challenge
	BetaZ  ristretto.Scalar // Response
}

// ErrInvalidProof is returned when DLEQ proof verification fails.
var ErrInvalidProof = errors.New("invalid DLEQ proof")

// VerifyProof verifies that the server's OPRF evaluation is correct.
// This ensures: log_G(publicKey) == log_blindedInput(blindedOutput)
func VerifyProof(blindedInput, blindedOutput *BlindedInput, publicKey *PublicKey, proof *Proof) error {
	// Compute v_t = G * beta_z - publicKey * c
	var vT ristretto.Point
	var gBetaZ, vC ristretto.Point
	gBetaZ.ScalarMultBase(&proof.BetaZ)
	vC.ScalarMult(&publicKey.point, &proof.C)
	vT.Sub(&gBetaZ, &vC)

	// Compute w_t = blindedInput * beta_z - blindedOutput * c
	var wT ristretto.Point
	var uBetaZ, wC ristretto.Point
	uBetaZ.ScalarMult(&blindedInput.point, &proof.BetaZ)
	wC.ScalarMult(blindedOutput.Point(), &proof.C)
	wT.Sub(&uBetaZ, &wC)

	// Recompute challenge
	c := hashToChallenge(
		blindedInput.Bytes(),
		publicKey.Bytes(),
		blindedOutput.Bytes(),
		vT.Bytes(),
		wT.Bytes(),
	)

	// Verify c == proof.C (constant time)
	if subtle.ConstantTimeCompare(c.Bytes(), proof.C.Bytes()) != 1 {
		return ErrInvalidProof
	}

	return nil
}

// hashToChallenge computes the Fiat-Shamir challenge for DLEQ.
func hashToChallenge(u, v, w, vT, wT []byte) ristretto.Scalar {
	h := sha512.New()
	h.Write([]byte("Juicebox_DLEQ_2023_1;"))
	h.Write(u)
	h.Write(v)
	h.Write(w)
	h.Write(vT)
	h.Write(wT)

	var hashBytes [64]byte
	copy(hashBytes[:], h.Sum(nil))

	var c ristretto.Scalar
	c.SetReduced(&hashBytes)
	return c
}
