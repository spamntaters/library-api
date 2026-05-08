package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/user/library-api/internal/config"
	"github.com/user/library-api/internal/db"
	"github.com/user/library-api/internal/graphql"
	"github.com/user/library-api/internal/openlibrary"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	store, err := db.NewStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer store.Close()

	olClient := openlibrary.NewClient(cfg.OpenLibraryURL)

	resolver := &graphql.Resolver{
		Store:    store,
		OLClient: olClient,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: resolver,
	}))

	http.Handle("/query", playground.Handler("GraphQL Playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("server started on http://localhost:%s", cfg.Port)
	log.Printf("graphql playground: http://localhost:%s/query", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
