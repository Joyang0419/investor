package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"apiserver/graphql/model"
	"context"
	"fmt"
)

// CreateInvestor is the resolver for the createInvestor field.
func (r *mutationResolver) CreateInvestor(ctx context.Context, input model.CreateInvestorInput) (*model.Investor, error) {
	panic(fmt.Errorf("not implemented: CreateInvestor - createInvestor"))
}

// Investors is the resolver for the investors field.
func (r *queryResolver) Investors(ctx context.Context) ([]*model.Investor, error) {
	panic(fmt.Errorf("not implemented: Investors - investors"))
}
