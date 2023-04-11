package metrics

import (
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/config"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/relabeling"
	"github.com/prometheus/client_golang/prometheus"
)

// Init initializes a metrics struct
func (m *Collection) Init(cfg *config.NamespaceConfig) {
	cfg.MustCompile()

	labels := cfg.OrderedLabelNames
	counterLabels := labels

	relabelings := relabeling.NewRelabelings(cfg.RelabelConfigs)
	relabelings = append(relabelings, relabeling.DefaultRelabelings...)
	relabelings = relabeling.UniqueRelabelings(relabelings)

	for _, r := range relabelings {
		if !r.OnlyCounter {
			labels = append(labels, r.TargetLabel)
		}
		counterLabels = append(counterLabels, r.TargetLabel)
	}

	m.CountTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_count_total",
		Help:        "Amount of processed HTTP requests",
	}, counterLabels)

	m.ResponseBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_size_bytes",
		Help:        "Total amount of transferred bytes",
	}, counterLabels)

	m.RequestBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_request_size_bytes",
		Help:        "Total amount of received bytes",
	}, counterLabels)

	m.UpstreamSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_upstream_time_seconds",
		Help:        "Time needed by upstream servers to handle requests",
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, counterLabels)

	m.UpstreamSecondsHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_upstream_time_seconds_hist",
		Help:        "Time needed by upstream servers to handle requests",
		Buckets:     cfg.HistogramBuckets,
	}, labels)

	m.UpstreamConnectSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_upstream_connect_time_seconds",
		Help:        "Time needed to connect to upstream servers",
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)

	m.UpstreamConnectSecondsHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_upstream_connect_time_seconds_hist",
		Help:        "Time needed to connect to upstream servers",
		Buckets:     cfg.HistogramBuckets,
	}, labels)

	m.ResponseSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_time_seconds",
		Help:        "Time needed by NGINX to handle requests",
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)

	m.ResponseSecondsHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_time_seconds_hist",
		Help:        "Time needed by NGINX to handle requests",
		Buckets:     cfg.HistogramBuckets,
	}, labels)

	m.ParseErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "parse_errors_total",
		Help:        "Total number of log file lines that could not be parsed",
	})
}
