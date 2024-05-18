package ginx

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors(origin string, methods string) gin.HandlerFunc {
	// methods: GET, POST, PUT, OPTIONS, HEAD
	if origin == "" {
		origin = "*"
	}

	allowHeaders := strings.Join([]string{"Content-Type", "Authorization"}, ", ")

	exposeHeaders := strings.Join([]string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Content-Type",
		"Content-Length",
	}, ", ")

	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", origin)

		ctx.Header("Access-Control-Allow-Headers", allowHeaders)
		// Content-Type, Authorization, X-CSRF-Token
		ctx.Header("Access-Control-Expose-Headers", exposeHeaders)
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", methods)

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func CacheControl(seconds int) gin.HandlerFunc {
	cc := fmt.Sprintf("public, max-age=%d", seconds)
	// strconv.FormatInt(time.Now().UnixMilli(), 10)
	etag := fmt.Sprintf(`"%d"`, time.Now().UnixMilli()) // must be a quoted string

	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}

		ctx.Header("Cache-Control", cc)
		// browser send If-None-Match: etag, if unchanged, response 304
		ctx.Header("ETag", etag)
		ctx.Next()
	}
}
