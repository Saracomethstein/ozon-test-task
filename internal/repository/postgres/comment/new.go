package comment

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

type comment struct {
	db DB
}

func New(db DB) repository.CommentUC {
	return &comment{
		db: db,
	}
}

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
