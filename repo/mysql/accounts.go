package mysql

import (
	"time"
)

type Account struct {
	ID          int       `gorm:"column:id"`            // 帳戶ID
	Email       string    `gorm:"column:email"`         // 電子郵件
	AccountName string    `gorm:"column:account_name"`  // 帳戶名稱
	Picture     string    `gorm:"column:picture"`       // 大頭貼URL
	Balance     float64   `gorm:"column:balance"`       // 餘額
	LastLoginAt time.Time `gorm:"column:last_login_at"` // 最後登入時間
	CreatedAt   time.Time `gorm:"column:created_at"`    // 建立時間
}

func (Account) TableName() string {
	return "accounts"
}
