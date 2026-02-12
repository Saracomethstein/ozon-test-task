package post

import "github.com/jackc/pgx/v4/pgxpool"

type post struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) UseCase {
	return &post{
		db: db,
	}
}
