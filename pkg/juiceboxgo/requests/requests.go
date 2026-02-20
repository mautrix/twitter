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

// Package requests implements Juicebox realm request/response types.
package requests

import (
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/noise"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/types"
)

// ClientRequestKind identifies the type of request.
type ClientRequestKind string

const (
	ClientRequestKindHandshakeOnly  ClientRequestKind = "HandshakeOnly"
	ClientRequestKindSecretsRequest ClientRequestKind = "SecretsRequest"
)

// ClientRequest wraps a request to a realm.
type ClientRequest struct {
	Realm     types.RealmID     `cbor:"realm"`
	AuthToken types.AuthToken   `cbor:"auth_token"`
	SessionID types.SessionID   `cbor:"session_id"`
	Kind      ClientRequestKind `cbor:"kind"`
	Encrypted NoiseRequest      `cbor:"encrypted"`
}

// NoiseHandshakeRequest wraps a handshake request (matches Rust enum variant).
type NoiseHandshakeRequest struct {
	Handshake *noise.HandshakeRequest `cbor:"handshake"`
}

// NoiseRequest wraps either a handshake or transport message.
type NoiseRequest struct {
	Handshake *NoiseHandshakeRequest `cbor:"Handshake,omitempty"`
	Transport *NoiseTransportRequest `cbor:"Transport,omitempty"`
}

// NoiseTransportRequest contains encrypted data.
type NoiseTransportRequest struct {
	Ciphertext []byte `cbor:"ciphertext"`
}

// ClientResponse is the response from a realm.
type ClientResponse struct {
	Ok                *NoiseResponse `cbor:"Ok,omitempty"`
	Unavailable       bool           `cbor:"Unavailable,omitempty"`
	InvalidAuth       bool           `cbor:"InvalidAuth,omitempty"`
	MissingSession    bool           `cbor:"MissingSession,omitempty"`
	SessionError      bool           `cbor:"SessionError,omitempty"`
	DecodingError     bool           `cbor:"DecodingError,omitempty"`
	PayloadTooLarge   bool           `cbor:"PayloadTooLarge,omitempty"`
	RateLimitExceeded bool           `cbor:"RateLimitExceeded,omitempty"`
}

// NoiseResponse wraps either a handshake or transport response.
type NoiseResponse struct {
	Handshake *NoiseHandshakeResponse `cbor:"Handshake,omitempty"`
	Transport *NoiseTransportResponse `cbor:"Transport,omitempty"`
}

// Duration represents a Rust std::time::Duration serialized via serde.
// Rust serializes Duration as {"secs": u64, "nanos": u32}.
type Duration struct {
	Secs  uint64 `cbor:"secs"`
	Nanos uint32 `cbor:"nanos"`
}

// ToDuration converts to Go's time.Duration.
func (d Duration) ToDuration() time.Duration {
	return time.Duration(d.Secs)*time.Second + time.Duration(d.Nanos)*time.Nanosecond
}

// NoiseHandshakeResponse contains handshake data.
type NoiseHandshakeResponse struct {
	Handshake       noise.HandshakeResponse `cbor:"handshake"`
	SessionLifetime Duration                `cbor:"session_lifetime"`
}

// NoiseTransportResponse contains encrypted data.
type NoiseTransportResponse struct {
	Ciphertext []byte `cbor:"ciphertext"`
}

// SecretsRequest is the inner request for secret operations.
// It uses custom CBOR marshaling to match the expected format:
// - Register1/Recover1: serialize as string "Register1"/"Recover1"
// - Register2/Recover2/Recover3: serialize as {"Register2": {...}}, etc.
type SecretsRequest struct {
	Register1 bool              `cbor:"-"`
	Register2 *Register2Request `cbor:"-"`
	Recover1  bool              `cbor:"-"`
	Recover2  *Recover2Request  `cbor:"-"`
	Recover3  *Recover3Request  `cbor:"-"`
}

// NeedsForwardSecrecy reports whether this request must be sent via a Noise
// transport message (not piggybacked on handshake payload).
func (s *SecretsRequest) NeedsForwardSecrecy() bool {
	if s == nil {
		return false
	}
	return s.Register2 != nil || s.Recover2 != nil || s.Recover3 != nil
}

func marshalNamedSecretsRequest[T any](name string, req *T) ([]byte, error) {
	return cbor.Marshal(map[string]*T{name: req})
}

// MarshalCBOR implements custom CBOR serialization for SecretsRequest.
// The Juicebox realm expects:
// - Empty structs (Register1/Recover1) as a string: "Register1"/"Recover1"
// - Structs with fields as a map: {"Register2": {...}}, {"Recover2": {...}}, etc.
func (s *SecretsRequest) MarshalCBOR() ([]byte, error) {
	if s.Register1 {
		return cbor.Marshal("Register1")
	}
	if s.Register2 != nil {
		return marshalNamedSecretsRequest("Register2", s.Register2)
	}
	if s.Recover1 {
		return cbor.Marshal("Recover1")
	}
	if s.Recover2 != nil {
		return marshalNamedSecretsRequest("Recover2", s.Recover2)
	}
	if s.Recover3 != nil {
		return marshalNamedSecretsRequest("Recover3", s.Recover3)
	}
	return nil, nil
}

// Register2Request is the request for phase 2 of registration.
type Register2Request struct {
	Version                   types.RegistrationVersion           `cbor:"version"`
	OprfPrivateKey            []byte                              `cbor:"oprf_private_key"`
	OprfSignedPublicKey       OprfSignedPublicKey                 `cbor:"oprf_signed_public_key"`
	UnlockKeyCommitment       types.UnlockKeyCommitment           `cbor:"unlock_key_commitment"`
	UnlockKeyTag              types.UnlockKeyTag                  `cbor:"unlock_key_tag"`
	EncryptionKeyScalarShare  []byte                              `cbor:"encryption_key_scalar_share"`
	EncryptedSecret           types.EncryptedUserSecret           `cbor:"encrypted_secret"`
	EncryptedSecretCommitment types.EncryptedUserSecretCommitment `cbor:"encrypted_secret_commitment"`
	Policy                    Policy                              `cbor:"policy"`
}

// Policy defines restrictions on how a secret may be accessed.
type Policy struct {
	NumGuesses uint16 `cbor:"num_guesses"`
}

// Recover2Request is the request for phase 2 of recovery.
type Recover2Request struct {
	Version          types.RegistrationVersion `cbor:"version"`
	OprfBlindedInput []byte                    `cbor:"oprf_blinded_input"`
}

// Recover3Request is the request for phase 3 of recovery.
type Recover3Request struct {
	Version      types.RegistrationVersion `cbor:"version"`
	UnlockKeyTag types.UnlockKeyTag        `cbor:"unlock_key_tag"`
}

// SecretsResponse is the response for secret operations.
type SecretsResponse struct {
	Register1 *Register1Response `cbor:"Register1,omitempty"`
	Register2 *Register2Response `cbor:"Register2,omitempty"`
	Recover1  *Recover1Response  `cbor:"Recover1,omitempty"`
	Recover2  *Recover2Response  `cbor:"Recover2,omitempty"`
	Recover3  *Recover3Response  `cbor:"Recover3,omitempty"`
}

// Register1Response is the response for phase 1 of registration.
type Register1Response struct {
	Ok bool `cbor:"Ok,omitempty"`
}

func parseOKUnitVariant(data []byte, responseName string) error {
	var variant string
	if err := cbor.Unmarshal(data, &variant); err == nil {
		if variant == "Ok" {
			return nil
		}
		return fmt.Errorf("unknown %s variant: %s", responseName, variant)
	}

	var parsed map[string]cbor.RawMessage
	if err := cbor.Unmarshal(data, &parsed); err != nil {
		return err
	}
	if _, ok := parsed["Ok"]; !ok {
		return fmt.Errorf("unknown %s payload", responseName)
	}
	return nil
}

// UnmarshalCBOR supports both string unit variants (e.g. "Ok") and map-encoded variants.
func (r *Register1Response) UnmarshalCBOR(data []byte) error {
	if err := parseOKUnitVariant(data, "Register1Response"); err != nil {
		return err
	}
	r.Ok = true
	return nil
}

// Register2Response is the response for phase 2 of registration.
type Register2Response struct {
	Ok bool `cbor:"Ok,omitempty"`
}

// UnmarshalCBOR supports both string unit variants (e.g. "Ok") and map-encoded variants.
func (r *Register2Response) UnmarshalCBOR(data []byte) error {
	if err := parseOKUnitVariant(data, "Register2Response"); err != nil {
		return err
	}
	r.Ok = true
	return nil
}

// Recover1Response is the response for phase 1 of recovery.
type Recover1Response struct {
	Ok            *Recover1ResponseOk `cbor:"Ok,omitempty"`
	NotRegistered bool                `cbor:"NotRegistered,omitempty"`
	NoGuesses     bool                `cbor:"NoGuesses,omitempty"`
}

// Recover1ResponseOk contains the version from phase 1.
type Recover1ResponseOk struct {
	Version types.RegistrationVersion `cbor:"version"`
}

// Recover2Response is the response for phase 2 of recovery.
type Recover2Response struct {
	Ok              *Recover2ResponseOk `cbor:"Ok,omitempty"`
	VersionMismatch bool                `cbor:"VersionMismatch,omitempty"`
	NotRegistered   bool                `cbor:"NotRegistered,omitempty"`
	NoGuesses       bool                `cbor:"NoGuesses,omitempty"`
}

// Recover2ResponseOk contains the OPRF result from phase 2.
type Recover2ResponseOk struct {
	OprfSignedPublicKey OprfSignedPublicKey       `cbor:"oprf_signed_public_key"`
	OprfBlindedResult   []byte                    `cbor:"oprf_blinded_result"`
	OprfProof           OprfProof                 `cbor:"oprf_proof"`
	UnlockKeyCommitment types.UnlockKeyCommitment `cbor:"unlock_key_commitment"`
	NumGuesses          uint16                    `cbor:"num_guesses"`
	GuessCount          uint16                    `cbor:"guess_count"`
}

// OprfSignedPublicKey is the server's signed OPRF public key.
type OprfSignedPublicKey struct {
	PublicKey    []byte   `cbor:"public_key"`
	VerifyingKey [32]byte `cbor:"verifying_key"`
	Signature    [64]byte `cbor:"signature"`
}

// OprfProof is the DLEQ proof from the server.
type OprfProof struct {
	C     []byte `cbor:"c"`
	BetaZ []byte `cbor:"beta_z"`
}

// Recover3Response is the response for phase 3 of recovery.
type Recover3Response struct {
	Ok              *Recover3ResponseOk `cbor:"Ok,omitempty"`
	VersionMismatch bool                `cbor:"VersionMismatch,omitempty"`
	NotRegistered   bool                `cbor:"NotRegistered,omitempty"`
	NoGuesses       bool                `cbor:"NoGuesses,omitempty"`
	BadUnlockKeyTag *BadUnlockKeyTag    `cbor:"BadUnlockKeyTag,omitempty"`
}

// Recover3ResponseOk contains the secret shares from phase 3.
type Recover3ResponseOk struct {
	EncryptionKeyScalarShare  []byte                              `cbor:"encryption_key_scalar_share"`
	EncryptedSecret           types.EncryptedUserSecret           `cbor:"encrypted_secret"`
	EncryptedSecretCommitment types.EncryptedUserSecretCommitment `cbor:"encrypted_secret_commitment"`
}

// BadUnlockKeyTag indicates wrong PIN with guesses remaining.
type BadUnlockKeyTag struct {
	GuessesRemaining uint16 `cbor:"guesses_remaining"`
}

// PaddedSecretsResponse is used for constant-size responses (hardware realms).
type PaddedSecretsResponse struct {
	UnpaddedLength uint16    `cbor:"unpadded_length"`
	PaddedBytes    [436]byte `cbor:"padded_bytes"`
}

// Marshal serializes to CBOR.
func Marshal(v interface{}) ([]byte, error) {
	return cbor.Marshal(v)
}

// Unmarshal deserializes from CBOR.
func Unmarshal(data []byte, v interface{}) error {
	return cbor.Unmarshal(data, v)
}

// UnmarshalSecretsResponse extracts the inner response from a padded response.
func UnmarshalSecretsResponse(padded *PaddedSecretsResponse) (*SecretsResponse, error) {
	var resp SecretsResponse
	err := cbor.Unmarshal(padded.PaddedBytes[:padded.UnpaddedLength], &resp)
	return &resp, err
}
