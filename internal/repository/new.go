package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) UseCase {
	return &repository{
		db: db,
	}
}
