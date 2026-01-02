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
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo"
)

// KeyBackupData represents the encryption keys stored in Juicebox.
type KeyBackupData struct {
	SigningKey        string `json:"signing_key"`
	SecretKey         string `json:"secret_key"`
	SigningKeyVersion string `json:"signing_key_version"`
}

// RecoverKeysFromJuicebox retrieves encryption keys using PIN from Juicebox.
// Called once during login as alternative to manual key entry.
// authTokens is a map of realm ID (hex string) to pre-fetched JWT auth token.
func RecoverKeysFromJuicebox(ctx context.Context, configJSON string, authTokens map[string]string, pin string, userInfo string, logger zerolog.Logger) (*KeyBackupData, error) {
	logger.Debug().
		Int("config_json_len", len(configJSON)).
		Msg("Creating Juicebox configuration from JSON")

	config, err := juiceboxgo.ConfigurationFromJSON(configJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse juicebox config: %w", err)
	}

	logger.Debug().Msg("Juicebox configuration created, creating client")

	// Convert string auth tokens to AuthToken type
	typedAuthTokens := make(map[string]juiceboxgo.AuthToken)
	for k, v := range authTokens {
		typedAuthTokens[k] = juiceboxgo.AuthToken(v)
	}

	client, err := juiceboxgo.NewClient(config, typedAuthTokens, nil, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create juicebox client: %w", err)
	}
	defer client.Close()

	logger.Debug().Msg("Juicebox client created, calling Recover")

	secret, err := client.Recover(ctx, []byte(pin), []byte(userInfo))
	if err != nil {
		return nil, fmt.Errorf("failed to recover keys: %w", err)
	}

	logger.Debug().
		Int("secret_len", len(secret)).
		Str("secret_hex_preview", fmt.Sprintf("%x", secret[:min(32, len(secret))])).
		Msg("Juicebox secret recovered")

	// Try JSON first (some implementations store JSON)
	if len(secret) > 0 && secret[0] == '{' {
		var data KeyBackupData
		if err := json.Unmarshal(secret, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal key data: %w", err)
		}
		return &data, nil
	}

	// Otherwise, it's raw binary key data (like reverse-xchat handles)
	// Format: concatenated key bytes, typically 64 bytes (32 byte private + 32 byte public)
	// or just 32 bytes for a single key
	var privBytes, pubBytes []byte
	if len(secret) >= 64 && len(secret)%2 == 0 {
		// Even length >= 64: split in half
		half := len(secret) / 2
		privBytes = secret[:half]
		pubBytes = secret[half:]
	} else if len(secret) > 32 {
		// More than 32 bytes: first 32 are private, rest are public
		privBytes = secret[:32]
		pubBytes = secret[32:]
	} else {
		// Just the private key
		privBytes = secret
	}

	logger.Debug().
		Int("secret_key_len", len(privBytes)).
		Int("signing_key_len", len(pubBytes)).
		Msg("Parsed raw key bytes")

	// Convert to base64 for storage (matching expected format)
	// First half is secret key, second half is signing key
	data := &KeyBackupData{
		SecretKey: base64.StdEncoding.EncodeToString(privBytes),
	}
	if len(pubBytes) > 0 {
		data.SigningKey = base64.StdEncoding.EncodeToString(pubBytes)
	}

	return data, nil
}
