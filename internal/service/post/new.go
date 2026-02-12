package post

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type postService struct {
	repo *repository.Container
}

func New(repo *repository.Container) UseCase {
	return &postService{
		repo: repo,
	}
}
