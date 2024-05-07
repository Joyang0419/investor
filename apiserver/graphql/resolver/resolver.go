package graphql

import (
	"apiserver/graphql"
	"tools/grpcx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	QueryResolver       graphql.QueryResolver
	MutationResolver    graphql.MutationResolver
	GrpcConnectionPools *GrpcConnectionPools
}

func NewResolver(
	queryResolver graphql.QueryResolver,
	mutationResolver graphql.MutationResolver,
	pools *GrpcConnectionPools,
) graphql.ResolverRoot {
	return &Resolver{
		QueryResolver:       queryResolver,
		MutationResolver:    mutationResolver,
		GrpcConnectionPools: pools,
	}
}

type GrpcConnectionPools struct {
	MicroAuthGrpcConnPool *grpcx.GrpcConnectionPool
}

func NewGrpcConnectionPools(
	microAuthGrpcConnPool *grpcx.GrpcConnectionPool,
) *GrpcConnectionPools {
	return &GrpcConnectionPools{
		MicroAuthGrpcConnPool: microAuthGrpcConnPool,
	}
}

func NewQueryResolver() graphql.QueryResolver { return &queryResolver{} }

func NewMutationResolver() graphql.MutationResolver { return &mutationResolver{} }
