package cmd

import (
	"fmt"

	"apiserver/conf"
	"apiserver/handler"
	oauthx "tools/oauth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	graphql "apiserver/graphql/resolver"
	"apiserver/router"
	"tools/grpcx"
	"tools/logger"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Run:   runServerCmd,
}

// TODO graceful shutdown: https://learnku.com/docs/gin-gonic/1.5/examples-graceful-restart-or-stop/6173
func runServerCmd(cmd *cobra.Command, _ []string) {
	microAuthGrpcConnPool := grpcx.NewGrpcConnectionPool(
		cmd.Context(),
		grpcx.GetGrpcAddress(conf.Config.GrpcServer.MicroAuth.Domain, conf.Config.GrpcServer.MicroAuth.Port),
		conf.Config.GrpcServer.MicroAuth.MaxConnectionNum,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer func() {
		microAuthGrpcConnPool.CloseAllConnectionsOfPool()
	}()

	googleOauthConfig := oauthx.NewGoogleOauth(
		conf.Config.Oauth2.Google.ClientId,
		conf.Config.Oauth2.Google.ClientSecret,
		conf.Config.Oauth2.Google.RedirectUrl,
		conf.Config.Oauth2.Google.Scopes,
	)

	r := router.NewGinRouter(
		[]gin.HandlerFunc{logger.GinLogger()},
		router.Handler{
			AuthHandler: handler.NewAuthHandler(googleOauthConfig, handler.NewGrpcConnectionPools(
				microAuthGrpcConnPool,
			)),
			GraphqlHandler: handler.NewGraphqlHandler(
				graphql.NewResolver(
					graphql.NewMutationResolver(),
					handler.NewGrpcConnectionPools(
						microAuthGrpcConnPool,
					),
				),
			),
		},
	)

	logger.Info("[runServerCmd]success on port: %d", conf.Config.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", conf.Config.Server.Port)); err != nil {
		logger.Fatal("[runServerCmd]r.Run err: %v", err)
	}
}
