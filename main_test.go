package main

import (
	"fmt"
	"testing"

	"github.com/martin-helmich/prometheus-nginxlog-exporter/log"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/config"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/metrics"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/relabeling"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

type stubFollower struct {
	lines chan string
}

func newStubFollower(lines ...string) *stubFollower {
	ch := make(chan string, len(lines))
	for _, line := range lines {
		ch <- line
	}
	close(ch)

	return &stubFollower{lines: ch}
}

func (f *stubFollower) Lines() chan string {
	return f.lines
}

func (f *stubFollower) OnError(func(error)) {}

type stubParser struct {
	fieldsByLine map[string]map[string]string
}

func (p *stubParser) ParseString(line string) (map[string]string, error) {
	fields, ok := p.fieldsByLine[line]
	if !ok {
		return nil, fmt.Errorf("unexpected line %q", line)
	}

	return fields, nil
}

func collectCounterSamples(t *testing.T, collector *prometheus.CounterVec) map[string]float64 {
	t.Helper()

	metricCh := make(chan prometheus.Metric, 16)
	collector.Collect(metricCh)
	close(metricCh)

	samples := make(map[string]float64)

	for metric := range metricCh {
		dtoMetric := &dto.Metric{}
		require.NoError(t, metric.Write(dtoMetric))

		labels := make(map[string]string, len(dtoMetric.GetLabel()))
		for _, label := range dtoMetric.GetLabel() {
			labels[label.GetName()] = label.GetValue()
		}

		key := fmt.Sprintf(
			"env=%s,page_type=%s,method=%s,status=%s",
			labels["env"],
			labels["page_type"],
			labels["method"],
			labels["status"],
		)
		samples[key] = dtoMetric.GetCounter().GetValue()
	}

	return samples
}

func TestProcessSourceDoesNotReuseRelabelValuesFromPreviousLine(t *testing.T) {
	logger, err := log.New("error", "console")
	require.NoError(t, err)

	nsCfg := &config.NamespaceConfig{
		Name: "nginx",
		Labels: map[string]string{
			"env": "prod",
		},
		RelabelConfigs: []config.RelabelConfig{
			{
				TargetLabel: "page_type",
				SourceValue: "page_type",
			},
		},
	}

	collection := &metrics.Collection{}
	collection.Init(nsCfg)
	relabelings := relabeling.BuildRelabelings(nsCfg.RelabelConfigs)

	parser := &stubParser{
		fieldsByLine: map[string]map[string]string{
			"line-1": {
				"page_type": "article",
				"request":   "GET /article HTTP/1.1",
				"status":    "200",
			},
			"line-2": {
				"request": "GET /landing HTTP/1.1",
				"status":  "200",
			},
		},
	}

	err = processSource(logger, nsCfg, newStubFollower("line-1", "line-2"), parser, collection, relabelings, false)
	require.NoError(t, err)

	samples := collectCounterSamples(t, collection.CountTotal)
	require.Len(t, samples, 2)
	require.Equal(t, 1.0, samples["env=prod,page_type=article,method=GET,status=200"])
	require.Equal(t, 1.0, samples["env=prod,page_type=,method=GET,status=200"])
}
