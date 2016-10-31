package config

import (
	"fmt"
	"io/ioutil"
)
import "github.com/hashicorp/hcl"

func LoadConfigFromFile(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := Config{}
	hclText := string(buf)

	err = hcl.Decode(&config, hclText)
	if err != nil {
		return nil, err
	}

	fmt.Println(config)
	return &config, nil
}

func LoadConfigFromFlags(flags *StartupFlags) (*Config, error) {
	config := Config{
		Port: flags.ListenPort,
	}

	return &config, nil
}
