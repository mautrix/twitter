package twittermeow

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrConnectSetEventHandler = errors.New("event handler must be set before connecting")
	ErrNotAuthenticatedYet    = errors.New("client has not been authenticated yet")

	ErrAlreadyPollingUpdates = errors.New("client is already polling for user updates")
	ErrNotPollingUpdates     = errors.New("client is not polling for user updates")
)

type TwitterError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	ErrCouldNotAuthenticate     error = TwitterError{Code: 32}
	ErrUserSuspended            error = TwitterError{Code: 63}
	ErrAccountSuspended         error = TwitterError{Code: 63}
	ErrNotActive                error = TwitterError{Code: 141}
	ErrAccountTemporarilyLocked error = TwitterError{Code: 326}
)

func IsAuthError(err error) bool {
	return errors.Is(err, ErrCouldNotAuthenticate) ||
		errors.Is(err, ErrUserSuspended) ||
		errors.Is(err, ErrAccountSuspended) ||
		errors.Is(err, ErrNotActive) ||
		errors.Is(err, ErrAccountTemporarilyLocked)
}

func (te TwitterError) Is(other error) bool {
	var ote TwitterError
	if errors.As(other, &ote) {
		return te.Code == ote.Code || te.Message == ote.Message
	}
	return false
}

func (te TwitterError) Error() string {
	return fmt.Sprintf("%d: %s", te.Code, te.Message)
}

type TwitterErrors struct {
	Errors []TwitterError `json:"errors"`
}

func (te *TwitterErrors) Error() string {
	if te == nil || len(te.Errors) == 0 {
		return "no errors"
	} else if len(te.Errors) == 1 {
		return te.Errors[0].Error()
	} else {
		errs := make([]string, len(te.Errors))
		for i, e := range te.Errors {
			errs[i] = e.Error()
		}
		return strings.Join(errs, ", ")
	}
}

func (te *TwitterErrors) Unwrap() []error {
	errs := make([]error, len(te.Errors))
	for i, e := range te.Errors {
		errs[i] = e
	}
	return errs
}
