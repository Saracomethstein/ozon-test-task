package post

import (
	"sync"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type post struct {
	mu    sync.RWMutex
	posts map[int64]*models.Post
	seq   int64
}

func New() repository.PostUC {
	return &post{
		posts: make(map[int64]*models.Post),
	}
}
