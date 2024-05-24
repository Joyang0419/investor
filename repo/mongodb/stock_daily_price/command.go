package stock_daily_price

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb"
	"repo/mongodb/schema"
	"tools/mongodbx"
)

type Command struct {
	client *mongo.Client
}

func NewCommand(client *mongo.Client) *Command {
	return &Command{client: client}
}

func (cmd *Command) InsertMany(
	ctx context.Context,
	timeout time.Duration,
	data []schema.StockDailyPrice,
) (*mongo.InsertManyResult, error) {
	return mongodbx.InsertMany(
		ctx,
		cmd.client,
		timeout,
		mongodb.StockDailyPriceStorage,
		data,
	)
}
