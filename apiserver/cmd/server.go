package cmd

import (
	"apiserver/conf"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	graphql "apiserver/graphql/resolver"
	"apiserver/middleware"
	"apiserver/router"
	"definition/micro_port"
	"tools/encryption"
	"tools/grpcx"
	"tools/logger"
	"tools/oauth"

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

	conf.Init()
}

// TODO graceful shutdown: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173
func runServerCmd(cmd *cobra.Command, _ []string) {
	jwtEncryption := encryption.NewJWTEncryption[middleware.TokenInfo](encryption.JWTRequirements{
		SecretKey:     []byte(viper.GetString("jwt.secret")),
		SigningMethod: encryption.JWTSigningMethodHS256,
	})

	googleOauth := oauth.NewGoogleOauth(
		conf.Config.Oauth2.Google.ClientId,
		conf.Config.Oauth2.Google.ClientSecret,
		conf.Config.Oauth2.Google.RedirectUrl,
		conf.Config.Oauth2.Google.Scopes,
	)

	// grpc connection pool init
	microAuthGrpcConnPool := grpcx.NewGrpcConnectionPool(
		cmd.Context(),
		micro_port.GetGrpcAddress("localhost", micro_port.MicroAuthPort), // TODO domain viper
		10, // TODO viper
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer func() {
		microAuthGrpcConnPool.CloseAllConnectionsOfPool()
	}()

	r := router.NewGinRouter(
		graphql.NewResolver(
			graphql.NewQueryResolver(),
			graphql.NewMutationResolver(),
			graphql.NewGrpcConnectionPools(
				microAuthGrpcConnPool,
			),
		),
		[]gin.HandlerFunc{logger.GinLogger()},
		jwtEncryption,
		googleOauth,
	)

	if err := r.Run(conf.Config.Port); err != nil {
		log.Fatalf("[runServerCmd]r.Run err: %v", err)
	}

	log.Print("[runServerCmd]success on port: 8080")
}
