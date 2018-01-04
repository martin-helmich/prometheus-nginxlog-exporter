package relabeling

import "github.com/martin-helmich/prometheus-nginxlog-exporter/config"

type Relabeling struct {
	config.RelabelConfig
}

type RelabelingSet []*Relabeling

func NewRelabelings(cfgs []config.RelabelConfig) RelabelingSet {
	r := make([]*Relabeling, len(cfgs))

	for i := range cfgs {
		r[i] = NewRelabeling(&cfgs[i])
	}

	return r
}

func NewRelabeling(cfg *config.RelabelConfig) *Relabeling {
	return &Relabeling{*cfg}
}