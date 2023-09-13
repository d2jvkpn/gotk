package ginx

import (
	"github.com/gin-gonic/gin"
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

func SetIdentity(ctx *gin.Context, kv map[string]any) {
	identity, e := Get[map[string]any](ctx, GIN_Identity)
	if e != nil {
		identity = make(map[string]any, 1)
		ctx.Set(GIN_Identity, identity)
	}

	for k := range kv {
		identity[k] = kv[k]
	}
}

func SetData(ctx *gin.Context, kv map[string]any) {
	data, e := Get[map[string]any](ctx, GIN_Data)
	if e != nil {
		data = make(map[string]any, 1)
		ctx.Set(GIN_Data, data)
	}

	for k := range kv {
		data[k] = kv[k]
	}
}
