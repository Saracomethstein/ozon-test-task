package service

import (
	"github.com/Saracomethstein/ozon-test-task/internal/service/comment"
	"github.com/Saracomethstein/ozon-test-task/internal/service/post"
)

type Container struct {
	PostService    post.UseCase
	CommentService comment.UseCase
}

func New(
	post post.UseCase,
	comment comment.UseCase,
) *Container {
	return &Container{
		PostService:    post,
		CommentService: comment,
	}
}
