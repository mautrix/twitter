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

package juiceboxgo

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// Configuration holds the realm configuration for Juicebox operations.
type Configuration struct {
	Realms            []Realm        `json:"realms"`
	RegisterThreshold uint32         `json:"register_threshold"`
	RecoverThreshold  uint32         `json:"recover_threshold"`
	PinHashingMode    PinHashingMode `json:"pin_hashing_mode"`
}

// realmJSON is used for JSON unmarshaling with hex-encoded fields.
type realmJSON struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	PublicKey string `json:"public_key,omitempty"`
}

// configJSON is used for JSON unmarshaling.
type configJSON struct {
	Realms            []realmJSON    `json:"realms"`
	RegisterThreshold uint32         `json:"register_threshold"`
	RecoverThreshold  uint32         `json:"recover_threshold"`
	PinHashingMode    PinHashingMode `json:"pin_hashing_mode"`
}

// ConfigurationFromJSON parses a configuration from JSON.
func ConfigurationFromJSON(jsonData string) (*Configuration, error) {
	var raw configJSON
	if err := json.Unmarshal([]byte(jsonData), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	config := &Configuration{
		RegisterThreshold: raw.RegisterThreshold,
		RecoverThreshold:  raw.RecoverThreshold,
		PinHashingMode:    raw.PinHashingMode,
		Realms:            make([]Realm, len(raw.Realms)),
	}

	for i, r := range raw.Realms {
		id, err := ParseRealmID(r.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid realm ID at index %d: %w", i, err)
		}

		realm := Realm{
			ID:      id,
			Address: r.Address,
		}

		if r.PublicKey != "" {
			pk, err := hex.DecodeString(r.PublicKey)
			if err != nil {
				return nil, fmt.Errorf("invalid public key at index %d: %w", i, err)
			}
			realm.PublicKey = pk
		}

		config.Realms[i] = realm
	}

	// Sort realms by ID to ensure consistent share indices.
	// This matches the Rust SDK behavior which sorts realms before assigning share indices.
	sort.Slice(config.Realms, func(i, j int) bool {
		return bytes.Compare(config.Realms[i].ID[:], config.Realms[j].ID[:]) < 0
	})

	// Validate threshold
	if config.RecoverThreshold == 0 {
		return nil, fmt.Errorf("recover_threshold must be > 0")
	}
	if int(config.RecoverThreshold) > len(config.Realms) {
		return nil, fmt.Errorf("recover_threshold (%d) exceeds number of realms (%d)",
			config.RecoverThreshold, len(config.Realms))
	}

	return config, nil
}

// ShareIndex returns the share index for a realm (1-based).
func (c *Configuration) ShareIndex(realmID RealmID) (uint32, bool) {
	for i, r := range c.Realms {
		if r.ID == realmID {
			return uint32(i + 1), true
		}
	}
	return 0, false
}
