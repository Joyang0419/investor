package cmd

import (
	"log"

	graphql "apiserver/graphql/resolver"
	"apiserver/middleware"
	"apiserver/router"
	"tools/encryption"
	"tools/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// viper get config from env.yaml
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// load env.yaml by viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	// default value
	{
		viper.SetDefault("server.port", ":8080")

		viper.SetDefault("jwt.secret", []byte(`kmkdmvqejmriqiwngijoqpw`))

		viper.SetDefault("oauth2.google.client_id", "client_id")
		viper.SetDefault("oauth2.google.client_secret", "client_secret")
		viper.SetDefault("oauth2.google.redirect_url", "http://localhost:8080/auth/google/callback")
		viper.SetDefault("oauth2.google.scopes", []string{})
	}
}

// TODO graceful shutdown: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173
func runServerCmd(_ *cobra.Command, _ []string) {
	jwtEncryption := encryption.NewJWTEncryption[middleware.TokenInfo](encryption.JWTRequirements{
		SecretKey:     []byte(viper.GetString("jwt.secret")),
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

	if err := r.Run(viper.GetString("server.port")); err != nil {
		log.Fatalf("[runServerCmd]r.Run err: %v", err)
	}

	log.Print("[runServerCmd]success on port: 8080")
}
