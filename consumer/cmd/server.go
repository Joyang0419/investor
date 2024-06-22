package cmd

import (
	"definition/kafka/transaction"

	"github.com/spf13/cobra"

	"consumer/conf"
	"consumer/service"
	kafka2 "tools/infra/kafka"
	"tools/infra/mysql"
	"tools/infra/redis"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

func runServerCmd(cmd *cobra.Command, _ []string) {
	// 註冊基礎設施
	mysqlConn := mysql.SetupConn(
		mysql.Config{
			Host:            conf.Config.MySQL.Host,
			Port:            conf.Config.MySQL.Port,
			Username:        conf.Config.MySQL.Username,
			Password:        conf.Config.MySQL.Password,
			Database:        conf.Config.MySQL.Database,
			MaxIdleConns:    conf.Config.MySQL.MaxIdleConns,
			MaxOpenConns:    conf.Config.MySQL.MaxOpenConns,
			ConnMaxLifeTime: conf.Config.MySQL.ConnMaxLifeTime,
		},
		nil,
	)

	redisClient := redis.SetupConn(
		redis.Config{
			Host:     conf.Config.Redis.Host,
			Port:     conf.Config.Redis.Port,
			Password: conf.Config.Redis.Password,
			DB:       conf.Config.Redis.DB,
		},
	)

	kafkaConn := kafka2.NewKafkaConn(kafka2.Config{
		Host: conf.Config.Kafka.Host,
		Port: conf.Config.Kafka.Port,
	})

	_, _ = mysqlConn, redisClient

	// 啟動連線池
	p := service.NewTaskPool(
		3,
		service.NewTask(
			cmd.Context(),
			kafka2.NewKafkaConsumer(kafkaConn, transaction.Topic, kafka2.WithGroupID("")),
			transaction.Handler(mysqlConn),
		),
	)
	defer p.Release()
	p.Start()
}
