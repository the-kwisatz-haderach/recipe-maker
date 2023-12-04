package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

func (p *Persistance) AddRecipe(ctx context.Context, input model.AddRecipeInput, userID string) (*model.Recipe, error) {
	var m = model.Recipe{Name: input.Name}
	q := `
		WITH new_recipe AS (
			INSERT INTO recipe (recipe_name) VALUES ($1) RETURNING id
		)
		INSERT INTO recipe_user (recipe_id, user_id, relation)
		SELECT id, $2, $3 FROM new_recipe 
			RETURNING recipe_id;
	`
	err := p.db.QueryRow(ctx, q, input.Name, userID, "owner").Scan(&m.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while creating recipe")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) UpdateRecipe(ctx context.Context, input model.UpdateRecipeInput) (*model.Recipe, error) {
	var m = model.Recipe{Name: input.Name}
	q := `
		UPDATE recipe SET recipe_name = $1 WHERE id = $2 RETURNING id;
	`
	err := p.db.QueryRow(ctx, q, input.Name, input.ID).Scan(&m.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while updating recipe")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) GetRecipes(ctx context.Context, userID string) ([]*model.Recipe, error) {
	q := `
		SELECT r.id, r.recipe_name FROM recipe r JOIN recipe_user ru ON ru.recipe_id = r.id WHERE ru.user_id = $1;
	`
	rows, err := p.db.Query(ctx, q, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while getting recipes")
		}
		return nil, err
	}
	defer rows.Close()

	var recipes []*model.Recipe
	for rows.Next() {
		var r model.Recipe
		err := rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, &r)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

func (p *Persistance) GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error) {
	var recipe model.Recipe
	q := `
		SELECT * FROM recipe WHERE id = $1;
	`
	err := p.db.QueryRow(ctx, q, recipeID).Scan(&recipe)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while getting recipe")
		}
		return nil, err
	}
	return &recipe, nil
}
