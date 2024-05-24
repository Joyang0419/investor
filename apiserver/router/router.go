package router

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"apiserver/graphql"
	"apiserver/handler"
)

func NewGinRouter(
	resolver graphql.ResolverRoot,
	middlewares []gin.HandlerFunc,
	googleOauth oauth2.Config,
) *gin.Engine {
	router := gin.New()

	for idx := range middlewares {
		router.Use(middlewares[idx])
	}

	router.POST("/query",
		handler.GraphqlHandler(resolver),
	)
	router.GET("/login", handler.Login(googleOauth))

	router.GET("/auth/google/callback", handler.GoogleCallback(googleOauth))

	return router
}
