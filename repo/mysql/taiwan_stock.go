package mysql

import (
	"gorm.io/gorm"

	"repo/mysql/schema"
)

type TaiwanStockRepo struct {
	mysqlConn *gorm.DB
}

func NewTaiwanStockRepo(mysqlConn *gorm.DB) *TaiwanStockRepo {
	return &TaiwanStockRepo{mysqlConn: mysqlConn}
}

func (r *TaiwanStockRepo) BulkInsertDailyPrices(dailyPrices []schema.DailyPrice) error {
	return r.mysqlConn.Create(&dailyPrices).Error
}
