package comment

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	getRootCommentsByPostQuery = `
		select id, post_id, parent_id, author, body, created_at
		from comments
		where post_id = $1 and parent_id is null
			and ($2::text is null or (created_at, id) < ($2::text, $3::bigint))
		order by created_at desc, id desc
		limit $4
	`

	getChildCommentsQuery = `
		select id, post_id, parent_id, author, body, created_at
		from comments
		where parent_id = $1
	  		and ($2::text is null or (created_at, id) < ($2::text, $3::bigint))
		order by created_at desc, id desc
		limit $4
	`

	getChildCommentsBatchQuery = `
		select parent_id, id, post_id, author, body, created_at
		from comments
		where parent_id = any($1)
		order by parent_id, created_at desc, id desc
	`

	countCommentsByPostQuery = `select count(*) from comments where post_id = $1`
)

func (r *comment) GetRootByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	rows, err := r.db.Query(ctx, getRootCommentsByPostQuery,
		postID,
		afterCreatedAt,
		afterID,
		limit,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Comment{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	comments := make([]*models.Comment, 0, limit)
	for rows.Next() {
		c := models.Comment{}

		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.Author,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &c)
	}

	return comments, rows.Err()
}

func (r *comment) GetChild(ctx context.Context, parentID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	rows, err := r.db.Query(ctx, getChildCommentsQuery,
		parentID,
		afterCreatedAt,
		afterID,
		limit,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Comment{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	comments := make([]*models.Comment, 0, limit)
	for rows.Next() {
		c := models.Comment{}

		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.Author,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &c)
	}
	return comments, rows.Err()
}

func (r *comment) GetChildBatch(ctx context.Context, parentIDs []int64) ([]*models.Comment, error) {
	rows, err := r.db.Query(ctx, getChildCommentsBatchQuery, parentIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Comment{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := models.Comment{}

		err := rows.Scan(
			&c.ParentID,
			&c.ID,
			&c.PostID,
			&c.Author,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &c)
	}

	return comments, rows.Err()
}

func (r *comment) TotalCount(ctx context.Context, postID int64) (int64, error) {
	var count int64

	err := r.db.QueryRow(ctx, countCommentsByPostQuery, postID).Scan(&count)

	return count, err
}
