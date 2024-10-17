package connector

import (
	"crypto/ecdh"
	"crypto/rand"

	"go.mau.fi/util/exerrors"
	"go.mau.fi/util/random"
)

type UserLoginMetadata struct {
	Cookies  string    `json:"cookies"`
	PushKeys *PushKeys `json:"push_keys,omitempty"`
}

type PushKeys struct {
	P256DH  []byte `json:"p256dh"`
	Auth    []byte `json:"auth"`
	Private []byte `json:"private"`
}

func (m *UserLoginMetadata) GeneratePushKeys() {
	privateKey := exerrors.Must(ecdh.P256().GenerateKey(rand.Reader))
	m.PushKeys = &PushKeys{
		P256DH:  privateKey.Public().(*ecdh.PublicKey).Bytes(),
		Auth:    random.Bytes(16),
		Private: privateKey.Bytes(),
	}
}
