package connector

import (
	"os"
	"path/filepath"
	"testing"

	up "go.mau.fi/util/configupgrade"
	"gopkg.in/yaml.v3"
)

func TestExampleConfigDisablesNativeLoginByDefault(t *testing.T) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(ExampleConfig), &raw); err != nil {
		t.Fatalf("failed to parse example config: %v", err)
	}
	value, ok := raw["native_login"]
	if !ok {
		t.Fatal("example config is missing native_login")
	}
	if value != false {
		t.Fatalf("native_login = %#v, want false", value)
	}

	var config Config
	if err := yaml.Unmarshal([]byte(ExampleConfig), &config); err != nil {
		t.Fatalf("failed to unmarshal example config: %v", err)
	}
	if config.NativeLogin {
		t.Fatal("Config.NativeLogin = true, want false")
	}
}

func TestConfigEnablesNativeLogin(t *testing.T) {
	var config Config
	if err := yaml.Unmarshal([]byte("native_login: true\n"), &config); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	if !config.NativeLogin {
		t.Fatal("Config.NativeLogin = false, want true")
	}
}

func TestConfigUpgradeHandlesNativeLogin(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{name: "missing defaults false", data: "x: true\n", want: false},
		{name: "explicit true is preserved", data: "native_login: true\n", want: true},
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
			if config.NativeLogin != test.want {
				t.Fatalf("Config.NativeLogin = %t, want %t", config.NativeLogin, test.want)
			}
		})
	}
}
