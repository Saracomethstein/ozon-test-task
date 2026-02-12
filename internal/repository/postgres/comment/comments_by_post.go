package comment

import (
	"context"
	"database/sql"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/pkg/errors"
)

const (
	getCommentsByPostQuery = `
		select id, post_id, parent_id, author, body, path, created_at
		from comments
		where post_id = $1
		  and ($2::text is null or (created_at, id) < ($2::text, $3::bigint))
		order by created_at desc, id desc
		limit $4
	`

	countCommentsByPostQuery = `select count(*) from comments where post_id = $1`
)

func (r *comment) GetCommentsByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	rows, err := r.db.Query(ctx, getCommentsByPostQuery,
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
			&c.Path,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *comment) TotalCountComments(ctx context.Context, postID int64) (int64, error) {
	var count int64

	err := r.db.QueryRow(ctx, countCommentsByPostQuery, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
