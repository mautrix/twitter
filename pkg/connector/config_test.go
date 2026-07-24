package connector

import (
	"os"
	"path/filepath"
	"testing"

	up "go.mau.fi/util/configupgrade"
	"gopkg.in/yaml.v3"
)

func TestExampleConfigUsesWebViewLoginByDefault(t *testing.T) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(ExampleConfig), &raw); err != nil {
		t.Fatalf("failed to parse example config: %v", err)
	}
	value, ok := raw["login_flow"]
	if !ok {
		t.Fatal("example config is missing login_flow")
	}
	if value != string(LoginFlowWebView) {
		t.Fatalf("login_flow = %#v, want %q", value, LoginFlowWebView)
	}

	var config Config
	if err := yaml.Unmarshal([]byte(ExampleConfig), &config); err != nil {
		t.Fatalf("failed to unmarshal example config: %v", err)
	}
	if config.LoginFlow != LoginFlowWebView {
		t.Fatalf("Config.LoginFlow = %q, want %q", config.LoginFlow, LoginFlowWebView)
	}
}

func TestConfigAcceptsLoginFlows(t *testing.T) {
	for _, flow := range []LoginFlow{LoginFlowWebView, LoginFlowNative, LoginFlowClientHTTP} {
		t.Run(string(flow), func(t *testing.T) {
			var config Config
			if err := yaml.Unmarshal([]byte("login_flow: "+string(flow)+"\n"), &config); err != nil {
				t.Fatalf("failed to unmarshal config: %v", err)
			}
			if config.LoginFlow != flow {
				t.Fatalf("Config.LoginFlow = %q, want %q", config.LoginFlow, flow)
			}
		})
	}

	var config Config
	if err := yaml.Unmarshal([]byte("login_flow: proxy_magic\n"), &config); err == nil {
		t.Fatal("invalid login_flow was accepted")
	}
}

func TestConfigUpgradeHandlesLoginFlow(t *testing.T) {
	tests := []struct {
		name string
		data string
		want LoginFlow
	}{
		{name: "missing defaults to webview", data: "x: true\n", want: LoginFlowWebView},
		{name: "webview is preserved", data: "login_flow: webview\n", want: LoginFlowWebView},
		{name: "native is preserved", data: "login_flow: native\n", want: LoginFlowNative},
		{name: "client HTTP is preserved", data: "login_flow: client_http\n", want: LoginFlowClientHTTP},
		{name: "legacy client HTTP true migrates", data: "client_http_login: true\n", want: LoginFlowClientHTTP},
		{name: "legacy client HTTP takes precedence", data: "client_http_login: true\nnative_login: true\n", want: LoginFlowClientHTTP},
		{name: "legacy native true migrates", data: "native_login: true\n", want: LoginFlowNative},
		{name: "legacy native false stays webview", data: "native_login: false\n", want: LoginFlowWebView},
		{name: "legacy client HTTP false stays webview", data: "client_http_login: false\n", want: LoginFlowWebView},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(configPath, []byte(test.data), 0o600); err != nil {
				t.Fatalf("failed to write source config: %v", err)
			}
			upgrader := &up.StructUpgrader{
				SimpleUpgrader: up.SimpleUpgrader(upgradeConfig),
				Base:           ExampleConfig,
			}
			output, _, err := up.Do(configPath, false, upgrader)
			if err != nil {
				t.Fatalf("failed to upgrade config: %v", err)
			}

			var config Config
			if err = yaml.Unmarshal(output, &config); err != nil {
				t.Fatalf("failed to unmarshal upgraded config: %v", err)
			}
			if config.LoginFlow != test.want {
				t.Fatalf("Config.LoginFlow = %q, want %q", config.LoginFlow, test.want)
			}
		})
	}
}
