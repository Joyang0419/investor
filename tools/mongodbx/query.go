package mongodbx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 寫共用

type FindOption func(*options.FindOptions)

// FilterDateRange 返回一個 它構造一個基於指定字段和日期範圍的過濾條件。
func FilterDateRange(field string, startDate, endDate time.Time) map[string]any {
	return map[string]any{
		field: map[string]any{
			"$gte": startDate.Unix(), // 大於或等於開始日期
			"$lte": endDate.Unix(),   // 小於或等於結束日期
		},
	}
}

func WithProjection(showFields []string) FindOption {
	return func(opts *options.FindOptions) {
		m := bson.M{}
		for _, trueField := range showFields {
			m[trueField] = 1
		}
		opts.SetProjection(m)
	}
}

// WithLimit 是一個選項函數，用於限制查詢結果的數量
func WithLimit(limit int64) FindOption {
	return func(opts *options.FindOptions) {
		opts.SetLimit(limit)
	}
}

// WithOffset 是一個選項函數，用於設置查詢結果的偏移量
func WithOffset(offset int64) FindOption {
	return func(opts *options.FindOptions) {
		opts.SetSkip(offset)
	}
}

func WithOrderBy(sortOrder map[string]int) FindOption {
	return func(opts *options.FindOptions) {
		opts.SetSort(sortOrder)
	}
}

func All[responseType any](
	ctx context.Context,
	mongoDBClient *mongo.Client,
	timeout time.Duration,
	storage Storage,
	filter map[string]any,
	opts ...FindOption,
) (resp responseType, err error) {
	// 檢查 responseType 是否為指針類型
	if reflect.ValueOf(resp).Kind() == reflect.Ptr {
		return resp, errors.New("responseType 不能是指針類型")
	}

	findOptions := options.Find()
	for _, opt := range opts {
		opt(findOptions)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	collection := mongoDBClient.Database(storage.Database).
		Collection(storage.Collection)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return resp, fmt.Errorf("[All]collection.Find err: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, ctx)

	if err = cursor.All(ctx, &resp); err != nil {
		return resp, fmt.Errorf("[All]cursor.All err: %w", err)
	}

	return resp, nil
}
