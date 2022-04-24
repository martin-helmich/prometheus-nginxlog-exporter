package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func loadConfigFromYAMLStream(config *Config, file io.Reader) error {
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return err
	}

	return nil
}
