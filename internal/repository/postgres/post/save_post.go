package post

import (
	"context"
	"strconv"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	savePostQuery = `
		insert into posts (title, body, author, allow_comments, created_at) 
		values ($1, $2, $3, $4, $5)
		returning id, title, body, author, allow_comments, created_at
	`
)

func (r *post) SavePost(ctx context.Context, post models.Post) (models.Post, error) {
	out := models.Post{}
	var id int64

	err := r.db.QueryRow(ctx, savePostQuery,
		post.Title,
		post.Body,
		post.Author,
		post.AllowComments,
		post.CreatedAt,
	).Scan(
		&id,
		&out.Title,
		&out.Body,
		&out.Author,
		&out.AllowComments,
		&out.CreatedAt,
	)
	if err != nil {
		return models.Post{}, err
	}

	out.ID = strconv.FormatInt(id, 10)
	return out, nil
}
