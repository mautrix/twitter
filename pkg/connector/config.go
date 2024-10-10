package connector

import (
	_ "embed"

	up "go.mau.fi/util/configupgrade"
)

//go:embed example-config.yaml
var ExampleConfig string

type Config struct {
	Proxy       string `yaml:"proxy"`
	GetProxyURL string `yaml:"get_proxy_url"`

	UsernameTemplate     string `yaml:"username_template"`
	DisplaynameTemplate  string `yaml:"displayname_template"`
	DisplayNameMaxLength int    `yaml:"displayname_max_length"`

	DeliveryReceipts bool `yaml:"delivery_receipts"`

	//displaynameTemplate *template.Template `yaml:"-"`
}

func upgradeConfig(helper up.Helper) {
	helper.Copy(up.Str|up.Null, "proxy")
	helper.Copy(up.Str|up.Null, "get_proxy_url")

	helper.Copy(up.Str, "username_template")
	helper.Copy(up.Str, "displayname_template")
	helper.Copy(up.Int, "displayname_max_length")

	helper.Copy(up.Bool, "delivery_receipts")
}

func (tc *TwitterConnector) GetConfig() (string, any, up.Upgrader) {
	return ExampleConfig, &tc.Config, &up.StructUpgrader{
		SimpleUpgrader: up.SimpleUpgrader(upgradeConfig),
		Base:           ExampleConfig,
	}
}
