package metrics

import "github.com/prometheus/client_golang/prometheus"

func (c *Collection) MustRegister(r *prometheus.Registry) {
	r.MustRegister(c.CountTotal)
	r.MustRegister(c.RequestBytesTotal)
	r.MustRegister(c.ResponseBytesTotal)
	r.MustRegister(c.UpstreamSeconds)
	r.MustRegister(c.UpstreamSecondsHist)
	r.MustRegister(c.UpstreamConnectSeconds)
	r.MustRegister(c.UpstreamConnectSecondsHist)
	r.MustRegister(c.ResponseSeconds)
	r.MustRegister(c.ResponseSecondsHist)
	r.MustRegister(c.ParseErrorsTotal)
}
