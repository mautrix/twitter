// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
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

package main

import (
	"maunium.net/go/mautrix/bridgev2/bridgeconfig"
	"maunium.net/go/mautrix/bridgev2/matrix/mxmain"

	"go.mau.fi/mautrix-twitter/pkg/connector"
)

// Information to find out exactly which commit the bridge was built from.
// These are filled at build time with the -X linker flag.
var (
	Tag       = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
)

var c = &connector.TwitterConnector{}
var m = mxmain.BridgeMain{
	Name:        "mautrix-twitter",
	URL:         "https://github.com/mautrix/twitter",
	Description: "A Matrix-Twitter puppeting bridge.",
	Version:     "0.4.3",
	Connector:   c,
}

func main() {
	bridgeconfig.HackyMigrateLegacyNetworkConfig = migrateLegacyConfig
	m.PostInit = func() {
		m.CheckLegacyDB(
			8,
			"v0.1.0",
			"v0.2.0",
			m.LegacyMigrateSimple(legacyMigrateRenameTables, legacyMigrateCopyData, 18),
			true,
		)
	}
	m.PostStart = func() {
		if m.Matrix.Provisioning != nil {
			m.Matrix.Provisioning.Router.HandleFunc("POST /v1/api/login", legacyProvLogin)
			m.Matrix.Provisioning.Router.HandleFunc("POST /v1/api/logout", legacyProvLogout)
		}
	}

	m.InitVersion(Tag, Commit, BuildTime)
	m.Run()
}
