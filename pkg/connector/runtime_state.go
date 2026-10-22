package connector

import (
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

// A down link stays down until that path recovers, so the other path cannot mark the bridge connected.
type runtimeLink struct {
	down bool
	auth bool
	err  error
}

func (tc *TwitterClient) resetRuntimeSession() uint64 {
	tc.runtimeMu.Lock()
	defer tc.runtimeMu.Unlock()
	tc.runtimeGen++
	tc.runtimeStarted = false
	tc.pollingLink = runtimeLink{}
	tc.xchatLink = runtimeLink{}
	return tc.runtimeGen
}

func (tc *TwitterClient) notePollingRuntime(evt *types.PollingError) {
	link := runtimeLink{}
	var gen uint64
	if evt != nil {
		gen = evt.Gen
		if evt.IsAuth || evt.Error != nil {
			link = runtimeLink{down: true, auth: evt.IsAuth, err: evt.Error}
		}
	}
	tc.setRuntimeLink(true, link, gen)
}

func (tc *TwitterClient) noteXChatRuntime(gen uint64, state twittermeow.XChatSocketState, err error) {
	link := runtimeLink{}
	switch state {
	case twittermeow.XChatSocketBadCredentials:
		link = runtimeLink{down: true, auth: true, err: err}
	case twittermeow.XChatSocketDisconnected:
		link = runtimeLink{down: true, err: err}
	}
	tc.setRuntimeLink(false, link, gen)
}

func (tc *TwitterClient) setRuntimeLink(polling bool, link runtimeLink, gen uint64) {
	tc.runtimeMu.Lock()
	defer tc.runtimeMu.Unlock()
	if gen != tc.runtimeGen {
		return
	}
	if polling {
		tc.pollingLink = link
	} else {
		tc.xchatLink = link
	}
	state := tc.runtimeBridgeStateLocked()
	if !tc.runtimeStarted && state.StateEvent != status.StateBadCredentials {
		return
	}
	tc.sendRuntimeState(state)
}

func (tc *TwitterClient) publishSyncInProgress(gen uint64) {
	tc.runtimeMu.Lock()
	defer tc.runtimeMu.Unlock()
	if gen != tc.runtimeGen {
		return
	}
	tc.runtimeStarted = true
	state := tc.runtimeBridgeStateLocked()
	if state.StateEvent == status.StateConnected {
		state.Reason = "sync_in_progress"
		state.Info = map[string]any{"sync": "in_progress"}
	}
	tc.sendRuntimeState(state)
}

func (tc *TwitterClient) publishRuntimeState() {
	tc.runtimeMu.Lock()
	defer tc.runtimeMu.Unlock()
	tc.runtimeStarted = true
	tc.sendRuntimeState(tc.runtimeBridgeStateLocked())
}

func (tc *TwitterClient) sendRuntimeState(state status.BridgeState) {
	if tc.runtimeSend != nil {
		tc.runtimeSend(state)
		return
	}
	if tc.userLogin == nil {
		return
	}
	tc.userLogin.BridgeState.Send(state)
}

func (tc *TwitterClient) runtimeBridgeStateLocked() status.BridgeState {
	if tc.xchatLink.auth || tc.pollingLink.auth {
		err := tc.pollingLink.err
		if tc.xchatLink.auth {
			err = tc.xchatLink.err
		}
		state := status.BridgeState{
			StateEvent: status.StateBadCredentials,
			Error:      authErrorCode(err),
			UserAction: authErrorUserAction(err),
		}
		if err != nil {
			state.Info = map[string]any{"go_error": err.Error()}
		}
		return state
	}
	if tc.xchatLink.down {
		return runtimeTransientState("twitter-xchat-disconnected", tc.xchatLink.err)
	}
	if tc.pollingLink.down {
		return runtimeTransientState("twitter-polling-error", tc.pollingLink.err)
	}
	return status.BridgeState{StateEvent: status.StateConnected}
}

func runtimeTransientState(code status.BridgeStateErrorCode, err error) status.BridgeState {
	state := status.BridgeState{
		StateEvent: status.StateTransientDisconnect,
		Error:      code,
	}
	if err != nil {
		state.Info = map[string]any{"go_error": err.Error()}
	}
	return state
}
