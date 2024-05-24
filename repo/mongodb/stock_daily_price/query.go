package stock_daily_price

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb"
	"repo/mongodb/schema"
	"tools/mongodbx"
)

type Query struct {
	client *mongo.Client
}

func (q *Query) DailyPrices(ctx context.Context, timeout time.Duration, opts ...mongodbx.FindOption) (data []schema.StockDailyPrice, err error) {
	return mongodbx.All[[]schema.StockDailyPrice](ctx, q.client, timeout, mongodb.StockDailyPriceStorage, nil, opts...)
}
