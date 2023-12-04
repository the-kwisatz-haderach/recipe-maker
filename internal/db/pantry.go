package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

func (p *Persistance) AddPantryItem(ctx context.Context, input model.AddPantryItemInput, userID string) (*model.PantryItem, error) {
	var pantryItem model.PantryItem
	q := `
		WITH new_pantry_item AS (
			INSERT INTO pantry_item (ingredient_id, quantity)
			VALUES ($1, $2)
			ON CONFLICT (ingredient_id) DO NOTHING
			RETURNING id
		)
		INSERT INTO pantry_item_user (pantry_item_id, user_id)
		SELECT id, $3 FROM new_pantry_item RETURNING pantry_item_id;
	`
	err := p.db.QueryRow(ctx, q, input.IngredientID, input.Quantity, userID).Scan(&pantryItem.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while adding pantry item")
		}
		return nil, err
	}
	err = p.db.
		QueryRow(ctx, "SELECT pi.id, i.name, pi.quantity, pi.unit FROM ingredient i INNER JOIN pantry_item pi ON i.id = pi.ingredient_id WHERE pi.id = $1", pantryItem.ID).
		Scan(&pantryItem.ID, &pantryItem.Name, &pantryItem.Quantity, &pantryItem.Unit)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while adding pantry item")
		}
		return nil, err
	}
	return &pantryItem, nil
}

func (p *Persistance) UpdatePantryItem(ctx context.Context, input model.UpdatePantryItemInput) (*model.PantryItem, error) {
	var pantryItem model.PantryItem
	p.db.QueryRow(ctx, "UPDATE pantry_item SET quantity = $1 WHERE id = $2;", input.Quantity, input.ID)
	return &pantryItem, nil
}

func (p *Persistance) GetPantryItems(ctx context.Context, userID string) ([]*model.PantryItem, error) {
	q := `
		SELECT pi.id, i.name, pi.quantity, pi.unit FROM pantry_item pi
		INNER JOIN pantry_item_user piu ON piu.pantry_item_id = pi.id
		INNER JOIN ingredient i ON i.id = pi.ingredient_id
		WHERE piu.user_id = $1;
	`
	rows, err := p.db.Query(ctx, q, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("no rows")
		} else {
			log.Error().Err(err).Msg("unknown error while getting pantry items")
		}
		return nil, err
	}
	defer rows.Close()

	var pantryItems []*model.PantryItem
	for rows.Next() {
		var pi model.PantryItem
		err := rows.Scan(&pi.ID, &pi.Name, &pi.Quantity, &pi.Unit)
		if err != nil {
			return nil, err
		}
		pantryItems = append(pantryItems, &pi)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pantryItems, nil
}
