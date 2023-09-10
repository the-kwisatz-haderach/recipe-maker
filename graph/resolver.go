package graph

import "github.com/the-kwisatz-haderach/recipemaker/internal/db"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db *db.Persistance
}
