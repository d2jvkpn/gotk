package ginx

import (
	"fmt"
	"testing"
	"time"
	// "errors"

	"github.com/spf13/viper"
)

func TestJwtHMAC(t *testing.T) {
	var (
		err     error
		bearar  string
		kind    string
		jwtData *JwtData
		jwt     *JwtHMAC
	)

	vp := viper.New()
	vp.Set("key", 123456)
	vp.Set("duration", time.Minute)
	vp.Set("method", 256)

	if jwt, err = NewJwtHMAC(vp, "app"); err != nil {
		t.Fatal(err)
	}

	bearar, err = jwt.Sign(&JwtData{
		ID:      "xx01",
		Subject: "acc01",
		Data:    map[string]string{"role": "admin"},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("==> bearar: %v\n", bearar)
	time.Sleep(11 * time.Second)

	if jwtData, kind, err = jwt.Auth(bearar); err == nil {
		t.Fatal("expect: token expired")
	}
	fmt.Printf("==> %v, %s, %+#v\n", jwtData, kind, err) // token has invalid claims: token is expired

	if jwtData, err = jwt.ParsePayload(bearar); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("==> jwt data: %#+v\n", jwtData)
}
