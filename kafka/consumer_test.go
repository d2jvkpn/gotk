package kafka

import (
	"fmt"
	"log"
	"testing"

	"github.com/IBM/sarama"
)

// go test -run TestConsumer -- -addrs=localhost:29091
func TestConsumer(t *testing.T) {
	var (
		count      int
		topics     []string
		partitions []int32
		err        error

		consumer  sarama.Consumer
		pconsumer sarama.PartitionConsumer
	)

	if consumer, err = sarama.NewConsumer(_TestAddrs, _TestConfig); err != nil {
		t.Fatal(err)
	}

	if topics, err = consumer.Topics(); err != nil {
		t.Fatal(err)
	}

	if len(topics) == 0 {
		t.Fatal("no topics")
	}
	fmt.Println("~~~ Avaialble topics:", topics)

	if partitions, err = consumer.Partitions(_TestTopic); err != nil {
		t.Fatal(err)
	}
	if len(partitions) == 0 {
		t.Fatalf("topic %s has no partitions\n", _TestTopic)
	}
	log.Printf("~~~ topic %s partitions: %v\n", _TestTopic, partitions)

	// topic string, partition int32, offset int64
	pconsumer, err = consumer.ConsumePartition(_TestTopic, partitions[0], _TestOffset)
	if err != nil {
		t.Fatal(err)
	}

	// go TestProducer(t)

	mc := pconsumer.Messages() // *sarama.ConsumerMessage

LOOP:
	for i := 0; i < _TestNum; i++ {
		select {
		case msg, ok := <-mc:
			if !ok {
				log.Println("!!! consumer is closed")
				break LOOP
			}
			_TestMsgHandle(msg)
			count += 1
		}
	}

	// pconsumer methods: Close, Pause, Resume
	log.Printf("~~~ Message consumed: %d\n", count)
	if err = consumer.Close(); err != nil {
		t.Fatal(err)
	}
}
