package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-twitter/pkg/connector"
)

func jsonResponse(w http.ResponseWriter, status int, response any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

type Error struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	ErrCode string `json:"errcode"`
}

type Response struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

func legacyProvLogin(w http.ResponseWriter, r *http.Request) {
	user := m.Matrix.Provisioning.GetUser(r)
	ctx := r.Context()
	var cookies map[string]string
	err := json.NewDecoder(r.Body).Decode(&cookies)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, Error{ErrCode: mautrix.MBadJSON.ErrCode, Error: err.Error()})
		return
	}

	newCookies := map[string]string{
		"auth_token": cookies["auth_token"],
		"ct0":        cookies["csrf_token"],
	}

	lp, err := c.CreateLogin(ctx, user, "cookies")
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to create login")
		jsonResponse(w, http.StatusInternalServerError, Error{ErrCode: "M_UNKNOWN", Error: "Internal error creating login"})
	} else if firstStep, err := lp.Start(ctx); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to start login")
		jsonResponse(w, http.StatusInternalServerError, Error{ErrCode: "M_UNKNOWN", Error: "Internal error starting login"})
	} else if firstStep.StepID != connector.LoginStepIDCookies {
		jsonResponse(w, http.StatusInternalServerError, Error{ErrCode: "M_UNKNOWN", Error: "Unexpected login step"})
	} else if finalStep, err := lp.(bridgev2.LoginProcessCookies).SubmitCookies(ctx, newCookies); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to log in")
		var respErr bridgev2.RespError
		if errors.As(err, &respErr) {
			jsonResponse(w, respErr.StatusCode, &respErr)
		} else {
			jsonResponse(w, http.StatusInternalServerError, Error{ErrCode: "M_UNKNOWN", Error: "Internal error logging in"})
		}
	} else if finalStep.StepID != connector.LoginStepIDComplete {
		jsonResponse(w, http.StatusInternalServerError, Error{ErrCode: "M_UNKNOWN", Error: "Unexpected login step"})
	} else {
		jsonResponse(w, http.StatusOK, map[string]any{})
		go handleLoginComplete(context.WithoutCancel(ctx), user, finalStep.CompleteParams.UserLogin)
	}
}

func handleLoginComplete(ctx context.Context, user *bridgev2.User, newLogin *bridgev2.UserLogin) {
	allLogins := user.GetUserLogins()
	for _, login := range allLogins {
		if login.ID != newLogin.ID {
			login.Delete(ctx, status.BridgeState{StateEvent: status.StateLoggedOut, Reason: "LOGIN_OVERRIDDEN"}, bridgev2.DeleteOpts{})
		}
	}
}

func legacyProvLogout(w http.ResponseWriter, r *http.Request) {
	user := m.Matrix.Provisioning.GetUser(r)
	logins := user.GetUserLogins()
	for _, login := range logins {
		// Intentionally don't delete the user login, only disconnect the client
		login.Client.(*connector.TwitterClient).LogoutRemote(r.Context())
	}
	jsonResponse(w, http.StatusOK, Response{
		Success: true,
		Status:  "logged_out",
	})
}
