package kafka

import (
	"context"
	// "fmt"
	"log"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// go test -run TestHandler -o TestHandler.out
// docker exec -it broker-1 bash
// ./TestHandler.out -- --help
// ./TestHandler.out -- -addrs=kafka-1:9092 -num=20
//
// go test -run TestHandler -- -addrs=localhost:29091
func TestHandler(t *testing.T) {
	var (
		err error
		ctx context.Context

		group   sarama.ConsumerGroup
		handler *Handler // impls sarama.ConsumerGroupHandler
	)

	_TestConfig.Consumer.Return.Errors = true
	_TestConfig.Consumer.Offsets.Initial = sarama.OffsetOldest // sarama.OffsetNewest
	group, err = sarama.NewConsumerGroup(_TestAddrs, _TestGroupId, _TestConfig)
	if err != nil {
		t.Fatal(err)
	}

	// go TestProducer(t)

	ctx = context.Background()
	handler = NewHandler(ctx, group, []string{_TestTopic})
	logger, _ := zap.NewProduction()
	err = handler.WithHandle(_TestMsgHandle).WithLogger(logger).Ok()
	if err != nil {
		t.Fatal(err)
	}

	handler.Consume()

	time.Sleep(15 * time.Second)
	log.Println("<<< Exit")

	if err = handler.Close(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}
