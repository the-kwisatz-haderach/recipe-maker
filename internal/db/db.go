package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-kwisatz-haderach/recipemaker/graph/model"
)

func ConnectDb(ctx context.Context, conStr string) *Persistance {
	pool, err := pgxpool.New(ctx, conStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	data, err := os.ReadFile("internal/db/main.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(data)); err != nil {
		log.Fatal(err)
	}
	return &Persistance{db: pool}
}

type Persistance struct {
	db *pgxpool.Pool
}

func (p *Persistance) CreateRecipe(ctx context.Context, name string) *model.Recipe {
	var m *model.Recipe
	row := p.db.QueryRow(ctx, "INSERT INTO recipes (recipe_name) VALUES ($1)", name).Scan()
	return m
}
