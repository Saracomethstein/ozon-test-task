package comment

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type comment struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.CommentUC {
	return &comment{
		db: db,
	}
}
