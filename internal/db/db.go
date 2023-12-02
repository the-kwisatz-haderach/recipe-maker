package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
	"github.com/the-kwisatz-haderach/recipemaker/internal/authservice"
)

func ConnectDb(ctx context.Context, conStr string) (*Persistance, func()) {
	pool, err := pgxpool.New(ctx, conStr)
	close := func() {
		pool.Close()
	}
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	return &Persistance{db: pool}, close
}

type Persistance struct {
	db *pgxpool.Pool
}

func (p *Persistance) CreateRecipe(ctx context.Context, recipeName string, userID string) (*model.Recipe, error) {
	var m model.Recipe
	m.RecipeName = recipeName
	q := `
		with new_recipe as (
			insert into recipes (recipe_name) values ($1) returning id
		)
		insert into recipe_roles (recipe_id, user_id, relation)
		select id, $2, $3 from new_recipe 
			returning recipe_id;
	`
	err := p.db.QueryRow(ctx, q, recipeName, userID, "owner").Scan(&m.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while creating recipe")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) GetRecipes(ctx context.Context, userID string) ([]*model.Recipe, error) {
	q := `
		select r.id, r.recipe_name from recipes r join recipe_roles rr on rr.recipe_id = r.id where rr.user_id = $1;
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
		err := rows.Scan(&r.ID, &r.RecipeName)
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

func (p *Persistance) CreateUser(ctx context.Context, input authservice.SignupInput) (*authservice.User, error) {
	var m = authservice.User{Email: input.Email, Username: input.Username}
	err := p.db.QueryRow(ctx, "insert into users (username, password, email) values ($1,$2,$3) returning id", input.Username, input.Password, input.Email).Scan(&m.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while creating user")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) FindUser(ctx context.Context, userID string) (*authservice.User, error) {
	var m authservice.User
	err := p.db.QueryRow(ctx, "select id, username, password, email from users where id = $1", userID).Scan(&m.ID, &m.Username, &m.Password, &m.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Msgf("couldn't find user with id %s", userID)
		} else {
			log.Error().Err(err).Msg("error when finding user")
		}
		return nil, err
	}
	return &m, nil
}
