package transaction

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"

	"tools/errorx"
	"tools/serialization"
)

func Handler(mysqlConn *gorm.DB) func(kafka.Message) error {
	return func(m kafka.Message) error {
		d, err := serialization.JsonUnmarshal[Data](m.Value)
		if errorx.IsErrorExist(err) {
			return fmt.Errorf("[kafka][transaction][Handler]JsonUnmarshal err: %w", err)
		}
		// do something
		_ = d
		panic("implement me")

		return nil
	}
}
