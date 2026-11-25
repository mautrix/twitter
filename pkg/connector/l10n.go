package connector

import (
	"go.mau.fi/util/l10n"
)

var L = l10n.NewNamespace("mautrix-twitter")

var withErrorComment = l10n.WithComment("error is a message explaining what went wrong; it may be a localized bridge message or untranslated text from X")

var (
	msgFlowCookiesName         = L.Define("login.flow.cookies.name", "Cookies")
	msgFlowCookiesDescription  = L.Define("login.flow.cookies.description", "Log in with your X account using your cookies")
	msgFlowPasswordName        = L.Define("login.flow.password.name", "Username/password")
	msgFlowPasswordDescription = L.Define("login.flow.password.description", "Log in with your X username, email, or phone number and password")
	msgFlowGoogleName          = L.Define("login.flow.google.name", "Google")
	msgFlowGoogleDescription   = L.Define("login.flow.google.description", "Sign in with the Google account linked to your X account")
	msgFlowAppleName           = L.Define("login.flow.apple.name", "Apple")
	msgFlowAppleDescription    = L.Define("login.flow.apple.description", "Sign in with the Apple account linked to your X account")

	msgFieldIdentifierName        = L.Define("login.field.identifier.name", "Username, email, or phone")
	msgFieldIdentifierDescription = L.Define("login.field.identifier.description", "The identifier you use to sign in to X.")
	msgFieldPasswordName          = L.Define("login.field.password.name", "Password")
	msgFieldVerificationName      = L.Define("login.field.verification_code.generic.name", "Verification")
	msgFieldVerificationPhoneName = L.Define("login.field.verification_code.phone.name", "Phone number")
	msgFieldVerificationCodeName  = L.Define("login.field.verification_code.name", "Verification code")
	msgFieldAuthMethodName        = L.Define("login.field.auth_method.name", "Verification method")
	msgFieldAuthMethodDescription = L.Define("login.field.auth_method.description", "Choose the method X should use for this login.")
	msgFieldPINName               = L.Define("login.field.pin.name", "Passcode")
	msgFieldPINSetupName          = L.Define("login.field.pin.setup.name", "Create your PIN code")

	msgStepCookiesInstructions = L.Define("login.step.cookies.instructions", "Open the Login URL in an Incognito/Private browsing mode. Then, extract the cookies as a JSON object/cURL command copied from the Network tab of your browser's DevTools. After that, close the browser **before** pasting the cookies.\n\nFor example: `{\"ct0\":\"123466-...\",\"auth_token\":\"abcde-...\"}`")
	msgStepPrepareInstructions = L.Define("login.step.prepare.instructions", "Preparing X login.")

	msgStepCredentialsInstructions       = L.Define("login.step.credentials.instructions", "Enter your X username, email, or phone number and password.")
	msgStepCredentialsErrorInstructions  = L.Define("login.step.credentials.error.instructions", "{error}\n\nEnter your X username, email, or phone number and password.", withErrorComment)
	msgStepCastleTokenInstructions       = L.Define("login.step.castle_token.instructions", "Generating an X browser token for this login.")
	msgStepCastleTokenErrorInstructions  = L.Define("login.step.castle_token.error.instructions", "{error}\n\nGenerating an X browser token for this login.", withErrorComment)
	msgStepSSOGoogleInstructions         = L.Define("login.step.sso.google.instructions", "Sign in with Google in the window that opens.")
	msgStepSSOGoogleErrorInstructions    = L.Define("login.step.sso.google.error.instructions", "{error}\n\nSign in with Google in the window that opens.", withErrorComment)
	msgStepSSOAppleInstructions          = L.Define("login.step.sso.apple.instructions", "Sign in with Apple in the window that opens.")
	msgStepSSOAppleErrorInstructions     = L.Define("login.step.sso.apple.error.instructions", "{error}\n\nSign in with Apple in the window that opens.", withErrorComment)
	msgStepVerificationInstructions      = L.Define("login.step.verification.instructions", "X needs additional verification for this login.")
	msgStepVerificationErrorInstructions = L.Define("login.step.verification.error.instructions", "{error}\n\nX needs additional verification for this login.", withErrorComment)
	msgStepAuthMethodInstructions        = L.Define("login.step.auth_method.instructions", "Choose how to verify this X login.")
	msgStepAuthMethodErrorInstructions   = L.Define("login.step.auth_method.error.instructions", "{error}\n\nChoose how to verify this X login.", withErrorComment)

	msgStepPINInstructions           = L.Define("login.step.pin.instructions", passcodeBodyRecover)
	msgStepPINErrorInstructions      = L.Define("login.step.pin.error.instructions", "{error}\n\n"+passcodeBodyRecover, withErrorComment)
	msgStepPINSetupInstructions      = L.Define("login.step.pin.setup.instructions", passcodeBodySetup)
	msgStepPINSetupErrorInstructions = L.Define("login.step.pin.setup.error.instructions", "{error}\n\n"+passcodeBodySetup, withErrorComment)

	msgStepCompleteInstructions = L.Define(
		"login.step.complete.instructions", "Successfully logged into X as @{username}",
		l10n.WithComment("username is the X handle of the account, without the @"),
	)

	msgErrClientHTTPFailure        = L.Define("login.error.client_http_failure", "The request did not complete on this device. Please try again.")
	msgErrSessionExpired           = L.Define("login.error.session_expired", "The X login session expired. Enter your X login details again.")
	msgErrCastleMetadataMissing    = L.Define("login.error.castle_metadata_missing", "X did not provide the browser-token metadata needed for native login. Try again.")
	msgErrSSOInvalidResult         = L.Define("login.error.sso_invalid_result", "The sign-in window didn't return a valid result. Try again.")
	msgErrSSOUnavailable           = L.Define("login.error.sso_unavailable", "X didn't offer Google or Apple sign-in. Try again later.")
	msgErrNoLoginStep              = L.Define("login.error.no_login_step", "X did not return a login step. Try again.")
	msgErrNoAuthMethods            = L.Define("login.error.no_auth_methods", "X returned a verification method chooser without any methods. Try again.")
	msgErrNeedsPassword            = L.Define("login.error.needs_password", "X still needs your password. Enter your login details again.")
	msgErrNeedsIdentifier          = L.Define("login.error.needs_identifier", "X still needs your username, email, or phone. Enter your login details again.")
	msgErrAccountIDUnreadable      = L.Define("login.error.account_id_unreadable", "Couldn't read your X account ID. Please enter your login details again.")
	msgErrChooseVerificationMethod = L.Define("login.error.choose_verification_method", "Choose a verification method.")
	msgErrEnterVerificationCode    = L.Define("login.error.enter_verification_code", "Enter the X verification code.")
	msgErrInvalidPasscode          = L.Define("login.error.invalid_passcode", "Invalid passcode.")
	msgErrInvalidPasscodeGuesses   = L.DefinePlural(
		"login.error.invalid_passcode_guesses", "count",
		"Invalid passcode. You have {count} guess remaining.",
		"Invalid passcode. You have {count} guesses remaining.",
	)
	msgErrLoginRejected           = L.Define("login.error.login_rejected", "X rejected this login. Please check the details and try again.")
	msgErrBadCredentials          = L.Define("login.error.bad_credentials", "X rejected the username or password. Check the details and try again.")
	msgErrRateLimited             = L.Define("login.error.rate_limited", "X rejected this login attempt. Wait a bit, then try again.")
	msgErrLoginRejectedWithReason = L.Define(
		"login.error.login_rejected_with_reason", "X rejected this login: {reason}",
		l10n.WithComment("reason is untranslated text from X"),
	)
	msgErrPasswordRejected      = L.Define("login.error.password_rejected", "X didn't accept this login. Check your password and try again.")
	msgErrUnsupportedChallenge  = L.Define("login.error.unsupported_challenge", "X returned a login challenge this bridge does not support yet.")
	msgErrAuthMethodUnavailable = L.Define("login.error.auth_method_unavailable", "That X verification method is not available for this login.")
	msgErrAuthMethodExpired     = L.Define("login.error.auth_method_expired", "The X verification method selection expired. Enter your X login details again.")
	msgErrLoginFailed           = L.Define(
		"login.error.login_failed", "X login failed: {error}",
		l10n.WithComment("error is an untranslated error message from the bridge or from X"),
	)

	msgErrWebviewInvalidToken       = L.Define("login.error.webview_invalid_token", "The X webview returned an invalid browser token.")
	msgErrWebviewNoToken            = L.Define("login.error.webview_no_token", "The X webview did not return a browser token.")
	msgErrWebviewInvalidFingerprint = L.Define("login.error.webview_invalid_fingerprint", "The X webview did not return a valid browser fingerprint.")
	msgErrFreshTokenNeeded          = L.Define("login.error.fresh_token_needed", "X needs a fresh browser token to continue this login.")
)
