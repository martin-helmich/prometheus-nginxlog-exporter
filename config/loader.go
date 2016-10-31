package config

import (
	"fmt"
	"io/ioutil"
)

import "github.com/hashicorp/hcl"

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
