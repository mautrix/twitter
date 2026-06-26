package twittermeow

import "testing"

func TestCastleRequestTokenShape(t *testing.T) {
	token, err := createCastleRequestToken()
	if err != nil {
		t.Fatalf("createCastleRequestToken() error = %v", err)
	}
	if !looksLikeCastleToken(token) {
		t.Fatalf("createCastleRequestToken() returned suspicious token length=%d", len(token))
	}
	t.Logf("Castle request token length=%d", len(token))
	second, err := createCastleRequestToken()
	if err != nil {
		t.Fatalf("second createCastleRequestToken() error = %v", err)
	}
	if token == second {
		t.Fatal("createCastleRequestToken() returned identical tokens")
	}
}
