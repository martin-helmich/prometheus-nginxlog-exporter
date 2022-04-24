package metrics

import (
	"github.com/martin-helmich/prometheus-nginxlog-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

type NamespaceMetrics struct {
	cfg      *config.NamespaceConfig
	registry *prometheus.Registry

	Collection
}

func NewForNamespace(cfg *config.NamespaceConfig) *NamespaceMetrics {
	m := &NamespaceMetrics{
		cfg:      cfg,
		registry: prometheus.NewRegistry(),
	}
	m.Init(cfg)
	m.MustRegister(m.registry)

	return m
}

func (m *NamespaceMetrics) Gatherer() prometheus.Gatherer {
	return m.registry
}
