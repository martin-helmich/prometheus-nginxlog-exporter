package config

// StartupFlags is a struct containing options that can be passed via the
// command line
type StartupFlags struct {
	ConfigFile                 string
	Filenames                  []string
	Format                     string
	Namespace                  string
	ListenPort                 int
	EnableExperimentalFeatures bool

	CPUProfile string
	MemProfile string
}

// Config models the application's configuration
type Config struct {
	Listen                     ListenConfig
	Consul                     ConsulConfig
	Namespaces                 []NamespaceConfig `hcl:"namespace"`
	EnableExperimentalFeatures bool              `hcl:"enable_experimental" yaml:"enable_experimental"`

	// In YAML, the EnableExperimentalFeatures property was originally set by the
	// "enableexperimentalfeatures" property (although documented as "enable_experimental").
	// This property is here for enabling the config to behave as documented, while keeping BC.
	EnableExperimentalFeaturesOld bool `yaml:"enableexperimentalfeatures"`
}

// ListenConfig is a struct describing the built-in webserver configuration
type ListenConfig struct {
	Port    int
	Address string
}

// ConsulConfig describes the connection to a Consul server that the exporter should
// register itself at
type ConsulConfig struct {
	Enable     bool
	Address    string
	Datacenter string
	Scheme     string
	Token      string
	Service    ConsulServiceConfig
}

// ConsulServiceConfig describes the Consul service that the exporter should use
type ConsulServiceConfig struct {
	ID   string
	Name string
	Tags []string
}

// StabilityWarnings tests if the Config or any of its sub-objects uses any
// configuration settings that are not yet declared "stable"
func (c *Config) StabilityWarnings() error {
	if c.EnableExperimentalFeatures || c.EnableExperimentalFeaturesOld {
		return nil
	}

	for i := range c.Namespaces {
		if err := c.Namespaces[i].StabilityWarnings(); err != nil {
			return err
		}
	}

	return nil
}
