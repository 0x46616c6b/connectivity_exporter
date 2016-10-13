package main

import (
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/0x46616c6b/connectivity_exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// Exporter collects Connectivity stats and exports them using
// the prometheus metrics package.
type Exporter struct {
	exporters []prometheus.Collector
}

// NewExporter returns an initialized Exporter wit all Collectors.
func NewExporter(hosts []string, timeout time.Duration) (*Exporter, error) {
	return &Exporter{
		[]prometheus.Collector{
			collector.NewHTTPExporter(hosts, timeout),
			collector.NewDNSExporter(hosts),
		},
	}, nil
}

// Describe describes all the metrics ever exported. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, exporter := range e.exporters {
		exporter.Describe(ch)
	}
}

// Collect fetches the stats from configured Modules and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, exporter := range e.exporters {
		exporter.Collect(ch)
	}
}

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9449", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		httpHosts     = flag.String("http.hosts", "google.com,facebook.com,github.com", "Comma seperated list with hosts to check")
		httpTimeout   = flag.String("http.timeout", "5s", "Timeout for the HTTP Checks")
	)
	flag.Parse()

	hosts := strings.Split(*httpHosts, ",")
	timeout, err := time.ParseDuration(*httpTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	exporter, err := NewExporter(hosts, timeout)
	if err != nil {
		log.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Connectivity Exporter</title></head>
             <body>
             <h1>Connectivity Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
