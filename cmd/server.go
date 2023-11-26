package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/the-kwisatz-haderach/recipemaker/graph"
	"github.com/the-kwisatz-haderach/recipemaker/internal/auth"
	"github.com/the-kwisatz-haderach/recipemaker/internal/config"
	db "github.com/the-kwisatz-haderach/recipemaker/internal/db"
)

var envFlag = flag.String("env", "development", "environment")

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if *envFlag == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx := context.Background()
	c := config.GetConfig()
	db, close := db.ConnectDb(ctx, c.DATABASE_URL)
	defer close()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Db: db, Auth: &auth.Authenticator{}}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", c.PORT)
	if err := http.ListenAndServe(":"+c.PORT, nil); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
