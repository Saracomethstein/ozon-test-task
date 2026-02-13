package post

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *post) SetCommentsAllowed(ctx context.Context, postID int64, allow bool) (*models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}

	post.AllowComments = allow
	clone := *post

	return &clone, nil
}
