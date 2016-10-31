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
	Port       int
	Namespaces []NamespaceConfig `hcl:"namespace"`
}

type NamespaceConfig struct {
	Name        string   `hcl:",key"`
	SourceFiles []string `hcl:"source_files"`
	Format      string   `hcl:"format"`
}
