/*
 * Copyright 2019 Martin Helmich <martin@helmich.me>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/martin-helmich/prometheus-nginxlog-exporter/syslog"

	"github.com/martin-helmich/prometheus-nginxlog-exporter/config"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/discovery"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/parser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/prof"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/relabeling"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/tail"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type NSMetrics struct {
	cfg      *config.NamespaceConfig
	registry *prometheus.Registry
	Metrics
}

func NewNSMetrics(cfg *config.NamespaceConfig) *NSMetrics {
	m := &NSMetrics{
		cfg:      cfg,
		registry: prometheus.NewRegistry(),
	}
	m.Init(cfg)

	m.registry.MustRegister(m.countTotal)
	m.registry.MustRegister(m.requestBytesTotal)
	m.registry.MustRegister(m.responseBytesTotal)
	m.registry.MustRegister(m.upstreamSeconds)
	m.registry.MustRegister(m.upstreamSecondsHist)
	m.registry.MustRegister(m.responseSeconds)
	m.registry.MustRegister(m.responseSecondsHist)
	m.registry.MustRegister(m.parseErrorsTotal)
	return m
}

// Metrics is a struct containing pointers to all metrics that should be
// exposed to Prometheus
type Metrics struct {
	countTotal          *prometheus.CounterVec
	responseBytesTotal  *prometheus.CounterVec
	requestBytesTotal   *prometheus.CounterVec
	upstreamSeconds     *prometheus.SummaryVec
	upstreamSecondsHist *prometheus.HistogramVec
	responseSeconds     *prometheus.SummaryVec
	responseSecondsHist *prometheus.HistogramVec
	parseErrorsTotal    prometheus.Counter
}

func inLabels(label string, labels []string) bool {
	for _, l := range labels {
		if label == l {
			return true
		}
	}
	return false
}

// Init initializes a metrics struct
func (m *Metrics) Init(cfg *config.NamespaceConfig) {
	cfg.MustCompile()

	labels := cfg.OrderedLabelNames
	counterLabels := labels

	for i := range cfg.RelabelConfigs {
		if !cfg.RelabelConfigs[i].OnlyCounter {
			labels = append(labels, cfg.RelabelConfigs[i].TargetLabel)
		}
		counterLabels = append(counterLabels, cfg.RelabelConfigs[i].TargetLabel)
	}

	for _, r := range relabeling.DefaultRelabelings {
		if !inLabels(r.TargetLabel, labels) {
			labels = append(labels, r.TargetLabel)
		}
		if !inLabels(r.TargetLabel, counterLabels) {
			counterLabels = append(counterLabels, r.TargetLabel)
		}
	}

	m.countTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_count_total",
		Help:        "Amount of processed HTTP requests",
	}, counterLabels)

	m.responseBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_size_bytes",
		Help:        "Total amount of transferred bytes",
	}, labels)

	m.requestBytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_request_size_bytes",
		Help:        "Total amount of received bytes",
	}, labels)

	m.upstreamSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_upstream_time_seconds",
		Help:        "Time needed by upstream servers to handle requests",
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)

	m.upstreamSecondsHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_upstream_time_seconds_hist",
		Help:        "Time needed by upstream servers to handle requests",
		Buckets:     cfg.HistogramBuckets,
	}, labels)

	m.responseSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_time_seconds",
		Help:        "Time needed by NGINX to handle requests",
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)

	m.responseSecondsHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "http_response_time_seconds_hist",
		Help:        "Time needed by NGINX to handle requests",
		Buckets:     cfg.HistogramBuckets,
	}, labels)

	m.parseErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   cfg.NamespacePrefix,
		ConstLabels: cfg.NamespaceLabels,
		Name:        "parse_errors_total",
		Help:        "Total number of log file lines that could not be parsed",
	})
}

func main() {
	var opts config.StartupFlags
	var cfg = config.Config{
		Listen: config.ListenConfig{
			Port:            4040,
			Address:         "0.0.0.0",
			MetricsEndpoint: "/metrics",
		},
	}
	nsGatherers := make(prometheus.Gatherers, 0)

	flag.IntVar(&opts.ListenPort, "listen-port", 4040, "HTTP port to listen on")
	flag.StringVar(&opts.ListenAddress, "listen-address", "0.0.0.0", "IP-address to bind")
	flag.StringVar(&opts.Parser, "parser", "text", "NGINX access log format parser. One of: [text, json]")
	flag.StringVar(&opts.Format, "format", `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"`, "NGINX access log format")
	flag.StringVar(&opts.Namespace, "namespace", "nginx", "namespace to use for metric names")
	flag.StringVar(&opts.ConfigFile, "config-file", "", "Configuration file to read from")
	flag.BoolVar(&opts.EnableExperimentalFeatures, "enable-experimental", false, "Set this flag to enable experimental features")
	flag.StringVar(&opts.CPUProfile, "cpuprofile", "", "write cpu profile to `file`")
	flag.StringVar(&opts.MemProfile, "memprofile", "", "write memory profile to `file`")
	flag.StringVar(&opts.MetricsEndpoint, "metrics-endpoint", cfg.Listen.MetricsEndpoint, "URL path at which to serve metrics")
	flag.BoolVar(&opts.VerifyConfig, "verify-config", false, "Enable this flag to check config file loads, then exit")
	flag.Parse()

	opts.Filenames = flag.Args()

	sigChan := make(chan os.Signal, 1)
	stopChan := make(chan bool)
	stopHandlers := sync.WaitGroup{}

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	go func() {
		sig := <-sigChan

		fmt.Printf("caught term %s. exiting\n", sig)

		close(stopChan)
		stopHandlers.Wait()

		os.Exit(0)
	}()

	defer func() {
		close(stopChan)
		stopHandlers.Wait()
	}()

	prof.SetupCPUProfiling(opts.CPUProfile, stopChan, &stopHandlers)
	prof.SetupMemoryProfiling(opts.MemProfile, stopChan, &stopHandlers)

	loadConfig(&opts, &cfg)

	fmt.Printf("using configuration %+v\n", cfg)

	if stabilityError := cfg.StabilityWarnings(); stabilityError != nil && !opts.EnableExperimentalFeatures {
		fmt.Fprintf(os.Stderr, "Your configuration file contains an option that is explicitly labeled as experimental feature:\n\n  %s\n\n", stabilityError.Error())
		fmt.Fprintln(os.Stderr, "Use the -enable-experimental flag or the enable_experimental option to enable these features. Use them at your own peril.")

		os.Exit(1)
	}

	if cfg.Consul.Enable {
		setupConsul(&cfg, stopChan, &stopHandlers)
	}

	for _, ns := range cfg.Namespaces {
		nsMetrics := NewNSMetrics(&ns)
		nsGatherers = append(nsGatherers, nsMetrics.registry)

		fmt.Printf("starting listener for namespace %s\n", ns.Name)
		go processNamespace(ns, &(nsMetrics.Metrics))
	}

	listenAddr := fmt.Sprintf("%s:%d", cfg.Listen.Address, cfg.Listen.Port)
	endpoint := cfg.Listen.MetricsEndpointOrDefault()

	fmt.Printf("running HTTP server on address %s, serving metrics at %s\n", listenAddr, endpoint)

	nsHandler := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(nsGatherers, promhttp.HandlerOpts{}),
	)

	http.Handle(endpoint, nsHandler)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		fmt.Printf("error while starting HTTP server: %s", err.Error())
	}
}

func loadConfig(opts *config.StartupFlags, cfg *config.Config) {
	if opts.ConfigFile != "" {
		fmt.Printf("loading configuration file %s\n", opts.ConfigFile)
		if err := config.LoadConfigFromFile(cfg, opts.ConfigFile); err != nil {
			panic(err)
		}
	} else if err := config.LoadConfigFromFlags(cfg, opts); err != nil {
		panic(err)
	}
	if opts.VerifyConfig {
		fmt.Printf("Configuration is valid")
		os.Exit(0)
	}
}

func setupConsul(cfg *config.Config, stopChan <-chan bool, stopHandlers *sync.WaitGroup) {
	registrator, err := discovery.NewConsulRegistrator(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("registering service in Consul\n")
	if err := registrator.RegisterConsul(); err != nil {
		panic(err)
	}

	go func() {
		<-stopChan
		fmt.Printf("unregistering service in Consul\n")

		if err := registrator.UnregisterConsul(); err != nil {
			fmt.Printf("error while unregistering from consul: %s\n", err.Error())
		}

		stopHandlers.Done()
	}()

	stopHandlers.Add(1)
}

func processNamespace(nsCfg config.NamespaceConfig, metrics *Metrics) {
	var followers []tail.Follower

	parser := parser.NewParser(nsCfg)

	for _, f := range nsCfg.SourceData.Files {
		t, err := tail.NewFileFollower(f)
		if err != nil {
			panic(err)
		}

		t.OnError(func(err error) {
			panic(err)
		})

		followers = append(followers, t)
	}

	if nsCfg.SourceData.Syslog != nil {
		slCfg := nsCfg.SourceData.Syslog

		fmt.Printf("running Syslog server on address %s\n", slCfg.ListenAddress)
		channel, server, err := syslog.Listen(slCfg.ListenAddress, slCfg.Format)
		if err != nil {
			panic(err)
		}

		for _, f := range slCfg.Tags {
			t, err := tail.NewSyslogFollower(f, server, channel)
			if err != nil {
				panic(err)
			}

			t.OnError(func(err error) {
				panic(err)
			})

			followers = append(followers, t)
		}
	}

	// determine once if there are any relabeling configurations for only the response counter
	hasCounterOnlyLabels := false
	for _, r := range nsCfg.RelabelConfigs {
		if r.OnlyCounter {
			hasCounterOnlyLabels = true
			break
		}
	}

	for _, f := range followers {
		go processSource(nsCfg, f, parser, metrics, hasCounterOnlyLabels)
	}

}

func processSource(nsCfg config.NamespaceConfig, t tail.Follower, parser parser.Parser, metrics *Metrics, hasCounterOnlyLabels bool) {
	relabelings := relabeling.NewRelabelings(nsCfg.RelabelConfigs)
	relabelings = append(relabelings, relabeling.DefaultRelabelings...)
	relabelings = relabeling.UniqueRelabelings(relabelings)

	staticLabelValues := nsCfg.OrderedLabelValues

	totalLabelCount := len(staticLabelValues) + len(relabelings)
	relabelLabelOffset := len(staticLabelValues)
	labelValues := make([]string, totalLabelCount)

	for i := range staticLabelValues {
		labelValues[i] = staticLabelValues[i]
	}

	for line := range t.Lines() {
		if nsCfg.PrintLog {
			fmt.Println(line)
		}

		fields, err := parser.ParseString(line)
		if err != nil {
			fmt.Printf("error while parsing line '%s': %s\n", line, err)
			metrics.parseErrorsTotal.Inc()
			continue
		}

		for i := range relabelings {
			if str, ok := fields[relabelings[i].SourceValue]; ok {
				mapped, err := relabelings[i].Map(str)
				if err == nil {
					labelValues[i+relabelLabelOffset] = mapped
				}
			}
		}

		var notCounterValues []string
		if hasCounterOnlyLabels {
			notCounterValues = relabeling.StripOnlyCounterValues(labelValues, relabelings)
		} else {
			notCounterValues = labelValues
		}

		metrics.countTotal.WithLabelValues(labelValues...).Inc()

		if bytes, ok, err := floatFromFields(fields, "body_bytes_sent"); ok {
			metrics.responseBytesTotal.WithLabelValues(notCounterValues...).Add(bytes)
		} else if err != nil {
			fmt.Printf("error while parsing $body_bytes_sent: %v\n", err)
			metrics.parseErrorsTotal.Inc()
		}

		if bytes, ok, err := floatFromFields(fields, "request_length"); ok {
			metrics.requestBytesTotal.WithLabelValues(notCounterValues...).Add(bytes)
		} else if err != nil {
			fmt.Printf("error while parsing $request_length: %v\n", err)
			metrics.parseErrorsTotal.Inc()
		}

		if upstreamTime, ok, err := floatFromFieldsMulti(fields, "upstream_response_time"); ok {
			metrics.upstreamSeconds.WithLabelValues(notCounterValues...).Observe(upstreamTime)
			metrics.upstreamSecondsHist.WithLabelValues(notCounterValues...).Observe(upstreamTime)
		} else if err != nil {
			fmt.Printf("error while parsing $upstream_response_time: %v\n", err)
			metrics.parseErrorsTotal.Inc()
		}

		if responseTime, ok, err := floatFromFields(fields, "request_time"); ok {
			metrics.responseSeconds.WithLabelValues(notCounterValues...).Observe(responseTime)
			metrics.responseSecondsHist.WithLabelValues(notCounterValues...).Observe(responseTime)
		} else if err != nil {
			fmt.Printf("error while parsing $request_time: %v\n", err)
			metrics.parseErrorsTotal.Inc()
		}
	}
}

func floatFromFieldsMulti(fields map[string]string, name string) (float64, bool, error) {
	f, ok, err := floatFromFields(fields, name)
	if err == nil {
		return f, ok, nil
	}

	val, ok := fields[name]
	if !ok {
		return 0, false, nil
	}

	sum := float64(0)

	for _, v := range strings.Split(val, ",") {
		v = strings.TrimSpace(v)

		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false, fmt.Errorf("value '%s' could not be parsed into float", val)
		}

		sum += f
	}

	return sum, true, nil
}

func floatFromFields(fields map[string]string, name string) (float64, bool, error) {
	val, ok := fields[name]
	if !ok {
		return 0, false, nil
	}

	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, false, fmt.Errorf("value '%s' could not be parsed into float", val)
	}

	return f, true, nil
}
