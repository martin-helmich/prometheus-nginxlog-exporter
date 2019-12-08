package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const HCLInput = `
listen {
  address = "10.0.0.1"
  port = 4040
}

consul {
  enable = true
  address = "localhost:8500"
  datacenter = "dc1"
  scheme = "https"
  token = "asdfasfdasf"

  service {
    id = "nginx-exporter"
    name = "nginx-exporter"
    tags = ["foo", "bar"]
  }
}

namespace "nginx" {
  source_files = [
    "test.log",
    "foo.log"
  ]
  format = "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""

  labels {
    app = "magicapp"
    foo = "bar"
  }

  relabel "user" {
    from = "remote_user"
    whitelist = ["-", "user1", "user2"]
  }

  relabel "request_uri" {
    from = "request"
    split = 2

    match "^/users/[0-9]+" {
      replacement = "/users/:id"
    }
    match "^/profile" {
      replacement = "/profile"
    }
  }
}
`

const YAMLInput = `
listen:
  address: "10.0.0.1"
  port: 4040

consul:
  enable: true
  address: "localhost:8500"
  datacenter: "dc1"
  scheme: "https"
  token: "asdfasfdasf"

  service:
    id: "nginx-exporter"
    name: "nginx-exporter"
    tags:
      - foo
      - bar

namespaces:
  - name: nginx
    source_files:
      - test.log
      - foo.log
    format: "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\""
    labels:
      app: "magicapp"
      foo: "bar"
    relabel_configs:
      - target_label: user
        from: "remote_user"
        whitelist: ["-", "user1", "user2"]
      - target_label: request_uri
        from: request
        split: 2
        matches:
          - regexp: "^/users/[0-9]+"
            replacement: "/users/:id"
          - regexp: "^/profile"
            replacement: "/profile"
`

func assertConfigContents(t *testing.T, cfg Config) {
	assert.Equal(t, "10.0.0.1", cfg.Listen.Address)
	assert.Equal(t, 4040, cfg.Listen.Port)

	assert.True(t, cfg.Consul.Enable)
	assert.Equal(t, "localhost:8500", cfg.Consul.Address)
	assert.Equal(t, "nginx-exporter", cfg.Consul.Service.ID)
	assert.Equal(t, "nginx-exporter", cfg.Consul.Service.Name)
	assert.Equal(t, []string{"foo", "bar"}, cfg.Consul.Service.Tags)
	assert.Equal(t, "dc1", cfg.Consul.Datacenter)
	assert.Equal(t, "https", cfg.Consul.Scheme)
	assert.Equal(t, "asdfasfdasf", cfg.Consul.Token)

	require.Len(t, cfg.Namespaces, 1)

	n := cfg.Namespaces[0]
	assert.Equal(t, "nginx", n.Name)
	assert.Equal(t, "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$http_x_forwarded_for\"", n.Format)
	assert.Equal(t, []string{"test.log", "foo.log"}, n.SourceFiles)
	assert.Equal(t, FileSource{"test.log", "foo.log"}, n.SourceData.Files)
	assert.Equal(t, "magicapp", n.Labels["app"])

	require.Len(t, n.RelabelConfigs, 2)
	assert.Equal(t, "user", n.RelabelConfigs[0].TargetLabel)
	assert.Equal(t, "request_uri", n.RelabelConfigs[1].TargetLabel)

	assert.Len(t, n.RelabelConfigs[1].Matches, 2)
	assert.Equal(t, "^/users/[0-9]+", n.RelabelConfigs[1].Matches[0].RegexpString)
}

func TestLoadsHCLConfigFile(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBufferString(HCLInput)
	cfg := Config{}

	err := LoadConfigFromStream(&cfg, buf, TypeHCL)
	assert.Nil(t, err, "unexpected error: %v", err)
	assertConfigContents(t, cfg)
}

func TestLoadsYAMLConfigFile(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBufferString(YAMLInput)
	cfg := Config{}

	err := LoadConfigFromStream(&cfg, buf, TypeYAML)
	assert.Nil(t, err, "unexpected error: %v", err)
	assertConfigContents(t, cfg)
}
