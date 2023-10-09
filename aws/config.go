package aws

import (
	"fmt"

	"github.com/spf13/viper"
)

type S3Config struct {
	AccessKeyId     string `mapstructure:"access_key_id"  json:"accessKeyId,omitempty"`
	SecretAccessKey string `mapstructure:"secrete_access_key" json:"secreteAccessKey,omitempty"`
	Region          string `mapstructure:"region" json:"region,omitempty"`
	Bucket          string `mapstructure:"bucket" json:"bucket,omitempty"`
	bucketUrl       string
	AssumeArn       string `mapstructure:"assume_arn" json:"assumeArn,omitempty"`
	DurationSeconds int64  `mapstructure:"duration_seconds" json:"site,omitempty"`
}

func (config *S3Config) Valid() (err error) {
	if config.AccessKeyId == "" || config.SecretAccessKey == "" {
		return fmt.Errorf("access_key_id or secret_access_key is unset")
	}

	if config.Region == "" || config.Bucket == "" {
		return fmt.Errorf("region_id or bucket is unset")
	}

	// check AssumeArn for Sts

	return nil
}

func NewS3Config(fp, field string) (config S3Config, err error) {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetConfigName("aws_s3_config")

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

	config.bucketUrl = bucketUrl(config.Bucket, config.Region)

	return config, nil
}

func bucketUrl(bucket, region string) string {
	return fmt.Sprintf("https://%s.s3.%s.%s", bucket, region, AWS_Domain)
}
