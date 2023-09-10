package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/the-kwisatz-haderach/recipemaker/graph"
	"github.com/the-kwisatz-haderach/recipemaker/internal/config"
	db "github.com/the-kwisatz-haderach/recipemaker/internal/db"
)

func main() {
	ctx := context.Background()
	c := config.GetConfig()
	db := db.ConnectDb(ctx, c.DATABASE_URL)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Db: db}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", c.PORT)
	log.Fatal(http.ListenAndServe(":"+c.PORT, nil))
}
