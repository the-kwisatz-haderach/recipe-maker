package graph

import (
	"context"

	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Persistance interface {
	AddRecipe(ctx context.Context, input model.AddRecipeInput, userID string) (*model.Recipe, error)
	UpdateRecipe(ctx context.Context, input model.UpdateRecipeInput) (*model.Recipe, error)
	GetRecipes(ctx context.Context, userID string) ([]*model.Recipe, error)
	GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error)
	GetIngredients(ctx context.Context, userID string) ([]*model.Ingredient, error)
	GetIngredient(ctx context.Context, ingredientID string) (*model.Ingredient, error)
	AddIngredient(ctx context.Context, input model.AddIngredientInput, userID string) (*model.Ingredient, error)
	AddPantryItem(ctx context.Context, input model.AddPantryItemInput, userID string) (*model.PantryItem, error)
	GetPantryItems(ctx context.Context, userID string) ([]*model.PantryItem, error)
	GetPantryItem(ctx context.Context, itemID string) (*model.PantryItem, error)
	UpdatePantryItem(ctx context.Context, input model.UpdatePantryItemInput) (*model.PantryItem, error)
}

type Resolver struct {
	Db Persistance
}
