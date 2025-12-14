package juicebox

import "errors"

var (
	ErrInvalidAuth       = errors.New("juicebox: invalid auth token")
	ErrUpgradeRequired   = errors.New("juicebox: SDK upgrade required")
	ErrRateLimitExceeded = errors.New("juicebox: rate limit exceeded")
	ErrAssertion         = errors.New("juicebox: assertion error")
	ErrTransient         = errors.New("juicebox: transient error")
	ErrInvalidPin        = errors.New("juicebox: invalid PIN")
	ErrNotRegistered     = errors.New("juicebox: secret not registered")
)

// RecoverError represents an error during secret recovery.
type RecoverError struct {
	Reason           error
	GuessesRemaining *uint16
}

func (e *RecoverError) Error() string {
	return e.Reason.Error()
}

func (e *RecoverError) Unwrap() error {
	return e.Reason
}
