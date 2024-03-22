package logger

import (
	"bytes"
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// 定義一個helper函數用於捕獲日誌輸出
func captureLogOutput(operation func()) string {
	// 創建一個buffer來存儲日誌輸出
	var buf bytes.Buffer
	log.SetOutput(&buf)
	//log.SetFormatter(&log.TextFormatter{
	//	DisableTimestamp: true,
	//})

	operation()

	// 將buffer的內容轉為string
	return buf.String()
}

func TestInfo(t *testing.T) {
	output := captureLogOutput(func() {
		Info("Hello, %s", "world")
	})
	assert.Contains(t, output, "Hello, world")
}

func TestError(t *testing.T) {
	output := captureLogOutput(func() {
		Error("Error: %s", "something went wrong")
	})
	assert.Contains(t, output, "Error: something went wrong")
}

func TestInfoWithContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key2", "value2")
	ctxKeys := []string{"key1", "key2"}

	InfoWithContext(ctx, ctxKeys, "This is a test message")

	// 這裡可以添加更多的斷言來檢查日誌的輸出
}

func TestErrorWithContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key2", "value2")
	ctxKeys := []string{"key1", "key2"}

	ErrorWithContext(ctx, ctxKeys, "This is an error message")
	// 這裡可以添加更多的斷言來檢查日誌的輸出
}

func TestFatal(t *testing.T) {
	Fatal("Hello, %s", "world")
}

func TestFatalWithContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key2", "value2")
	ctxKeys := []string{"key1", "key2"}

	FatalWithContext(ctx, ctxKeys, "This is an error message")
}
