package relabeling

import "github.com/martin-helmich/prometheus-nginxlog-exporter/config"

type Relabeling struct {
	config.RelabelConfig
}

func NewRelabelings(cfgs []config.RelabelConfig) []*Relabeling {
	r := make([]*Relabeling, len(cfgs))

	for i := range cfgs {
		r[i] = NewRelabeling(&cfgs[i])
	}

	return r
}

func NewRelabeling(cfg *config.RelabelConfig) *Relabeling {
	return &Relabeling{*cfg}
}