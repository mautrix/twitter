package connector

import (
	_ "embed"
	"strings"
	"text/template"

	up "go.mau.fi/util/configupgrade"
	"gopkg.in/yaml.v3"
)

//go:embed example-config.yaml
var ExampleConfig string

type Config struct {
	Proxy       string `yaml:"proxy"`
	GetProxyURL string `yaml:"get_proxy_url"`

	DisplaynameTemplate string `yaml:"displayname_template"`

	ConversationSyncLimit int `yaml:"conversation_sync_limit"`

	X bool `yaml:"x"`

	displaynameTemplate *template.Template `yaml:"-"`
}

type umConfig Config

func (c *Config) UnmarshalYAML(node *yaml.Node) error {
	err := node.Decode((*umConfig)(c))
	if err != nil {
		return err
	}
	return c.PostProcess()
}

func (c *Config) PostProcess() error {
	var err error
	c.displaynameTemplate, err = template.New("displayname").Parse(c.DisplaynameTemplate)
	return err
}

func upgradeConfig(helper up.Helper) {
	helper.Copy(up.Str|up.Null, "proxy")
	helper.Copy(up.Str|up.Null, "get_proxy_url")
	helper.Copy(up.Str, "displayname_template")
	helper.Copy(up.Int, "conversation_sync_limit")
	helper.Copy(up.Bool, "x")
}

type DisplaynameParams struct {
	Username    string
	DisplayName string
}

func (c *Config) FormatDisplayname(username string, displayname string) string {
	var nameBuf strings.Builder
	err := c.displaynameTemplate.Execute(&nameBuf, &DisplaynameParams{
		Username:    username,
		DisplayName: displayname,
	})
	if err != nil {
		panic(err)
	}
	return nameBuf.String()
}

func (tc *TwitterConnector) GetConfig() (string, any, up.Upgrader) {
	return ExampleConfig, &tc.Config, &up.StructUpgrader{
		SimpleUpgrader: up.SimpleUpgrader(upgradeConfig),
		Base:           ExampleConfig,
	}
}
