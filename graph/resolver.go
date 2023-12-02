package graph

import (
	"context"

	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Persistance interface {
	CreateRecipe(ctx context.Context, recipeName string, userID string) (*model.Recipe, error)
	GetRecipes(ctx context.Context, username string) ([]*model.Recipe, error)
}

type Resolver struct {
	Db Persistance
}
