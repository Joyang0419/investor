package accounts

import (
	"context"

	"gorm.io/gorm"

	"definition/mysql/schema"
)

func UpdateBalance(ctx context.Context, db *gorm.DB, accountID int64, amount float64) (err error) {
	return db.WithContext(ctx).Model(new(schema.Account)).Where("id = ?", accountID).Updates(
		map[string]interface{}{
			"balance": gorm.Expr("balance + ?", amount),
		},
	).Error
}
