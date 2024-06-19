package kafka

import (
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"

	"tools/errorx"
	"tools/reflectx"
	"tools/serialization"
)

func JsonFormatToKafkaMsg[dataType any](data dataType) (kafka.Message, error) {
	if reflectx.IsNil(data) {
		return kafka.Message{}, errors.New("[kafka][JsonFormatToKafkaMsg]data is nil")
	}

	bytes, err := serialization.JsonMarshal(data)
	if errorx.IsErrorExist(err) {
		return kafka.Message{}, fmt.Errorf("[kafka][JsonFormatToKafkaMsg]JsonMarshal err: %w", err)
	}

	return kafka.Message{
		Value: bytes,
	}, nil
}
