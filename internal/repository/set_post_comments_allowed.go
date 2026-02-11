package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	setPostCommentsAllowedQuery = `
		update posts
		set allow_comments = $2
		where id = $1
		returning id, title, body, author, allow_comments, created_at
	`
)

func (r *repository) SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error) {
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	var out models.Post
	var dbID int64
	err = r.db.QueryRow(ctx, setPostCommentsAllowedQuery, id, allow).Scan(
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
