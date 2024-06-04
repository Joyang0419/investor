package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config MongoDB配置結構
type Config struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	ConnectTimeout  time.Duration
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}

// SetupConn 用於建立與MongoDB的連線
func SetupConn(config Config) *mongo.Client {
	// 建立MongoDB客戶端配置
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(config.ConnectTimeout).
		SetMaxPoolSize(config.MaxPoolSize).        // 設置連線池的最大連線數
		SetMaxConnIdleTime(config.MaxConnIdleTime) // 設置連線的最大空閒時間

	// 建立MongoDB連線
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(fmt.Sprintf("[SetupConn]mongo.Connect err: %v", err))
	}

	// 檢查連線
	if err = client.Ping(context.TODO(), nil); err != nil {
		panic(fmt.Sprintf("[SetupConn]client.Ping err: %v", err))
	}

	return client
}
