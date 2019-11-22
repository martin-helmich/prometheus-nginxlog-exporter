package config

import (
	"reflect"
	"testing"
)

func configFromFlags(t *testing.T, flags StartupFlags) Config {
	cfg := Config{}

	err := LoadConfigFromFlags(&cfg, &flags)
	if err != nil {
		t.Error("unexpected error", err)
	}

	return cfg
}

func TestConfigContainsPortFromFlags(t *testing.T) {
	t.Parallel()

	cfg := configFromFlags(t, StartupFlags{
		ListenPort: 1234,
	})

	if cfg.Listen.Port != 1234 {
		t.Error("unexpected listen port", "expected", 1234, "actual", cfg.Listen.Port)
	}
}

func TestConfigContainsFilenamesFromFlags(t *testing.T) {
	t.Parallel()

	sf := []string{"/foo.log", "/bar.log"}
	cfg := configFromFlags(t, StartupFlags{
		ListenPort: 1234,
		Filenames:  sf,
	})

	if len(cfg.Namespaces) != 1 {
		t.Error("unexpected namespace count", "expected", 1, "actual", len(cfg.Namespaces))
	}

	if !reflect.DeepEqual(cfg.Namespaces[0].SourceData.Files, sf) {
		t.Error("unexpected source files", "expected", sf, "actual", cfg.Namespaces[0].SourceData.Files)
	}
}
