package logger

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

func init() {
	Logger.SetFormatter(&log.JSONFormatter{})
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
