package service

import (
	"github.com/Saracomethstein/ozon-test-task/internal/service/comment"
	"github.com/Saracomethstein/ozon-test-task/internal/service/post"
)

type Container struct {
	PostService    post.PostService
	CommentService comment.CommentService
}

func New() *Container {
	return &Container{
		PostService:    *post.New(),
		CommentService: *comment.New(),
	}
}
