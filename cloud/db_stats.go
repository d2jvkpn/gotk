package cloud

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/d2jvkpn/gotk"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

func SetupDBStatsProm(db *sql.DB) (tickerDBConn *gotk.Ticker, err error) {
	var gauge *prometheus.GaugeVec

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	//
	gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_conns",
		Help: "Number of connections of database",
	}, []string{"type"})

	if err = prometheus.Register(gauge); err != nil {
		return nil, err
	}

	dbStats := func() {
		var stats sql.DBStats

		stats = db.Stats()
		// fmt.Printf("==> DBStats: %+v\n", stats)
		gauge.With(prometheus.Labels{"type": "in_user"}).Set(float64(stats.InUse))
		gauge.With(prometheus.Labels{"type": "idle"}).Set(float64(stats.Idle))

		gauge.With(prometheus.Labels{"type": "open_connections"}).
			Set(float64(stats.OpenConnections))
	}

	tickerDBConn = gotk.NewTicker([]func(){dbStats}, 15*time.Second)

	return tickerDBConn, nil
}

func SetupDBStatsOtel(db *sql.DB, meter otelmetric.Meter) (err error) {
	var gauge otelmetric.Float64ObservableGauge

	if db == nil {
		return fmt.Errorf("db is nil")
	}

	//
	gauge, err = meter.Float64ObservableGauge(
		"db_conns",
		otelmetric.WithDescription("Number of connections of database"),
	)

	_, err = meter.RegisterCallback(func(_ context.Context,
		o otelmetric.Observer) error {
		// println("==> DBStatsMetrics")

		var stats sql.DBStats

		stats = db.Stats()

		o.ObserveFloat64(
			gauge, float64(stats.InUse),
			otelmetric.WithAttributes(attribute.Key("type").String("in_user")),
		)

		o.ObserveFloat64(
			gauge, float64(stats.Idle),
			otelmetric.WithAttributes(attribute.Key("type").String("idle")),
		)

		o.ObserveFloat64(
			gauge, float64(stats.OpenConnections),
			otelmetric.WithAttributes(attribute.Key("type").String("open_connections")),
		)

		return nil
	}, gauge)

	if err != nil {
		return nil
	}

	return nil
}
