package connector

import (
	"context"
	"errors"
	"strings"
	"testing"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo"
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

func makeJuiceboxKeyCandidate(version, configJSON, realmID, authToken string) response.PublicKeyWithTokenMap {
	return response.PublicKeyWithTokenMap{
		PublicKeyWithMetadata: response.XChatPublicKeyWithMeta{Version: version},
		TokenMap: response.KeyStoreTokenMap{
			KeyStoreTokenMapJSON: configJSON,
			TokenMap: []response.KeyStoreTokenEntry{{
				Key: realmID,
				Value: response.KeyStoreToken{
					Token: authToken,
				},
			}},
		},
	}
}

func makePublicKeysResponse(candidates ...response.PublicKeyWithTokenMap) *response.GetPublicKeysResponse {
	data := &response.GetPublicKeysResponse{}
	data.Data.UserResultsByRestIDs = []response.UserResultsWithPublicKeys{{}}
	data.Data.UserResultsByRestIDs[0].Result.GetPublicKeys.PublicKeysWithTokenMap = candidates
	return data
}

func TestSelectRegisteredJuiceboxKeyChoosesNewestRegisteredCandidate(t *testing.T) {
	data := makePublicKeysResponse(
		makeJuiceboxKeyCandidate("1768512248739", "oldest-config", "realm-a", "token-a"),
		makeJuiceboxKeyCandidate("1768512270000", "newest-config", "realm-c", "token-c"),
		makeJuiceboxKeyCandidate("1768512266300", "middle-config", "realm-b", "token-b"),
	)
	var checked []string
	selected, hasCandidates, err := selectRegisteredJuiceboxKey(
		context.Background(),
		data,
		func(_ context.Context, configJSON string, _ map[string]string) error {
			checked = append(checked, configJSON)
			return nil
		},
	)
	if err != nil || !hasCandidates {
		t.Fatalf("selectRegisteredJuiceboxKey() = (_, %t, %v), want registered candidate", hasCandidates, err)
	}
	if got, want := selected.PublicKeyWithMetadata.Version, "1768512270000"; got != want {
		t.Fatalf("selected version = %q, want %q", got, want)
	}
	if got, want := strings.Join(checked, ","), "newest-config"; got != want {
		t.Fatalf("checked configs = %q, want %q", got, want)
	}
}

func TestSelectRegisteredJuiceboxKeySkipsUnregisteredNewerCandidate(t *testing.T) {
	data := makePublicKeysResponse(
		makeJuiceboxKeyCandidate("1768512248739", "oldest-config", "realm-a", "token-a"),
		makeJuiceboxKeyCandidate("1768512270000", "newest-config", "realm-c", "token-c"),
		makeJuiceboxKeyCandidate("1768512266300", "current-config", "realm-b", "token-b"),
	)
	var checked []string
	selected, hasCandidates, err := selectRegisteredJuiceboxKey(
		context.Background(),
		data,
		func(_ context.Context, configJSON string, _ map[string]string) error {
			checked = append(checked, configJSON)
			if configJSON == "newest-config" {
				return juiceboxgo.ErrNotRegistered
			}
			return nil
		},
	)
	if err != nil || !hasCandidates {
		t.Fatalf("selectRegisteredJuiceboxKey() = (_, %t, %v), want registered candidate", hasCandidates, err)
	}
	if got, want := selected.PublicKeyWithMetadata.Version, "1768512266300"; got != want {
		t.Fatalf("selected version = %q, want %q", got, want)
	}
	if got, want := strings.Join(checked, ","), "newest-config,current-config"; got != want {
		t.Fatalf("checked configs = %q, want %q", got, want)
	}
}

func TestSelectRegisteredJuiceboxKeyStopsOnOperationalError(t *testing.T) {
	data := makePublicKeysResponse(
		makeJuiceboxKeyCandidate("1768512248739", "oldest-config", "realm-a", "token-a"),
		makeJuiceboxKeyCandidate("1768512270000", "newest-config", "realm-b", "token-b"),
	)
	checks := 0
	_, hasCandidates, err := selectRegisteredJuiceboxKey(
		context.Background(),
		data,
		func(context.Context, string, map[string]string) error {
			checks++
			return juiceboxgo.ErrTransient
		},
	)
	if !hasCandidates {
		t.Fatal("selectRegisteredJuiceboxKey() hasCandidates = false, want true")
	}
	if !errors.Is(err, juiceboxgo.ErrTransient) {
		t.Fatalf("selectRegisteredJuiceboxKey() error = %v, want ErrTransient", err)
	}
	if checks != 1 {
		t.Fatalf("registration checks = %d, want 1", checks)
	}
}

func TestSelectRegisteredJuiceboxKeyReportsAllUnregistered(t *testing.T) {
	data := makePublicKeysResponse(
		makeJuiceboxKeyCandidate("1768512248739", "oldest-config", "realm-a", "token-a"),
		makeJuiceboxKeyCandidate("1768512270000", "newest-config", "realm-b", "token-b"),
	)
	checks := 0
	_, hasCandidates, err := selectRegisteredJuiceboxKey(
		context.Background(),
		data,
		func(context.Context, string, map[string]string) error {
			checks++
			return juiceboxgo.ErrNotRegistered
		},
	)
	if !hasCandidates || !errors.Is(err, juiceboxgo.ErrNotRegistered) {
		t.Fatalf("selectRegisteredJuiceboxKey() = (_, %t, %v), want candidates and ErrNotRegistered", hasCandidates, err)
	}
	if checks != 2 {
		t.Fatalf("registration checks = %d, want 2", checks)
	}
}

func TestSelectRegisteredJuiceboxKeyIgnoresIncompleteTokenMaps(t *testing.T) {
	data := makePublicKeysResponse(
		makeJuiceboxKeyCandidate("blank-config", " ", "realm-a", "token-a"),
		makeJuiceboxKeyCandidate("blank-realm", "config-a", " ", "token-a"),
		makeJuiceboxKeyCandidate("blank-token", "config-b", "realm-b", " "),
	)
	checks := 0
	_, hasCandidates, err := selectRegisteredJuiceboxKey(
		context.Background(),
		data,
		func(context.Context, string, map[string]string) error {
			checks++
			return nil
		},
	)
	if err != nil || hasCandidates {
		t.Fatalf("selectRegisteredJuiceboxKey() = (_, %t, %v), want no candidates", hasCandidates, err)
	}
	if checks != 0 {
		t.Fatalf("registration checks = %d, want 0", checks)
	}
}
