package twittermeow

import "strings"

func looksLikeCastleToken(token string) bool {
	return len(token) > 500 && !strings.Contains(token, "=")
}
