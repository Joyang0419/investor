package accounting

import (
	"context"

	"gorm.io/gorm"
)

type IQuery interface {
	IsAccountIDsExist(ctx context.Context, accountIDs []uint64) (bool, error)
}

type Query struct {
	mysqlDB *gorm.DB
}

func (q *Query) IsAccountIDsExist(ctx context.Context, accountIDs []uint64) (bool, error) {
	panic("implement me")
}

func NewQuery(mysqlDB *gorm.DB) IQuery {
	return &Query{
		mysqlDB: mysqlDB,
	}
}
