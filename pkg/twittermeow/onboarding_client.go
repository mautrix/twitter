package twittermeow

import (
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"
)

type OnboardingClient struct {
	client       *Client
	flowToken    string
	currentTasks *types.TaskResponse
}

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
