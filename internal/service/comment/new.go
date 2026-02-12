package comment

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type commentService struct {
	repo *repository.Container
}

func New(repo *repository.Container) UseCase {
	return &commentService{
		repo: repo,
	}
}
