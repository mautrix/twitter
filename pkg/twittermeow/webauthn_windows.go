//go:build windows

package twittermeow

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	windowsWebAuthnClientDataVersion                  = 1
	windowsWebAuthnCredentialVersion                  = 1
	windowsWebAuthnGetAssertionOptionsVersion         = 2
	windowsWebAuthnUserVerificationDiscouraged        = 3
	windowsWebAuthnDefaultAssertionTimeoutMS          = 60000
	windowsWebAuthnCredentialTypePublicKey            = "public-key"
	windowsWebAuthnHashAlgorithmSHA256                = "SHA-256"
	windowsWebAuthnCredentialAttachmentAny     uint32 = 0
)

var (
	windowsWebAuthnDLL               = windows.NewLazySystemDLL("webauthn.dll")
	windowsWebAuthnGetAssertionProc  = windowsWebAuthnDLL.NewProc("WebAuthNAuthenticatorGetAssertion")
	windowsWebAuthnFreeAssertionProc = windowsWebAuthnDLL.NewProc("WebAuthNFreeAssertion")
)

type windowsWebAuthnClientData struct {
	Version         uint32
	ClientDataJSONN uint32
	ClientDataJSON  *byte
	HashAlgID       *uint16
}

type windowsWebAuthnCredential struct {
	Version        uint32
	IDN            uint32
	ID             *byte
	CredentialType *uint16
}

type windowsWebAuthnCredentials struct {
	CredentialsN uint32
	Credentials  *windowsWebAuthnCredential
}

type windowsWebAuthnExtension struct {
	ExtensionIdentifier *uint16
	ExtensionN          uint32
	Extension           unsafe.Pointer
}

type windowsWebAuthnExtensions struct {
	ExtensionsN uint32
	Extensions  *windowsWebAuthnExtension
}

type windowsWebAuthnGetAssertionOptionsV2 struct {
	Version                     uint32
	TimeoutMilliseconds         uint32
	CredentialList              windowsWebAuthnCredentials
	Extensions                  windowsWebAuthnExtensions
	AuthenticatorAttachment     uint32
	UserVerificationRequirement uint32
	Flags                       uint32
	U2fAppID                    *uint16
	U2fAppIDUsed                *int32
}

type windowsWebAuthnAssertion struct {
	Version                     uint32
	AuthenticatorDataN          uint32
	AuthenticatorData           *byte
	SignatureN                  uint32
	Signature                   *byte
	Credential                  windowsWebAuthnCredential
	UserIDN                     uint32
	UserID                      *byte
	Extensions                  windowsWebAuthnExtensions
	CredLargeBlobN              uint32
	CredLargeBlob               *byte
	CredLargeBlobStatus         uint32
	HmacSecret                  unsafe.Pointer
	UsedTransport               uint32
	UnsignedExtensionOutputsN   uint32
	UnsignedExtensionOutputs    *byte
	ClientDataJSONN             uint32
	ClientDataJSON              *byte
	AuthenticationResponseJSONN uint32
	AuthenticationResponseJSON  *byte
}

func platformWebAuthnGetAssertion(_ context.Context, challenge webAuthnChallenge, clientDataJSON []byte) (*webAuthnAssertion, error) {
	rpID, err := windows.UTF16PtrFromString(challenge.rpID())
	if err != nil {
		return nil, err
	}
	hashAlg, err := windows.UTF16PtrFromString(windowsWebAuthnHashAlgorithmSHA256)
	if err != nil {
		return nil, err
	}
	credentialType, err := windows.UTF16PtrFromString(windowsWebAuthnCredentialTypePublicKey)
	if err != nil {
		return nil, err
	}

	credentialIDs, credentials, err := windowsWebAuthnCredentialsFromChallenge(challenge, credentialType)
	if err != nil {
		return nil, err
	}
	var credentialPtr *windowsWebAuthnCredential
	if len(credentials) > 0 {
		credentialPtr = &credentials[0]
	}

	clientData := windowsWebAuthnClientData{
		Version:         windowsWebAuthnClientDataVersion,
		ClientDataJSONN: uint32(len(clientDataJSON)),
		ClientDataJSON:  byteSlicePtr(clientDataJSON),
		HashAlgID:       hashAlg,
	}
	opts := windowsWebAuthnGetAssertionOptionsV2{
		Version:                     windowsWebAuthnGetAssertionOptionsVersion,
		TimeoutMilliseconds:         windowsWebAuthnDefaultAssertionTimeoutMS,
		CredentialList:              windowsWebAuthnCredentials{CredentialsN: uint32(len(credentials)), Credentials: credentialPtr},
		AuthenticatorAttachment:     windowsWebAuthnCredentialAttachmentAny,
		UserVerificationRequirement: windowsWebAuthnUserVerificationDiscouraged,
	}
	var u2fAppIDUsed int32
	if appID := strings.TrimSpace(challenge.RequestOptions.Extensions.AppID); appID != "" {
		u2fAppID, err := windows.UTF16PtrFromString(appID)
		if err != nil {
			return nil, err
		}
		opts.U2fAppID = u2fAppID
		opts.U2fAppIDUsed = &u2fAppIDUsed
	}

	var assertionPtr *windowsWebAuthnAssertion
	result, _, _ := windowsWebAuthnGetAssertionProc.Call(
		0,
		uintptr(unsafe.Pointer(rpID)),
		uintptr(unsafe.Pointer(&clientData)),
		uintptr(unsafe.Pointer(&opts)),
		uintptr(unsafe.Pointer(&assertionPtr)),
	)
	runtime.KeepAlive(clientDataJSON)
	runtime.KeepAlive(credentialIDs)
	runtime.KeepAlive(credentials)
	runtime.KeepAlive(credentialType)
	runtime.KeepAlive(opts)
	if result != 0 {
		return nil, fmt.Errorf("windows webauthn assertion failed: HRESULT 0x%08x", uint32(result))
	}
	if assertionPtr == nil {
		return nil, fmt.Errorf("windows webauthn assertion was empty")
	}
	defer windowsWebAuthnFreeAssertionProc.Call(uintptr(unsafe.Pointer(assertionPtr)))

	assertion := *assertionPtr
	assertionClientData := copyWindowsBytes(assertion.ClientDataJSON, assertion.ClientDataJSONN)
	if len(assertionClientData) == 0 {
		assertionClientData = append([]byte(nil), clientDataJSON...)
	}
	return &webAuthnAssertion{
		CredentialID:      copyWindowsBytes(assertion.Credential.ID, assertion.Credential.IDN),
		ClientDataJSON:    assertionClientData,
		AuthenticatorData: copyWindowsBytes(assertion.AuthenticatorData, assertion.AuthenticatorDataN),
		Signature:         copyWindowsBytes(assertion.Signature, assertion.SignatureN),
	}, nil
}

func windowsWebAuthnCredentialsFromChallenge(challenge webAuthnChallenge, credentialType *uint16) ([][]byte, []windowsWebAuthnCredential, error) {
	ids := make([][]byte, 0, len(challenge.RequestOptions.AllowCredentials))
	credentials := make([]windowsWebAuthnCredential, 0, len(challenge.RequestOptions.AllowCredentials))
	for _, allowed := range challenge.RequestOptions.AllowCredentials {
		if allowed.Type != "" && allowed.Type != windowsWebAuthnCredentialTypePublicKey {
			continue
		}
		id, err := webAuthnDecodeBase64(allowed.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("decode X security-key credential ID: %w", err)
		}
		if len(id) == 0 {
			continue
		}
		ids = append(ids, id)
		credentials = append(credentials, windowsWebAuthnCredential{
			Version:        windowsWebAuthnCredentialVersion,
			IDN:            uint32(len(id)),
			ID:             byteSlicePtr(id),
			CredentialType: credentialType,
		})
	}
	return ids, credentials, nil
}

func byteSlicePtr(data []byte) *byte {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

func copyWindowsBytes(ptr *byte, size uint32) []byte {
	if ptr == nil || size == 0 {
		return nil
	}
	return append([]byte(nil), unsafe.Slice(ptr, int(size))...)
}
