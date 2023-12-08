package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
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

func (p *Persistance) IsHealthy(ctx context.Context) bool {
	err := p.db.Ping(ctx)
	return err == nil
}
