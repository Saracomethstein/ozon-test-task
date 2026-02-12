package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/cfg"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers"
	"github.com/Saracomethstein/ozon-test-task/internal/pkg/db"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
	"github.com/Saracomethstein/ozon-test-task/internal/repository/postgres"
	"github.com/Saracomethstein/ozon-test-task/internal/service"
	"github.com/Saracomethstein/ozon-test-task/internal/service/comment"
	"github.com/Saracomethstein/ozon-test-task/internal/service/post"
	"github.com/vektah/gqlparser/v2/ast"

	commentRepository "github.com/Saracomethstein/ozon-test-task/internal/repository/postgres/comment"
	postRepository "github.com/Saracomethstein/ozon-test-task/internal/repository/postgres/post"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cfg := cfg.New()

	pgpool := db.SetupDB(*cfg)

	postRepo := postRepository.New(pgpool)
	commentRepo := commentRepository.New(pgpool)

	postgresRepository := postgres.New(commentRepo, postRepo)

	repository := repository.New(postgresRepository)

	postService := post.New(repository)
	commentService := comment.New(repository)

	service := service.New(postService, commentService)

	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolvers.New(service)}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
