package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

func (p *Persistance) GetIngredients(ctx context.Context, userID string) ([]*model.Ingredient, error) {
	q := `
		SELECT i.id, i.name FROM ingredient i JOIN ingredient_user iu ON iu.ingredient_id = i.id WHERE iu.user_id = $1;
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

	var ingredients []*model.Ingredient
	for rows.Next() {
		var ing model.Ingredient
		err := rows.Scan(&ing.ID, &ing.Name)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, &ing)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (p *Persistance) GetIngredient(ctx context.Context, ingredientID string) (*model.Ingredient, error) {
	var ingredient model.Ingredient
	q := `
		SELECT (id, name) FROM ingredient WHERE id = $1;
	`
	err := p.db.QueryRow(ctx, q, ingredientID).Scan(&ingredient.ID, &ingredient.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while getting recipes")
		}
		return nil, err
	}
	return &ingredient, nil
}

func (p *Persistance) AddIngredient(ctx context.Context, input model.AddIngredientInput, userID string) (*model.Ingredient, error) {
	var ingredient = model.Ingredient{Name: input.Name}
	q := `
		WITH new_ingredient AS (
			INSERT INTO ingredient (name) VALUES ($1) RETURNING id
		)
		INSERT INTO ingredient_user (ingredient_id, user_id)
		SELECT id, $2 FROM new_ingredient RETURNING ingredient_id;
	`
	err := p.db.QueryRow(ctx, q, input.Name, userID).Scan(&ingredient.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while getting recipes")
		}
		return nil, err
	}
	return &ingredient, nil
}
