package conf

import (
	"time"
)

var Config struct {
	Server struct {
		Port int
	}
	Jwt struct {
		Secret         string
		ExpireDuration time.Duration
	}
	MongoDB struct {
		Host            string
		Port            int
		Username        string
		Password        string
		Database        string
		ConnectTimeout  time.Duration
		MaxPoolSize     uint64
		MaxConnIdleTime time.Duration
	}
	App struct {
		DBTimeout time.Duration
	}
}
