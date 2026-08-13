package twittermeow

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func newTestXChatWebsocketClient() *xchatWebsocketClient {
	return &xchatWebsocketClient{
		client:            &Client{Logger: zerolog.Nop()},
		initialRetryDelay: time.Millisecond,
		maximumRetryDelay: 5 * time.Millisecond,
	}
}

func blockingSuccessfulXChatConnection(
	ctx context.Context,
	_ string,
	_ zerolog.Logger,
	onConnected func(context.Context) error,
) (bool, error) {
	if err := onConnected(ctx); err != nil {
		return false, err
	}
	<-ctx.Done()
	return false, ctx.Err()
}

func TestXChatWebsocketStartupRetriesTokenBeforeReady(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	var attempts atomic.Int32
	xc.tokenProvider = func(context.Context, bool) (string, error) {
		if attempts.Add(1) == 1 {
			return "", errors.New("temporary token failure")
		}
		return "token", nil
	}
	xc.connectionRunner = blockingSuccessfulXChatConnection

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := xc.start(ctx); err != nil {
		t.Fatalf("start() error = %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Fatalf("token attempts = %d, want 2", got)
	}
	xc.stop()
}

func TestXChatWebsocketStartupWaitsForRealConnection(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	tokenAttempted := make(chan struct{})
	allowToken := make(chan struct{})
	xc.tokenProvider = func(context.Context, bool) (string, error) {
		close(tokenAttempted)
		<-allowToken
		return "token", nil
	}
	xc.connectionRunner = blockingSuccessfulXChatConnection

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startResult := make(chan error, 1)
	go func() {
		startResult <- xc.start(ctx)
	}()

	<-tokenAttempted
	select {
	case err := <-startResult:
		t.Fatalf("start() returned before a connection was established: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	close(allowToken)
	select {
	case err := <-startResult:
		if err != nil {
			t.Fatalf("start() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("start() did not return after a connection was established")
	}
	xc.stop()
}

func TestXChatWebsocketStartupCancellationBeforeReady(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	firstAttempt := make(chan struct{})
	var notified atomic.Bool
	xc.tokenProvider = func(context.Context, bool) (string, error) {
		if notified.CompareAndSwap(false, true) {
			close(firstAttempt)
		}
		return "", errors.New("token unavailable")
	}
	xc.connectionRunner = blockingSuccessfulXChatConnection

	ctx, cancel := context.WithCancel(context.Background())
	startResult := make(chan error, 1)
	go func() {
		startResult <- xc.start(ctx)
	}()
	<-firstAttempt
	cancel()

	select {
	case err := <-startResult:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("start() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("start() did not stop after context cancellation")
	}
}

func TestXChatWebsocketStartupRetriesEmptyToken(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	var attempts atomic.Int32
	var forcedRetry atomic.Bool
	xc.tokenProvider = func(_ context.Context, forceRefresh bool) (string, error) {
		if attempts.Add(1) == 1 {
			return "", nil
		}
		forcedRetry.Store(forceRefresh)
		return "token", nil
	}
	xc.connectionRunner = blockingSuccessfulXChatConnection

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := xc.start(ctx); err != nil {
		t.Fatalf("start() error = %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Fatalf("token attempts = %d, want 2", got)
	}
	if !forcedRetry.Load() {
		t.Fatal("empty cached token did not force a refresh on retry")
	}
	xc.stop()
}

func TestXChatReconnectBackoffRequiresStableConnection(t *testing.T) {
	connectedAt := time.Unix(100, 0)
	if wait, next := xchatReconnectBackoff(
		4*time.Second,
		time.Second,
		8*time.Second,
		true,
		connectedAt,
		connectedAt.Add(stableConnectionDuration-time.Second),
	); wait != 4*time.Second || next != 8*time.Second {
		t.Fatalf("short-lived connection backoff = %s/%s, want 4s/8s", wait, next)
	}
	if wait, next := xchatReconnectBackoff(
		8*time.Second,
		time.Second,
		8*time.Second,
		true,
		connectedAt,
		connectedAt.Add(stableConnectionDuration),
	); wait != time.Second || next != time.Second {
		t.Fatalf("stable connection backoff = %s/%s, want 1s/1s", wait, next)
	}
	if wait, next := xchatReconnectBackoff(
		4*time.Second,
		time.Second,
		8*time.Second,
		false,
		time.Time{},
		connectedAt,
	); wait != 4*time.Second || next != 8*time.Second {
		t.Fatalf("failed attempt backoff = %s/%s, want 4s/8s", wait, next)
	}
}

func TestXChatWebsocketRefreshesTokenWhenConnectionRequestsIt(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	var tokenCalls atomic.Int32
	var forceRefreshCalls atomic.Int32
	xc.tokenProvider = func(_ context.Context, forceRefresh bool) (string, error) {
		tokenCalls.Add(1)
		if forceRefresh {
			forceRefreshCalls.Add(1)
			return "fresh-token", nil
		}
		return "cached-token", nil
	}
	var connectionCalls atomic.Int32
	xc.connectionRunner = func(
		ctx context.Context,
		token string,
		_ zerolog.Logger,
		onConnected func(context.Context) error,
	) (bool, error) {
		if connectionCalls.Add(1) == 1 {
			if token != "cached-token" {
				t.Fatalf("first token = %q, want cached-token", token)
			}
			return true, errors.New("token rejected")
		}
		if token != "fresh-token" {
			t.Fatalf("replacement token = %q, want fresh-token", token)
		}
		return false, onConnected(ctx)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := xc.start(ctx); err != nil {
		t.Fatalf("start() error = %v", err)
	}
	if tokenCalls.Load() != 2 || forceRefreshCalls.Load() != 1 {
		t.Fatalf("token/refresh calls = %d/%d, want 2/1", tokenCalls.Load(), forceRefreshCalls.Load())
	}
	xc.stop()
}

func TestXChatWebsocketRecoversPanickedConnectionAttempt(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	xc.tokenProvider = func(context.Context, bool) (string, error) {
		return "token", nil
	}
	var attempts atomic.Int32
	xc.connectionRunner = func(
		ctx context.Context,
		_ string,
		_ zerolog.Logger,
		onConnected func(context.Context) error,
	) (bool, error) {
		if attempts.Add(1) == 1 {
			panic("simulated connection panic")
		}
		if err := onConnected(ctx); err != nil {
			return false, err
		}
		<-ctx.Done()
		return false, ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := xc.start(ctx); err != nil {
		t.Fatalf("start() error = %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Fatalf("connection attempts = %d, want 2", got)
	}
	xc.stop()
}

func TestXChatWebsocketRunsCatchupBeforeReadingEachConnection(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	xc.tokenProvider = func(context.Context, bool) (string, error) {
		return "token", nil
	}

	catchupStarted := make(chan int, 2)
	allowCatchup := make(chan struct{})
	readerStarted := make(chan int, 2)
	var catchupCalls atomic.Int32
	xc.client.SetXChatConnectHandler(func(ctx context.Context) error {
		call := int(catchupCalls.Add(1))
		catchupStarted <- call
		select {
		case <-allowCatchup:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	var connections atomic.Int32
	xc.connectionRunner = func(
		ctx context.Context,
		_ string,
		_ zerolog.Logger,
		onConnected func(context.Context) error,
	) (bool, error) {
		attempt := connections.Add(1)
		if err := onConnected(ctx); err != nil {
			return false, err
		}
		readerStarted <- int(attempt)
		if attempt == 1 {
			return false, errors.New("simulated disconnect")
		}
		<-ctx.Done()
		return false, ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startResult := make(chan error, 1)
	go func() {
		startResult <- xc.start(ctx)
	}()

	select {
	case call := <-catchupStarted:
		if call != 1 {
			t.Fatalf("first catch-up call = %d, want 1", call)
		}
	case <-time.After(time.Second):
		t.Fatal("initial socket catch-up did not start")
	}
	select {
	case attempt := <-readerStarted:
		t.Fatalf("socket reader %d started before catch-up completed", attempt)
	case <-time.After(20 * time.Millisecond):
	}
	allowCatchup <- struct{}{}

	select {
	case attempt := <-readerStarted:
		if attempt != 1 {
			t.Fatalf("first reader attempt = %d, want 1", attempt)
		}
	case <-time.After(time.Second):
		t.Fatal("initial reader did not start after catch-up completed")
	}
	select {
	case err := <-startResult:
		if err != nil {
			t.Fatalf("start() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("start() did not return after initial catch-up")
	}

	select {
	case call := <-catchupStarted:
		if call != 2 {
			t.Fatalf("replacement catch-up call = %d, want 2", call)
		}
	case <-time.After(time.Second):
		t.Fatal("replacement socket catch-up did not start")
	}
	select {
	case attempt := <-readerStarted:
		t.Fatalf("replacement reader %d started before catch-up completed", attempt)
	case <-time.After(20 * time.Millisecond):
	}
	allowCatchup <- struct{}{}
	select {
	case attempt := <-readerStarted:
		if attempt != 2 {
			t.Fatalf("replacement reader attempt = %d, want 2", attempt)
		}
	case <-time.After(time.Second):
		t.Fatal("replacement reader did not start after catch-up completed")
	}
	if got := catchupCalls.Load(); got != 2 {
		t.Fatalf("catch-up calls = %d, want 2", got)
	}
	xc.stop()
}

func TestXChatWebsocketReconnectRetriesFailedCatchup(t *testing.T) {
	xc := newTestXChatWebsocketClient()
	xc.tokenProvider = func(context.Context, bool) (string, error) {
		return "token", nil
	}

	catchupSucceeded := make(chan struct{})
	var catchupCalls atomic.Int32
	xc.client.SetXChatConnectHandler(func(context.Context) error {
		if catchupCalls.Add(1) == 1 {
			return errors.New("temporary catch-up failure")
		}
		close(catchupSucceeded)
		return nil
	})

	var connections atomic.Int32
	xc.connectionRunner = func(
		ctx context.Context,
		_ string,
		_ zerolog.Logger,
		onConnected func(context.Context) error,
	) (bool, error) {
		attempt := connections.Add(1)
		if err := onConnected(ctx); err != nil {
			return false, err
		}
		if attempt == 1 {
			return false, errors.New("simulated disconnect")
		}
		<-ctx.Done()
		return false, ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := xc.start(ctx); err != nil {
		t.Fatalf("start() error = %v", err)
	}
	select {
	case <-catchupSucceeded:
	case <-time.After(time.Second):
		t.Fatal("failed reconnect catch-up was not retried")
	}
	if got := catchupCalls.Load(); got != 2 {
		t.Fatalf("catch-up calls = %d, want 2", got)
	}
	xc.stop()
}
