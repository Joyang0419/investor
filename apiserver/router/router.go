package router

import (
	"apiserver/graphql"
	"apiserver/middleware"
	"tools/encryption"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func graphqlHandler(resolver graphql.ResolverRoot) gin.HandlerFunc {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(
		graphql.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func NewGinRouter(
	resolver graphql.ResolverRoot,
	middlewares []gin.HandlerFunc,
	jwtEncryption *encryption.JWTEncryption[middleware.TokenInfo],
) *gin.Engine {
	router := gin.New()

	for idx := range middlewares {
		router.Use(middlewares[idx])
	}

	// 加入JWT middleware
	router.POST("/query",
		//middleware.JWTMiddleware(jwtEncryption), // TODO 要分別針對 Schema 做權限控管
		graphqlHandler(resolver),
	)
	router.GET("/", playgroundHandler())

	return router
}
