package juicebox

/*
#include "./juicebox-sdk-ffi.h"
#include <stdlib.h>
#include <string.h>

// C helper to invoke the HTTP response callback
static void invokeHttpResponseCallback(JuiceboxHttpResponseFn fn, JuiceboxHttpClientState *ctx, JuiceboxHttpResponse *resp) {
	fn(ctx, resp);
}

// C helper to invoke the auth token callback
static void invokeAuthTokenCallback(JuiceboxAuthTokenGetCallbackFn fn, JuiceboxAuthTokenManager *ctx, uint64_t ctx_id, JuiceboxAuthToken *token) {
	fn(ctx, ctx_id, token);
}
*/
import "C"

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"sync"
	"unsafe"

	"github.com/rs/zerolog"
)

// SDKVersion is the Juicebox SDK version. Must match the compiled Rust SDK.
// This is sent as the X-Juicebox-Version header which the FFI bridge doesn't add automatically.
const SDKVersion = "0.3.4"

// ClientState holds references for FFI callbacks.
type ClientState struct {
	HTTPClient *http.Client
	// AuthTokens maps realm ID (hex string, lowercase) to pre-generated auth token string.
	// Used when we have pre-fetched tokens from the API.
	AuthTokens map[string]string
	Logger     zerolog.Logger
}

// Global state - the SDK doesn't support passing context through callbacks,
// so we use a package-level variable. This is safe because:
// 1. Juicebox operations are blocking (we wait for result)
// 2. We only have one active client at a time per process
var (
	currentState   *ClientState
	currentStateMu sync.RWMutex
)

func setCurrentState(state *ClientState) {
	currentStateMu.Lock()
	defer currentStateMu.Unlock()
	currentState = state
}

func getCurrentState() *ClientState {
	currentStateMu.RLock()
	defer currentStateMu.RUnlock()
	return currentState
}

func clearCurrentState() {
	currentStateMu.Lock()
	defer currentStateMu.Unlock()
	currentState = nil
}

// C-allocated memory for non-null empty array pointers.
// The Rust SDK asserts that array data pointers are non-null, even for zero-length arrays.
// We use C.malloc because Go pointers cannot be passed to C code.
var (
	emptyHeaderPtr *C.JuiceboxHttpHeader
	emptyBytePtr   *C.uint8_t
)

func init() {
	// Allocate C memory once at startup (never freed, lives for program duration)
	emptyHeaderPtr = (*C.JuiceboxHttpHeader)(C.malloc(C.size_t(unsafe.Sizeof(C.JuiceboxHttpHeader{}))))
	emptyBytePtr = (*C.uint8_t)(C.malloc(1))
}

// allocResponse allocates a C response struct and initializes with non-null pointers.
// The caller must free it with C.free after the callback returns.
func allocResponse(reqID [16]byte, statusCode uint16) *C.JuiceboxHttpResponse {
	resp := (*C.JuiceboxHttpResponse)(C.malloc(C.size_t(unsafe.Sizeof(C.JuiceboxHttpResponse{}))))
	copy((*[16]byte)(unsafe.Pointer(&resp.id[0]))[:], reqID[:])
	resp.status_code = C.uint16_t(statusCode)
	// Provide non-null C pointers for empty arrays
	resp.headers.data = emptyHeaderPtr
	resp.headers.length = 0
	resp.body.data = emptyBytePtr
	resp.body.length = 0
	return resp
}

//export goHttpSendCallback
func goHttpSendCallback(ctx *C.JuiceboxHttpClientState, req *C.JuiceboxHttpRequest, callback C.JuiceboxHttpResponseFn) {
	state := getCurrentState()
	if state == nil {
		// Return error response with C-allocated struct
		var reqID [16]byte
		copy(reqID[:], (*[16]byte)(unsafe.Pointer(&req.id[0]))[:])
		resp := allocResponse(reqID, 500)
		C.invokeHttpResponseCallback(callback, ctx, resp)
		C.free(unsafe.Pointer(resp))
		return
	}

	// Extract request data (must copy before returning)
	var reqID [16]byte
	copy(reqID[:], (*[16]byte)(unsafe.Pointer(&req.id[0]))[:])
	method := httpMethodToString(req.method)
	url := C.GoString(req.url)
	body := unmanagedArrayToBytes(req.body)

	// Extract headers
	headers := make(map[string]string)
	if req.headers.length > 0 {
		headerSlice := unsafe.Slice(req.headers.data, req.headers.length)
		for _, h := range headerSlice {
			headers[C.GoString(h.name)] = C.GoString(h.value)
		}
	}

	// Make HTTP request in a goroutine
	go func() {
		httpReq, err := http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			sendErrorResponse(ctx, callback, reqID, 500)
			return
		}

		// Add required Juicebox version header (FFI bridge doesn't add it automatically)
		httpReq.Header.Set("X-Juicebox-Version", SDKVersion)
		httpReq.Header.Set("User-Agent", "JuiceboxSdk-Go/"+SDKVersion)

		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}

		resp, err := state.HTTPClient.Do(httpReq)
		if err != nil {
			sendErrorResponse(ctx, callback, reqID, 503)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		sendSuccessResponse(ctx, callback, reqID, uint16(resp.StatusCode), respBody)
	}()
}

func sendErrorResponse(ctx *C.JuiceboxHttpClientState, callback C.JuiceboxHttpResponseFn, reqID [16]byte, statusCode uint16) {
	resp := allocResponse(reqID, statusCode)
	C.invokeHttpResponseCallback(callback, ctx, resp)
	C.free(unsafe.Pointer(resp))
}

func sendSuccessResponse(ctx *C.JuiceboxHttpClientState, callback C.JuiceboxHttpResponseFn, reqID [16]byte, statusCode uint16, body []byte) {
	resp := allocResponse(reqID, statusCode)

	var cBody unsafe.Pointer
	if len(body) > 0 {
		// Allocate body that will be valid for the callback
		cBody = C.CBytes(body)
		resp.body.data = (*C.uint8_t)(cBody)
		resp.body.length = C.size_t(len(body))
	}

	C.invokeHttpResponseCallback(callback, ctx, resp)

	// Free after callback returns
	if cBody != nil {
		C.free(cBody)
	}
	C.free(unsafe.Pointer(resp))
}

//export goAuthTokenGetCallback
func goAuthTokenGetCallback(ctx *C.JuiceboxAuthTokenManager, ctxID C.uint64_t, realmID *C.uint8_t, callback C.JuiceboxAuthTokenGetCallbackFn) {
	state := getCurrentState()
	if state == nil {
		C.invokeAuthTokenCallback(callback, ctx, ctxID, nil)
		return
	}

	// Copy realm ID and convert to lowercase hex string (matching reverse-xchat)
	var rid [16]byte
	copy(rid[:], (*[16]byte)(unsafe.Pointer(realmID))[:])
	realmIDHex := strings.ToLower(hex.EncodeToString(rid[:]))

	state.Logger.Debug().
		Str("realm_id", realmIDHex).
		Int("available_tokens", len(state.AuthTokens)).
		Msg("Juicebox auth token requested")

	// Look up pre-fetched token for this realm
	if tokenStr, ok := state.AuthTokens[realmIDHex]; ok {
		state.Logger.Debug().
			Str("realm_id", realmIDHex).
			Str("token_preview", tokenStr[:min(50, len(tokenStr))]+"...").
			Msg("Found auth token for realm")
		cTokenStr := C.CString(tokenStr)
		token := C.juicebox_auth_token_create(cTokenStr)
		C.free(unsafe.Pointer(cTokenStr))
		C.invokeAuthTokenCallback(callback, ctx, ctxID, token)
		return
	}

	// Log available keys for debugging
	var availableKeys []string
	for k := range state.AuthTokens {
		availableKeys = append(availableKeys, k)
	}
	state.Logger.Error().
		Str("requested_realm_id", realmIDHex).
		Strs("available_realm_ids", availableKeys).
		Msg("No auth token found for requested realm")
	C.invokeAuthTokenCallback(callback, ctx, ctxID, nil)
}
