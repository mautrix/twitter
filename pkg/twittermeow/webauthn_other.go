//go:build !windows

package twittermeow

import (
	"context"
	"fmt"
)

func platformWebAuthnGetAssertion(context.Context, webAuthnChallenge, []byte) (*webAuthnAssertion, error) {
	return nil, fmt.Errorf("native security-key/passkey login requires Windows WebAuthn support")
}
