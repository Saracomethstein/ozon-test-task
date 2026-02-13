package post

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type post struct {
	db DB
}

func New(db DB) repository.PostUC {
	return &post{
		db: db,
	}
}

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
