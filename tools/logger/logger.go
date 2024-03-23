package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

func init() {
	Logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.DateTime,
	})
	Logger.SetLevel(log.InfoLevel)
}

func Info(format string, args ...any) {
	Logger.Infof(format, args...)
}

func Error(format string, args ...any) {
	Logger.Errorf(format, args...)
}

func Fatal(format string, args ...any) {
	Logger.Fatalf(format, args...)
}

func InfoWithContext(ctx context.Context, ctxKeys []string, format string, args ...any) {
	Logger.WithFields(extractValuesFromCtx(ctx, ctxKeys)).Infof(format, args...)
}

func ErrorWithContext(ctx context.Context, ctxKeys []string, format string, args ...any) {
	Logger.WithFields(extractValuesFromCtx(ctx, ctxKeys)).Infof(format, args...)
}

func FatalWithContext(ctx context.Context, ctxKeys []string, format string, args ...any) {
	Logger.WithFields(extractValuesFromCtx(ctx, ctxKeys)).Fatalf(format, args...)
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
		Logger.WithFields(log.Fields{
			"status_code": statusCode,
			"total_time":  latencyTime,
			"ip":          clientIP,
			"method":      reqMethod,
			"uri":         reqUri,
		}).Info("access")
	}
}
