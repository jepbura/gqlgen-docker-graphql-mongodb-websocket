package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

import (
	"context"
	"fmt"

	"github.com/jepbura/go-server/pkg/infrastructure/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// CreateBook is the resolver for the createBook field.
func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	panic(fmt.Errorf("not implemented: CreateBook - createBook"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	// panic(fmt.Errorf("not implemented: Users - users"))
	fmt.Print("*********************************************\n")
	fmt.Print("Resolvers Users\n")
	fmt.Print("*********************************************\n")

	// Check if Usecase is nil
	if r.Usecase == nil {
		return nil, fmt.Errorf("not implemented: Users - users")
	}

	users, err := r.Usecase.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Books is the resolver for the books field.
func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	panic(fmt.Errorf("not implemented: Books - books"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }