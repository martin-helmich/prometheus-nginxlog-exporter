/*
 * Copyright 2019-2022 Martin Helmich <martin@helmich.me>
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

	"github.com/martin-helmich/prometheus-nginxlog-exporter/log"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/config"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/discovery"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/metrics"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/parser"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/prof"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/relabeling"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/syslog"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/pkg/tail"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const maxStaticLabels = 128

func main() {
	var opts config.StartupFlags
	var cfg = config.Config{
		Listen: config.ListenConfig{
			Port:            4040,
			Address:         "0.0.0.0",
			MetricsEndpoint: "/metrics",
		},
	}

	versionMetrics := prometheus.NewRegistry()
	versionMetrics.MustRegister(version.NewCollector("prometheus_nginxlog_exporter"))

	gatherers := prometheus.Gatherers{versionMetrics}

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
	flag.StringVar(&opts.LogLevel, "log-level", "info", "level of logs. Allowed values: error, warning, info, debug")
	flag.StringVar(&opts.LogFormat, "log-format", "console", "Define log format. Allowed values: console, json")
	flag.BoolVar(&opts.VerifyConfig, "verify-config", false, "Enable this flag to check config file loads, then exit")
	flag.BoolVar(&opts.Version, "version", false, "set to print version information")
	flag.Parse()

	logger, err := log.New(opts.LogLevel, opts.LogFormat)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts.Filenames = flag.Args()

	sigChan := make(chan os.Signal, 1)
	stopChan := make(chan bool)
	stopHandlers := sync.WaitGroup{}

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	go func() {
		sig := <-sigChan

		logger.Infof("caught term %s. exiting", sig)

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

	loadConfig(logger, &opts, &cfg)

	logger.Debugf("using configuration %+v", cfg)

	if stabilityError := cfg.StabilityWarnings(); stabilityError != nil && !opts.EnableExperimentalFeatures {
		logger.Error("Your configuration file contains an option that is explicitly labeled as experimental feature")
		logger.Error(stabilityError.Error())
		logger.Error("Use the -enable-experimental flag or the enable_experimental option to enable these features. Use them at your own peril.")

		os.Exit(1)
	}

	if cfg.Consul.Enable {
		setupConsul(logger, &cfg, stopChan, &stopHandlers)
	}

	for i := range cfg.Namespaces {
		namespace := &cfg.Namespaces[i]

		nsMetrics := metrics.NewForNamespace(namespace)
		gatherers = append(gatherers, nsMetrics.Gatherer())

		logger.Infof("starting listener for namespace %s", namespace.Name)
		go func(ns *config.NamespaceConfig) {
			processNamespace(logger, ns, &(nsMetrics.Collection), stopChan, &stopHandlers)
		}(namespace)
	}

	listenAddr := fmt.Sprintf("%s:%d", cfg.Listen.Address, cfg.Listen.Port)
	endpoint := cfg.Listen.MetricsEndpointOrDefault()

	logger.Infof("running HTTP server on address %s, serving metrics at %s", listenAddr, endpoint)

	nsHandler := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{}),
	)

	http.Handle(endpoint, nsHandler)

	logger.Fatal(http.ListenAndServe(listenAddr, nil))
}

func loadConfig(logger *log.Logger, opts *config.StartupFlags, cfg *config.Config) {
	if opts.ConfigFile != "" {
		logger.Infof("loading configuration file %s", opts.ConfigFile)
		if err := config.LoadConfigFromFile(logger, cfg, opts.ConfigFile); err != nil {
			logger.Fatal(err)
		}
	} else if err := config.LoadConfigFromFlags(cfg, opts); err != nil {
		logger.Fatal(err)
	}

	if opts.VerifyConfig {
		fmt.Printf("Configuration is valid")
		os.Exit(0)
	}
}

func setupConsul(logger *log.Logger, cfg *config.Config, stopChan <-chan bool, stopHandlers *sync.WaitGroup) {
	registrator, err := discovery.NewConsulRegistrator(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("registering service in Consul")
	if err = registrator.RegisterConsul(); err != nil {
		logger.Fatal(err)
	}

	go func() {
		<-stopChan
		logger.Info("unregistering service in Consul")

		if err := registrator.UnregisterConsul(); err != nil {
			logger.Errorf("error while unregistering from consul: %s", err.Error())
		}

		stopHandlers.Done()
	}()

	stopHandlers.Add(1)
}

func processNamespace(logger *log.Logger, nsCfg *config.NamespaceConfig, metrics *metrics.Collection, stopChan <-chan bool, stopHandlers *sync.WaitGroup) error {
	var followers []tail.Follower

	logParser := parser.NewParser(nsCfg)

	for _, f := range nsCfg.SourceData.Files {
		t, err := tail.NewFileFollower(logger, f)
		if err != nil {
			logger.Fatal(err)
		}

		t.OnError(func(err error) {
			logger.Fatal(err)
		})

		followers = append(followers, t)
	}

	if nsCfg.SourceData.Syslog != nil {
		slCfg := nsCfg.SourceData.Syslog

		logger.Infof("running Syslog server on address %s", slCfg.ListenAddress)
		channel, server, closeServer, err := syslog.Listen(slCfg.ListenAddress, slCfg.Format)
		if err != nil {
			panic(err)
		}

		stopHandlers.Add(1)

		go func() {
			<-stopChan

			if err := closeServer(); err != nil {
				fmt.Printf("error while closing syslog server: %s\n", err.Error())
			}

			stopHandlers.Done()
		}()

		for _, f := range slCfg.Tags {
			t, err := tail.NewSyslogFollower(f, server, channel)
			if err != nil {
				logger.Fatal(err)
			}

			t.OnError(func(err error) {
				logger.Fatal(err)
			})

			followers = append(followers, t)
		}
	}

	if nsCfg.SourceData.Docker != nil {
		t, err := tail.NewDockerFollower(nsCfg.SourceData.Docker.ContainerName)
		if err != nil {
			panic(err)
		}

		followers = append(followers, t)
	}

	// determine once if there are any relabeling configurations for only the response counter
	hasCounterOnlyLabels := false
	for _, r := range nsCfg.RelabelConfigs {
		if r.OnlyCounter {
			hasCounterOnlyLabels = true
			break
		}
	}

	errs := make(chan error)
	defer close(errs)

	for _, follower := range followers {
		go func(f tail.Follower) {
			if err := processSource(logger, nsCfg, f, logParser, metrics, hasCounterOnlyLabels); err != nil {
				errs <- err
			}
		}(follower)
	}

	return <-errs
}

func processSource(logger *log.Logger, nsCfg *config.NamespaceConfig, t tail.Follower, parser parser.Parser, metrics *metrics.Collection, hasCounterOnlyLabels bool) error {
	relabelings := relabeling.NewRelabelings(nsCfg.RelabelConfigs)
	relabelings = append(relabelings, relabeling.DefaultRelabelings...)
	relabelings = relabeling.UniqueRelabelings(relabelings)

	staticLabelValues := nsCfg.OrderedLabelValues

	totalLabelCount := len(staticLabelValues) + len(relabelings)
	relabelLabelOffset := len(staticLabelValues)

	if totalLabelCount > maxStaticLabels {
		return errors.Errorf("configured label count exceeds the maximum count of %d", maxStaticLabels)
	}

	labelValues := make([]string, totalLabelCount)

	copy(labelValues, staticLabelValues)

	for line := range t.Lines() {
		if nsCfg.PrintLog {
			fmt.Println(line)
		}

		fields, err := parser.ParseString(line)
		if err != nil {
			logger.Errorf("error while parsing line '%s': %s", line, err)
			metrics.ParseErrorsTotal.Inc()
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

		metrics.CountTotal.WithLabelValues(labelValues...).Inc()

		if v, ok := observeMetrics(logger, fields, "body_bytes_sent", floatFromFields, metrics.ParseErrorsTotal); ok {
			metrics.ResponseBytesTotal.WithLabelValues(notCounterValues...).Add(v)
		}

		if v, ok := observeMetrics(logger, fields, "request_length", floatFromFields, metrics.ParseErrorsTotal); ok {
			metrics.RequestBytesTotal.WithLabelValues(notCounterValues...).Add(v)
		}

		if v, ok := observeMetrics(logger, fields, "upstream_response_time", floatFromFieldsMulti, metrics.ParseErrorsTotal); ok {
			metrics.UpstreamSeconds.WithLabelValues(notCounterValues...).Observe(v)
			metrics.UpstreamSecondsHist.WithLabelValues(notCounterValues...).Observe(v)
		}

		if v, ok := observeMetrics(logger, fields, "upstream_connect_time", floatFromFieldsMulti, metrics.ParseErrorsTotal); ok {
			metrics.UpstreamConnectSeconds.WithLabelValues(notCounterValues...).Observe(v)
			metrics.UpstreamConnectSecondsHist.WithLabelValues(notCounterValues...).Observe(v)
		}

		if v, ok := observeMetrics(logger, fields, "request_time", floatFromFields, metrics.ParseErrorsTotal); ok {
			metrics.ResponseSeconds.WithLabelValues(notCounterValues...).Observe(v)
			metrics.ResponseSecondsHist.WithLabelValues(notCounterValues...).Observe(v)
		}
	}

	return nil
}

func observeMetrics(logger *log.Logger, fields map[string]string, name string, extractor func(map[string]string, string) (float64, bool, error), parseErrors prometheus.Counter) (float64, bool) {
	if observation, ok, err := extractor(fields, name); ok {
		return observation, true
	} else if err != nil {
		logger.Errorf("error while parsing $%s: %v", name, err)
		parseErrors.Inc()
	}

	return 0, false
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

	for _, v := range strings.FieldsFunc(val, func(r rune) bool { return r == ',' || r == ':' }) {
		v = strings.TrimSpace(v)

		if v == "-" {
			continue
		}

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

	if val == "-" {
		return 0, false, nil
	}

	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, false, fmt.Errorf("value '%s' could not be parsed into float", val)
	}

	return f, true, nil
}
