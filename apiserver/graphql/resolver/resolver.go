package graphql

import (
	"apiserver/graphql"
	"apiserver/handler"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MutationResolver    graphql.MutationResolver
	GrpcConnectionPools handler.GrpcConnectionPools
}

func NewResolver(
	mutationResolver graphql.MutationResolver,
	pools handler.GrpcConnectionPools,
) graphql.ResolverRoot {
	return &Resolver{
		MutationResolver:    mutationResolver,
		GrpcConnectionPools: pools,
	}
}

func NewMutationResolver() graphql.MutationResolver { return &mutationResolver{} }
