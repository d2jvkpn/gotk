package impls

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

/*
test_tracer0001:
- call01:
  - step01
  - call02

- call02
- job01
*/
func TestTracer(t *testing.T) {
	ctx := context.Background()

	// Write telemetry data to a file.
	_ = os.MkdirAll("wk", 0755)
	file, err := os.Create("wk/tracing.out")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(file),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		// stdouttrace.WithoutTimestamps(),
	)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("==>", "TestTracer")

	reso, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("TestTracer"),
			semconv.ServiceVersionKey.String("0.1.0"),
			attribute.String("what", "demo"),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(reso),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)

	//
	tracer := otel.Tracer("test_tracer")

	tCtx, span := tracer.Start(ctx, "test_tracer0001")
	fmt.Println(
		"~~~ test_tracer0001:",
		span.SpanContext().TraceID().String(),
		span.SpanContext().SpanID().String(),
	)

	tCtx1, span1 := tracer.Start(ctx, "test_tracer0001")
	call01(tCtx1)
	span1.End()

	call02(tCtx)

	_, span2 := tracer.Start(tCtx, "job01")
	time.Sleep(time.Second)
	span2.End()

	span.End()
}

func call01(ctx context.Context) {
	fmt.Println("==>", "call01")
	fmt.Printf("~~~ Context: %v\n", ctx)
	parentSpan := trace.SpanFromContext(ctx)

	labels := []attribute.KeyValue{
		attribute.String("orderId", "order0001"),
	}
	parentSpan.SetAttributes(labels...)

	fmt.Println(
		"~~~ ParentSpan:",
		parentSpan.SpanContext().TraceID().String(),
		parentSpan.SpanContext().SpanID().String(),
	)

	tracer := otel.Tracer("call01")

	// Create a span to track `childFunction()` - this is a nested span whose parent is `parentSpan`
	tCtx, span := tracer.Start(ctx, "step01")
	defer span.End()

	fmt.Println(
		"~~~ CurrentSpan:",
		span.SpanContext().TraceID().String(),
		span.SpanContext().SpanID().String(),
	)

	call02(tCtx)
	time.Sleep(3 * time.Second)

	opts := []trace.EventOption{
		trace.WithAttributes(attribute.Int64("count", 42)),
	}

	span.AddEvent("call01 is done", opts...)
	span.SetStatus(codes.Ok, "ok")
}

func call02(ctx context.Context) {
	fmt.Println("==>", "call02")
	fmt.Printf("~~~ Context: %v\n", ctx)

	parentSpan := trace.SpanFromContext(ctx)

	fmt.Println(
		"~~~ ParentSpan:",
		parentSpan.SpanContext().TraceID().String(),
		parentSpan.SpanContext().SpanID().String(),
	)

	time.Sleep(3 * time.Second)
	parentSpan.AddEvent("call02 is done")
}
