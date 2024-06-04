package infra_conn

import (
	"context"
	"fmt"
	"time"
	"tools/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaCfg struct {
	Host          string
	Port          int
	Password      string
	Network       string
	WriteDeadline time.Duration
	ReadDeadline  time.Duration
}

func NewKafkaConn(config KafkaCfg, topic string) *kafka.Conn {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := kafka.DialLeader(context.Background(), config.Network, address, topic, 0)
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

func WriteMessage(c *kafka.Conn, msg []byte) {
	_, err := c.WriteMessages(
		kafka.Message{Value: msg},
	)
	if err != nil {
		logger.Fatal("[WriteMessage]c.WriteMessages err: %v", err)
	}
}

func ReadMessage(c *kafka.Conn, maxBytes int) []byte {
	msg, err := c.ReadMessage(maxBytes)
	if err != nil {
		logger.Fatal("[ReadMessage]c.ReadMessage err: %v", err)
	}

	return msg.Value
}

func ReadBatch(c *kafka.Conn, minBytes int, maxBytes int) []string {
	batchMsg := make([]string, 0)
	batch := c.ReadBatch(maxBytes, maxBytes)

	b := make([]byte, minBytes)
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		batchMsg = append(batchMsg, string(b[:n]))
	}

	err := batch.Close()
	if err != nil {
		logger.Fatal("[ReadBatch]batch.Close err: %v", err)
	}

	return batchMsg
}
