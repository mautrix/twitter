package juicebox

/*
#include "./juicebox-sdk-ffi.h"
#include <stdlib.h>

// Forward declarations for Go callbacks
extern void goHttpSendCallback(JuiceboxHttpClientState *ctx, JuiceboxHttpRequest *req, JuiceboxHttpResponseFn callback);
extern void goAuthTokenGetCallback(JuiceboxAuthTokenManager *ctx, uint64_t ctx_id, uint8_t *realm_id, JuiceboxAuthTokenGetCallbackFn callback);
extern void goRecoverCallback(void *ctx, JuiceboxUnmanagedDataArray secret, JuiceboxRecoverError *err);

// Typedef for recover callback (not in header, defined inline)
typedef void (*JuiceboxRecoverResponseFn)(const void *context, JuiceboxUnmanagedDataArray secret, const JuiceboxRecoverError *error);

// Helper to call recover with our callback
static void callRecover(JuiceboxClient *client, void *ctx, JuiceboxUnmanagedDataArray pin, JuiceboxUnmanagedDataArray info) {
	juicebox_client_recover(client, ctx, pin, info, (void (*)(const void*, JuiceboxUnmanagedDataArray, const JuiceboxRecoverError*))goRecoverCallback);
}
*/
import "C"

import (
	"context"
	"net/http"
	"runtime"
	"runtime/cgo"
	"strings"
	"sync"
	"unsafe"

	"github.com/rs/zerolog"
)

// Client is a Juicebox SDK client for PIN-protected secret storage.
type Client struct {
	ptr    *C.JuiceboxClient
	config *Configuration
	state  *ClientState

	mu sync.Mutex
}

// NewClient creates a new Juicebox client.
// authTokens is a map of realm ID (hex string) to pre-fetched JWT auth token.
func NewClient(config *Configuration, authTokens map[string]string, httpClient *http.Client, logger zerolog.Logger) (*Client, error) {
	if config == nil || config.ptr == nil {
		return nil, ErrAssertion
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Normalize auth token keys to lowercase for consistent lookup
	normalizedTokens := make(map[string]string, len(authTokens))
	for k, v := range authTokens {
		normalizedTokens[strings.ToLower(k)] = v
	}

	// Create callback state with pre-fetched tokens
	state := &ClientState{
		HTTPClient: httpClient,
		AuthTokens: normalizedTokens,
		Logger:     logger,
	}

	// Set global state for callbacks (SDK doesn't support passing context)
	setCurrentState(state)

	// Create empty previous configurations array
	// Use C-allocated memory to provide non-null pointer (Rust SDK requires non-null data)
	emptyPtr := C.malloc(C.size_t(unsafe.Sizeof(uintptr(0))))
	defer C.free(emptyPtr)
	*(**C.JuiceboxConfiguration)(emptyPtr) = nil

	prevConfigs := C.JuiceboxUnmanagedConfigurationArray{
		data:   (**C.JuiceboxConfiguration)(emptyPtr),
		length: 0,
	}

	// Create the client
	clientPtr := C.juicebox_client_create(
		config.ptr,
		prevConfigs,
		C.JuiceboxAuthTokenGetFn(C.goAuthTokenGetCallback),
		C.JuiceboxHttpSendFn(C.goHttpSendCallback),
	)

	if clientPtr == nil {
		clearCurrentState()
		return nil, ErrAssertion
	}

	return &Client{
		ptr:    clientPtr,
		config: config,
		state:  state,
	}, nil
}

// Close releases all resources associated with the client.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear global state
	clearCurrentState()

	if c.ptr != nil {
		C.juicebox_client_destroy(c.ptr)
		c.ptr = nil
	}
}

// recoverResult holds the result of a recover operation.
type recoverResult struct {
	secret []byte
	err    error
}

// Recover retrieves a PIN-protected secret.
func (c *Client) Recover(ctx context.Context, pin Pin, info UserInfo) (Secret, error) {
	c.mu.Lock()
	if c.ptr == nil {
		c.mu.Unlock()
		return nil, ErrAssertion
	}
	clientPtr := c.ptr
	c.mu.Unlock()

	// Create result channel
	resultChan := make(chan recoverResult, 1)

	// Create context for the callback
	callbackCtx := &recoverCallbackContext{
		resultChan: resultChan,
	}
	handle := cgo.NewHandle(callbackCtx)

	// Convert pin and info to C arrays
	pinArr := bytesToUnmanagedArray(pin)
	infoArr := bytesToUnmanagedArray(info)

	// Call the FFI function via helper
	C.callRecover(
		clientPtr,
		unsafe.Pointer(&handle),
		pinArr,
		infoArr,
	)

	// Keep references alive until callback completes
	runtime.KeepAlive(pin)
	runtime.KeepAlive(info)

	// Wait for result or context cancellation
	select {
	case result := <-resultChan:
		handle.Delete()
		return result.secret, result.err
	case <-ctx.Done():
		// Note: The callback may still fire, but we'll ignore it
		handle.Delete()
		return nil, ctx.Err()
	}
}

// recoverCallbackContext holds context for the recover callback.
type recoverCallbackContext struct {
	resultChan chan recoverResult
}

//export goRecoverCallback
func goRecoverCallback(ctx unsafe.Pointer, secret C.JuiceboxUnmanagedDataArray, recoverErr *C.JuiceboxRecoverError) {
	handle := *(*cgo.Handle)(ctx)
	callbackCtx := handle.Value().(*recoverCallbackContext)

	var result recoverResult

	if recoverErr != nil {
		// Map error reason to Go error
		switch recoverErr.reason {
		case C.JuiceboxRecoverErrorReasonInvalidPin:
			err := &RecoverError{Reason: ErrInvalidPin}
			if recoverErr.guesses_remaining != nil {
				guesses := uint16(*recoverErr.guesses_remaining)
				err.GuessesRemaining = &guesses
			}
			result.err = err
		case C.JuiceboxRecoverErrorReasonNotRegistered:
			result.err = ErrNotRegistered
		case C.JuiceboxRecoverErrorReasonInvalidAuth:
			result.err = ErrInvalidAuth
		case C.JuiceboxRecoverErrorReasonUpgradeRequired:
			result.err = ErrUpgradeRequired
		case C.JuiceboxRecoverErrorReasonRateLimitExceeded:
			result.err = ErrRateLimitExceeded
		case C.JuiceboxRecoverErrorReasonAssertion:
			result.err = ErrAssertion
		case C.JuiceboxRecoverErrorReasonTransient:
			result.err = ErrTransient
		default:
			result.err = ErrAssertion
		}
	} else {
		// Success - copy the secret
		result.secret = unmanagedArrayToBytes(secret)
	}

	// Send result (non-blocking in case context was cancelled)
	select {
	case callbackCtx.resultChan <- result:
	default:
	}
}
