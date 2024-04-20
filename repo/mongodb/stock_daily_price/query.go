package stock_daily_price

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb"
	"tools/mongodbx"
)

type Query struct {
	client *mongo.Client
}

func (q *Query) DailyPrices(ctx context.Context, timeout time.Duration, opts ...mongodbx.FindOption) (prices []mongodb.StockDailyPriceSchema, err error) {
	return mongodbx.All[[]mongodb.StockDailyPriceSchema](ctx, q.client, timeout, mongodb.StockDailyPriceStorage, nil, opts...)
}
