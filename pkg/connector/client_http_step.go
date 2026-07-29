package connector

import (
	"bytes"
	"fmt"

	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
)

func makeClientHTTPRequestStep(request *twittermeow.ClientHTTPRequest) (*bridgev2.LoginStep, error) {
	if request == nil {
		return nil, fmt.Errorf("client HTTP request is missing")
	}
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeClientHTTP,
		StepID:       LoginStepIDClientHTTPRequest,
		Instructions: "Sending an X sign-in request from this device. Client HTTP is in beta.",
		ClientHTTPParams: &bridgev2.LoginClientHTTPParams{
			RequestID: request.ID,
			Method:    request.Method,
			URL:       request.URL,
			Headers:   request.Headers.Clone(),
			Body:      bytes.Clone(request.Body),
		},
	}, nil
}
