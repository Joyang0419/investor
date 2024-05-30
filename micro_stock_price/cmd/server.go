package cmd

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"micro_stock_price/server/taiwan_stock"
	"protos/micro_stock_price"
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

// TODO , 因為業務定義不明確，還要思考一下，先不花時間, 還沒砍掉的原因，是因為 要把排程的東西，移動到 Scheduler
func runServerCmd(_ *cobra.Command, _ []string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 0))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	s := grpc.NewServer()
	// 註冊服務
	micro_stock_price.RegisterTaiwanPriceServer(
		s,
		taiwan_stock.NewServer(),
	)

	logger.Info("gRPC server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("s.Serve err: %v", err)
	}
}
