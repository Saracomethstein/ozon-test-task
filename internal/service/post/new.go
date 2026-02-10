package post

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type postService struct {
	repo repository.UseCase
}

func New(repo repository.UseCase) UseCase {
	return &postService{
		repo: repo,
	}
}
