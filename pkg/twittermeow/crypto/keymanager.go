package crypto

import (
	"context"
	"crypto/ecdsa"
)

// KeyStore provides persistent storage for cryptographic keys.
// Implementations should be thread-safe.
type KeyStore interface {
	// Conversation key operations
	GetConversationKey(ctx context.Context, conversationID, keyVersion string) (*ConversationKey, error)
	PutConversationKey(ctx context.Context, key *ConversationKey) error
	DeleteConversationKey(ctx context.Context, conversationID, keyVersion string) error
	GetLatestConversationKey(ctx context.Context, conversationID string) (*ConversationKey, error)

	// Own signing key (the user's private key for signing outgoing messages)
	GetOwnSigningKey(ctx context.Context) (*SigningKeyPair, error)
	PutOwnSigningKey(ctx context.Context, key *SigningKeyPair) error

	// Public key lookup (for signature verification of others)
	GetPublicKey(ctx context.Context, userID, keyVersion string) (*PublicKeyInfo, error)
	PutPublicKey(ctx context.Context, key *PublicKeyInfo) error

	// Conversation token operations (server-provided tokens for sending messages)
	GetConversationToken(ctx context.Context, conversationID string) (string, error)
	PutConversationToken(ctx context.Context, conversationID, token string) error
}

// NoOpKeyStore is a KeyStore that stores nothing - useful for per-call key usage.
type NoOpKeyStore struct{}

func (n *NoOpKeyStore) GetConversationKey(ctx context.Context, conversationID, keyVersion string) (*ConversationKey, error) {
	return nil, ErrKeyNotFound
}

func (n *NoOpKeyStore) PutConversationKey(ctx context.Context, key *ConversationKey) error {
	return nil
}

func (n *NoOpKeyStore) DeleteConversationKey(ctx context.Context, conversationID, keyVersion string) error {
	return nil
}

func (n *NoOpKeyStore) GetLatestConversationKey(ctx context.Context, conversationID string) (*ConversationKey, error) {
	return nil, ErrKeyNotFound
}

func (n *NoOpKeyStore) GetOwnSigningKey(ctx context.Context) (*SigningKeyPair, error) {
	return nil, ErrKeyNotFound
}

func (n *NoOpKeyStore) PutOwnSigningKey(ctx context.Context, key *SigningKeyPair) error {
	return nil
}

func (n *NoOpKeyStore) GetPublicKey(ctx context.Context, userID, keyVersion string) (*PublicKeyInfo, error) {
	return nil, ErrKeyNotFound
}

func (n *NoOpKeyStore) PutPublicKey(ctx context.Context, key *PublicKeyInfo) error {
	return nil
}

func (n *NoOpKeyStore) GetConversationToken(ctx context.Context, conversationID string) (string, error) {
	return "", ErrKeyNotFound
}

func (n *NoOpKeyStore) PutConversationToken(ctx context.Context, conversationID, token string) error {
	return nil
}

// KeyManager manages cryptographic keys with in-memory caching.
// All methods are thread-safe.
type KeyManager struct {
	store KeyStore
}

// NewKeyManager creates a KeyManager with the given store and config.
// If store is nil, a NoOpKeyStore is used.
func NewKeyManager(store KeyStore) *KeyManager {
	if store == nil {
		store = &NoOpKeyStore{}
	}
	return &KeyManager{
		store: store,
	}
}

// GetConversationKey retrieves a conversation key, checking cache first.
func (km *KeyManager) GetConversationKey(ctx context.Context, conversationID, keyVersion string) (*ConversationKey, error) {
	return km.store.GetConversationKey(ctx, conversationID, keyVersion)
}

// GetLatestConversationKey retrieves the latest conversation key for a conversation.
func (km *KeyManager) GetLatestConversationKey(ctx context.Context, conversationID string) (*ConversationKey, error) {
	return km.store.GetLatestConversationKey(ctx, conversationID)
}

// PutConversationKey stores a conversation key in both cache and persistent store.
func (km *KeyManager) PutConversationKey(ctx context.Context, key *ConversationKey) error {
	return km.store.PutConversationKey(ctx, key)
}

// GetOwnSigningKey returns the user's own signing key for message signing.
func (km *KeyManager) GetOwnSigningKey(ctx context.Context) (*SigningKeyPair, error) {
	return km.store.GetOwnSigningKey(ctx)
}

// SetOwnSigningKey sets the user's signing key from a base64-encoded private scalar.
func (km *KeyManager) SetOwnSigningKey(ctx context.Context, privateKeyB64 string) error {
	key, err := LoadSigningKeyPair("", "", "", privateKeyB64)
	if err != nil {
		return err
	}

	return km.store.PutOwnSigningKey(ctx, key)
}

// SetOwnSigningKeyPair sets the user's signing key directly.
func (km *KeyManager) SetOwnSigningKeyPair(ctx context.Context, key *SigningKeyPair) error {
	return km.store.PutOwnSigningKey(ctx, key)
}

// GetPublicKeyForVerification retrieves a public key for verifying a signature.
// If spkiB64 is provided, it parses and caches it; otherwise loads from store.
func (km *KeyManager) GetPublicKeyForVerification(ctx context.Context, userID, keyVersion, spkiB64 string) (*ecdsa.PublicKey, error) {
	var info *PublicKeyInfo
	var err error

	if spkiB64 != "" {
		info, err = LoadPublicKeyInfo(userID, keyVersion, spkiB64)
		if err != nil {
			return nil, err
		}
		// Store for future use (ignore error, caching is best-effort)
		_ = km.store.PutPublicKey(ctx, info)
	} else {
		info, err = km.store.GetPublicKey(ctx, userID, keyVersion)
		if err != nil {
			return nil, err
		}
	}

	return info.PublicKey, nil
}

// GetConversationToken retrieves a server-provided conversation token.
func (km *KeyManager) GetConversationToken(ctx context.Context, conversationID string) (string, error) {
	return km.store.GetConversationToken(ctx, conversationID)
}

// PutConversationToken stores a server-provided conversation token.
func (km *KeyManager) PutConversationToken(ctx context.Context, conversationID, token string) error {
	return km.store.PutConversationToken(ctx, conversationID, token)
}
