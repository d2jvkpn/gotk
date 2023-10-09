package aliyun

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type OssConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id" json:"accessKeyId,omitempty"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"accessKeySecret,omitempty"`
	Region          string `mapstructure:"region" json:"region,omitempty"`
	Bucket          string `mapstructure:"bucket" json:"bucket,omitempty"`
	bucketUrl       string

	// use custom domain instead of https://BUCKET.oss-{REGION_ID}.aliyunc.com
	Site            string `mapstructure:"site" json:"site,omitempty"`
	RoleArn         string `mapstructure:"role_arn" json:"roleArn,omitempty"`                 // sts
	DurationSeconds int    `mapstructure:"duration_seconds" json:"durationSeconds,omitempty"` // sts
}

func NewOssConfig(fp, field string) (config OssConfig, err error) {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetConfigName("aliyun_config")

	//	switch {
	//	case strings.HasSuffix(fp, ".toml"):
	//		conf.SetConfigType("toml")
	//	case strings.HasSuffix(fp, ".yaml"):
	//		conf.SetConfigType("yaml")
	//	default:
	//		return config, fmt.Errorf("unkonw config file, use .yaml or .toml")
	//	}
	conf.SetConfigFile(fp)

	if err = conf.ReadInConfig(); err != nil {
		return config, err
	}

	if err = conf.UnmarshalKey(field, &config); err != nil {
		return config, err
	}

	if err = config.Valid(); err != nil {
		return config, err
	}

	if config.Site != "" {
		config.Site = strings.TrimRight(config.Site, "/")
	}

	config.bucketUrl = bucketUrl(config.Bucket, config.Region)

	return config, nil
}

func bucketUrl(bucket, regionId string) string {
	return fmt.Sprintf("https://%s.oss-%s.%s", bucket, regionId, ALIYUN_Domain)
}

func regionUrl(regionId string) string {
	return fmt.Sprintf("https://oss-%s.%s", regionId, ALIYUN_Domain)
}

func (config *OssConfig) Valid() (err error) {
	if config.AccessKeyId == "" || config.AccessKeySecret == "" {
		return fmt.Errorf("access_key_id or access_key_secret is empty")
	}

	if config.Region == "" || config.Bucket == "" {
		return fmt.Errorf("region or bucket is empty")
	}

	return nil
}

func (config *OssConfig) Url(ps ...string) string {
	if len(ps) == 0 {
		return config.Site
	}

	p := strings.TrimLeft(ps[0], "/")
	return fmt.Sprintf("%s/%s", config.bucketUrl, p)
}

func ValidUrlFilepath(p string) (out string, err error) {
	if out = strings.Trim(p, "/"); out == "" {
		return "", fmt.Errorf("invalid url filepath")
	}

	return out, nil
}
