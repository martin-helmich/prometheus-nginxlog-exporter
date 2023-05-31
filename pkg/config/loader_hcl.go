package config

import (
	"io"

	"github.com/hashicorp/hcl"
)

func loadConfigFromHCLStream(config *Config, file io.Reader) error {
	buf, err := io.ReadAll(file)
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
