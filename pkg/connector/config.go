package connector

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	up "go.mau.fi/util/configupgrade"
	"gopkg.in/yaml.v3"
)

//go:embed example-config.yaml
var ExampleConfig string

type LoginFlow string

const (
	LoginFlowWebView    LoginFlow = "webview"
	LoginFlowNative     LoginFlow = "native"
	LoginFlowClientHTTP LoginFlow = "client_http"
)

type Config struct {
	Proxy       string `yaml:"proxy"`
	GetProxyURL string `yaml:"get_proxy_url"`

	DisplaynameTemplate   string    `yaml:"displayname_template"`
	ConversationSyncLimit int       `yaml:"conversation_sync_limit"`
	CacheSession          bool      `yaml:"cache_session"`
	LoginFlow             LoginFlow `yaml:"login_flow"`

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
	if c.LoginFlow == "" {
		c.LoginFlow = LoginFlowWebView
	}
	switch c.LoginFlow {
	case LoginFlowWebView, LoginFlowNative, LoginFlowClientHTTP:
	default:
		return fmt.Errorf("invalid login_flow %q", c.LoginFlow)
	}
	var err error
	c.displaynameTemplate, err = template.New("displayname").Parse(c.DisplaynameTemplate)
	return err
}

func (c *Config) EffectiveLoginFlow() LoginFlow {
	if c == nil || c.LoginFlow == "" {
		return LoginFlowWebView
	}
	return c.LoginFlow
}

func upgradeConfig(helper up.Helper) {
	helper.Copy(up.Str|up.Null, "proxy")
	helper.Copy(up.Str|up.Null, "get_proxy_url")
	helper.Copy(up.Str, "displayname_template")
	helper.Copy(up.Int, "conversation_sync_limit")
	helper.Copy(up.Bool, "cache_session")
	if _, ok := helper.Get(up.Str, "login_flow"); ok {
		helper.Copy(up.Str, "login_flow")
	} else if clientHTTPLogin, ok := helper.Get(up.Bool, "client_http_login"); ok && clientHTTPLogin == "true" {
		helper.Set(up.Str, string(LoginFlowClientHTTP), "login_flow")
	} else if nativeLogin, ok := helper.Get(up.Bool, "native_login"); ok && nativeLogin == "true" {
		helper.Set(up.Str, string(LoginFlowNative), "login_flow")
	}
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
