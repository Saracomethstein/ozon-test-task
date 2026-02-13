package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *post) Save(ctx context.Context, post models.Post) (models.Post, error) {
	return models.Post{}, nil
}

func (r *post) GetByID(ctx context.Context, postID int64) (*models.Post, error) { return nil, nil }

func (r *post) Get(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error) {
	return nil, nil
}

func (r *post) TotalCount(ctx context.Context) (int64, error) { return 0, nil }

func (r *post) SetCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error) {
	return nil, nil
}
