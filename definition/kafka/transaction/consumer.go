package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"

	"tools/errorx"
	kafka2 "tools/infra/kafka"
	"tools/logger"
	"tools/reflectx"
	"tools/serialization"
)

func ReceiveMessages(ctx context.Context, kafkaConn *kafka.Conn, topic string, processMessage func(data Data) error, opts ...kafka2.OptionConsumer) error {
	if reflectx.IsNil(kafkaConn) {
		return errors.New("[transaction][ReceiveMessages]kafkaConn is nil")
	}

	consumer := kafka2.NewKafkaConsumer(kafkaConn, topic, opts...)

	defer func() {
		if errClose := consumer.Close(); errorx.IsErrorExist(errClose) {
			logger.Error("[transaction][ReceiveMessages]consumer.Close err: %v", errClose)
		}
	}()

	for {
		msg, err := consumer.ReadMessage(ctx)
		if errorx.IsErrorExist(err) {
			if errors.Is(err, context.Canceled) {
				logger.Info("[transaction][ReceiveMessages]consumer.ReadMessage context canceled")
				return nil
			}

			return fmt.Errorf("[transaction][ReceiveMessages]consumer.ReadMessage err: %w", err)
		}

		result, err := serialization.JsonUnmarshal[Data](msg.Value)
		if errorx.IsErrorExist(err) {
			return fmt.Errorf("[transaction][ReceiveMessages]serialization.JsonUnmarshal err: %w", err)
		}

		if err = processMessage(result); err != nil {
			return fmt.Errorf("[transaction][ReceiveMessages]processMessage err: %w", err)
		}
	}
}
