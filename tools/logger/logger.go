package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
)

var logger = log.New()

func init() {
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.DateTime,
	})
	logger.SetLevel(log.InfoLevel)
}

func GormInfoLogger() gormLogger.Interface {
	return gormLogger.Default.LogMode(gormLogger.Info)
}

func loggerWithFuncNameAndFilename() *log.Entry {
	pc, file, line, ok := runtime.Caller(2) // Caller(2) 往上找兩層的呼叫者
	if !ok {
		file = "unknown"
		pc = 0
		line = 0
	}

	funcName := runtime.FuncForPC(pc).Name() // 獲得完整的函數名
	// 創建一個 Entry
	return logger.WithFields(log.Fields{
		"funcName": path.Base(funcName),              // 函數名的最後一部分
		"filepath": fmt.Sprintf("%s:%d", file, line), // 文件名和行號
	},
	)
}

func Info(format string, args ...any) {
	loggerWithFuncNameAndFilename().Infof(format, args...)
}

func Error(format string, args ...any) {
	loggerWithFuncNameAndFilename().Errorf(format, args...)
}

func Fatal(format string, args ...any) {
	loggerWithFuncNameAndFilename().Fatalf(format, args...)
}

func InfoWithContext(ctx context.Context, ctxKeys []string, format string, args ...any) {
	loggerWithFuncNameAndFilename().WithFields(extractValuesFromCtx(ctx, ctxKeys)).Infof(format, args...)
}

func ErrorWithContext(ctx context.Context, ctxKeys []string, format string, args ...any) {
	loggerWithFuncNameAndFilename().WithFields(extractValuesFromCtx(ctx, ctxKeys)).Infof(format, args...)
}

func FatalWithContext(ctx context.Context, ctxKeys []string, format string, args ...any) {
	loggerWithFuncNameAndFilename().WithFields(extractValuesFromCtx(ctx, ctxKeys)).Fatalf(format, args...)
}

func extractValuesFromCtx(ctx context.Context, keys []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range keys {
		if value := ctx.Value(key); value != nil {
			result[key] = value
			continue
		}
		result[key] = "undefined" // 如果key不存在，設置為"undefined"
	}

	return result
}

// GinLogger
// https://juejin.cn/post/6974640757374189575
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := fmt.Sprintf("%6v", endTime.Sub(startTime))

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		loggerWithFuncNameAndFilename().WithFields(log.Fields{
			"status_code": statusCode,
			"total_time":  latencyTime,
			"ip":          clientIP,
			"method":      reqMethod,
			"uri":         reqUri,
		}).Info("access")
	}
}
