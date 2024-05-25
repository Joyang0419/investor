package schema

type StockDailyPrice struct {
	StockCode     string  `bson:"stockCode"`     // 股票代碼
	Volume        int64   `bson:"volume"`        // 當日成交股數
	HighestPrice  float64 `bson:"highestPrice"`  // 當日最高價
	LowestPrice   float64 `bson:"lowestPrice"`   // 當日最低價
	OpeningPrice  float64 `bson:"openingPrice"`  // 當日開盤價
	ClosingPrice  float64 `bson:"closingPrice"`  // 當日收盤價
	Change        float64 `bson:"change"`        // 漲跌價差
	DateTimestamp int64   `bson:"dateTimestamp"` // 資料歸屬時間點
}
