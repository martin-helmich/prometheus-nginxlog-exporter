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
	} else {
		return fmt.Errorf("config file '%s' has unsupported file type", filename)
	}

	return LoadConfigFromStream(config, reader, typ)
}

func LoadConfigFromStream(config *Config, stream io.Reader, typ ConfigType) error {
	switch typ {
	case TYPE_HCL:
		return loadConfigFromHCLStream(config, stream)
	default:
		return fmt.Errorf("unsupported config type %d", typ)
	}
}