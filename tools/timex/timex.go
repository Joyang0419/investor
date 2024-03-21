package timex

import (
	"time"
)

func Int64ToTime(timestamp int64) time.Time {
	// 將毫秒轉換為秒（取整）和納秒（餘數）
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1000000 // 將餘數轉換為納秒

	// 使用time.Unix將時間戳轉換為time.Time
	return time.Unix(seconds, nanoseconds)
}
