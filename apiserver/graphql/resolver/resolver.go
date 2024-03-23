package graphql

import (
	"apiserver/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	QueryResolver    graphql.QueryResolver
	MutationResolver graphql.MutationResolver
}

func NewResolver(queryResolver graphql.QueryResolver, mutationResolver graphql.MutationResolver) graphql.ResolverRoot {
	return &Resolver{QueryResolver: queryResolver, MutationResolver: mutationResolver}
}

func (r *Resolver) Mutation() graphql.MutationResolver {
	return r.MutationResolver
}

func (r *Resolver) Query() graphql.QueryResolver {
	return r.QueryResolver
}

func NewQueryResolver() graphql.QueryResolver {
	return new(queryResolver)
}

func NewMutationResolver() graphql.MutationResolver {
	return new(mutationResolver)
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
