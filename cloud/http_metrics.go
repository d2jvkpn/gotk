package cloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/ginx"
	"github.com/gin-gonic/gin"
	// "github.com/prometheus/client_golang/prometheus/promauto" // ?
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
	// engine.Use(Cors(vp.GetString("*")))
	// engine.NoRoute(no_route())
	router = &engine.RouterGroup

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	if enableProm {
		router.GET("/prometheus", gin.WrapH(promhttp.Handler()))
	}

	if enableDebug {
		debug := router.Group("/debug", ginx.Cors(vp.GetString("*"), "GET"))

		debug.GET("/meta", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, meta)
		})

		for k, fn := range gotk.PprofHandlerFuncs() {
			debug.GET(fmt.Sprintf("/pprof/%s", k), gin.WrapF(fn))
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
