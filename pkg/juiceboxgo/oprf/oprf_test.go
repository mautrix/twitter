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
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/bwesterb/go-ristretto"
)

// Test vectors from Rust SDK: rust/oprf/src/test_vectors.json

type testVector struct {
	Name               string
	Input              string // hex
	PrivateKeySeed     string // hex (64 bytes)
	BlindingFactorSeed string // hex (64 bytes)

	PrivateKey     string // hex (32 bytes)
	PublicKey      string // hex (32 bytes)
	BlindingFactor string // hex (32 bytes)
	BlindedInput   string // hex (32 bytes)
	BlindedOutput  string // hex (32 bytes)
	ProofC         string // hex (32 bytes)
	ProofBetaZ     string // hex (32 bytes)
	Output         string // hex (64 bytes)
}

var testVectors = []testVector{
	{
		Name:               "random-test-01",
		Input:              "840841056a7b40f363df5a25827a",
		PrivateKeySeed:     "7319eafd6f8a54a59bb9d450a4cbfc04a6f0bb42797d58fc808e8514b31423e0c81b0f143403a175210497ae2090596b34a6476d602317f5c1ecda80f82050f0",
		BlindingFactorSeed: "0f0b48a38ad125109d3888c6fc57232a5afc0d7bbf7013c7afd170385d4e47631522b62cc117592337758d86e6a1eaf89eae9afdb16c0027dbb8830d50d88f77",
		PrivateKey:         "cca1a0304b113ec01cafa2545c0428497fd65a4924b4697033f5c19aaec2ac0a",
		PublicKey:          "9e2bc4e246e540092324937ed33fd01caf0297137e35345c32ecf49e87e35056",
		BlindingFactor:     "36fce774182f86c3d373bf632e7056e92371a49fbe27b01ef9d2688605446703",
		BlindedInput:       "a8767323a469385742eb85b73a3d51372f4e15d336f72567eb12d3410fa6815c",
		BlindedOutput:      "5020caea110ad8333a059d848b33f5607c56eca6c93df929999df5362ed00c63",
		ProofC:             "d1b6200350187462bc499b9761b960560d300f9b03873a3ba8bb31b68e13150e",
		ProofBetaZ:         "a3be5bc291c6cc07c86de7b6e0d12a49475663bd4b4805f7a535f96900fa7802",
		Output:             "0c12ee3f037370f7c9cf203d144f3653c96174bfbfd63b6c07d67e9a99e4e59953fbfa90232dff216ea2290e915fb8f2818be91f81902994aa50cf9aa77ed843",
	},
	{
		Name:               "random-test-02",
		Input:              "4031",
		PrivateKeySeed:     "81a5851536ddf8ab2e406cdfa4212239bb9fa70b3371fc210d34989c9021b73502e7ba25204359bda3a7acb1b5d078c9b3591313e30198301d26d2293520c47c",
		BlindingFactorSeed: "f162c59f57c748bffa0d7e60eab204fb43c667b144219e5164602ef34c190d8ac91bb3337601ec3391a1b04583d768909c716035e955ba5012506e271279a6f8",
		PrivateKey:         "b93aae4cdc5d28f51aedd681218f1dab4ff3a4c1836234a5e291af12794e3c0b",
		PublicKey:          "70fd2bd8bc7a9ee328d366d0a97d696a30b1a369663afa43f45de26280300f49",
		BlindingFactor:     "33cafe98e426b123375558e5f8da85bb89be38af5f7c4461f790328459e71305",
		BlindedInput:       "2a9e7fc1fd010bde24673255559f21bafcf83221251947ef9f497e15b3190c0b",
		BlindedOutput:      "26330cde20259f5f134e6d3291a44e23d1762be006e520008aa173aa5e0e6826",
		ProofC:             "ff6fbb60b5d9ea65c34febdcd73a25181f7e25ec818b8cef101623eaa028bc02",
		ProofBetaZ:         "d119a7c682324f6fd4e79aaba1c32e5dd494626d5d115fd2a319bccd474d7f04",
		Output:             "bf283a626e65927026c4cfe9e55250b252fc8a0883cd307cc463d498991372584b8b689697d36366a22f0fe10be92309e1a17bea37360205d0448cd26a71f3ac",
	},
	{
		Name:               "random-test-03",
		Input:              "c5523613f851591b30b74d",
		PrivateKeySeed:     "5c1ed0f935d10d8927262ec850750c9dd2cd68e3dd8e20f39ae19e7fdb226a786528e24b8f8ee96846084e72faa161b0172d5495eca19f487d252e2b7525ecf6",
		BlindingFactorSeed: "3e7308496c6944a9bdf2e4efeb0d068e8a2033bc5072d0b526aa3bafe86a311553c1a4a5266bb0793f57e13e873158e2109130c36121f6e8e49eb4cd37d42f47",
		PrivateKey:         "9ffd1a9832e13dd942bdd75633453421bcaa8540d377fb37b3a52f921b43010a",
		PublicKey:          "a837f11591ba58f358f3fa02124454b3b0bdf3f910e92d1796f2e80a1c33d85d",
		BlindingFactor:     "d3dde93aca36bd0188b4c2335cbfc95b4aefdacd7fc491e89eeb023dfd21d50e",
		BlindedInput:       "ea1957d1930a99b88b6c0de809e0d7a143226a70f9909865320e177e77bb5121",
		BlindedOutput:      "f499b4eb255f081866c11e52cb2d9de2ed44f97474bdaf8f67af05896de28d79",
		ProofC:             "0f6c67d6c153751ab07830128d5da2e2181c50db8ce995d423e52c99b42a9707",
		ProofBetaZ:         "b6aa330aebbc812de855f0d7433902a32b46ffc46b60b2e11e2d0f206ecc4407",
		Output:             "eac0d507b41bc5d0b91d60fc6981c30d03b563d5dbb101997f29e7100d91b68431fabe192a9dde0264ca11c3e436f79e472d6725dfab78d97ab41b18c7ae64b1",
	},
}

// TestHashToPoint verifies that DeriveDalek produces the same hash-to-point as Rust's hash_from_bytes<Sha512>.
// This is CRITICAL - if this doesn't match, the entire OPRF will produce wrong results.
func TestHashToPoint(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/hash_to_point", func(t *testing.T) {
			input, err := hex.DecodeString(tv.Input)
			if err != nil {
				t.Fatalf("failed to decode input: %v", err)
			}

			// Hash input to point using DeriveDalek
			var inputPoint ristretto.Point
			inputPoint.DeriveDalek(input)

			t.Logf("Input: %s", tv.Input)
			t.Logf("Hash-to-point result: %s", hex.EncodeToString(inputPoint.Bytes()))

			// We can verify the hash-to-point indirectly by checking the blinded_input
			// blinded_input = inputPoint * blinding_factor
			// So: inputPoint = blinded_input * (1/blinding_factor)
			blindingFactorBytes, _ := hex.DecodeString(tv.BlindingFactor)
			blindedInputBytes, _ := hex.DecodeString(tv.BlindedInput)

			var blindingFactor ristretto.Scalar
			var buf32 [32]byte
			copy(buf32[:], blindingFactorBytes)
			blindingFactor.SetBytes(&buf32)

			var blindedInputPoint ristretto.Point
			copy(buf32[:], blindedInputBytes)
			if !blindedInputPoint.SetBytes(&buf32) {
				t.Fatal("failed to decode blinded_input point")
			}

			// Compute inverse of blinding factor
			var invBlindingFactor ristretto.Scalar
			invBlindingFactor.Inverse(&blindingFactor)

			// Recover input point: blindedInput * (1/blindingFactor)
			var recoveredInputPoint ristretto.Point
			recoveredInputPoint.ScalarMult(&blindedInputPoint, &invBlindingFactor)

			t.Logf("Recovered input point: %s", hex.EncodeToString(recoveredInputPoint.Bytes()))

			if !bytes.Equal(inputPoint.Bytes(), recoveredInputPoint.Bytes()) {
				t.Errorf("Hash-to-point MISMATCH:\nGo DeriveDalek:  %s\nRust recovered:  %s",
					hex.EncodeToString(inputPoint.Bytes()),
					hex.EncodeToString(recoveredInputPoint.Bytes()))
			} else {
				t.Logf("✓ Hash-to-point matches Rust")
			}
		})
	}
}

// TestBlindingFactor verifies that blinding factor derivation matches Rust.
func TestBlindingFactor(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/blinding_factor", func(t *testing.T) {
			seed, _ := hex.DecodeString(tv.BlindingFactorSeed)
			expected, _ := hex.DecodeString(tv.BlindingFactor)

			// Rust uses Scalar::from_bytes_mod_order_wide on 64-byte seed
			var seedBytes [64]byte
			copy(seedBytes[:], seed)

			var scalar ristretto.Scalar
			scalar.SetReduced(&seedBytes)

			if !bytes.Equal(scalar.Bytes(), expected) {
				t.Errorf("Blinding factor mismatch:\nGo:   %s\nRust: %s",
					hex.EncodeToString(scalar.Bytes()),
					hex.EncodeToString(expected))
			} else {
				t.Logf("✓ Blinding factor matches Rust: %s", hex.EncodeToString(scalar.Bytes()))
			}
		})
	}
}

// TestBlindedInput verifies the full blinding operation.
func TestBlindedInput(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/blinded_input", func(t *testing.T) {
			input, _ := hex.DecodeString(tv.Input)
			blindingFactorBytes, _ := hex.DecodeString(tv.BlindingFactor)
			expectedBlindedInput, _ := hex.DecodeString(tv.BlindedInput)

			// Hash input to point
			var inputPoint ristretto.Point
			inputPoint.DeriveDalek(input)

			// Set blinding factor
			var blindingFactor ristretto.Scalar
			var buf32 [32]byte
			copy(buf32[:], blindingFactorBytes)
			blindingFactor.SetBytes(&buf32)

			// Compute blinded input
			var blindedInput ristretto.Point
			blindedInput.ScalarMult(&inputPoint, &blindingFactor)

			if !bytes.Equal(blindedInput.Bytes(), expectedBlindedInput) {
				t.Errorf("Blinded input mismatch:\nGo:   %s\nRust: %s",
					hex.EncodeToString(blindedInput.Bytes()),
					hex.EncodeToString(expectedBlindedInput))
			} else {
				t.Logf("✓ Blinded input matches Rust: %s", hex.EncodeToString(blindedInput.Bytes()))
			}
		})
	}
}

// TestPrivateKey verifies private key derivation.
func TestPrivateKey(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/private_key", func(t *testing.T) {
			seed, _ := hex.DecodeString(tv.PrivateKeySeed)
			expected, _ := hex.DecodeString(tv.PrivateKey)

			var seedBytes [64]byte
			copy(seedBytes[:], seed)

			var scalar ristretto.Scalar
			scalar.SetReduced(&seedBytes)

			if !bytes.Equal(scalar.Bytes(), expected) {
				t.Errorf("Private key mismatch:\nGo:   %s\nRust: %s",
					hex.EncodeToString(scalar.Bytes()),
					hex.EncodeToString(expected))
			} else {
				t.Logf("✓ Private key matches Rust: %s", hex.EncodeToString(scalar.Bytes()))
			}
		})
	}
}

// TestPublicKey verifies public key derivation (G * private_key).
func TestPublicKey(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/public_key", func(t *testing.T) {
			privateKeyBytes, _ := hex.DecodeString(tv.PrivateKey)
			expected, _ := hex.DecodeString(tv.PublicKey)

			var privateKey ristretto.Scalar
			var buf32 [32]byte
			copy(buf32[:], privateKeyBytes)
			privateKey.SetBytes(&buf32)

			var publicKey ristretto.Point
			publicKey.ScalarMultBase(&privateKey)

			if !bytes.Equal(publicKey.Bytes(), expected) {
				t.Errorf("Public key mismatch:\nGo:   %s\nRust: %s",
					hex.EncodeToString(publicKey.Bytes()),
					hex.EncodeToString(expected))
			} else {
				t.Logf("✓ Public key matches Rust: %s", hex.EncodeToString(publicKey.Bytes()))
			}
		})
	}
}

// TestBlindedOutput verifies server-side OPRF evaluation.
func TestBlindedOutput(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/blinded_output", func(t *testing.T) {
			blindedInputBytes, _ := hex.DecodeString(tv.BlindedInput)
			privateKeyBytes, _ := hex.DecodeString(tv.PrivateKey)
			expected, _ := hex.DecodeString(tv.BlindedOutput)

			var blindedInput ristretto.Point
			var buf32 [32]byte
			copy(buf32[:], blindedInputBytes)
			if !blindedInput.SetBytes(&buf32) {
				t.Fatal("failed to decode blinded input")
			}

			var privateKey ristretto.Scalar
			copy(buf32[:], privateKeyBytes)
			privateKey.SetBytes(&buf32)

			// Server computes: blinded_output = blinded_input * private_key
			var blindedOutput ristretto.Point
			blindedOutput.ScalarMult(&blindedInput, &privateKey)

			if !bytes.Equal(blindedOutput.Bytes(), expected) {
				t.Errorf("Blinded output mismatch:\nGo:   %s\nRust: %s",
					hex.EncodeToString(blindedOutput.Bytes()),
					hex.EncodeToString(expected))
			} else {
				t.Logf("✓ Blinded output matches Rust: %s", hex.EncodeToString(blindedOutput.Bytes()))
			}
		})
	}
}

// TestFullOPRF tests the complete OPRF flow: Start -> (server eval) -> Finalize.
func TestFullOPRF(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/full_oprf", func(t *testing.T) {
			input, _ := hex.DecodeString(tv.Input)
			blindingFactorBytes, _ := hex.DecodeString(tv.BlindingFactor)
			blindedOutputBytes, _ := hex.DecodeString(tv.BlindedOutput)
			expectedOutput, _ := hex.DecodeString(tv.Output)

			// Set up blinding factor (normally random, but we use test vector)
			var blindingFactor ristretto.Scalar
			var buf32 [32]byte
			copy(buf32[:], blindingFactorBytes)
			blindingFactor.SetBytes(&buf32)

			// Set up blinded output from server
			var blindedOutputPoint BlindedOutput
			copy(buf32[:], blindedOutputBytes)
			blindedOutputPoint.point.SetBytes(&buf32)

			// Finalize OPRF
			bf := &BlindingFactor{scalar: blindingFactor}
			output := Finalize(input, bf, &blindedOutputPoint)

			if !bytes.Equal(output[:], expectedOutput) {
				t.Errorf("OPRF output mismatch:\nGo:   %s\nRust: %s",
					hex.EncodeToString(output[:]),
					hex.EncodeToString(expectedOutput))
			} else {
				t.Logf("✓ OPRF output matches Rust: %s", hex.EncodeToString(output[:]))
			}
		})
	}
}

// TestDLEQProofVerification tests DLEQ proof verification.
func TestDLEQProofVerification(t *testing.T) {
	for _, tv := range testVectors {
		t.Run(tv.Name+"/dleq_proof", func(t *testing.T) {
			blindedInputBytes, _ := hex.DecodeString(tv.BlindedInput)
			blindedOutputBytes, _ := hex.DecodeString(tv.BlindedOutput)
			publicKeyBytes, _ := hex.DecodeString(tv.PublicKey)
			proofCBytes, _ := hex.DecodeString(tv.ProofC)
			proofBetaZBytes, _ := hex.DecodeString(tv.ProofBetaZ)

			var buf32 [32]byte

			var blindedInput BlindedInput
			copy(buf32[:], blindedInputBytes)
			blindedInput.point.SetBytes(&buf32)

			var blindedOutput BlindedInput // Using BlindedInput for VerifyProof signature
			copy(buf32[:], blindedOutputBytes)
			blindedOutput.point.SetBytes(&buf32)

			var publicKey PublicKey
			copy(buf32[:], publicKeyBytes)
			publicKey.point.SetBytes(&buf32)

			var proofC, proofBetaZ ristretto.Scalar
			copy(buf32[:], proofCBytes)
			proofC.SetBytes(&buf32)
			copy(buf32[:], proofBetaZBytes)
			proofBetaZ.SetBytes(&buf32)

			proof := &Proof{C: proofC, BetaZ: proofBetaZ}

			err := VerifyProof(&blindedInput, &blindedOutput, &publicKey, proof)
			if err != nil {
				t.Errorf("DLEQ proof verification failed: %v", err)
			} else {
				t.Logf("✓ DLEQ proof verification passed")
			}
		})
	}
}
