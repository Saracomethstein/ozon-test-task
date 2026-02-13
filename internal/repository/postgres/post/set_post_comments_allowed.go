package post

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

func (r *post) SetCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error) {
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	var out models.Post
	err = r.db.QueryRow(ctx, setPostCommentsAllowedQuery, id, allow).Scan(
		&out.ID,
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

	return &out, nil
}
