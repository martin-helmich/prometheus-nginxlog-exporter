package config

import (
	"fmt"
	"regexp"
)

// RelabelConfig is a struct describing a single re-labeling configuration for taking
// over label values from an access log line into a Prometheus metric
type RelabelConfig struct {
	TargetLabel string              `hcl:",key" yaml:"target_label"`
	SourceValue string              `hcl:"from" yaml:"from"`
	Whitelist   []string            `hcl:"whitelist"`
	Matches     []RelabelValueMatch `hcl:"match"`
	Split       int                 `hcl:"split"`
	Separator   string              `hcl:"separator"`
	OnlyCounter bool                `hcl:"only_counter" yaml:"only_counter"`

	WhitelistExists bool
	WhitelistMap    map[string]interface{}
}

// RelabelValueMatch describes a single label match statement
type RelabelValueMatch struct {
	RegexpString   string `hcl:",key" yaml:"regexp"`
	Replacement    string `hcl:"replacement"`
	IgnoreString   string `hcl:"ignore" yaml:"ignore"`
	CompiledRegexp *regexp.Regexp
	CompiledIgnore *regexp.Regexp
}

// Compile compiles expressions and lookup tables for efficient later use
func (c *RelabelConfig) Compile() error {
	c.WhitelistMap = make(map[string]interface{})
	c.WhitelistExists = len(c.Whitelist) > 0

	for i := range c.Whitelist {
		c.WhitelistMap[c.Whitelist[i]] = nil
	}

	for i := range c.Matches {
		if c.Matches[i].RegexpString != "" {
			r, err := regexp.Compile(c.Matches[i].RegexpString)
			if err != nil {
				return fmt.Errorf("could not compile regexp '%s': %s", c.Matches[i].RegexpString, err.Error())
			}

			c.Matches[i].CompiledRegexp = r
		}

		if c.Matches[i].IgnoreString != "" {
			r, err := regexp.Compile(c.Matches[i].IgnoreString)
			if err != nil {
				return fmt.Errorf("could not compile regexp '%s': %s", c.Matches[i].IgnoreString, err.Error())
			}

			c.Matches[i].CompiledIgnore = r
		}
	}

	return nil
}
