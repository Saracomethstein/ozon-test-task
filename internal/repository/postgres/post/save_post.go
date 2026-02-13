package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	savePostQuery = `
		insert into posts (title, body, author, allow_comments, created_at) 
		values ($1, $2, $3, $4, $5)
		returning id, title, body, author, allow_comments, created_at
	`
)

func (r *post) Save(ctx context.Context, post models.Post) (models.Post, error) {
	out := models.Post{}

	err := r.db.QueryRow(ctx, savePostQuery,
		post.Title,
		post.Body,
		post.Author,
		post.AllowComments,
		post.CreatedAt,
	).Scan(
		&out.ID,
		&out.Title,
		&out.Body,
		&out.Author,
		&out.AllowComments,
		&out.CreatedAt,
	)
	if err != nil {
		return models.Post{}, err
	}

	return out, nil
}
