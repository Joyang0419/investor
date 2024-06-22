package service

import (
	"context"

	"github.com/segmentio/kafka-go"

	"tools/errorx"
	kafka2 "tools/infra/kafka"
	"tools/logger"
	"tools/reflectx"
)

type Task struct {
	ctx      context.Context
	Consumer *kafka.Reader
	handler  func(kafka.Message) error
}

func NewTask(ctx context.Context, consumer *kafka.Reader, handler func(kafka.Message) error) *Task {
	return &Task{
		ctx:      ctx,
		Consumer: consumer,
		handler:  handler,
	}
}

func (t *Task) Do() {
	if reflectx.IsNil(t.Consumer) {
		logger.Fatal("[Task]Consumer is nil")
	}

	if err := kafka2.ReadMsgs(t.ctx, t.Consumer, t.handler); errorx.IsErrorExist(err) {
		c := t.Consumer.Config()
		logger.Error("[Task]ReadMsgs err: %v, topic: %s, groupID: %s",
			err,
			c.Topic, c.GroupID,
		)
	}
}
