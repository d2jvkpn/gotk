package ginx

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCtx(t *testing.T) {
	ctx := new(gin.Context)

	SetData(ctx, "a", 1)
	SetData(ctx, "b", 2)

	data, ok := ctx.Get(GIN_Data)
	fmt.Printf("ok: %t, data: %+v\n", ok, data)
}
