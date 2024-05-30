package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"micro_auth/conf"
	"micro_auth/server"
	"protos/micro_auth"
	investor2 "repo/mongodb/investor"
	"tools/encryption"
	"tools/infra_conn"
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
	mongoDbConn, err := infra_conn.SetupMongoDB(
		infra_conn.MongoDBCfg{
			Host:            conf.Config.MongoDB.Host,
			Port:            conf.Config.MongoDB.Port,
			Username:        conf.Config.MongoDB.Username,
			Password:        conf.Config.MongoDB.Password,
			Database:        conf.Config.MongoDB.Database,
			ConnectTimeout:  conf.Config.MongoDB.ConnectTimeout,
			MaxPoolSize:     conf.Config.MongoDB.MaxPoolSize,
			MaxConnIdleTime: conf.Config.MongoDB.MaxConnIdleTime,
		},
	)
	if err != nil {
		logger.Fatal("[runDailyPriceCmd]infra_conn.SetupMongoDB err: %v", err)
	}

	// 註冊gRPC服務
	grpcServer := grpc.NewServer()
	micro_auth.RegisterAuthServiceServer(
		grpcServer,
		server.NewAuth(
			investor2.NewQuery(mongoDbConn),
			investor2.NewCommand(mongoDbConn),
			conf.Config.App.DBTimeout,
			encryption.NewJWTEncryption[server.TokenInfo](
				encryption.JWTRequirements{
					SecretKey:      conf.Config.Jwt.Secret,
					SigningMethod:  encryption.JWTSigningMethodHS256,
					ExpireDuration: conf.Config.Jwt.ExpireDuration,
				},
			),
		),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Config.Server.Port))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	logger.Info("gRPC server listening at %v", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
