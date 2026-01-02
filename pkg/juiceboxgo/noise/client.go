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

// Package noise implements the Noise NK protocol for Juicebox hardware realms.
// Protocol: Noise_NK_25519_ChaChaPoly_BLAKE2s
package noise

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

const (
	protocolName = "Noise_NK_25519_ChaChaPoly_BLAKE2s"
	hashLen      = 32
)

var (
	ErrHandshakeFailed = errors.New("noise handshake failed")
	ErrDecryption      = errors.New("noise decryption failed")
	ErrEncryption      = errors.New("noise encryption failed")
)

// HandshakeRequest is sent from client to server during handshake.
type HandshakeRequest struct {
	ClientEphemeralPublic []byte `cbor:"client_ephemeral_public"`
	PayloadCiphertext     []byte `cbor:"payload_ciphertext"`
}

// HandshakeResponse is sent from server to client during handshake.
type HandshakeResponse struct {
	ServerEphemeralPublic []byte `cbor:"server_ephemeral_public"`
	PayloadCiphertext     []byte `cbor:"payload_ciphertext"`
}

// Transport handles encrypted communication after handshake.
type Transport struct {
	inbound  *cipherState
	outbound *cipherState
}

// Decrypt decrypts a message from the server.
func (t *Transport) Decrypt(ciphertext []byte) ([]byte, error) {
	return t.inbound.decryptWithAD(ciphertext, nil)
}

// Encrypt encrypts a message to the server.
func (t *Transport) Encrypt(plaintext []byte) ([]byte, error) {
	return t.outbound.encryptWithAD(plaintext, nil)
}

// Handshake holds the state during an NK handshake.
type Handshake struct {
	clientEphemeralSecret [32]byte
	h                     [hashLen]byte
	ck                    [hashLen]byte
}

// Start initiates an NK handshake with the server.
// serverStaticPublic is the server's known static public key (32 bytes).
// payloadPlaintext is the optional payload to send (may be empty, but won't have forward secrecy).
func Start(serverStaticPublic []byte, payloadPlaintext []byte) (*Handshake, *HandshakeRequest, error) {
	if len(serverStaticPublic) != 32 {
		return nil, nil, ErrHandshakeFailed
	}

	// Generate client ephemeral key pair
	var clientEphemeralSecret, clientEphemeralPublic [32]byte
	if _, err := rand.Read(clientEphemeralSecret[:]); err != nil {
		return nil, nil, err
	}
	curve25519.ScalarBaseMult(&clientEphemeralPublic, &clientEphemeralSecret)

	// Initialize h and ck
	h := blake2s.Sum256([]byte(protocolName))
	ck := h

	// mix_hash(prologue) - empty prologue
	// mix_hash(serverStaticPublic)
	h = mixHash(h, serverStaticPublic)
	// mix_hash(clientEphemeralPublic)
	h = mixHash(h, clientEphemeralPublic[:])

	// DH(clientEphemeral, serverStatic)
	var serverPub [32]byte
	copy(serverPub[:], serverStaticPublic)
	sharedSecret, err := curve25519.X25519(clientEphemeralSecret[:], serverPub[:])
	if err != nil {
		return nil, nil, err
	}

	// mix_key(sharedSecret)
	var cipher *cipherState
	ck, cipher = mixKey(ck, sharedSecret)

	// Encrypt payload with AD = h
	ciphertext, err := cipher.encryptWithAD(payloadPlaintext, h[:])
	if err != nil {
		return nil, nil, err
	}

	// mix_hash(ciphertext)
	h = mixHash(h, ciphertext)

	handshake := &Handshake{
		clientEphemeralSecret: clientEphemeralSecret,
		h:                     h,
		ck:                    ck,
	}

	request := &HandshakeRequest{
		ClientEphemeralPublic: clientEphemeralPublic[:],
		PayloadCiphertext:     ciphertext,
	}

	return handshake, request, nil
}

// Finish completes the handshake with the server's response.
// Returns the Transport for further communication and the decrypted response payload.
func (hs *Handshake) Finish(response *HandshakeResponse) (*Transport, []byte, error) {
	if len(response.ServerEphemeralPublic) != 32 {
		return nil, nil, ErrHandshakeFailed
	}

	h := hs.h
	ck := hs.ck

	// mix_hash(serverEphemeralPublic)
	h = mixHash(h, response.ServerEphemeralPublic)

	// DH(clientEphemeral, serverEphemeral)
	var serverPub [32]byte
	copy(serverPub[:], response.ServerEphemeralPublic)
	sharedSecret, err := curve25519.X25519(hs.clientEphemeralSecret[:], serverPub[:])
	if err != nil {
		return nil, nil, err
	}

	// mix_key(sharedSecret)
	var cipher *cipherState
	ck, cipher = mixKey(ck, sharedSecret)

	// Decrypt payload with AD = h
	plaintext, err := cipher.decryptWithAD(response.PayloadCiphertext, h[:])
	if err != nil {
		return nil, nil, ErrDecryption
	}

	// Split to get Transport (client role)
	transport := split(ck, roleClient)

	return transport, plaintext, nil
}

// cipherState corresponds to Noise CipherState.
type cipherState struct {
	key   [32]byte
	nonce uint64
}

func newCipherState(key [32]byte) *cipherState {
	return &cipherState{
		key:   key,
		nonce: 0,
	}
}

func (cs *cipherState) nextNonce() []byte {
	// Noise uses 4 bytes padding + 8 bytes little-endian counter
	var nonce [12]byte
	binary.LittleEndian.PutUint64(nonce[4:], cs.nonce)
	cs.nonce++
	return nonce[:]
}

func (cs *cipherState) encryptWithAD(plaintext, ad []byte) ([]byte, error) {
	cipher, _ := chacha20poly1305.New(cs.key[:])
	nonce := cs.nextNonce()
	return cipher.Seal(nil, nonce, plaintext, ad), nil
}

func (cs *cipherState) decryptWithAD(ciphertext, ad []byte) ([]byte, error) {
	cipher, _ := chacha20poly1305.New(cs.key[:])
	nonce := cs.nextNonce()
	return cipher.Open(nil, nonce, ciphertext, ad)
}

type role int

const (
	roleClient role = iota
	roleServer
)

func mixHash(h [hashLen]byte, data []byte) [hashLen]byte {
	hasher, _ := blake2s.New256(nil)
	hasher.Write(h[:])
	hasher.Write(data)
	var result [hashLen]byte
	copy(result[:], hasher.Sum(nil))
	return result
}

func mixKey(ck [hashLen]byte, dh []byte) ([hashLen]byte, *cipherState) {
	ck1, k := hkdfPair(ck[:], dh)
	return ck1, newCipherState(k)
}

func split(ck [hashLen]byte, r role) *Transport {
	k1, k2 := hkdfPair(ck[:], nil)
	switch r {
	case roleServer:
		return &Transport{
			inbound:  newCipherState(k1),
			outbound: newCipherState(k2),
		}
	default: // roleClient
		return &Transport{
			inbound:  newCipherState(k2),
			outbound: newCipherState(k1),
		}
	}
}

func hkdfPair(salt, ikm []byte) ([hashLen]byte, [hashLen]byte) {
	reader := hkdf.New(newBlake2sHash, ikm, salt, nil)
	var okm [hashLen * 2]byte
	io.ReadFull(reader, okm[:])

	var first, second [hashLen]byte
	copy(first[:], okm[:hashLen])
	copy(second[:], okm[hashLen:])
	return first, second
}

// newBlake2sHash returns a new blake2s-256 hash for use with hkdf.
func newBlake2sHash() hash.Hash {
	h, _ := blake2s.New256(nil)
	return h
}
