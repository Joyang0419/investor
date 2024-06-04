package kafka_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func makeTopic() string {
	return fmt.Sprintf("kafka-go-%016x", rand.Int63())
}

func makeGroupID() string {
	return fmt.Sprintf("kafka-go-group-%016x", rand.Int63())
}

func makeTransactionalID() string {
	return fmt.Sprintf("kafka-go-transactional-id-%016x", rand.Int63())
}

func TestConn(t *testing.T) {
	d := new(kafka.Dialer)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dialContext, err := d.DialContext(ctx, "tcp", "localhost:9092")
	if err != nil {
		panic(err)
	}

	dialContext.Brokers()

}
