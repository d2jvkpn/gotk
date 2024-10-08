package ginx

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors0(origin string, methods string) gin.HandlerFunc {
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

/*
HTTP/1.1 204 No Content
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Origin,Content-Type,Authorization
Access-Control-Allow-Methods: GET,POST,OPTIONS,HEAD
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Max-Age: 3600
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
Date: Thu, 08 Aug 2024 00:18:50 GMT
*/
func Cors(origins []string, maxAges ...time.Duration) gin.HandlerFunc {
	maxAge := 12 * time.Hour
	if len(maxAges) > 0 {
		maxAge = maxAges[0]
	}

	return cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "OPTIONS", "HEAD"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
		},
		AllowCredentials: true,
		// AllowOriginFunc:  func(origin string) bool { return origin == "https://github.com" },
		MaxAge: maxAge,
	})
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
