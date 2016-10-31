/*
 * Copyright 2016 Martin Helmich <kontakt@martin-helmich.de>
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
	"strings"

	"github.com/hpcloud/tail"
	"github.com/martin-helmich/prometheus-nginxlog-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/satyrius/gonx"
)

// Metrics is a struct containing pointers to all metrics that should be
// exposed to Prometheus
type Metrics struct {
	countTotal      *prometheus.CounterVec
	bytesTotal      *prometheus.CounterVec
	upstreamSeconds *prometheus.SummaryVec
	responseSeconds *prometheus.SummaryVec
}

// Init initializes a metrics struct
func (m *Metrics) Init(cfg *config.NamespaceConfig) {
	labels := []string{"method", "status"}

	m.countTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: cfg.Name,
		Name:      "http_response_count_total",
		Help:      "Amount of processed HTTP requests",
	}, labels)

	m.bytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: cfg.Name,
		Name:      "http_response_size_bytes",
		Help:      "Total amount of transferred bytes",
	}, labels)

	m.upstreamSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: cfg.Name,
		Name:      "http_upstream_time_seconds",
		Help:      "Time needed by upstream servers to handle requests",
	}, labels)

	m.responseSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: cfg.Name,
		Name:      "http_response_time_seconds",
		Help:      "Time needed by NGINX to handle requests",
	}, labels)

	prometheus.MustRegister(m.countTotal)
	prometheus.MustRegister(m.bytesTotal)
	prometheus.MustRegister(m.upstreamSeconds)
	prometheus.MustRegister(m.responseSeconds)
}

func main() {
	var opts config.StartupFlags
	var cfg = config.Config{
		Port: 4040,
	}

	flag.IntVar(&opts.ListenPort, "listen-port", 4040, "HTTP port to listen on")
	flag.StringVar(&opts.Format, "format", `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"`, "NGINX access log format")
	flag.StringVar(&opts.Namespace, "namespace", "nginx", "namespace to use for metric names")
	flag.StringVar(&opts.ConfigFile, "config-file", "", "Configuration file to read from")
	flag.Parse()

	opts.Filenames = flag.Args()

	if opts.ConfigFile != "" {
		fmt.Printf("loading configuration file %s\n", opts.ConfigFile)
		if err := config.LoadConfigFromFile(&cfg, opts.ConfigFile); err != nil {
			panic(err)
		}
	} else if err := config.LoadConfigFromFlags(&cfg, &opts); err != nil {
		panic(err)
	}

	fmt.Printf("using configuration %s\n", cfg)

	for _, ns := range cfg.Namespaces {
		fmt.Printf("starting listener for namespace %s\n", ns.Name)

		go func(nsCfg *config.NamespaceConfig) {
			parser := gonx.NewParser(nsCfg.Format)

			metrics := Metrics{}
			metrics.Init(nsCfg)

			for _, f := range nsCfg.SourceFiles {
				t, err := tail.TailFile(f, tail.Config{
					Follow: true,
					ReOpen: true,
					Poll:   true,
				})
				if err != nil {
					panic(err)
				}

				go func() {
					for line := range t.Lines {
						entry, err := parser.ParseString(line.Text)
						if err != nil {
							fmt.Printf("error while parsing line '%s': %s", line.Text, err)
							continue
						}

						method := "UNKNOWN"
						status := "0"

						if request, err := entry.Field("request"); err == nil {
							f := strings.Split(request, " ")
							method = f[0]
						}

						if s, err := entry.Field("status"); err == nil {
							status = s
						}

						metrics.countTotal.WithLabelValues(method, status).Inc()

						if bytes, err := entry.FloatField("body_bytes_sent"); err == nil {
							metrics.bytesTotal.WithLabelValues(method, status).Add(bytes)
						}

						if upstreamTime, err := entry.FloatField("upstream_response_time"); err == nil {
							metrics.upstreamSeconds.WithLabelValues(method, status).Observe(upstreamTime)
						}

						if responseTime, err := entry.FloatField("request_time"); err == nil {
							metrics.responseSeconds.WithLabelValues(method, status).Observe(responseTime)
						}
					}
				}()
			}
		}(&ns)
	}

	listenAddr := fmt.Sprintf("%s:%d", "0.0.0.0", opts.ListenPort)
	fmt.Printf("running HTTP server on address %s\n", listenAddr)

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(listenAddr, nil)
}
