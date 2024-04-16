package cmd

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"micro_stock_price/server/taiwan_stock"
	"protos/micro_stock_price"
	"repo/mysql"
	"tools/logger"

	"github.com/spf13/cobra"
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	s := grpc.NewServer()
	// 註冊服務
	micro_stock_price.RegisterTaiwanPriceServer(
		s,
		taiwan_stock.NewTaiwanPriceServer(mysql.NewTaiwanStockRepo(nil)),
	)

	logger.Info("gRPC server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("s.Serve err: %v", err)
	}
}
