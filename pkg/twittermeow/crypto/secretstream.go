package crypto

import (
	"errors"
	"fmt"

	"github.com/openziti/secretstream"
)

const (
	// SecretstreamChunkSize is the plaintext chunk size used for secretstream encryption.
	SecretstreamChunkSize = 1024

	// SecretstreamABytes is the authentication tag overhead per chunk (16 bytes MAC + 1 byte tag).
	SecretstreamABytes = 17

	// SecretstreamHeaderBytes is the header size for secretstream.
	SecretstreamHeaderBytes = 24
)

// SecretstreamEncrypt encrypts plaintext using XChaCha20-Poly1305 secretstream.
// The plaintext is split into 1024-byte chunks, each chunk is encrypted with an auth tag.
// Returns: header (24 bytes) || encrypted_chunks
func SecretstreamEncrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != secretboxKeySize {
		return nil, fmt.Errorf("secretstream key must be %d bytes", secretboxKeySize)
	}

	enc, header, err := secretstream.NewEncryptor(key)
	if err != nil {
		return nil, fmt.Errorf("create encryptor: %w", err)
	}

	// Pre-allocate output buffer: header + plaintext + auth tags for each chunk
	numChunks := (len(plaintext) + SecretstreamChunkSize - 1) / SecretstreamChunkSize
	if numChunks == 0 {
		numChunks = 1 // Handle empty plaintext
	}
	output := make([]byte, 0, len(header)+len(plaintext)+numChunks*SecretstreamABytes)
	output = append(output, header...)

	for offset := 0; offset < len(plaintext); offset += SecretstreamChunkSize {
		end := offset + SecretstreamChunkSize
		isLast := end >= len(plaintext)
		if isLast {
			end = len(plaintext)
		}

		tag := byte(secretstream.TagMessage)
		if isLast {
			tag = byte(secretstream.TagFinal)
		}

		chunk, err := enc.Push(plaintext[offset:end], tag)
		if err != nil {
			return nil, fmt.Errorf("encrypt chunk at offset %d: %w", offset, err)
		}
		output = append(output, chunk...)
	}

	// Handle empty plaintext case - push an empty final chunk
	if len(plaintext) == 0 {
		chunk, err := enc.Push(nil, byte(secretstream.TagFinal))
		if err != nil {
			return nil, fmt.Errorf("encrypt empty chunk: %w", err)
		}
		output = append(output, chunk...)
	}

	return output, nil
}

// SecretstreamDecrypt decrypts secretstream ciphertext (header || encrypted_chunks).
// Returns the decrypted plaintext.
func SecretstreamDecrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != secretboxKeySize {
		return nil, fmt.Errorf("secretstream key must be %d bytes", secretboxKeySize)
	}

	if len(ciphertext) < SecretstreamHeaderBytes {
		return nil, errors.New("secretstream ciphertext too short for header")
	}

	header := ciphertext[:SecretstreamHeaderBytes]
	dec, err := secretstream.NewDecryptor(key, header)
	if err != nil {
		return nil, fmt.Errorf("create decryptor: %w", err)
	}

	var plaintext []byte
	const encryptedChunkSize = SecretstreamChunkSize + SecretstreamABytes

	offset := SecretstreamHeaderBytes // skip header
	for offset < len(ciphertext) {
		end := offset + encryptedChunkSize
		if end > len(ciphertext) {
			end = len(ciphertext)
		}

		chunk, tag, err := dec.Pull(ciphertext[offset:end])
		if err != nil {
			return nil, fmt.Errorf("decrypt chunk at offset %d: %w", offset, err)
		}
		plaintext = append(plaintext, chunk...)

		if tag == secretstream.TagFinal {
			break
		}
		offset = end
	}

	return plaintext, nil
}
