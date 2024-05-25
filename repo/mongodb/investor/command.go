package investor

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb"
	"repo/mongodb/schema"
	"tools/mongodbx"
)

var _ ICommand = (*Command)(nil)

type Command struct {
	client *mongo.Client
}

func (cmd *Command) Upsert(ctx context.Context, timeout time.Duration, data schema.Investor) (*mongo.UpdateResult, error) {
	return mongodbx.Upsert[schema.Investor](
		ctx,
		cmd.client,
		timeout,
		mongodb.InvestorStorage,
		bson.M{"_id": data.ID},
		data,
	)
}

func NewCommand(client *mongo.Client) *Command {
	return &Command{client: client}
}
