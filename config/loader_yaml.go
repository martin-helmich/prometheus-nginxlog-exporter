package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
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
