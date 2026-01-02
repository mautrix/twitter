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

// Package oprf implements the Oblivious Pseudorandom Function used by Juicebox.
package oprf

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"

	"github.com/bwesterb/go-ristretto"
)

// ErrInvalidPoint is returned when a point cannot be decoded.
var ErrInvalidPoint = errors.New("invalid ristretto point")

// BlindingFactor is a random scalar used to blind the OPRF input.
type BlindingFactor struct {
	scalar ristretto.Scalar
}

// BlindedInput is the blinded OPRF input sent to the server.
type BlindedInput struct {
	point ristretto.Point
}

// Bytes returns the compressed point bytes.
func (b *BlindedInput) Bytes() []byte {
	return b.point.Bytes()
}

// SetBytes sets the blinded input from compressed point bytes.
func (b *BlindedInput) SetBytes(data []byte) error {
	if len(data) != 32 {
		return ErrInvalidPoint
	}
	var buf [32]byte
	copy(buf[:], data)
	if !b.point.SetBytes(&buf) {
		return ErrInvalidPoint
	}
	return nil
}

// Point returns the underlying point.
func (b *BlindedInput) Point() *ristretto.Point {
	return &b.point
}

// BlindedOutput is the server's blinded OPRF result.
type BlindedOutput struct {
	point ristretto.Point
}

// Bytes returns the compressed point bytes.
func (b *BlindedOutput) Bytes() []byte {
	return b.point.Bytes()
}

// SetBytes sets the blinded output from compressed point bytes.
func (b *BlindedOutput) SetBytes(data []byte) error {
	if len(data) != 32 {
		return ErrInvalidPoint
	}
	var buf [32]byte
	copy(buf[:], data)
	if !b.point.SetBytes(&buf) {
		return ErrInvalidPoint
	}
	return nil
}

// Point returns the underlying point.
func (b *BlindedOutput) Point() *ristretto.Point {
	return &b.point
}

// Output is the final OPRF result after finalization (64 bytes).
type Output [64]byte

// PublicKey is the OPRF server's public key.
type PublicKey struct {
	point ristretto.Point
}

// Bytes returns the compressed point bytes.
func (p *PublicKey) Bytes() []byte {
	return p.point.Bytes()
}

// SetBytes sets the public key from compressed point bytes.
func (p *PublicKey) SetBytes(data []byte) error {
	if len(data) != 32 {
		return ErrInvalidPoint
	}
	var buf [32]byte
	copy(buf[:], data)
	if !p.point.SetBytes(&buf) {
		return ErrInvalidPoint
	}
	return nil
}

// Point returns the underlying point.
func (p *PublicKey) Point() *ristretto.Point {
	return &p.point
}

// Start begins the OPRF protocol on the client.
// Returns the blinding factor (keep secret) and blinded input (send to server).
func Start(input []byte) (*BlindingFactor, *BlindedInput, error) {
	// Hash input to a point using DeriveDalek for compatibility with Rust's hash_from_bytes
	var inputPoint ristretto.Point
	inputPoint.DeriveDalek(input)

	// Generate random blinding factor
	var blindingFactor ristretto.Scalar
	var randomBytes [64]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		return nil, nil, err
	}
	blindingFactor.SetReduced(&randomBytes)

	// Compute blinded input: inputPoint * blindingFactor
	var blindedPoint ristretto.Point
	blindedPoint.ScalarMult(&inputPoint, &blindingFactor)

	return &BlindingFactor{scalar: blindingFactor},
		&BlindedInput{point: blindedPoint},
		nil
}

// Finalize completes the OPRF protocol on the client.
// Returns the final OPRF output.
func Finalize(input []byte, blindingFactor *BlindingFactor, blindedOutput *BlindedOutput) Output {
	// Compute 1/blindingFactor
	var invBlindingFactor ristretto.Scalar
	invBlindingFactor.Inverse(&blindingFactor.scalar)

	// Compute result: blindedOutput * (1/blindingFactor)
	var result ristretto.Point
	result.ScalarMult(&blindedOutput.point, &invBlindingFactor)

	// Hash to output: SHA512("Juicebox_OPRF_2023_1;" || input || compress(result))
	return hashToOutput(input, &result)
}

func hashToOutput(input []byte, result *ristretto.Point) Output {
	h := sha512.New()
	h.Write([]byte("Juicebox_OPRF_2023_1;"))
	h.Write(input)
	h.Write(result.Bytes())

	var output Output
	copy(output[:], h.Sum(nil))
	return output
}
