package ginx

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/json"
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

func JSONStatic(data any) gin.HandlerFunc {
	bts, _ := json.Marshal(data)

	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Writer.Write(bts)
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
	if listDir {
		return func(rg *gin.RouterGroup) {
			rg.StaticFS(httpDir, http.Dir(local))
		}
	} else {
		return func(rg *gin.RouterGroup) {
			rg.Static(httpDir, local)
		}
	}
}

/*
go:embed static
efs embed.FS

p = "static"
secs = 3600
*/
func ServeStaticFS(rg *gin.RouterGroup, efs embed.FS, p string, secs int) (err error) {
	var fsys fs.FS

	if fsys, err = fs.Sub(efs, p); err != nil {
		return err
	}

	static := rg.Group("/"+p, CacheControl(secs))
	static.StaticFS("/", http.FS(fsys))

	return nil
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

func ResponseFile(ctx *gin.Context, buf *bytes.Buffer, filename, typ string) {
	var contextType string

	switch typ {
	case "xls", "xlsx":
		contextType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "doc", "docx":
		contextType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "pdf":
		contextType = "application/pdf"
	default:
		contextType = "application/octet-stream"
	}

	// ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)

	ctx.Data(http.StatusOK, contextType, buf.Bytes())
}
