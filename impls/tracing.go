package impls

import (
	// "fmt"
	"context"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func LoadOtel(addr, service string, secure bool) (closeOtel func(), err error) {
	var (
		client   otlptrace.Client
		exporter *otlptrace.Exporter
		reso     *resource.Resource
		provider *sdktrace.TracerProvider
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(addr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	}
	if !secure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	client = otlptracegrpc.NewClient(opts...)

	if exporter, err = otlptrace.New(ctx, client); err != nil {
		return nil, err
	}

	reso, err = resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(service),
		),
	)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	provider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(reso),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(provider)

	return func() {
		var (
			ctx    context.Context
			cancel func()
		)
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := exporter.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}, nil
}
