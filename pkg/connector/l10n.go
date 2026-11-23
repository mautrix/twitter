package connector

import (
	"go.mau.fi/util/l10n"
)

var L = l10n.NewNamespace("mautrix-twitter")

var withErrorComment = l10n.WithComment("error is an English error message from the bridge or from X")

var (
	msgFlowCookiesName         = L.Define("login.flow.cookies.name", "Cookies")
	msgFlowCookiesDescription  = L.Define("login.flow.cookies.description", "Log in with your X account using your cookies")
	msgFlowPasswordName        = L.Define("login.flow.password.name", "Username/password")
	msgFlowPasswordDescription = L.Define("login.flow.password.description", "Log in with your X username, email, or phone number and password")

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
)
