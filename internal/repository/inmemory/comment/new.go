package comment

import (
	"sync"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type comment struct {
	mu       sync.Mutex
	comments map[int64]*models.Comment
	seq      int64
}

func New() repository.CommentUC {
	return &comment{
		comments: make(map[int64]*models.Comment),
	}
}
