package relabeling

import (
	"testing"

	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/config"
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

func TestUnlimitedCardinalityWhenUnset(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "a", "a")
	assertMapping(t, r, "b", "b")
	assertMapping(t, r, "c", "c")
}

func TestMaxCardinalityMapsExcessValuesToOverflowValue(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{
		MaxCardinality: 2,
		OverflowValue:  "overflow",
	})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "a", "a")
	assertMapping(t, r, "b", "b")
	assertMapping(t, r, "a", "a")
	assertMapping(t, r, "c", "overflow")
	assertMapping(t, r, "d", "overflow")
}

func TestMaxCardinalityDefaultsOverflowValueToOther(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{
		MaxCardinality: 1,
	})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "a", "a")
	assertMapping(t, r, "b", "other")
}

func TestMaxCardinalityTracksMappedValuesAfterRegexReplacement(t *testing.T) {
	t.Parallel()

	r, err := buildRelabeling(config.RelabelConfig{
		Split:          2,
		MaxCardinality: 1,
		Matches: []config.RelabelValueMatch{
			{RegexpString: "^/users/[0-9]+", Replacement: "/users/:id"},
			{RegexpString: "^/orders/[0-9]+", Replacement: "/orders/:id"},
		},
	})
	if err != nil {
		t.Error(err)
	}

	assertMapping(t, r, "GET /users/123 HTTP/1.1", "/users/:id")
	assertMapping(t, r, "GET /users/456 HTTP/1.1", "/users/:id")
	assertMapping(t, r, "GET /orders/999 HTTP/1.1", "other")
}
