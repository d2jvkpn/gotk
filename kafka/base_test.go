package kafka

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/IBM/sarama"
)

var (
	_TestKafkaVersion string
	_TestAddrs        []string
	_TestTopic        string
	_TestGroupId      string
	_TestConfig       *sarama.Config

	_TestIndex, _TestNum int
	_TestOffset          int64

	_TestFlag *flag.FlagSet
)

// $ go test -run TestHandler -- -index=10
func TestMain(m *testing.M) {
	var (
		addrs string
		err   error
	)

	_TestFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	_TestFlag.StringVar(&addrs, "addrs", "127.0.0.1:29091", "kakfa brokers address seperated by comma")
	_TestFlag.StringVar(&_TestTopic, "topic", "test", "kafka topic")
	_TestFlag.StringVar(&_TestGroupId, "groupId", "default", "kakfa group id")
	_TestFlag.StringVar(&_TestKafkaVersion, "kafkaVersion", "3.4.0", "kakfa version")

	_TestFlag.IntVar(&_TestIndex, "index", 0, "first message index")
	_TestFlag.IntVar(&_TestNum, "num", 10, "number of messages")
	_TestFlag.Int64Var(&_TestOffset, "offset", 0, "offset number")

	_TestFlag.Parse(flag.Args())

	_TestAddrs = strings.Fields(strings.Replace(addrs, ",", " ", -1))

	fmt.Printf(
		"==> TestMain: TestAddrs=%v, TestTopic=%q, TestGroupId=%q, TestKafkaVersion=%q,\n"+
			"    TestIndex=%d, TestNum=%d, TestOffset=%d\n",
		_TestAddrs, _TestTopic, _TestGroupId, _TestKafkaVersion, _TestIndex, _TestNum, _TestOffset,
	)

	// if testNum == 0 {
	// 	fmt.Println("invalid num:", testNum)
	// 	os.Exit(1)
	// }

	_TestConfig = sarama.NewConfig()
	_TestConfig.Version, err = sarama.ParseKafkaVersion(_TestKafkaVersion)
	if err != nil {
		fmt.Printf("!!! Invalid kafka version: %v\n", err)
		os.Exit(1)
	}

	m.Run()
}

func _TestMsgHandle(msg *sarama.ConsumerMessage) {
	tmpl := "<-- MSG: Timestamp=%q, Topic=%q, Partition=%d, Offset=%v; key=%q, value=%q\n"
	// msg.Headers []*RecordHeader
	// msg.BlockTimestamp

	log.Printf(
		tmpl, msg.Timestamp.Format(RFC3339ms), msg.Topic, msg.Partition, msg.Offset,
		msg.Key, msg.Value,
	)
}
