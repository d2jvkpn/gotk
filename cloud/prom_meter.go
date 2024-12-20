package cloud

import (
	// "fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func PromMeterHttp(labels []string) (fn func(string, float64, []string), err error) {
	var (
		codeCounter    *prometheus.CounterVec
		requestLatency *prometheus.HistogramVec
		// concurrentRequests prometheus.Gauge
	)

	// http_code, http response code summary
	codeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_code_total",
			Help: "HTTP response codes counter",
		},
		append([]string{"api"}, labels...),
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_latency",
			Help:    "HTTP response latency milliseconds",
			Buckets: []float64{10.0, 100.0, 200.0, 500.0, 1000.0, 5000.0}, // prometheus.DefBuckets,
		},
		[]string{"api"},
	)
	/*
		concurrentRequests = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_concurrent",
			Help: "Number of concurrent requests being processed",
		})
	*/

	for _, v := range []prometheus.Collector{codeCounter, requestLatency} {
		if err = prometheus.Register(v); err != nil {
			return nil, err
		}
	}

	return func(api string, latency float64, values []string) {
		labelValues := append([]string{"api"}, values...)

		codeCounter.WithLabelValues(labelValues...).Inc()
		requestLatency.WithLabelValues(api).Observe(latency)
		// concurrentRequests.Dec()
	}, nil
}
