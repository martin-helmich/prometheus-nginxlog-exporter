package config

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
	Namespaces []NamespaceConfig `hcl:"namespace"`
}

// ListenConfig is a struct describing the built-in webserver configuration
type ListenConfig struct {
	Port    int
	Address string
}

// NamespaceConfig is a struct describing single metric namespaces
type NamespaceConfig struct {
	Name        string            `hcl:",key"`
	SourceFiles []string          `hcl:"source_files"`
	Format      string            `hcl:"format"`
	Labels      map[string]string `hcl:"labels"`
}
