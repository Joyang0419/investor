package kafka

import (
	"fmt"
	"time"

	"tools/errorx"
	"tools/logger"
	"tools/reflectx"
	"tools/slicex"

	"github.com/segmentio/kafka-go"
)

const defaultTimeout = 3 * time.Second

type Config struct {
	Host string
	Port int
}

func NewKafkaConn(config Config) *kafka.Conn {
	kafkaConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		logger.Fatal("[NewKafkaConn]kafka.Dial err: %v", err)
	}

	return kafkaConn
}

func NewKafkaSyncProducer(conn *kafka.Conn, topic string, requiredAck kafka.RequiredAcks, writeTimeout ...time.Duration) *kafka.Writer {
	if reflectx.IsNil(conn) {
		logger.Fatal("[NewKafkaSyncProducer]conn is nil")
	}
	if slicex.IsElementNotInSlice([]kafka.RequiredAcks{kafka.RequireNone, kafka.RequireOne, kafka.RequireAll}, requiredAck) {
		logger.Fatal("[NewKafkaSyncProducer]requiredAck is invalid")
	}

	wt := defaultTimeout
	if len(writeTimeout) > 0 {
		wt = writeTimeout[0]
	}

	brokers, err := conn.Brokers()
	if errorx.IsErrorExist(err) {
		logger.Fatal("[NewKafkaSyncProducer]conn.Brokers err: %v", err)
	}

	addrs := make([]string, 0, len(brokers))
	for idx := range brokers {
		addrs = append(addrs, fmt.Sprintf("%s:%d", brokers[idx].Host, brokers[idx].Port))
	}

	return &kafka.Writer{
		Addr:         kafka.TCP(addrs...),
		Topic:        topic,
		RequiredAcks: requiredAck,
		Async:        false,
		WriteTimeout: wt,
	}
}

//func NewKafkaConsumerGroup(config Config, groupID string, topics []string) *kafka.ConsumerGroup {
//	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
//		ID:      groupID,
//		Brokers: config.Brokers,
//		Topics:  topics,
//	})
//	if err != nil {
//		logger.Fatal("[NewKafkaConsumerGroup]kafka.NewConsumerGroup err: %v", err)
//	}
//
//	return group
//}
//
//// TODO ERIC 要用的，讓他自己調整, 我先把Producer搞好一版
//func NewKafkaReader(config Config, topic string, partition int) *kafka.Reader {
//	r := kafka.NewReader(kafka.ReaderConfig{
//		Brokers:   config.Brokers,
//		Topic:     topic,
//		Partition: partition,
//		MinBytes:  10e3, // 10KB
//		MaxBytes:  10e6, // 10MB
//	})
//	return r
//}
