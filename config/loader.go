package config

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// LoadConfigFromFile fills a configuration object (passed as parameter) with
// values read from a configuration file (pass as parameter by filename). The
// configuration file needs to be in HCL format.
func LoadConfigFromFile(config *Config, filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	hclText := string(buf)

	err = hcl.Decode(config, hclText)
	if err != nil {
		return err
	}

	fmt.Println(config)

	return nil
}

// LoadConfigFromFlags fills a configuration object (passed as parameter) with
// values from command-line flags.
func LoadConfigFromFlags(config *Config, flags *StartupFlags) error {
	config.Listen = ListenConfig{
		Port:    flags.ListenPort,
		Address: "0.0.0.0",
	}
	config.Namespaces = []NamespaceConfig{
		NamespaceConfig{
			Format:      flags.Format,
			SourceFiles: flags.Filenames,
			Name:        flags.Namespace,
		},
	}

	return nil
}
