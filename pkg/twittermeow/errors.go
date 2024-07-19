package twittermeow

import "errors"

var (
	// Connection errors
	ErrConnectPleaseSetEventHandler = errors.New("please set event handler in client before connecting")
	ErrNotAuthenticatedYet          = errors.New("client has not been authenticated yet")

	// Polling errors
	ErrAlreadyPollingUpdates = errors.New("client is already polling for user updates")
	ErrNotPollingUpdates     = errors.New("client is not polling for user updates")

	// Api errors
	ErrFailedMarkConversationRead = errors.New("failed to mark conversation as read")
)
