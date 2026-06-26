//go:build liveprobe

package twittermeow

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/cookies"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func TestLiveJetfuelLoginLandingProbe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	if err := client.loadPage(ctx, jetfuelProbePageURL()); err != nil {
		t.Fatalf("loadPage() failed: %v", err)
	}
	if os.Getenv("TWITTER_DNT_PROBE") == "1" {
		client.cookies.Set(cookies.XCookieName("dnt"), "1")
		t.Logf("using injected dnt cookie: true")
	}
	if cuid := strings.TrimSpace(os.Getenv("TWITTER_CUID_PROBE")); cuid != "" {
		client.cookies.Set(cookies.XCookieName("__cuid"), cuid)
		t.Logf("using injected __cuid cookie: true")
	}
	t.Logf("transaction token state after login page: verification=%t animation=%t",
		client.session.VerificationToken != "",
		client.session.AnimationToken != "",
	)
	t.Logf("cookie presence after login page: cf_bm=%t cuid=%t dnt=%t guest_id=%t guest_id_ads=%t guest_id_marketing=%t personalization_id=%t gt=%t ct0=%t att=%t",
		!client.cookies.IsCookieEmpty(cookies.XCookieName("__cf_bm")),
		!client.cookies.IsCookieEmpty(cookies.XCookieName("__cuid")),
		!client.cookies.IsCookieEmpty(cookies.XCookieName("dnt")),
		!client.cookies.IsCookieEmpty(cookies.XGuestID),
		!client.cookies.IsCookieEmpty(cookies.XGuestIDAds),
		!client.cookies.IsCookieEmpty(cookies.XGuestIDMarketing),
		!client.cookies.IsCookieEmpty(cookies.XPersonalizationID),
		!client.cookies.IsCookieEmpty(cookies.XGuestToken),
		!client.cookies.IsCookieEmpty(cookies.XCt0),
		!client.cookies.IsCookieEmpty(cookies.XAtt),
	)
	if _, err := client.jetfuelGet(ctx, endpoints.JETFUEL_LANDING_PATH); err != nil {
		t.Logf("landing preflight failed: %v", err)
	}
	t.Logf("cookie presence after Jetfuel landing: cf_bm=%t cuid=%t dnt=%t guest_id=%t guest_id_ads=%t guest_id_marketing=%t personalization_id=%t gt=%t ct0=%t att=%t",
		!client.cookies.IsCookieEmpty(cookies.XCookieName("__cf_bm")),
		!client.cookies.IsCookieEmpty(cookies.XCookieName("__cuid")),
		!client.cookies.IsCookieEmpty(cookies.XCookieName("dnt")),
		!client.cookies.IsCookieEmpty(cookies.XGuestID),
		!client.cookies.IsCookieEmpty(cookies.XGuestIDAds),
		!client.cookies.IsCookieEmpty(cookies.XGuestIDMarketing),
		!client.cookies.IsCookieEmpty(cookies.XPersonalizationID),
		!client.cookies.IsCookieEmpty(cookies.XGuestToken),
		!client.cookies.IsCookieEmpty(cookies.XCt0),
		!client.cookies.IsCookieEmpty(cookies.XAtt),
	)

	body, err := client.jetfuelGet(ctx, endpoints.JETFUEL_LOGIN_PATH)
	if err != nil {
		t.Fatalf("jetfuelGet() failed: %v", err)
	}
	parsed := parseJetfuelLoginResponse(body)
	t.Logf("Jetfuel landing strings=%d paths=%v fields=%v", len(parsed.strings), parsed.paths, parsed.fields)
	if !parsed.hasPath(endpoints.JETFUEL_BEGIN_LOGIN_PATH) && !parsed.hasField("username_or_email") {
		t.Fatalf("Jetfuel landing did not expose begin_login or username_or_email")
	}
}

func TestLiveJetfuelIdentifierMetadataProbe(t *testing.T) {
	identifier := strings.TrimSpace(os.Getenv("TWITTER_IDENTIFIER_PROBE"))
	if identifier == "" {
		t.Skip("TWITTER_IDENTIFIER_PROBE is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	if err := client.loadPage(ctx, jetfuelProbePageURL()); err != nil {
		t.Fatalf("loadPage() failed: %v", err)
	}
	if os.Getenv("TWITTER_DNT_PROBE") == "1" {
		client.cookies.Set(cookies.XCookieName("dnt"), "1")
		t.Logf("using injected dnt cookie: true")
	}
	if cuid := strings.TrimSpace(os.Getenv("TWITTER_CUID_PROBE")); cuid != "" {
		client.cookies.Set(cookies.XCookieName("__cuid"), cuid)
		t.Logf("using injected __cuid cookie: true")
	}
	t.Logf("transaction token state after login page: verification=%t animation=%t",
		client.session.VerificationToken != "",
		client.session.AnimationToken != "",
	)
	if _, err := client.jetfuelGet(ctx, endpoints.JETFUEL_LANDING_PATH); err != nil {
		t.Logf("landing preflight failed: %v", err)
	}
	if _, err := client.jetfuelGet(ctx, endpoints.JETFUEL_LOGIN_PATH); err != nil {
		t.Fatalf("login graph GET failed: %v", err)
	}

	form := url.Values{
		"username_or_email": {identifier},
	}
	if password := os.Getenv("TWITTER_PASSWORD_PROBE"); password != "" {
		form.Set("password", password)
		t.Logf("using injected password field: true")
	}
	if os.Getenv("TWITTER_CURRENT_CASTLE_TOKEN_PROBE") == "1" {
		castleToken, err := createCurrentCastleRequestToken("091a1f32-3826-4bad-9250-aa14e3c0a2b2")
		if err != nil {
			t.Fatalf("createCurrentCastleRequestToken() failed: %v", err)
		}
		form.Set("$castle_token", castleToken)
		t.Logf("using generated current Castle token: true")
	}
	if castleToken := strings.TrimSpace(os.Getenv("TWITTER_CASTLE_TOKEN_PROBE")); castleToken != "" {
		form.Set("$castle_token", castleToken)
		t.Logf("using injected castle token: true")
	}
	body, err := client.jetfuelPostForm(ctx, endpoints.JETFUEL_BEGIN_LOGIN_PATH, form)
	if err != nil {
		t.Fatalf("begin_login POST failed: %v", err)
	}
	parsed := parseJetfuelLoginResponse(body)
	t.Logf("begin_login strings=%d paths=%v fields=%v", len(parsed.strings), parsed.paths, redactJetfuelDebugList(parsed.fields, identifier))
	for _, line := range filteredJetfuelDebugStrings(parsed, identifier) {
		t.Logf("begin_login string: %s", line)
	}
}

func TestLiveJetfuelVerificationResponseProbe(t *testing.T) {
	identifier := strings.TrimSpace(os.Getenv("TWITTER_LIVE_IDENTIFIER"))
	password := os.Getenv("TWITTER_LIVE_PASSWORD")
	code := strings.TrimSpace(os.Getenv("TWITTER_LIVE_VERIFICATION_CODE"))
	if identifier == "" || password == "" || code == "" {
		t.Skip("TWITTER_LIVE_IDENTIFIER, TWITTER_LIVE_PASSWORD, and TWITTER_LIVE_VERIFICATION_CODE are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	session := NewWebLoginSession(client)
	result, err := session.Start(ctx)
	logWebLoginStage(t, "start", result, err)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	result, err = session.SubmitCredentials(ctx, identifier, password)
	logWebLoginStage(t, "credentials", result, err)
	if err != nil {
		t.Fatalf("SubmitCredentials() failed: %v", err)
	}
	if result.Status == WebLoginStatusNeedsPassword {
		result, err = session.SubmitPassword(ctx, password)
		logWebLoginStage(t, "password", result, err)
		if err != nil {
			t.Fatalf("SubmitPassword() failed: %v", err)
		}
	}
	if result.Status != WebLoginStatusNeedsText {
		t.Fatalf("SubmitCredentials() status = %s, want %s", result.Status, WebLoginStatusNeedsText)
	}
	t.Logf("verification action present=%t fields=%v", session.jetfuel != nil && session.jetfuel.verificationAction != "", session.jetfuelVerificationFields())

	form := url.Values{}
	for _, field := range session.jetfuelVerificationFields() {
		form.Set(field, code)
	}
	if session.jetfuel.sessionToken != "" {
		form.Set("session_token", session.jetfuel.sessionToken)
	}
	if session.jetfuel.preludeDispatchID != "" {
		form.Set("prelude_dispatch_id", session.jetfuel.preludeDispatchID)
	}
	body, err := client.jetfuelPostForm(ctx, session.jetfuel.verificationAction, form)
	if err != nil {
		t.Fatalf("verification POST failed: %v", err)
	}
	parsed := parseJetfuelLoginResponse(body)
	t.Logf("verification response: logged_in=%t complete=%t strings=%d paths=%v fields=%v password_action=%q begin_2fa_action=%q verification_action=%q",
		client.IsLoggedIn(),
		parsed.isComplete(),
		len(parsed.strings),
		redactJetfuelDebugList(parsed.paths, identifier, password, code),
		redactJetfuelDebugList(parsed.fields, identifier, password, code),
		parsed.passwordAction(),
		parsed.beginTwoFactorAction(),
		parsed.verificationAction(),
	)
	if err := parsed.loginError(); err != nil {
		var webErr *WebLoginError
		if errors.As(err, &webErr) {
			t.Logf("verification response login error: code=%d message=%q", webErr.Code, redactJetfuelDebugString(webErr.Message, identifier, password, code))
		} else {
			t.Logf("verification response login error: %T", err)
		}
	}
	for _, line := range filteredJetfuelDebugStrings(parsed, identifier, password, code) {
		t.Logf("verification response string: %s", line)
	}
}

func TestLiveJetfuelPasswordResponseProbe(t *testing.T) {
	identifier := strings.TrimSpace(os.Getenv("TWITTER_LIVE_IDENTIFIER"))
	password := os.Getenv("TWITTER_LIVE_PASSWORD")
	if identifier == "" || password == "" {
		t.Skip("TWITTER_LIVE_IDENTIFIER and TWITTER_LIVE_PASSWORD are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	session := NewWebLoginSession(client)
	result, err := session.Start(ctx)
	logWebLoginStage(t, "start", result, err)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	result, err = session.SubmitIdentifier(ctx, identifier)
	logWebLoginStage(t, "identifier", result, err)
	if err != nil {
		t.Fatalf("SubmitIdentifier() failed: %v", err)
	}
	if result.Status != WebLoginStatusNeedsPassword {
		t.Fatalf("SubmitIdentifier() status = %s, want %s", result.Status, WebLoginStatusNeedsPassword)
	}

	form := url.Values{"password": {password}}
	if session.jetfuel.identifier != "" {
		form.Set("username", session.jetfuel.identifier)
	}
	if session.jetfuel.sessionToken != "" {
		form.Set("session_token", session.jetfuel.sessionToken)
	}
	body, err := client.jetfuelPostForm(ctx, session.jetfuel.passwordAction, form)
	if err != nil {
		t.Fatalf("password POST failed: %v", err)
	}
	parsed := parseJetfuelLoginResponse(body)
	t.Logf("password response: logged_in=%t complete=%t strings=%d paths=%v fields=%v password_action=%q begin_2fa_action=%q verification_action=%q",
		client.IsLoggedIn(),
		parsed.isComplete(),
		len(parsed.strings),
		redactJetfuelDebugList(parsed.paths, identifier, password),
		redactJetfuelDebugList(parsed.fields, identifier, password),
		parsed.passwordAction(),
		parsed.beginTwoFactorAction(),
		parsed.verificationAction(),
	)
	if err := parsed.loginError(); err != nil {
		var webErr *WebLoginError
		if errors.As(err, &webErr) {
			t.Logf("password response login error: code=%d message=%q", webErr.Code, redactJetfuelDebugString(webErr.Message, identifier, password))
		} else {
			t.Logf("password response login error: %T", err)
		}
	}
	for _, line := range filteredJetfuelDebugStrings(parsed, identifier, password) {
		t.Logf("password response string: %s", line)
	}
}

func jetfuelProbePageURL() string {
	if os.Getenv("TWITTER_ROOT_PAGE_PROBE") == "1" {
		return endpoints.BASE_URL + "/"
	}
	return endpoints.BASE_FLOW_LOGIN_URL
}

func TestLiveCloudflareJSDProbe(t *testing.T) {
	if os.Getenv("TWITTER_CLOUDFLARE_JSD_PROBE") != "1" {
		t.Skip("TWITTER_CLOUDFLARE_JSD_PROBE=1 is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	extraHeaders := map[string]string{
		"upgrade-insecure-requests": "1",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"sec-fetch-dest":            "document",
	}
	resp, body, err := client.MakeRequest(ctx, endpoints.BASE_FLOW_LOGIN_URL, http.MethodGet, client.buildHeaders(HeaderOpts{Extra: extraHeaders, WithCookies: true}), nil, types.ContentTypeNone)
	if resp != nil {
		client.cookies.UpdateFromResponse(resp)
	}
	if err != nil {
		t.Fatalf("login page request failed: %v", err)
	}
	pageURL := resp.Request.URL
	scriptURL := methods.ParseCloudflareJSDURL(string(body))
	t.Logf("cloudflare jsd script present=%t", scriptURL != "")
	if scriptURL == "" {
		return
	}
	parsedScriptURL, err := pageURL.Parse(scriptURL)
	if err != nil {
		t.Fatalf("parse script URL failed: %v", err)
	}
	originalCheckRedirect := client.HTTP.CheckRedirect
	client.disableRedirects()
	defer func() { client.HTTP.CheckRedirect = originalCheckRedirect }()
	for i := 0; i < 4; i++ {
		resp, _, err = client.MakeRequest(ctx, parsedScriptURL.String(), http.MethodGet, client.buildHeaders(HeaderOpts{
			Extra: map[string]string{
				"accept":         "*/*",
				"sec-fetch-dest": "script",
				"sec-fetch-mode": "no-cors",
				"sec-fetch-site": "same-origin",
			},
			Referer:     pageURL.String(),
			WithCookies: true,
		}), nil, types.ContentTypeNone)
		names := setCookieNames(resp)
		if resp != nil {
			client.cookies.UpdateFromResponse(resp)
			t.Logf("jsd hop=%d status=%d set_cookie_names=%v cf_bm_present=%t", i, resp.StatusCode, names, !client.cookies.IsCookieEmpty(cookies.XCookieName("__cf_bm")))
		}
		if !errors.Is(err, ErrRedirectAttempted) {
			if err != nil {
				t.Fatalf("jsd request failed: %v", err)
			}
			return
		}
		location := resp.Header.Get("Location")
		if location == "" {
			t.Fatalf("redirect without location")
		}
		parsedScriptURL, err = parsedScriptURL.Parse(location)
		if err != nil {
			t.Fatalf("parse redirect URL failed: %v", err)
		}
	}
}

func setCookieNames(resp *http.Response) []string {
	if resp == nil {
		return nil
	}
	names := make([]string, 0)
	for _, cookie := range resp.Cookies() {
		names = append(names, cookie.Name)
	}
	return names
}

func TestLiveJetfuelFakeCredentialsProbe(t *testing.T) {
	if os.Getenv("TWITTER_FAKE_CREDENTIALS_PROBE") != "1" {
		t.Skip("TWITTER_FAKE_CREDENTIALS_PROBE=1 is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	session := NewWebLoginSession(client)
	result, err := session.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	if result.Status != WebLoginStatusNeedsIdentifier {
		t.Fatalf("Start() status = %s, want %s", result.Status, WebLoginStatusNeedsIdentifier)
	}
	result, err = session.SubmitCredentials(ctx, "codex_probe_20260625_noacct", "not-a-real-password")
	if err == nil {
		t.Fatalf("SubmitCredentials() returned result %#v, want missing-account error", result)
	}
	var webErr *WebLoginError
	if !errors.As(err, &webErr) {
		t.Fatalf("SubmitCredentials() error = %T, want *WebLoginError: %v", err, err)
	}
	if webErr.Code != 32 {
		t.Fatalf("WebLoginError.Code = %d, want missing-account/credential code 32 (%s)", webErr.Code, webErr.Message)
	}
	if strings.Contains(strings.ToLower(webErr.Message), "temporarily limited") {
		t.Fatalf("SubmitCredentials() hit temporary-limit branch: %s", webErr.Message)
	}
}

func TestLiveOCFFakeCredentialsProbe(t *testing.T) {
	if os.Getenv("TWITTER_OCF_FAKE_CREDENTIALS_PROBE") != "1" {
		t.Skip("TWITTER_OCF_FAKE_CREDENTIALS_PROBE=1 is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := NewClient(cookies.NewCookies(nil), nil, zerolog.Nop())
	session := NewWebLoginSession(client)
	result, err := session.startOCF(ctx)
	if err != nil {
		t.Fatalf("startOCF() failed: %v", err)
	}
	if result.Status != WebLoginStatusNeedsIdentifier {
		t.Fatalf("startOCF() status = %s, want %s", result.Status, WebLoginStatusNeedsIdentifier)
	}
	result, err = session.SubmitCredentials(ctx, "codex_probe_20260625_noacct", "not-a-real-password")
	if err == nil {
		t.Fatalf("SubmitCredentials() returned result %#v, want credential/missing-account error", result)
	}
	var webErr *WebLoginError
	if !errors.As(err, &webErr) {
		t.Fatalf("SubmitCredentials() error = %T, want *WebLoginError: %v", err, err)
	}
	if strings.Contains(strings.ToLower(webErr.Message), "temporarily limited") {
		t.Fatalf("OCF SubmitCredentials() hit temporary-limit branch: %s", webErr.Message)
	}
	t.Logf("OCF fake credentials error: code=%d message=%q", webErr.Code, webErr.Message)
}

func filteredJetfuelDebugStrings(parsed jetfuelLoginResponse, secrets ...string) []string {
	var out []string
	needles := []string{
		"password",
		"login",
		"verification",
		"challenge",
		"two",
		"factor",
		"error",
		"message",
		"limited",
		"temporarily",
		"castle",
		"action",
		"flow",
		"next",
	}
	for _, str := range parsed.strings {
		lower := strings.ToLower(str)
		matched := false
		for _, needle := range needles {
			if strings.Contains(lower, needle) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		redacted := redactJetfuelDebugString(str, secrets...)
		if len(redacted) > 300 {
			redacted = redacted[:300] + "..."
		}
		out = append(out, redacted)
		if len(out) >= 30 {
			break
		}
	}
	return out
}

func redactJetfuelDebugList(values []string, secrets ...string) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = redactJetfuelDebugString(value, secrets...)
	}
	return out
}

func redactJetfuelDebugString(value string, secrets ...string) string {
	for _, secret := range secrets {
		if secret != "" {
			value = strings.ReplaceAll(value, secret, "<secret>")
		}
	}
	value = uuidDebugRegex.ReplaceAllString(value, "<uuid>")
	return value
}

func logWebLoginStage(t *testing.T, stage string, result *WebLoginResult, err error) {
	t.Helper()
	if result != nil {
		t.Logf("%s result: status=%s subtask=%s", stage, result.Status, result.CurrentSubtaskID)
	}
	if err != nil {
		var webErr *WebLoginError
		if errors.As(err, &webErr) {
			t.Logf("%s error: code=%d message=%q", stage, webErr.Code, webErr.Message)
			return
		}
		t.Logf("%s error: %T", stage, err)
	}
}

var uuidDebugRegex = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
