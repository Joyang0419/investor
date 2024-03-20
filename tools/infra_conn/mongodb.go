package infra_conn

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBCfg MongoDB配置結構
type MongoDBCfg struct {
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	ConnectTimeout  time.Duration
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}

// SetupMongoDB 用於建立與MongoDB的連線
func SetupMongoDB(config MongoDBCfg) (*mongo.Client, error) {
	// 建立MongoDB客戶端配置
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(config.ConnectTimeout).
		SetMaxPoolSize(config.MaxPoolSize).        // 設置連線池的最大連線數
		SetMaxConnIdleTime(config.MaxConnIdleTime) // 設置連線的最大空閒時間

	// 建立MongoDB連線
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("[SetupMongoDB]mongo.Connect err: %w", err)
	}

	// 檢查連線
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("[SetupMongoDB]client.Ping err: %w", err)
	}

	return client, nil
}
