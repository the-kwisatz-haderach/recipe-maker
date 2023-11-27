package graph

import (
	"context"

	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Persistance interface {
	CreateRecipe(ctx context.Context, name string) (*model.Recipe, error)
	GetRecipes(ctx context.Context) ([]*model.Recipe, error)
	// Signup(ctx context.Context, input model.SignupInput) (*model.User, error)
	// Login(ctx context.Context, input model.LoginInput) (*model.User, error)
	// Logout(ctx context.Context) error
}

type Resolver struct {
	Db Persistance
}
