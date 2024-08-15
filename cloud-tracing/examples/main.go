// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const meterName = "go.opentelemetry.io/otel/example/prometheus"

func main() {
	var (
		err       error
		ctx       context.Context
		rng       *rand.Rand
		exporter  *prometheus.Exporter
		provider  *sdkmetric.MeterProvider
		meter     otelmetric.Meter
		opt       otelmetric.MeasurementOption
		counter   otelmetric.Float64Counter
		gauge     otelmetric.Float64ObservableGauge
		histogram otelmetric.Float64Histogram
	)

	ctx = context.Background()
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	if exporter, err = prometheus.New(); err != nil {
		log.Fatal(err)
	}
	provider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter = provider.Meter(meterName)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	opt = otelmetric.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	// This is the equivalent of prometheus.NewCounterVec
	counter, err = meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, opt)

	gauge, err = meter.Float64ObservableGauge(
		"bar",
		otelmetric.WithDescription("a fun little gauge"),
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o otelmetric.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err = meter.Float64Histogram(
		"baz",
		otelmetric.WithDescription("a histogram with custom buckets and rename"),
		otelmetric.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
	)
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 136, opt)
	histogram.Record(ctx, 64, opt)
	histogram.Record(ctx, 701, opt)
	histogram.Record(ctx, 830, opt)

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
