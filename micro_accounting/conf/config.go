package conf

import (
	"time"
)

var Config struct {
	Server struct {
		Port int
	}
	MySQL struct {
		Host            string
		Port            int
		Username        string
		Password        string
		Database        string
		MaxIdleConns    uint8
		MaxOpenConns    uint8
		ConnMaxLifeTime time.Duration
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}
