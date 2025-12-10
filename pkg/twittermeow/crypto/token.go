package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// GenerateConversationToken creates a JWT for SendMessageMutation.
// The token uses HS256 (HMAC-SHA256) signing.
// secretKey should be either the conversation key or signing private key (32 bytes).
func GenerateConversationToken(requestingUser, recipient string, validSinceMSec int64, secretKey []byte) string {
	header := `{"alg":"HS256"}`
	payload := fmt.Sprintf(`{"requestingUser": "%s", "recipient": "%s", "validSinceMSec": "%d"}`,
		requestingUser, recipient, validSinceMSec)

	headerB64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadB64 := base64.RawURLEncoding.EncodeToString([]byte(payload))

	signingInput := headerB64 + "." + payloadB64

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(signingInput))
	signature := mac.Sum(nil)
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + signatureB64
}

// GetRecipientFromConversationID extracts the other user's ID from a conversation ID.
// For "user1:user2", returns the ID that is NOT the requesting user.
// For group chats or invalid format, returns the requesting user as fallback.
func GetRecipientFromConversationID(conversationID, requestingUser string) string {
	parts := strings.Split(conversationID, ":")
	if len(parts) != 2 {
		return requestingUser // fallback for group chats or invalid format
	}
	if parts[0] == requestingUser {
		return parts[1]
	}
	return parts[0]
}
