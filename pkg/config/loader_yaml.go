package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

func loadConfigFromYAMLStream(config *Config, file io.Reader) error {
	buf, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return err
	}

	return nil
}
