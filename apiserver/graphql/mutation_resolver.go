package graphql

import (
	"context"
)

type ImplMutationResolver struct {
	// todo grpc 連線要inject
}

func (r *ImplMutationResolver) CreateInvestor(ctx context.Context, input CreateInvestorInput) (*Investor, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ImplMutationResolver) CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func NewMutationResolver() MutationResolver {
	return new(ImplMutationResolver)
}
