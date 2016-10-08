package collector

import (
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "connectivity"
)

var (
	httpRequestSuccessful = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "http_request_successful"),
		"Boolean Metric with 1 if the request was successful",
		[]string{"host"}, nil,
	)
	httpRequestTimeNS = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "http_request_time_ns"),
		"Duration of the request",
		[]string{"host"}, nil,
	)
)

// HTTPExporter collects HTTP stats and exports them using
// the prometheus metrics package.
type HTTPExporter struct {
	hosts  []string
	client *http.Client
}

// NewHTTPExporter returns an initialized Exporter.
func NewHTTPExporter(hosts []string) *HTTPExporter {
	return &HTTPExporter{
		hosts:  hosts,
		client: &http.Client{},
	}
}

// Describe describes all the metrics ever exported. It
// implements prometheus.Collector.
func (p *HTTPExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- httpRequestSuccessful
	ch <- httpRequestTimeNS
}

// Collect fetches the Fastping stats and delivers them as
// Prometheus metrics. It implements promethues.Collector.
func (p *HTTPExporter) Collect(ch chan<- prometheus.Metric) {
	for _, host := range p.hosts {
		var s float64
		s = 1
		url := host

		if !strings.HasPrefix(host, "http") {
			url = "https://" + host
		}

		t := time.Now()
		_, err := p.client.Get(url)
		if err != nil {
			log.Errorln(err)
			s = 0
		}
		d := time.Since(t)

		ch <- prometheus.MustNewConstMetric(
			httpRequestSuccessful, prometheus.GaugeValue, s, host,
		)

		ch <- prometheus.MustNewConstMetric(
			httpRequestTimeNS, prometheus.GaugeValue, float64(d.Nanoseconds()), host,
		)
	}
}
