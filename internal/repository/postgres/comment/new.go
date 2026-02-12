package comment

import "github.com/jackc/pgx/v4/pgxpool"

type comment struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) UseCase {
	return &comment{
		db: db,
	}
}
