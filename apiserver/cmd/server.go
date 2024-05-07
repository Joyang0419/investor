package cmd

import (
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

	// TODO config dir
	// viper get config from env.yaml
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// load env.yaml by viper
	// TODO 寫一下readme, 我不知道env 要放哪，也懶得找。建議放在跟config 一樣的位置, 可以學一下 env.template.yaml
	// TODO 加入makefile, make RunApiServer時，先自動複製 env.template.yaml 到 env.yaml
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Fatalf("error reading config file, %s", err)
	//}

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
func runServerCmd(cmd *cobra.Command, _ []string) {
	jwtEncryption := encryption.NewJWTEncryption[middleware.TokenInfo](encryption.JWTRequirements{
		SecretKey:     []byte(viper.GetString("jwt.secret")),
		SigningMethod: encryption.JWTSigningMethodHS256,
	})

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

	// gin router init
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
	)

	if err := r.Run(viper.GetString("server.port")); err != nil {
		log.Fatalf("[runServerCmd]r.Run err: %v", err)
	}

	log.Print("[runServerCmd]success on port: 8080")
}
