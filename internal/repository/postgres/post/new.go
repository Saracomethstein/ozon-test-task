package post

import (
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type post struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.PostUC {
	return &post{
		db: db,
	}
}
