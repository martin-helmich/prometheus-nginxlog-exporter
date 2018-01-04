package config

// RelabelConfig is a struct describing a single re-labeling configuration for taking
// over label values from an access log line into a Prometheus metric
type RelabelConfig struct {
	TargetLabel string   `hcl:",key"`
	SourceValue string   `hcl:"from"`
	Whitelist   []string `hcl:"whitelist"`

	WhitelistExists bool
	WhitelistMap    map[string]interface{}
}

func (c *RelabelConfig) Compile() error {
	c.WhitelistMap = make(map[string]interface{})
	c.WhitelistExists = len(c.Whitelist) > 0

	for i := range c.Whitelist {
		c.WhitelistMap[c.Whitelist[i]] = nil
	}

	return nil
}
