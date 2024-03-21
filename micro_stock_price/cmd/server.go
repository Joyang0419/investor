package cmd

import (
	"log"
	"net"

	"micro_stock_price/server"
	"protos/micro_stock_price"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
	micro_stock_price.RegisterTaiwanPriceServer(s, server.NewTaiwanPriceServer())

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("s.Serve err: %v", err)
	}
}
