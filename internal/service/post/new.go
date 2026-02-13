package post

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type Post struct {
	repo repository.PostUC
}

func New(repo repository.PostUC) *Post {
	return &Post{
		repo: repo,
	}
}
