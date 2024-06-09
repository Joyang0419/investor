package schema

import (
	"time"
)

type Transaction struct {
	ID              int64     `gorm:"column:id"`                // 交易ID
	Type string `gorm:"column:type"` // 交易類型
	Amount          float64   `gorm:"column:amount"`            // 交易金額
	AccountID       int64     `gorm:"column:account_id"`        // 交易帳戶ID
	TargetAccountID int64     `gorm:"column:target_account_id"` // 目標帳戶ID
	CreatedAt       time.Time `gorm:"column:created_at"`        // 交易時間
}

func (*Transaction) TableName() string {
	return "transactions"
}
