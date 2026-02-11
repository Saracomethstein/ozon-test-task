package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/pkg/errors"
)

const (
	getPostByIdQuery = `
		SELECT id, title, body, author, allow_comments, created_at
		FROM posts
		WHERE id = $1
	`
)

func (r *repository) GetPostById(ctx context.Context, postID int64) (*models.Post, error) {
	var out models.Post

	var dbID int64
	err := r.db.QueryRow(ctx, getPostByIdQuery, postID).Scan(
		&dbID,
		&out.Title,
		&out.Body,
		&out.Author,
		&out.AllowComments,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	out.ID = strconv.FormatInt(dbID, 10)
	return &out, nil
}
