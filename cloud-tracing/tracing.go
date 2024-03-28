package tracing

import (
	// "fmt"
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
)

func LoadOtelGrpc(addr, service string, secure bool, attrs ...attribute.KeyValue) (
	closeOtel func() error, err error) {
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
	defer func() {
		if err == nil {
			return
		}

		_ = exporter.Shutdown(context.TODO())
	}()

	attrs = append(attrs, semconv.ServiceNameKey.String(service))
	reso, err = resource.New(ctx,
		resource.WithFromEnv(),
		// resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(attrs...),
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

	closeOtel = func() error {
		var (
			e1, e2 error
			ctx    context.Context
			cancel func()
		)
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if e1 = provider.Shutdown(ctx); e1 != nil {
			otel.Handle(e1)
		}

		if e2 = exporter.Shutdown(ctx); e2 != nil {
			otel.Handle(e2)
		}

		return errors.Join(e1, e2)
	}

	return closeOtel, nil
}

func LoadOtelFile(fp, service string, attrs ...attribute.KeyValue) (
	closeOtel func() error, err error) {
	var (
		file     *os.File
		exporter *stdouttrace.Exporter
		reso     *resource.Resource
		provider *sdktrace.TracerProvider
	)

	if err = os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			return
		}

		if exporter != nil {
			_ = exporter.Shutdown(context.TODO())
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

	provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(reso),
	)

	otel.SetTracerProvider(provider)

	closeOtel = func() error {
		var (
			e1, e2, e3 error
			ctx        context.Context
			cancel     func()
		)
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

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
