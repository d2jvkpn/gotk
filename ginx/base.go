package ginx

import (
	"github.com/gin-gonic/gin"
)

const (
	GIN_RequestId = "GIN_RequestId"
	GIN_Error     = "GIN_Error"
	GIN_AccountId = "GIN_AccountId" // string
	GIN_Data      = "GIN_Data"      // map[string]any

)

func SetRequestId(ctx *gin.Context, requestId string) {
	ctx.Set(GIN_RequestId, requestId)
}

func SetError(ctx *gin.Context, err any) {
	ctx.Set(GIN_Error, err)
}

func SetData(ctx *gin.Context, kvs map[string]any) {
	data, e := Get[map[string]any](ctx, GIN_Data)
	if e != nil {
		data = make(map[string]any, 1)
		ctx.Set(GIN_Data, data)
	}

	for k := range kvs {
		data[k] = kvs[k]
	}
}

func SetKV(ctx *gin.Context, key string, value any) {
	var (
		e    error
		data map[string]any
	)

	if data, e = Get[map[string]any](ctx, GIN_Data); e != nil {
		data = make(map[string]any, 1)
		ctx.Set(GIN_Data, data)
	}

	data[key] = value
}
