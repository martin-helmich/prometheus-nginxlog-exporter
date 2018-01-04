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
}

// Config models the application's configuration
type Config struct {
	Listen                     ListenConfig
	Consul                     ConsulConfig
	Namespaces                 []NamespaceConfig `hcl:"namespace"`
	EnableExperimentalFeatures bool              `hcl:"enable_experimental"`
}

// ListenConfig is a struct describing the built-in webserver configuration
type ListenConfig struct {
	Port    int
	Address string
}

type ConsulConfig struct {
	Enable     bool
	Address    string
	Datacenter string
	Scheme     string
	Token      string
	Service    ConsulServiceConfig
}

type ConsulServiceConfig struct {
	ID   string
	Name string
	Tags []string
}

func (c *Config) StabilityWarnings() error {
	if c.EnableExperimentalFeatures {
		return nil
	}

	for i := range c.Namespaces {
		if err := c.Namespaces[i].StabilityWarnings(); err != nil {
			return err
		}
	}

	return nil
}
