package timex

import (
	"math/rand"
	"time"
)

func Int64ToTime(timestamp int64) time.Time {
	// 將毫秒轉換為秒（取整）和納秒（餘數）
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1000000 // 將餘數轉換為納秒

	// 使用time.Unix將時間戳轉換為time.Time
	return time.Unix(seconds, nanoseconds)
}

// SleepRandomSeconds 隨機在指定的秒數範圍內暫停執行。
// 參數 minSec 和 maxSec 分別表示睡眠時間的最小值和最大值（單位：秒）。
func SleepRandomSeconds(minSec, maxSec int) {
	// 為 rand 函數提供一個種子值以生成真正的隨機數
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成一個介於 minSec 和 maxSec 之間的隨機數
	duration := r.Intn(maxSec-minSec+1) + minSec
	// 將秒轉換為時間持續量 Duration 並睡眠
	time.Sleep(time.Duration(duration) * time.Second)
}

// GetCurrentTimestampSeconds returns Current Timestamp (10 digits)
func GetCurrentTimestampSeconds() uint64 {
	return uint64(time.Now().Unix())
}
