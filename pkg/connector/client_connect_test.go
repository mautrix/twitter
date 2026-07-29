package connector

import (
	"context"
	"testing"
	"time"
)

func waitForConnectSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", name)
	}
}

func TestStartConnectRunsInBackgroundAndCancelsPreviousWorker(t *testing.T) {
	tc := &TwitterClient{}
	firstStarted := make(chan struct{})
	firstCanceled := make(chan struct{})
	firstStartReturned := make(chan struct{})
	go func() {
		tc.startConnect(context.Background(), func(ctx context.Context) {
			close(firstStarted)
			<-ctx.Done()
			close(firstCanceled)
		})
		close(firstStartReturned)
	}()
	waitForConnectSignal(t, firstStartReturned, "first startConnect call to return")
	waitForConnectSignal(t, firstStarted, "first connect worker to start")

	secondStarted := make(chan struct{})
	secondCanceled := make(chan struct{})
	tc.startConnect(context.Background(), func(ctx context.Context) {
		close(secondStarted)
		<-ctx.Done()
		close(secondCanceled)
	})
	waitForConnectSignal(t, firstCanceled, "first connect worker to be canceled")
	waitForConnectSignal(t, secondStarted, "second connect worker to start")

	tc.cancelConnect()
	waitForConnectSignal(t, secondCanceled, "second connect worker to be canceled")
}

func TestStartConnectWorkerStopsWithParentContext(t *testing.T) {
	tc := &TwitterClient{}
	parentCtx, cancelParent := context.WithCancel(context.Background())
	workerStarted := make(chan struct{})
	workerCanceled := make(chan struct{})
	tc.startConnect(parentCtx, func(ctx context.Context) {
		close(workerStarted)
		<-ctx.Done()
		close(workerCanceled)
	})
	waitForConnectSignal(t, workerStarted, "connect worker to start")

	cancelParent()
	waitForConnectSignal(t, workerCanceled, "connect worker to observe parent cancellation")
}
