package taiwan_stock

import (
	"repo/mysql/schema"
)

type ITaiwanStockRepo interface {
	BulkInsertDailyPrices(dailyPrices []schema.DailyPrice) error
}
