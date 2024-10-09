package connector

import (
	"go.mau.fi/util/configupgrade"
)

type TwitterConfig struct {
}

func (tc *TwitterConnector) GetConfig() (example string, data any, upgrader configupgrade.Upgrader) {
	return "", nil, configupgrade.NoopUpgrader
}
