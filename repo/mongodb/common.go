package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 寫共用

type QueryOption func(*options.FindOptions)

func WithDateRange(field string, startDate, endDate time.Time) QueryOption {
	return func(opts *options.FindOptions) {
		opts.SetHint(bson.M{field: bson.M{"$gte": startDate.Unix(), "$lte": endDate.Unix()}})
	}
}

func WithProjection(trueFields, falseFields []string) QueryOption {
	return func(opts *options.FindOptions) {
		m := bson.M{}
		for _, trueField := range trueFields {
			m[trueField] = 1
		}
		for _, falseField := range falseFields {
			m[falseField] = 0
		}
		opts.SetProjection(m)
	}
}

// WithLimit 是一個選項函數，用於限制查詢結果的數量
func WithLimit(limit int64) QueryOption {
	return func(opts *options.FindOptions) {
		opts.SetLimit(limit)
	}
}

// WithOffset 是一個選項函數，用於設置查詢結果的偏移量
func WithOffset(offset int64) QueryOption {
	return func(opts *options.FindOptions) {
		opts.SetSkip(offset)
	}
}

// todo test

func All[responseType any](
	ctx context.Context,
	mongoDBClient *mongo.Client,
	timeout time.Duration,
	storage MongoStorage,
	opts ...QueryOption,
) (resp responseType, err error) {
	findOptions := options.Find()
	for _, opt := range opts {
		opt(findOptions)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	collection := mongoDBClient.Database(storage.Database).
		Collection(storage.Collection)
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return resp, fmt.Errorf("[All]collection.Find err: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	if err = cursor.All(ctx, resp); err != nil {
		return resp, fmt.Errorf("[All]cursor.All err: %w", err)
	}

	return resp, nil
}
