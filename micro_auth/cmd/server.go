package cmd

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"definition/micro_port"
	"micro_auth/server/investor"
	"protos/micro_auth"
	investor2 "repo/mongodb/investor"
	"tools/infra_conn"
	"tools/logger"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServerCmd(_ *cobra.Command, _ []string) {
	// 註冊基礎設施
	mongoDbConn, err := infra_conn.SetupMongoDB(
		infra_conn.MongoDBCfg{
			Host:            "127.0.0.1",
			Port:            27017,
			Username:        "root",
			Password:        "root",
			Database:        "admin",
			ConnectTimeout:  20 * time.Second,
			MaxPoolSize:     20,
			MaxConnIdleTime: 15 * time.Minute,
		},
	)
	if err != nil {
		logger.Fatal("[runDailyPriceCmd]infra_conn.SetupMongoDB err: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", micro_port.MicroAuthPort))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	s := grpc.NewServer()
	// 註冊Grpc服務
	// TODO viper timeout 放入環境變數
	micro_auth.RegisterInvestorServiceServer(
		s,
		investor.NewServer(
			investor2.NewQuery(mongoDbConn),
			investor2.NewCommand(mongoDbConn),
			30*time.Second,
		),
	)

	logger.Info("gRPC server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("s.Serve err: %v", err)
	}
}
