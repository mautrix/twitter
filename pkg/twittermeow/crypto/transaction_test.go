package crypto_test

import (
	"log"
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
)

func TestXClientTransactionId(t *testing.T) {
	verificationToken := ""
	v, err := crypto.SignTransaction(verificationToken, "/1.1/jot/client_event.json", "POST")
	if err != nil {
		log.Fatalf("failed to sign client transaction id: %s", err.Error())
	}
	log.Println(v)
}
