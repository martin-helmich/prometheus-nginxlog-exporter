package config

// LoadConfigFromFlags fills a configuration object (passed as parameter) with
// values from command-line flags.
func LoadConfigFromFlags(config *Config, flags *StartupFlags) error {
	config.Listen = ListenConfig{
		Port:            flags.ListenPort,
		Address:         flags.ListenAddress,
		MetricsEndpoint: flags.MetricsEndpoint,
	}
	config.Namespaces = []NamespaceConfig{
		{
			Format: flags.Format,
			Parser: flags.Parser,
			Name:   flags.Namespace,
			SourceData: SourceData{
				Files: flags.Filenames,
			},
		},
	}

	return nil
}
