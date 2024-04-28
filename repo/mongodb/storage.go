package mongodb

import (
	"tools/mongodbx"
)

var (
	StockDailyPriceStorage = mongodbx.Storage{
		Database:   InvestorDatabase,
		Collection: stockDailyPriceCollection,
	}
	InvestorStorage = mongodbx.Storage{
		Database:   InvestorDatabase,
		Collection: investorCollection,
	}
)
