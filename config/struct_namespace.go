package config

import (
	"errors"
	"sort"
)

// NamespaceConfig is a struct describing single metric namespaces
type NamespaceConfig struct {
	Name string `hcl:",key"`

	NamespaceLabelName string `hcl:"namespace_label" yaml:"namespace_label"`
	NamespaceLabels    map[string]string

	MetricsOverride *struct {
		Prefix string `hcl:"prefix" yaml:"prefix"`
	} `hcl:"metrics_override" yaml:"metrics_override"`
	NamespacePrefix string

	SourceFiles      []string          `hcl:"source_files" yaml:"source_files"`
	SourceData       SourceData        `hcl:"source" yaml:"source"`
	Format           string            `hcl:"format"`
	Labels           map[string]string `hcl:"labels"`
	RelabelConfigs   []RelabelConfig   `hcl:"relabel" yaml:"relabel_configs"`
	HistogramBuckets []float64         `hcl:"histogram_buckets" yaml:"histogram_buckets"`

	OrderedLabelNames  []string
	OrderedLabelValues []string
}

type SourceData struct {
	Files  FileSource    `hcl:"files" yaml:"files"`
	Syslog *SyslogSource `hcl:"syslog" yaml:"syslog"`
	Docker *DockerSource `hcl:"docker" yaml:"docker"`
}

type FileSource []string

type SyslogSource struct {
	ListenAddress string   `hcl:"listen_address" yaml:"listen_address"`
	Format        string   `hcl:"format" yaml:"format"`
	Tags          []string `hcl:"tags" yaml:"tags"`
}

type DockerSource struct {
	ContainerName string `hcl:"container" yaml:"container"`
}

// StabilityWarnings tests if the NamespaceConfig uses any configuration settings
// that are not yet declared "stable"
func (c *NamespaceConfig) StabilityWarnings() error {
	if len(c.RelabelConfigs) > 0 {
		return errors.New("you are using the 'relabel' configuration parameter")
	}

	if c.SourceData.Docker != nil {
		return errors.New("you are using the 'docker' log source")
	}

	return nil
}

// DeprecationWarnings tests if the NamespaceConfig uses any deprecated
// configuration settings
func (c *NamespaceConfig) DeprecationWarnings() error {
	if len(c.SourceFiles) > 0 {
		return errors.New("you are using the 'source_files' configuration parameter")
	}

	return nil
}

// MustCompile compiles the configuration (mostly regular expressions that are used
// in configuration variables) for later use
func (c *NamespaceConfig) MustCompile() {
	err := c.Compile()
	if err != nil {
		panic(err)
	}
}

// ResolveDeprecations converts any values from depreated fields into the new
// structures
func (c *NamespaceConfig) ResolveDeprecations() {
	if len(c.SourceFiles) > 0 {
		c.SourceData.Files = FileSource(c.SourceFiles)
	}
}

// Compile compiles the configuration (mostly regular expressions that are used
// in configuration variables) for later use
func (c *NamespaceConfig) Compile() error {
	for i := range c.RelabelConfigs {
		if err := c.RelabelConfigs[i].Compile(); err != nil {
			return nil
		}
	}
	if c.NamespaceLabelName != "" {
		c.NamespaceLabels = make(map[string]string)
		c.NamespaceLabels[c.NamespaceLabelName] = c.Name
	}

	c.OrderLabels()
	c.NamespacePrefix = c.Name
	if c.MetricsOverride != nil {
		c.NamespacePrefix = c.MetricsOverride.Prefix
	}

	return nil
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
