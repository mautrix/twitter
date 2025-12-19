package connector

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
)

// userLoginKeyStore stores cryptographic keys in the user login metadata.
// Signing keys and conversation keys are persisted; public keys are not.
// All methods are thread-safe.
type userLoginKeyStore struct {
	login *bridgev2.UserLogin
	meta  *UserLoginMetadata
	mu    sync.RWMutex
}

// conversationCacheKey builds a cache key for conversation key storage.
func conversationCacheKey(conversationID, keyVersion string) string {
	return conversationID + ":" + keyVersion
}

func newUserLoginKeyStore(login *bridgev2.UserLogin) *userLoginKeyStore {
	return &userLoginKeyStore{
		login: login,
		meta:  ensureUserLoginMetadata(login),
	}
}

func ensureUserLoginMetadata(login *bridgev2.UserLogin) *UserLoginMetadata {
	if meta, ok := login.Metadata.(*UserLoginMetadata); ok && meta != nil {
		if meta.ConversationKeys == nil {
			meta.ConversationKeys = make(map[string]*ConversationKeyData)
		}
		if meta.ConversationTokens == nil {
			meta.ConversationTokens = make(map[string]string)
		}
		if meta.MaxUserSequenceID == "" && meta.MaxSequenceID != "" {
			// Migrate old field name to the new one.
			meta.MaxUserSequenceID = meta.MaxSequenceID
			meta.MaxSequenceID = ""
		}
		if meta.UserID == "" {
			meta.UserID = string(login.ID)
		}
		return meta
	}
	meta := &UserLoginMetadata{
		ConversationKeys:   make(map[string]*ConversationKeyData),
		ConversationTokens: make(map[string]string),
		UserID:             string(login.ID),
	}
	login.Metadata = meta
	return meta
}

func (ks *userLoginKeyStore) GetConversationKey(ctx context.Context, conversationID, keyVersion string) (*crypto.ConversationKey, error) {
	log := zerolog.Ctx(ctx)
	cacheKey := conversationCacheKey(conversationID, keyVersion)

	ks.mu.RLock()
	data, ok := ks.meta.ConversationKeys[cacheKey]
	totalKeys := len(ks.meta.ConversationKeys)
	ks.mu.RUnlock()

	if !ok {
		log.Info().
			Str("conversation_id", conversationID).
			Str("key_version", keyVersion).
			Str("cache_key", cacheKey).
			Int("total_stored_keys", totalKeys).
			Msg("Conversation key not found in keystore")
		return nil, crypto.ErrKeyNotFound
	}
	log.Info().
		Str("conversation_id", conversationID).
		Str("key_version", keyVersion).
		Str("cache_key", cacheKey).
		Int("key_length", len(data.Key)).
		Str("key_prefix", truncateKeyHex(data.Key, 8)).
		Time("created_at", data.CreatedAt).
		Msg("Retrieved conversation key from keystore")
	return &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     data.KeyVersion,
		Key:            data.Key,
		CreatedAt:      data.CreatedAt,
		ExpiresAt:      data.ExpiresAt,
	}, nil
}

func (ks *userLoginKeyStore) PutConversationKey(ctx context.Context, key *crypto.ConversationKey) error {
	log := zerolog.Ctx(ctx)
	if key == nil {
		return fmt.Errorf("conversation key cannot be nil")
	}
	cacheKey := conversationCacheKey(key.ConversationID, key.KeyVersion)

	ks.mu.Lock()
	defer ks.mu.Unlock()

	log.Info().
		Str("conversation_id", key.ConversationID).
		Str("key_version", key.KeyVersion).
		Str("cache_key", cacheKey).
		Int("key_length", len(key.Key)).
		Str("key_prefix", truncateKeyHex(key.Key, 8)).
		Int("total_stored_keys_before", len(ks.meta.ConversationKeys)).
		Msg("Storing conversation key in keystore")
	ks.meta.ConversationKeys[cacheKey] = &ConversationKeyData{
		KeyVersion: key.KeyVersion,
		Key:        key.Key,
		CreatedAt:  key.CreatedAt,
		ExpiresAt:  key.ExpiresAt,
	}

	err := ks.login.Save(ctx)
	if err != nil {
		log.Err(err).
			Str("conversation_id", key.ConversationID).
			Str("key_version", key.KeyVersion).
			Msg("Failed to save conversation key to database")
	} else {
		log.Info().
			Str("conversation_id", key.ConversationID).
			Str("key_version", key.KeyVersion).
			Int("total_stored_keys_after", len(ks.meta.ConversationKeys)).
			Msg("Successfully saved conversation key to database")
	}
	return err
}

func (ks *userLoginKeyStore) DeleteConversationKey(ctx context.Context, conversationID, keyVersion string) error {
	cacheKey := conversationCacheKey(conversationID, keyVersion)

	ks.mu.Lock()
	defer ks.mu.Unlock()

	delete(ks.meta.ConversationKeys, cacheKey)
	return ks.login.Save(ctx)
}

func (ks *userLoginKeyStore) GetLatestConversationKey(ctx context.Context, conversationID string) (*crypto.ConversationKey, error) {
	log := zerolog.Ctx(ctx)
	var latest *ConversationKeyData
	var latestVersion string
	prefix := conversationID + ":"

	ks.mu.RLock()
	var matchingKeys []string
	for cacheKey, data := range ks.meta.ConversationKeys {
		if !strings.HasPrefix(cacheKey, prefix) {
			continue
		}
		matchingKeys = append(matchingKeys, cacheKey)
		if latest == nil || data.CreatedAt.After(latest.CreatedAt) {
			latest = data
			latestVersion = data.KeyVersion
		}
	}
	totalKeys := len(ks.meta.ConversationKeys)
	ks.mu.RUnlock()

	if latest == nil {
		log.Info().
			Str("conversation_id", conversationID).
			Int("total_stored_keys", totalKeys).
			Msg("No conversation keys found for conversation")
		return nil, crypto.ErrKeyNotFound
	}

	log.Info().
		Str("conversation_id", conversationID).
		Str("latest_key_version", latestVersion).
		Strs("matching_keys", matchingKeys).
		Int("key_length", len(latest.Key)).
		Str("key_prefix", truncateKeyHex(latest.Key, 8)).
		Time("created_at", latest.CreatedAt).
		Msg("Retrieved latest conversation key from keystore")

	return &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     latestVersion,
		Key:            latest.Key,
		CreatedAt:      latest.CreatedAt,
		ExpiresAt:      latest.ExpiresAt,
	}, nil
}

func (ks *userLoginKeyStore) GetOwnSigningKey(ctx context.Context) (*crypto.SigningKeyPair, error) {
	log := zerolog.Ctx(ctx)

	ks.mu.RLock()
	secretKey := ks.meta.SecretKey
	signingKey := ks.meta.SigningKey
	signingKeyVersion := ks.meta.SigningKeyVersion
	loginID := string(ks.login.ID)
	ks.mu.RUnlock()

	log.Info().
		Bool("has_secret_key", secretKey != "").
		Bool("has_signing_key", signingKey != "").
		Str("signing_key_version", signingKeyVersion).
		Str("login_id", loginID).
		Msg("GetOwnSigningKey called")

	if secretKey == "" && signingKey == "" {
		log.Warn().Msg("No signing keys found in metadata")
		return nil, crypto.ErrKeyNotFound
	}
	key, err := crypto.LoadSigningKeyPair(loginID, signingKeyVersion, signingKey, secretKey)
	if err != nil {
		log.Err(err).Msg("Failed to load signing key pair")
		return nil, err
	}
	log.Info().
		Bool("has_signing_key_obj", key.SigningKey != nil).
		Bool("has_decrypt_key_obj", key.DecryptKey != nil).
		Msg("Successfully loaded signing key pair")
	return key, nil
}

func (ks *userLoginKeyStore) PutOwnSigningKey(ctx context.Context, key *crypto.SigningKeyPair) error {
	if key == nil {
		return fmt.Errorf("signing key cannot be nil")
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.meta.SecretKey = key.DecryptKeyB64
	ks.meta.SigningKey = key.SigningKeyB64
	ks.meta.SigningKeyVersion = key.KeyVersion
	return ks.login.Save(ctx)
}

func (ks *userLoginKeyStore) GetPublicKey(_ context.Context, _, _ string) (*crypto.PublicKeyInfo, error) {
	return nil, crypto.ErrKeyNotFound
}

func (ks *userLoginKeyStore) PutPublicKey(_ context.Context, _ *crypto.PublicKeyInfo) error {
	return nil
}

func (ks *userLoginKeyStore) GetConversationToken(_ context.Context, conversationID string) (string, error) {
	ks.mu.RLock()
	token, ok := ks.meta.ConversationTokens[conversationID]
	ks.mu.RUnlock()

	if !ok {
		return "", crypto.ErrKeyNotFound
	}
	return token, nil
}

func (ks *userLoginKeyStore) PutConversationToken(ctx context.Context, conversationID, token string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.meta.ConversationTokens[conversationID] = token
	return ks.login.Save(ctx)
}

// truncateKeyHex returns a hex representation of the first n bytes of a key for logging.
func truncateKeyHex(key []byte, n int) string {
	if len(key) <= n {
		return hex.EncodeToString(key)
	}
	return hex.EncodeToString(key[:n]) + "..."
}
