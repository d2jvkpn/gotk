package cloud

import (
	// "fmt"
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
)

func SetupOtelTracing(appName string, vp *viper.Viper) (
	shutdown func(context.Context) error, err error) {
	var (
		address  string
		grpcConn *grpc.ClientConn
	)

	address = vp.GetString("address")

	opts := []grpc.DialOption{grpc.WithTimeout(time.Second * 3)}
	//opts := []otlptracegrpc.Option{
	//	otlptracegrpc.WithEndpoint(addr),
	//	otlptracegrpc.WithDialOption(grpc.WithBlock()),
	//}
	if !vp.GetBool("tls") {
		opts = append(opts, grpc.WithInsecure())
	}
	grpcConn, err = grpc.DialContext(context.Background(), address, opts...)

	if err != nil {
		return nil, err
	}

	if shutdown, err = setupOtelTracing(grpcConn, appName); err != nil {
		return nil, err
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

	ctx = context.Background()

	client = otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	if exporter, err = otlptrace.New(ctx, client); err != nil {
		return nil, err
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
		return nil, err
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

func SetupOtelTracingFile(ctx context.Context, fp, service string, attrs ...attribute.KeyValue) (
	closeOtel func(context.Context) error, err error) {
	var (
		file     *os.File
		exporter *stdouttrace.Exporter
		reso     *resource.Resource
		provider *trace.TracerProvider
	)

	if err = os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			return
		}

		if exporter != nil {
			_ = exporter.Shutdown(ctx)
		}
		if file != nil {
			_ = file.Close()
		}
	}()

	if file, err = os.Create(fp); err != nil {
		return nil, err
	}

	exporter, err = stdouttrace.New(
		stdouttrace.WithWriter(file),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		// stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		return nil, err
	}

	attrs = append(attrs, semconv.ServiceNameKey.String(service))

	// reso, err := resource.Merge(resource.Default(), b)
	reso = resource.NewWithAttributes(semconv.SchemaURL, attrs...)

	provider = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(reso),
	)

	otel.SetTracerProvider(provider)

	closeOtel = func(ctx context.Context) error {
		var e1, e2, e3 error

		if e1 = provider.Shutdown(ctx); e1 != nil {
			otel.Handle(e1)
		}

		if e2 = exporter.Shutdown(ctx); e2 != nil {
			otel.Handle(e2)
		}

		if e3 = file.Close(); e3 != nil {
			otel.Handle(e3)
		}

		return errors.Join(e1, e2, e3)
	}

	return closeOtel, nil
}
