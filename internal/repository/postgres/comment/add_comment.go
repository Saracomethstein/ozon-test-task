package comment

import (
	"context"
	"database/sql"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/pkg/errors"
)

const (
	addCommentQuery = `
		insert into comments (post_id, parent_id, author, body, created_at)
		values ($1, $2, $3, $4, $5)
		returning id
	`

	checkPostAllowsCommentsQuery = `
		select allow_comments from posts where id = $1
	`

	checkParentCommentQuery = `
		select post_id from comments where id = $1
	`
)

func (r *comment) AddComment(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	err := r.db.QueryRow(ctx, addCommentQuery,
		comment.PostID,
		comment.ParentID,
		comment.Author,
		comment.Text,
		comment.CreatedAt,
	).Scan(&comment.ID)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *comment) CheckPostAllowComments(ctx context.Context, postID int64) (bool, error) {
	var allow bool

	err := r.db.QueryRow(ctx, checkPostAllowsCommentsQuery, postID).Scan(&allow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("post not found")
		}
		return false, err
	}

	return allow, nil
}

func (r *comment) CheckParentCommentExists(ctx context.Context, parentID int64) (int64, error) {
	var postID int64

	err := r.db.QueryRow(ctx, checkParentCommentQuery, parentID).Scan(&postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("parent comment not found")
		}
		return 0, err
	}

	return postID, nil
}
