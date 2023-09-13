package ginx

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCtx(t *testing.T) {
	ctx := new(gin.Context)

	SetDataField(ctx, "a", 1)
	SetDataField(ctx, "b", 2)

	data, ok := ctx.Get(GIN_Data)
	fmt.Printf("ok: %t, data: %+v\n", ok, data)
}
