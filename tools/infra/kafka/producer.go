package kafka

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"tools/errorx"
	"tools/logger"
	"tools/reflectx"
	"tools/slicex"
)

type OptionProducer func(*kafka.Writer)

func WithBalancer(balancer kafka.Balancer) OptionProducer {
	return func(w *kafka.Writer) {
		w.Balancer = balancer
	}
}

func WithTimeout(timeout time.Duration) OptionProducer {
	return func(w *kafka.Writer) {
		w.WriteTimeout = timeout
	}
}

func NewKafkaSyncProducer(conn *kafka.Conn, topic string, requiredAck kafka.RequiredAcks, opts ...OptionProducer) *kafka.Writer {
	if reflectx.IsNil(conn) {
		logger.Fatal("[NewKafkaSyncProducer]conn is nil")
	}
	if slicex.IsElementNotInSlice([]kafka.RequiredAcks{kafka.RequireNone, kafka.RequireOne, kafka.RequireAll}, requiredAck) {
		logger.Fatal("[NewKafkaSyncProducer]requiredAck is invalid")
	}

	brokers, err := conn.Brokers()
	if errorx.IsErrorExist(err) {
		logger.Fatal("[NewKafkaSyncProducer]conn.Brokers err: %v", err)
	}

	addrs := make([]string, 0, len(brokers))
	for idx := range brokers {
		addrs = append(addrs, fmt.Sprintf("%s:%d", brokers[idx].Host, brokers[idx].Port))
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(addrs...),
		Topic:        topic,
		RequiredAcks: requiredAck,
		Async:        false,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}
