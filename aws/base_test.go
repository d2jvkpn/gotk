package aws

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/d2jvkpn/gotk"
)

var (
	testContext  context.Context = context.Background()
	testConfig   *S3Config
	testFlag     *flag.FlagSet
	testS3Client *S3Client
)

// default config: wk_config/aws_test.yaml
func TestMain(m *testing.M) {
	var (
		configFile string
		field      string
		err        error
	)

	testFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	testFlag.StringVar(&configFile, "config", "configs/aws_test.yaml", "config filepath")
	testFlag.StringVar(&field, "field", "aws_s3", "aws s3 field in config")
	testFlag.Parse(flag.Args())

	if configFile, err = gotk.RootFile(configFile); err != nil {
		panic(err)
	}
	fmt.Printf("~~~ config %s::%s\n", configFile, field)

	if testS3Client, err = NewS3Client(configFile, field); err != nil {
		panic(err)
	}

	m.Run()
}
