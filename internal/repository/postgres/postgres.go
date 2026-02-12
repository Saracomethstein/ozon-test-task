package postgres

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository/postgres/comment"
	"github.com/Saracomethstein/ozon-test-task/internal/repository/postgres/post"
)

type Container struct {
	Comment comment.UseCase
	Post    post.UseCase
}

func New(
	c comment.UseCase,
	p post.UseCase,
) *Container {
	return &Container{
		Comment: c,
		Post:    p,
	}
}
