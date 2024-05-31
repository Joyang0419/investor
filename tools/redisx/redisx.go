package redisx

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"

	"tools/reflectx"
)

var (
	errNilClient = errors.New("redis client is nil")
)

// setKeyNX 如果鍵不存在則設置一個帶有值和 TTL 的鍵
// 用途：用於分布式鎖
// 說明: 如果鍵不存在，則設置一個帶有值和 TTL 的鍵，返回 true；
// 如果鍵已存在，則不設置，返回 false。
func setKeyNX(
	ctx context.Context,
	redisClient *redis.Client,
	key string,
	value any,
	expiration ...time.Duration,
) (bool, error) {
	if reflectx.IsNil(redisClient) {
		return false, errNilClient
	}

	noExpiration := time.Duration(0)
	setTTL := noExpiration
	if len(expiration) > 0 {
		setTTL = expiration[0]
	}

	return redisClient.SetNX(ctx, key, value, setTTL).Result()
}
