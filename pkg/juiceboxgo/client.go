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

package juiceboxgo

import (
	"context"
	"crypto/ed25519"
	crand "crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"net/http"
	"sync"

	"github.com/bwesterb/go-ristretto"
	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/crypto"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/oprf"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/pin"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/realm"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/requests"
	"go.mau.fi/mautrix-twitter/pkg/juiceboxgo/secretsharing"
)

// Client is used to recover PIN-protected secrets from Juicebox realms.
type Client struct {
	config     *Configuration
	authTokens map[string]AuthToken
	httpClient *http.Client
	logger     zerolog.Logger

	realmClients map[RealmID]*realm.Client
}

// NewClient creates a new Juicebox client.
func NewClient(config *Configuration, authTokens map[string]AuthToken, httpClient *http.Client, logger zerolog.Logger) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	normalizedConfig, err := NormalizeConfiguration(config)
	if err != nil {
		return nil, err
	}

	authProvider := func(realmID RealmID) (AuthToken, bool) {
		token, ok := authTokens[realmID.String()]
		return token, ok
	}

	realmClients := make(map[RealmID]*realm.Client)
	for _, r := range normalizedConfig.Realms {
		realmClients[r.ID] = realm.NewClient(r, httpClient, authProvider, logger)
	}

	return &Client{
		config:       normalizedConfig,
		authTokens:   authTokens,
		httpClient:   httpClient,
		logger:       logger,
		realmClients: realmClients,
	}, nil
}

// Close releases resources.
func (c *Client) Close() {
	// No persistent resources to clean up in pure Go implementation
}

// Register stores a PIN-protected secret in Juicebox realms.
func (c *Client) Register(ctx context.Context, pinBytes Pin, secret Secret, userInfo UserInfo, numGuesses uint16) error {
	threshold := int(c.config.RegisterThreshold)
	if threshold <= 0 || threshold > len(c.config.Realms) {
		return ErrAssertion
	}
	if len(c.config.Realms) == 0 {
		return ErrAssertion
	}

	c.logger.Debug().Msg("Starting registration phase 1")
	if err := c.registerPhase1(ctx, threshold); err != nil {
		return err
	}
	c.logger.Debug().Msg("Registration phase 1 complete")

	var version RegistrationVersion
	if _, err := crand.Read(version[:]); err != nil {
		return ErrAssertion
	}

	hashResult := pin.HashPIN(pinBytes, pin.HashingMode(c.config.PinHashingMode), [16]byte(version), userInfo)

	oprfPrivateKey, err := randomScalar()
	if err != nil {
		return ErrAssertion
	}
	oprfPrivateKeyShares, err := createScalarShares(oprfPrivateKey, c.config.RecoverThreshold, uint32(len(c.config.Realms)))
	if err != nil {
		return ErrAssertion
	}

	encryptionKeyScalar, err := randomScalar()
	if err != nil {
		return ErrAssertion
	}
	encryptionKeyScalarShares, err := createScalarShares(encryptionKeyScalar, c.config.RecoverThreshold, uint32(len(c.config.Realms)))
	if err != nil {
		return ErrAssertion
	}

	oprfOutput := unobliviousEvaluate(oprfPrivateKey, hashResult.AccessKey[:])
	unlockKey, unlockKeyCommitment := crypto.DeriveUnlockKeyAndCommitment(oprfOutput)

	encryptionKey := crypto.DeriveEncryptionKey(hashResult.EncryptionKeySeed, &encryptionKeyScalar)
	encryptedSecret, err := crypto.EncryptSecret(secret, encryptionKey)
	if err != nil {
		return err
	}

	verifyingKey, signingKey, err := ed25519.GenerateKey(crand.Reader)
	if err != nil {
		return ErrAssertion
	}

	c.logger.Debug().Msg("Starting registration phase 2")
	if err := c.registerPhase2(
		ctx,
		threshold,
		version,
		oprfPrivateKeyShares,
		encryptionKeyScalarShares,
		unlockKey,
		unlockKeyCommitment,
		EncryptedUserSecret(encryptedSecret),
		numGuesses,
		verifyingKey,
		signingKey,
	); err != nil {
		return err
	}
	c.logger.Debug().Msg("Registration phase 2 complete")

	return nil
}

func (c *Client) registerPhase1(ctx context.Context, threshold int) error {
	return c.runParallelRegisterRequestsWithThreshold(threshold, func(_ int, r Realm) error {
		client := c.realmClients[r.ID]
		resp, err := client.MakeRequest(ctx, &requests.SecretsRequest{Register1: true})
		if err != nil {
			return err
		}
		if resp.Register1 == nil || !resp.Register1.Ok {
			return ErrAssertion
		}
		return nil
	})
}

func (c *Client) runParallelRegisterRequestsWithThreshold(threshold int, run func(i int, r Realm) error) error {
	results := make(chan error, len(c.config.Realms))
	var wg sync.WaitGroup

	for i, r := range c.config.Realms {
		wg.Add(1)
		go func(i int, r Realm) {
			defer wg.Done()
			results <- run(i, r)
		}(i, r)
	}

	wg.Wait()
	close(results)

	return waitForThresholdSuccess(results, threshold)
}

func (c *Client) registerPhase2(
	ctx context.Context,
	threshold int,
	version RegistrationVersion,
	oprfPrivateKeyShares []ristretto.Scalar,
	encryptionKeyScalarShares []ristretto.Scalar,
	unlockKey [32]byte,
	unlockKeyCommitment UnlockKeyCommitment,
	encryptedSecret EncryptedUserSecret,
	numGuesses uint16,
	verifyingKey ed25519.PublicKey,
	signingKey ed25519.PrivateKey,
) error {
	if len(oprfPrivateKeyShares) != len(c.config.Realms) || len(encryptionKeyScalarShares) != len(c.config.Realms) {
		return ErrAssertion
	}
	if len(verifyingKey) != ed25519.PublicKeySize {
		return ErrAssertion
	}

	var verifyingKeyFixed [32]byte
	copy(verifyingKeyFixed[:], verifyingKey)

	return c.runParallelRegisterRequestsWithThreshold(threshold, func(i int, r Realm) error {
		var oprfPublicKey ristretto.Point
		oprfPublicKey.ScalarMultBase(&oprfPrivateKeyShares[i])
		oprfPublicKeyBytes := oprfPublicKey.Bytes()

		sigMsg := buildOprfSignatureMessage(r.ID[:], oprfPublicKeyBytes)
		signature := ed25519.Sign(signingKey, sigMsg)
		if len(signature) != ed25519.SignatureSize {
			return ErrAssertion
		}

		var signatureFixed [64]byte
		copy(signatureFixed[:], signature)

		unlockKeyTag := crypto.DeriveUnlockKeyTag(unlockKey, [16]byte(r.ID))
		encryptedSecretCommitment := crypto.DeriveEncryptedUserSecretCommitment(
			unlockKey, [16]byte(r.ID), &encryptionKeyScalarShares[i], encryptedSecret[:],
		)

		req := &requests.SecretsRequest{
			Register2: &requests.Register2Request{
				Version:                   version,
				OprfPrivateKey:            append([]byte(nil), oprfPrivateKeyShares[i].Bytes()...),
				OprfSignedPublicKey:       requests.OprfSignedPublicKey{PublicKey: append([]byte(nil), oprfPublicKeyBytes...), VerifyingKey: verifyingKeyFixed, Signature: signatureFixed},
				UnlockKeyCommitment:       unlockKeyCommitment,
				UnlockKeyTag:              unlockKeyTag,
				EncryptionKeyScalarShare:  append([]byte(nil), encryptionKeyScalarShares[i].Bytes()...),
				EncryptedSecret:           encryptedSecret,
				EncryptedSecretCommitment: encryptedSecretCommitment,
				Policy:                    requests.Policy{NumGuesses: numGuesses},
			},
		}

		client := c.realmClients[r.ID]
		resp, err := client.MakeRequest(ctx, req)
		if err != nil {
			return err
		}
		if resp.Register2 == nil || !resp.Register2.Ok {
			return ErrAssertion
		}
		return nil
	})
}

func waitForThresholdSuccess(results <-chan error, threshold int) error {
	successCount := 0
	var errs []error
	for err := range results {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		successCount++
	}
	if successCount >= threshold {
		return nil
	}
	if len(errs) > 0 {
		return bestError(errs, rankRegisterError)
	}
	return ErrAssertion
}

func bestError(errs []error, rank func(error) int) error {
	best := errs[0]
	bestRank := rank(best)
	for _, err := range errs[1:] {
		errRank := rank(err)
		if errRank < bestRank {
			best = err
			bestRank = errRank
		}
	}
	return best
}

func rankRegisterError(err error) int {
	switch {
	case errors.Is(err, ErrInvalidAuth):
		return 0
	case errors.Is(err, ErrUpgradeRequired):
		return 1
	case errors.Is(err, ErrRateLimitExceeded):
		return 2
	case errors.Is(err, ErrAssertion):
		return 3
	case errors.Is(err, ErrTransient):
		return 4
	default:
		return 5
	}
}

func rankRecoverError(err error) int {
	var recoverErr *RecoverError
	switch {
	case errors.As(err, &recoverErr):
		return 0
	case errors.Is(err, ErrNotRegistered):
		return 1
	case errors.Is(err, ErrInvalidAuth):
		return 2
	case errors.Is(err, ErrUpgradeRequired):
		return 3
	case errors.Is(err, ErrRateLimitExceeded):
		return 4
	case errors.Is(err, ErrAssertion):
		return 5
	case errors.Is(err, ErrTransient):
		return 6
	default:
		return 7
	}
}

func rankDeleteError(err error) int {
	switch {
	case errors.Is(err, ErrInvalidAuth):
		return 0
	case errors.Is(err, ErrUpgradeRequired):
		return 1
	case errors.Is(err, ErrRateLimitExceeded):
		return 2
	case errors.Is(err, ErrAssertion):
		return 3
	case errors.Is(err, ErrTransient):
		return 4
	default:
		return 5
	}
}

// Recover retrieves a PIN-protected secret from the realms.
func (c *Client) Recover(ctx context.Context, pinBytes Pin, userInfo UserInfo) (Secret, error) {
	return c.recoverWithConfiguration(ctx, pinBytes, userInfo, c.config, c.realmClients)
}

func (c *Client) recoverWithConfiguration(ctx context.Context, pinBytes Pin, userInfo UserInfo, config *Configuration, realmClients map[RealmID]*realm.Client) (Secret, error) {
	// Phase 1: Query version from realms
	c.logger.Debug().Msg("Starting recovery phase 1")
	version, realms, err := c.recoverPhase1(ctx, config, realmClients)
	if err != nil {
		return nil, err
	}
	c.logger.Debug().Str("version", string(version[:])).Int("realms", len(realms)).Msg("Phase 1 complete")

	// Hash PIN with Argon2
	hashResult := pin.HashPIN(pinBytes, pin.HashingMode(config.PinHashingMode), [16]byte(version), userInfo)

	// DEBUG: Log PIN hash results for comparison with Rust
	c.logger.Trace().
		Hex("access_key", hashResult.AccessKey[:]).
		Hex("encryption_key_seed", hashResult.EncryptionKeySeed[:]).
		Msg("PIN hash complete")

	// Phase 2: OPRF evaluation
	c.logger.Debug().Msg("Starting recovery phase 2")
	unlockKey, err := c.recoverPhase2(ctx, config, realmClients, version, realms, hashResult.AccessKey)
	if err != nil {
		return nil, err
	}
	c.logger.Debug().Msg("Phase 2 complete")

	// Phase 3: Retrieve encrypted secret
	c.logger.Debug().Msg("Starting recovery phase 3")
	secret, err := c.recoverPhase3(ctx, config, realmClients, version, realms, unlockKey, hashResult.EncryptionKeySeed)
	if err != nil {
		return nil, err
	}
	c.logger.Debug().Int("secret_len", len(secret)).Msg("Phase 3 complete - recovery successful")

	return secret, nil
}

// Delete removes any secret stored on the currently configured realms.
func (c *Client) Delete(ctx context.Context) error {
	results := make(chan error, len(c.config.Realms))
	var wg sync.WaitGroup

	for _, r := range c.config.Realms {
		wg.Add(1)
		go func(r Realm) {
			defer wg.Done()
			client := c.realmClients[r.ID]
			resp, err := client.MakeRequest(ctx, &requests.SecretsRequest{Delete: true})
			if err != nil {
				results <- err
				return
			}
			if resp.Delete == nil || !resp.Delete.Ok {
				results <- ErrAssertion
				return
			}
			results <- nil
		}(r)
	}

	wg.Wait()
	close(results)

	var errs []error
	for err := range results {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return bestError(errs, rankDeleteError)
	}
	return nil
}

// recoverPhase1 queries realms for the registration version.
func (c *Client) recoverPhase1(ctx context.Context, config *Configuration, realmClients map[RealmID]*realm.Client) (RegistrationVersion, []Realm, error) {
	type phase1Result struct {
		version RegistrationVersion
		realm   Realm
		err     error
	}

	results := make(chan phase1Result, len(config.Realms))
	var wg sync.WaitGroup

	for _, r := range config.Realms {
		wg.Add(1)
		go func(r Realm) {
			defer wg.Done()
			client := realmClients[r.ID]
			resp, err := client.MakeRequest(ctx, &requests.SecretsRequest{Recover1: true})
			if err != nil {
				results <- phase1Result{err: err, realm: r}
				return
			}
			if resp.Recover1 == nil {
				results <- phase1Result{err: ErrAssertion, realm: r}
				return
			}
			if resp.Recover1.NotRegistered {
				results <- phase1Result{err: ErrNotRegistered, realm: r}
				return
			}
			if resp.Recover1.NoGuesses {
				results <- phase1Result{err: ErrInvalidPin(0), realm: r}
				return
			}
			if resp.Recover1.Ok == nil {
				results <- phase1Result{err: ErrAssertion, realm: r}
				return
			}
			results <- phase1Result{version: resp.Recover1.Ok.Version, realm: r}
		}(r)
	}

	wg.Wait()
	close(results)

	// Collect successful results
	versionCounts := make(map[RegistrationVersion][]Realm)
	var errs []error
	successCount := 0

	for result := range results {
		if result.err != nil {
			errs = append(errs, result.err)
			continue
		}
		successCount++
		versionCounts[result.version] = append(versionCounts[result.version], result.realm)
	}

	// Find version with threshold agreement
	threshold := int(config.RecoverThreshold)
	if successCount < threshold {
		return RegistrationVersion{}, nil, bestError(errs, rankRecoverError)
	}
	for version, realms := range versionCounts {
		if len(realms) >= threshold {
			return version, realms, nil
		}
	}

	return RegistrationVersion{}, nil, ErrNotRegistered
}

// recoverPhase2 performs OPRF evaluation across realms.
func (c *Client) recoverPhase2(ctx context.Context, config *Configuration, realmClients map[RealmID]*realm.Client, version RegistrationVersion, realms []Realm, accessKey [32]byte) (UnlockKey, error) {
	// Start OPRF
	blindingFactor, blindedInput, err := oprf.Start(accessKey[:])
	if err != nil {
		return UnlockKey{}, err
	}

	// DEBUG: Log OPRF start results
	c.logger.Trace().
		Hex("oprf_input", accessKey[:]).
		Hex("blinded_input", blindedInput.Bytes()).
		Msg("OPRF started")

	type phase2Result struct {
		realm              Realm
		blindedResultShare secretsharing.PointShare
		commitment         UnlockKeyCommitment
		verifyingKey       [32]byte
		guessesRemaining   uint16
		err                error
	}

	results := make(chan phase2Result, len(realms))
	var wg sync.WaitGroup

	for _, r := range realms {
		wg.Add(1)
		go func(r Realm) {
			defer wg.Done()
			client := realmClients[r.ID]

			req := &requests.SecretsRequest{
				Recover2: &requests.Recover2Request{
					Version:          version,
					OprfBlindedInput: blindedInput.Bytes(),
				},
			}

			resp, err := client.MakeRequest(ctx, req)
			if err != nil {
				results <- phase2Result{err: err, realm: r}
				return
			}
			if resp.Recover2 == nil {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}
			if resp.Recover2.NotRegistered {
				results <- phase2Result{err: ErrNotRegistered, realm: r}
				return
			}
			if resp.Recover2.NoGuesses {
				results <- phase2Result{err: ErrInvalidPin(0), realm: r}
				return
			}
			if resp.Recover2.VersionMismatch {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}
			if resp.Recover2.Ok == nil {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}

			ok := resp.Recover2.Ok

			// Parse blinded result as point
			var blindedResultPoint ristretto.Point
			if len(ok.OprfBlindedResult) != 32 {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}
			var blindedResultBytes [32]byte
			copy(blindedResultBytes[:], ok.OprfBlindedResult)
			if !blindedResultPoint.SetBytes(&blindedResultBytes) {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}

			// Parse DLEQ proof
			var proofC, proofBetaZ ristretto.Scalar
			if len(ok.OprfProof.C) != 32 || len(ok.OprfProof.BetaZ) != 32 {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}
			var cBytes, betaZBytes [32]byte
			copy(cBytes[:], ok.OprfProof.C)
			copy(betaZBytes[:], ok.OprfProof.BetaZ)
			proofC.SetBytes(&cBytes)
			proofBetaZ.SetBytes(&betaZBytes)

			proof := &oprf.Proof{C: proofC, BetaZ: proofBetaZ}

			// Verify Ed25519 signature of OPRF public key
			signatureMsg := buildOprfSignatureMessage(r.ID[:], ok.OprfSignedPublicKey.PublicKey)
			if !ed25519.Verify(ok.OprfSignedPublicKey.VerifyingKey[:], signatureMsg, ok.OprfSignedPublicKey.Signature[:]) {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}

			// Parse public key
			var publicKey oprf.PublicKey
			if err := publicKey.SetBytes(ok.OprfSignedPublicKey.PublicKey); err != nil {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}

			// Verify DLEQ proof
			blindedOutput := &oprf.BlindedInput{}
			blindedOutput.SetBytes(ok.OprfBlindedResult)
			if err := oprf.VerifyProof(blindedInput, blindedOutput, &publicKey, proof); err != nil {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}

			// Get share index
			index, found := config.ShareIndex(r.ID)
			if !found {
				results <- phase2Result{err: ErrAssertion, realm: r}
				return
			}

			// DEBUG: Log realm response
			c.logger.Trace().
				Str("realm", r.ID.String()).
				Uint32("share_index", index).
				Hex("blinded_result", ok.OprfBlindedResult).
				Hex("commitment", ok.UnlockKeyCommitment[:]).
				Msg("Realm phase 2 response")

			results <- phase2Result{
				realm: r,
				blindedResultShare: secretsharing.PointShare{
					Index:  secretsharing.Index(index),
					Secret: blindedResultPoint,
				},
				commitment:       ok.UnlockKeyCommitment,
				verifyingKey:     ok.OprfSignedPublicKey.VerifyingKey,
				guessesRemaining: ok.NumGuesses - ok.GuessCount,
			}
		}(r)
	}

	wg.Wait()
	close(results)

	// Collect results grouped by (commitment, verifyingKey)
	type resultKey struct {
		commitment   UnlockKeyCommitment
		verifyingKey [32]byte
	}
	grouped := make(map[resultKey][]phase2Result)
	var errs []error
	successCount := 0

	for result := range results {
		if result.err != nil {
			errs = append(errs, result.err)
			continue
		}
		successCount++
		key := resultKey{commitment: result.commitment, verifyingKey: result.verifyingKey}
		grouped[key] = append(grouped[key], result)
	}

	// Find group with threshold agreement
	threshold := int(config.RecoverThreshold)
	if successCount < threshold {
		return UnlockKey{}, bestError(errs, rankRecoverError)
	}
	var selectedResults []phase2Result
	var unlockKeyCommitment UnlockKeyCommitment

	for key, groupResults := range grouped {
		if len(groupResults) >= threshold {
			selectedResults = groupResults
			unlockKeyCommitment = key.commitment
			break
		}
	}

	if len(selectedResults) < threshold {
		return UnlockKey{}, ErrAssertion
	}

	// Combine OPRF blinded results using Lagrange interpolation
	shares := make([]secretsharing.PointShare, len(selectedResults))
	for i, r := range selectedResults {
		shares[i] = r.blindedResultShare
	}

	combinedBlindedResult, err := secretsharing.RecoverPoint(shares)
	if err != nil {
		return UnlockKey{}, ErrAssertion
	}

	// Finalize OPRF
	blindedOutput := &oprf.BlindedOutput{}
	blindedOutput.SetBytes(combinedBlindedResult.Bytes())
	oprfOutput := oprf.Finalize(accessKey[:], blindingFactor, blindedOutput)

	// Derive unlock key
	unlockKeyRaw, ourCommitment := crypto.DeriveUnlockKeyAndCommitment(oprfOutput)

	// DEBUG: Log OPRF finalization results for comparison with Rust
	c.logger.Trace().
		Hex("combined_blinded_result", combinedBlindedResult.Bytes()).
		Hex("oprf_output", oprfOutput[:]).
		Hex("derived_commitment", ourCommitment[:]).
		Hex("expected_commitment", unlockKeyCommitment[:]).
		Hex("unlock_key", unlockKeyRaw[:]).
		Msg("OPRF finalized")

	// Verify commitment (constant time)
	if subtle.ConstantTimeCompare(unlockKeyCommitment[:], ourCommitment[:]) != 1 {
		// Wrong PIN
		minGuesses := uint16(0xFFFF)
		for _, result := range selectedResults {
			if result.guessesRemaining < minGuesses {
				minGuesses = result.guessesRemaining
			}
		}
		return UnlockKey{}, ErrInvalidPin(minGuesses)
	}

	return UnlockKey(unlockKeyRaw), nil
}

// recoverPhase3 retrieves the encrypted secret from realms.
func (c *Client) recoverPhase3(ctx context.Context, config *Configuration, realmClients map[RealmID]*realm.Client, version RegistrationVersion, realms []Realm, unlockKey UnlockKey, encryptionKeySeed [32]byte) (Secret, error) {
	type phase3Result struct {
		realm           Realm
		scalarShare     secretsharing.ScalarShare
		encryptedSecret EncryptedUserSecret
		commitment      EncryptedUserSecretCommitment
		skipped         bool
		err             error
	}

	results := make(chan phase3Result, len(realms))
	var wg sync.WaitGroup

	for _, r := range realms {
		wg.Add(1)
		go func(r Realm) {
			defer wg.Done()
			client := realmClients[r.ID]

			unlockKeyTag := crypto.DeriveUnlockKeyTag([32]byte(unlockKey), [16]byte(r.ID))

			req := &requests.SecretsRequest{
				Recover3: &requests.Recover3Request{
					Version:      version,
					UnlockKeyTag: unlockKeyTag,
				},
			}

			resp, err := client.MakeRequest(ctx, req)
			if err != nil {
				results <- phase3Result{err: err, realm: r}
				return
			}
			if resp.Recover3 == nil {
				results <- phase3Result{err: ErrAssertion, realm: r}
				return
			}
			if resp.Recover3.NotRegistered {
				results <- phase3Result{err: ErrNotRegistered, realm: r}
				return
			}
			if resp.Recover3.NoGuesses {
				results <- phase3Result{err: ErrInvalidPin(0), realm: r}
				return
			}
			if resp.Recover3.BadUnlockKeyTag != nil {
				results <- phase3Result{err: ErrInvalidPin(resp.Recover3.BadUnlockKeyTag.GuessesRemaining), realm: r}
				return
			}
			if resp.Recover3.VersionMismatch {
				results <- phase3Result{err: ErrAssertion, realm: r}
				return
			}
			if resp.Recover3.Ok == nil {
				results <- phase3Result{err: ErrAssertion, realm: r}
				return
			}

			ok := resp.Recover3.Ok

			// Parse scalar share
			if len(ok.EncryptionKeyScalarShare) != 32 {
				results <- phase3Result{err: ErrAssertion, realm: r}
				return
			}
			var scalarBytes [32]byte
			copy(scalarBytes[:], ok.EncryptionKeyScalarShare)
			var scalar ristretto.Scalar
			scalar.SetBytes(&scalarBytes)

			// Get share index
			index, found := config.ShareIndex(r.ID)
			if !found {
				results <- phase3Result{err: ErrAssertion, realm: r}
				return
			}

			// Verify commitment
			ourCommitment := crypto.DeriveEncryptedUserSecretCommitment(
				[32]byte(unlockKey), [16]byte(r.ID), &scalar, ok.EncryptedSecret[:],
			)
			if subtle.ConstantTimeCompare(ourCommitment[:], ok.EncryptedSecretCommitment[:]) != 1 {
				// Skip this share - commitment doesn't match
				results <- phase3Result{skipped: true, realm: r}
				return
			}

			results <- phase3Result{
				realm: r,
				scalarShare: secretsharing.ScalarShare{
					Index:  secretsharing.Index(index),
					Secret: scalar,
				},
				encryptedSecret: ok.EncryptedSecret,
				commitment:      ok.EncryptedSecretCommitment,
			}
		}(r)
	}

	wg.Wait()
	close(results)

	// Collect results grouped by encrypted secret
	grouped := make(map[EncryptedUserSecret][]phase3Result)
	var errs []error
	successCount := 0

	for result := range results {
		if result.err != nil {
			errs = append(errs, result.err)
			continue
		}
		if result.skipped {
			continue
		}
		successCount++
		grouped[result.encryptedSecret] = append(grouped[result.encryptedSecret], result)
	}

	// Find group with threshold agreement
	threshold := int(config.RecoverThreshold)
	if successCount < threshold {
		if len(errs) > 0 {
			return nil, bestError(errs, rankRecoverError)
		}
		return nil, ErrAssertion
	}
	var selectedResults []phase3Result
	var encryptedSecret EncryptedUserSecret

	for secret, groupResults := range grouped {
		if len(groupResults) >= threshold {
			selectedResults = groupResults
			encryptedSecret = secret
			break
		}
	}

	if len(selectedResults) < threshold {
		return nil, ErrAssertion
	}

	// Combine scalar shares using Lagrange interpolation
	shares := make([]secretsharing.ScalarShare, len(selectedResults))
	for i, r := range selectedResults {
		shares[i] = r.scalarShare
	}

	combinedScalar, err := secretsharing.RecoverScalar(shares)
	if err != nil {
		return nil, ErrAssertion
	}

	// Derive encryption key
	encryptionKey := crypto.DeriveEncryptionKey(encryptionKeySeed, &combinedScalar)

	// Decrypt secret
	secret, err := crypto.DecryptSecret(encryptedSecret[:], encryptionKey)
	if err != nil {
		return nil, ErrAssertion
	}

	return secret, nil
}

func randomScalar() (ristretto.Scalar, error) {
	var scalar ristretto.Scalar
	var randomBytes [64]byte
	if _, err := crand.Read(randomBytes[:]); err != nil {
		return scalar, err
	}
	scalar.SetReduced(&randomBytes)
	return scalar, nil
}

func createScalarShares(secret ristretto.Scalar, threshold, count uint32) ([]ristretto.Scalar, error) {
	if threshold == 0 || count == 0 || threshold > count {
		return nil, ErrAssertion
	}

	coefficients := make([]ristretto.Scalar, threshold-1)
	for i := range coefficients {
		coeff, err := randomScalar()
		if err != nil {
			return nil, err
		}
		coefficients[i] = coeff
	}

	shares := make([]ristretto.Scalar, count)
	for i := uint32(1); i <= count; i++ {
		x := scalarFromUint32(i)

		var share ristretto.Scalar
		share.SetZero()

		for j := range coefficients {
			share.Add(&share, &coefficients[j])
			share.Mul(&share, &x)
		}
		share.Add(&share, &secret)

		shares[i-1] = share
	}

	return shares, nil
}

func scalarFromUint32(value uint32) ristretto.Scalar {
	var scalar ristretto.Scalar
	var bytes [32]byte
	binary.LittleEndian.PutUint32(bytes[:4], value)
	scalar.SetBytes(&bytes)
	return scalar
}

func unobliviousEvaluate(privateKey ristretto.Scalar, input []byte) oprf.Output {
	var inputPoint ristretto.Point
	inputPoint.DeriveDalek(input)

	var result ristretto.Point
	result.ScalarMult(&inputPoint, &privateKey)

	h := sha512.New()
	h.Write([]byte("Juicebox_OPRF_2023_1;"))
	h.Write(input)
	h.Write(result.Bytes())

	var output oprf.Output
	copy(output[:], h.Sum(nil))
	return output
}

// buildOprfSignatureMessage constructs the message for Ed25519 signature verification.
// Format: BE2(len(realm_id)) || realm_id || BE2(len(public_key)) || public_key
func buildOprfSignatureMessage(realmID, publicKey []byte) []byte {
	msg := make([]byte, 2+len(realmID)+2+len(publicKey))
	binary.BigEndian.PutUint16(msg[0:2], uint16(len(realmID)))
	copy(msg[2:], realmID)
	offset := 2 + len(realmID)
	binary.BigEndian.PutUint16(msg[offset:offset+2], uint16(len(publicKey)))
	copy(msg[offset+2:], publicKey)
	return msg
}
