# Ozon Test Task — GraphQL Posts & Comments System

Данный проект представляет собой систему блога с постами и иерархическими комментариями, реализованную на **Go** с использованием **GraphQL** (библиотека [gqlgen](https://gqlgen.com/)). Система поддерживает два типа хранилищ: **in-memory** и **PostgreSQL**.

## Build

#### Postgres storage
```sh
make docker-up
```

#### Inmemory storage
```sh
make run
```

#### Tests
```sh
make test
```

### Project Tree
```sh
├── cmd
│   └── service
│       └── main.go
├── docker-compose.yaml
├── Dockerfile
├── generated
│   └── graphql
│       ├── models_generated.go
│       ├── prelude_generated.go
│       ├── root_.generated.go
│       └── schema_generated.go
├── go.mod
├── go.sum
├── gqlgen.yml
├── internal
│   ├── cfg
│   │   └── config.go
│   ├── graphql
│   │   ├── dataloader
│   │   │   ├── dataloader.go
│   │   │   └── new.go
│   │   └── middleware
│   │       └── middleware.go
│   ├── handler
│   │   └── graphql
│   │       └── resolvers
│   │           ├── comment
│   │           │   ├── children.go
│   │           │   ├── comment.go
│   │           │   └── comment_test.go
│   │           ├── mutation
│   │           │   ├── add_comment.go
│   │           │   ├── create_post.go
│   │           │   ├── mutation.go
│   │           │   ├── mutation_test.go
│   │           │   └── set_post_comments_allowed.go
│   │           ├── query
│   │           │   ├── comment_by_post.go
│   │           │   ├── post.go
│   │           │   ├── posts.go
│   │           │   ├── query.go
│   │           │   └── query_test.go
│   │           ├── resolver.go
│   │           └── subscription
│   │               ├── comment_added.go
│   │               └── subscription.go
│   ├── models
│   │   └── models.go
│   ├── pkg
│   │   └── db
│   │       └── setup_repository.go
│   ├── repository
│   │   ├── comment_interface.go
│   │   ├── inmemory
│   │   │   ├── comment
│   │   │   │   ├── add_comment.go
│   │   │   │   ├── comments_by_post.go
│   │   │   │   ├── comment_test.go
│   │   │   │   └── new.go
│   │   │   └── post
│   │   │       ├── new.go
│   │   │       ├── post.go
│   │   │       ├── posts.go
│   │   │       ├── post_test.go
│   │   │       ├── save_post.go
│   │   │       └── set_comments_allowed.go
│   │   ├── interface.go
│   │   ├── mocks
│   │   │   ├── mock_CommentUC.go
│   │   │   └── mock_PostUC.go
│   │   ├── postgres
│   │   │   ├── comment
│   │   │   │   ├── add_comment.go
│   │   │   │   ├── comments_by_post.go
│   │   │   │   ├── comment_test.go
│   │   │   │   ├── mocks
│   │   │   │   │   └── mock_DB.go
│   │   │   │   └── new.go
│   │   │   └── post
│   │   │       ├── mocks
│   │   │       │   └── mock_DB.go
│   │   │       ├── new.go
│   │   │       ├── post.go
│   │   │       ├── posts.go
│   │   │       ├── post_test.go
│   │   │       ├── save_post.go
│   │   │       └── set_comments_allowed.go
│   │   └── repository.go
│   ├── service
│   │   ├── comment
│   │   │   ├── add_comment.go
│   │   │   ├── children.go
│   │   │   ├── comments_by_post.go
│   │   │   ├── comment_test.go
│   │   │   ├── interface.go
│   │   │   ├── mocks
│   │   │   │   └── mock_UseCase.go
│   │   │   └── new.go
│   │   ├── post
│   │   │   ├── create_post.go
│   │   │   ├── interface.go
│   │   │   ├── mocks
│   │   │   │   └── mock_UseCase.go
│   │   │   ├── new.go
│   │   │   ├── post.go
│   │   │   ├── posts.go
│   │   │   ├── post_test.go
│   │   │   └── set_post_comments_allowed.go
│   │   └── service.go
│   └── utils
│       └── cursor
│           └── cursor.go
├── Makefile
├── migrations
│   ├── 001-add-post.sql
│   └── 002-add-comment.sql
├── README.md
└── schema
    └── schema.graphqls
```