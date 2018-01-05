package config

import (
	"strings"
	"fmt"
	"io"
	"os"
)

type ConfigType int

const (
	TYPE_HCL ConfigType = iota
	TYPE_YAML
)

// LoadConfigFromFile fills a configuration object (passed as parameter) with
// values read from a configuration file (pass as parameter by filename). The
// configuration file needs to be in HCL format.
func LoadConfigFromFile(config *Config, filename string) error {
	var typ ConfigType

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer reader.Close()

	if strings.HasSuffix(filename, ".hcl") {
		typ = TYPE_HCL
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		typ = TYPE_YAML
	} else {
		return fmt.Errorf("config file '%s' has unsupported file type", filename)
	}

	return LoadConfigFromStream(config, reader, typ)
}

// LoadConfigFromStream fills a configuration object (passed as parameter) with
// values read from a Reader interface (passed as parameter).
func LoadConfigFromStream(config *Config, stream io.Reader, typ ConfigType) error {
	switch typ {
	case TYPE_HCL:
		return loadConfigFromHCLStream(config, stream)
	case TYPE_YAML:
		return loadConfigFromYAMLStream(config, stream)
	default:
		return fmt.Errorf("unsupported config type %d", typ)
	}
}