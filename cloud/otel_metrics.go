package cloud

import (
	// "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/d2jvkpn/gotk/ginx"
	"github.com/d2jvkpn/gotk/trace_error"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func OtelMetrics(meter otelmetric.Meter) (hf gin.HandlerFunc, err error) {
	var (
		codeCounter    otelmetric.Float64Counter
		requestLatency otelmetric.Float64Histogram
	)

	// http_code, http response code summary
	codeCounter, err = meter.Float64Counter(
		"http_code", // suffix _total added
		otelmetric.WithDescription("HTTP response code and kind counter"),
	)
	if err != nil {
		return nil, err
	}

	requestLatency, err = meter.Float64Histogram(
		"http_latency",
		otelmetric.WithDescription("HTTP response latency milliseconds"),
		otelmetric.WithExplicitBucketBoundaries(10.0, 100.0, 200.0, 500.0, 1000.0, 5000.0),
	)
	if err != nil {
		return nil, err
	}

	hf = func(ctx *gin.Context) {
		var (
			e           error
			api         string
			status      int
			start       time.Time
			err         *trace_error.Error
			labelValues [2]string
		)

		// concurrentRequests.Inc()
		api = fmt.Sprintf("%s@%s", ctx.Request.Method, ctx.Request.URL.Path)

		start = time.Now()
		ctx.Next()

		status = ctx.Writer.Status()
		latency := float64(time.Since(start).Microseconds()) / 1e3

		labelValues[0], labelValues[1] = "OK", ""
		if status != http.StatusOK {
			if err, e = ginx.Get[*trace_error.Error](ctx, "Error"); e == nil {
				labelValues[0], labelValues[1] = err.Code, err.Kind
			}
		}

		codeCounter.Add(ctx, 1, otelmetric.WithAttributes(
			attribute.Key("code").String(labelValues[0]),
			attribute.Key("kind").String(labelValues[1]),
		))

		requestLatency.Record(ctx, latency, otelmetric.WithAttributes(
			attribute.Key("api").String(api),
		))
		// concurrentRequests.Dec()
	}

	return hf, nil
}
