package infra_conn

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLCfg struct {
	Host            string        `mapstructure:"Host" default:"localhost"`
	Port            int           `mapstructure:"Port" default:"3306"`
	Username        string        `mapstructure:"Username" default:"root"`
	Password        string        `mapstructure:"Password" default:"root"`
	Database        string        `mapstructure:"Database" default:"dev"`
	MaxIdleConns    uint8         `mapstructure:"MaxIdleConns" default:"10"`
	MaxOpenConns    uint8         `mapstructure:"MaxOpenConns" default:"100"`
	ConnMaxLifeTime time.Duration `mapstructure:"ConnMaxLifeTime" default:"60m"`
}

func SetupMySQL(config MySQLCfg, logger logger.Interface) *gorm.DB {
	// dataSourceName
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(fmt.Sprintf("[SetupMySQL]gorm.Open err: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("[SetupMySQL]db.DB err: %v", err))
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(int(config.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(config.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifeTime)

	return db
}
