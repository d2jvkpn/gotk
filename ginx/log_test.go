package ginx

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCtx(t *testing.T) {
	ctx := new(gin.Context)

	SetData(ctx, map[string]any{"a": 1})
	SetData(ctx, map[string]any{"b": 2})

	data, ok := ctx.Get(GIN_Data)
	fmt.Printf("ok: %t, data: %+v\n", ok, data)

	call(ctx)
	data, ok = ctx.Get(GIN_Data)
	fmt.Printf("ok: %t, data: %+v\n", ok, data) // has no ans: 42
}

func call(ctx context.Context) {
	ctx = context.WithValue(ctx, "ans", 42)
}
