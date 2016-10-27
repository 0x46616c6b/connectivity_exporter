package collector

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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
func NewHTTPExporter(hosts []string, timeout time.Duration) *HTTPExporter {
	return &HTTPExporter{
		hosts: hosts,
		client: &http.Client{
			Timeout: timeout,
		},
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
	var wg sync.WaitGroup

	for _, host := range p.hosts {
		wg.Add(1)

		go func(host string, ch chan<- prometheus.Metric) {
			defer wg.Done()

			s, url := 1, host
			if !strings.HasPrefix(host, "http") {
				url = "https://" + host
			}

			t := time.Now()
			res, err := p.client.Get(url)
			if err != nil {
				log.Errorln(err)
				s = 0
			} else {
				defer res.Body.Close()
			}
			d := time.Since(t)

			ch <- prometheus.MustNewConstMetric(
				httpRequestSuccessful, prometheus.GaugeValue, float64(s), host,
			)

			ch <- prometheus.MustNewConstMetric(
				httpRequestTimeNS, prometheus.GaugeValue, float64(d.Nanoseconds()), host,
			)
		}(host, ch)
	}

	wg.Wait()
}
