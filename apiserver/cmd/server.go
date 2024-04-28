package cmd

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"

	graphql "apiserver/graphql/resolver"
	"apiserver/middleware"
	"apiserver/router"
	"tools/encryption"
	"tools/logger"
)

// TODO 改成server
var serverCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

// todo: graceful shutdown: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173

func runServerCmd(_ *cobra.Command, _ []string) {
	// TODO secretkey 變字串
	jwtEncryption := encryption.NewJWTEncryption[middleware.TokenInfo](encryption.JWTRequirements{
		SecretKey:     nil,
		SigningMethod: nil,
	})

	r := router.NewGinRouter(
		graphql.NewResolver(
			graphql.NewQueryResolver(),
			graphql.NewMutationResolver(),
		),
		[]gin.HandlerFunc{logger.GinLogger()},
		jwtEncryption,
	)

	// 啟動服務
	// todo viper 環境變數 :8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[runServerCmd]r.Run err: %v", err)
	}
	log.Print("[runServerCmd]success on port: 8080")
}
