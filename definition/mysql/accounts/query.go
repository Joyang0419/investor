package accounts

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"definition/mysql/schema"
	"tools/slicex"
)

func GetBalance(ctx context.Context, db *gorm.DB, accountID int64) (balance float64, err error) {
	if err = db.WithContext(ctx).Model(new(schema.Account)).Select("balance").Where("id = ?", accountID).Scan(&balance).Error; err != nil {
		return 0, fmt.Errorf("[accounts][GetBalance]Scan err: %w, accountID: %d", err, accountID)
	}

	return balance, nil
}

func IsCorrectAccountIDs(ctx context.Context, db *gorm.DB, accountIDs ...int64) (result bool, err error) {
	if len(accountIDs) == 0 {
		return false, errors.New("[accounts][IsCorrectAccountIDs]Empty accountIDs")
	}

	uniqueAccountIDs := slicex.RemoveDuplicateElement(accountIDs)
	var count int64
	if err = db.WithContext(ctx).Model(new(schema.Account)).Where("id IN ?", uniqueAccountIDs).Count(&count).Error; err != nil {
		return false, fmt.Errorf("[accounts][IsCorrectAccountIDs]Count err: %w", err)
	}

	return int64(len(uniqueAccountIDs)) == count, nil
}
