package post

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	getPostsQuery = `
		select id, title, body, author, allow_comments, created_at
		from posts
		where ($1::text is null or (created_at, id) < ($1::text, $2::bigint))
		order by created_at desc, id desc
		limit $3
	`

	totalCountQuery = `select count(*) from posts`
)

func (r *post) Get(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, getPostsQuery, afterCreatedAt, afterID, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*models.Post{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0, limit)
	for rows.Next() {
		var p models.Post

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Body,
			&p.Author,
			&p.AllowComments,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *post) TotalCount(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.QueryRow(ctx, totalCountQuery).Scan(&count)

	return count, err
}
