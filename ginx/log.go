package ginx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/d2jvkpn/gotk"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	GIN_RequestId = "GIN_RequestId"
	GIN_Error     = "GIN_Error"
	GIN_Identity  = "GIN_Identity"
	GIN_Data      = "GIN_Data"
)

func SetRequestId(ctx *gin.Context, requestId string) {
	ctx.Set(GIN_RequestId, requestId)
}

func SetError(ctx *gin.Context, err any) {
	ctx.Set(GIN_Error, err)
}

func SetIdentity(ctx *gin.Context, key string, val any) {
	identity, e := Get[map[string]any](ctx, GIN_Identity)
	if e != nil {
		identity = make(map[string]any, 1)
	}

	identity[key] = val

	if e != nil {
		ctx.Set(GIN_Identity, identity)
	}
}

func SetData(ctx *gin.Context, key string, val any) {
	data, e := Get[map[string]any](ctx, GIN_Data)
	if e != nil {
		data = make(map[string]any, 1)
	}

	data[key] = val
	if e != nil {
		ctx.Set(GIN_Data, data)
	}
}

// skip = 5
func APILog(lgg *zap.Logger, name string, skip int) gin.HandlerFunc {
	mod, _ := gotk.RootModule()
	logger := lgg.Named(name)

	return func(ctx *gin.Context) {
		var (
			start time.Time
			msg   string
		)

		start = time.Now()
		msg = fmt.Sprintf("%s@%s", ctx.Request.Method, ctx.Request.URL.Path)

		defer func() {
			var recoverData any

			if recoverData = recover(); recoverData == nil {
				return
			}
			// fmt.Println("!!! 3")

			// TODO: alerting
			stacks := gotk.Stack(mod)
			SetError(ctx, map[string]any{
				"recover": recoverData,
				"stacks":  stacks,
			})

			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"request_id": ctx.GetString(GIN_RequestId), "code": 1, "msg": "panic"},
			)
			apiLogEnd(ctx, start, msg, logger)
		}()

		ctx.Next()
		apiLogEnd(ctx, start, msg, logger)
	}
}

func apiLogEnd(ctx *gin.Context, start time.Time, msg string, logger *zap.Logger) {
	var (
		statusCode int
		fields     []zap.Field
	)

	fields = make([]zap.Field, 0, 8)
	push := func(key, field string) {
		fields = append(fields, zap.String(key, field))
	}

	push("request_id", ctx.GetString(GIN_RequestId))
	push("ip", ctx.ClientIP())
	push("query", ctx.Request.URL.RawQuery)

	statusCode = ctx.Writer.Status()
	fields = append(fields, zap.Int("status_code", statusCode))

	latency := fmt.Sprintf("%fms", float64(time.Since(start).Microseconds())/1e3)
	push("latency", latency)

	if identity, e := Get[map[string]any](ctx, GIN_Identity); e == nil {
		fields = append(fields, zap.Any("identity", identity))
	}

	if data, e := Get[map[string]any](ctx, GIN_Data); e == nil {
		fields = append(fields, zap.Any("data", data))
	}

	// ki, _ = Get[*api_key.KeyInvocation](ctx, KEY_AKI)
	if err, ok := ctx.Get(GIN_Error); ok {
		fields = append(fields, zap.Any("error", err))
	}

	switch {
	case statusCode < 400:
		logger.Info(msg, fields...)
	case statusCode >= 400 && statusCode < 500:
		logger.Warn(msg, fields...)
	default:
		logger.Error(msg, fields...)
	}
}
