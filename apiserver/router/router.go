package router

import (
	"apiserver/graphql"
	"apiserver/handler"
	"apiserver/middleware"
	"tools/encryption"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(
	resolver graphql.ResolverRoot,
	middlewares []gin.HandlerFunc,
	jwtEncryption *encryption.JWTEncryption[middleware.TokenInfo],
) *gin.Engine {
	router := gin.New()

	for idx := range middlewares {
		router.Use(middlewares[idx])
	}

	// TODO 要分別針對 Schema 做權限控管
	router.POST("/query",
		//middleware.JWT(jwtEncryption),
		handler.GraphqlHandler(resolver),
	)
	router.GET("/", handler.PlayGroundHandler())

	return router
}
