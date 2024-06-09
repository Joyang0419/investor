package accounting

import (
	"context"

	"gorm.io/gorm"

	"definition/mysql/accounts"
)

type IQuery interface {
	IsAccountIDsExist(ctx context.Context, accountIDs ...int64) (bool, error)
}

type Query struct {
	mysqlDB *gorm.DB
}

func NewQuery(mysqlDB *gorm.DB) IQuery {
	return &Query{mysqlDB: mysqlDB}
}

func (q *Query) IsAccountIDsExist(ctx context.Context, accountIDs ...int64) (bool, error) {
	return accounts.IsCorrectAccountIDs(ctx, q.mysqlDB, accountIDs...)
}
