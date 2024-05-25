package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"

	"apiserver/graphql"
)

type GraphqlHandler struct {
	Resolver graphql.ResolverRoot
}

func NewGraphqlHandler(resolver graphql.ResolverRoot) GraphqlHandler {
	return GraphqlHandler{Resolver: resolver}
}

func (h *GraphqlHandler) HandleGraphql() gin.HandlerFunc {
	graphqlServer := handler.NewDefaultServer(graphql.NewExecutableSchema(
		graphql.Config{Resolvers: h.Resolver}))

	return func(c *gin.Context) {
		graphqlServer.ServeHTTP(c.Writer, c.Request)
	}
}
