package investor

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb"
	"repo/mongodb/schema/investor"
	"tools/mapx"
	"tools/mongodbx"
	"tools/numberx"
	"tools/slicex"
)

type Query struct {
	client *mongo.Client
}

func NewQuery(client *mongo.Client) *Query {
	return &Query{
		client: client,
	}
}

func (filter *GetInvestorsOptFilter) GetFindOptions() []mongodbx.FindOption {
	var findOptions []mongodbx.FindOption
	if numberx.CheckNotZero(filter.Page) {
		offset := (filter.Page - 1) * filter.PageSize
		findOptions = append(findOptions, mongodbx.WithOffset(int64(offset)))
	}
	if numberx.CheckNotZero(filter.PageSize) {
		findOptions = append(findOptions, mongodbx.WithLimit(int64(filter.PageSize)))
	}

	return findOptions
}
func (filter *GetInvestorsOptFilter) GetFilter() map[string]interface{} {
	var readyToCombinedMaps []map[string]any
	if slicex.IsNotEmpty(filter.InvestorIDs) {
		readyToCombinedMaps = append(readyToCombinedMaps, mongodbx.FilterFieldInValues("_id", filter.InvestorIDs))
	}
	if slicex.IsNotEmpty(filter.LoginAccounts) {
		readyToCombinedMaps = append(readyToCombinedMaps, mongodbx.FilterFieldInValues("loginAccount", filter.LoginAccounts))
	}

	return mapx.CombineMaps(readyToCombinedMaps...)
}

func (q *Query) GetInvestors(ctx context.Context, timeout time.Duration, OptFilter GetInvestorsOptFilter) ([]investor.Schema, error) {
	return mongodbx.All[[]investor.Schema](
		ctx,
		q.client,
		timeout,
		mongodb.InvestorStorage,
		OptFilter.GetFilter(),
		OptFilter.GetFindOptions()...,
	)
}
