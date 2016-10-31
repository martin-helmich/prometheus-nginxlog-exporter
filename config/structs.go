package config

// StartOptions is a struct containing options that can be passed via the
// command line
type StartupFlags struct {
	ConfigFile string
	Filenames  []string
	Format     string
	Namespace  string
	ListenPort int
}

type Config struct {
	Listen     ListenConfig
	Namespaces []NamespaceConfig `hcl:"namespace"`
}

type ListenConfig struct {
	Port    int
	Address string
}

type NamespaceConfig struct {
	Name        string   `hcl:",key"`
	SourceFiles []string `hcl:"source_files"`
	Format      string   `hcl:"format"`
}
