package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/cfg"
	"github.com/Saracomethstein/ozon-test-task/internal/graphql/dataloader"
	"github.com/Saracomethstein/ozon-test-task/internal/graphql/middleware"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers"
	"github.com/Saracomethstein/ozon-test-task/internal/pkg/db"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
	"github.com/Saracomethstein/ozon-test-task/internal/service"
	"github.com/Saracomethstein/ozon-test-task/internal/service/comment"
	"github.com/Saracomethstein/ozon-test-task/internal/service/post"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

var (
	production = flag.Bool("production", false, "use PostgreSQL storage")
)

func main() {
	flag.Parse()

	rContainer := GetRepositoryContainer()

	postSvc := post.New(rContainer.Post)
	commentSvc := comment.New(rContainer.Comment)
	commentLoader := dataloader.NewCommentLoader(rContainer.Comment)

	allSvc := service.New(postSvc, commentSvc)

	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolvers.New(allSvc)}))

	handlerWithDataloader := middleware.DataloaderMiddleware(*commentLoader)(srv)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handlerWithDataloader)
	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}

func GetRepositoryContainer() *repository.Container {
	cfg := cfg.New()

	if *production {
		log.Println("Starting with PostgreSQL storage")
		pgpool := db.SetupDB(*cfg)
		return db.NewPostgresContainer(pgpool)
	}

	log.Println("Starting with inmemory storage")
	return db.NewInmemoryContainer()
}
