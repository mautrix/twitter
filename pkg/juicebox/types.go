package juicebox

// Pin represents a user's PIN for protecting secrets (max 128 bytes).
type Pin []byte

// Secret represents the secret to be stored/recovered (max 128 bytes).
type Secret []byte

// UserInfo represents additional context for the secret.
type UserInfo []byte

// RealmID is a 16-byte realm identifier.
type RealmID [16]byte

// SecretID is a 16-byte secret identifier.
type SecretID [16]byte
