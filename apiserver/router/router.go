package router

import (
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"

	"apiserver/graphql"
	"apiserver/handler"
)

func NewGinRouter(
	resolver graphql.ResolverRoot,
	middlewares []gin.HandlerFunc,
	googleOauthConfig oauth2.Config,
) *gin.Engine {
	router := gin.New()

	for idx := range middlewares {
		router.Use(middlewares[idx])
	}

	router.POST("/query",
		//middleware.JWT(jwtEncryption),
		handler.GraphqlHandler(resolver),
	)

	// 尚未進到系統前，都使用Restful api(example: login, callback, ...)
	router.GET("/auth/google/login", handler.GoogleOauthLoginHandler(googleOauthConfig))
	router.GET("/auth/google/login/callback", handler.GoogleOauthCallbackHandler(googleOauthConfig))

	return router
}
