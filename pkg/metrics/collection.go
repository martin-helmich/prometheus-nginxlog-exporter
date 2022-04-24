package metrics

import "github.com/prometheus/client_golang/prometheus"

// Collection is a struct containing pointers to all metrics that should be
// exposed to Prometheus
type Collection struct {
	CountTotal                 *prometheus.CounterVec
	ResponseBytesTotal         *prometheus.CounterVec
	RequestBytesTotal          *prometheus.CounterVec
	UpstreamSeconds            *prometheus.SummaryVec
	UpstreamSecondsHist        *prometheus.HistogramVec
	UpstreamConnectSeconds     *prometheus.SummaryVec
	UpstreamConnectSecondsHist *prometheus.HistogramVec
	ResponseSeconds            *prometheus.SummaryVec
	ResponseSecondsHist        *prometheus.HistogramVec
	ParseErrorsTotal           prometheus.Counter
}
