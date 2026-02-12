package dataloader

import repository "github.com/Saracomethstein/ozon-test-task/internal/repository"

type CommentLoader struct {
	repo *repository.Container
}

func NewCommentLoader(repo *repository.Container) *CommentLoader {
	return &CommentLoader{repo: repo}
}
