package collector

import (
	"net"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	dnsRequestSuccessful = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dns_request_successful"),
		"Boolean Metric with 1 if the request was successful",
		[]string{"host"}, nil,
	)
	dnsRequestTimeNS = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dns_request_time_ns"),
		"Duration of the request",
		[]string{"host"}, nil,
	)
)

// DNSExporter collects DNS stats and exports them using
// the prometheus metrics package.
type DNSExporter struct {
	hosts []string
}

// NewDNSExporter returns an initialized Exporter.
func NewDNSExporter(hosts []string) *DNSExporter {
	return &DNSExporter{
		hosts: hosts,
	}
}

// Describe describes all the metrics ever exported. It
// implements prometheus.Collector.
func (p *DNSExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- dnsRequestSuccessful
	ch <- dnsRequestTimeNS
}

// Collect fetches the Fastping stats and delivers them as
// Prometheus metrics. It implements promethues.Collector.
func (p *DNSExporter) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for _, host := range p.hosts {
		wg.Add(1)

		go func(host string, ch chan<- prometheus.Metric) {
			defer wg.Done()

			s := 1
			t := time.Now()
			_, err := net.LookupHost(host)
			if err != nil {
				log.Errorln(err)
				s = 0
			}
			d := time.Since(t)

			ch <- prometheus.MustNewConstMetric(
				dnsRequestSuccessful, prometheus.GaugeValue, float64(s), host,
			)

			ch <- prometheus.MustNewConstMetric(
				dnsRequestTimeNS, prometheus.GaugeValue, float64(d.Nanoseconds()), host,
			)
		}(host, ch)
	}

	wg.Wait()
}
