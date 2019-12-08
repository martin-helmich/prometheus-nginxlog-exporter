package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSourceFilesAreMappedToNewSourceConfig(t *testing.T) {
	c := &NamespaceConfig{
		Name:        "foo",
		SourceFiles: []string{"bar.log", "baz.log"},
	}

	c.ResolveDeprecations()

	require.Equal(t, FileSource{"bar.log", "baz.log"}, c.SourceData.Files)
}
