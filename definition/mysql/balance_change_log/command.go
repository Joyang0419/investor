package balance_change_log

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"definition/mysql/schema"
)

func Create(
	ctx context.Context,
	db *gorm.DB,
	accountID int64, transactionID int64,
	beforeBalance, afterBalance float64,
) (err error) {
	readyToInsert := schema.BalanceChangeLog{
		AccountID:     accountID,
		TransactionID: transactionID,
		BeforeBalance: beforeBalance,
		AfterBalance:  afterBalance,
	}

	if db.WithContext(ctx).Create(&readyToInsert).Error != nil {
		return fmt.Errorf("[transactions][Create]Create err: %w", err)
	}

	return nil
}
