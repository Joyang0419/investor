package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Database   string
	Collection string
}

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

func Upsert[docType any](
	ctx context.Context,
	mongoDBClient *mongo.Client,
	timeout time.Duration,
	storage Storage,
	filter map[string]interface{},
	doc docType,
) (*mongo.UpdateResult, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()
	return mongoDBClient.Database(storage.Database).Collection(storage.Collection).UpdateOne(
		ctx,
		filter,
		bson.M{"$set": doc},
		options.Update().SetUpsert(true),
	)
}

func convertToInterfaceSlice[docType any](slice []docType) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
