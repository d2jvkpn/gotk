package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type Config struct {
	Addrs   []string `mapstructure:"addrs"`
	Version string   `mapstructure:"version"` // 3.4.0
	Topic   string   `mapstructure:"topic"`

	// consumer
	GroupId string `mapstructure:"group_id"` // default

	// producer
	Key string `mapstructure:"key"`
}

func NewConfigFromViper(vp *viper.Viper, field string) (
	config *Config, scfg *sarama.Config, err error) {

	config = new(Config)

	if err = vp.UnmarshalKey(field, config); err != nil {
		return nil, nil, err
	}

	if len(config.Addrs) == 0 || config.Version == "" {
		return nil, nil, fmt.Errorf("invalid addrs or version")
	}

	if config.Topic == "" {
		return nil, nil, fmt.Errorf("invalid topic")
	}

	scfg = sarama.NewConfig()
	if scfg.Version, err = sarama.ParseKafkaVersion(config.Version); err != nil {
		return nil, nil, err
	}

	return config, scfg, nil
}
