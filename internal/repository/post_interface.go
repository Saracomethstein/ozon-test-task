package repository

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type PostUC interface {
	Save(ctx context.Context, post models.Post) (models.Post, error)
	SetCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error)
	GetByID(ctx context.Context, postID int64) (*models.Post, error)
	Get(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error)
	TotalCount(ctx context.Context) (int64, error)
}
