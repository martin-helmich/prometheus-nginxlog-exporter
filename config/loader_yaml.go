package config

import (
	"io/ioutil"
	"io"
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
