package config

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// FileFormat describes which kind of configuration file the exporter was started with
type FileFormat int

const (
	// TypeHCL describes the HCL (Hashicorp configuration language) file format
	TypeHCL FileFormat = iota
	// TypeYAML describes the YAML file format
	TypeYAML
)

// LoadConfigFromFile fills a configuration object (passed as parameter) with
// values read from a configuration file (pass as parameter by filename). The
// configuration file needs to be in HCL format.
func LoadConfigFromFile(config *Config, filename string) error {
	var typ FileFormat

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer reader.Close()

	if strings.HasSuffix(filename, ".hcl") {
		typ = TypeHCL
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		typ = TypeYAML
	} else {
		return fmt.Errorf("config file '%s' has unsupported file type", filename)
	}

	return LoadConfigFromStream(config, reader, typ)
}

// LoadConfigFromStream fills a configuration object (passed as parameter) with
// values read from a Reader interface (passed as parameter).
func LoadConfigFromStream(config *Config, stream io.Reader, typ FileFormat) error {
	switch typ {
	case TypeHCL:
		if err := loadConfigFromHCLStream(config, stream); err != nil {
			return err
		}
	case TypeYAML:
		if err := loadConfigFromYAMLStream(config, stream); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported config type %d", typ)
	}

	for i := range config.Namespaces {
		config.Namespaces[i].ResolveDeprecations()
	}

	return nil
}
