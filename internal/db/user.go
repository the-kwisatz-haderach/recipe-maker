package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/internal/authservice"
)

func (p *Persistance) CreateUser(ctx context.Context, input authservice.SignupInput) (*authservice.User, error) {
	var m = authservice.User{Email: input.Email, Username: input.Username}
	err := p.db.QueryRow(ctx, "INSERT INTO \"user\" (username, password, email) VALUES ($1,$2,$3) RETURNING id", input.Username, input.Password, input.Email).Scan(&m.ID)
	if err != nil {
		log.Debug().Err(err).Msg("error while creating user")
		return nil, err
	}
	return &m, nil
}

func (p *Persistance) FindUser(ctx context.Context, userID, username string) (*authservice.User, error) {
	var m authservice.User
	arg := userID
	q := "SELECT id, username, password, email FROM \"user\" WHERE id = $1;"
	if userID == "" {
		q = "SELECT id, username, password, email FROM \"user\" WHERE username = $1;"
		arg = username
	}
	err := p.db.QueryRow(ctx, q, arg).Scan(&m.ID, &m.Username, &m.Password, &m.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Msgf("couldn't find user with id %s or username %s", userID, username)
		} else {
			log.Error().Err(err).Msg("error when finding user")
		}
		return nil, err
	}
	return &m, nil
}
