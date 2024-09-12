package gotk

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestViper(t *testing.T) {
	var (
		err     error
		vp, sub *viper.Viper
	)

	if vp, err = LoadYamlConfig("viper_test.yaml", "test"); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("==> 1. vp: %v\n", vp)
	fmt.Printf("==> 2. hello.world=%d\n", vp.GetInt("hello.world"))

	fmt.Printf("==> 3. prometheus.enabled=%t\n", vp.GetBool("prometheus.enabled"))

	//sub = vp.Sub("prometheus")
	//fmt.Printf("==> sub.enabled=%t\n", sub.GetBool("enabled")) // invalid memory address or nil pointer dereference

	vp.SetDefault("prometheus", map[string]any{"hello": 1024})
	sub = vp.Sub("prometheus")
	fmt.Printf("==> 4.  sub.enabled=%t\n", sub.GetBool("enabled")) // sub.enabled=false

	vp.SetDefault("hello", map[string]any{})
	fmt.Printf("==> 5. hello.world=%d\n", vp.GetInt("hello.world")) // 42
}
