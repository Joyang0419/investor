package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"tools/slicex"
)

type Account struct {
	ID          uint64    `gorm:"column:id"`            // 帳戶ID
	Email       string    `gorm:"column:email"`         // 電子郵件
	Name        string    `gorm:"column:name"`          // 帳戶名稱
	Picture     string    `gorm:"column:picture"`       // 大頭貼URL
	Balance     float64   `gorm:"column:balance"`       // 餘額
	LastLoginAt time.Time `gorm:"column:last_login_at"` // 最後登入時間
	CreatedAt   time.Time `gorm:"column:created_at"`    // 建立時間
	UpdatedAt   time.Time `gorm:"column:updated_at"`    // 更新時間
}

func (*Account) TableName() string {
	return "accounts"
}

func (model *Account) UpdateBalance(ctx context.Context, db *gorm.DB, accountID uint64, amount float64, updatedAt time.Time) (err error) {
	return db.WithContext(ctx).Model(new(Account)).Where("id = ?", accountID).Updates(
		map[string]interface{}{
			"balance":    gorm.Expr("balance + ?", amount),
			"updated_at": updatedAt,
		},
	).Error
}

func (model *Account) GetBalance(ctx context.Context, db *gorm.DB, accountID uint64) (balance float64, err error) {
	if err = db.WithContext(ctx).Model(new(Account)).Select("balance").Where("id = ?", accountID).Scan(&balance).Error; err != nil {
		return 0, fmt.Errorf("[updateBalance]Update err: %w, accountID: %d", err, accountID)
	}

	return balance, nil
}

func (model *Account) IsCorrectAccountIDs(ctx context.Context, db *gorm.DB, accountIDs ...uint64) (result bool, err error) {
	if len(accountIDs) == 0 {
		return false, errors.New("[model.Account][IsCorrectAccountIDs]Empty accountIDs")
	}

	uniqueAccountIDs := slicex.RemoveDuplicateElement(accountIDs)
	var count int64
	if err = db.WithContext(ctx).Model(&Account{}).Where("id IN ?", uniqueAccountIDs).Count(&count).Error; err != nil {
		return false, fmt.Errorf("[model.Account][IsCorrectAccountIDs]Query err: %w", err)
	}

	return int64(len(uniqueAccountIDs)) == count, nil
}
