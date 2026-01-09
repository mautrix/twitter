package connector

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
)

// userLoginKeyStore stores cryptographic keys.
// User-level keys (signing key) are stored in UserLoginMetadata.
// Conversation-level keys (conversation keys, tokens) are stored in PortalMetadata.
// All methods are thread-safe.
type userLoginKeyStore struct {
	login     *bridgev2.UserLogin
	connector *TwitterConnector
	meta      *UserLoginMetadata
	mu        sync.RWMutex
}

func newUserLoginKeyStore(login *bridgev2.UserLogin, connector *TwitterConnector) *userLoginKeyStore {
	return &userLoginKeyStore{
		login:     login,
		connector: connector,
		meta:      login.Metadata.(*UserLoginMetadata),
	}
}

// getPortalMetadata looks up a portal by conversation ID and returns its metadata.
func (ks *userLoginKeyStore) getPortalMetadata(ctx context.Context, conversationID string) (*PortalMetadata, *bridgev2.Portal, error) {
	portalKey := MakePortalKeyForConversation(conversationID, ks.login.ID, ks.connector.br.Config.SplitPortals)
	portal, err := ks.connector.br.GetPortalByKey(ctx, portalKey)
	if err != nil {
		return nil, nil, err
	}
	return portal.Metadata.(*PortalMetadata), portal, nil
}

func (ks *userLoginKeyStore) GetConversationKey(ctx context.Context, conversationID, keyVersion string) (*crypto.ConversationKey, error) {
	log := zerolog.Ctx(ctx)

	meta, _, err := ks.getPortalMetadata(ctx, conversationID)
	if err != nil {
		log.Warn().Err(err).
			Str("conversation_id", conversationID).
			Str("key_version", keyVersion).
			Msg("Failed to get portal metadata for conversation key")
		return nil, err
	}

	if meta.ConversationKeys == nil {
		log.Info().
			Str("conversation_id", conversationID).
			Str("key_version", keyVersion).
			Msg("No conversation keys in portal metadata")
		return nil, crypto.ErrKeyNotFound
	}

	data, ok := meta.ConversationKeys[keyVersion]
	if !ok {
		log.Info().
			Str("conversation_id", conversationID).
			Str("key_version", keyVersion).
			Int("total_keys", len(meta.ConversationKeys)).
			Msg("Conversation key not found in portal metadata")
		return nil, crypto.ErrKeyNotFound
	}

	log.Info().
		Str("conversation_id", conversationID).
		Str("key_version", keyVersion).
		Int("key_length", len(data.Key)).
		Str("key_prefix", truncateKeyHex(data.Key, 8)).
		Time("created_at", data.CreatedAt).
		Msg("Retrieved conversation key from portal metadata")

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

	meta, portal, err := ks.getPortalMetadata(ctx, key.ConversationID)
	if err != nil {
		log.Err(err).
			Str("conversation_id", key.ConversationID).
			Str("key_version", key.KeyVersion).
			Msg("Failed to get portal metadata for storing conversation key")
		return err
	}

	if meta.ConversationKeys == nil {
		meta.ConversationKeys = make(map[string]*ConversationKeyData)
	}

	log.Info().
		Str("conversation_id", key.ConversationID).
		Str("key_version", key.KeyVersion).
		Int("key_length", len(key.Key)).
		Str("key_prefix", truncateKeyHex(key.Key, 8)).
		Int("total_keys_before", len(meta.ConversationKeys)).
		Msg("Storing conversation key in portal metadata")

	meta.ConversationKeys[key.KeyVersion] = &ConversationKeyData{
		KeyVersion: key.KeyVersion,
		Key:        key.Key,
		CreatedAt:  key.CreatedAt,
		ExpiresAt:  key.ExpiresAt,
	}

	err = portal.Save(ctx)
	if err != nil {
		log.Err(err).
			Str("conversation_id", key.ConversationID).
			Str("key_version", key.KeyVersion).
			Msg("Failed to save portal metadata with conversation key")
	} else {
		log.Info().
			Str("conversation_id", key.ConversationID).
			Str("key_version", key.KeyVersion).
			Int("total_keys_after", len(meta.ConversationKeys)).
			Msg("Successfully saved conversation key to portal metadata")
	}
	return err
}

func (ks *userLoginKeyStore) DeleteConversationKey(ctx context.Context, conversationID, keyVersion string) error {
	meta, portal, err := ks.getPortalMetadata(ctx, conversationID)
	if err != nil {
		return err
	}

	if meta.ConversationKeys != nil {
		delete(meta.ConversationKeys, keyVersion)
	}
	return portal.Save(ctx)
}

func (ks *userLoginKeyStore) GetLatestConversationKey(ctx context.Context, conversationID string) (*crypto.ConversationKey, error) {
	log := zerolog.Ctx(ctx)

	meta, _, err := ks.getPortalMetadata(ctx, conversationID)
	if err != nil {
		log.Warn().Err(err).
			Str("conversation_id", conversationID).
			Msg("Failed to get portal metadata for latest conversation key")
		return nil, err
	}

	if len(meta.ConversationKeys) == 0 {
		log.Info().
			Str("conversation_id", conversationID).
			Msg("No conversation keys found in portal metadata")
		return nil, crypto.ErrKeyNotFound
	}

	// Find key with latest CreatedAt
	var latest *ConversationKeyData
	for _, data := range meta.ConversationKeys {
		if latest == nil || data.CreatedAt.After(latest.CreatedAt) {
			latest = data
		}
	}

	log.Info().
		Str("conversation_id", conversationID).
		Str("latest_key_version", latest.KeyVersion).
		Int("total_keys", len(meta.ConversationKeys)).
		Int("key_length", len(latest.Key)).
		Str("key_prefix", truncateKeyHex(latest.Key, 8)).
		Time("created_at", latest.CreatedAt).
		Msg("Retrieved latest conversation key from portal metadata")

	return &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     latest.KeyVersion,
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
	loginID := ParseUserLoginID(ks.login.ID)
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

func (ks *userLoginKeyStore) GetConversationToken(ctx context.Context, conversationID string) (string, error) {
	meta, _, err := ks.getPortalMetadata(ctx, conversationID)
	if err != nil {
		return "", err
	}

	if meta.ConversationToken == "" {
		return "", crypto.ErrKeyNotFound
	}
	return meta.ConversationToken, nil
}

func (ks *userLoginKeyStore) PutConversationToken(ctx context.Context, conversationID, token string) error {
	log := zerolog.Ctx(ctx)

	meta, portal, err := ks.getPortalMetadata(ctx, conversationID)
	if err != nil {
		log.Err(err).
			Str("conversation_id", conversationID).
			Msg("Failed to get portal metadata for storing conversation token")
		return err
	}

	meta.ConversationToken = token
	return portal.Save(ctx)
}

// truncateKeyHex returns a hex representation of the first n bytes of a key for logging.
func truncateKeyHex(key []byte, n int) string {
	if len(key) <= n {
		return hex.EncodeToString(key)
	}
	return hex.EncodeToString(key[:n]) + "..."
}
