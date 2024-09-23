package vocechat

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var (
	_TestFlag   *flag.FlagSet
	_TestCtx    context.Context = context.Background()
	_TestConfig *viper.Viper
)

// go test --run TestName -- --config=configs/local.yaml
func TestMain(m *testing.M) {
	var (
		// release bool
		config string
		err    error
	)

	_TestFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	// _TestFlag.BoolVar(&release, "release", false, "run in release mode")
	_TestFlag.StringVar(&config, "config", "configs/local.yaml", "config filepath")

	_TestFlag.Parse(flag.Args())
	fmt.Printf("==> load config: %q\n", config)

	defer func() {
		if err != nil {
			fmt.Printf("!!! TestMain: %v\n", err)
			os.Exit(1)
		}
	}()

	_TestConfig = viper.New()
	_TestConfig.SetConfigType("yaml")
	_TestConfig.SetConfigFile(config)

	if err = _TestConfig.ReadInConfig(); err != nil {
		return
	}

	m.Run()
}
