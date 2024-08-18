package cloud

import (
	// "fmt"
	"context"
	"errors"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
)

func SetupOtelTracing(appName string, vp *viper.Viper) (
	shutdown func(context.Context) error, err error) {
	var grpcConn *grpc.ClientConn

	shutdown = func(context.Context) error { return nil }
	if !vp.GetBool("trace") {
		return shutdown, nil
	}

	opts := []grpc.DialOption{grpc.WithTimeout(time.Second * 3)}
	if !vp.GetBool("tls") {
		opts = append(opts, grpc.WithInsecure())
	}

	grpcConn, err = grpc.DialContext(
		context.Background(), vp.GetString("address"), opts...,
	)
	if err != nil {
		return shutdown, err
	}

	if shutdown, err = setupOtelTracing(grpcConn, appName); err != nil {
		return shutdown, err
	}

	return shutdown, nil
}

// conn, err := grpc.DialContext(ctx, "collector:4317", grpc.WithInsecure())
func setupOtelTracing(conn *grpc.ClientConn, service string, attrs ...attribute.KeyValue) (
	shutdown func(context.Context) error, err error) {
	var (
		ctx      context.Context
		client   otlptrace.Client
		exporter *otlptrace.Exporter
		reso     *resource.Resource
		provider *trace.TracerProvider
	)

	shutdown = func(context.Context) error { return nil }
	ctx = context.Background()

	client = otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	if exporter, err = otlptrace.New(ctx, client); err != nil {
		return shutdown, err
	}
	defer func() {
		if err == nil {
			return
		}

		_ = exporter.Shutdown(ctx)
	}()

	attrs = append(attrs, semconv.ServiceNameKey.String(service))
	reso, err = resource.New(ctx,
		// resource.WithFromEnv(),
		// resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return shutdown, err
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	provider = trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(reso),
		trace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(provider)

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
