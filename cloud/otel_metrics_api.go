package cloud

import (
	"context"
	// "fmt"
	"time"

	"github.com/d2jvkpn/gotk/trace_error"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func OtelMetricsAPI(meter otelmetric.Meter) (func(string, float64, *trace_error.Error), error) {
	var (
		e              error
		codeCounter    otelmetric.Float64Counter
		requestLatency otelmetric.Float64Histogram
	)

	// http_code, http response code summary
	codeCounter, e = meter.Float64Counter(
		"http_code", // suffix _total added
		otelmetric.WithDescription("HTTP response code and kind counter"),
	)
	if e != nil {
		return nil, e
	}

	requestLatency, e = meter.Float64Histogram(
		"http_latency",
		otelmetric.WithDescription("HTTP response latency milliseconds"),
		otelmetric.WithExplicitBucketBoundaries(10.0, 100.0, 200.0, 500.0, 1000.0, 5000.0),
	)
	if e != nil {
		return nil, e
	}

	return func(api string, latency float64, err *trace_error.Error) {
		var labelValues [3]string

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		labelValues[0], labelValues[1], labelValues[2] = "OK", "", api
		if err != nil {
			labelValues[0], labelValues[1] = err.Code, err.Kind
		}

		codeCounter.Add(ctx, 1, otelmetric.WithAttributes(
			attribute.Key("code").String(labelValues[0]),
			attribute.Key("kind").String(labelValues[1]),
			attribute.Key("api").String(labelValues[2]),
		))

		requestLatency.Record(ctx, latency, otelmetric.WithAttributes(
			attribute.Key("api").String(api),
		))
		// concurrentRequests.Dec()
	}, nil
}
