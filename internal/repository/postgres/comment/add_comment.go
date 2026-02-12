package comment

import (
	"context"
	"database/sql"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/pkg/errors"
)

const (
	addCommentQuery = `
		insert into comments (post_id, parent_id, author, body, path, created_at)
		values ($1, $2, $3, $4, $5, $6)
		returning id
	`

	checkPostAllowsCommentsQuery = `
		select allow_comments from posts where id = $1
	`

	checkParentCommentQuery = `
		select post_id from comments where id = $1
	`

	getCommentPathQuery    = `select path from comments where id = $1`
	updateCommentPathQuery = `update comments set path = $2 where id = $1`
)

func (r *comment) AddComment(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	err := r.db.QueryRow(ctx, addCommentQuery,
		comment.PostID,
		comment.ParentID,
		comment.Author,
		comment.Text,
		comment.Path,
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

func (r *comment) GetCommentPath(ctx context.Context, id int64) (string, error) {
	var path string

	err := r.db.QueryRow(ctx, getCommentPathQuery, id).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("comment not found")
		}
		return "", err
	}

	return path, nil
}

func (r *comment) SetCommentPath(ctx context.Context, id int64, path string) error {
	_, err := r.db.Exec(ctx, updateCommentPathQuery, id, path)
	return err
}
