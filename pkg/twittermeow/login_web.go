package twittermeow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

const (
	webLoginSubtaskJSInstrumentation = "LoginJsInstrumentationSubtask"
	webLoginSubtaskIdentifier        = "LoginEnterUserIdentifierSSO"
	webLoginSubtaskTwoFactor         = "LoginTwoFactorAuthChallenge"

	webLoginLinkNext = "next_link"
)

var (
	ErrWebLoginUnexpectedSubtask = errors.New("unexpected X login subtask")
	ErrWebLoginMissingFlowToken  = errors.New("x login response did not include a flow token")
	ErrWebLoginMissingGuestToken = errors.New("x guest activation response did not include a guest token")
)

type WebLoginStatus string

const (
	WebLoginStatusNeedsIdentifier WebLoginStatus = "needs_identifier"
	WebLoginStatusNeedsPassword   WebLoginStatus = "needs_password"
	WebLoginStatusNeedsText       WebLoginStatus = "needs_text"
	WebLoginStatusComplete        WebLoginStatus = "complete"
	WebLoginStatusUnsupported     WebLoginStatus = "unsupported"
)

type WebLoginError struct {
	Code    int
	Message string
}

func (wle *WebLoginError) Error() string {
	if wle == nil {
		return ""
	}
	if wle.Code != 0 {
		return fmt.Sprintf("X login failed (%d): %s", wle.Code, wle.Message)
	}
	return fmt.Sprintf("X login failed: %s", wle.Message)
}

func (wle *WebLoginError) UserMessage() string {
	if wle == nil {
		return "X rejected this login. Please check the details and try again."
	}
	switch wle.Code {
	case 32, 64, 89, 99:
		return "X rejected the username or password. Check the details and try again."
	case 88, 226, 326, 399:
		return "X rejected this login attempt. Wait a bit, then try again."
	default:
		msg := strings.TrimSpace(wle.Message)
		if msg == "" {
			return "X rejected this login. Please check the details and try again."
		}
		return fmt.Sprintf("X rejected this login: %s", msg)
	}
}

type WebLoginChallenge struct {
	SubtaskID   string
	Hint        string
	Description string
	IsTwoFactor bool
}

type WebLoginResult struct {
	Status           WebLoginStatus
	Challenge        *WebLoginChallenge
	CurrentSubtaskID string
}

type WebLoginSession struct {
	client    *Client
	flowToken string
	subtasks  []onboardingSubtask
	backend   webLoginBackend
	jetfuel   *jetfuelLoginState
}

func NewWebLoginSession(client *Client) *WebLoginSession {
	return &WebLoginSession{client: client}
}

func (wls *WebLoginSession) Client() *Client {
	return wls.client
}

func (wls *WebLoginSession) Start(ctx context.Context) (*WebLoginResult, error) {
	if result, err := wls.startJetfuel(ctx); err == nil {
		return result, nil
	} else {
		wls.client.Logger.Warn().Err(err).Msg("Jetfuel login start failed, falling back to OCF login")
	}
	return wls.startOCF(ctx)
}

func (wls *WebLoginSession) startOCF(ctx context.Context) (*WebLoginResult, error) {
	wls.backend = webLoginBackendOCF
	if err := wls.client.loadPage(ctx, endpoints.BASE_FLOW_LOGIN_URL); err != nil {
		return nil, fmt.Errorf("failed to load x login page: %w", err)
	}
	if err := wls.client.activateGuest(ctx); err != nil {
		return nil, err
	}

	startPayload := newWebLoginStartPayload(wls.client.session.Country)
	resp, err := wls.client.sendOnboardingTask(ctx, endpoints.ONBOARDING_LOGIN_TASK_URL, startPayload)
	if err != nil {
		return nil, err
	}
	if err := wls.update(resp); err != nil {
		return nil, err
	}
	return wls.advanceJSInstrumentation(ctx)
}

func newWebLoginStartPayload(countryCode string) onboardingTaskRequest {
	return onboardingTaskRequest{
		InputFlowData: &onboardingInputFlowData{
			FlowContext: onboardingFlowContext{
				DebugOverrides: map[string]any{},
				StartLocation:  map[string]string{"location": "manual_link"},
			},
			CountryCode: countryCode,
		},
		SubtaskVersions: webLoginSubtaskVersions(),
	}
}

func (wls *WebLoginSession) SubmitIdentifier(ctx context.Context, identifier string) (*WebLoginResult, error) {
	if wls.backend == webLoginBackendJetfuel {
		return wls.submitJetfuelIdentifier(ctx, identifier)
	}
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, fmt.Errorf("x username, email, or phone is required")
	}
	st := wls.currentSubtask()
	if st == nil || st.SettingsList == nil {
		return nil, fmt.Errorf("%w: expected identifier settings_list, got %s", ErrWebLoginUnexpectedSubtask, subtaskName(st))
	}
	link := st.SettingsList.nextLinkID()
	payload := onboardingTaskRequest{
		FlowToken: wls.flowToken,
		SubtaskInputs: []onboardingSubtaskInput{{
			SubtaskID: st.SubtaskID,
			SettingsList: &settingsListInput{
				SettingResponses: []settingResponseInput{{
					Key: "user_identifier",
					ResponseData: map[string]resultInput{
						"text_data": {Result: identifier},
					},
				}},
				Link: link,
			},
		}},
	}
	resp, err := wls.client.sendOnboardingTask(ctx, endpoints.ONBOARDING_TASK_URL, payload)
	if err != nil {
		return nil, err
	}
	if err := wls.update(resp); err != nil {
		return nil, err
	}
	return wls.advanceJSInstrumentation(ctx)
}

func (wls *WebLoginSession) SubmitCredentials(ctx context.Context, identifier, password string) (*WebLoginResult, error) {
	if wls.backend == webLoginBackendJetfuel {
		return wls.submitJetfuelCredentials(ctx, identifier, password)
	}
	result, err := wls.SubmitIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if result.Status != WebLoginStatusNeedsPassword {
		return result, nil
	}
	return wls.SubmitPassword(ctx, password)
}

func (wls *WebLoginSession) SubmitPassword(ctx context.Context, password string) (*WebLoginResult, error) {
	if wls.backend == webLoginBackendJetfuel {
		return wls.submitJetfuelPassword(ctx, password)
	}
	if password == "" {
		return nil, fmt.Errorf("x password is required")
	}
	st := wls.currentSubtask()
	if st == nil || st.EnterPassword == nil {
		return nil, fmt.Errorf("%w: expected enter_password, got %s", ErrWebLoginUnexpectedSubtask, subtaskName(st))
	}
	link := st.EnterPassword.nextLinkID()
	payload := onboardingTaskRequest{
		FlowToken: wls.flowToken,
		SubtaskInputs: []onboardingSubtaskInput{{
			SubtaskID: st.SubtaskID,
			EnterPassword: &enterPasswordInput{
				Password: password,
				Link:     link,
			},
		}},
	}
	resp, err := wls.client.sendOnboardingTask(ctx, endpoints.ONBOARDING_TASK_URL, payload)
	if err != nil {
		return nil, err
	}
	if err := wls.update(resp); err != nil {
		return nil, err
	}
	return wls.advanceJSInstrumentation(ctx)
}

func (wls *WebLoginSession) SubmitText(ctx context.Context, text string) (*WebLoginResult, error) {
	if wls.backend == webLoginBackendJetfuel {
		return wls.submitJetfuelText(ctx, text)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("x verification code is required")
	}
	st := wls.currentSubtask()
	if st == nil || st.EnterText == nil {
		return nil, fmt.Errorf("%w: expected enter_text, got %s", ErrWebLoginUnexpectedSubtask, subtaskName(st))
	}
	link := st.EnterText.nextLinkID()
	payload := onboardingTaskRequest{
		FlowToken: wls.flowToken,
		SubtaskInputs: []onboardingSubtaskInput{{
			SubtaskID: st.SubtaskID,
			EnterText: &enterTextInput{
				Text: text,
				Link: link,
			},
		}},
	}
	resp, err := wls.client.sendOnboardingTask(ctx, endpoints.ONBOARDING_TASK_URL, payload)
	if err != nil {
		return nil, err
	}
	if err := wls.update(resp); err != nil {
		return nil, err
	}
	return wls.advanceJSInstrumentation(ctx)
}

func (wls *WebLoginSession) advanceJSInstrumentation(ctx context.Context) (*WebLoginResult, error) {
	for range 3 {
		st := wls.currentSubtask()
		if st == nil || st.JSInstrumentation == nil {
			return wls.result(), nil
		}
		metrics := "{}"
		payload := onboardingTaskRequest{
			FlowToken: wls.flowToken,
			SubtaskInputs: []onboardingSubtaskInput{{
				SubtaskID: st.SubtaskID,
				JSInstrumentation: &jsInstrumentationInput{
					Response: metrics,
					Link:     webLoginLinkNext,
				},
			}},
		}
		resp, err := wls.client.sendOnboardingTask(ctx, endpoints.ONBOARDING_TASK_URL, payload)
		if err != nil {
			return nil, err
		}
		if err := wls.update(resp); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("%w: JS instrumentation loop did not settle", ErrWebLoginUnexpectedSubtask)
}

func (wls *WebLoginSession) update(resp *onboardingTaskResponse) error {
	if resp == nil {
		return fmt.Errorf("x login response was empty")
	}
	if resp.FlowToken == "" && !wls.client.IsLoggedIn() {
		return ErrWebLoginMissingFlowToken
	}
	if resp.FlowToken != "" {
		wls.flowToken = resp.FlowToken
	}
	wls.subtasks = resp.Subtasks
	return nil
}

func (wls *WebLoginSession) result() *WebLoginResult {
	if wls.client != nil && wls.client.IsLoggedIn() {
		return &WebLoginResult{Status: WebLoginStatusComplete}
	}
	st := wls.currentSubtask()
	if st == nil {
		return &WebLoginResult{Status: WebLoginStatusUnsupported}
	}
	result := &WebLoginResult{CurrentSubtaskID: st.SubtaskID}
	switch {
	case st.JSInstrumentation != nil:
		result.Status = WebLoginStatusUnsupported
	case st.SettingsList != nil:
		result.Status = WebLoginStatusNeedsIdentifier
		result.Challenge = &WebLoginChallenge{
			SubtaskID:   st.SubtaskID,
			Hint:        st.SettingsList.identifierHint(),
			Description: richTextText(st.SettingsList.DetailText),
		}
	case st.EnterPassword != nil:
		result.Status = WebLoginStatusNeedsPassword
		result.Challenge = &WebLoginChallenge{
			SubtaskID:   st.SubtaskID,
			Hint:        st.EnterPassword.hint(),
			Description: richTextText(st.EnterPassword.SecondaryText),
		}
	case st.EnterText != nil:
		result.Status = WebLoginStatusNeedsText
		result.Challenge = &WebLoginChallenge{
			SubtaskID:   st.SubtaskID,
			Hint:        st.EnterText.hint(),
			Description: richTextText(st.EnterText.DetailText),
			IsTwoFactor: st.SubtaskID == webLoginSubtaskTwoFactor,
		}
	case st.OpenAccount != nil || st.OpenHomeTimeline != nil || st.EndFlow != nil:
		result.Status = WebLoginStatusComplete
	default:
		result.Status = WebLoginStatusUnsupported
	}
	return result
}

func (wls *WebLoginSession) currentSubtask() *onboardingSubtask {
	for i := range wls.subtasks {
		st := &wls.subtasks[i]
		if st.JSInstrumentation != nil || st.SettingsList != nil || st.EnterPassword != nil ||
			st.EnterText != nil || st.OpenAccount != nil || st.OpenHomeTimeline != nil || st.EndFlow != nil {
			return st
		}
	}
	if len(wls.subtasks) == 0 {
		return nil
	}
	return &wls.subtasks[0]
}

func subtaskName(st *onboardingSubtask) string {
	if st == nil {
		return "<none>"
	}
	if st.SubtaskID == "" {
		return "<unnamed>"
	}
	return st.SubtaskID
}

func (c *Client) activateGuest(ctx context.Context) error {
	resp, respBody, err := c.MakeRequest(ctx, endpoints.GUEST_ACTIVATE_URL, http.MethodPost, c.buildHeaders(HeaderOpts{
		WithNonAuthBearer:   true,
		WithCookies:         true,
		WithXTwitterHeaders: true,
		Origin:              endpoints.BASE_URL,
		Referer:             endpoints.BASE_FLOW_LOGIN_URL,
	}), []byte(`{}`), types.ContentTypeJSON)
	if resp != nil {
		c.cookies.UpdateFromResponse(resp)
	}
	if err != nil {
		return fmt.Errorf("failed to activate X guest session: %w", err)
	}
	var guest struct {
		GuestToken string `json:"guest_token"`
	}
	if err = json.Unmarshal(respBody, &guest); err != nil {
		return fmt.Errorf("failed to parse X guest activation response: %w", err)
	}
	if guest.GuestToken == "" {
		return ErrWebLoginMissingGuestToken
	}
	c.cookies.Set(cookies.XGuestToken, guest.GuestToken)
	return nil
}

func (c *Client) sendOnboardingTask(ctx context.Context, url string, payload onboardingTaskRequest) (*onboardingTaskResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode X login task: %w", err)
	}
	txID, err := crypto.SignTransaction(c.session.AnimationToken, c.session.VerificationToken, url, http.MethodPost)
	if err != nil {
		c.Logger.Trace().Err(err).Msg("Failed to create X login client transaction ID")
		txID = "e:"
	}
	resp, respBody, err := c.makeRequestDirect(ctx, url, http.MethodPost, c.buildHeaders(HeaderOpts{
		WithNonAuthBearer:   true,
		WithCookies:         true,
		WithXTwitterHeaders: true,
		WithXGuestToken:     true,
		Origin:              endpoints.BASE_URL,
		Referer:             endpoints.BASE_FLOW_LOGIN_URL,
		Extra: map[string]string{
			"x-client-transaction-id": txID,
			"accept":                  "*/*",
			"sec-fetch-dest":          "empty",
			"sec-fetch-mode":          "cors",
			"sec-fetch-site":          "same-site",
		},
	}), body, types.ContentTypeJSON)
	if resp != nil {
		c.cookies.UpdateFromResponse(resp)
	}
	taskResp := &onboardingTaskResponse{}
	if len(respBody) > 0 {
		if unmarshalErr := json.Unmarshal(respBody, taskResp); unmarshalErr != nil && err == nil {
			return nil, fmt.Errorf("failed to parse X login task response: %w", unmarshalErr)
		}
	}
	if len(taskResp.Errors) > 0 {
		return taskResp, &WebLoginError{
			Code:    taskResp.Errors[0].Code,
			Message: taskResp.Errors[0].Message,
		}
	}
	if err != nil {
		return taskResp, err
	}
	return taskResp, nil
}

type onboardingTaskRequest struct {
	FlowToken       string                   `json:"flow_token,omitempty"`
	InputFlowData   *onboardingInputFlowData `json:"input_flow_data,omitempty"`
	SubtaskVersions map[string]int           `json:"subtask_versions,omitempty"`
	SubtaskInputs   []onboardingSubtaskInput `json:"subtask_inputs,omitempty"`
}

type onboardingInputFlowData struct {
	FlowContext onboardingFlowContext `json:"flow_context"`
	CountryCode string                `json:"country_code,omitempty"`
}

type onboardingFlowContext struct {
	DebugOverrides map[string]any    `json:"debug_overrides"`
	StartLocation  map[string]string `json:"start_location"`
}

type onboardingSubtaskInput struct {
	SubtaskID         string                  `json:"subtask_id"`
	JSInstrumentation *jsInstrumentationInput `json:"js_instrumentation,omitempty"`
	SettingsList      *settingsListInput      `json:"settings_list,omitempty"`
	EnterPassword     *enterPasswordInput     `json:"enter_password,omitempty"`
	EnterText         *enterTextInput         `json:"enter_text,omitempty"`
}

type jsInstrumentationInput struct {
	Response string `json:"response"`
	Link     string `json:"link"`
}

type settingsListInput struct {
	SettingResponses []settingResponseInput `json:"setting_responses"`
	Link             string                 `json:"link"`
	CastleToken      string                 `json:"castle_token,omitempty"`
}

type settingResponseInput struct {
	Key          string                 `json:"key"`
	ResponseData map[string]resultInput `json:"response_data"`
}

type resultInput struct {
	Result any `json:"result"`
}

type enterPasswordInput struct {
	Password    string `json:"password"`
	Link        string `json:"link"`
	CastleToken string `json:"castle_token,omitempty"`
}

type enterTextInput struct {
	Text        string `json:"text"`
	Link        string `json:"link"`
	CastleToken string `json:"castle_token,omitempty"`
}

type onboardingTaskResponse struct {
	FlowToken string              `json:"flow_token"`
	Subtasks  []onboardingSubtask `json:"subtasks"`
	Errors    []TwitterError      `json:"errors"`
}

type onboardingSubtask struct {
	SubtaskID         string                    `json:"subtask_id"`
	JSInstrumentation *jsInstrumentationSubtask `json:"js_instrumentation,omitempty"`
	SettingsList      *settingsListSubtask      `json:"settings_list,omitempty"`
	EnterPassword     *enterPasswordSubtask     `json:"enter_password,omitempty"`
	EnterText         *enterTextSubtask         `json:"enter_text,omitempty"`
	OpenAccount       *struct{}                 `json:"open_account,omitempty"`
	OpenHomeTimeline  *struct{}                 `json:"open_home_timeline,omitempty"`
	EndFlow           *struct{}                 `json:"end_flow,omitempty"`
}

type jsInstrumentationSubtask struct {
	URL string `json:"url"`
}

type settingsListSubtask struct {
	Settings   []settingsListSetting `json:"settings"`
	NextLink   *navigationLink       `json:"next_link,omitempty"`
	DetailText *richText             `json:"detail_text,omitempty"`
}

type settingsListSetting struct {
	ValueIdentifier string            `json:"value_identifier"`
	ValueType       string            `json:"value_type"`
	ValueData       *settingValueData `json:"value_data,omitempty"`
}

type settingValueData struct {
	TextField *textFieldData `json:"text_field,omitempty"`
	Button    *buttonData    `json:"button,omitempty"`
}

type textFieldData struct {
	HintText string `json:"hint_text"`
}

type buttonData struct {
	NavigationLink *navigationLink `json:"navigation_link,omitempty"`
}

type enterPasswordSubtask struct {
	Hint          string          `json:"hint"`
	PasswordField *textFieldData  `json:"password_field,omitempty"`
	NextLink      *navigationLink `json:"next_link,omitempty"`
	SecondaryText *richText       `json:"secondary_text,omitempty"`
}

type enterTextSubtask struct {
	HintText   string          `json:"hint_text"`
	DetailText *richText       `json:"detail_text,omitempty"`
	NextLink   *navigationLink `json:"next_link,omitempty"`
}

type navigationLink struct {
	LinkID string `json:"link_id"`
}

type richText struct {
	Text string `json:"text"`
}

func (sls *settingsListSubtask) nextLinkID() string {
	if sls == nil {
		return webLoginLinkNext
	}
	if sls.NextLink != nil && sls.NextLink.LinkID != "" {
		return sls.NextLink.LinkID
	}
	for _, setting := range sls.Settings {
		if setting.ValueIdentifier == "next_button" && setting.ValueData != nil &&
			setting.ValueData.Button != nil && setting.ValueData.Button.NavigationLink != nil &&
			setting.ValueData.Button.NavigationLink.LinkID != "" {
			return setting.ValueData.Button.NavigationLink.LinkID
		}
	}
	return webLoginLinkNext
}

func (sls *settingsListSubtask) identifierHint() string {
	if sls == nil {
		return ""
	}
	for _, setting := range sls.Settings {
		if setting.ValueIdentifier == "user_identifier" && setting.ValueData != nil && setting.ValueData.TextField != nil {
			return setting.ValueData.TextField.HintText
		}
	}
	return ""
}

func (eps *enterPasswordSubtask) nextLinkID() string {
	if eps != nil && eps.NextLink != nil && eps.NextLink.LinkID != "" {
		return eps.NextLink.LinkID
	}
	return webLoginLinkNext
}

func (eps *enterPasswordSubtask) hint() string {
	if eps == nil {
		return ""
	}
	if eps.PasswordField != nil && eps.PasswordField.HintText != "" {
		return eps.PasswordField.HintText
	}
	return eps.Hint
}

func (ets *enterTextSubtask) nextLinkID() string {
	if ets != nil && ets.NextLink != nil && ets.NextLink.LinkID != "" {
		return ets.NextLink.LinkID
	}
	return webLoginLinkNext
}

func (ets *enterTextSubtask) hint() string {
	if ets == nil {
		return ""
	}
	return ets.HintText
}

func richTextText(rt *richText) string {
	if rt == nil {
		return ""
	}
	return strings.TrimSpace(rt.Text)
}

func webLoginSubtaskVersions() map[string]int {
	return map[string]int{
		"action_list":                          2,
		"alert_dialog":                         1,
		"app_download_cta":                     1,
		"check_logged_in_account":              1,
		"choice_selection":                     3,
		"contacts_live_sync_permission_prompt": 0,
		"cta":                                  7,
		"email_verification":                   2,
		"end_flow":                             1,
		"enter_date":                           1,
		"enter_email":                          2,
		"enter_password":                       5,
		"enter_phone":                          2,
		"enter_recaptcha":                      1,
		"enter_text":                           5,
		"enter_username":                       2,
		"generic_urt":                          3,
		"in_app_notification":                  1,
		"interest_picker":                      3,
		"js_instrumentation":                   1,
		"menu_dialog":                          1,
		"notifications_permission_prompt":      2,
		"open_account":                         2,
		"open_home_timeline":                   1,
		"open_link":                            1,
		"phone_verification":                   4,
		"privacy_options":                      1,
		"security_key":                         3,
		"select_avatar":                        4,
		"select_banner":                        2,
		"settings_list":                        7,
		"show_code":                            1,
		"sign_up":                              2,
		"sign_up_review":                       4,
		"tweet_selection_urt":                  1,
		"update_users":                         1,
		"upload_media":                         1,
		"user_recommendations_list":            4,
		"user_recommendations_urt":             1,
		"wait_spinner":                         3,
		"web_modal":                            1,
	}
}
