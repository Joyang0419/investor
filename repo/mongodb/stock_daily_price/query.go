package stock_daily_price

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb"
	"repo/mongodb/schema/investor"
	"tools/mongodbx"
)

type Query struct {
	client *mongo.Client
}

func (q *Query) DailyPrices(ctx context.Context, timeout time.Duration, opts ...mongodbx.FindOption) (data []investor.Schema, err error) {
	return mongodbx.All[[]investor.Schema](ctx, q.client, timeout, mongodb.StockDailyPriceStorage, nil, opts...)
}
