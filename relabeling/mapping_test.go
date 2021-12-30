package relabeling

import (
	"github.com/martin-helmich/prometheus-nginxlog-exporter/config"
	"testing"
)

func buildRelabeling(cfg config.RelabelConfig) (*Relabeling, error) {
	if err := cfg.Compile(); err != nil {
		return nil, err
	}

	return NewRelabeling(&cfg), nil
}

func assertMapping(t *testing.T, r *Relabeling, in string, expected string) {
	mapped, err := r.Map(in)
	if err != nil {
		t.Error(err)
	}

	if mapped != expected {
		t.Errorf("expected '%s', but got '%s'", expected, mapped)
	}
}

func TestSplitMapping(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{Split: 2})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "foo bar", "bar")
}

func TestRequestURIMapping(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{
		Split: 2,
		Matches: []config.RelabelValueMatch{
			{RegexpString: "^/users/[0-9]+", Replacement: "/users/:id"},
		},
	})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "GET /users/12345 HTTP/1.1", "/users/:id")
	assertMapping(t, r, "GET /users/12345/about HTTP/1.1", "/users/:id/about")
	assertMapping(t, r, "GET /v1/users/12345 HTTP/1.1", "")
}

func TestIgnore(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{
		Split: 2,
		Matches: []config.RelabelValueMatch{
			{
				RegexpString: ".*",
				Replacement:  "unknown",
				IgnoreString: "^/api",
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "GET /index.html HTTP/1.1", "unknown")
	assertMapping(t, r, "GET /api/v1/index.html HTTP/1.1", "/api/v1/index.html")
}
