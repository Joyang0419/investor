package kafka

import (
	"context"
	"fmt"
	"time"

	"tools/logger"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Network       string
	Host          string
	Port          int
	WriteDeadline time.Duration
	ReadDeadline  time.Duration
	Brokers       []string // for consumer groups only
}

func NewKafkaConn(config Config, topic string, partition int) *kafka.Conn {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := kafka.DialLeader(context.Background(), config.Network, address, topic, partition)
	if err != nil {
		logger.Fatal("[NewKafkaConn]kafka.DialLeader err: %v", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(config.WriteDeadline))
	if err != nil {
		logger.Fatal("[NewKafkaConn]conn.SetWriteDeadline err: %v", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(config.ReadDeadline))
	if err != nil {
		logger.Fatal("[NewKafkaConn]conn.SetReadDeadline err: %v", err)
	}

	return conn
}

func NewKafkaConsumerGroup(config Config, groupID string, topics []string) *kafka.ConsumerGroup {
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      groupID,
		Brokers: config.Brokers,
		Topics:  topics,
	})
	if err != nil {
		logger.Fatal("[NewKafkaConsumerGroup]kafka.NewConsumerGroup err: %v", err)
	}

	return group
}

func NewKafkaReader(config Config, topic string, partition int) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   config.Brokers,
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	return r
}

func NewKafkaWriter(config Config, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: config.Brokers,
		Topic:   topic,
	})
}

// TODO
/*
import (
    "fmt"

    "github.com/IBM/sarama"

    "tools/errorx"
    "tools/logger"
    "tools/slicex"
)


// Config 結構用於存儲 Kafka 配置
type Config struct {
    Brokers []string
}

// KafkaConnection 包含 Kafka 配置，並允許動態創建 Producer
type KafkaConnection struct {
    cfg Config
}

// SetupKafkaClient 函數用於設置 Kafka 連線並確認 Kafka Broker 可用
func SetupKafkaClient(cfg Config) sarama.Client {
    // 檢查 Kafka Broker 是否可用
    c, err := sarama.NewClient(cfg.Brokers, sarama.NewConfig())
    if err != nil {
        panic(fmt.Sprintf("[tools][infra_conn][SetupKafkaConnection]: %v", err))
    }
    defer func(client sarama.Client) {
        if errClose := client.Close(); errorx.IsErrorExist(errClose) {
            logger.Error("[tools][infra_conn][SetupKafkaConnection]c.Close err: %v", errClose)
        }
    }(c)
    // 檢查每個 Broker 的可用性
    for _, broker := range c.Brokers() {
        connected, errConnected := broker.Connected()
        if errorx.IsErrorExist(errConnected) {
            panic(fmt.Sprintf("[tools][infra_conn][SetupKafkaConnection]broker.Connected err: %v", errConnected))
        }
        if !connected {
            panic(fmt.Sprintf("[tools][infra_conn][SetupKafkaConnection]broker not connected: %v", broker.Addr()))
        }
    }

    return c
}

func SetupProducer(client sarama.Client, IsSync ...bool) sarama.SyncProducer {
    addrs := make([]string, 0, len(client.Brokers()))
    for _, broker := range client.Brokers() {
        addrs = append(addrs, broker.Addr())
    }

    cfg := sarama.NewConfig()
    if slicex.IsEmpty(IsSync) {
        cfg.Producer.Return.Successes = IsSync[0]
    }

    p, err := sarama.NewSyncProducer(addrs, cfg)
    if errorx.IsErrorExist(err) {
        panic(fmt.Sprintf("[tools][infra_conn][SetupProducer]NewSyncProducer err: %v", err))
    }

    return p
}

var validOffsetPolicy = []int64{sarama.OffsetOldest, sarama.OffsetNewest}

func SetupConsumerGroup(client sarama.Client, groupID string, offsetPolicy ...int64) sarama.ConsumerGroup {
    addrs := make([]string, 0, len(client.Brokers()))
    for _, broker := range client.Brokers() {
        addrs = append(addrs, broker.Addr())
    }

    cfg := sarama.NewConfig()
    if slicex.IsEmpty(offsetPolicy) {
        if slicex.IsElementNotInSlice(validOffsetPolicy, offsetPolicy[0]) {
            panic(fmt.Sprintf("[tools][infra_conn][SetupConsumerGroup]invalid offset policy: %v", offsetPolicy[0]))
        }
        cfg.Consumer.Offsets.Initial = offsetPolicy[0]
    }

    cg, err := sarama.NewConsumerGroup(addrs, groupID, cfg)
    if errorx.IsErrorExist(err) {
        panic(fmt.Sprintf("[tools][infra_conn][SetupConsumerGroup]NewConsumerGroup err: %v", err))
    }

    return cg
}

func SetupConsumer(client sarama.Client, offsetPolicy ...int64) sarama.Consumer {
    addrs := make([]string, 0, len(client.Brokers()))
    for _, broker := range client.Brokers() {
        addrs = append(addrs, broker.Addr())
    }

    cfg := sarama.NewConfig()
    if slicex.IsEmpty(offsetPolicy) {
        if slicex.IsElementNotInSlice(validOffsetPolicy, offsetPolicy[0]) {
            panic(fmt.Sprintf("[tools][infra_conn][SetupConsumerGroup]invalid offset policy: %v", offsetPolicy[0]))
        }
        cfg.Consumer.Offsets.Initial = offsetPolicy[0]
    }

    c, err := sarama.NewConsumer(addrs, cfg)
    if errorx.IsErrorExist(err) {
        panic(fmt.Sprintf("[tools][infra_conn][SetupConsumer]NewConsumer err: %v", err))
    }

    return c
}

*/
