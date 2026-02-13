package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *post) GetByID(ctx context.Context, postID int64) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[postID]
	if !ok {
		return nil, ErrPostNotFound
	}

	clone := *post
	return &clone, nil
}
