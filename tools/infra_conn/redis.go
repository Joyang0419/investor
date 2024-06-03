package infra_conn

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisCfg struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func SetupRedis(config RedisCfg) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		panic(fmt.Sprintf("[SetupRedis]client.Ping err: %v", err))
	}

	return client
}
