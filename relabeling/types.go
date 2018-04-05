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
