package transaction

const Topic = "transaction"

type Data struct {
	ID              int64   // 交易ID
	Type            string  // 交易類型
	Amount          float64 // 交易金額
	AccountID       int64   // 交易帳戶ID
	TargetAccountID int64   // 目標帳戶ID
}
