package ginx

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// errors: value is unset, type not match
func Get[T any](ctx *gin.Context, key string) (item T, err error) {
	var (
		ok    bool
		value any
	)

	if value, ok = ctx.Get(key); !ok {
		// return item, fmt.Errorf("value is unset: %s", key)
		return item, fmt.Errorf("value is unset")
	}

	if item, ok = value.(T); !ok {
		// return item, fmt.Errorf("type of value doesn't match: %s", key)
		return item, fmt.Errorf("type not match")
	}

	return item, nil
}

func GetDate(ctx *gin.Context) (*time.Time, error) {
	var (
		header string
		ans    time.Time
		err    error
	)

	header = ctx.GetHeader("Date")

	if ans, err = time.Parse(time.RFC1123, header); err != nil {
		return nil, err
	}

	// ?? ans.Location == *time.Local

	return &ans, nil
}

func GetAPI(ctx *gin.Context) string {
	return fmt.Sprintf("%s@%s", ctx.Request.Method, ctx.Request.URL.Path)
}
