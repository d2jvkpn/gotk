package impls

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
func TestLoadOtelFile(t *testing.T) {
	shutdown, err := LoadOtelFile(
		"wk/tracing.out",
		"TestTracer",
		semconv.ServiceVersionKey.String("0.1.0"),
		attribute.String("what", "demo"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := shutdown(); err != nil {
			t.Fatal("shutdwon:", err)
		}
	}()

	tracer := otel.Tracer("test_tracer")

	tCtx, span := tracer.Start(context.Background(), "test_tracer0001")
	fmt.Println(
		"~~~ test_tracer0001:",
		span.SpanContext().TraceID().String(),
		span.SpanContext().SpanID().String(),
	)

	tCtx1, span1 := tracer.Start(tCtx, "test_tracer0001")
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
