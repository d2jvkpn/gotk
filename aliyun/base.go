package aliyun

import (
	_ "embed"
)

const (
	ALIYUN_Domain  = "aliyuncs.com"
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
