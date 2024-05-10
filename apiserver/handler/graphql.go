package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"

	"apiserver/graphql"
)

func GraphqlHandler(resolver graphql.ResolverRoot) gin.HandlerFunc {
	h := handler.NewDefaultServer(
		graphql.NewExecutableSchema(
			graphql.Config{Resolvers: resolver}),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
