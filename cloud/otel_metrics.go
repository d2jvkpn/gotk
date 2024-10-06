package cloud

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// not export to otel-collector, but export metrics to promethus http handler(/metrics)
func OtelMetrics2Prom(appName string, vp *viper.Viper) (otelmetric.Meter, error) {
	var (
		err      error
		exporter *otelprometheus.Exporter
		provider *sdkmetric.MeterProvider
	)

	if exporter, err = otelprometheus.New(); err != nil {
		return nil, err
	}
	provider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))

	//	if withRuntime {
	//		err = runtime.Start(
	//			runtime.WithMeterProvider(provider),
	//			runtime.WithMinimumReadMemStatsInterval(15*time.Second),
	//		)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}

	return provider.Meter(appName), nil
}

// https://opentelemetry.io/docs/languages/go/getting-started/
// get otelmetric.Meter by otel.GetMeterProvider()
func OtelMetricsGrpc(appName string, vp *viper.Viper, withRuntime bool) (
	shutdown func(context.Context) error, err error) {
	var (
		ctx      context.Context
		exporter *otlpmetricgrpc.Exporter
		reso     *resource.Resource
		provider *sdkmetric.MeterProvider
	)

	ctx = context.Background()
	shutdown = func(context.Context) error { return nil }

	reso, err = resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(appName)),
	)
	if err != nil {
		return shutdown, fmt.Errorf("resource.New: %w", err) // nil, shutdown, err
	}

	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(vp.GetString("address"))}
	if !vp.GetBool("tls") {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	if exporter, err = otlpmetricgrpc.New(ctx, opts...); err != nil {
		return shutdown, fmt.Errorf("otlpmetricgrpc.New: %w", err) // nil, shutdown, err
	}

	provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			exporter, sdkmetric.WithInterval(15*time.Second),
		)),
		sdkmetric.WithResource(reso),
	)

	// set global
	otel.SetMeterProvider(provider)

	if withRuntime {
		err = runtime.Start(
			runtime.WithMeterProvider(provider),
			runtime.WithMinimumReadMemStatsInterval(15*time.Second),
		)
		if err != nil {
			return shutdown, fmt.Errorf("runtime.Start: %w", err)
		}
	}

	shutdown = func(ctx context.Context) error {
		var e1, e2 error

		if e1 = provider.Shutdown(ctx); e1 != nil {
			otel.Handle(e1)
		}

		if e2 = exporter.Shutdown(ctx); e2 != nil {
			otel.Handle(e2)
		}

		return errors.Join(e1, e2)
	}

	return shutdown, nil
}

func OtelMetricsHttp(meter otelmetric.Meter, attrs []string) (
	fn func(string, float64, []string), err error) {
	var (
		codeCounter    otelmetric.Float64Counter
		requestLatency otelmetric.Float64Histogram
	)

	// http_code, http response code summary
	codeCounter, err = meter.Float64Counter(
		"http_code", // suffix _total added
		otelmetric.WithDescription("HTTP response codes counter"),
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

	return func(api string, latency float64, values []string) {
		var attributes []attribute.KeyValue

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		attributes = make([]attribute.KeyValue, 0, 1+len(attrs))

		attributes = append(attributes, attribute.Key("api").String(api))

		for i := 0; i < min(len(attrs), len(values)); i++ {
			attributes = append(attributes, attribute.Key(attrs[i]).String(values[i]))
		}

		codeCounter.Add(ctx, 1, otelmetric.WithAttributes(attributes...))

		requestLatency.Record(ctx, latency, otelmetric.WithAttributes(
			attribute.Key("api").String(api),
		))
		// concurrentRequests.Dec()
	}, nil
}
