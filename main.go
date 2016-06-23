package main

import "fmt"
import (
	"flag"
	"github.com/hpcloud/tail"
	"github.com/satyrius/gonx"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
)

type StartOptions struct {
	Filenames  []string
	Format     string
	Namespace  string
	ListenPort int
}

type Metrics struct {
	countTotal      *prometheus.CounterVec
	bytesTotal      *prometheus.CounterVec
	upstreamSeconds *prometheus.SummaryVec
	responseSeconds *prometheus.SummaryVec
}

func (m *Metrics) Init(opts *StartOptions) {
	m.countTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: opts.Namespace,
		Name: "http_response_count_total",
		Help: "Amount of processed HTTP requests",
	}, []string{"method"})

	m.bytesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: opts.Namespace,
		Name: "http_response_size_bytes",
		Help: "Total amount of transferred bytes",
	}, []string{"method"})

	m.upstreamSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: opts.Namespace,
		Name: "http_upstream_time_seconds",
		Help: "Time needed by upstream servers to handle requests",
	}, []string{"method"})

	m.responseSeconds = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: opts.Namespace,
		Name: "http_response_time_seconds",
		Help: "Time needed by NGINX to handle requests",
	}, []string{"method"})

	prometheus.MustRegister(m.countTotal)
	prometheus.MustRegister(m.bytesTotal)
	prometheus.MustRegister(m.upstreamSeconds)
	prometheus.MustRegister(m.responseSeconds)
}

func main() {
	var opts StartOptions

	flag.IntVar(&opts.ListenPort, "listen-port", 4040, "HTTP port to listen on")
	flag.StringVar(&opts.Format, "format", `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"`, "NGINX access log format")
	flag.StringVar(&opts.Namespace, "namespace", "nginx", "namespace to use for metric names")
	flag.Parse()

	opts.Filenames = flag.Args()
	parser := gonx.NewParser(opts.Format)

	metrics := Metrics{}
	metrics.Init(&opts)

	for _, f := range opts.Filenames {
		t, err := tail.TailFile(f, tail.Config{
			Follow: true,
			ReOpen: true,
			Poll: true,
		})
		if err != nil {
			panic(err)
		}

		go func() {
			for line := range t.Lines {
				fmt.Printf("read from %s: %s\n", f, line)
				entry, err := parser.ParseString(line.Text)
				if err != nil {
					fmt.Printf("error while parsing line '%s': %s", line.Text, err)
					continue
				}

				var method = "UNKNOWN"
				if request, err := entry.Field("request"); err == nil {
					f := strings.Split(request, " ")
					method = f[0]
				}

				metrics.countTotal.WithLabelValues(method).Inc()

				if bytes, err := entry.FloatField("body_bytes_sent"); err == nil {
					metrics.bytesTotal.WithLabelValues(method).Add(bytes)
				}

				if upstreamTime, err := entry.FloatField("upstream_response_time"); err == nil {
					metrics.upstreamSeconds.WithLabelValues(method).Observe(upstreamTime)
				}

				if responseTime, err := entry.FloatField("request_time"); err == nil {
					metrics.responseSeconds.WithLabelValues(method).Observe(responseTime)
				}

				fmt.Println(entry)
			}
		}()
	}

	listenAddr := fmt.Sprintf("%s:%d", "0.0.0.0", opts.ListenPort)

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(listenAddr, nil)
}