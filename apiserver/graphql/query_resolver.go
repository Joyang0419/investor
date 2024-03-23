package graphql

import (
	"context"
)

type ImplQueryResolver struct {
	// todo grpc 連線要inject
}

func (r *ImplQueryResolver) Investors(ctx context.Context) ([]*Investor, error) {
	investor := new(Investor)

	return []*Investor{investor}, nil
}

func (r *ImplQueryResolver) Orders(ctx context.Context) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func NewQueryResolver() QueryResolver {
	return new(ImplQueryResolver)
}
