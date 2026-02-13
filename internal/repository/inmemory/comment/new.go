package comment

import (
	"sync"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type comment struct {
	mu       sync.RWMutex
	comments map[int64]*models.Comment
	seq      int64

	byPost   map[int64][]int64
	byParent map[int64][]int64
	repoPost repository.PostUC
}

func New(repoPost repository.PostUC) repository.CommentUC {
	return &comment{
		comments: make(map[int64]*models.Comment),
		byPost:   make(map[int64][]int64),
		byParent: make(map[int64][]int64),
		repoPost: repoPost,
	}
}
