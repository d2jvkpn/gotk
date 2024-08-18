package cloud

import (
	"fmt"
	"net/http"
	"time"

	"github.com/d2jvkpn/gotk/ginx"
	"github.com/d2jvkpn/gotk/trace_error"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// https://prometheus.io/docs/prometheus/latest/querying/examples/
// https://robert-scherbarth.medium.com/measure-request-duration-with-prometheus-and-golang-adc6f4ca05fe
func PromMetrics(errKey string) (hf gin.HandlerFunc, err error) {
	var (
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

			inflight := prometheus.NewGauge(prometheus.GaugeOpts{
				Name:      "http_inflight",
				Help:      "http inflight",
		})
	*/

	hf = func(ctx *gin.Context) {
		var (
			status      int
			e           error
			latency     float64
			api         string
			labelValues [2]string
			start       time.Time
			err         *trace_error.Error
		)

		// inflight.Inc()
		// defer inflight.Dec()

		api = fmt.Sprintf("%s@%s", ctx.Request.Method, ctx.Request.URL.Path)
		start = time.Now()

		ctx.Next()

		status = ctx.Writer.Status()
		latency = float64(time.Since(start).Microseconds()) / 1e3

		labelValues[0], labelValues[1] = "OK", ""
		if status != http.StatusOK {
			if err, e = ginx.Get[*trace_error.Error](ctx, errKey); e == nil {
				labelValues[0], labelValues[1] = err.Code, err.Kind
			}
		}

		codeCounter.WithLabelValues(labelValues[:]...).Inc()
		requestLatency.WithLabelValues(api).Observe(latency)
		// concurrentRequests.Dec()
	}

	// concurrentRequests
	for _, v := range []prometheus.Collector{codeCounter, requestLatency} {
		if err = prometheus.Register(v); err != nil {
			return nil, err
		}
	}

	return hf, nil
}
