package accounting

import (
	"context"

	"gorm.io/gorm"

	"definition/db_schema/mysql"
)

type IQuery interface {
	IsAccountIDsExist(ctx context.Context, accountIDs []uint64) (bool, error)
}

type Query struct {
	mysqlDB *gorm.DB
}

func (q *Query) IsAccountIDsExist(ctx context.Context, accountIDs []uint64) (bool, error) {
	return new(mysql.Account).IsCorrectAccountIDs(ctx, q.mysqlDB, accountIDs...)
}

func NewQuery(mysqlDB *gorm.DB) IQuery {
	return &Query{
		mysqlDB: mysqlDB,
	}
}
