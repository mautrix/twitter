package crypto

import (
	"fmt"
	"io"

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
func SecretstreamDecrypt(ciphertext io.Reader, key []byte, writeTo io.Writer) (int, error) {
	if len(key) != secretboxKeySize {
		return 0, fmt.Errorf("secretstream key must be %d bytes", secretboxKeySize)
	}

	header := make([]byte, SecretstreamHeaderBytes)
	_, err := io.ReadFull(ciphertext, header)
	if err != nil {
		return 0, fmt.Errorf("read secretstream header: %w", err)
	}

	dec, err := secretstream.NewDecryptor(key, header)
	if err != nil {
		return 0, fmt.Errorf("create decryptor: %w", err)
	}

	const encryptedChunkSize = SecretstreamChunkSize + SecretstreamABytes
	buf := make([]byte, encryptedChunkSize)
	var size int
	for {
		// This will return ErrUnexpectedEOF on the last chunk, but we ignore it because n won't be 0
		n, err := io.ReadFull(ciphertext, buf)
		if err != nil && n == 0 {
			return 0, fmt.Errorf("read chunk at %d: %w", size, err)
		}
		chunk, tag, err := dec.Pull(buf[:n])
		if err != nil {
			return 0, fmt.Errorf("decrypt chunk at %d: %w", size, err)
		}
		_, err = writeTo.Write(chunk)
		if err != nil {
			return 0, fmt.Errorf("write decrypted chunk at %d: %w", size, err)
		}
		size += n
		if tag == secretstream.TagFinal {
			break
		}
	}
	return size, nil
}
