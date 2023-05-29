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

// #### Usage demo:
//
//	```go
//	import (
//		"go.opentelemetry.io/otel/attribute"
//		"go.opentelemetry.io/otel/trace"
//	)
//
//	//==> funcA(ctx context.Context)
//	span := trace.SpanFromContext(ctx)
//	traceId := span.SpanContext().TraceID().String()
//	spanId := span.SpanContext().SpanID().String()
//	fmt.Println("~~~ funcA:", traceId, spanId)
//	funcB(ctx, "order0001")
//
//	//==> funcB(ctx, oid string)
//	span := trace.SpanFromContext(ctx)
//	labels := []attribute.KeyValue{
//		attribute.String("orderId", oid),
//	}
//	span.SetAttributes(labels...)
//	funcC(ctx)
//
//	//==> funcC(ctx context.Context)
//	tracer := otel.Tracer("service-c")
//	_, span := tracer.Start(ctx, "c1")
//	defer span.End()
//	time.Sleep(3*time.Second)
//	opts := []trace.EventOption{
//		trace.WithAttributes(attribute.Int64("count", 42)),
//	}
//	span.AddEvent("successfully finished call service-c", opts...)
//	```
//
//	```go
//	engi.Use(otelgin.Middleware("service"))
//	```
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
