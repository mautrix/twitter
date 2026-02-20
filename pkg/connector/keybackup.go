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

package connector

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo"
	twcrypto "go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
)

// KeyBackupData represents the encryption keys stored in Juicebox.
type KeyBackupData struct {
	SigningKey        string `json:"signing_key"`
	SecretKey         string `json:"secret_key"`
	SigningKeyVersion string `json:"signing_key_version"`
}

// FirstTimePINBootstrapData contains generated key material for first-time PIN setup.
type FirstTimePINBootstrapData struct {
	SecretKey                  string
	SigningKey                 string
	PublicKeySPKI              string
	SigningPublicKeySPKI       string
	IdentityPublicKeySignature string
	RawSecret                  []byte
}

func generatePrivateScalarAndPublicSPKI() (string, string, *ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate P-256 key pair: %w", err)
	}

	scalarBytes := key.D.FillBytes(make([]byte, 32))
	scalarB64 := base64.StdEncoding.EncodeToString(scalarBytes)
	spkiB64, err := twcrypto.EncodePublicKeySPKI(&key.PublicKey)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to encode SPKI: %w", err)
	}

	return scalarB64, spkiB64, key, nil
}

func decodePrivateScalarB64(scalarB64 string) ([]byte, error) {
	scalar, err := base64.StdEncoding.DecodeString(scalarB64)
	if err != nil {
		scalar, err = base64.RawStdEncoding.DecodeString(scalarB64)
		if err != nil {
			return nil, err
		}
	}
	if len(scalar) != 32 {
		return nil, fmt.Errorf("expected 32-byte scalar, got %d bytes", len(scalar))
	}
	return scalar, nil
}

// ComposeRawSecret builds the raw 64-byte Juicebox secret format: decrypt||signing.
func ComposeRawSecret(decryptScalarB64, signingScalarB64 string) ([]byte, error) {
	decryptScalar, err := decodePrivateScalarB64(decryptScalarB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode decrypt scalar: %w", err)
	}
	signingScalar, err := decodePrivateScalarB64(signingScalarB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signing scalar: %w", err)
	}

	raw := make([]byte, 64)
	copy(raw[:32], decryptScalar)
	copy(raw[32:], signingScalar)
	return raw, nil
}

// GenerateFirstTimePINBootstrapData generates key material for first-time PIN setup.
func GenerateFirstTimePINBootstrapData() (*FirstTimePINBootstrapData, error) {
	decryptScalarB64, publicKeySPKI, _, err := generatePrivateScalarAndPublicSPKI()
	if err != nil {
		return nil, err
	}
	signingScalarB64, signingPublicKeySPKI, signingPriv, err := generatePrivateScalarAndPublicSPKI()
	if err != nil {
		return nil, err
	}

	publicKeyDER, err := base64.StdEncoding.DecodeString(publicKeySPKI)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key SPKI: %w", err)
	}
	identitySig, err := twcrypto.Sign(signingPriv, publicKeyDER)
	if err != nil {
		return nil, fmt.Errorf("failed to sign identity public key: %w", err)
	}

	rawSecret, err := ComposeRawSecret(decryptScalarB64, signingScalarB64)
	if err != nil {
		return nil, err
	}

	return &FirstTimePINBootstrapData{
		SecretKey:                  decryptScalarB64,
		SigningKey:                 signingScalarB64,
		PublicKeySPKI:              publicKeySPKI,
		SigningPublicKeySPKI:       signingPublicKeySPKI,
		IdentityPublicKeySignature: identitySig,
		RawSecret:                  rawSecret,
	}, nil
}

func newJuiceboxClient(
	configJSON string,
	authTokens map[string]string,
	logger zerolog.Logger,
) (*juiceboxgo.Client, error) {
	config, err := juiceboxgo.ConfigurationFromJSON(configJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse juicebox config: %w", err)
	}

	typedAuthTokens := make(map[string]juiceboxgo.AuthToken, len(authTokens))
	for k, v := range authTokens {
		typedAuthTokens[strings.ToLower(k)] = juiceboxgo.AuthToken(v)
	}
	client, err := juiceboxgo.NewClient(config, typedAuthTokens, nil, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create juicebox client: %w", err)
	}
	return client, nil
}

func splitRawRecoveredSecret(secret []byte) (privBytes, pubBytes []byte) {
	switch {
	case len(secret) >= 64 && len(secret)%2 == 0:
		// Even length >= 64: split in half.
		half := len(secret) / 2
		return secret[:half], secret[half:]
	case len(secret) > 32:
		// More than 32 bytes: first 32 are private, rest are public.
		return secret[:32], secret[32:]
	default:
		// Just the private key.
		return secret, nil
	}
}

func parseRecoveredKeyBackupData(secret []byte, logger zerolog.Logger) (*KeyBackupData, error) {
	// Try JSON first (some implementations store JSON).
	if len(secret) > 0 && secret[0] == '{' {
		var data KeyBackupData
		if err := json.Unmarshal(secret, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal key data: %w", err)
		}
		return &data, nil
	}

	// Otherwise, it's raw binary key data (like reverse-xchat handles).
	// Format: concatenated key bytes, typically 64 bytes (32 byte private + 32 byte public)
	// or just 32 bytes for a single key.
	privBytes, pubBytes := splitRawRecoveredSecret(secret)
	logger.Debug().
		Int("secret_key_len", len(privBytes)).
		Int("signing_key_len", len(pubBytes)).
		Msg("Parsed raw key bytes")

	data := &KeyBackupData{
		SecretKey: base64.StdEncoding.EncodeToString(privBytes),
	}
	if len(pubBytes) > 0 {
		data.SigningKey = base64.StdEncoding.EncodeToString(pubBytes)
	}
	return data, nil
}

// RegisterSecretToJuicebox stores a freshly generated 64-byte secret using a first-time PIN.
func RegisterSecretToJuicebox(ctx context.Context, configJSON string, authTokens map[string]string, pin string, userInfo string, rawSecret []byte, maxGuessCount int, logger zerolog.Logger) error {
	if len(rawSecret) != 64 {
		return fmt.Errorf("invalid raw secret size: expected 64 bytes, got %d", len(rawSecret))
	}
	if maxGuessCount <= 0 {
		maxGuessCount = 20
	}
	if maxGuessCount > (1<<16 - 1) {
		return fmt.Errorf("invalid max guess count: %d", maxGuessCount)
	}
	if len(authTokens) == 0 {
		return fmt.Errorf("no auth tokens available for registration")
	}

	client, err := newJuiceboxClient(configJSON, authTokens, logger)
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Register(ctx, []byte(pin), rawSecret, []byte(userInfo), uint16(maxGuessCount))
	if err != nil {
		return fmt.Errorf("failed to register keys: %w", err)
	}
	return nil
}

// RecoverKeysFromJuicebox retrieves encryption keys using PIN from Juicebox.
// Called once during login as alternative to manual key entry.
// authTokens is a map of realm ID (hex string) to pre-fetched JWT auth token.
func RecoverKeysFromJuicebox(ctx context.Context, configJSON string, authTokens map[string]string, pin string, userInfo string, logger zerolog.Logger) (*KeyBackupData, error) {
	logger.Debug().
		Int("config_json_len", len(configJSON)).
		Msg("Creating Juicebox configuration from JSON")

	client, err := newJuiceboxClient(configJSON, authTokens, logger)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	logger.Debug().Msg("Juicebox client created, calling Recover")

	secret, err := client.Recover(ctx, []byte(pin), []byte(userInfo))
	if err != nil {
		return nil, fmt.Errorf("failed to recover keys: %w", err)
	}

	logger.Debug().
		Int("secret_len", len(secret)).
		Msg("Juicebox secret recovered")

	return parseRecoveredKeyBackupData(secret, logger)
}
