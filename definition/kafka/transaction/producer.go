package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"

	"tools/errorx"
	kafka2 "tools/infra/kafka"
	"tools/reflectx"
	"tools/serialization"
	"tools/slicex"
)

func (d *Data) ToKafkaMessage() (kafka.Message, error) {
	if reflectx.IsNil(d) {
		return kafka.Message{}, errors.New("[transaction][Data][ToKafkaMessage]data is nil")
	}
	bytes, err := serialization.JsonMarshal(d)
	if errorx.IsErrorExist(err) {
		return kafka.Message{}, fmt.Errorf("[transaction][Data][ToKafkaMessage]JsonMarshal err: %w", err)
	}

	return kafka.Message{
		Value: bytes,
	}, nil
}

func WriteMessages(ctx context.Context, kafkaConn *kafka.Conn, data []Data, requiredAck ...kafka.RequiredAcks) error {
	if slicex.IsEmpty(data) {
		return nil
	}
	if reflectx.IsNil(kafkaConn) {
		return errors.New("[transaction][WriteMessages]kafkaConn is nil")
	}
	ack := kafka.RequireNone
	if slicex.IsNotEmpty(requiredAck) {
		ack = requiredAck[0]
	}

	msgs := make([]kafka.Message, 0, len(data))
	for idx := range data {
		msg, err := data[idx].ToKafkaMessage()
		if err != nil {
			return fmt.Errorf("[transaction][WriteMessages]data.ToKafkaMessage err: %w", err)
		}
		msgs = append(msgs, msg)
	}

	return kafka2.NewKafkaSyncProducer(kafkaConn, Topic, ack).WriteMessages(ctx, msgs...)
}
