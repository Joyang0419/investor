package keys

import (
	"fmt"
)

type LockKeyF = string

const (
	accountingLockKey LockKeyF = "AccountingLockKey: %d"
)

func GetAccountingLockKey(accountID uint64) string {
	return fmt.Sprintf(accountingLockKey, accountID)
}
