package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	secretboxKeySize   = 32
	secretboxNonceSize = 24
)

// SecretboxEncrypt encrypts plaintext with a random 24-byte nonce using XSalsa20-Poly1305
// (libsodium secretbox). The returned slice is nonce||ciphertext.
func SecretboxEncrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != secretboxKeySize {
		return nil, fmt.Errorf("secretbox key must be %d bytes", secretboxKeySize)
	}
	var nonce [secretboxNonceSize]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	var k [secretboxKeySize]byte
	copy(k[:], key)

	ct := secretbox.Seal(nil, plaintext, &nonce, &k)
	out := make([]byte, 0, len(nonce)+len(ct))
	out = append(out, nonce[:]...)
	out = append(out, ct...)
	return out, nil
}

// SecretboxDecrypt opens a libsodium secretbox message where the nonce is
// prefixed to the ciphertext (nonce||ciphertext). Returns an error on failure.
func SecretboxDecrypt(nonceCiphertext, key []byte) ([]byte, error) {
	if len(key) != secretboxKeySize {
		return nil, fmt.Errorf("secretbox key must be %d bytes", secretboxKeySize)
	}
	if len(nonceCiphertext) < secretboxNonceSize+secretbox.Overhead {
		return nil, errors.New("secretbox payload too short")
	}
	var nonce [secretboxNonceSize]byte
	copy(nonce[:], nonceCiphertext[:secretboxNonceSize])
	ciphertext := nonceCiphertext[secretboxNonceSize:]

	var k [secretboxKeySize]byte
	copy(k[:], key)

	plaintext, ok := secretbox.Open(nil, ciphertext, &nonce, &k)
	if !ok {
		return nil, errors.New("secretbox decrypt failed")
	}
	return plaintext, nil
}

// UnwrapConversationKey mirrors the WebCrypto flow in reverse-xchat:
//  1) keyB64 = base64(pub(65 bytes) || AES-GCM(ct||tag))
//  2) privScalarB64 = base64(32-byte P-256 private scalar)
//  3) derive shared = ECDH(priv, ephPub) (32-byte x coordinate)
//  4) KDF2-SHA256(shared || counter || ephPub) → 32 bytes
//     first 16 = AES key, last 16 = IV
//  5) AES-GCM decrypt ciphertext, expecting a 32-byte conversation key.
func UnwrapConversationKey(keyB64, privScalarB64 string) ([]byte, error) {
	blob, err := decodeBase64Flexible(keyB64)
	if err != nil {
		return nil, fmt.Errorf("decode key blob: %w", err)
	}
	if len(blob) < 65+16 {
		return nil, fmt.Errorf("key blob too short: %d", len(blob))
	}
	ephPub := blob[:65]
	cipherAndTag := blob[65:]

	privScalar, err := decodeBase64Flexible(privScalarB64)
	if err != nil {
		return nil, fmt.Errorf("decode private scalar: %w", err)
	}
	if len(privScalar) != 32 {
		return nil, fmt.Errorf("private scalar must be 32 bytes, got %d", len(privScalar))
	}

	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, ephPub)
	if x == nil || y == nil {
		return nil, errors.New("invalid ephemeral public key")
	}

	// Derive shared secret (x coordinate, 32 bytes big-endian).
	sx, _ := curve.ScalarMult(x, y, privScalar)
	shared := sx.Bytes()
	if len(shared) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(shared):], shared)
		shared = padded
	}

	kdfOut, err := kdf2SHA256(shared, ephPub, 32)
	if err != nil {
		return nil, fmt.Errorf("kdf: %w", err)
	}
	aesKey := kdfOut[:16]
	iv := kdfOut[16:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("aes init: %w", err)
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return nil, fmt.Errorf("gcm init: %w", err)
	}

	plaintext, err := gcm.Open(nil, iv, cipherAndTag, nil)
	if err != nil {
		return nil, fmt.Errorf("aes-gcm decrypt: %w", err)
	}
	if len(plaintext) != 32 {
		return nil, fmt.Errorf("unexpected conversation key length: %d", len(plaintext))
	}
	return plaintext, nil
}

func kdf2SHA256(shared, other []byte, length int) ([]byte, error) {
	var (
		counter uint32 = 1
		out             = make([]byte, 0, length)
	)
	for len(out) < length {
		counterBytes := []byte{0, 0, 0, 0}
		counterBytes[0] = byte(counter >> 24)
		counterBytes[1] = byte(counter >> 16)
		counterBytes[2] = byte(counter >> 8)
		counterBytes[3] = byte(counter)

		h := sha256.New()
		if _, err := h.Write(shared); err != nil {
			return nil, err
		}
		if _, err := h.Write(counterBytes); err != nil {
			return nil, err
		}
		if _, err := h.Write(other); err != nil {
			return nil, err
		}
		out = h.Sum(out)
		counter++
	}
	return out[:length], nil
}

// decodeBase64Flexible trims whitespace and tries standard and URL-safe base64
// decodings (with and without padding).
func decodeBase64Flexible(s string) ([]byte, error) {
	clean := strings.TrimSpace(s)
	clean = strings.ReplaceAll(clean, "\n", "")
	clean = strings.ReplaceAll(clean, "\r", "")
	clean = strings.ReplaceAll(clean, " ", "")

	encodings := []*base64.Encoding{
		base64.StdEncoding,
		base64.URLEncoding,
		base64.RawStdEncoding,
		base64.RawURLEncoding,
	}
	var lastErr error
	for _, enc := range encodings {
		if dec, err := enc.DecodeString(clean); err == nil {
			return dec, nil
		} else {
			lastErr = err
		}
	}
	return nil, lastErr
}
