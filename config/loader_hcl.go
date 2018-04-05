package config

import (
	"github.com/hashicorp/hcl"
	"io"
	"io/ioutil"
)

func loadConfigFromHCLStream(config *Config, file io.Reader) error {
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	hclText := string(buf)

	err = hcl.Decode(config, hclText)
	if err != nil {
		return err
	}

	return nil
}
