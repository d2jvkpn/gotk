package gotk

import (
	"context"
	"net/http"
)

type HttpContext struct {
	key string
	*http.Request
}

func NewHttpContext(key string, r *http.Request) *HttpContext {
	ctx := context.WithValue(r.Context(), key, make(map[string]any, 0))
	*r = *(r.WithContext(ctx))

	return &HttpContext{
		key:     key,
		Request: r,
	}
}

func (ctx *HttpContext) Set(key string, value any) bool {
	data, _ := ctx.Request.Context().Value(key).(map[string]any)

	_, ok := data[key]
	data[key] = value
	return ok
}

func (ctx *HttpContext) Get(key string) (any, bool) {
	data, _ := ctx.Request.Context().Value(ctx.key).(map[string]any)
	value, ok := data[key]
	return value, ok
}

func (ctx *HttpContext) GetData() map[string]any {
	data, _ := ctx.Request.Context().Value(ctx.key).(map[string]any)
	return data
}

func MiddlewareHandler(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Content-Type, Authorization, X-CSRF-Token
		w.Header().Set(
			"Access-Control-Expose-Headers",
			"Access-Control-Allow-Origin, Access-Control-Allow-Headers, "+
				"Content-Type, Content-Length",
		)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, HEAD")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MiddlewareFunc(origin string, hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Content-Type, Authorization, X-CSRF-Token
		w.Header().Set(
			"Access-Control-Expose-Headers",
			"Access-Control-Allow-Origin, Access-Control-Allow-Headers, "+
				"Content-Type, Content-Length",
		)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, HEAD")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		hf(w, r)
	}
}
