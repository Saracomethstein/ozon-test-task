package dataloader

import repository "github.com/Saracomethstein/ozon-test-task/internal/repository"

type CommentLoader struct {
	repo repository.CommentUC
}

func NewCommentLoader(repo repository.CommentUC) *CommentLoader {
	return &CommentLoader{repo: repo}
}
