package mysql

import (
	"errors"
	"fmt"
	"time"

	"tools/slicex"

	"gorm.io/gorm"
)

var validTransactionTypes = []string{"deposit", "withdraw", "transfer"}

var (
	errWrongTransactionType = errors.New("wrong transaction type")
)

type Transaction struct {
	ID              uint64    `gorm:"column:id"`                // 交易ID
	Type            string    `gorm:"column:type"`              // 交易類型
	Amount          float64   `gorm:"column:amount"`            // 交易金額
	AccountID       uint64    `gorm:"column:account_id"`        // 交易帳戶ID
	TargetAccountID uint64    `gorm:"column:target_account_id"` // 目標帳戶ID
	CreatedAt       time.Time `gorm:"column:created_at"`        // 交易時間
}

func (*Transaction) TableName() string {
	return "transactions"
}

// BeforeSave 鉤子函數，在保存前進行驗證
func (t *Transaction) BeforeSave(tx *gorm.DB) (err error) {
	_ = tx
	if slicex.IsElementInSlice(validTransactionTypes, t.Type) {
		return fmt.Errorf("[Transaction][BeforeSave]err: %w, type: %s", errWrongTransactionType, t.Type)
	}

	return nil
}
