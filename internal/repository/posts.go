package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/pkg/errors"
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

func (r *repository) GetPosts(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error) {
	rows, err := r.db.Query(ctx, getPostsQuery, afterCreatedAt, afterID, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Post{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0, limit)
	for rows.Next() {
		var p models.Post
		var dbID int64

		err := rows.Scan(
			&dbID,
			&p.Title,
			&p.Body,
			&p.Author,
			&p.AllowComments,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.ID = strconv.FormatInt(dbID, 10)
		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *repository) TotalCount(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.QueryRow(ctx, totalCountQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
