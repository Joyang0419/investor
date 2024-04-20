package mongodbx

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertMany[docType any](
	ctx context.Context,
	mongoDBClient *mongo.Client,
	timeout time.Duration,
	storage Storage,
	docs []docType,
) (*mongo.InsertManyResult, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()

	return mongoDBClient.
		Database(storage.Database).
		Collection(storage.Collection).
		InsertMany(ctx, convertToInterfaceSlice(docs))
}

func convertToInterfaceSlice[docType any](slice []docType) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
