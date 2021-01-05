package relabeling

import "github.com/martin-helmich/prometheus-nginxlog-exporter/config"

// Relabeling contains a relabeling configuration and is responsible for
// executing the rules specified in the original configuration
type Relabeling struct {
	config.RelabelConfig
}

// NewRelabelings creates a new set of relabelling runners from a list of
// configurations (which are typically read from the config file)
func NewRelabelings(cfgs []config.RelabelConfig) []*Relabeling {
	r := make([]*Relabeling, len(cfgs))

	for i := range cfgs {
		r[i] = NewRelabeling(&cfgs[i])
	}

	return r
}

// NewRelabeling creates a single new relabelling runner
func NewRelabeling(cfg *config.RelabelConfig) *Relabeling {
	return &Relabeling{*cfg}
}

// UniqueRelabelings creates a unique relabelings, the duplicated one at the end will discard.
func UniqueRelabelings(relabelings []*Relabeling) []*Relabeling {
	result := make([]*Relabeling, 0, len(relabelings))
	found := make(map[string]struct{})
	for _, r := range relabelings {
		if _, ok := found[r.TargetLabel]; ok {
			continue
		}
		found[r.TargetLabel] = struct{}{}
		result = append(result, r)
	}
	return result
}

// StripOnlyCounterValues strips all values that are associated to relabelings only intended for the request counter
func StripOnlyCounterValues(values []string, relabelings []*Relabeling) []string {
	result := make([]string, 0, len(values))
	offset := len(values) - len(relabelings)
	for i := range values {
		if i >= offset && relabelings[i-offset].OnlyCounter {
			// skip if relabeling and only enabled for counter
			continue
		}
		result = append(result, values[i])
	}
	return result
}
