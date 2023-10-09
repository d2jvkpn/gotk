package kafka

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/IBM/sarama"
	"github.com/d2jvkpn/gotk"
	"github.com/spf13/viper"
)

// go test -run TestProducer01 -- -addrs=localhost:29091
func TestProducer01(t *testing.T) {
	var (
		err      error
		producer sarama.AsyncProducer
	)

	if producer, err = sarama.NewAsyncProducer(_TestAddrs, _TestConfig); err != nil {
		t.Fatal(err)
	}

	/*
		#### https://silverback-messaging.net/concepts/broker/kafka/kafka-partitioning.html?tabs=destination-partition-fluent%2Cenricher-fluent%2Cconcurrency-fluent%2Cassignment-fluent
		Kafka can guarantee ordering only inside the same partition and it is therefore important to
		be able to route correlated messages into the same partition. To do so you need to specify a
		key for each message and Kafka will put all messages with the same key in the same partition.
	*/
	for i := _TestIndex; i < _TestIndex+_TestNum; i++ {
		msg := fmt.Sprintf("hello message: %d", i)
		log.Println("--> send msg:", msg)

		pmsg1 := sarama.ProducerMessage{
			Topic: _TestTopic,
			Key:   sarama.StringEncoder("key0001"),
			Value: sarama.ByteEncoder([]byte(msg)),
		}

		producer.Input() <- &pmsg1

		// pmsg2 := pmsg1
		// pmsg2.Key = sarama.StringEncoder("key0002")
		// producer.Input() <- &pmsg2
	}

	if err = producer.Close(); err != nil {
		t.Fatal(err)
	}
}

// go test -run TestProducer02
func TestProducer02(t *testing.T) {
	var (
		err      error
		vp       *viper.Viper
		producer *KafkaProducer
	)

	if vp, err = gotk.LoadYamlConfig("../../configs/local.yaml", "TestConfig"); err != nil {
		t.Fatal(err)
	}

	if producer, err = NewKafkaProducer(vp, "kafka"); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("==> %#v\n", producer)

	msg := producer.SendMsg(context.TODO(), []byte("hello world"))
	fmt.Printf("~~~ msg: %#v\n", msg)

	if err = producer.Close(); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("~~~ msg: %#v\n", msg)
}
