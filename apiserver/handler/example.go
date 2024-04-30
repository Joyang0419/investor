package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"apiserver/graphql"
)

// TODO filename example,go 該改了吧
func GraphqlHandler(resolver graphql.ResolverRoot) gin.HandlerFunc {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(
		graphql.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PlayGroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
