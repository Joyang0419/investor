package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"micro_accounting/conf"
	"micro_accounting/service/accounting"
	"protos/micro_accounting"
	"tools/infra/mysql"
	"tools/infra/redis"
	"tools/logger"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

func runServerCmd(_ *cobra.Command, _ []string) {
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

	// Query
	accountQuery := accounting.NewQuery(mysqlConn)
	// Command
	accountCommand := accounting.NewCommand(mysqlConn, redisClient)

	// 註冊gRPC服務
	grpcServer := grpc.NewServer()
	micro_accounting.RegisterAccountingServiceServer(
		grpcServer,
		accounting.NewService(
			accountQuery,
			accountCommand,
		),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Config.Server.Port))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	logger.Info("gRPC service listening at %v", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
