package post

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	getPostByIdQuery = `
		select id, title, body, author, allow_comments, created_at
		from posts
		where id = $1
	`
)

func (r *post) GetByID(ctx context.Context, postID int64) (*models.Post, error) {
	var out models.Post

	err := r.db.QueryRow(ctx, getPostByIdQuery, postID).Scan(
		&out.ID,
		&out.Title,
		&out.Body,
		&out.Author,
		&out.AllowComments,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &out, nil
}
