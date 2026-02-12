package repository

import "github.com/Saracomethstein/ozon-test-task/internal/repository/postgres"

type Container struct {
	DB *postgres.Container
	// mem inmemory.UseCase
}

func New(
	db *postgres.Container,
	// mem inmemory.UseCase,
) *Container {
	return &Container{
		DB: db,
		// mem: mem,
	}
}
