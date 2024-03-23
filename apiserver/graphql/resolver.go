package graphql

type Resolver struct {
	QueryResolver    QueryResolver
	MutationResolver MutationResolver
}

func NewResolver(queryResolver QueryResolver, mutationResolver MutationResolver) ResolverRoot {
	return &Resolver{QueryResolver: queryResolver, MutationResolver: mutationResolver}
}

func (r *Resolver) Mutation() MutationResolver {
	return r.MutationResolver
}

func (r *Resolver) Query() QueryResolver {
	return r.QueryResolver
}
