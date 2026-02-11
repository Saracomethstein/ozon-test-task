package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	CreatePost(ctx context.Context, in models.CreatePostInput) (*models.Post, error)
	SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error)
}
