package connector

import (
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func TestLatestPublicKeyWithJuiceboxTokens(t *testing.T) {
	makeKey := func(version string, hasTokens bool) response.PublicKeyWithTokenMap {
		key := response.PublicKeyWithTokenMap{}
		key.PublicKeyWithMetadata.Version = version
		if hasTokens {
			key.TokenMap.TokenMap = []response.KeyStoreTokenEntry{{}}
		}
		return key
	}

	data := &response.GetPublicKeysResponse{}
	user := response.UserResultsWithPublicKeys{}
	user.Result.GetPublicKeys.PublicKeysWithTokenMap = []response.PublicKeyWithTokenMap{
		makeKey("1768512248739", true),
		makeKey("1768512270000", false),
		makeKey("1768512266300", true),
	}
	data.Data.UserResultsByRestIDs = []response.UserResultsWithPublicKeys{user}

	key, ok := latestPublicKeyWithJuiceboxTokens(data)
	if !ok {
		t.Fatal("latestPublicKeyWithJuiceboxTokens() found no key")
	}
	if got, want := key.PublicKeyWithMetadata.Version, "1768512266300"; got != want {
		t.Fatalf("latestPublicKeyWithJuiceboxTokens() version = %q, want %q", got, want)
	}
}

func TestLatestPublicKeyWithJuiceboxTokensHandlesMissingData(t *testing.T) {
	for _, data := range []*response.GetPublicKeysResponse{nil, {}} {
		if key, ok := latestPublicKeyWithJuiceboxTokens(data); ok {
			t.Fatalf("latestPublicKeyWithJuiceboxTokens() = %#v, true; want no key", key)
		}
	}
}
