package cloud

import (
	"context"
	// "fmt"
	"testing"

	// otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"google.golang.org/grpc"
)

func TestOtelGrpcDial(t *testing.T) {
	var (
		addr string
		err  error
		// exporter *otlpmetricgrpc.Exporter
	)

	addr = "localhost:4317"

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(addr),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithInsecure(),
	}
	if _, err = otlpmetricgrpc.New(context.Background(), opts...); err != nil {
		t.Fatal(err)
	}

	// grpc.Block: Deprecated: this DialOption is not supported by NewClient. Will be supported
}
