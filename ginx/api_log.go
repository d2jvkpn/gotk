package ginx

import (
	// "errors"
	"fmt"
	"net/http"
	// "os"
	"time"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/trace_error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
)

type Logger[T any] interface {
	Debug(string, ...T)
	Info(string, ...T)
	Warn(string, ...T)
	Error(string, ...T)
}

func NewAPILog(logger Logger[zap.Field], debug bool,
	metrics ...func(string, float64, codes []string)) (hf gin.HandlerFunc) {
	gomod, _ := gotk.RootModule()
	// debug := logger.Level() == zapcore.DebugLevel

	hf = func(ctx *gin.Context) {
		var (
			e           error
			api         string
			requestId   string
			panicField  map[string]any
			start       time.Time
			err         *trace_error.Error
			fields      []zap.Field
			data        any
			labelValues [2]string
		)

		// concurrentRequests.Inc()
		api = fmt.Sprintf("%s@%s", ctx.Request.Method, ctx.Request.URL.Path)

		fields = make([]zap.Field, 0, 11)
		appendString := func(key, val string) {
			fields = append(fields, zap.String(key, val))
		}

		start = time.Now()
		requestId = uuid.New().String()
		// ctx.Set("RequestId", requestId) // CONTEXT_RequestId
		ctx.Header("x-request-id", requestId)
		// ctx.Header("x-server", server) // HEADER_Server

		// HEADER_Client
		// client := ctx.GetHeader("x-client")

		appendString("ip", ctx.ClientIP())
		appendString("request", api)
		appendString("query", ctx.Request.URL.RawQuery)

		final := func() {
			// ctx.Request.Referer(), ctx.GetHeader("User-Agent")

			tokenId := ctx.GetString("TokenId") // CONTEXT_TokenId
			appendString("tokenId", tokenId)

			accountId := ctx.GetString("AccountId") // CONTEXT_AccountId
			appendString("accountId", accountId)

			latency := float64(time.Since(start).Microseconds()) / 1e3
			fields = append(fields, zap.Float64("latencyMilli", latency))

			status := ctx.Writer.Status()
			fields = append(fields, zap.Int("status", status))

			labelValues[0], labelValues[1] = "OK", ""
			if status != http.StatusOK {
				if err, e = Get[*trace_error.Error](ctx, "Error"); e == nil { // CONTEXT_Error
					fields = append(fields, zap.Any("error", &err))
					labelValues[0], labelValues[1] = err.Code, err.Kind
				}
			}

			// CONTEXT_Data
			if data, e = Get[any](ctx, "Data"); e == nil {
				fields = append(fields, zap.Any("data", &data))
			}

			if panicField != nil {
				fields = append(fields, zap.Any("panic", panicField))
			}

			switch {
			case status < 400:
				logger.Info(requestId, fields...)
				// labelValues[0] = "200"
			case status >= 400 && status < 500:
				logger.Warn(requestId, fields...)
				// labelValues[0] = "400"
			default:
				logger.Error(requestId, fields...)
				// labelValues[0] = "500"
			}

			for i := range errHandlers {
				errHandlers[i](api, latency, err)
			}
		}

		defer func() {
			var (
				panicData any
				panicErr  *trace_error.Error
			)

			if panicData = recover(); panicData == nil {
				return
			}

			// errx.NoTrace(), errx.Skip(5)
			panicErr = trace_error.NewError(
				fmt.Errorf("panic"),
				"internal_error",
				"InternalError",
				trace_error.NoTrace(),
			)
			ctx.JSON(http.StatusInternalServerError, panicErr)

			stacks := gotk.Stack(gomod)
			panicField = map[string]any{"error": panicErr, "data": &panicData, "stacks": stacks}

			final()
		}()

		// ctx.Status(1000)
		ctx.Next()

		/*
			select {
			case <-ctx.Done():
				logger.Named("http_timeout").Warn(requestId)
			default:
			}
		*/

		final()
	}

	return hf
}
