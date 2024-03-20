package cmd

import (
	"log"
	"time"

	"gateway/handler"
	"gateway/router"
	"gateway/service"
	"tools/infra_conn"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func runGateway(_ *cobra.Command, _ []string) {
	// 設定Infra 連線 && 可以用makefile中: UpDevInfra 建置
	// todo 缺少一個 mongoDB init.db .js 自動建置 investor database
	mysqlDbConn, err := infra_conn.SetupMySQL(infra_conn.MySQLCfg{
		Host:            "localhost",
		Port:            "3306",
		Username:        "root",
		Password:        "root",
		Database:        "investor",
		MaxIdleConns:    20,
		MaxOpenConns:    20,
		ConnMaxLifeTime: 15 * time.Minute,
	}, nil) // todo logger logru
	if err != nil {
		log.Fatalf("[runGateway]infra_conn.SetupMySQL err: %v", err)
	}

	mongoDBConn, err := infra_conn.SetupMongoDB(infra_conn.MongoDBCfg{
		Host:            "localhost",
		Port:            "27017",
		Username:        "root",
		Password:        "root",
		Database:        "admin",
		ConnectTimeout:  20 * time.Second,
		MaxPoolSize:     20,
		MaxConnIdleTime: 10 * time.Minute,
	})

	_, _ = mysqlDbConn, mongoDBConn
	// 準備好實作的service
	exampleService := service.NewExampleService()

	// todo logger 使用 logru
	r := router.NewGinRouter(router.NewHandlers(
		handler.NewExampleHandler(exampleService),
	), gin.Logger())

	// 啟動服務
	// todo viper 環境變數 :8080
	if err = r.Run(":8080"); err != nil {
		log.Fatalf("[runGateway] route.Run err: %v", err)
	}
	log.Print("[runGateway]success on port: 8080")

}

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "gateway short description",
	Long:  "gateway long description",
	Run:   runGateway,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
}
