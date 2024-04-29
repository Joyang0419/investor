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

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

// TODO viper to get secret and port
var secret = []byte("ejkorjqwiejriwejri")
var port = ":8080"

func init() {
	rootCmd.AddCommand(serverCmd)
}

// TODO graceful shutdown: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173
func runServerCmd(_ *cobra.Command, _ []string) {
	jwtEncryption := encryption.NewJWTEncryption[middleware.TokenInfo](encryption.JWTRequirements{
		SecretKey:     secret,
		SigningMethod: encryption.JWTSigningMethodHS256,
	})

	r := router.NewGinRouter(
		graphql.NewResolver(
			graphql.NewQueryResolver(),
			graphql.NewMutationResolver(),
		),
		[]gin.HandlerFunc{logger.GinLogger()},
		jwtEncryption,
	)

	if err := r.Run(port); err != nil {
		log.Fatalf("[runServerCmd]r.Run err: %v", err)
	}

	log.Print("[runServerCmd]success on port: 8080")
}
