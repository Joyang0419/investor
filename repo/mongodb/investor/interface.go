package investor

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"repo/mongodb/schema"
)

type ICommand interface {
	InsertMany(ctx context.Context, timeout time.Duration, data []schema.Investor) (*mongo.InsertManyResult, error)
}

type IQuery interface {
	GetInvestors(ctx context.Context, timeout time.Duration, OptFilter GetInvestorsOptFilter) ([]schema.Investor, error)
}

type GetInvestorsOptFilter struct {
	InvestorIDs   []string
	LoginAccounts []string
	Page          uint32
	PageSize      uint32
}
