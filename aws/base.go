package aws

import (
	_ "embed"
)

const (
	AWS_Domain     = "amazonaws.com"
	CODE_NoSuchKey = "NoSuchKey"
)

var (
	//go:embed config.yaml
	configDemo string
)

func init() {
}

func ConfigDemo() string {
	return configDemo
}
