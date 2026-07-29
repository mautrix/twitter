package twittermeow

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

func TestRefreshConversationKeysDoesNotBlockOnConversationDataCallback(t *testing.T) {
	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	client.HTTP = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(
				`{"data":{"get_inbox_page_conversation_data":{"data":{"conversation_detail":{"conversation_id":"1:2"}}}}}`,
			)),
		}, nil
	})}

	callbackStarted := make(chan struct{})
	callbackRelease := make(chan struct{})
	client.SetConversationDataCallback(func(context.Context, string, *response.XChatInboxItem) {
		close(callbackStarted)
		<-callbackRelease
	})

	refreshDone := make(chan error, 1)
	go func() {
		refreshDone <- client.RefreshConversationKeys(context.Background(), "1:2")
	}()

	select {
	case <-callbackStarted:
	case <-time.After(time.Second):
		t.Fatal("conversation data callback was not invoked")
	}

	select {
	case err := <-refreshDone:
		if err != nil {
			t.Fatalf("RefreshConversationKeys() error = %v", err)
		}
		close(callbackRelease)
	case <-time.After(250 * time.Millisecond):
		close(callbackRelease)
		err := <-refreshDone
		if err != nil {
			t.Fatalf("RefreshConversationKeys() error after releasing callback = %v", err)
		}
		t.Fatal("RefreshConversationKeys blocked on the conversation data callback")
	}
}
