package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"

	"tools/reflectx"
)

var (
	ErrNilClient = errors.New("redis client is nil")
)

// SetLock 如果鍵不存在則設置一個帶有值和 TTL 的鍵
// 用途：用於分布式鎖
// 說明: 如果鍵不存在，則設置一個帶有值和 TTL 的鍵，返回 true；
// 如果鍵已存在，則不設置，返回 false。
func SetLock(
	ctx context.Context,
	redisClient *redis.Client,
	key string,
	value any,
	expiration ...time.Duration,
) (bool, error) {
	if reflectx.IsNil(redisClient) {
		return false, ErrNilClient
	}

	noExpiration := time.Duration(0)
	setTTL := noExpiration
	if len(expiration) > 0 {
		setTTL = expiration[0]
	}

	return redisClient.SetNX(ctx, key, value, setTTL).Result()
}

// ReleaseLock removes a lock set by SetNX.
func ReleaseLock(
	ctx context.Context,
	redisClient *redis.Client,
	key string,
) error {
	if reflectx.IsNil(redisClient) {
		return ErrNilClient
	}

	// 删除键来释放锁
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
