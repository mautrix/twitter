package twittermeow

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
)

var (
	ErrConnectSetEventHandler = errors.New("event handler must be set before connecting")
	ErrNotAuthenticatedYet    = errors.New("client has not been authenticated yet")

	ErrAlreadyPollingUpdates = errors.New("client is already polling for user updates")
	ErrNotPollingUpdates     = errors.New("client is not polling for user updates")

	ErrCloudflareChallenge = errors.New("cloudflare challenge page")

	ErrXChatMediaSizeUnsupported = errors.New("media size rejected by X Chat")
)

type TwitterError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	ErrCouldNotAuthenticate     error = TwitterError{Code: 32}
	ErrClientNotPermitted       error = TwitterError{Code: 61}
	ErrRatelimitExceeded        error = TwitterError{Code: 88, Message: "Rate limit exceeded"}
	ErrUserSuspended            error = TwitterError{Code: 63}
	ErrAccountSuspended         error = TwitterError{Code: 63}
	ErrNotActive                error = TwitterError{Code: 141}
	ErrConversationDoesntExist  error = TwitterError{Code: 279}
	ErrAccountTemporarilyLocked error = TwitterError{Code: 326}
	ErrCSRFMismatch             error = TwitterError{Code: 353}
)

var ErrCloudflareRateLimited error = cloudflareRateLimitError{}

type cloudflareRateLimitError struct{}

func (cloudflareRateLimitError) Error() string {
	return "rate limited by Cloudflare (error 1015)"
}

func (cloudflareRateLimitError) Is(target error) bool {
	return errors.Is(ErrRatelimitExceeded, target)
}

type RateLimitError struct {
	Reset time.Time
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("rate limited until %s", e.Reset.Format(time.RFC3339))
}

func (e RateLimitError) Is(target error) bool {
	return errors.Is(ErrRatelimitExceeded, target)
}

type XChatSendFailure struct {
	Type payload.FailureType
}

func (f XChatSendFailure) Error() string {
	switch f.Type {
	case payload.FailureTypeInternalError, payload.FailureTypeGenericError, payload.FailureTypeEmptyDetail:
		return "X couldn't deliver this message, try again later"
	case payload.FailureTypeContentsTooLarge:
		return "X rejected this message: contents too large"
	case payload.FailureTypeTooManyMessages:
		return "X rejected this message: you're sending messages too fast"
	case payload.FailureTypeInvalidSenderSignature:
		return "X rejected this message: invalid sender signature"
	case payload.FailureTypeNonLatestCkeyVersion:
		return "X rejected this message: the chat's encryption key changed"
	case payload.FailureTypeRecipientHasNotTrustedConversation:
		return "X rejected this message: recipient has not trusted this conversation"
	case payload.FailureTypeRecipientKeyHasChanged:
		return "X rejected this message: the recipient's encryption key changed"
	case payload.FailureTypeOnlyEncryptedMessagesAllowed:
		return "X rejected this message: only encrypted messages are allowed in this chat"
	case payload.FailureTypeRequesterNotAdmin:
		return "X rejected this message: only group admins can do that"
	case payload.FailureTypeFlaggedAsSpam:
		return "X rejected this message: flagged as spam"
	case payload.FailureTypeRateLimitUpsell:
		return "X rejected this message: you've reached the message request limit"
	case payload.FailureTypeSignatureVerificationFailed:
		return "X rejected this message: signature verification failed"
	case payload.FailureTypeSenderNotGroupMember:
		return "X rejected this message: you're not a member of this group"
	case payload.FailureTypeInvalidSignatureVersion:
		return "X rejected this message: invalid signature version"
	case payload.FailureTypeInvalidPinRequest:
		return "X rejected this message: invalid pin request"
	case payload.FailureTypeTooManyPins:
		return "X rejected this message: too many pinned messages"
	case payload.FailureTypeInvalidXMoneyToken:
		return "X rejected this message: invalid payment token"
	case payload.FailureTypeRecipientNotMessageable:
		return "X rejected this message: this account can't receive your messages"
	default:
		return fmt.Sprintf("X rejected this message: failure type %d", f.Type)
	}
}

func (f XChatSendFailure) Retriable() bool {
	switch f.Type {
	case payload.FailureTypeEmptyDetail, payload.FailureTypeInternalError, payload.FailureTypeTooManyMessages,
		payload.FailureTypeNonLatestCkeyVersion, payload.FailureTypeGenericError:
		return true
	default:
		return false
	}
}

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
		return te.Code == ote.Code
	}
	return false
}

func (te TwitterError) TwitterErrorCode() int {
	return te.Code
}

func (te TwitterError) Error() string {
	return fmt.Sprintf("%d: %s", te.Code, te.Message)
}

type TwitterErrors struct {
	Errors     []TwitterError `json:"errors"`
	StatusCode int            `json:"-"`
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

// clientHTTPErrorPrefix mirrors bridgev2/matrix.ErrorFromClientPrefix, as documented there
const clientHTTPErrorPrefix = "error from client: "

func IsClientHTTPError(err error) bool {
	return strings.Contains(err.Error(), clientHTTPErrorPrefix)
}
