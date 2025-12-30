package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// Common errors for key operations.
var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrInvalidKeyFormat = errors.New("invalid key format")
	ErrKeyExpired       = errors.New("key has expired")
	ErrSignatureInvalid = errors.New("signature verification failed")
)

// ConversationKey represents a decrypted conversation key for XChat.
type ConversationKey struct {
	ConversationID string     `json:"conversation_id"`
	KeyVersion     string     `json:"key_version"`
	Key            []byte     `json:"key"` // 32-byte secretbox key
	CreatedAt      time.Time  `json:"created_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// SigningKeyPair represents the cryptographic keys for XChat.
// It contains two separate keys:
// - SigningKey: ECDSA P-256 key pair for signing messages
// - DecryptKey: ECDSA P-256 private key for decrypting conversation keys
type SigningKeyPair struct {
	UserID     string `json:"user_id"`
	KeyVersion string `json:"key_version"`

	// Parsed key objects (not serialized)
	SigningKey *ecdsa.PrivateKey `json:"-"` // For signing messages
	DecryptKey *ecdsa.PrivateKey `json:"-"` // For decrypting conversation keys

	// Serialized forms for storage (both are base64-encoded 32-byte scalars)
	SigningKeyB64 string `json:"signing_key_b64"`
	DecryptKeyB64 string `json:"decrypt_key_b64"`

	CreatedAt time.Time `json:"created_at"`
}

// PublicKeyInfo stores a public key for verifying others' signatures.
type PublicKeyInfo struct {
	UserID     string `json:"user_id"`
	KeyVersion string `json:"key_version"`

	// Parsed key object (not serialized)
	PublicKey *ecdsa.PublicKey `json:"-"`

	// Serialized form for storage
	PublicKeySPKI string `json:"public_key_spki"`

	CreatedAt time.Time `json:"created_at"`
}

// ParsePublicKeySPKI parses a base64-encoded SPKI (SubjectPublicKeyInfo) public key.
func ParsePublicKeySPKI(spkiB64 string) (*ecdsa.PublicKey, error) {
	der, err := decodeBase64Flexible(spkiB64)
	if err != nil {
		return nil, fmt.Errorf("decode SPKI base64: %w", err)
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse SPKI: %w", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an ECDSA public key", ErrInvalidKeyFormat)
	}

	if ecdsaPub.Curve != elliptic.P256() {
		return nil, fmt.Errorf("%w: not a P-256 key", ErrInvalidKeyFormat)
	}

	return ecdsaPub, nil
}

// ParsePrivateKeyScalar parses a base64-encoded 32-byte P-256 private scalar.
// This is the format used by X/Twitter for storing private keys.
func ParsePrivateKeyScalar(scalarB64 string) (*ecdsa.PrivateKey, error) {
	scalar, err := decodeBase64Flexible(scalarB64)
	if err != nil {
		return nil, fmt.Errorf("decode scalar base64: %w", err)
	}

	if len(scalar) != 32 {
		return nil, fmt.Errorf("%w: scalar must be 32 bytes, got %d", ErrInvalidKeyFormat, len(scalar))
	}

	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(scalar)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(scalar)

	return priv, nil
}

// ParsePublicKeyUncompressed parses a 65-byte uncompressed P-256 public key.
// Format: 0x04 || X (32 bytes) || Y (32 bytes)
func ParsePublicKeyUncompressed(data []byte) (*ecdsa.PublicKey, error) {
	if len(data) != 65 {
		return nil, fmt.Errorf("%w: expected 65 bytes, got %d", ErrInvalidKeyFormat, len(data))
	}
	if data[0] != 0x04 {
		return nil, fmt.Errorf("%w: expected uncompressed point (0x04 prefix)", ErrInvalidKeyFormat)
	}

	curve := elliptic.P256()
	x := new(big.Int).SetBytes(data[1:33])
	y := new(big.Int).SetBytes(data[33:65])

	if !curve.IsOnCurve(x, y) {
		return nil, fmt.Errorf("%w: point is not on P-256 curve", ErrInvalidKeyFormat)
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

// EncodePublicKeySPKI encodes an ECDSA public key to base64 SPKI format.
func EncodePublicKeySPKI(pub *ecdsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("marshal SPKI: %w", err)
	}
	return base64.StdEncoding.EncodeToString(der), nil
}

// EncodePrivateKeyScalar encodes the private scalar as base64.
// The result is a 32-byte scalar in big-endian format, base64-encoded.
func EncodePrivateKeyScalar(priv *ecdsa.PrivateKey) string {
	scalar := priv.D.Bytes()
	// Pad to 32 bytes if necessary
	if len(scalar) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(scalar):], scalar)
		scalar = padded
	}
	return base64.StdEncoding.EncodeToString(scalar)
}

// EncodePublicKeyUncompressed encodes an ECDSA public key to uncompressed format.
// Returns: 0x04 || X (32 bytes) || Y (32 bytes) = 65 bytes
func EncodePublicKeyUncompressed(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.Curve == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	byteLen := (pub.Curve.Params().BitSize + 7) / 8
	out := make([]byte, 1+2*byteLen)
	out[0] = 0x04
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()
	copy(out[1+byteLen-len(xBytes):1+byteLen], xBytes)
	copy(out[1+2*byteLen-len(yBytes):], yBytes)
	return out
}

// LoadSigningKeyPair creates a SigningKeyPair from stored base64 scalar values.
// signingKeyB64: the signing key for signing messages (base64 32-byte scalar)
// decryptKeyB64: the decryption key for decrypting conversation keys (base64 32-byte scalar)
func LoadSigningKeyPair(userID, keyVersion, signingKeyB64, decryptKeyB64 string) (*SigningKeyPair, error) {
	var signingKey *ecdsa.PrivateKey
	var decryptKey *ecdsa.PrivateKey
	var err error

	if signingKeyB64 != "" {
		signingKey, err = ParsePrivateKeyScalar(signingKeyB64)
		if err != nil {
			return nil, fmt.Errorf("parse signing key: %w", err)
		}
	}

	if decryptKeyB64 != "" {
		decryptKey, err = ParsePrivateKeyScalar(decryptKeyB64)
		if err != nil {
			return nil, fmt.Errorf("parse decrypt key: %w", err)
		}
	}

	return &SigningKeyPair{
		UserID:        userID,
		KeyVersion:    keyVersion,
		SigningKey:    signingKey,
		DecryptKey:    decryptKey,
		SigningKeyB64: signingKeyB64,
		DecryptKeyB64: decryptKeyB64,
		CreatedAt:     time.Now(),
	}, nil
}

// LoadPublicKeyInfo creates a PublicKeyInfo from a stored base64 SPKI value.
func LoadPublicKeyInfo(userID, keyVersion, publicKeySPKI string) (*PublicKeyInfo, error) {
	pubKey, err := ParsePublicKeySPKI(publicKeySPKI)
	if err != nil {
		return nil, err
	}

	return &PublicKeyInfo{
		UserID:        userID,
		KeyVersion:    keyVersion,
		PublicKey:     pubKey,
		PublicKeySPKI: publicKeySPKI,
		CreatedAt:     time.Now(),
	}, nil
}
