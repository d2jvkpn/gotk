package ginx

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

func Healthz(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
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
func ServeStaticFile(bts []byte, name string) gin.HandlerFunc {
	ct := ""
	if len(bts) > 512 {
		ct = http.DetectContentType(bts[:512])
	} else {
		ct = http.DetectContentType(bts)
	}
	// ext := mini.ExtensionsByType(ct)

	return func(ctx *gin.Context) {
		reader := bytes.NewReader(bts)
		ctx.Header("Content-Type", ct)
		http.ServeContent(ctx.Writer, ctx.Request, name, time.Now(), reader)
	}
}

func WsUpgrade(ctx *gin.Context) {
	if ctx.GetHeader("Upgrade") != "websocket" && ctx.GetHeader("Connection") != "Upgrade" {
		// fmt.Printf("~~~~ Headers: %v\n", ctx.Request.Header)
		ctx.String(http.StatusUpgradeRequired, "Upgrade Required")
		ctx.Abort()
		return
	}

	ctx.Next()
}
