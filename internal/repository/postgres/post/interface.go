package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	// TODO: rename func for package 'post'

	SavePost(ctx context.Context, post models.Post) (models.Post, error)
	SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error)
	GetPostById(ctx context.Context, postID int64) (*models.Post, error)
	GetPosts(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error)
	TotalCount(ctx context.Context) (int64, error)
}
