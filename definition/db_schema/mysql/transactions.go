package mysql

import (
	"context"
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

func (model *Transaction) Create(
	ctx context.Context,
	db *gorm.DB,
	transactionType string, amount float64,
	accountID, targetAccountID uint64, createdAt time.Time,
) (insertedID uint64, err error) {
	readyToInsert := Transaction{
		Type:            transactionType,
		Amount:          amount,
		AccountID:       accountID,
		TargetAccountID: targetAccountID,
		CreatedAt:       createdAt,
	}

	if slicex.IsElementInSlice(validTransactionTypes, readyToInsert.Type) {
		return 0, fmt.Errorf("[model.Transaction][Create]err: %w", errWrongTransactionType)
	}
	if db.WithContext(ctx).Create(&readyToInsert).Error != nil {
		return 0, fmt.Errorf("[model.Transaction][Create]Create err: %w", err)
	}

	return readyToInsert.ID, nil
}
