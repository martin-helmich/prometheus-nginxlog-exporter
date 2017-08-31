package config

import (
	"regexp"
	"sort"
)

// StartupFlags is a struct containing options that can be passed via the
// command line
type StartupFlags struct {
	ConfigFile string
	Filenames  []string
	Format     string
	Namespace  string
	ListenPort int
}

// Config models the application's configuration
type Config struct {
	Listen     ListenConfig
	Consul     ConsulConfig
	Namespaces []NamespaceConfig `hcl:"namespace"`
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

// NamespaceConfig is a struct describing single metric namespaces
type NamespaceConfig struct {
	Name        string            `hcl:",key"`
	SourceFiles []string          `hcl:"source_files"`
	Format      string            `hcl:"format"`
	Labels      map[string]string `hcl:"labels"`
	Routes      []string          `hcl:"routes"`

	OrderedLabelNames  []string
	OrderedLabelValues []string

	CompiledRouteRegexps []*regexp.Regexp
}

func (c *NamespaceConfig) MustCompileRouteRegexps() {
	err := c.CompileRouteRegexps()
	if err != nil {
		panic(err)
	}
}

func (c *NamespaceConfig) CompileRouteRegexps() error {
	c.CompiledRouteRegexps = make([]*regexp.Regexp, len(c.Routes))

	for i, e := range c.Routes {
		r, err := regexp.Compile(e)
		if err != nil {
			return err
		}

		c.CompiledRouteRegexps[i] = r
	}

	return nil
}

func (c *NamespaceConfig) HasRouteMatchers() bool {
	return len(c.Routes) > 0
}

// OrderLabels builds two lists of label keys and values, ordered by label name
func (c *NamespaceConfig) OrderLabels() {
	keys := make([]string, 0, len(c.Labels))
	values := make([]string, len(c.Labels))

	for k := range c.Labels {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i, k := range keys {
		values[i] = c.Labels[k]
	}

	c.OrderedLabelNames = keys
	c.OrderedLabelValues = values
}
