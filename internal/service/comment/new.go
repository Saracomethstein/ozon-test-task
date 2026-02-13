package comment

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type Service struct {
	repo repository.CommentUC
}

func New(repo repository.CommentUC) *Service {
	return &Service{
		repo: repo,
	}
}
