package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamespaceCardinalityWarnings(t *testing.T) {
	t.Parallel()

	namespace := NamespaceConfig{
		Name: "nginx",
		RelabelConfigs: []RelabelConfig{
			{
				TargetLabel: "page_type",
				SourceValue: "sent_http_x_page_type",
			},
			{
				TargetLabel: "user",
				SourceValue: "remote_user",
				Whitelist:   []string{"-", "user1"},
			},
			{
				TargetLabel:    "request_uri",
				SourceValue:    "request",
				MaxCardinality: 100,
			},
			{
				TargetLabel: "route",
				SourceValue: "request",
				Matches: []RelabelValueMatch{
					{RegexpString: "^/users/[0-9]+", Replacement: "/users/:id"},
				},
			},
		},
	}

	warnings := namespace.CardinalityWarnings()
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], `namespace "nginx"`)
	assert.Contains(t, warnings[0], `relabel "page_type"`)
	assert.Contains(t, warnings[0], "max_cardinality")
}
