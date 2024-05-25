package graphql

import (
	"apiserver/graphql"
	"apiserver/handler"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	QueryResolver       graphql.QueryResolver
	MutationResolver    graphql.MutationResolver
	GrpcConnectionPools handler.GrpcConnectionPools
}

func NewResolver(
	queryResolver graphql.QueryResolver,
	mutationResolver graphql.MutationResolver,
	pools handler.GrpcConnectionPools,
) graphql.ResolverRoot {
	return &Resolver{
		QueryResolver:       queryResolver,
		MutationResolver:    mutationResolver,
		GrpcConnectionPools: pools,
	}
}

func NewQueryResolver() graphql.QueryResolver { return &queryResolver{} }

func NewMutationResolver() graphql.MutationResolver { return &mutationResolver{} }
