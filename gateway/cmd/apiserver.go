package cmd

import (
	"log"

	"apiserver/handler"
	"apiserver/router"
	"apiserver/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func runApiServer(_ *cobra.Command, _ []string) {
	// cmd run

	// 準備好實作的service
	exampleService := service.NewExampleService()

	// todo logger 使用 logru
	r := router.NewGinRouter(router.NewHandlers(
		handler.NewExampleHandler(exampleService),
	), gin.Logger())

	// 啟動服務
	// todo viper 環境變數 :8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[runApiServer] route.Run err: %v", err)
	}
	log.Print("[runApiServer]success on port: 8080")

}

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "gateway short description",
	Long:  "gateway long description",
	Run:   runApiServer,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
}
