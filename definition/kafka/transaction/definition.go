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

const Topic = "transaction"

type Data struct {
	ID              int64   // 交易ID
	Type            string  // 交易類型
	Amount          float64 // 交易金額
	AccountID       int64   // 交易帳戶ID
	TargetAccountID int64   // 目標帳戶ID
}

func (d *Data) ToKafkaMessage() (kafka.Message, error) {
	if reflectx.IsNil(d) {
		return kafka.Message{}, errors.New("[Data][ToKafkaMessage]data is nil")
	}
	bytes, err := serialization.JsonMarshal(d)
	if errorx.IsErrorExist(err) {
		return kafka.Message{}, fmt.Errorf("[Data][ToKafkaMessage]JsonMarshal err: %w", err)
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
		return errors.New("[WriteMessages]kafkaConn is nil")
	}
	ack := kafka.RequireNone
	if slicex.IsNotEmpty(requiredAck) {
		ack = requiredAck[0]
	}

	msgs := make([]kafka.Message, 0, len(data))
	for idx := range data {
		msg, err := data[idx].ToKafkaMessage()
		if err != nil {
			return fmt.Errorf("[WriteMessages]data.ToKafkaMessage err: %w", err)
		}
		msgs = append(msgs, msg)
	}

	return kafka2.NewKafkaSyncProducer(kafkaConn, Topic, ack).WriteMessages(ctx, msgs...)
}
