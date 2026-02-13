package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *post) Save(ctx context.Context, post models.Post) (models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	id := r.seq

	clone := models.Post{
		ID:            id,
		Title:         post.Title,
		Body:          post.Body,
		Author:        post.Author,
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt,
	}

	r.posts[id] = &clone

	return clone, nil
}
