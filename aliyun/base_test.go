package aliyun

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/d2jvkpn/gotk"
)

var (
	testOssConfig *OssConfig
	testFlag      *flag.FlagSet
	testOssClient *OssClient
)

// default config: wk_config/aliyun_test.yaml
func TestMain(m *testing.M) {
	var (
		configFile string
		field      string
		err        error
	)

	testFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	testFlag.StringVar(&configFile, "config", "configs/aliyun_test.yaml", "config filepath")
	testFlag.StringVar(&field, "oss", "aliyun_oss", "aliyun oss field in config")

	testFlag.Parse(flag.Args())
	fmt.Printf("~~~ config %s::[%s]\n", configFile, field)

	defer func() {
		if err != nil {
			fmt.Printf("!!! TestMain: %v\n", err)
			os.Exit(1)
		}
	}()

	if configFile, err = gotk.RootFile(configFile); err != nil {
		return
	}

	//	if testConfig, err = NewConfig(configFile, field); err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}

	if testOssClient, err = NewOssClient(configFile, field); err != nil {
		return
	}

	m.Run()
}

func TestConfig(t *testing.T) {
	var (
		config OssConfig
		err    error
	)

	if config, err = NewOssConfig("config.yaml", "aliyun_oss"); err != nil {
		t.Fatal(err)
	}

	if config, err = NewOssConfig("config.yaml", "aliyun_sts"); err != nil {
		t.Fatal(err)
	}

	fmt.Println(config)
}

func TestConfigDemo(t *testing.T) {
	fmt.Printf(">>> TestConfigDemo:\n%s\n", ConfigDemo())
}
