package twittermeow

import (
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"
)

type OnboardingClient struct {
	//lint:ignore U1000 TODO fix unused field
	client       *Client
	flowToken    string
	currentTasks *types.TaskResponse
}

//lint:ignore U1000 TODO fix unused method
func (c *Client) newOnboardingClient() *OnboardingClient {
	return &OnboardingClient{
		client:       c,
		currentTasks: &types.TaskResponse{},
	}
}

func (o *OnboardingClient) SetFlowToken(flowToken string) *OnboardingClient {
	o.flowToken = flowToken
	return o
}

func (o *OnboardingClient) SetCurrentTasks(tasks *types.TaskResponse) {
	o.currentTasks = tasks
}
