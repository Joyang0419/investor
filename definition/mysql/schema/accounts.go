package schema

import (
	"time"
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

func (Account) TableName() string {
	return "accounts"
}
