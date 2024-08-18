package cloud

import (
	// "fmt"

	. "github.com/d2jvkpn/gotk/trace_error"

	"github.com/prometheus/client_golang/prometheus"
)

func PromMetricsAPI() (func(string, float64, *Error), error) {
	var (
		e              error
		codeCounter    *prometheus.CounterVec
		requestLatency *prometheus.HistogramVec
		// concurrentRequests prometheus.Gauge
	)

	// http_code, http response code summary
	codeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_code",
			Help: "HTTP response code and kind counter",
		},
		[]string{"code", "kind"},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_latency",
			Help:    "HTTP response latency milliseconds",
			Buckets: []float64{10, 100, 200, 500, 1000, 5000}, // prometheus.DefBuckets,
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
		if e = prometheus.Register(v); e != nil {
			return nil, e
		}
	}

	return func(api string, latency float64, err *Error) {
		var labelValues [2]string

		labelValues[0], labelValues[1] = "OK", ""
		if err != nil {
			labelValues[0], labelValues[1] = err.Code, err.Kind
		}

		codeCounter.WithLabelValues(labelValues[:]...).Inc()
		requestLatency.WithLabelValues(api).Observe(latency)
		// concurrentRequests.Dec()
	}, nil
}
