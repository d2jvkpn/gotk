package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type KafkaProducer struct {
	config   Config
	producer sarama.AsyncProducer
}

func NewKafkaProducer(vp *viper.Viper, field string) (producer *KafkaProducer, err error) {
	var (
		config *Config
		scfg   *sarama.Config
	)

	if config, scfg, err = NewConfigFromViper(vp, field); err != nil {
		return nil, err
	}

	if config.Topic == "" || config.Key == "" {
		return nil, fmt.Errorf("neither the topic nor the key is set")
	}

	producer = &KafkaProducer{config: *config}

	producer.producer, err = sarama.NewAsyncProducer(producer.config.Addrs, scfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (producer *KafkaProducer) SendMsg(ctx context.Context, bts []byte) (msg *sarama.ProducerMessage) {
	msg = &sarama.ProducerMessage{
		Topic: producer.config.Topic,
		Key:   sarama.StringEncoder(producer.config.Key),
		Value: sarama.ByteEncoder(bts),
	}

	select {
	case <-ctx.Done():
	case producer.producer.Input() <- msg:
	}
	return msg
}

func (producer *KafkaProducer) Close() (err error) {
	return producer.producer.Close()
}
