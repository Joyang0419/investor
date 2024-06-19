package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"

	"tools/logger"
)

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
