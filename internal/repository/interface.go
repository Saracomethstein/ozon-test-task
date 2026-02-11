package repository

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	SavePost(ctx context.Context, post models.Post) (models.Post, error)
	SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error)
}
