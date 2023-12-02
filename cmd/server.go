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
	"github.com/the-kwisatz-haderach/recipemaker/internal/authservice"
	"github.com/the-kwisatz-haderach/recipemaker/internal/config"
	db "github.com/the-kwisatz-haderach/recipemaker/internal/db"
)

var envFlag = flag.String("env", "development", "environment")

func init() {
	config.InitConfiguration()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if *envFlag == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx := context.Background()

	db, close := db.ConnectDb(ctx, config.Config.DATABASE_URL)
	defer close()

	authService := authservice.NewAuthService(db)
	router := http.NewServeMux()

	// Authentication service
	router.HandleFunc("/login", authService.LoginHandler)
	router.HandleFunc("/signup", authService.SignupHandler)
	router.HandleFunc("/logout", authService.LogoutHandler)

	// GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Db: db}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", authService.Middleware(srv))

	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Config.PORT)
	if err := http.ListenAndServe(":"+config.Config.PORT, router); err != nil {
		log.Fatal().Err(err).Msg("server interrupted")
	}
}
