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

// Package secretsharing implements Shamir's Secret Sharing with Lagrange interpolation.
package secretsharing

import (
	"errors"

	"github.com/bwesterb/go-ristretto"
)

// Index is a 1-based share index (must be non-zero).
type Index uint32

// ErrDuplicateShares is returned when duplicate share indices are provided.
var ErrDuplicateShares = errors.New("duplicate share indices")

// ScalarShare is a share of a secret scalar.
type ScalarShare struct {
	Index  Index
	Secret ristretto.Scalar
}

// PointShare is a share of a secret point.
type PointShare struct {
	Index  Index
	Secret ristretto.Point
}

// RecoverScalar recovers a secret scalar from shares using Lagrange interpolation.
func RecoverScalar(shares []ScalarShare) (ristretto.Scalar, error) {
	var result ristretto.Scalar
	result.SetZero()

	for i, share := range shares {
		// Check for duplicates
		for j := 0; j < i; j++ {
			if shares[j].Index == share.Index {
				return result, ErrDuplicateShares
			}
		}

		// Compute Lagrange coefficient
		lagrange := computeLagrangeCoefficient(shares, i)

		// Add share.Secret * lagrange to result
		var term ristretto.Scalar
		term.Mul(&share.Secret, &lagrange)
		result.Add(&result, &term)
	}

	return result, nil
}

// RecoverPoint recovers a secret point from shares using Lagrange interpolation.
func RecoverPoint(shares []PointShare) (ristretto.Point, error) {
	var result ristretto.Point
	result.SetZero()

	for i, share := range shares {
		// Check for duplicates
		for j := 0; j < i; j++ {
			if shares[j].Index == share.Index {
				return result, ErrDuplicateShares
			}
		}

		// Compute Lagrange coefficient
		lagrange := computeLagrangeCoefficientPoint(shares, i)

		// Add share.Secret * lagrange to result
		var term ristretto.Point
		term.ScalarMult(&share.Secret, &lagrange)
		result.Add(&result, &term)
	}

	return result, nil
}

// computeLagrangeCoefficient computes the Lagrange coefficient for share at index i.
func computeLagrangeCoefficient(shares []ScalarShare, i int) ristretto.Scalar {
	var numerator, denominator ristretto.Scalar
	numerator.SetOne()
	denominator.SetOne()

	xi := indexToScalar(shares[i].Index)

	for j, share := range shares {
		if i == j {
			continue
		}

		xj := indexToScalar(share.Index)

		// numerator *= xj
		numerator.Mul(&numerator, &xj)

		// denominator *= (xj - xi)
		var diff ristretto.Scalar
		diff.Sub(&xj, &xi)
		denominator.Mul(&denominator, &diff)
	}

	// Return numerator / denominator
	var result ristretto.Scalar
	var invDenom ristretto.Scalar
	invDenom.Inverse(&denominator)
	result.Mul(&numerator, &invDenom)
	return result
}

// computeLagrangeCoefficientPoint is the same as computeLagrangeCoefficient but for PointShare.
func computeLagrangeCoefficientPoint(shares []PointShare, i int) ristretto.Scalar {
	var numerator, denominator ristretto.Scalar
	numerator.SetOne()
	denominator.SetOne()

	xi := indexToScalar(shares[i].Index)

	for j, share := range shares {
		if i == j {
			continue
		}

		xj := indexToScalar(share.Index)

		// numerator *= xj
		numerator.Mul(&numerator, &xj)

		// denominator *= (xj - xi)
		var diff ristretto.Scalar
		diff.Sub(&xj, &xi)
		denominator.Mul(&denominator, &diff)
	}

	// Return numerator / denominator
	var result ristretto.Scalar
	var invDenom ristretto.Scalar
	invDenom.Inverse(&denominator)
	result.Mul(&numerator, &invDenom)
	return result
}

// indexToScalar converts a share index to a scalar.
func indexToScalar(index Index) ristretto.Scalar {
	var s ristretto.Scalar
	// Index is a uint32, convert to bytes in little-endian
	var bytes [32]byte
	bytes[0] = byte(index)
	bytes[1] = byte(index >> 8)
	bytes[2] = byte(index >> 16)
	bytes[3] = byte(index >> 24)
	s.SetBytes(&bytes)
	return s
}
