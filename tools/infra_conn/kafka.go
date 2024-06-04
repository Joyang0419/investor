package infra_conn

import (
	"fmt"

	"github.com/IBM/sarama"

	"tools/errorx"
	"tools/logger"
	"tools/slicex"
)

// KafkaCfg 結構用於存儲 Kafka 配置
type KafkaCfg struct {
	Brokers []string
}

// KafkaConnection 包含 Kafka 配置，並允許動態創建 Producer
type KafkaConnection struct {
	cfg KafkaCfg
}

// SetupKafkaClient 函數用於設置 Kafka 連線並確認 Kafka Broker 可用
func SetupKafkaClient(cfg KafkaCfg) sarama.Client {
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
