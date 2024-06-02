package mysql

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"tools/numberx"
	"tools/slicex"
)

var validTransactionTypes = []string{"deposit", "withdraw", "transfer"}

var (
	ErrInvalidAmount = errors.New("invalid amount, amount should be greater than 0")
	ErrInvalidType   = fmt.Errorf("invalid transaction type, current valid types should be %v", validTransactionTypes)
)

type Transaction struct {
	ID              uint64    `gorm:"column:id"`                // 交易ID
	Type            string    `gorm:"column:type"`              // 交易類型
	Amount          float64   `gorm:"column:amount"`            // 交易金額
	AccountID       uint64    `gorm:"column:account_id"`        // 交易帳戶ID
	TargetAccountID uint64    `gorm:"column:target_account_id"` // 目標帳戶ID
	CreatedAt       time.Time `gorm:"column:created_at"`        // 交易時間
}

func (Transaction) TableName() string {
	return "transactions"
}

type ITransactionQuery interface {
	CheckValid(trx Transaction) error
}

type TransactionQuery struct {
	db           *gorm.DB
	AccountQuery IAccountQuery
}

func NewTransactionQuery(db *gorm.DB) ITransactionQuery {
	return &TransactionQuery{db: db}
}

func (query *TransactionQuery) CheckValid(trx Transaction) error {
	if numberx.IsLTE(trx.Amount, 0) {
		return ErrInvalidAmount
	}
	if slicex.IsElementNotInSlice(validTransactionTypes, trx.Type) {
		return ErrInvalidType
	}
	if err := query.AccountQuery.CheckAccountIDsExist([]uint64{trx.AccountID, trx.TargetAccountID}); err != nil {
		return fmt.Errorf("[Transaction][CheckValid]CheckAccountIDsExist err: %w", err)
	}

	return nil
}

type ITransactionCommand interface {
}

type TransactionCommand struct {
	db *gorm.DB
}

func NewTransactionCommand(db *gorm.DB) ITransactionCommand {
	return &TransactionCommand{db: db}
}
