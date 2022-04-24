package config

import (
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl"
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
