package transactions

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"definition/mysql/schema"
	"tools/slicex"
)

var validTransactionTypes = []string{"deposit", "withdraw", "transfer"}

var (
	errWrongTransactionType = errors.New("wrong transaction type")
)

func Create(
	ctx context.Context,
	db *gorm.DB,
	transactionType string, amount float64,
	accountID, targetAccountID int64,
) (insertedID int64, err error) {
	readyToInsert := schema.Transaction{
		Type:            transactionType,
		Amount:          amount,
		AccountID:       accountID,
		TargetAccountID: targetAccountID,
	}

	if slicex.IsElementInSlice(validTransactionTypes, readyToInsert.Type) {
		return 0, fmt.Errorf("[transactions][Create]err: %w", errWrongTransactionType)
	}
	if db.WithContext(ctx).Create(&readyToInsert).Error != nil {
		return 0, fmt.Errorf("[transactions][Create]Create err: %w", err)
	}

	return readyToInsert.ID, nil
}
