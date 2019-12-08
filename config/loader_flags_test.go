package config

import (
	"testing"

	"github.com/stretchr/testify/require"
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

	require.Len(t, cfg.Namespaces, 1)
	require.Equal(t, FileSource(sf), cfg.Namespaces[0].SourceData.Files)
}
