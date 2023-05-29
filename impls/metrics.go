package impls

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/d2jvkpn/gotk"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

/*
	yaml config:

```yaml
prometheus: true
debug: true
addr: :8080
```

routes:
- /healthz
- /metrics
- /debug/meta
- /debug/pprof/expvar
- /debug/pprof/
- /debug/pprof/profile
- /debug/pprof/trace
- /debug/pprof/cmdline
- /debug/pprof/symbol
- /debug/pprof/allocs
- /debug/pprof/block
- /debug/pprof/goroutine
- /debug/pprof/heap
- /debug/pprof/mutex
- /debug/pprof/threadcreate
*/
func HttpMetrics(vp *viper.Viper, meta map[string]any, opts ...func(*http.Server)) (
	shutdown func() error, err error) {
	var (
		enableProm  bool
		enableDebug bool
		addr        string
		listener    net.Listener
		// mux      *http.ServeMux
		engine *gin.Engine
		router *gin.RouterGroup
		server *http.Server
	)

	enableProm = vp.GetBool("prometheus")
	enableDebug = vp.GetBool("debug")
	addr = vp.GetString("addr")

	if listener, err = net.Listen("tcp", addr); err != nil {
		return nil, err
	}

	// mux = http.NewServeMux()
	// mux.Handle("/metrics", promhttp.Handler())

	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	// engine.Use(gin.Recovery())
	engine.Use(Cors(vp.GetString("*")))
	router = &engine.RouterGroup

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	if enableProm {
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	if enableDebug {
		debug := router.Group("/debug")

		debug.GET("/meta", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, meta)
		})

		for k, f := range gotk.PprofHandlerFuncs() {
			debug.GET(fmt.Sprintf("/pprof/%s", k), gin.WrapF(f))
		}
	}

	server = &http.Server{
		// ReadTimeout:       time.Second * 30,
		// WriteTimeout:      time.Minute * 5,
		// ReadHeaderTimeout: time.Second * 2,
		// MaxHeaderBytes:    1 << 20,
		// Addr:              addr,
		Handler: engine,
	}

	for i := range opts {
		opts[i](server)
	}

	shutdown = func() error {
		ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
		err := server.Shutdown(ctx)
		cancel()
		return err
	}

	go func() {
		if server.TLSConfig != nil {
			_ = server.ServeTLS(listener, "", "")
		} else {
			_ = server.Serve(listener)
		}
	}()

	return shutdown, nil
}

// https://prometheus.io/docs/prometheus/latest/querying/examples/
// https://robert-scherbarth.medium.com/measure-request-duration-with-prometheus-and-golang-adc6f4ca05fe
func PromMetrics() (hf gin.HandlerFunc, err error) {
	requestsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "http",
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"status"})

	// promauto.NewHistogramVec
	requestsDuration := promauto.NewHistogram(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "http_requests_duration",
		Help:      "Request duration(milliseconds) of HTTP requests",
		Buckets:   prometheus.DefBuckets,
	})

	/*
		requestDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: "http",
			Name:      "http_requests_duration",
			Help:      "Request duration(milliseconds) of HTTP requests",
			Buckets:   prometheus.DefBuckets,
		}, []string{"path"})
	*/

	inflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: "http",
		Name:      "http_requests_inflight",
		Help:      "Total inflight of HTTP requests",
	})

	if err = prometheus.Register(requestsTotal); err != nil {
		return nil, err
	}

	if err = prometheus.Register(requestsDuration); err != nil {
		return nil, err
	}

	if err = prometheus.Register(inflight); err != nil {
		return nil, err
	}

	hf = func(ctx *gin.Context) {
		// req := ctx.Request
		inflight.Inc()
		defer inflight.Dec()

		start := time.Now()
		// timer := prometheus.NewTimer(requestDuration.WithLabelValues(req.Method, req.URL.Path))
		ctx.Next()

		status := strconv.Itoa(ctx.Writer.Status())
		requestsTotal.WithLabelValues(status).Inc()

		// timer.ObserveDuration()
		ms := time.Since(start).Milliseconds()
		msFloat := float64(ms)
		if ms == 0 {
			msFloat = 0.1
		}
		// requestsDuration.WithLabelValues(req.URL.Path).Observe(msFloat)
		requestsDuration.Observe(msFloat)
	}

	return hf, nil
}
