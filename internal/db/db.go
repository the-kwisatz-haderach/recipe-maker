package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
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

func (p *Persistance) CreateRecipe(ctx context.Context, name string) (*model.Recipe, error) {
	var m model.Recipe
	m.RecipeName = name
	err := p.db.QueryRow(ctx, "insert into recipes (recipe_name) values ($1) returning id", name).Scan(&m.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while creating recipe")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) GetRecipes(ctx context.Context) ([]*model.Recipe, error) {
	log.Panic().Msg("not implemented")
	return nil, nil
}

func (p *Persistance) CreateUser(ctx context.Context, input model.SignupInput) (*model.User, error) {
	var m = model.User{Email: input.Email, Username: input.Username}
	err := p.db.QueryRow(ctx, "insert into users (username, password, email) values ($1,$2,$3) returning id", input.Username, input.Password, input.Email).Scan(&m.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while creating user")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) FindUser(ctx context.Context, username string) (*model.User, error) {
	var m model.User
	err := p.db.QueryRow(ctx, "select id, username, password, email from users where username = $1", username).Scan(&m.ID, &m.Username, &m.Password, &m.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Msgf("couldn't find user with username %s", username)
		} else {
			log.Error().Err(err).Msg("error when finding user")
		}
		return nil, err
	}
	return &m, nil
}
