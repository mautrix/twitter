// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2025 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import "maunium.net/go/mautrix/bridgev2/status"

func init() {
	status.BridgeStateHumanErrors.Update(status.BridgeStateErrorMap{
		"twitter-not-logged-in":       "You're not logged into X. Please log in again.",
		"twitter-migration-reauth":    "X Chat needs a passcode update. Please log in and enter your passcode.",
		"twitter-invalid-credentials": "Your X session has expired. Please log in again.",
		"twitter-load-error":          "Couldn't connect to X.",
		"twitter-xchat-fetch-error":   "Couldn't sync X Chat.",
	})
}
