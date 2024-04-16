package schema

// DailyPrice 存放每日股價
type DailyPrice struct {
	ID            uint    `gorm:"primarykey"`
	StockCode     string  `gorm:"type:varchar(255);not null"`
	HighestPrice  float64 `gorm:"type:float;not null"`
	LowestPrice   float64 `gorm:"type:float;not null"`
	OpeningPrice  float64 `gorm:"type:float;not null"`
	ClosingPrice  float64 `gorm:"type:float;not null"`
	Volume        int64   `gorm:"type:bigint;not null"`
	Change        float64 `gorm:"type:float;not null"`
	DateTimestamp int64   `gorm:"type:bigint;not null"`
}

func (DailyPrice) TableName() string {
	return TableNameDailyPrice
}
