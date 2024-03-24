package mongodb

type MongoStorage struct {
	Database   string
	Collection string
}

var (
	StockDailyPriceStorage = MongoStorage{
		Database:   "inventory",
		Collection: "stock_daily_price",
	}
)
