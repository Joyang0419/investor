package schema

type BalanceChangeLog struct {
	ID            int64   `gorm:"column:id"`             // ID
	AccountID     int64   `gorm:"column:account_id"`     // 交易帳戶ID
	TransactionID int64   `gorm:"column:transaction_id"` // 订单ID
	BeforeBalance float64 `gorm:"column:before_balance"` // 交易前余额
	AfterBalance  float64 `gorm:"column:after_balance"`  // 交易后余额
	CreatedAt     int64   `gorm:"column:created_at"`     // log產生时间
}

func (BalanceChangeLog) TableName() string {
	return "balance_change_log"
}
