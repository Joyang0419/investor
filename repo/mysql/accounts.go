package mysql

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"tools/slicex"
)

var (
	ErrInvalidAccountIDs = errors.New("invalid account ids")
)

type Account struct {
	ID          uint64    `gorm:"column:id"`            // 帳戶ID
	Email       string    `gorm:"column:email"`         // 電子郵件
	Name        string    `gorm:"column:name"`          // 帳戶名稱
	Picture     string    `gorm:"column:picture"`       // 大頭貼URL
	Balance     float64   `gorm:"column:balance"`       // 餘額
	LastLoginAt time.Time `gorm:"column:last_login_at"` // 最後登入時間
	CreatedAt   time.Time `gorm:"column:created_at"`    // 建立時間
}

func (Account) TableName() string {
	return "accounts"
}

type IAccountQuery interface {
	CheckAccountIDsExist(ids []uint64) error
}

type AccountQuery struct {
	db *gorm.DB
}

func NewAccountQuery(db *gorm.DB) *AccountQuery {
	return &AccountQuery{db: db}
}

func (query *AccountQuery) CheckAccountIDsExist(ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	ids = slicex.RemoveDuplicateElement(ids)
	var count int64
	if err := query.db.Model(&Account{}).Where("id IN ?", ids).Count(&count).Error; err != nil {
		return fmt.Errorf("[AccountQuery][CheckAccountIDsExist]Count error: %w", err)
	}

	if count != int64(len(ids)) {
		return fmt.Errorf("[AccountQuery][CheckAccountIDsExist]err: %w, searchIDs: %v", ErrInvalidAccountIDs, ids)
	}

	return nil
}
