package impls

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func IndexStaticFiles(router *gin.RouterGroup, d string) (err error) {
	var files []fs.FileInfo

	if files, err = ioutil.ReadDir(d); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		router.StaticFile("/"+f.Name(), filepath.Join(d, f.Name()))
	}

	return nil
}

func Cors(origin string) gin.HandlerFunc {
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
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// handle key: no_token, invalid_token, incorrect_token, User:XXXX
func BasicAuth(username, password string, handle func(*gin.Context, string)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			bts []byte
			err error
		)

		key := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(key, "Basic ") {
			ctx.Header("Www-Authenticate", `Basic realm="login required"`)
			handle(ctx, "no_token")
			ctx.Abort()
			return
		}

		if bts, err = base64.StdEncoding.DecodeString(key[6:]); err != nil {
			handle(ctx, "invalid_token")
			ctx.Abort()
			return
		}

		if u, p, found := strings.Cut(string(bts), ":"); !found {
			handle(ctx, "invalid_token")
			ctx.Abort()
			return
		} else if u != username || p != password {
			handle(ctx, "incorrect_token")
			ctx.Abort()
			return
		}

		handle(ctx, fmt.Sprintf("User:%s", username))
		ctx.Next()
	}
}

// handle key: no_token, invalid_token, incorrect_token, User:XXXX
func BasicBcrypt(username, password string, handle func(*gin.Context, ...string)) gin.HandlerFunc {
	passwordBts := []byte(password)

	return func(ctx *gin.Context) {
		var (
			key []byte
			err error
		)

		key = []byte(ctx.GetHeader("Authorization"))

		if !bytes.HasPrefix(key, []byte("Basic ")) {
			ctx.Header("Www-Authenticate", `Basic realm="login required"`)
			handle(ctx, "no_token")
			ctx.Abort()
			return
		}
		key = key[6:]

		if key, err = base64.StdEncoding.DecodeString(string(key)); err != nil {
			handle(ctx, "no_token")
			ctx.Abort()
			return
		}

		u, p, found := bytes.Cut(key, []byte{':'})
		if !found {
			handle(ctx, "invalid_token")
			ctx.Abort()
			return
		}
		if string(u) != username {
			_ = bcrypt.CompareHashAndPassword(passwordBts, p)
			handle(ctx, "incorrect_token")
			ctx.Abort()
			return
		}

		if err = bcrypt.CompareHashAndPassword(passwordBts, p); err != nil {
			handle(ctx, "incorrect_token")
			ctx.Abort()
			return
		}

		handle(ctx, "user", username)
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

func Healthz(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func ServeStaticDir(httpDir, local string, listDir bool) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		if listDir {
			rg.StaticFS(httpDir, http.Dir(local))
		} else {
			rg.Static(httpDir, local)
		}
		return
	}
}

// name: filename, e.g. favicon.ico
// ct: Content-Type, e.g. image/x-icon
func ServeStaticFile(bts []byte, name, ct string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reader := bytes.NewReader(bts)
		ctx.Header("Content-Type", ct)
		http.ServeContent(ctx.Writer, ctx.Request, name, time.Now(), reader)
	}
}
