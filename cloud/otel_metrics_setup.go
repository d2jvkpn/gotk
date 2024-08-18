package cloud

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// without export to otel-collector
func setupOtelPrometheus(appName string, vp *viper.Viper) (otelmetric.Meter, error) {
	var (
		err      error
		exporter *otelprometheus.Exporter
		provider *sdkmetric.MeterProvider
	)

	if exporter, err = otelprometheus.New(); err != nil {
		return nil, err
	}
	provider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))

	return provider.Meter(appName), nil
}

// https://opentelemetry.io/docs/languages/go/getting-started/
func SetupOtelMetrics(appName string, vp *viper.Viper) (
	otelmetric.Meter, func(context.Context) error, error) {
	var (
		err      error
		ctx      context.Context
		exporter *otlpmetricgrpc.Exporter
		reso     *resource.Resource
		provider *sdkmetric.MeterProvider
		shutdown func(context.Context) error
		meter    otelmetric.Meter
	)

	ctx = context.Background()
	shutdown = func(context.Context) error { return nil }

	// println("==>", vp.GetBool("metrics"))
	if !vp.GetBool("metrics") {
		meter, err = setupOtelPrometheus(appName, vp)
		return meter, shutdown, err
	}

	reso, err = resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(appName)),
	)
	if err != nil {
		return nil, shutdown, err
	}

	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(vp.GetString("address"))}
	if !vp.GetBool("tls") {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	if exporter, err = otlpmetricgrpc.New(ctx, opts...); err != nil {
		return nil, shutdown, err
	}

	provider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			exporter, sdkmetric.WithInterval(10*time.Second),
		)),
		sdkmetric.WithResource(reso),
	)
	otel.SetMeterProvider(provider)

	return provider.Meter(appName), exporter.Shutdown, nil
}
