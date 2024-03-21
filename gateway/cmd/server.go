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
	// 設定Infra 連線 && 可以用makefile中: UpDevInfra 建置
	// todo viper db 設定
	// todo 缺少一個 mongoDB init.db .js 自動建置 investor database
	// todo 這只是一個範例，基本上，這邊不會有db連線，全部應該是 grpc Client
	// todo logger logru
	mysqlDbConn, err := infra_conn.SetupMySQL(infra_conn.MySQLCfg{
		Host:            "localhost",
		Port:            "3306",
		Username:        "root",
		Password:        "root",
		Database:        "investor",
		MaxIdleConns:    20,
		MaxOpenConns:    20,
		ConnMaxLifeTime: 15 * time.Minute,
	}, nil)
	if err != nil {
		log.Fatalf("[runServerCmd]infra_conn.SetupMySQL err: %v", err)
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
	if err != nil {
		log.Fatalf("[runGateway]infra_conn.SetupMongoDB err: %v", err)
	}

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
		log.Fatalf("[runServerCmd] route.Run err: %v", err)
	}
	log.Print("[runServerCmd]success on port: 8080")
}
