package comment

import "github.com/Saracomethstein/ozon-test-task/internal/repository"

type commentService struct {
	repo repository.UseCase
}

func New(repo repository.UseCase) UseCase {
	return &commentService{
		repo: repo,
	}
}
