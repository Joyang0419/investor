package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"tools/errorx"
	"tools/logger"
	"tools/reflectx"
)

type OptionConsumer func(*kafka.ReaderConfig)

func WithGroupID(groupID string) OptionConsumer {
	return func(c *kafka.ReaderConfig) {
		c.GroupID = groupID
	}
}

func WithMinBytes(minBytes int) OptionConsumer {
	return func(c *kafka.ReaderConfig) {
		c.MinBytes = minBytes
	}
}

func WithMaxBytes(maxBytes int) OptionConsumer {
	return func(c *kafka.ReaderConfig) {
		c.MaxBytes = maxBytes
	}
}

func NewKafkaConsumer(conn *kafka.Conn, topic string, opts ...OptionConsumer) *kafka.Reader {
	if reflectx.IsNil(conn) {
		logger.Fatal("[NewKafkaConsumer]conn is nil")
	}

	brokers, err := conn.Brokers()
	if errorx.IsErrorExist(err) {
		logger.Fatal("[NewKafkaConsumer]conn.Brokers err: %v", err)
	}

	addrs := make([]string, 0, len(brokers))
	for idx := range brokers {
		addrs = append(addrs, fmt.Sprintf("%s:%d", brokers[idx].Host, brokers[idx].Port))
	}

	config := kafka.ReaderConfig{
		Brokers: addrs,
		Topic:   topic,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return kafka.NewReader(config)
}

func ReadMsgs(ctx context.Context, r *kafka.Reader, handler func(kafka.Message) error) error {
	for {
		m, err := r.ReadMessage(ctx)
		if errorx.IsErrorExist(err) {
			return fmt.Errorf("[kafka][ReadMsgs]ReadMessage err: %w", err)
		}

		if err = handler(m); errorx.IsErrorExist(err) {
			return fmt.Errorf("[kafka][ReadMsgs]handler err: %w", err)
		}
	}
}
